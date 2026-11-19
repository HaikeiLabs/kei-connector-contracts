// Copyright 2026 Haikei Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mockapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"

	contracts "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/contracts/v1"
)

// Route is one route in the mock's table. Routes mirror the manifest: one
// per operation with an HTTP transport, plus POST /graphql per GraphQL root
// field. SQL operations have no route.
type Route struct {
	// OpID is the operation id, for example github.issues.get.
	OpID string
	// Method is the HTTP method ("POST" for GraphQL).
	Method string
	// Template is the path template ("/repos/{owner}/{repository}/issues/{issue_number}")
	// or "/graphql" for GraphQL operations.
	Template string
	// RootField is the GraphQL document's root field, empty for REST.
	RootField string
}

// route is the compiled form of Route.
type route struct {
	Route
	re           *regexp.Regexp
	connector    *contracts.Connector
	req          protoreflect.MessageDescriptor
	resp         protoreflect.MessageDescriptor
	op           *contracts.Operation
	cursorWire   string
	nextWire     string
	cursorIsInt  bool
	cursorInBody bool
	multiSegment []string
}

// manifestFile mirrors the manifest's envelope; only what the mock needs.
type manifestFile struct {
	Connectors []struct {
		ID           string `json:"id"`
		Capabilities []struct {
			ID       string `json:"id"`
			Request  string `json:"request"`
			Response string `json:"response"`
		} `json:"capabilities"`
	} `json:"connectors"`
}

var rootFieldRe = regexp.MustCompile(`\{\s*([A-Za-z_][A-Za-z0-9_]*)`)

// buildRoutes resolves every manifest capability against the registered
// descriptors and compiles the route table. An operation with an HTTP
// transport that the manifest lacks, or a manifest capability whose messages
// are not linked in, is a contracts drift and fails New.
func buildRoutes(raw []byte) ([]*route, error) {
	var mf manifestFile
	if err := json.Unmarshal(raw, &mf); err != nil {
		return nil, fmt.Errorf("mockapi: parse manifest: %w", err)
	}
	var routes []*route
	inManifest := map[string]bool{}
	for _, mc := range mf.Connectors {
		connector, ok := contracts.LookupConnector(mc.ID)
		if !ok {
			return nil, fmt.Errorf("mockapi: manifest connector %s is not linked into this binary (import its generated package)", mc.ID)
		}
		for _, cap := range mc.Capabilities {
			if inManifest[cap.ID] {
				return nil, fmt.Errorf("mockapi: manifest lists %s twice", cap.ID)
			}
			inManifest[cap.ID] = true
			reqDesc, err := findMessage(cap.Request)
			if err != nil {
				return nil, fmt.Errorf("mockapi: %s: %w", cap.ID, err)
			}
			respDesc, err := findMessage(cap.Response)
			if err != nil {
				return nil, fmt.Errorf("mockapi: %s: %w", cap.ID, err)
			}
			op := operationOf(reqDesc)
			if op == nil {
				return nil, fmt.Errorf("mockapi: %s: request %s carries no operation option", cap.ID, cap.Request)
			}
			if op.GetId() != cap.ID {
				return nil, fmt.Errorf("mockapi: manifest op %s but request %s declares %s", cap.ID, cap.Request, op.GetId())
			}
			r, err := newRoute(op, connector, reqDesc, respDesc)
			if err != nil {
				return nil, fmt.Errorf("mockapi: %s: %w", cap.ID, err)
			}
			if r != nil {
				routes = append(routes, r)
			}
		}
	}
	// Every registered HTTP/GraphQL operation must be in the manifest.
	for _, op := range contracts.All() {
		if op.GetHttpMethod() == "" && op.GetGraphqlDocument() == "" {
			continue
		}
		if !inManifest[op.GetId()] {
			return nil, fmt.Errorf("mockapi: operation %s has an HTTP transport but is not in the manifest: regenerate it", op.GetId())
		}
	}
	return routes, nil
}

// newRoute compiles one operation's route; nil for SQL operations.
func newRoute(op *contracts.Operation, connector *contracts.Connector, req, resp protoreflect.MessageDescriptor) (*route, error) {
	if op.GetHttpMethod() == "" && op.GetGraphqlDocument() == "" {
		return nil, nil // SQL: runs over a DSN, never HTTP.
	}
	r := &route{Route: Route{OpID: op.GetId(), Method: op.GetHttpMethod(), Template: op.GetPathTemplate()}, op: op, connector: connector, req: req, resp: resp}
	if doc := op.GetGraphqlDocument(); doc != "" {
		r.Method = http.MethodPost
		r.Template = "/graphql"
		m := rootFieldRe.FindStringSubmatch(doc)
		if m == nil {
			return nil, fmt.Errorf("graphql document %q has no root field", doc)
		}
		r.RootField = m[1]
	} else if r.Method == "" || r.Template == "" {
		return nil, fmt.Errorf("operation has neither an HTTP method+path nor a GraphQL document")
	}
	// Path variables and their shapes.
	multi := map[string]bool{}
	for _, name := range pathVars(r.Template) {
		f := req.Fields().ByName(protoreflect.Name(name))
		if f == nil {
			return nil, fmt.Errorf("path variable {%s} is not a field of %s", name, req.FullName())
		}
		if p := paramOf(f); p != nil && p.GetMultiSegment() {
			multi[name] = true
			r.multiSegment = append(r.multiSegment, name)
		}
	}
	// Cursor wiring: where the cursor arrives (query for GET/DELETE, body or
	// GraphQL variables otherwise) and where the next cursor goes.
	if cf := op.GetCursorField(); cf != "" {
		f := req.Fields().ByName(protoreflect.Name(cf))
		if f == nil {
			return nil, fmt.Errorf("cursor field %s is not a field of %s", cf, req.FullName())
		}
		r.cursorWire = wireName(f, r.RootField != "")
		switch r.Method {
		case http.MethodGet, http.MethodDelete, http.MethodHead:
			// cursor arrives in the query string
		default:
			r.cursorInBody = true
		}
		nf := resp.Fields().ByName(protoreflect.Name(op.GetNextCursorField()))
		if nf == nil {
			return nil, fmt.Errorf("next cursor field %s is not a field of %s", op.GetNextCursorField(), resp.FullName())
		}
		r.nextWire = wireName(nf, false)
		r.cursorIsInt = nf.Kind() == protoreflect.Int32Kind || nf.Kind() == protoreflect.Int64Kind ||
			nf.Kind() == protoreflect.Uint32Kind || nf.Kind() == protoreflect.Uint64Kind
	}
	// Compile the matcher.
	if r.RootField != "" {
		r.re = regexp.MustCompile(`^/graphql$`)
	} else {
		pattern := ""
		for _, seg := range strings.Split(strings.TrimPrefix(r.Template, "/"), "/") {
			pattern += "/"
			if name, ok := pathVar(seg); ok {
				if multi[name] {
					pattern += `(?P<` + name + `>.+)`
				} else {
					pattern += `(?P<` + name + `>[^/]+)`
				}
			} else {
				pattern += regexp.QuoteMeta(seg)
			}
		}
		re, err := regexp.Compile(`^` + pattern + `$`)
		if err != nil {
			return nil, fmt.Errorf("template %s: %w", r.Template, err)
		}
		r.re = re
	}
	return r, nil
}

// cursorInBody marks cursor arrival in the body (POST/PUT/PATCH) instead of
// the query string. is the name an operation's wire carries the field under: the
// declared override, else the JSON name for GraphQL variables and the proto
// field name for REST (query and body), matching the connection runtime.
func wireName(f protoreflect.FieldDescriptor, graphql bool) string {
	if p := paramOf(f); p != nil && p.GetName() != "" {
		return p.GetName()
	}
	if graphql {
		return f.JSONName()
	}
	return string(f.Name())
}

// paramOf reads the kei.contracts.v1.param field option, nil when unset.
func paramOf(f protoreflect.FieldDescriptor) *contracts.Param {
	opts, ok := f.Options().(*descriptorpb.FieldOptions)
	if !ok || opts == nil || !proto.HasExtension(opts, contracts.E_Param) {
		return nil
	}
	p, _ := proto.GetExtension(opts, contracts.E_Param).(*contracts.Param)
	return p
}

// operationOf reads the kei.contracts.v1.operation message option, nil when
// unset.
func operationOf(md protoreflect.MessageDescriptor) *contracts.Operation {
	opts, ok := md.Options().(*descriptorpb.MessageOptions)
	if !ok || opts == nil || !proto.HasExtension(opts, contracts.E_Operation) {
		return nil
	}
	op, _ := proto.GetExtension(opts, contracts.E_Operation).(*contracts.Operation)
	return op
}

// findMessage resolves a full message name in the linked descriptors.
func findMessage(fullName string) (protoreflect.MessageDescriptor, error) {
	var found protoreflect.MessageDescriptor
	var walk func(md protoreflect.MessageDescriptor)
	walk = func(md protoreflect.MessageDescriptor) {
		if found != nil {
			return
		}
		if string(md.FullName()) == fullName {
			found = md
			return
		}
		for i := 0; i < md.Messages().Len(); i++ {
			walk(md.Messages().Get(i))
		}
	}
	protoregistry.GlobalFiles.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		for i := 0; i < fd.Messages().Len(); i++ {
			walk(fd.Messages().Get(i))
		}
		return found == nil
	})
	if found == nil {
		return nil, fmt.Errorf("message %s is not linked into this binary", fullName)
	}
	return found, nil
}

// pathVars lists the {variables} of a path template, in order.
func pathVars(template string) []string {
	var names []string
	for _, seg := range strings.Split(strings.TrimPrefix(template, "/"), "/") {
		if name, ok := pathVar(seg); ok {
			names = append(names, name)
		}
	}
	return names
}

// pathVar extracts the name of a {variable} segment, ok=false for literals.
func pathVar(seg string) (string, bool) {
	if len(seg) >= 3 && strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}") {
		return seg[1 : len(seg)-1], true
	}
	return "", false
}

// makeRoutes snapshots the public view of the route table.
func makeRoutes(rs []*route) []Route {
	out := make([]Route, len(rs))
	for i, r := range rs {
		out[i] = r.Route
	}
	return out
}
