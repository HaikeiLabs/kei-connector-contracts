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
	"sort"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	contracts "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/contracts/v1"
	"github.com/HaikeiLabs/kei-connector-contracts/tools/protoc-gen-kei/internal/model"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

//go:embed python.py.tmpl
var pythonTemplateText string
var pythonTemplate = template.Must(template.New("python").Parse(pythonTemplateText))

type pyModule struct{ Package, Names string }
type pyOperation struct {
	Const, ID, Class, Destructive, Snake, Method string
	MethodJSON, PathJSON, Params, Errors         string
	GraphQL, CursorField, NextCursorField        string
	ResourceTemplate, MethodName, Connections    string
	Transport                                    string
	Req, Resp, PosArgs, PosBuild                 string
	SQL, Helper                                  bool
}
type pyFile struct {
	Source, Name, Runtime string
	SQL, Positional       bool
	Modules               []pyModule
	Ops                   []pyOperation
}

// pythonKeywords cannot be argument names or constructor keyword arguments.
var pythonKeywords = map[string]bool{
	"False": true, "None": true, "True": true, "and": true, "as": true,
	"assert": true, "async": true, "await": true, "break": true, "class": true,
	"continue": true, "def": true, "del": true, "elif": true, "else": true,
	"except": true, "finally": true, "for": true, "from": true, "global": true,
	"if": true, "import": true, "in": true, "is": true, "lambda": true,
	"nonlocal": true, "not": true, "or": true, "pass": true, "raise": true,
	"return": true, "try": true, "while": true, "with": true, "yield": true,
}

// pyStr renders a python string literal. Like the runtime's encoder it does
// not escape HTML characters, so the golden output matches what the runtime
// produces on the wire.
// pyConnections renders the connector's connection names as a tuple
// literal; empty stays an empty tuple.
func pyConnections(names []string) string {
	var b strings.Builder
	b.WriteString("(")
	for i, n := range names {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(pyStr(n))
	}
	if len(names) == 1 {
		b.WriteString(",")
	}
	b.WriteString(")")
	return b.String()
}

// pyTransport maps a model transport to the generated Python enum member.
func pyTransport(t model.Transport) string {
	switch t {
	case model.GraphQL:
		return "Transport.GRAPHQL"
	case model.SQL:
		return "Transport.SQL"
	case model.MCP:
		return "Transport.MCP"
	default:
		return "Transport.HTTP"
	}
}

func pyStr(s string) string {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(s)
	return strings.TrimSuffix(b.String(), "\n")
}

// pyLiteral renders a python literal for the generated spec dictionaries:
// deterministic key order, no HTML escaping.
func pyLiteral(v any) string {
	switch x := v.(type) {
	case string:
		return pyStr(x)
	case bool:
		if x {
			return "True"
		}
		return "False"
	case int:
		return strconv.Itoa(x)
	case []any:
		parts := make([]string, 0, len(x))
		for _, e := range x {
			parts = append(parts, pyLiteral(e))
		}
		return "[" + strings.Join(parts, ", ") + "]"
	case map[string]any:
		keys := make([]string, 0, len(x))
		for k := range x {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		parts := make([]string, 0, len(x))
		for _, k := range keys {
			parts = append(parts, pyStr(k)+": "+pyLiteral(x[k]))
		}
		return "{" + strings.Join(parts, ", ") + "}"
	default:
		panic(fmt.Sprintf("protoc-gen-kei-ops: pyLiteral: unsupported type %T", v))
	}
}

func pyScalarKind(k protoreflect.Kind) string {
	switch k {
	case protoreflect.BoolKind:
		return "bool"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
		protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind,
		protoreflect.Uint32Kind, protoreflect.Fixed32Kind, protoreflect.Uint64Kind,
		protoreflect.Fixed64Kind:
		return "int"
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		return "float"
	case protoreflect.StringKind:
		return "str"
	case protoreflect.BytesKind:
		return "bytes"
	case protoreflect.EnumKind:
		return "int"
	default:
		panic(fmt.Sprintf("protoc-gen-kei-ops: pyScalarKind: unsupported kind %s", k))
	}
}

// pyMessageRef returns the reference to a message in generated python:
// things_pb2.GetThingRequest, or things_pb2.Outer.Inner for a nested one.
func pyMessageRef(m protoreflect.MessageDescriptor) string {
	name := string(m.Name())
	for p := m.Parent(); p != nil; p = p.Parent() {
		pm, ok := p.(protoreflect.MessageDescriptor)
		if !ok {
			break
		}
		name = string(pm.Name()) + "." + name
	}
	mod := path.Base(strings.TrimSuffix(m.ParentFile().Path(), ".proto") + "_pb2")
	return mod + "." + name
}

func pyType(f *protogen.Field) string {
	if f.Desc.IsMap() {
		val := f.Desc.MapValue()
		if val.Kind() == protoreflect.MessageKind || val.Kind() == protoreflect.GroupKind {
			return "dict[str, " + pyMessageRef(val.Message()) + "]"
		}
		return "dict[str, " + pyScalarKind(val.Kind()) + "]"
	}
	if f.Desc.IsList() {
		elem := ""
		if f.Desc.Kind() == protoreflect.MessageKind || f.Desc.Kind() == protoreflect.GroupKind {
			elem = pyMessageRef(f.Desc.Message())
		} else if f.Desc.Kind() == protoreflect.EnumKind {
			elem = "int"
		} else {
			elem = pyScalarKind(f.Desc.Kind())
		}
		return "list[" + elem + "]"
	}
	if f.Desc.Kind() == protoreflect.MessageKind || f.Desc.Kind() == protoreflect.GroupKind {
		return pyMessageRef(f.Desc.Message())
	}
	return pyScalarKind(f.Desc.Kind())
}

// pySnakeMethod converts a camelCase operation method name to the snake_case
// python convention: GetPullRequest -> get_pull_request.
func pySnakeMethod(name string) string {
	var b strings.Builder
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i > 0 {
				b.WriteByte('_')
			}
			b.WriteRune(unicode.ToLower(r))
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func renderPython(gen *protogen.Plugin, c *model.Connector, opts Options) error {
	dir := path.Dir(c.File.Desc.Path())
	d := pyFile{Source: c.File.Desc.Path(), Name: tsPascal(c.ID), SQL: c.IsSQL(), Positional: opts.positional()}
	rt := map[string]bool{"Op": true, "OpClass": true, "Transport": true}
	if d.SQL {
		rt["SqlConnection"] = true
		rt["call_sql"] = true
	} else {
		rt["HttpConnection"] = true
		rt["call_http"] = true
	}
	modules := map[string]map[string]bool{}
	addModuleDesc := func(md protoreflect.MessageDescriptor) {
		file := strings.TrimSuffix(md.ParentFile().Path(), ".proto") + "_pb2"
		pkg := strings.ReplaceAll(path.Dir(file), "/", ".")
		if modules[pkg] == nil {
			modules[pkg] = map[string]bool{}
		}
		modules[pkg][path.Base(file)] = true
	}
	addModule := func(m *protogen.Message) { addModuleDesc(m.Desc) }

	for _, op := range c.Ops {
		o := pyOperation{
			Const:       "OP_" + tsSnake(op.MethodName),
			ID:          pyStr(op.ID),
			Class:       "READ",
			Destructive: "False",
			Method:      pySnakeMethod(op.MethodName),
			Req:         pyMessageRef(op.Request.Desc),
			Resp:        pyMessageRef(op.Response.Desc),
			SQL:         op.Transport == model.SQL,
		}
		if op.Opt.GetClass() == contracts.OperationClass_OPERATION_CLASS_WRITE {
			o.Class = "WRITE"
		}
		if op.Opt.GetDestructive() {
			o.Destructive = "True"
		}
		o.ResourceTemplate = pyStr(op.Opt.GetResourceTemplate())
		o.MethodName = pyStr(op.MethodName)
		o.Connections = pyConnections(op.Opt.GetConnections())
		o.Transport = pyTransport(op.Transport)
		o.Snake = strings.TrimPrefix(o.Const, "OP_")
		addModule(op.Request)
		addModule(op.Response)
		o.Helper = !opts.positional() && op.Opt.GetPositionalHelper()
		if opts.positional() || o.Helper {
			for _, oneof := range op.Request.Oneofs {
				if !oneof.Desc.IsSynthetic() {
					return fmt.Errorf("protoc-gen-kei-ops: %s: positional signatures do not support oneof %s", op.ID, oneof.Desc.Name())
				}
			}
			var args, buildParts, dictParts []string
			used := map[string]bool{"ops": true, "self": true}
			kwOK := true
			for _, f := range op.Request.Fields {
				name := string(f.Desc.Name())
				if pythonKeywords[name] {
					kwOK = false
				}
				arg := name
				if pythonKeywords[arg] || arg == "self" || arg == "ops" {
					arg += "_"
				}
				for used[arg] {
					arg += "_"
				}
				used[arg] = true
				args = append(args, arg+": "+pyType(f))
				buildParts = append(buildParts, name+"="+arg)
				dictParts = append(dictParts, pyStr(name)+": "+arg)
			}
			o.PosArgs = strings.Join(args, ", ")
			switch {
			case len(buildParts) == 0:
				o.PosBuild = o.Req + "()"
			case kwOK:
				o.PosBuild = o.Req + "(" + strings.Join(buildParts, ", ") + ")"
			default:
				o.PosBuild = o.Req + "(**{" + strings.Join(dictParts, ", ") + "})"
			}
		}
		params := []any{}
		for _, p := range op.Params {
			loc := "body"
			if p.In == contracts.ParamIn_PARAM_IN_PATH {
				loc = "path"
			}
			if p.In == contracts.ParamIn_PARAM_IN_QUERY {
				loc = "query"
			}
			params = append(params, map[string]any{"field": string(p.Field.Desc.Name()), "in": loc, "wire": p.Wire, "multi_segment": p.MultiSegment})
		}
		o.Params = pyLiteral(params)
		mappings := []any{}
		for _, e := range op.Opt.GetErrors() {
			mappings = append(mappings, map[string]any{"status": int(e.GetStatus()), "kind": strings.ToLower(strings.TrimPrefix(e.GetKind().String(), "ERROR_KIND_")), "retryable": e.GetRetryable()})
		}
		o.Errors = pyLiteral(mappings)
		o.MethodJSON = pyStr(strings.ToUpper(op.Opt.GetHttpMethod()))
		o.PathJSON = pyStr(op.Opt.GetPathTemplate())
		o.GraphQL = pyStr(op.Opt.GetGraphqlDocument())
		o.CursorField = pyStr(op.Opt.GetCursorField())
		o.NextCursorField = pyStr(op.Opt.GetNextCursorField())
		d.Ops = append(d.Ops, o)
	}
	// Positional signatures annotate helper arguments with field types, so
	// the modules of those field types must be imported too.
	needsFieldTypes := opts.positional()
	for _, o := range d.Ops {
		if o.Helper {
			needsFieldTypes = true
		}
	}
	if needsFieldTypes {
		for _, op := range c.Ops {
			for _, m := range []*protogen.Message{op.Request, op.Response} {
				addModule(m)
				for _, f := range m.Fields {
					fd := f.Desc
					if fd.IsMap() {
						if mv := fd.MapValue(); mv.Kind() == protoreflect.MessageKind || mv.Kind() == protoreflect.GroupKind {
							addModuleDesc(mv.Message())
						}
						continue
					}
					if fd.Kind() == protoreflect.MessageKind || fd.Kind() == protoreflect.GroupKind {
						addModuleDesc(fd.Message())
					}
				}
			}
		}
	}
	pkgs := make([]string, 0, len(modules))
	for p := range modules {
		pkgs = append(pkgs, p)
	}
	sort.Strings(pkgs)
	for _, p := range pkgs {
		names := make([]string, 0, len(modules[p]))
		for n := range modules[p] {
			names = append(names, n)
		}
		sort.Strings(names)
		d.Modules = append(d.Modules, pyModule{Package: p, Names: strings.Join(names, ", ")})
	}
	rts := make([]string, 0, len(rt))
	for n := range rt {
		rts = append(rts, n)
	}
	sort.Strings(rts)
	d.Runtime = strings.Join(rts, ", ")
	var b bytes.Buffer
	if err := pythonTemplate.Execute(&b, d); err != nil {
		return err
	}
	g := gen.NewGeneratedFile(dir+"/ops.py", "")
	_, err := g.Write(b.Bytes())
	return err
}
