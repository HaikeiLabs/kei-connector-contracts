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
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"unicode"

	contracts "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/contracts/v1"
	"github.com/HaikeiLabs/kei-connector-contracts/tools/protoc-gen-kei/internal/model"
	"google.golang.org/protobuf/compiler/protogen"
)

//go:embed ts.ts.tmpl
var tsTemplateText string
var tsTemplate = template.Must(template.New("ts").Parse(tsTemplateText))

type tsImport struct{ Path, Names string }
type tsOperation struct {
	Const, Descriptor, Method, Req, Resp, Spec string
	Params, Build                              string
	SQL, Helper                                bool
}
type tsFile struct {
	Source, Runtime, Name        string
	SQL, Positional, NeedsCreate bool
	Imports                      []tsImport
	Ops                          []tsOperation
}

func tsRelative(from, to string) string {
	r, _ := filepath.Rel(from, to)
	r = filepath.ToSlash(r)
	if !strings.HasPrefix(r, ".") {
		r = "./" + r
	}
	return r
}
func tsJSON(v any) string { b, _ := json.Marshal(v); return string(b) }

// tsConnections renders the connector's connection names as a JSON array.
func tsConnections(names []string) string { return tsJSON(names) }

// transportKind maps a model transport to the generated TS literal.
func transportKind(t model.Transport) string {
	switch t {
	case model.GraphQL:
		return "graphql"
	case model.SQL:
		return "sql"
	case model.MCP:
		return "mcp"
	default:
		return "http"
	}
}
func tsPascal(s string) string {
	var out strings.Builder
	for _, p := range strings.Split(s, "_") {
		if p != "" {
			out.WriteString(strings.ToUpper(p[:1]) + p[1:])
		}
	}
	return out.String()
}
func tsSnake(s string) string {
	var out strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			out.WriteByte('_')
		}
		out.WriteRune(unicode.ToUpper(r))
	}
	return out.String()
}
func renderTS(gen *protogen.Plugin, c *model.Connector, opts Options) error {
	dir := path.Dir(c.File.Desc.Path())
	d := tsFile{Source: c.File.Desc.Path(), Runtime: tsRelative(dir, "kei/connection/index.js"), Name: tsPascal(c.ID), SQL: c.IsSQL(), Positional: opts.positional()}
	imports := map[string]map[string]bool{}
	for _, op := range c.Ops {
		for _, m := range []*protogen.Message{op.Request, op.Response} {
			file := strings.TrimSuffix(m.Desc.ParentFile().Path(), ".proto") + "_pb.js"
			p := tsRelative(dir, file)
			if imports[p] == nil {
				imports[p] = map[string]bool{}
			}
			imports[p][m.GoIdent.GoName] = true
			imports[p][m.GoIdent.GoName+"Schema"] = true
		}
		class := "read"
		if op.Opt.GetClass() == contracts.OperationClass_OPERATION_CLASS_WRITE {
			class = "write"
		}
		descriptor := fmt.Sprintf("{ id: %s, class: %s, destructive: %t, resourceTemplate: %s, methodName: %s, connections: %s, transport: %s }", tsJSON(op.ID), tsJSON(class), op.Opt.GetDestructive(), tsJSON(op.Opt.GetResourceTemplate()), tsJSON(op.MethodName), tsConnections(op.Opt.GetConnections()), tsJSON(transportKind(op.Transport)))
		name := "OP_" + tsSnake(op.MethodName)
		o := tsOperation{Const: name, Descriptor: descriptor, Method: strings.ToLower(op.MethodName[:1]) + op.MethodName[1:], Req: op.Request.GoIdent.GoName, Resp: op.Response.GoIdent.GoName, SQL: op.Transport == model.SQL}
		o.Helper = !opts.positional() && op.Opt.GetPositionalHelper()
		if opts.positional() || o.Helper {
			d.NeedsCreate = true
			for _, oneof := range op.Request.Oneofs {
				if !oneof.Desc.IsSynthetic() {
					return fmt.Errorf("protoc-gen-kei-ops: %s: positional signatures do not support oneof %s", op.ID, oneof.Desc.Name())
				}
			}
			var args, properties []string
			used := map[string]bool{"req": true, "ops": true}
			for _, f := range op.Request.Fields {
				property := string(f.Desc.Name())
				parts := strings.Split(property, "_")
				for i := 1; i < len(parts); i++ {
					if parts[i] != "" {
						parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
					}
				}
				property = strings.Join(parts, "")
				arg := property
				switch arg {
				case "default", "class", "function", "new", "return", "var", "let", "const", "delete", "this", "export", "import", "extends", "super", "void", "typeof", "in", "instanceof", "await", "yield":
					arg += "_"
				}
				for used[arg] {
					arg += "_"
				}
				used[arg] = true
				args = append(args, fmt.Sprintf("%s: %s[%s]", arg, o.Req, tsJSON(property)))
				properties = append(properties, fmt.Sprintf("%s: %s", tsJSON(property), arg))
			}
			o.Params = strings.Join(args, ", ")
			o.Build = "create(" + o.Req + "Schema, { " + strings.Join(properties, ", ") + " })"
		}
		params := []map[string]any{}
		for _, p := range op.Params {
			loc := "body"
			if p.In == contracts.ParamIn_PARAM_IN_PATH {
				loc = "path"
			}
			if p.In == contracts.ParamIn_PARAM_IN_QUERY {
				loc = "query"
			}
			params = append(params, map[string]any{"field": string(p.Field.Desc.Name()), "in": loc, "wire": p.Wire, "multiSegment": p.MultiSegment})
		}
		mappings := []map[string]any{}
		for _, e := range op.Opt.GetErrors() {
			mappings = append(mappings, map[string]any{"status": e.GetStatus(), "kind": strings.ToLower(strings.TrimPrefix(e.GetKind().String(), "ERROR_KIND_")), "retryable": e.GetRetryable()})
		}
		spec := map[string]any{"op": json.RawMessage(descriptorJSON(op, class)), "method": strings.ToUpper(op.Opt.GetHttpMethod()), "pathTemplate": op.Opt.GetPathTemplate(), "params": params, "errors": mappings, "graphqlDocument": op.Opt.GetGraphqlDocument(), "cursorField": op.Opt.GetCursorField(), "nextCursorField": op.Opt.GetNextCursorField()}
		o.Spec = tsJSON(spec)
		d.Ops = append(d.Ops, o)
	}
	paths := make([]string, 0, len(imports))
	for p := range imports {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	for _, p := range paths {
		names := make([]string, 0, len(imports[p]))
		for n := range imports[p] {
			names = append(names, n)
		}
		sort.Strings(names)
		d.Imports = append(d.Imports, tsImport{p, strings.Join(names, ", ")})
	}
	var b bytes.Buffer
	if err := tsTemplate.Execute(&b, d); err != nil {
		return err
	}
	g := gen.NewGeneratedFile(dir+"/ops.ts", "")
	_, err := g.Write(b.Bytes())
	return err
}
func descriptorJSON(op *model.Operation, class string) string {
	return tsJSON(map[string]any{"id": op.ID, "class": class, "destructive": op.Opt.GetDestructive(), "resourceTemplate": op.Opt.GetResourceTemplate(), "methodName": op.MethodName, "connections": op.Opt.GetConnections(), "transport": transportKind(op.Transport)})
}
