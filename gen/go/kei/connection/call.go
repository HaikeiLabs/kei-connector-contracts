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

package connection

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	contracts "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/contracts/v1"
)

// ParamIn says where a request field goes.
type ParamIn int

const (
	// InPath fills a path-template variable.
	InPath ParamIn = iota + 1
	// InQuery adds a query parameter; unset fields are omitted.
	InQuery
	// InBody adds a key to the JSON body (or to the GraphQL variables);
	// unset fields are omitted.
	InBody
)

// Param places one request field. Field is the proto field name; Wire is
// the name on the wire. MultiSegment (path only) keeps "/" between escaped
// segments; see RenderPath.
type Param struct {
	Field        string
	In           ParamIn
	Wire         string
	MultiSegment bool
}

// HTTPSpec is everything a generated HTTP operation passes to CallHTTP. The
// generator fills it from the operation option; every placement is explicit.
type HTTPSpec struct {
	Op           Op
	Method       string
	PathTemplate string
	Params       []Param
	Errors       []*contracts.ErrorMapping
	// GraphQLDocument, when set, makes the call a GraphQL request: the body
	// is {"query": document, "variables": {...}} and the response's data
	// is unwrapped (see CallHTTP).
	GraphQLDocument string
	// CursorField is the request field carrying the page cursor, and
	// NextCursorField the response field returning the next one.
	CursorField     string
	NextCursorField string
}

// ErrDecode is wrapped by every response decoding error. Decoding errors
// never include the response body.
var ErrDecode = errors.New("response does not match the contract")

// CallHTTP runs one HTTP operation: it renders the path, places the request
// fields, sends the request, maps a non-2xx status through spec.Errors and
// decodes the JSON response into resp, discarding unknown fields.
//
//   - A JSON array response fills resp's only repeated field (GitHub lists).
//   - GraphQL: an "errors" array maps extensions.code to an error kind
//     (NOT_FOUND, UNAUTHENTICATED, FORBIDDEN, RATELIMITED; anything else is
//     invalid). When data has one key that is not a field of resp, that
//     key's value is decoded; otherwise data itself is.
//   - When NextCursorField is unset after decoding, it is filled from the
//     Link header's rel="next" URL, from the cursor field's query parameter.
func CallHTTP(ctx context.Context, conn HTTP, spec *HTTPSpec, req, resp proto.Message) error {
	if conn == nil {
		return fmt.Errorf("connection: %s: nil connection", spec.Op.ID)
	}
	path, query, body, err := encodeRequest(spec, req.ProtoReflect())
	if err != nil {
		return fmt.Errorf("connection: %s: %w", spec.Op.ID, err)
	}
	status, header, respBody, err := conn.Do(ctx, spec.Method, path, query, body)
	if err != nil {
		return fmt.Errorf("connection: %s: %w", spec.Op.ID, err)
	}
	if err := MapError(spec.Op, status, spec.Errors); err != nil {
		return err
	}
	if spec.GraphQLDocument != "" {
		if err := decodeGraphQL(spec, status, respBody, resp); err != nil {
			return err
		}
	} else if err := decodeJSON(spec, respBody, resp); err != nil {
		return err
	}
	fillNextCursor(spec, header, resp.ProtoReflect())
	return nil
}

func encodeRequest(spec *HTTPSpec, m protoreflect.Message) (string, url.Values, []byte, error) {
	fields := m.Descriptor().Fields()
	pathValues := map[string]string{}
	var multiSegment []string
	query := url.Values{}
	body := object{}
	hasBody := spec.GraphQLDocument != ""
	names := nameFunc(WireName)
	if hasBody {
		names = jsonName
	}
	for _, p := range spec.Params {
		fd := fields.ByName(protoreflect.Name(p.Field))
		if fd == nil {
			return "", nil, nil, fmt.Errorf("spec names unknown request field %q", p.Field)
		}
		switch p.In {
		case InPath:
			if fd.IsList() || fd.IsMap() {
				return "", nil, nil, fmt.Errorf("field %s: a list cannot be a path value", fd.Name())
			}
			s, err := scalarString(fd, m.Get(fd))
			if err != nil {
				return "", nil, nil, err
			}
			pathValues[p.Wire] = s
			if p.MultiSegment {
				multiSegment = append(multiSegment, p.Wire)
			}
		case InQuery:
			if !m.Has(fd) {
				continue
			}
			if fd.IsMap() {
				return "", nil, nil, fmt.Errorf("field %s: a map cannot be a query value", fd.Name())
			}
			if fd.IsList() {
				l := m.Get(fd).List()
				for i := 0; i < l.Len(); i++ {
					s, err := scalarString(fd, l.Get(i))
					if err != nil {
						return "", nil, nil, err
					}
					query.Add(p.Wire, s)
				}
				continue
			}
			s, err := scalarString(fd, m.Get(fd))
			if err != nil {
				return "", nil, nil, err
			}
			query.Set(p.Wire, s)
		case InBody:
			hasBody = true
			if !m.Has(fd) {
				continue
			}
			raw, err := encodeField(fd, m.Get(fd), names)
			if err != nil {
				return "", nil, nil, err
			}
			body[p.Wire] = raw
		default:
			return "", nil, nil, fmt.Errorf("field %s: unsupported placement %d", fd.Name(), p.In)
		}
	}
	path, err := RenderPath(spec.PathTemplate, pathValues, multiSegment...)
	if err != nil {
		return "", nil, nil, err
	}
	if !hasBody {
		return path, query, nil, nil
	}
	if spec.GraphQLDocument != "" {
		return path, query, object{"query": jsonString(spec.GraphQLDocument), "variables": body.encode()}.encode(), nil
	}
	return path, query, body.encode(), nil
}

func decodeError(spec *HTTPSpec, resp proto.Message) error {
	return fmt.Errorf("connection: %s: decoding %s: %w", spec.Op.ID, resp.ProtoReflect().Descriptor().Name(), ErrDecode)
}

func decodeJSON(spec *HTTPSpec, body []byte, resp proto.Message) error {
	body = bytes.TrimSpace(body)
	if len(body) == 0 {
		return nil
	}
	if body[0] == '[' {
		fd := soleRepeatedField(resp.ProtoReflect().Descriptor(), spec.NextCursorField)
		if fd == nil {
			return decodeError(spec, resp)
		}
		body = object{fd.JSONName(): body}.encode()
	}
	if err := (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(body, resp); err != nil {
		return decodeError(spec, resp)
	}
	return nil
}

// soleRepeatedField is md's only list field, not counting the next-cursor
// field, or nil.
func soleRepeatedField(md protoreflect.MessageDescriptor, nextCursor string) protoreflect.FieldDescriptor {
	var found protoreflect.FieldDescriptor
	fields := md.Fields()
	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		if !fd.IsList() || string(fd.Name()) == nextCursor {
			continue
		}
		if found != nil {
			return nil
		}
		found = fd
	}
	return found
}

type graphQLEnvelope struct {
	Data   json.RawMessage `json:"data"`
	Errors []struct {
		Extensions struct {
			Code string `json:"code"`
		} `json:"extensions"`
	} `json:"errors"`
}

func decodeGraphQL(spec *HTTPSpec, status int, body []byte, resp proto.Message) error {
	var env graphQLEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		return decodeError(spec, resp)
	}
	if len(env.Errors) > 0 {
		kind, retryable := graphQLKind(env.Errors[0].Extensions.Code)
		return &Error{Kind: kind, Status: status, Retryable: retryable, OpID: spec.Op.ID}
	}
	data := bytes.TrimSpace(env.Data)
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		return decodeError(spec, resp)
	}
	if len(obj) == 1 {
		fields := resp.ProtoReflect().Descriptor().Fields()
		for k, v := range obj {
			if fields.ByJSONName(k) == nil && fields.ByName(protoreflect.Name(k)) == nil {
				data = v
			}
		}
	}
	return decodeJSON(spec, data, resp)
}

func graphQLKind(code string) (contracts.ErrorKind, bool) {
	switch code {
	case "NOT_FOUND":
		return contracts.ErrorKind_ERROR_KIND_NOT_FOUND, false
	case "UNAUTHENTICATED":
		return contracts.ErrorKind_ERROR_KIND_UNAUTHORIZED, false
	case "FORBIDDEN":
		return contracts.ErrorKind_ERROR_KIND_FORBIDDEN, false
	case "RATELIMITED":
		return contracts.ErrorKind_ERROR_KIND_RATE_LIMITED, true
	default:
		return contracts.ErrorKind_ERROR_KIND_INVALID, false
	}
}

var linkRe = regexp.MustCompile(`<([^>]*)>((?:\s*;\s*[^;,]*)*)`)

// nextLink returns the rel="next" URL of a Link header, or "".
func nextLink(header http.Header) string {
	for _, v := range header.Values("Link") {
		for _, m := range linkRe.FindAllStringSubmatch(v, -1) {
			for _, param := range strings.Split(m[2], ";") {
				k, val, ok := strings.Cut(strings.TrimSpace(param), "=")
				if !ok || !strings.EqualFold(strings.TrimSpace(k), "rel") {
					continue
				}
				for _, rel := range strings.Fields(strings.Trim(strings.TrimSpace(val), `"`)) {
					if strings.EqualFold(rel, "next") {
						return m[1]
					}
				}
			}
		}
	}
	return ""
}

// fillNextCursor sets an unset next-cursor field from the Link header. A
// value that does not fit the field's type is ignored.
func fillNextCursor(spec *HTTPSpec, header http.Header, resp protoreflect.Message) {
	if spec.NextCursorField == "" || spec.CursorField == "" {
		return
	}
	fd := resp.Descriptor().Fields().ByName(protoreflect.Name(spec.NextCursorField))
	if fd == nil || fd.IsList() || fd.IsMap() || resp.Has(fd) {
		return
	}
	link := nextLink(header)
	if link == "" {
		return
	}
	u, err := url.Parse(link)
	if err != nil {
		return
	}
	wire := spec.CursorField
	for _, p := range spec.Params {
		if p.Field == spec.CursorField {
			wire = p.Wire
		}
	}
	raw := u.Query().Get(wire)
	if raw == "" {
		return
	}
	switch fd.Kind() {
	case protoreflect.StringKind:
		resp.Set(fd, protoreflect.ValueOfString(raw))
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		if n, err := strconv.ParseInt(raw, 10, 32); err == nil {
			resp.Set(fd, protoreflect.ValueOfInt32(int32(n)))
		}
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		if n, err := strconv.ParseInt(raw, 10, 64); err == nil {
			resp.Set(fd, protoreflect.ValueOfInt64(n))
		}
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		if n, err := strconv.ParseUint(raw, 10, 32); err == nil {
			resp.Set(fd, protoreflect.ValueOfUint32(uint32(n)))
		}
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		if n, err := strconv.ParseUint(raw, 10, 64); err == nil {
			resp.Set(fd, protoreflect.ValueOfUint64(n))
		}
	}
}
