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

// Package model turns proto descriptors carrying the kei.contracts.v1
// connector and operation options into a language-neutral model that the
// validator and every renderer share. Placement defaults are resolved here,
// once, so all four languages place request fields the same way.
package model

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connection"
	contracts "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/contracts/v1"
)

// Set is every connector and operation found, each sorted by id.
type Set struct {
	Connectors []*Connector
	Ops        []*Operation
}

// Connector is one connector option and the operations that name it.
type Connector struct {
	ID string
	// Opt is the connector option as declared.
	Opt *contracts.Connector
	// File is the file that declares the option.
	File *protogen.File
	// Ops are the connector's operations, sorted by id.
	Ops []*Operation
}

// IsSQL reports whether every connection of the connector reaches a SQL
// engine: a SQL DSN pool, or a key-pair JWT connection that names the SQL
// engine it reaches (Snowflake, whose SQL API v2 runs statements over HTTP
// but keeps the contract's SQL request and response shape).
func (c *Connector) IsSQL() bool {
	conns := c.Opt.GetConnections()
	if len(conns) == 0 {
		return false
	}
	for _, conn := range conns {
		switch conn.GetKind() {
		case contracts.ConnectionKind_CONNECTION_KIND_SQL_DSN:
		case contracts.ConnectionKind_CONNECTION_KIND_KEY_PAIR_JWT:
			if conn.GetEngine() == "" {
				return false
			}
		default:
			return false
		}
	}
	return true
}

// Transport is how an operation reaches the provider.
type Transport int

const (
	// HTTP is a REST call: method, rendered path, query and JSON body.
	HTTP Transport = iota + 1
	// GraphQL posts the operation's fixed document with variables.
	GraphQL
	// SQL runs a statement over a SQL connection.
	SQL
	// MCP invokes a tool on a remote MCP server.
	MCP
)

// Operation is one message carrying the operation option.
type Operation struct {
	ID  string
	Opt *contracts.Operation
	// Connector is the connector named by Opt.Connector, nil when no file
	// declares it.
	Connector *Connector
	Request   *protogen.Message
	// Response is <Op>Response in the request's package, nil when missing.
	Response *protogen.Message
	// MethodName is Opt.MethodName, or the request name without "Request".
	MethodName string
	Transport  Transport
	// Params places every request field (HTTP and GraphQL), in field order.
	Params []*Param
}

// Param places one request field.
type Param struct {
	Field *protogen.Field
	In    contracts.ParamIn
	// Wire is the path variable, query key or body key.
	Wire string
	// Explicit is true when the field carries a param option.
	Explicit bool
	// MultiSegment keeps "/" in a path value, escaping each segment.
	MultiSegment bool
}

var pathVarRe = regexp.MustCompile(`\{([^{}]*)\}`)

// PathVars returns the variable names of a path template, in order.
func PathVars(template string) []string {
	var vars []string
	for _, m := range pathVarRe.FindAllStringSubmatch(template, -1) {
		vars = append(vars, m[1])
	}
	return vars
}

// Load reads the options from files (and only from files). It fails on a
// structural problem the validator cannot report per operation: two
// connector options with the same id.
func Load(files []*protogen.File) (*Set, error) {
	set := &Set{}
	byID := map[string]*Connector{}
	messagesByPkg := map[protoreflect.FullName]map[string]*protogen.Message{}

	for _, f := range files {
		if c, ok := proto.GetExtension(f.Desc.Options(), contracts.E_Connector).(*contracts.Connector); ok && c != nil && c.GetId() != "" {
			if prev, dup := byID[c.GetId()]; dup {
				return nil, fmt.Errorf("connector %q is declared in both %s and %s", c.GetId(), prev.File.Desc.Path(), f.Desc.Path())
			}
			conn := &Connector{ID: c.GetId(), Opt: c, File: f}
			byID[c.GetId()] = conn
			set.Connectors = append(set.Connectors, conn)
		}
		pkg := f.Desc.Package()
		if messagesByPkg[pkg] == nil {
			messagesByPkg[pkg] = map[string]*protogen.Message{}
		}
		for _, m := range f.Messages {
			messagesByPkg[pkg][string(m.Desc.Name())] = m
		}
	}

	for _, f := range files {
		for _, m := range f.Messages {
			opt, ok := proto.GetExtension(m.Desc.Options(), contracts.E_Operation).(*contracts.Operation)
			if !ok || opt == nil || opt.GetId() == "" {
				continue
			}
			op := &Operation{ID: opt.GetId(), Opt: opt, Request: m, Connector: byID[opt.GetConnector()]}
			name := string(m.Desc.Name())
			if base, ok := strings.CutSuffix(name, "Request"); ok {
				op.Response = messagesByPkg[f.Desc.Package()][base+"Response"]
			}
			op.MethodName = opt.GetMethodName()
			if op.MethodName == "" {
				op.MethodName = strings.TrimSuffix(name, "Request")
			}
			switch {
			case opt.GetGraphqlDocument() != "":
				op.Transport = GraphQL
			case op.Connector != nil && op.Connector.IsSQL() && opt.GetHttpMethod() == "":
				op.Transport = SQL
			default:
				op.Transport = HTTP
			}
			if op.Transport != SQL {
				op.Params = params(op)
			}
			set.Ops = append(set.Ops, op)
			if op.Connector != nil {
				op.Connector.Ops = append(op.Connector.Ops, op)
			}
		}
	}

	sort.Slice(set.Connectors, func(i, j int) bool { return set.Connectors[i].ID < set.Connectors[j].ID })
	sort.SliceStable(set.Ops, func(i, j int) bool { return set.Ops[i].ID < set.Ops[j].ID })
	for _, c := range set.Connectors {
		sort.SliceStable(c.Ops, func(i, j int) bool { return c.Ops[i].ID < c.Ops[j].ID })
	}
	return set, nil
}

// params resolves the placement of every request field:
//
//   - an explicit param option wins;
//   - a field named by a path-template variable goes in the path;
//   - otherwise GraphQL fields are variables (body), GET, DELETE and HEAD
//     fields go in the query, and the rest in the JSON body.
//
// The wire name is the param option's name, else the field's JSON name for
// a GraphQL variable, else connection.WireName. A path field's wire name is
// always its proto name, which the template uses.
func params(op *Operation) []*Param {
	inPath := map[string]bool{}
	for _, v := range PathVars(op.Opt.GetPathTemplate()) {
		inPath[v] = true
	}
	method := strings.ToUpper(op.Opt.GetHttpMethod())
	var out []*Param
	for _, f := range op.Request.Fields {
		p := &Param{Field: f, Wire: connection.WireName(f.Desc)}
		if op.Transport == GraphQL {
			// GraphQL variables follow the schema's camelCase: the JSON name.
			p.Wire = f.Desc.JSONName()
		}
		opt, _ := proto.GetExtension(f.Desc.Options(), contracts.E_Param).(*contracts.Param)
		if opt != nil && opt.GetIn() != contracts.ParamIn_PARAM_IN_UNSPECIFIED {
			p.In, p.Explicit = opt.GetIn(), true
		}
		if opt != nil && opt.GetName() != "" {
			p.Wire = opt.GetName()
		}
		p.MultiSegment = opt.GetMultiSegment()
		name := string(f.Desc.Name())
		switch {
		case p.Explicit:
		case inPath[name]:
			p.In = contracts.ParamIn_PARAM_IN_PATH
		case op.Transport == GraphQL:
			p.In = contracts.ParamIn_PARAM_IN_BODY
		case method == "GET" || method == "DELETE" || method == "HEAD":
			p.In = contracts.ParamIn_PARAM_IN_QUERY
		default:
			p.In = contracts.ParamIn_PARAM_IN_BODY
		}
		if p.In == contracts.ParamIn_PARAM_IN_PATH {
			p.Wire = name
		}
		out = append(out, p)
	}
	return out
}

// Field returns the request field named name, or nil.
func (op *Operation) Field(name string) *protogen.Field {
	for _, f := range op.Request.Fields {
		if string(f.Desc.Name()) == name {
			return f
		}
	}
	return nil
}
