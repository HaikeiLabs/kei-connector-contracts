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

package render

import (
	"bytes"
	_ "embed"
	"fmt"
	"path"
	"sort"
	"strings"
	"text/template"
	"unicode"

	contracts "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/contracts/v1"
	"github.com/HaikeiLabs/kei-connector-contracts/tools/protoc-gen-kei/internal/model"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

//go:embed rust.rs.tmpl
var rustTemplateText string
var rustTemplate = template.Must(template.New("rust").Parse(rustTemplateText))

type rustOperation struct {
	// ConstBlock is the pub const OP_X: Op block (with its doc comment).
	ConstBlock string
	// ErrorsBlock is the per-operation error table, empty when the
	// operation declares no error mappings.
	ErrorsBlock string
	// ResponseFieldsBlock is the GraphQL response's field-name table.
	ResponseFieldsBlock string
	Method, Doc         string
	// Deprecated is the #[deprecated] attribute line, empty when the
	// operation is not deprecated.
	Deprecated string
	Req, Resp  string
	Body, Call string
}

type rustFile struct {
	Source, ID, Name, Uses string
	// Consts is the joined OP_X descriptor blocks, blank-line separated.
	Consts string
	// ConstTables is the joined error/response-field tables, blank-line
	// separated.
	ConstTables string
	SQL         bool
	Ops         []rustOperation
}

// rustPascal is a Rust type name for a connector id: testconn → Testconn,
// http_api → HttpApi.
func rustPascal(s string) string {
	var out strings.Builder
	for _, p := range strings.Split(s, "_") {
		if p != "" {
			out.WriteString(strings.ToUpper(p[:1]) + p[1:])
		}
	}
	return out.String()
}

// rustSnake is SCREAMING_SNAKE_CASE: GetIssue → GET_ISSUE.
func rustSnake(s string) string {
	var out strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			out.WriteByte('_')
		}
		out.WriteRune(unicode.ToUpper(r))
	}
	return out.String()
}

// rustKindPascal is a Rust error-kind variant: ERROR_KIND_RATE_LIMITED →
// RateLimited.
func rustKindPascal(s string) string {
	var out strings.Builder
	for _, p := range strings.Split(s, "_") {
		if p != "" {
			out.WriteString(strings.ToUpper(p[:1]) + strings.ToLower(p[1:]))
		}
	}
	return out.String()
}

// rustIntType is the Rust integer type of a proto integer kind.
func rustIntType(k protoreflect.Kind) string {
	switch k {
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return "i32"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return "i64"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "u32"
	default:
		return "u64"
	}
}

func isIntKind(k protoreflect.Kind) bool {
	switch k {
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
		protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind,
		protoreflect.Uint32Kind, protoreflect.Fixed32Kind,
		protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return true
	}
	return false
}

// rustQuote renders a Rust string literal (ASCII text only).
func rustQuote(s string) string {
	var out strings.Builder
	out.WriteByte('"')
	for _, r := range s {
		switch r {
		case '"':
			out.WriteString("\\\"")
		case '\\':
			out.WriteString("\\\\")
		case '\n':
			out.WriteString("\\n")
		case '\r':
			out.WriteString("\\r")
		case '\t':
			out.WriteString("\\t")
		default:
			out.WriteRune(r)
		}
	}
	out.WriteByte('"')
	return out.String()
}

func renderRust(gen *protogen.Plugin, c *model.Connector, opts Options) error {
	dir := path.Dir(c.File.Desc.Path())
	d := rustFile{Source: c.File.Desc.Path(), ID: c.ID, Name: rustPascal(c.ID), SQL: c.IsSQL()}
	var hasGraphQL, hasHTTP, hasCursor, hasBytes, hasErrors bool
	for _, op := range c.Ops {
		o, err := rustOp(op)
		if err != nil {
			return err
		}
		if op.Transport == model.GraphQL {
			hasGraphQL = true
		} else {
			hasHTTP = true
		}
		if op.Opt.GetCursorField() != "" {
			hasCursor = true
		}
		if len(op.Opt.GetErrors()) > 0 {
			hasErrors = true
		}
		for _, p := range op.Params {
			if p.Field.Desc.Kind() == protoreflect.BytesKind {
				hasBytes = true
			}
		}
		d.Ops = append(d.Ops, o)
	}
	var conn []string
	if d.SQL {
		conn = []string{"Class", "Op", "OperationError", "Sql", "call_sql"}
	} else {
		conn = []string{"Class", "HttpConnection", "Op", "OperationError", "map_error", "render_path"}
		if hasHTTP {
			conn = append(conn, "decode_response")
		}
		if hasGraphQL {
			conn = append(conn, "decode_graphql_response")
		}
		if hasCursor {
			conn = append(conn, "next_cursor_from_link")
		}
		if hasBytes {
			conn = append(conn, "encode_bytes")
		}
		sort.Strings(conn)
	}
	// Message types, grouped by their defining proto package (prost lays
	// a package's messages flat in one module).
	msgs := map[string]map[string]bool{}
	for _, op := range c.Ops {
		for _, m := range []*protogen.Message{op.Request, op.Response} {
			p := strings.ReplaceAll(string(m.Desc.ParentFile().Package()), ".", "/")
			if msgs[p] == nil {
				msgs[p] = map[string]bool{}
			}
			msgs[p][m.GoIdent.GoName] = true
		}
	}
	paths := make([]string, 0, len(msgs))
	for p := range msgs {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	var uses strings.Builder
	fmt.Fprintf(&uses, "use crate::connection::{%s};\n", strings.Join(conn, ", "))
	if hasErrors {
		fmt.Fprintf(&uses, "use crate::kei::contracts::v1::{ErrorKind, ErrorMapping};\n")
	}
	for _, p := range paths {
		names := make([]string, 0, len(msgs[p]))
		for n := range msgs[p] {
			names = append(names, n)
		}
		sort.Strings(names)
		fmt.Fprintf(&uses, "use crate::%s::{%s};\n", strings.ReplaceAll(p, "/", "::"), strings.Join(names, ", "))
	}
	d.Uses = uses.String()
	var constBlocks, constTables []string
	for _, o := range d.Ops {
		constBlocks = append(constBlocks, o.ConstBlock)
		if o.ErrorsBlock != "" {
			constTables = append(constTables, o.ErrorsBlock)
		}
		if o.ResponseFieldsBlock != "" {
			constTables = append(constTables, o.ResponseFieldsBlock)
		}
	}
	d.Consts = strings.Join(constBlocks, "\n")
	d.ConstTables = strings.Join(constTables, "\n")
	var b bytes.Buffer
	if err := rustTemplate.Execute(&b, d); err != nil {
		return fmt.Errorf("protoc-gen-kei-ops: %s: %w", c.ID, err)
	}
	g := gen.NewGeneratedFile(dir+"/ops.rs", "")
	_, err := g.Write(b.Bytes())
	return err
}

// rustOp builds one operation's template data.
func rustOp(op *model.Operation) (rustOperation, error) {
	o := rustOperation{
		Method: strings.ToLower(rustSnake(op.MethodName)),
		Req:    op.Request.GoIdent.GoName,
		Resp:   op.Response.GoIdent.GoName,
	}
	route := strings.ToUpper(op.Opt.GetHttpMethod()) + " " + op.Opt.GetPathTemplate()
	switch op.Transport {
	case model.GraphQL:
		route = "GraphQL " + route
	case model.SQL:
		route = "SQL query"
	}
	o.Doc = fmt.Sprintf("%s runs %s (%s).", o.Method, op.ID, route)
	if op.Opt.GetDeprecated() {
		note := "this operation is deprecated."
		if r := op.Opt.GetReplacedBy(); r != "" {
			note = "use " + r + " instead."
		}
		o.Deprecated = "    #[deprecated(note = " + rustQuote(note) + ")]\n"
	}
	class := "Read"
	if op.Opt.GetClass() == contracts.OperationClass_OPERATION_CLASS_WRITE {
		class = "Write"
	}
	o.ConstBlock = fmt.Sprintf(`/// Operation descriptor for %s, for authorization and audit: the
/// capability a policy grants is connector:%s.
pub const %s: Op = Op {
    id: %s,
    class: Class::%s,
    destructive: %t,
};
`, op.ID, op.ID, "OP_"+rustSnake(op.MethodName), rustQuote(op.ID), class, op.Opt.GetDestructive())
	if len(op.Opt.GetErrors()) > 0 {
		var entries []string
		for _, e := range op.Opt.GetErrors() {
			entries = append(entries, fmt.Sprintf("    ErrorMapping {\n        status: %d,\n        kind: ErrorKind::%s as i32,\n        retryable: %t,\n    },",
				e.GetStatus(), rustKindPascal(strings.TrimPrefix(e.GetKind().String(), "ERROR_KIND_")), e.GetRetryable()))
		}
		o.ErrorsBlock = fmt.Sprintf("const ERRORS_%s: &[ErrorMapping] = &[\n%s\n];\n", rustSnake(op.MethodName), strings.Join(entries, "\n"))
	}
	if op.Transport == model.GraphQL {
		var pairs []string
		for _, f := range op.Response.Fields {
			pairs = append(pairs, fmt.Sprintf("(%s, %s)", rustQuote(string(f.Desc.Name())), rustQuote(f.Desc.JSONName())))
		}
		o.ResponseFieldsBlock = fmt.Sprintf("const RESPONSE_FIELDS_%s: &[(&str, &str)] = &[%s];\n", rustSnake(op.MethodName), strings.Join(pairs, ", "))
	}
	if op.Transport == model.SQL {
		// SQL requests are RunQueryRequest-shaped (sql, params, max_rows);
		// the model's SQL transport has no other shape today.
		o.Call = fmt.Sprintf("call_sql(&self.conn, &%s, &req.sql, &req.params, req.max_rows.unwrap_or(0)).await", "OP_"+rustSnake(op.MethodName))
		return o, nil
	}
	body, err := rustFnBody(op)
	if err != nil {
		return o, err
	}
	o.Body = body
	return o, nil
}

// rustFnBody renders one generated operation's async-fn body: it places
// every request field (path, query or JSON body), calls the connection and
// maps the errors, exactly like the Go runtime's callHTTP.
func rustFnBody(op *model.Operation) (string, error) {
	opConst := "OP_" + rustSnake(op.MethodName)
	var pathParams, multi, queryParams, bodyParams []string
	for _, p := range op.Params {
		f := p.Field
		name := string(f.Desc.Name())
		optional := f.Desc.HasPresence() && !f.Desc.IsList()
		switch p.In {
		case contracts.ParamIn_PARAM_IN_PATH:
			if p.MultiSegment {
				multi = append(multi, rustQuote(p.Wire))
			}
			switch {
			case f.Desc.IsList():
				return "", fmt.Errorf("protoc-gen-kei-ops: %s: field %s: a list cannot be a path value", op.ID, name)
			case optional && f.Desc.Kind() == protoreflect.StringKind:
				pathParams = append(pathParams, fmt.Sprintf("if let Some(v) = &req.%s { path_values.push((%s, v.as_str())); }", name, rustQuote(p.Wire)))
			case optional:
				pathParams = append(pathParams, fmt.Sprintf("let %s_s = req.%s.map(|v| v.to_string());", name, name))
				pathParams = append(pathParams, fmt.Sprintf("if let Some(v) = &%s_s { path_values.push((%s, v.as_str())); }", name, rustQuote(p.Wire)))
			case f.Desc.Kind() == protoreflect.StringKind:
				pathParams = append(pathParams, fmt.Sprintf("path_values.push((%s, req.%s.as_str()));", rustQuote(p.Wire), name))
			case isIntKind(f.Desc.Kind()):
				pathParams = append(pathParams, fmt.Sprintf("let %s_s = req.%s.to_string();", name, name))
				pathParams = append(pathParams, fmt.Sprintf("path_values.push((%s, %s_s.as_str()));", rustQuote(p.Wire), name))
			default:
				return "", fmt.Errorf("protoc-gen-kei-ops: %s: field %s: only string and integer path params are supported by the Rust template yet", op.ID, name)
			}
		case contracts.ParamIn_PARAM_IN_QUERY:
			switch {
			case f.Desc.IsList():
				item := "v.clone()"
				if f.Desc.Kind() != protoreflect.StringKind {
					item = "v.to_string()"
				}
				queryParams = append(queryParams, fmt.Sprintf("for v in &req.%s { query_values.push((%s.to_string(), %s)); }", name, rustQuote(p.Wire), item))
			case optional && f.Desc.Kind() == protoreflect.StringKind:
				queryParams = append(queryParams, fmt.Sprintf("if let Some(v) = &req.%s { query_values.push((%s.to_string(), v.clone())); }", name, rustQuote(p.Wire)))
			case optional && isIntKind(f.Desc.Kind()):
				queryParams = append(queryParams, fmt.Sprintf("if let Some(v) = req.%s { query_values.push((%s.to_string(), v.to_string())); }", name, rustQuote(p.Wire)))
			case optional:
				queryParams = append(queryParams, fmt.Sprintf("if let Some(v) = req.%s { query_values.push((%s.to_string(), v.to_string())); }", name, rustQuote(p.Wire)))
			case f.Desc.Kind() == protoreflect.StringKind:
				queryParams = append(queryParams, fmt.Sprintf("query_values.push((%s.to_string(), req.%s.clone()));", rustQuote(p.Wire), name))
			case isIntKind(f.Desc.Kind()):
				queryParams = append(queryParams, fmt.Sprintf("query_values.push((%s.to_string(), req.%s.to_string()));", rustQuote(p.Wire), name))
			case f.Desc.Kind() == protoreflect.BoolKind:
				queryParams = append(queryParams, fmt.Sprintf("if req.%s { query_values.push((%s.to_string(), \"true\".to_string())); }", name, rustQuote(p.Wire)))
			default:
				return "", fmt.Errorf("protoc-gen-kei-ops: %s: field %s: only string, integer and bool query params are supported by the Rust template yet", op.ID, name)
			}
		default: // body
			switch {
			case !f.Desc.IsMap() && f.Desc.IsList() && (f.Desc.Kind() == protoreflect.MessageKind || f.Desc.Kind() == protoreflect.GroupKind):
				// Repeated message: a JSON array of the elements' protojson
				// shape; pbjson-types serializes Struct/ListValue directly.
				bodyParams = append(bodyParams, fmt.Sprintf("let mut %s_json = Vec::with_capacity(req.%s.len());", name, name))
				bodyParams = append(bodyParams, fmt.Sprintf("for v in &req.%s { %s_json.push(serde_json::to_value(v).map_err(|e| OperationError::wrap(%s.id, e.to_string()))?); }", name, name, opConst))
				bodyParams = append(bodyParams, fmt.Sprintf("body.insert(%s.to_string(), serde_json::Value::Array(%s_json));", rustQuote(p.Wire), name))
			case !f.Desc.IsMap() && !f.Desc.IsList() && (f.Desc.Kind() == protoreflect.MessageKind || f.Desc.Kind() == protoreflect.GroupKind):
				// Singular message (prost models it as Option<T>): the whole
				// object in its protojson shape.
				bodyParams = append(bodyParams, fmt.Sprintf("if let Some(v) = &req.%s { body.insert(%s.to_string(), serde_json::to_value(v).map_err(|e| OperationError::wrap(%s.id, e.to_string()))?); }", name, rustQuote(p.Wire), opConst))
			case f.Desc.IsList():
				if f.Desc.Kind() != protoreflect.StringKind {
					return "", fmt.Errorf("protoc-gen-kei-ops: %s: field %s: only string list body params are supported by the Rust template yet", op.ID, name)
				}
				bodyParams = append(bodyParams, fmt.Sprintf("body.insert(%s.to_string(), serde_json::Value::Array(req.%s.iter().map(|v| serde_json::Value::String(v.clone())).collect()));", rustQuote(p.Wire), name))
			case optional && f.Desc.Kind() == protoreflect.StringKind:
				bodyParams = append(bodyParams, fmt.Sprintf("if let Some(v) = &req.%s { body.insert(%s.to_string(), serde_json::Value::String(v.clone())); }", name, rustQuote(p.Wire)))
			case optional && isIntKind(f.Desc.Kind()):
				bodyParams = append(bodyParams, fmt.Sprintf("if let Some(v) = req.%s { body.insert(%s.to_string(), serde_json::Value::Number(v.into())); }", name, rustQuote(p.Wire)))
			case optional && f.Desc.Kind() == protoreflect.BoolKind:
				bodyParams = append(bodyParams, fmt.Sprintf("if let Some(v) = req.%s { body.insert(%s.to_string(), serde_json::Value::Bool(v)); }", name, rustQuote(p.Wire)))
			case optional && f.Desc.Kind() == protoreflect.BytesKind:
				bodyParams = append(bodyParams, fmt.Sprintf("if let Some(v) = &req.%s { body.insert(%s.to_string(), serde_json::Value::String(encode_bytes(v))); }", name, rustQuote(p.Wire)))
			case f.Desc.Kind() == protoreflect.StringKind:
				bodyParams = append(bodyParams, fmt.Sprintf("body.insert(%s.to_string(), serde_json::Value::String(req.%s.clone()));", rustQuote(p.Wire), name))
			case isIntKind(f.Desc.Kind()):
				bodyParams = append(bodyParams, fmt.Sprintf("body.insert(%s.to_string(), serde_json::Value::Number(req.%s.into()));", rustQuote(p.Wire), name))
			case f.Desc.Kind() == protoreflect.BoolKind:
				bodyParams = append(bodyParams, fmt.Sprintf("if req.%s { body.insert(%s.to_string(), serde_json::Value::Bool(true)); }", name, rustQuote(p.Wire)))
			case f.Desc.Kind() == protoreflect.BytesKind:
				bodyParams = append(bodyParams, fmt.Sprintf("body.insert(%s.to_string(), serde_json::Value::String(encode_bytes(&req.%s)));", rustQuote(p.Wire), name))
			default:
				return "", fmt.Errorf("protoc-gen-kei-ops: %s: field %s: only string, integer, bool and bytes body params are supported by the Rust template yet", op.ID, name)
			}
		}
	}
	var lines []string
	valuesArg := "&[]"
	if len(pathParams) > 0 {
		lines = append(lines, "let mut path_values: Vec<(&str, &str)> = Vec::new();")
		lines = append(lines, pathParams...)
		valuesArg = "&path_values"
	}
	lines = append(lines, fmt.Sprintf(`let path = render_path(
    %s,
    %s,
    &[%s],
)
.map_err(|e| OperationError::wrap(%s.id, format!("render path: {e}")))?;`,
		rustQuote(op.Opt.GetPathTemplate()), valuesArg, strings.Join(multi, ", "), opConst))
	queryArg := "&[]"
	if len(queryParams) > 0 {
		lines = append(lines, "let mut query_values: Vec<(String, String)> = Vec::new();")
		lines = append(lines, queryParams...)
		queryArg = "&query_values"
	}
	bodyArg := "None"
	if len(bodyParams) > 0 || op.Transport == model.GraphQL {
		lines = append(lines, "let mut body = serde_json::Map::new();")
		lines = append(lines, bodyParams...)
		if op.Transport == model.GraphQL {
			lines = append(lines, fmt.Sprintf("let payload = serde_json::json!({ \"query\": %s, \"variables\": serde_json::Value::Object(body) });", rustQuote(op.Opt.GetGraphqlDocument())))
			lines = append(lines, fmt.Sprintf("let body_bytes = serde_json::to_vec(&payload).map_err(|e| OperationError::wrap(%s.id, e.to_string()))?;", opConst))
		} else {
			lines = append(lines, fmt.Sprintf("let body_bytes = serde_json::to_vec(&body).map_err(|e| OperationError::wrap(%s.id, e.to_string()))?;", opConst))
		}
		bodyArg = "Some(body_bytes.as_slice())"
	}
	lines = append(lines, fmt.Sprintf(`let out = self
    .conn
    .do_request(%s, &path, %s, %s)
    .await
    .map_err(|e| OperationError::wrap(%s.id, e.to_string()))?;`,
		rustQuote(strings.ToUpper(op.Opt.GetHttpMethod())), queryArg, bodyArg, opConst))
	errArg := "&[]"
	if len(op.Opt.GetErrors()) > 0 {
		errArg = "&ERRORS_" + rustSnake(op.MethodName)
	}
	lines = append(lines, fmt.Sprintf(`if let Some(err) = map_error(&%s, out.status, %s) {
    return Err(err);
}`, opConst, errArg))
	// The response's only repeated field (not counting the next-cursor
	// field), whose JSON name wraps a top-level array.
	var repeated []protoreflect.FieldDescriptor
	for _, f := range op.Response.Fields {
		if f.Desc.IsList() && string(f.Desc.Name()) != op.Opt.GetNextCursorField() {
			repeated = append(repeated, f.Desc)
		}
	}
	arrayArg := "None"
	if len(repeated) == 1 {
		arrayArg = "Some(" + rustQuote(repeated[0].JSONName()) + ")"
	}
	respDecl := "let resp"
	if cf := op.Opt.GetCursorField(); cf != "" {
		respDecl = "let mut resp"
	}
	if op.Transport == model.GraphQL {
		lines = append(lines, fmt.Sprintf("%s: %s = decode_graphql_response(%s.id, out.status, &out.body, &RESPONSE_FIELDS_%s, %s)?;",
			respDecl, op.Response.GoIdent.GoName, opConst, rustSnake(op.MethodName), arrayArg))
	} else {
		lines = append(lines, fmt.Sprintf("%s: %s = decode_response(%s.id, &out.body, %s)?;",
			respDecl, op.Response.GoIdent.GoName, opConst, arrayArg))
	}
	if cf := op.Opt.GetCursorField(); cf != "" {
		var wire string
		found := false
		for _, p := range op.Params {
			if string(p.Field.Desc.Name()) == cf {
				wire, found = p.Wire, true
				break
			}
		}
		if !found {
			return "", fmt.Errorf("protoc-gen-kei-ops: %s: cursor field %q is not a request param", op.ID, cf)
		}
		var nf *protogen.Field
		for _, f := range op.Response.Fields {
			if string(f.Desc.Name()) == op.Opt.GetNextCursorField() {
				nf = f
				break
			}
		}
		if nf == nil {
			return "", fmt.Errorf("protoc-gen-kei-ops: %s: next_cursor_field %q is not a field of %s", op.ID, op.Opt.GetNextCursorField(), op.Response.GoIdent.GoName)
		}
		name := string(nf.Desc.Name())
		optional := nf.Desc.HasPresence()
		switch {
		case nf.Desc.Kind() == protoreflect.StringKind && optional:
			lines = append(lines, fmt.Sprintf(`if resp.%s.is_none() {
    if let Some(cursor) = next_cursor_from_link(&out.headers, %s) {
        resp.%s = Some(cursor);
    }
}`, name, rustQuote(wire), name))
		case nf.Desc.Kind() == protoreflect.StringKind:
			lines = append(lines, fmt.Sprintf(`if resp.%s.is_empty() {
    if let Some(cursor) = next_cursor_from_link(&out.headers, %s) {
        resp.%s = cursor;
    }
}`, name, rustQuote(wire), name))
		case isIntKind(nf.Desc.Kind()) && optional:
			lines = append(lines, fmt.Sprintf(`if resp.%s.is_none() {
    if let Some(cursor) = next_cursor_from_link(&out.headers, %s) {
        resp.%s = cursor.parse::<%s>().ok();
    }
}`, name, rustQuote(wire), name, rustIntType(nf.Desc.Kind())))
		case isIntKind(nf.Desc.Kind()):
			lines = append(lines, fmt.Sprintf(`if resp.%s == 0 {
    if let Some(cursor) = next_cursor_from_link(&out.headers, %s) {
        if let Ok(v) = cursor.parse::<%s>() {
            resp.%s = v;
        }
    }
}`, name, rustQuote(wire), rustIntType(nf.Desc.Kind()), name))
		default:
			return "", fmt.Errorf("protoc-gen-kei-ops: %s: next cursor field %s must be a string or integer", op.ID, name)
		}
	}
	lines = append(lines, "Ok(resp)")
	joined := strings.Join(lines, "\n")
	var indented []string
	for _, l := range strings.Split(joined, "\n") {
		if l != "" {
			l = "        " + l
		}
		indented = append(indented, l)
	}
	return strings.Join(indented, "\n") + "\n", nil
}
