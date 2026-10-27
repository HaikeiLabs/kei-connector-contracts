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

// Package authorizerequest defines the wire contract of the
// kei.authorize-request/v1 request a harness sends to the tenant-local
// kei-proxy daemon (POST /v1/authorize). It carries metadata only: the tool or
// connector capability being called, a resource name and a digest of the
// arguments. Tool arguments, message content and credentials never appear here.
//
// The optional Caller object (harness kind, origin and an opaque session
// reference) lets the audit trail say who called. A request without it is
// valid and its caller is unknown, so existing senders keep working.
package authorizerequest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
)

// Schema is the value of the request's schema field.
const Schema = "kei.authorize-request/v1"

// Caller origins.
const (
	OriginInteractive = "interactive"
	OriginAutomation  = "automation"
)

var (
	harnessPattern    = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,63}$`)
	sessionRefPattern = regexp.MustCompile(`^[A-Za-z0-9._:/-]{1,128}$`)
)

// Request is a kei.authorize-request/v1 object. Exactly one of Tool and
// Connector is set.
type Request struct {
	Schema          string     `json:"schema"`
	Tool            string     `json:"tool,omitempty"`
	Action          string     `json:"action,omitempty"`
	Resource        string     `json:"resource,omitempty"`
	AgentID         string     `json:"agent_id,omitempty"`
	Connector       *Connector `json:"connector,omitempty"`
	SpanID          string     `json:"span_id,omitempty"`
	ParentSpan      string     `json:"parent_span,omitempty"`
	DelegationDepth int        `json:"delegation_depth,omitempty"`
	Framework       string     `json:"framework,omitempty"`
	ToolArgsDigest  string     `json:"tool_args_digest,omitempty"`
	OnBehalfOf      string     `json:"on_behalf_of,omitempty"`
	// ShellCommand is accepted only over the tenant-local socket and is used
	// solely to evaluate shell: policies. It is never copied to audit.
	ShellCommand string `json:"shell_command,omitempty"`
	// Caller is optional caller metadata for the audit trail. Nil means the
	// sender did not say.
	Caller *Caller `json:"caller,omitempty"`
}

// Connector names a connector capability.
type Connector struct {
	ConnectorID string `json:"connector_id"`
	Capability  string `json:"capability"`
}

// Caller says who made the call. All fields are optional; nil means not sent. Harness is a short
// lowercase label (for example openwebui), Origin is interactive for a person
// in a live session or automation for an unattended run, and SessionRef is an
// opaque session or chat reference chosen by the sender. SessionRef is a
// pointer only: it must never carry message content.
type Caller struct {
	Harness    *string `json:"harness,omitempty"`
	Origin     *string `json:"origin,omitempty"`
	SessionRef *string `json:"session_ref,omitempty"`
}

// Validate checks the caller values. Errors never echo the submitted value.
func (c Caller) Validate() error {
	if c.Harness != nil && !harnessPattern.MatchString(*c.Harness) {
		return errors.New("caller.harness must be a short lowercase token")
	}
	if c.Origin != nil && *c.Origin != OriginInteractive && *c.Origin != OriginAutomation {
		return errors.New("caller.origin must be interactive or automation")
	}
	if c.SessionRef != nil && !sessionRefPattern.MatchString(*c.SessionRef) {
		return errors.New("caller.session_ref must be a short opaque token")
	}
	return nil
}

// Validate checks the request shape.
func (r Request) Validate() error {
	if r.Schema != Schema {
		return fmt.Errorf("schema must be %s", Schema)
	}
	if (r.Tool == "") == (r.Connector == nil) {
		return errors.New("exactly one of tool or connector is required")
	}
	if r.Connector != nil && (r.Connector.ConnectorID == "" || r.Connector.Capability == "") {
		return errors.New("connector needs connector_id and capability")
	}
	if r.Caller != nil {
		return r.Caller.Validate()
	}
	return nil
}

// Decode parses and validates a request. Unknown fields and trailing data are
// rejected, matching the daemon.
func Decode(raw []byte) (Request, error) {
	var req Request
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		return Request{}, errors.New("request is not a valid kei.authorize-request/v1 object")
	}
	if _, err := dec.Token(); !errors.Is(err, io.EOF) {
		return Request{}, errors.New("trailing data after the request object")
	}
	if err := req.Validate(); err != nil {
		return Request{}, err
	}
	return req, nil
}
