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

package contracts_test

import (
	"regexp"
	"slices"
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"

	databasev1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/database/v1"
	httpapiv1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/http_api/v1"
	linearv1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/linear/v1"
	contracts "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/contracts/v1"
)

// Importing every connector package puts its files, and so its connector
// and operations, into protoregistry.GlobalFiles.
var (
	_ = databasev1.RunQueryRequest{}
	_ = httpapiv1.ListRecordsRequest{}
	_ = linearv1.GetIssueRequest{}
)

// annotatedOp pairs an operation with the request message that declares it,
// so tests can check the option against the message's real fields.
type annotatedOp struct {
	op  *contracts.Operation
	req protoreflect.MessageDescriptor
}

func annotatedOps(t *testing.T) []annotatedOp {
	t.Helper()
	var out []annotatedOp
	protoregistry.GlobalFiles.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		msgs := fd.Messages()
		for i := 0; i < msgs.Len(); i++ {
			md := msgs.Get(i)
			opts, ok := md.Options().(*descriptorpb.MessageOptions)
			if !ok || !proto.HasExtension(opts, contracts.E_Operation) {
				continue
			}
			out = append(out, annotatedOp{proto.GetExtension(opts, contracts.E_Operation).(*contracts.Operation), md})
		}
		return true
	})
	if len(out) == 0 {
		t.Fatal("no annotated operations registered")
	}
	return out
}

var templateVar = regexp.MustCompile(`\{([a-z0-9_]+)\}`)

func templateVars(tmpl string) []string {
	var vs []string
	for _, m := range templateVar.FindAllStringSubmatch(tmpl, -1) {
		vs = append(vs, m[1])
	}
	return vs
}

// The pilot connectors and their connection names (plan Task 2).
func TestPilotConnectorsDeclareConnections(t *testing.T) {
	want := map[string][]string{
		"database": {"dsn", "snowflake_key_pair"},
		"github":   {"user_oauth", "pat"},
		"http_api": {"api_key"},
		"linear":   {"user_oauth", "api_key"},
	}
	for id, names := range want {
		c, ok := contracts.LookupConnector(id)
		if !ok {
			t.Errorf("connector %q is not declared", id)
			continue
		}
		var got []string
		for _, conn := range c.GetConnections() {
			got = append(got, conn.GetName())
			if conn.GetKind() == contracts.ConnectionKind_CONNECTION_KIND_UNSPECIFIED {
				t.Errorf("%s/%s: connection kind unspecified", id, conn.GetName())
			}
			if conn.GetAccountModel() == contracts.AccountModel_ACCOUNT_MODEL_UNSPECIFIED {
				t.Errorf("%s/%s: account model unspecified", id, conn.GetName())
			}
			if conn.GetAuthMethod() == "" {
				t.Errorf("%s/%s: no auth_method", id, conn.GetName())
			}
			if conn.GetKind() == contracts.ConnectionKind_CONNECTION_KIND_HTTP_HEADER && conn.GetHeaderName() == "" {
				t.Errorf("%s/%s: HTTP_HEADER connection names no header", id, conn.GetName())
			}
			// A shared connection is configured by an admin: it must say
			// which setup field holds its secret.
			if conn.GetAccountModel() == contracts.AccountModel_ACCOUNT_MODEL_SHARED &&
				!slices.ContainsFunc(conn.GetSetupFields(), (*contracts.SetupField).GetSecret) {
				t.Errorf("%s/%s: shared connection has no secret setup field", id, conn.GetName())
			}
		}
		if !slices.Equal(got, names) {
			t.Errorf("connector %q connections = %v, want %v", id, got, names)
		}
	}
}

// Every operation names only connections its connector declares, and at
// least one: an operation with no connection can never be called.
func TestEveryOperationNamesDeclaredConnections(t *testing.T) {
	for _, op := range contracts.All() {
		c, ok := contracts.LookupConnector(op.GetConnector())
		if !ok {
			t.Errorf("%s: connector %q is not declared", op.GetId(), op.GetConnector())
			continue
		}
		if len(op.GetConnections()) == 0 {
			t.Errorf("%s: names no connections", op.GetId())
		}
		for _, name := range op.GetConnections() {
			if !slices.ContainsFunc(c.GetConnections(), func(conn *contracts.Connection) bool { return conn.GetName() == name }) {
				t.Errorf("%s: connection %q is not declared by connector %q", op.GetId(), name, c.GetId())
			}
		}
	}
}

func TestGetIssueAnnotated(t *testing.T) {
	op, ok := contracts.Lookup("github.issues.get")
	if !ok {
		t.Fatal("github.issues.get not found")
	}
	if op.GetMethodName() != "GetIssue" {
		t.Errorf("method_name = %q, want GetIssue", op.GetMethodName())
	}
	if op.GetResourceTemplate() != "repos/{owner}/{repository}" {
		t.Errorf("resource_template = %q", op.GetResourceTemplate())
	}
	if !slices.Equal(op.GetConnections(), []string{"user_oauth", "pat"}) {
		t.Errorf("connections = %v, want [user_oauth pat]", op.GetConnections())
	}
}

// Every path-template variable is a request field carrying param.in = PATH,
// and every PATH field appears in the template: the generator renders the
// path from exactly these fields.
func TestPathVariablesArePathParams(t *testing.T) {
	for _, a := range annotatedOps(t) {
		vars := templateVars(a.op.GetPathTemplate())
		fields := a.req.Fields()
		for _, v := range vars {
			fd := fields.ByName(protoreflect.Name(v))
			if fd == nil {
				t.Errorf("%s: path variable {%s} has no request field", a.op.GetId(), v)
				continue
			}
			if paramIn(fd) != contracts.ParamIn_PARAM_IN_PATH {
				t.Errorf("%s: field %s is a path variable but param.in = %v", a.op.GetId(), v, paramIn(fd))
			}
		}
		for i := 0; i < fields.Len(); i++ {
			fd := fields.Get(i)
			if paramIn(fd) == contracts.ParamIn_PARAM_IN_PATH && !slices.Contains(vars, string(fd.Name())) {
				t.Errorf("%s: field %s has param.in = PATH but is not in %q", a.op.GetId(), fd.Name(), a.op.GetPathTemplate())
			}
		}
	}
}

func paramIn(fd protoreflect.FieldDescriptor) contracts.ParamIn {
	opts, ok := fd.Options().(*descriptorpb.FieldOptions)
	if !ok || !proto.HasExtension(opts, contracts.E_Param) {
		return contracts.ParamIn_PARAM_IN_UNSPECIFIED
	}
	return proto.GetExtension(opts, contracts.E_Param).(*contracts.Param).GetIn()
}

// The resource template is what policy decides on: each variable must be a
// request field, or the runner cannot render the resource.
func TestResourceTemplatesRenderFromRequest(t *testing.T) {
	for _, a := range annotatedOps(t) {
		rt := a.op.GetResourceTemplate()
		if rt == "" {
			t.Errorf("%s: no resource_template", a.op.GetId())
			continue
		}
		if strings.HasPrefix(rt, "/") {
			t.Errorf("%s: resource_template %q must be relative", a.op.GetId(), rt)
		}
		for _, v := range templateVars(rt) {
			if a.req.Fields().ByName(protoreflect.Name(v)) == nil {
				t.Errorf("%s: resource variable {%s} has no request field", a.op.GetId(), v)
			}
		}
	}
}

// method_name is the generated function name; it must be a Go-exported
// identifier and unique within a connector.
func TestMethodNamesUniquePerConnector(t *testing.T) {
	ident := regexp.MustCompile(`^[A-Z][A-Za-z0-9]*$`)
	seen := map[string]string{}
	for _, op := range contracts.All() {
		name := op.GetMethodName()
		if !ident.MatchString(name) {
			t.Errorf("%s: method_name %q is not an exported identifier", op.GetId(), name)
			continue
		}
		key := op.GetConnector() + "." + name
		if prev, dup := seen[key]; dup {
			t.Errorf("%s and %s share method_name %q", prev, op.GetId(), name)
		}
		seen[key] = op.GetId()
	}
}

// Semantics are required on every operation and must agree with its class:
// a READ operation maps only to the read action and a WRITE operation never
// does. A disagreement would let a semantic "read issues" rule cover a write.
func TestSemanticsAgreeWithClass(t *testing.T) {
	for _, op := range contracts.All() {
		s := op.GetSemantics()
		if s == nil || s.GetAction() == contracts.SemanticAction_SEMANTIC_ACTION_UNSPECIFIED {
			t.Errorf("%s: no semantic action", op.GetId())
			continue
		}
		if s.GetResourceType() == "" || s.GetVersion() == 0 {
			t.Errorf("%s: semantics needs resource_type and version >= 1, got %v", op.GetId(), s)
		}
		isRead := s.GetAction() == contracts.SemanticAction_SEMANTIC_ACTION_READ
		if isRead != (op.GetClass() == contracts.OperationClass_OPERATION_CLASS_READ) {
			t.Errorf("%s: semantic action %v disagrees with class %v", op.GetId(), s.GetAction(), op.GetClass())
		}
	}
}

// Every HTTP operation maps the standard provider statuses; only a rate
// limit is retryable. Operations without an HTTP method (SQL) are exempt.
func TestHTTPOperationsMapStandardErrors(t *testing.T) {
	want := map[int32]struct {
		kind      contracts.ErrorKind
		retryable bool
	}{
		401: {contracts.ErrorKind_ERROR_KIND_UNAUTHORIZED, false},
		403: {contracts.ErrorKind_ERROR_KIND_FORBIDDEN, false},
		404: {contracts.ErrorKind_ERROR_KIND_NOT_FOUND, false},
		429: {contracts.ErrorKind_ERROR_KIND_RATE_LIMITED, true},
	}
	for _, op := range contracts.All() {
		if op.GetHttpMethod() == "" {
			continue
		}
		got := map[int32]*contracts.ErrorMapping{}
		for _, e := range op.GetErrors() {
			if _, dup := got[e.GetStatus()]; dup {
				t.Errorf("%s: status %d mapped twice", op.GetId(), e.GetStatus())
			}
			got[e.GetStatus()] = e
		}
		for status, w := range want {
			e, ok := got[status]
			if !ok {
				t.Errorf("%s: status %d not mapped", op.GetId(), status)
				continue
			}
			if e.GetKind() != w.kind || e.GetRetryable() != w.retryable {
				t.Errorf("%s: status %d = %v retryable=%v, want %v retryable=%v",
					op.GetId(), status, e.GetKind(), e.GetRetryable(), w.kind, w.retryable)
			}
		}
	}
}

// The database operation runs over a SQL connection, not HTTP. Every
// connection reaches a SQL engine: a SQL DSN pool, or a key-pair JWT
// connection that names the SQL engine it reaches (Snowflake).
func TestDatabaseQueryHasNoHTTPTransport(t *testing.T) {
	op, ok := contracts.Lookup("database.query.run")
	if !ok {
		t.Fatal("database.query.run not found")
	}
	if op.GetHttpMethod() != "" || op.GetPathTemplate() != "" {
		t.Errorf("database.query.run has HTTP transport %q %q; want none", op.GetHttpMethod(), op.GetPathTemplate())
	}
	c, _ := contracts.LookupConnector("database")
	for _, conn := range c.GetConnections() {
		switch conn.GetKind() {
		case contracts.ConnectionKind_CONNECTION_KIND_SQL_DSN:
		case contracts.ConnectionKind_CONNECTION_KIND_KEY_PAIR_JWT:
			if conn.GetEngine() == "" {
				t.Errorf("database connection %q kind = KEY_PAIR_JWT with no engine", conn.GetName())
			}
		default:
			t.Errorf("database connection %q kind = %v, want SQL_DSN or KEY_PAIR_JWT with an engine", conn.GetName(), conn.GetKind())
		}
	}
}

// multi_segment is only meaningful on a path param, and the GitHub contents
// path (a file path such as docs/a.md) needs it: escaping its slashes would
// address a different file.
func TestMultiSegmentPathParams(t *testing.T) {
	for _, a := range annotatedOps(t) {
		fields := a.req.Fields()
		for i := 0; i < fields.Len(); i++ {
			fd := fields.Get(i)
			p := paramOf(fd)
			if p.GetMultiSegment() && p.GetIn() != contracts.ParamIn_PARAM_IN_PATH {
				t.Errorf("%s: field %s is multi_segment but not a path param", a.op.GetId(), fd.Name())
			}
		}
	}
	for _, c := range []struct{ op, field string }{{"github.contents.put", "path"}} {
		var found bool
		for _, a := range annotatedOps(t) {
			if a.op.GetId() != c.op {
				continue
			}
			found = true
			if fd := a.req.Fields().ByName(protoreflect.Name(c.field)); fd == nil || !paramOf(fd).GetMultiSegment() {
				t.Errorf("%s: field %s must be multi_segment", c.op, c.field)
			}
		}
		if !found {
			t.Errorf("%s not found", c.op)
		}
	}
}

func paramOf(fd protoreflect.FieldDescriptor) *contracts.Param {
	opts, ok := fd.Options().(*descriptorpb.FieldOptions)
	if !ok || !proto.HasExtension(opts, contracts.E_Param) {
		return nil
	}
	return proto.GetExtension(opts, contracts.E_Param).(*contracts.Param)
}

var (
	graphqlOperation = regexp.MustCompile(`^(query|mutation)\s+[A-Za-z_][A-Za-z0-9_]*\s*(\(|\{)`)
	graphqlVariable  = regexp.MustCompile(`\$([A-Za-z_][A-Za-z0-9_]*)`)
)

// Decision 4: every Linear operation sends one fixed document. A read sends
// a query and a write a mutation, so a document can never write under a
// read grant. Every variable is a request field (by its JSON name, which is
// what the request serializes to), and every request field is a variable,
// so no request value is silently dropped.
func TestLinearOperationsHaveGraphQLDocument(t *testing.T) {
	var n int
	for _, a := range annotatedOps(t) {
		if a.op.GetConnector() != "linear" {
			if a.op.GetGraphqlDocument() != "" {
				t.Errorf("%s: graphql_document on a non-GraphQL connector", a.op.GetId())
			}
			continue
		}
		n++
		doc := strings.TrimSpace(a.op.GetGraphqlDocument())
		m := graphqlOperation.FindStringSubmatch(doc)
		if m == nil {
			t.Errorf("%s: graphql_document %q does not start with a named query or mutation", a.op.GetId(), doc)
			continue
		}
		wantKind := "query"
		if a.op.GetClass() == contracts.OperationClass_OPERATION_CLASS_WRITE {
			wantKind = "mutation"
		}
		if m[1] != wantKind {
			t.Errorf("%s: %s document for class %v, want %s", a.op.GetId(), m[1], a.op.GetClass(), wantKind)
		}
		if strings.Count(doc, "{") != strings.Count(doc, "}") || strings.Count(doc, "(") != strings.Count(doc, ")") {
			t.Errorf("%s: unbalanced braces or parentheses", a.op.GetId())
		}
		vars := map[string]bool{}
		for _, v := range graphqlVariable.FindAllStringSubmatch(doc, -1) {
			vars[v[1]] = true
		}
		jsonNames := map[string]bool{}
		fields := a.req.Fields()
		for i := 0; i < fields.Len(); i++ {
			jsonNames[fields.Get(i).JSONName()] = true
			if !vars[fields.Get(i).JSONName()] {
				t.Errorf("%s: request field %s is not a document variable", a.op.GetId(), fields.Get(i).JSONName())
			}
		}
		for v := range vars {
			if !jsonNames[v] {
				t.Errorf("%s: variable $%s is not a request field", a.op.GetId(), v)
			}
		}
	}
	if n == 0 {
		t.Fatal("no Linear operations registered")
	}
}
