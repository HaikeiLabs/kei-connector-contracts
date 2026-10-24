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

// Package calltemplate validates tool call templates (HAI-856, ADR-039).
//
// A template is a kei.tool-definition/v1 entry (schemas/tool-definition.v1.schema.json)
// read under a stricter profile: it must carry a closed parameter schema, and
// its binding.transport may hold placeholders only. A harness imports the
// templates it needs from a connector-contracts version pinned into its image,
// registers names and schemas from them, and templates the provider call from
// the same entries after an allow. The credential is injected by the harness
// after the allow and never appears in a template.
package calltemplate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Template is one validated tool call template.
type Template struct {
	Schema      string          `json:"schema"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  json.RawMessage `json:"parameters"`
	Permission  string          `json:"permission"`
	Category    string          `json:"category,omitempty"`
	Service     string          `json:"service,omitempty"`
	Version     string          `json:"version,omitempty"`
	Tags        []string        `json:"tags,omitempty"`
	Binding     *Binding        `json:"binding,omitempty"`
}

// Binding is the non-secret routing metadata of a template.
type Binding struct {
	ConnectorID      string          `json:"connector_id"`
	Config           json.RawMessage `json:"config,omitempty"`
	DelegatedContext []string        `json:"delegated_context,omitempty"`
	Transport        *Transport      `json:"transport,omitempty"`
}

// Transport describes the provider call. Absent for action tools the harness
// executes itself.
type Transport struct {
	Method      string         `json:"method"`
	URLTemplate string         `json:"url_template"`
	Auth        *Auth          `json:"auth,omitempty"`
	Args        map[string]Arg `json:"args,omitempty"`
}

// Auth is credential placement only; the value is injected after the allow.
type Auth struct {
	Placement string `json:"placement,omitempty"`
	Name      string `json:"name,omitempty"`
	Scheme    string `json:"scheme,omitempty"`
}

// Arg maps a tool parameter onto the provider call.
type Arg struct {
	In       string `json:"in"`
	Name     string `json:"name,omitempty"`
	Encoding string `json:"encoding,omitempty"`
}

var (
	namePattern        = regexp.MustCompile(`^[a-z0-9_]+(\.[a-z0-9_]+)*$`)
	paramNamePattern   = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)
	placeholderPattern = regexp.MustCompile(`\{([^{}]+)\}`)
	methods            = map[string]bool{"GET": true, "POST": true, "PUT": true, "PATCH": true, "DELETE": true}
	argLocations       = map[string]bool{"path": true, "query": true, "body": true, "header": true}
	placements         = map[string]bool{"header": true, "query": true, "none": true}
	encodings          = map[string]bool{"base64": true}
)

// LoadFile reads and validates a template file.
func LoadFile(path string) ([]Template, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	tpls, err := Parse(data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return tpls, nil
}

// Parse validates a template file body (a JSON array of templates) and returns
// the templates. Every problem found is reported, joined into one error.
func Parse(data []byte) ([]Template, error) {
	var rawList []json.RawMessage
	if err := json.Unmarshal(data, &rawList); err != nil {
		return nil, fmt.Errorf("templates must be a JSON array: %w", err)
	}
	if len(rawList) == 0 {
		return nil, errors.New("templates must contain at least one entry")
	}
	var (
		errs  []error
		tpls  []Template
		names = map[string]bool{}
	)
	for i, raw := range rawList {
		dec := json.NewDecoder(bytes.NewReader(raw))
		dec.DisallowUnknownFields()
		var t Template
		if err := dec.Decode(&t); err != nil {
			errs = append(errs, fmt.Errorf("[%d]: %w (unknown field or wrong type)", i, err))
			continue
		}
		label := fmt.Sprintf("[%d] %s", i, t.Name)
		for _, e := range validate(t, raw) {
			errs = append(errs, fmt.Errorf("%s: %w", label, e))
		}
		if names[t.Name] {
			errs = append(errs, fmt.Errorf("%s: duplicate tool name", label))
		}
		names[t.Name] = true
		tpls = append(tpls, t)
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return tpls, nil
}

func validate(t Template, raw json.RawMessage) []error {
	var errs []error
	add := func(format string, a ...any) { errs = append(errs, fmt.Errorf(format, a...)) }

	if t.Schema != "kei.tool-definition/v1" {
		add("schema must be kei.tool-definition/v1, got %q", t.Schema)
	}
	if !namePattern.MatchString(t.Name) {
		add("name %q must match %s", t.Name, namePattern)
	}
	if t.Description == "" {
		add("description is required")
	}
	if t.Permission == "" {
		add("permission is required")
	}
	errs = append(errs, validateParameters(t.Parameters)...)
	if t.Binding != nil {
		if t.Binding.ConnectorID == "" {
			add("binding.connector_id is required")
		}
		if t.Binding.Transport != nil {
			errs = append(errs, validateTransport(t.Binding.Transport)...)
		}
	}
	var tree any
	if err := json.Unmarshal(raw, &tree); err == nil {
		for _, msg := range findSecrets("", tree) {
			add("secret in template: %s", msg)
		}
	}
	return errs
}

// validateParameters requires a parameter schema: a list of named parameters,
// or a closed JSON Schema object.
func validateParameters(raw json.RawMessage) []error {
	if len(raw) == 0 || string(raw) == "null" {
		return []error{errors.New("parameters (the tool schema) is required")}
	}
	var list []struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(raw, &list); err == nil {
		var errs []error
		for _, p := range list {
			if !paramNamePattern.MatchString(p.Name) {
				errs = append(errs, fmt.Errorf("parameters: name %q must match %s", p.Name, paramNamePattern))
			}
		}
		return errs
	}
	var obj struct {
		Type                 string                     `json:"type"`
		Properties           map[string]json.RawMessage `json:"properties"`
		Required             []string                   `json:"required"`
		AdditionalProperties *bool                      `json:"additionalProperties"`
	}
	if err := json.Unmarshal(raw, &obj); err != nil {
		return []error{fmt.Errorf("parameters must be a list or a JSON Schema object: %w", err)}
	}
	var errs []error
	if obj.Type != "object" {
		errs = append(errs, errors.New("parameters schema must have type object"))
	}
	if obj.AdditionalProperties == nil || *obj.AdditionalProperties {
		errs = append(errs, errors.New("parameters schema must set additionalProperties: false"))
	}
	for _, r := range obj.Required {
		if _, ok := obj.Properties[r]; !ok {
			errs = append(errs, fmt.Errorf("parameters: required %q is not in properties", r))
		}
	}
	return errs
}

func validateTransport(tr *Transport) []error {
	var errs []error
	add := func(format string, a ...any) { errs = append(errs, fmt.Errorf(format, a...)) }
	if !methods[tr.Method] {
		add("transport.method %q must be one of GET, POST, PUT, PATCH, DELETE", tr.Method)
	}
	if tr.URLTemplate == "" {
		add("transport.url_template is required")
	}
	if tr.Auth != nil {
		if tr.Auth.Placement != "" && !placements[tr.Auth.Placement] {
			add("transport.auth.placement %q must be header, query or none", tr.Auth.Placement)
		}
		if tr.Auth.Placement != "" && tr.Auth.Placement != "none" && tr.Auth.Name == "" {
			add("transport.auth.name is required for placement %q", tr.Auth.Placement)
		}
	}
	for name, a := range tr.Args {
		if !argLocations[a.In] {
			add("transport.args.%s.in %q must be path, query, body or header", name, a.In)
		}
		if a.Encoding != "" && !encodings[a.Encoding] {
			add("transport.args.%s.encoding %q must be base64", name, a.Encoding)
		}
	}
	for _, m := range placeholderPattern.FindAllStringSubmatch(tr.URLTemplate, -1) {
		if a, ok := tr.Args[m[1]]; !ok || a.In != "path" {
			add("url_template placeholder {%s} needs a transport.args entry with in: path", m[1])
		}
	}
	return errs
}

var (
	secretValuePatterns = []*regexp.Regexp{
		regexp.MustCompile(`gh[pousr]_[A-Za-z0-9]{16,}`),
		regexp.MustCompile(`github_pat_[A-Za-z0-9_]{16,}`),
		regexp.MustCompile(`\bsk-[A-Za-z0-9_-]{16,}`),
		regexp.MustCompile(`xox[abprs]-[A-Za-z0-9-]{8,}`),
		regexp.MustCompile(`AKIA[0-9A-Z]{16}`),
		regexp.MustCompile(`eyJ[A-Za-z0-9_-]{8,}\.[A-Za-z0-9_-]{8,}\.`),
		regexp.MustCompile(`-----BEGIN [A-Z ]*PRIVATE KEY`),
		regexp.MustCompile(`(?i)\b(bearer|basic)\s+[A-Za-z0-9._~+/=-]{8,}`),
		regexp.MustCompile(`://[^/@\s{}]+:[^/@\s{}]*@`),
		regexp.MustCompile(`(?i)[?&](access_token|token|api[_-]?key|key|secret|password)=[^&{\s]`),
	}
	secretKeyPattern = regexp.MustCompile(`(?i)(secret|token|password|passwd|api[_-]?key|credential|private[_-]?key)`)
)

// findSecrets walks a decoded template and reports values that look like
// credentials, and secret-named keys that hold a literal value. Placeholders
// ({name}) are fine.
func findSecrets(path string, v any) []string {
	var out []string
	switch x := v.(type) {
	case map[string]any:
		for k, child := range x {
			p := path + "." + k
			if s, ok := child.(string); ok && secretKeyPattern.MatchString(k) && s != "" && !isPlaceholder(s) {
				out = append(out, fmt.Sprintf("%s holds a literal value (credentials are injected by the harness after the allow)", p))
			}
			out = append(out, findSecrets(p, child)...)
		}
	case []any:
		for i, child := range x {
			out = append(out, findSecrets(fmt.Sprintf("%s[%d]", path, i), child)...)
		}
	case string:
		for _, re := range secretValuePatterns {
			if re.MatchString(x) {
				out = append(out, fmt.Sprintf("%s matches secret pattern %s", path, re))
				break
			}
		}
	}
	return out
}

func isPlaceholder(s string) bool {
	return strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}")
}
