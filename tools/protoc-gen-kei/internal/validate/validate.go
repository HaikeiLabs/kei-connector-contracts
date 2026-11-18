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

// Package validate is the validation pass both plugins run before they
// generate anything (proposal 0002, section 5.1). Every error starts with
// its rule name and a colon.
package validate

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"

	contracts "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/contracts/v1"
	"github.com/HaikeiLabs/kei-connector-contracts/tools/protoc-gen-kei/internal/model"
)

// Rule names. The first six are section 5.1's; the rest keep generated code
// sound.
const (
	RuleDuplicateID            = "duplicate_id"
	RuleConnectorMismatch      = "connector_mismatch"
	RuleWriteWithoutScope      = "write_without_scope"
	RuleUnknownPathVar         = "unknown_path_var"
	RuleUndeclaredConnection   = "undeclared_connection"
	RuleMissingResponse        = "missing_response"
	RuleMissingSemantics       = "missing_semantics"
	RuleSemanticsClassMismatch = "semantics_class_mismatch"
	RuleGraphQLWithoutDocument = "graphql_without_document"
	RuleUnknownConnector       = "unknown_connector"
	RuleMissingClass           = "missing_class"
	RuleMissingRoute           = "missing_route"
	RuleDuplicateMethod        = "duplicate_method"
	RuleUnsupportedParam       = "unsupported_param"
	RuleUnknownCursorField     = "unknown_cursor_field"
	RuleQueryShape             = "query_shape"
)

func errorf(rule, subject, format string, args ...any) error {
	return fmt.Errorf("%s: %s: %s", rule, subject, fmt.Sprintf(format, args...))
}

// Check returns every rule violation in set, sorted.
func Check(set *model.Set) []error {
	var errs []error
	seen := map[string]int{}
	for _, op := range set.Ops {
		seen[op.ID]++
		if seen[op.ID] == 2 {
			errs = append(errs, errorf(RuleDuplicateID, op.ID, "declared more than once"))
		}
		errs = append(errs, checkOp(op)...)
	}
	for _, c := range set.Connectors {
		methods := map[string]string{}
		for _, op := range c.Ops {
			if prev, dup := methods[op.MethodName]; dup && prev != op.ID {
				errs = append(errs, errorf(RuleDuplicateMethod, op.ID, "method %s is also %s", op.MethodName, prev))
			}
			methods[op.MethodName] = op.ID
		}
	}
	sort.Slice(errs, func(i, j int) bool { return errs[i].Error() < errs[j].Error() })
	return errs
}

func checkOp(op *model.Operation) []error {
	var errs []error
	opt := op.Opt
	add := func(rule, format string, args ...any) { errs = append(errs, errorf(rule, op.ID, format, args...)) }

	if first, _, _ := strings.Cut(op.ID, "."); first != opt.GetConnector() {
		add(RuleConnectorMismatch, "connector %q is not the id's first segment %q", opt.GetConnector(), first)
	}
	if op.Connector == nil {
		add(RuleUnknownConnector, "no file declares connector %q", opt.GetConnector())
	} else {
		declared := map[string]bool{}
		for _, c := range op.Connector.Opt.GetConnections() {
			declared[c.GetName()] = true
		}
		for _, name := range opt.GetConnections() {
			if !declared[name] {
				add(RuleUndeclaredConnection, "connector %q declares no connection %q", opt.GetConnector(), name)
			}
		}
		if len(opt.GetConnections()) == 0 {
			add(RuleUndeclaredConnection, "names no connection")
		}
	}

	switch opt.GetClass() {
	case contracts.OperationClass_OPERATION_CLASS_WRITE:
		if len(opt.GetScopes()) == 0 {
			add(RuleWriteWithoutScope, "a write operation must declare a scope")
		}
	case contracts.OperationClass_OPERATION_CLASS_READ:
	default:
		add(RuleMissingClass, "class must be read or write")
	}

	if op.Response == nil {
		add(RuleMissingResponse, "%s has no %sResponse message", op.Request.Desc.Name(), strings.TrimSuffix(string(op.Request.Desc.Name()), "Request"))
	}

	sem := opt.GetSemantics()
	switch {
	case sem == nil || sem.GetAction() == contracts.SemanticAction_SEMANTIC_ACTION_UNSPECIFIED || sem.GetResourceType() == "" || sem.GetVersion() == 0:
		add(RuleMissingSemantics, "semantics needs an action, a resource type and a version")
	case (sem.GetAction() == contracts.SemanticAction_SEMANTIC_ACTION_READ) != (opt.GetClass() == contracts.OperationClass_OPERATION_CLASS_READ):
		add(RuleSemanticsClassMismatch, "action %s does not match class %s", sem.GetAction(), opt.GetClass())
	}

	switch op.Transport {
	case model.SQL:
		errs = append(errs, checkQueryShape(op)...)
	default:
		if opt.GetPathTemplate() == "/graphql" && opt.GetGraphqlDocument() == "" {
			add(RuleGraphQLWithoutDocument, "POST /graphql needs a graphql_document")
		}
		if opt.GetHttpMethod() == "" || opt.GetPathTemplate() == "" {
			add(RuleMissingRoute, "an HTTP operation needs http_method and path_template")
		}
		for _, v := range model.PathVars(opt.GetPathTemplate()) {
			if op.Field(v) == nil {
				add(RuleUnknownPathVar, "path variable {%s} has no request field", v)
			}
		}
		for _, p := range op.Params {
			errs = append(errs, checkParam(op, p)...)
		}
	}

	if name := opt.GetCursorField(); name != "" && op.Field(name) == nil {
		add(RuleUnknownCursorField, "cursor_field %q is not a request field", name)
	}
	if name := opt.GetNextCursorField(); name != "" && (op.Response == nil || field(op.Response.Desc, name) == nil) {
		add(RuleUnknownCursorField, "next_cursor_field %q is not a response field", name)
	}
	return errs
}

func checkParam(op *model.Operation, p *model.Param) []error {
	fd := p.Field.Desc
	bad := func(format string, args ...any) []error {
		return []error{errorf(RuleUnsupportedParam, op.ID, "field %s: "+format, append([]any{fd.Name()}, args...)...)}
	}
	if p.MultiSegment && (p.In != contracts.ParamIn_PARAM_IN_PATH || fd.Kind() != protoreflect.StringKind) {
		return bad("multi_segment applies only to a string path field")
	}
	switch p.In {
	case contracts.ParamIn_PARAM_IN_HEADER:
		return bad("header parameters are not supported yet")
	case contracts.ParamIn_PARAM_IN_PATH:
		if fd.IsList() || fd.IsMap() || fd.Kind() == protoreflect.MessageKind {
			return bad("a path value must be a scalar")
		}
		if !slices.Contains(model.PathVars(op.Opt.GetPathTemplate()), string(fd.Name())) {
			return bad("placed in the path but the template has no {%s}", fd.Name())
		}
	case contracts.ParamIn_PARAM_IN_QUERY:
		if fd.IsMap() || fd.Kind() == protoreflect.MessageKind || fd.Kind() == protoreflect.GroupKind {
			return bad("a query value must be a scalar or a list of scalars")
		}
	}
	return nil
}

// checkQueryShape enforces the request and response shape CallSQL reads.
func checkQueryShape(op *model.Operation) []error {
	var errs []error
	add := func(format string, args ...any) { errs = append(errs, errorf(RuleQueryShape, op.ID, format, args...)) }
	req := op.Request.Desc
	want := func(md protoreflect.MessageDescriptor, name string, ok func(protoreflect.FieldDescriptor) bool, what string) {
		if fd := field(md, name); fd == nil || !ok(fd) {
			add("%s.%s must be %s", md.Name(), name, what)
		}
	}
	isMsg := func(full string) func(protoreflect.FieldDescriptor) bool {
		return func(fd protoreflect.FieldDescriptor) bool {
			return fd.Kind() == protoreflect.MessageKind && string(fd.Message().FullName()) == full
		}
	}
	want(req, "sql", func(fd protoreflect.FieldDescriptor) bool {
		return fd.Kind() == protoreflect.StringKind && !fd.IsList()
	}, "a string")
	want(req, "params", func(fd protoreflect.FieldDescriptor) bool { return fd.IsList() && isMsg("google.protobuf.Value")(fd) }, "repeated google.protobuf.Value")
	want(req, "max_rows", func(fd protoreflect.FieldDescriptor) bool { return fd.Kind() == protoreflect.Int64Kind && !fd.IsList() }, "an int64")
	if op.Response == nil {
		return errs
	}
	resp := op.Response.Desc
	want(resp, "columns", func(fd protoreflect.FieldDescriptor) bool {
		return fd.IsList() && fd.Kind() == protoreflect.MessageKind && field(fd.Message(), "name") != nil && field(fd.Message(), "type") != nil
	}, "repeated {name, type}")
	want(resp, "rows", func(fd protoreflect.FieldDescriptor) bool {
		return fd.IsList() && isMsg("google.protobuf.ListValue")(fd)
	}, "repeated google.protobuf.ListValue")
	want(resp, "truncated", func(fd protoreflect.FieldDescriptor) bool { return fd.Kind() == protoreflect.BoolKind && !fd.IsList() }, "a bool")
	return errs
}

func field(md protoreflect.MessageDescriptor, name string) protoreflect.FieldDescriptor {
	return md.Fields().ByName(protoreflect.Name(name))
}
