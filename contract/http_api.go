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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"
)

type HTTPMethod string

const (
	HTTPMethodGet  HTTPMethod = "GET"
	HTTPMethodHead HTTPMethod = "HEAD"
	// HTTPMethodPost and HTTPMethodPatch are writes. They carry a JSON
	// object body (HTTPInvocation.Body) sent as application/json.
	HTTPMethodPost  HTTPMethod = "POST"
	HTTPMethodPatch HTTPMethod = "PATCH"
)

// HTTPContentTypeJSON is the only request content type http_api sends.
const HTTPContentTypeJSON = "application/json"

// httpMethodActions classifies each supported method the way other providers
// classify their capabilities: GET and HEAD read, POST creates, PATCH updates.
// A method absent here is not supported.
var httpMethodActions = map[HTTPMethod]Action{
	HTTPMethodGet:   ActionRead,
	HTTPMethodHead:  ActionRead,
	HTTPMethodPost:  ActionCreate,
	HTTPMethodPatch: ActionUpdate,
}

// HTTPMethodAction returns the action a supported HTTP method performs, and
// whether the method is supported.
func HTTPMethodAction(method HTTPMethod) (Action, bool) {
	action, ok := httpMethodActions[method]
	return action, ok
}

// httpMethodHasBody reports whether a supported method carries a request body.
func httpMethodHasBody(method HTTPMethod) bool {
	return method == HTTPMethodPost || method == HTTPMethodPatch
}

type CredentialInjectionLocation string

const CredentialInjectionHeader CredentialInjectionLocation = "header"

// HTTPAPI is destination and transport policy only. CredentialRef remains an
// opaque reference in Metadata; no credential material is representable here.
type HTTPAPI struct {
	RegisteredBaseURL       string              `json:"registered_base_url"`
	RegisteredHostname      string              `json:"registered_hostname"`
	AllowedMethods          []HTTPMethod        `json:"allowed_methods"`
	AllowedPathPatterns     []string            `json:"allowed_path_patterns"`
	QueryAllowlist          []string            `json:"query_allowlist"`
	RequestMaxBytes         int64               `json:"request_max_bytes"`
	ResponseMaxBytes        int64               `json:"response_max_bytes"`
	TimeoutMS               int64               `json:"timeout_ms"`
	AllowRedirects          bool                `json:"allow_redirects"`
	MaxRedirects            int                 `json:"max_redirects"`
	AllowCrossHostRedirects bool                `json:"allow_cross_host_redirects"`
	CredentialInjection     CredentialInjection `json:"credential_injection"`
	Redaction               RedactionRules      `json:"redaction"`
}

type CredentialInjection struct {
	Location CredentialInjectionLocation `json:"location"`
	Name     string                      `json:"name"`
	Scheme   string                      `json:"scheme,omitempty"`
}
type RedactionRules struct {
	Headers         []string `json:"headers"`
	QueryParameters []string `json:"query_parameters"`
	ResponseFields  []string `json:"response_fields"`
}

// HTTPInvocation is the full HTTP request the tenant runtime executes. It stays
// in the runtime (ADR-011): the control plane authorizes on the capability's
// method and the resource path (AuthorizeCall) and never needs Body.
type HTTPInvocation struct {
	Method HTTPMethod          `json:"method"`
	Path   string              `json:"path"`
	Query  map[string][]string `json:"query,omitempty"`
	// ContentType and Body are set exactly when Method is POST or PATCH:
	// ContentType is application/json and Body is a JSON object no larger
	// than the connector's RequestMaxBytes.
	ContentType string          `json:"content_type,omitempty"`
	Body        json.RawMessage `json:"body,omitempty"`
}

const (
	maxHTTPBodyBytes int64 = 10 << 20
	maxHTTPTimeoutMS int64 = 30_000
	maxHTTPRedirects       = 3
)

var safeHTTPName = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,127}$`)

func (h HTTPAPI) Validate() error {
	if containsSecretMaterial(h.RegisteredBaseURL) || containsSecretMaterial(h.RegisteredHostname) || containsSecretMaterial(h.CredentialInjection.Scheme) {
		return errors.New("metadata must not contain secret material")
	}
	u, err := url.Parse(h.RegisteredBaseURL)
	if err != nil || u.Scheme != "https" || u.User != nil || u.Host == "" || u.RawQuery != "" || u.Fragment != "" {
		return errors.New("registered_base_url must be an https URL without credentials, query, or fragment")
	}
	if !validPublicHostname(u.Hostname()) || !strings.EqualFold(u.Hostname(), h.RegisteredHostname) {
		return errors.New("registered hostname is invalid or does not match registered base URL")
	}
	if len(h.AllowedMethods) == 0 {
		return errors.New("http_api must allow at least one of GET, HEAD, POST, or PATCH")
	}
	for _, method := range h.AllowedMethods {
		if _, ok := HTTPMethodAction(method); !ok {
			return fmt.Errorf("http method %q is not allowed", method)
		}
	}
	if len(h.AllowedPathPatterns) == 0 {
		return errors.New("allowed_path_patterns is required")
	}
	for _, pattern := range h.AllowedPathPatterns {
		if err := validateHTTPPathPattern(pattern); err != nil {
			return err
		}
	}
	for _, name := range h.QueryAllowlist {
		if !safeHTTPName.MatchString(name) {
			return fmt.Errorf("query parameter %q is invalid", name)
		}
	}
	if h.RequestMaxBytes <= 0 || h.RequestMaxBytes > maxHTTPBodyBytes || h.ResponseMaxBytes <= 0 || h.ResponseMaxBytes > maxHTTPBodyBytes {
		return errors.New("HTTP request and response limits are invalid")
	}
	if h.TimeoutMS <= 0 || h.TimeoutMS > maxHTTPTimeoutMS {
		return errors.New("HTTP timeout is invalid")
	}
	if h.MaxRedirects < 0 || h.MaxRedirects > maxHTTPRedirects || (!h.AllowRedirects && h.MaxRedirects != 0) {
		return errors.New("HTTP redirect policy is invalid")
	}
	if h.AllowCrossHostRedirects {
		return errors.New("cross-host redirects are not permitted")
	}
	if h.CredentialInjection.Location != CredentialInjectionHeader || !safeHTTPName.MatchString(h.CredentialInjection.Name) {
		return errors.New("credential injection must use a valid header name")
	}
	for _, name := range append(append([]string{}, h.Redaction.Headers...), h.Redaction.QueryParameters...) {
		if !safeHTTPName.MatchString(name) {
			return fmt.Errorf("redaction name %q is invalid", name)
		}
	}
	for _, field := range h.Redaction.ResponseFields {
		if field == "" || strings.ContainsAny(field, "\r\n") {
			return errors.New("response redaction field is invalid")
		}
	}
	if !containsFold(h.Redaction.Headers, h.CredentialInjection.Name) {
		return errors.New("credential header must be redacted")
	}
	return nil
}

// httpCapabilityMethod reads the HTTP method an http_api capability performs
// from its name (http.get -> GET), so no capability table is needed. Whether
// that method is allowed is the connector's configured AllowedMethods.
func httpCapabilityMethod(capability string) (HTTPMethod, bool) {
	verb, ok := strings.CutPrefix(capability, "http.")
	if !ok || verb == "" {
		return "", false
	}
	return HTTPMethod(strings.ToUpper(verb)), true
}

// ValidateInvocation checks the full request the runtime executes: the
// request target (validateTarget) and, for POST and PATCH, a JSON object body
// within RequestMaxBytes sent as application/json. GET and HEAD carry no body.
func (h HTTPAPI) ValidateInvocation(in HTTPInvocation) error {
	if err := h.validateTarget(in); err != nil {
		return err
	}
	if !httpMethodHasBody(in.Method) {
		if in.ContentType != "" || len(in.Body) != 0 {
			return errors.New("HTTP method does not take a request body")
		}
		return nil
	}
	if in.ContentType != HTTPContentTypeJSON {
		return errors.New("HTTP request content type must be application/json")
	}
	if int64(len(in.Body)) > h.RequestMaxBytes {
		return errors.New("HTTP request body exceeds request_max_bytes")
	}
	if !json.Valid(in.Body) || !bytes.HasPrefix(bytes.TrimSpace(in.Body), []byte("{")) {
		return errors.New("HTTP request body must be a JSON object")
	}
	return nil
}

// validateTarget checks the method, path, and query against the connector's
// configuration. It is what the control plane can check without the body.
func (h HTTPAPI) validateTarget(in HTTPInvocation) error {
	if _, ok := HTTPMethodAction(in.Method); !ok || !containsHTTPMethod(h.AllowedMethods, in.Method) {
		return errors.New("HTTP method is not allowed")
	}
	if err := validateRelativeHTTPPath(in.Path); err != nil {
		return err
	}
	matched := false
	for _, pattern := range h.AllowedPathPatterns {
		if pathMatches(pattern, in.Path) {
			matched = true
			break
		}
	}
	if !matched {
		return errors.New("HTTP path is not allowed")
	}
	for name := range in.Query {
		if !containsFold(h.QueryAllowlist, name) {
			return fmt.Errorf("query parameter %q is not allowed", name)
		}
	}
	return nil
}

func validPublicHostname(host string) bool {
	host = strings.TrimSuffix(strings.ToLower(host), ".")
	if host == "" || host == "localhost" || host == "metadata" || strings.HasPrefix(host, "metadata.") || strings.HasSuffix(host, ".localhost") || strings.HasSuffix(host, ".local") || strings.HasSuffix(host, ".internal") {
		return false
	}
	if ip := net.ParseIP(host); ip != nil {
		return !ip.IsLoopback() && !ip.IsPrivate() && !ip.IsLinkLocalUnicast() && !ip.IsUnspecified()
	}
	return !strings.Contains(host, "..") && !strings.HasPrefix(host, ".")
}
func validateHTTPPathPattern(pattern string) error {
	if err := validateRelativeHTTPPath(pattern); err != nil {
		return err
	}
	if strings.ContainsAny(pattern, "?#") {
		return errors.New("HTTP path pattern must not contain query or fragment")
	}
	return nil
}
func validateRelativeHTTPPath(path string) error {
	if path == "" || !strings.HasPrefix(path, "/") || strings.HasPrefix(path, "//") || strings.ContainsAny(path, "?#[\\\r\n") {
		return errors.New("HTTP path must be relative")
	}
	for _, part := range strings.Split(path, "/") {
		decoded, err := url.PathUnescape(part)
		if err != nil || strings.ContainsAny(decoded, "/\\") {
			return errors.New("HTTP path contains an invalid escape")
		}
		if decoded == ".." || decoded == "." {
			return errors.New("HTTP path traversal is not allowed")
		}
	}
	u, err := url.Parse(path)
	if err != nil || u.IsAbs() || u.Host != "" {
		return errors.New("HTTP path must be relative")
	}
	return nil
}
func pathMatches(pattern, path string) bool {
	prefix := strings.TrimSuffix(pattern, "*")
	return path == prefix || (strings.HasSuffix(pattern, "*") && strings.HasPrefix(path, prefix))
}
func containsHTTPMethod(values []HTTPMethod, want HTTPMethod) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
func containsFold(values []string, want string) bool {
	for _, value := range values {
		if strings.EqualFold(value, want) {
			return true
		}
	}
	return false
}
