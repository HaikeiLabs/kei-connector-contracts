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
	// result_required rules: a write must declare the response paths that
	// prove the write happened, and every path must resolve to a bool,
	// string or message field. Writes whose response has no fields are
	// exempt. Reads never declare the option.
	RuleMissingResultFields    = "missing_result_fields"
	RuleUnsupportedResultField = "unsupported_result_field"
	RuleResultFieldsOnRead     = "result_fields_on_read"
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
	errs = append(errs, checkResultRequired(op)...)
	return errs
}

// checkResultRequired enforces the result_required option: every write with
// a fielded response declares the response paths that prove the write
// happened, every path resolves to a bool (true), string (non-empty),
// numeric (non-zero) or message (set) field, and a read never declares the
// option. A write whose response has no fields is exempt: nothing in the
// response could carry a proof.
func checkResultRequired(op *model.Operation) []error {
	opt := op.Opt
	switch opt.GetClass() {
	case contracts.OperationClass_OPERATION_CLASS_WRITE:
	case contracts.OperationClass_OPERATION_CLASS_READ:
		if len(opt.GetResultRequired()) > 0 {
			return []error{errorf(RuleResultFieldsOnRead, op.ID, "result_required applies only to write operations")}
		}
		return nil
	default:
		return nil
	}
	if op.Response == nil || op.Response.Desc.Fields().Len() == 0 {
		return nil
	}
	var errs []error
	if len(opt.GetResultRequired()) == 0 {
		errs = append(errs, errorf(RuleMissingResultFields, op.ID,
			"a write operation must declare result_required (response paths that prove the write happened)"))
	}
	for _, path := range opt.GetResultRequired() {
		if err := resolveResultPath(op.Response.Desc, path); err != nil {
			errs = append(errs, errorf(RuleUnsupportedResultField, op.ID, "result_required %q: %v", path, err))
		}
	}
	return errs
}

// resolveResultPath walks a result_required dot path through the response
// message and checks that the final field proves the write: a bool (must be
// true), a string (must be non-empty), a plain numeric scalar (must be
// non-zero, matching proto3's no-presence semantics) or a message (must be
// set). Every segment must name a field; an intermediate segment must be a
// singular message the walk can continue through. Enums, bytes, lists,
// maps and presence-carrying scalars cannot prove a write.
func resolveResultPath(md protoreflect.MessageDescriptor, path string) error {
	var fd protoreflect.FieldDescriptor
	segs := strings.Split(path, ".")
	for i, seg := range segs {
		if fd = field(md, seg); fd == nil {
			return fmt.Errorf("%s has no field %s", md.FullName(), seg)
		}
		if i == len(segs)-1 {
			break
		}
		if (fd.Kind() != protoreflect.MessageKind && fd.Kind() != protoreflect.GroupKind) || fd.IsList() || fd.IsMap() {
			return fmt.Errorf("field %s does not lead to a message", seg)
		}
		md = fd.Message()
	}
	if fd.IsMap() || fd.IsList() {
		return fmt.Errorf("field %s is %s, not a bool, a string, a number or a message", fd.Name(), kindWord(fd))
	}
	switch fd.Kind() {
	case protoreflect.BoolKind, protoreflect.StringKind, protoreflect.MessageKind, protoreflect.GroupKind:
		return nil
	case protoreflect.DoubleKind, protoreflect.FloatKind, protoreflect.Int32Kind, protoreflect.Int64Kind,
		protoreflect.Uint32Kind, protoreflect.Uint64Kind, protoreflect.Sint32Kind, protoreflect.Sint64Kind,
		protoreflect.Fixed32Kind, protoreflect.Fixed64Kind, protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind:
		if fd.HasPresence() {
			return fmt.Errorf("field %s is optional; a numeric result field must be a plain scalar", fd.Name())
		}
		return nil
	default:
		return fmt.Errorf("field %s is %s, not a bool, a string, a number or a message", fd.Name(), kindWord(fd))
	}
}

func kindWord(fd protoreflect.FieldDescriptor) string {
	switch {
	case fd.IsMap():
		return "a map"
	case fd.IsList():
		return "a repeated " + fd.Kind().String()
	case fd.Kind() == protoreflect.EnumKind:
		return "an enum"
	default:
		return "a " + fd.Kind().String()
	}
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
