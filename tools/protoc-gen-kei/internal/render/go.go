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
	"go/token"
	"path"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	contracts "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/contracts/v1"
	"github.com/HaikeiLabs/kei-connector-contracts/tools/protoc-gen-kei/internal/model"
)

//go:embed go.go.tmpl
var goTemplateText string

var goTemplate = template.Must(template.New("go").Parse(goTemplateText))

const (
	connectionImport = protogen.GoImportPath("github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connection")
	contractsImport  = protogen.GoImportPath("github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/contracts/v1")
)

// GoPackageName is the generated Go package of a connector: its id without
// underscores, plus "ops" (github → githubops, http_api → httpapiops).
func GoPackageName(connectorID string) string {
	return strings.ReplaceAll(connectorID, "_", "") + "ops"
}

type goFileData struct {
	Pkg, ID, Source string
	Positional      bool
	SQL             bool
	// Qualifiers, such as "context." and "connection.".
	Ctx, Conn, Errors string
	Ops               []goOpData
	// Imports is the import block's body, one spec per line.
	Imports []string
}

// goImports names the packages a generated file uses: each by its real Go
// package name (github, contracts, structpb), made unique when two collide.
type goImports struct {
	names  map[protogen.GoImportPath]string // known package names
	byPath map[protogen.GoImportPath]string
	used   map[string]bool
}

func newGoImports(gen *protogen.Plugin, self string) *goImports {
	im := &goImports{
		names:  map[protogen.GoImportPath]string{connectionImport: "connection", contractsImport: "contracts"},
		byPath: map[protogen.GoImportPath]string{},
		used:   map[string]bool{self: true},
	}
	for _, f := range gen.Files {
		if _, ok := im.names[f.GoImportPath]; !ok {
			im.names[f.GoImportPath] = string(f.GoPackageName)
		}
	}
	return im
}

// q returns the qualifier ("name.") for an import path, importing it.
func (im *goImports) q(ip protogen.GoImportPath) string {
	if name, ok := im.byPath[ip]; ok {
		return name + "."
	}
	name, ok := im.names[ip]
	if !ok {
		name = path.Base(string(ip))
	}
	for orig, i := name, 2; im.used[name]; i++ {
		name = orig + strconv.Itoa(i)
	}
	im.used[name] = true
	im.byPath[ip] = name
	return name + "."
}

func (im *goImports) ident(id protogen.GoIdent) string { return im.q(id.GoImportPath) + id.GoName }

// specs renders the import specs, sorted by path.
func (im *goImports) specs() []string {
	paths := make([]string, 0, len(im.byPath))
	for p := range im.byPath {
		paths = append(paths, string(p))
	}
	sort.Strings(paths)
	out := make([]string, len(paths))
	for i, p := range paths {
		name := im.byPath[protogen.GoImportPath(p)]
		if name == path.Base(p) {
			out[i] = strconv.Quote(p)
		} else {
			out[i] = name + " " + strconv.Quote(p)
		}
	}
	return out
}

type goOpData struct {
	ID, Method, Route    string
	OpVar, SpecVar       string
	OpLiteral, SpecValue string
	Deprecated           string // the deprecation note, empty when not deprecated
	Req, Resp            string
	// Params is the method's parameter list after ctx; Args forwards them.
	Params, Args string
	// Build constructs the request from positional parameters.
	Build      string
	Helper     bool
	HelperArgs string // the positional parameters, for the helper
	SQL        bool
}

func renderGo(gen *protogen.Plugin, c *model.Connector, opts Options) error {
	pkg := GoPackageName(c.ID)
	dir := path.Dir(c.File.GeneratedFilenamePrefix)
	importPath := c.File.GoImportPath + protogen.GoImportPath("/"+pkg)

	files := map[string]string{"ops": "ops.go", "mock": "ops_mock.go"}
	for _, name := range []string{"ops", "mock"} {
		g := gen.NewGeneratedFile(dir+"/"+pkg+"/"+files[name], importPath)
		im := newGoImports(gen, pkg)
		data := goFileData{
			Pkg: pkg, ID: c.ID, Source: dir,
			Positional: opts.positional(), SQL: c.IsSQL(),
			Ctx: im.q("context"),
		}
		if name == "ops" {
			data.Conn = im.q(connectionImport)
		} else {
			data.Errors = im.q("errors")
		}
		for _, op := range c.Ops {
			od, err := goOp(im, op, opts, name == "ops")
			if err != nil {
				return err
			}
			data.Ops = append(data.Ops, od)
		}
		data.Imports = im.specs()
		var buf bytes.Buffer
		if err := goTemplate.ExecuteTemplate(&buf, name, data); err != nil {
			return fmt.Errorf("protoc-gen-kei-ops: %s: %w", c.ID, err)
		}
		if _, err := g.Write(buf.Bytes()); err != nil {
			return err
		}
	}
	return nil
}

// goOp builds one operation's template data. withSpec is false for the mock
// file, which needs only the signatures (and so imports nothing else).
func goOp(im *goImports, op *model.Operation, opts Options, withSpec bool) (goOpData, error) {
	opt := op.Opt
	d := goOpData{
		ID: op.ID, Method: op.MethodName,
		OpVar: "Op" + op.MethodName, SpecVar: "spec" + op.MethodName,
		Req:  im.ident(op.Request.GoIdent),
		Resp: im.ident(op.Response.GoIdent),
		SQL:  op.Transport == model.SQL,
	}
	switch op.Transport {
	case model.HTTP:
		d.Route = opt.GetHttpMethod() + " " + opt.GetPathTemplate()
	case model.GraphQL:
		d.Route = "GraphQL " + opt.GetHttpMethod() + " " + opt.GetPathTemplate()
	case model.SQL:
		d.Route = "SQL query"
	}
	if opt.GetDeprecated() {
		d.Deprecated = "this operation is deprecated."
		if r := opt.GetReplacedBy(); r != "" {
			d.Deprecated = "use " + r + " instead."
		}
	}

	if withSpec {
		class := "Read"
		if opt.GetClass() == contracts.OperationClass_OPERATION_CLASS_WRITE {
			class = "Write"
		}
		conn := im.q(connectionImport)
		d.OpLiteral = fmt.Sprintf("%sOp{ID: %s, Class: %s%s", conn, strconv.Quote(op.ID), conn, class)
		if opt.GetDestructive() {
			d.OpLiteral += ", Destructive: true"
		}
		d.OpLiteral += fmt.Sprintf(", ResourceTemplate: %s, MethodName: %s", strconv.Quote(opt.GetResourceTemplate()), strconv.Quote(op.MethodName))
		if conns := opt.GetConnections(); len(conns) > 0 {
			d.OpLiteral += ", Connections: []string{"
			for i, c := range conns {
				if i > 0 {
					d.OpLiteral += ", "
				}
				d.OpLiteral += strconv.Quote(c)
			}
			d.OpLiteral += "}"
		} else {
			d.OpLiteral += ", Connections: nil"
		}
		d.OpLiteral += fmt.Sprintf(", Transport: %s%s}", conn, transportConst(op.Transport))
		if op.Transport != model.SQL {
			d.SpecValue = goSpec(op, d.OpVar, im)
		}
	}

	needPositional := opts.positional() || opt.GetPositionalHelper()
	var params, args, build []string
	if needPositional {
		if len(op.Request.Oneofs) > 0 {
			for _, o := range op.Request.Oneofs {
				if !o.Desc.IsSynthetic() {
					return d, fmt.Errorf("protoc-gen-kei-ops: %s: positional signatures do not support oneof %s", op.ID, o.Desc.Name())
				}
			}
		}
		used := map[string]bool{"ctx": true, "ops": true, "q": true, "m": true, "req": true, "resp": true}
		for _, f := range op.Request.Fields {
			name := paramName(f.GoName, used)
			params = append(params, name+" "+goType(im, f))
			args = append(args, name)
			build = append(build, f.GoName+": "+name)
		}
		d.Build = "&" + d.Req + "{" + strings.Join(build, ", ") + "}"
		d.HelperArgs = strings.Join(params, ", ")
	}
	if opts.positional() {
		d.Params, d.Args = strings.Join(params, ", "), strings.Join(args, ", ")
	} else {
		d.Params, d.Args = "req *"+d.Req, "req"
		d.Helper = opt.GetPositionalHelper()
	}
	return d, nil
}

// goSpec renders the connection.HTTPSpec literal of an HTTP or GraphQL op.
// transportConst maps a model transport to the generated Go constant name.
func transportConst(t model.Transport) string {
	switch t {
	case model.GraphQL:
		return "TransportGraphQL"
	case model.SQL:
		return "TransportSQL"
	case model.MCP:
		return "TransportMCP"
	default:
		return "TransportHTTP"
	}
}

func goSpec(op *model.Operation, opVar string, im *goImports) string {
	opt := op.Opt
	conn := im.q(connectionImport)
	var b strings.Builder
	fmt.Fprintf(&b, "&%sHTTPSpec{\n", conn)
	fmt.Fprintf(&b, "Op: %s,\n", opVar)
	fmt.Fprintf(&b, "Method: %s,\n", strconv.Quote(strings.ToUpper(opt.GetHttpMethod())))
	fmt.Fprintf(&b, "PathTemplate: %s,\n", strconv.Quote(opt.GetPathTemplate()))
	if len(op.Params) > 0 {
		b.WriteString("Params: []" + conn + "Param{\n")
		for _, p := range op.Params {
			multi := ""
			if p.MultiSegment {
				multi = ", MultiSegment: true"
			}
			fmt.Fprintf(&b, "{Field: %s, In: %s%s, Wire: %s%s},\n",
				strconv.Quote(string(p.Field.Desc.Name())), conn, goParamIn(p.In), strconv.Quote(p.Wire), multi)
		}
		b.WriteString("},\n")
	}
	if len(opt.GetErrors()) > 0 {
		ct := im.q(contractsImport)
		b.WriteString("Errors: []*" + ct + "ErrorMapping{\n")
		for _, e := range opt.GetErrors() {
			fmt.Fprintf(&b, "{Status: %d, Kind: %sErrorKind_%s", e.GetStatus(), ct, e.GetKind())
			if e.GetRetryable() {
				b.WriteString(", Retryable: true")
			}
			b.WriteString("},\n")
		}
		b.WriteString("},\n")
	}
	if doc := opt.GetGraphqlDocument(); doc != "" {
		fmt.Fprintf(&b, "GraphQLDocument: %s,\n", strconv.Quote(doc))
	}
	if f := opt.GetCursorField(); f != "" {
		fmt.Fprintf(&b, "CursorField: %s,\n", strconv.Quote(f))
	}
	if f := opt.GetNextCursorField(); f != "" {
		fmt.Fprintf(&b, "NextCursorField: %s,\n", strconv.Quote(f))
	}
	b.WriteString("}")
	return b.String()
}

func goParamIn(in contracts.ParamIn) string {
	switch in {
	case contracts.ParamIn_PARAM_IN_PATH:
		return "InPath"
	case contracts.ParamIn_PARAM_IN_QUERY:
		return "InQuery"
	default:
		return "InBody"
	}
}

// paramName is a Go parameter name for a field: lowerCamelCase, made unique
// and kept clear of keywords and the generated code's own names.
func paramName(goName string, used map[string]bool) string {
	name := strings.ToLower(goName[:1]) + goName[1:]
	if token.IsKeyword(name) || used[name] {
		name += "_"
	}
	for used[name] {
		name += "_"
	}
	used[name] = true
	return name
}

// goType is the Go type of a generated struct field.
func goType(im *goImports, f *protogen.Field) string {
	if f.Desc.IsMap() {
		return "map[" + goType(im, f.Message.Fields[0]) + "]" + goType(im, f.Message.Fields[1])
	}
	var t string
	switch f.Desc.Kind() {
	case protoreflect.BoolKind:
		t = "bool"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		t = "int32"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		t = "int64"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		t = "uint32"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		t = "uint64"
	case protoreflect.FloatKind:
		t = "float32"
	case protoreflect.DoubleKind:
		t = "float64"
	case protoreflect.StringKind:
		t = "string"
	case protoreflect.BytesKind:
		t = "[]byte"
	case protoreflect.EnumKind:
		t = im.ident(f.Enum.GoIdent)
	case protoreflect.MessageKind, protoreflect.GroupKind:
		t = "*" + im.ident(f.Message.GoIdent)
	}
	switch {
	case f.Desc.IsList():
		return "[]" + t
	case f.Desc.HasPresence() && f.Desc.Kind() != protoreflect.MessageKind && f.Desc.Kind() != protoreflect.GroupKind && f.Desc.Kind() != protoreflect.BytesKind:
		return "*" + t
	}
	return t
}
