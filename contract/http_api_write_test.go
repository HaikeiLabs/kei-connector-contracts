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

package contract

import (
	"encoding/json"
	"strings"
	"testing"
)

// writeHTTPFixture is httpFixture with POST and PATCH configured.
func writeHTTPFixture() Metadata {
	m := httpFixture()
	m.HTTPAPI.AllowedMethods = []HTTPMethod{HTTPMethodGet, HTTPMethodHead, HTTPMethodPost, HTTPMethodPatch}
	return m
}

func httpWriteCall(capability string, action Action, method HTTPMethod, body string) Invocation {
	in := validCall("c-http", capability, "/contacts/42")
	in.Action = action
	in.HTTP = &HTTPInvocation{Method: method, Path: "/contacts/42", ContentType: HTTPContentTypeJSON, Body: json.RawMessage(body)}
	return in
}

// TestHTTPAPIWriteCapabilitiesAreWrites: http.post and http.patch are classified
// like every other provider's writes, create and update, never read.
func TestHTTPAPIWriteCapabilitiesAreWrites(t *testing.T) {
	for method, want := range map[HTTPMethod]Action{HTTPMethodGet: ActionRead, HTTPMethodHead: ActionRead, HTTPMethodPost: ActionCreate, HTTPMethodPatch: ActionUpdate} {
		if got, ok := HTTPMethodAction(method); !ok || got != want {
			t.Errorf("%s action = %q, %v; want %q", method, got, ok, want)
		}
	}
	for _, method := range []HTTPMethod{"DELETE", "PUT", "OPTIONS", "post"} {
		if _, ok := HTTPMethodAction(method); ok {
			t.Errorf("%s is supported", method)
		}
	}

	// A write method cannot be declared as a read, so policy on the read
	// action can never admit it.
	for _, c := range []Capability{{Name: "http.post", Action: ActionRead}, {Name: "http.patch", Action: ActionCreate}, {Name: "http.get", Action: ActionUpdate}} {
		m := writeHTTPFixture()
		m.Capabilities = []Capability{c}
		if err := m.Validate(); err == nil || !strings.Contains(err.Error(), "must declare action") {
			t.Errorf("%s as %s: error = %v", c.Name, c.Action, err)
		}
	}
}

func TestHTTPAPIValidateAllowsWriteMethods(t *testing.T) {
	m := writeHTTPFixture()
	if err := m.Validate(); err != nil {
		t.Fatalf("POST and PATCH rejected: %v", err)
	}
	m.HTTPAPI.AllowedMethods = []HTTPMethod{HTTPMethodPost}
	if err := m.Validate(); err != nil {
		t.Fatalf("POST-only connector rejected: %v", err)
	}
	for _, method := range []HTTPMethod{"DELETE", "PUT", "patch", ""} {
		m.HTTPAPI.AllowedMethods = []HTTPMethod{HTTPMethodGet, method}
		if err := m.Validate(); err == nil {
			t.Errorf("method %q accepted", method)
		}
	}
}

func TestHTTPAPIWriteCallAcceptsJSONObjectBody(t *testing.T) {
	m := writeHTTPFixture()
	for _, in := range []Invocation{
		httpWriteCall("http.post", ActionCreate, HTTPMethodPost, `{"name":"Ada"}`),
		httpWriteCall("http.patch", ActionUpdate, HTTPMethodPatch, ` {"stage":"won","tags":["a"]} `),
		httpWriteCall("http.post", ActionCreate, HTTPMethodPost, `{}`),
	} {
		if err := ValidateCall(m, in); err != nil {
			t.Errorf("%s %s rejected: %v", in.HTTP.Method, in.HTTP.Body, err)
		}
	}
}

func TestHTTPAPIWriteCallRejectsInvalidBody(t *testing.T) {
	m := writeHTTPFixture()
	m.HTTPAPI.RequestMaxBytes = 32
	for name, tc := range map[string]struct {
		mutate func(*HTTPInvocation)
		want   string
	}{
		"missing body":        {func(h *HTTPInvocation) { h.Body = nil }, "HTTP request body must be a JSON object"},
		"array body":          {func(h *HTTPInvocation) { h.Body = json.RawMessage(`[1]`) }, "HTTP request body must be a JSON object"},
		"string body":         {func(h *HTTPInvocation) { h.Body = json.RawMessage(`"x"`) }, "HTTP request body must be a JSON object"},
		"null body":           {func(h *HTTPInvocation) { h.Body = json.RawMessage(`null`) }, "HTTP request body must be a JSON object"},
		"malformed body":      {func(h *HTTPInvocation) { h.Body = json.RawMessage(`{"a":`) }, "HTTP request body must be a JSON object"},
		"trailing data":       {func(h *HTTPInvocation) { h.Body = json.RawMessage(`{} {}`) }, "HTTP request body must be a JSON object"},
		"oversized body":      {func(h *HTTPInvocation) { h.Body = json.RawMessage(`{"a":"` + strings.Repeat("x", 32) + `"}`) }, "HTTP request body exceeds request_max_bytes"},
		"no content type":     {func(h *HTTPInvocation) { h.ContentType = "" }, "HTTP request content type must be application/json"},
		"form content type":   {func(h *HTTPInvocation) { h.ContentType = "application/x-www-form-urlencoded" }, "HTTP request content type must be application/json"},
		"text content type":   {func(h *HTTPInvocation) { h.ContentType = "text/plain" }, "HTTP request content type must be application/json"},
		"path outside":        {func(h *HTTPInvocation) { h.Path = "/admin" }, "HTTP path is not allowed"},
		"query not allowlist": {func(h *HTTPInvocation) { h.Query = map[string][]string{"secret": {"x"}} }, `query parameter "secret" is not allowed`},
	} {
		t.Run(name, func(t *testing.T) {
			in := httpWriteCall("http.post", ActionCreate, HTTPMethodPost, `{"a":1}`)
			tc.mutate(in.HTTP)
			if err := ValidateCall(m, in); err == nil || err.Error() != tc.want {
				t.Fatalf("error = %v, want %q", err, tc.want)
			}
		})
	}
}

// TestHTTPAPIReadsStayBodyless: GET and HEAD behave as before and take no body.
func TestHTTPAPIReadsStayBodyless(t *testing.T) {
	m := writeHTTPFixture()
	in := validCall("c-http", "http.get", "/contacts/42")
	in.HTTP = &HTTPInvocation{Method: HTTPMethodGet, Path: "/contacts/42"}
	if err := ValidateCall(m, in); err != nil {
		t.Fatalf("GET rejected: %v", err)
	}
	for _, mutate := range []func(*HTTPInvocation){
		func(h *HTTPInvocation) { h.Body = json.RawMessage(`{}`) },
		func(h *HTTPInvocation) { h.ContentType = HTTPContentTypeJSON },
	} {
		get := in
		h := *in.HTTP
		mutate(&h)
		get.HTTP = &h
		if err := ValidateCall(m, get); err == nil || err.Error() != "HTTP method does not take a request body" {
			t.Errorf("GET with body error = %v", err)
		}
	}
}

// TestHTTPAPIRequestMethodMatchesCapability: a request's method is the one its
// capability names, so a read capability cannot carry a write request.
func TestHTTPAPIRequestMethodMatchesCapability(t *testing.T) {
	m := writeHTTPFixture()
	for _, in := range []Invocation{
		httpWriteCall("http.get", ActionRead, HTTPMethodPost, `{}`),
		httpWriteCall("http.post", ActionCreate, HTTPMethodPatch, `{}`),
	} {
		if err := ValidateCall(m, in); err == nil || err.Error() != "HTTP method does not match capability" {
			t.Errorf("%s via %s error = %v", in.HTTP.Method, in.Capability, err)
		}
	}
}

// TestHTTPAPIWritesAreOptIn: a write is reachable only when the connector both
// declares the write capability and configures the method.
func TestHTTPAPIWritesAreOptIn(t *testing.T) {
	post := httpWriteCall("http.post", ActionCreate, HTTPMethodPost, `{"a":1}`)

	readOnly := writeHTTPFixture()
	readOnly.Capabilities = []Capability{{Name: "http.get", Action: ActionRead}}
	if err := ValidateCall(readOnly, post); err == nil || err.Error() != "capability is not allowed" {
		t.Errorf("undeclared write error = %v", err)
	}

	unconfigured := httpFixture() // declares http.post; allows only GET and HEAD
	if err := ValidateCall(unconfigured, post); err == nil || err.Error() != "HTTP method is not allowed" {
		t.Errorf("unconfigured method error = %v", err)
	}

	wrongAction := post
	wrongAction.Action = ActionRead
	if err := ValidateCall(writeHTTPFixture(), wrongAction); err == nil || err.Error() != "capability action mismatch" {
		t.Errorf("write invoked as read error = %v", err)
	}
}

// TestAuthorizeCallAuthorizesHTTPWritesWithoutBody: the control plane never
// sees the request body, so it authorizes a write on the capability's method
// and the resource path alone.
func TestAuthorizeCallAuthorizesHTTPWritesWithoutBody(t *testing.T) {
	m := writeHTTPFixture()
	for _, tc := range []struct {
		capability string
		action     Action
	}{{"http.post", ActionCreate}, {"http.patch", ActionUpdate}} {
		in := validCall("c-http", tc.capability, "/contacts/42")
		in.Action = tc.action
		if err := AuthorizeCall(m, in); err != nil {
			t.Errorf("%s rejected: %v", tc.capability, err)
		}
		in.Resource = "/admin"
		if err := AuthorizeCall(m, in); err == nil || err.Error() != "HTTP path is not allowed" {
			t.Errorf("%s outside patterns error = %v", tc.capability, err)
		}
	}

	in := validCall("c-http", "http.post", "/contacts/42")
	in.Action = ActionCreate
	if err := AuthorizeCall(httpFixture(), in); err == nil || err.Error() != "HTTP method is not allowed" {
		t.Errorf("unconfigured POST error = %v", err)
	}
}
