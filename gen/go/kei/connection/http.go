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

package connection

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	contracts "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/contracts/v1"
)

// DefaultHeaderName is the header an HTTP_HEADER connection puts the
// credential in when the connection declares no header name.
const DefaultHeaderName = "X-Api-Key"

// MaxResponseBytes bounds the response body Do reads.
const MaxResponseBytes = 32 << 20

// HTTP is an authenticated HTTP connection to one provider origin. The
// credential is already attached; callers pass a rendered, escaped path.
type HTTP interface {
	// Do sends one request. path must be an absolute path (from RenderPath)
	// with no query or fragment; it is appended to the base URL's path. A
	// non-2xx status is not an error: map it with MapError.
	Do(ctx context.Context, method, path string, query url.Values, body []byte) (status int, header http.Header, respBody []byte, err error)
}

// NewHTTP builds an HTTP connection for conn's kind:
//
//   - CONNECTION_KIND_HTTP_BEARER: Authorization: Bearer <credential>;
//   - CONNECTION_KIND_HTTP_HEADER: the credential in the connection's
//     header_name (DefaultHeaderName when it is empty);
//   - CONNECTION_KIND_HTTP_BASIC: the credential is "user:password".
//
// baseURL must be an absolute http or https URL with no user info, query or
// fragment. client is required: the connection makes no request except
// through it.
func NewHTTP(baseURL string, conn *contracts.Connection, credential string, client *http.Client) (HTTP, error) {
	if conn == nil {
		return nil, errors.New("connection: NewHTTP: nil connection")
	}
	if client == nil {
		return nil, errors.New("connection: NewHTTP: nil *http.Client")
	}
	base, err := parseBase(baseURL)
	if err != nil {
		return nil, err
	}
	if credential == "" {
		return nil, fmt.Errorf("connection: NewHTTP %q: empty credential", conn.GetName())
	}
	if strings.ContainsAny(credential, "\r\n\x00") {
		return nil, fmt.Errorf("connection: NewHTTP %q: credential has a control character", conn.GetName())
	}
	h := &httpConn{base: base, client: client, name: conn.GetName(), kind: conn.GetKind()}
	switch conn.GetKind() {
	case contracts.ConnectionKind_CONNECTION_KIND_HTTP_BEARER:
		h.header, h.value = "Authorization", "Bearer "+credential
	case contracts.ConnectionKind_CONNECTION_KIND_HTTP_HEADER:
		name := conn.GetHeaderName()
		if name == "" {
			name = DefaultHeaderName
		}
		if !validHeaderName(name) {
			return nil, fmt.Errorf("connection: NewHTTP %q: invalid header name %q", conn.GetName(), name)
		}
		h.header, h.value = name, credential
	case contracts.ConnectionKind_CONNECTION_KIND_HTTP_BASIC:
		user, pass, ok := strings.Cut(credential, ":")
		if !ok || user == "" {
			return nil, fmt.Errorf("connection: NewHTTP %q: a basic credential must be user:password", conn.GetName())
		}
		h.basicUser, h.basicPass = user, pass
	default:
		return nil, fmt.Errorf("connection: NewHTTP %q: kind %s is not an HTTP connection", conn.GetName(), conn.GetKind())
	}
	return h, nil
}

// validHeaderName reports whether name is an RFC 9110 token.
func validHeaderName(name string) bool {
	if name == "" {
		return false
	}
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
		case strings.ContainsRune("!#$%&'*+-.^_`|~", r):
		default:
			return false
		}
	}
	return true
}

func parseBase(raw string) (*url.URL, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("connection: invalid base URL: %w", err)
	}
	switch {
	case u.Scheme != "http" && u.Scheme != "https":
		return nil, fmt.Errorf("connection: base URL %q must be http or https", u.Redacted())
	case u.Host == "":
		return nil, fmt.Errorf("connection: base URL %q has no host", u.Redacted())
	case u.User != nil:
		return nil, errors.New("connection: base URL must not carry user info")
	case u.RawQuery != "" || u.ForceQuery || u.Fragment != "":
		return nil, fmt.Errorf("connection: base URL %q must not have a query or fragment", u.Redacted())
	}
	return u, nil
}

// httpConn holds the credential in unexported fields; String, GoString and
// Format never print them.
type httpConn struct {
	base   *url.URL
	client *http.Client
	name   string
	kind   contracts.ConnectionKind

	header, value        string
	basicUser, basicPass string
}

func (h *httpConn) String() string {
	return fmt.Sprintf("connection.HTTP{name=%s kind=%s base=%s}", h.name, h.kind, h.base.Redacted())
}

func (h *httpConn) GoString() string { return h.String() }

// Format makes every verb (%v, %+v, %#v, %s, %q) print String.
func (h *httpConn) Format(f fmt.State, verb rune) {
	if verb == 'q' {
		fmt.Fprintf(f, "%q", h.String())
		return
	}
	_, _ = io.WriteString(f, h.String())
}

func (h *httpConn) Do(ctx context.Context, method, path string, query url.Values, body []byte) (int, http.Header, []byte, error) {
	target, err := h.resolve(path)
	if err != nil {
		return 0, nil, nil, err
	}
	target.RawQuery = query.Encode()
	var rd io.Reader
	if body != nil {
		rd = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, target.String(), rd)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("connection: %s %s: %w", method, path, err)
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if h.header != "" {
		req.Header.Set(h.header, h.value)
	} else {
		req.SetBasicAuth(h.basicUser, h.basicPass)
	}
	resp, err := h.client.Do(req)
	if err != nil {
		// *url.Error carries the URL, which never holds a credential: the
		// base URL has no user info and the credential is only in headers.
		return 0, nil, nil, fmt.Errorf("connection: %s %s: %w", method, path, err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(io.LimitReader(resp.Body, MaxResponseBytes+1))
	if err != nil {
		return 0, nil, nil, fmt.Errorf("connection: %s %s: reading response: %w", method, path, err)
	}
	if len(respBody) > MaxResponseBytes {
		return 0, nil, nil, fmt.Errorf("connection: %s %s: response larger than %d bytes", method, path, MaxResponseBytes)
	}
	return resp.StatusCode, resp.Header, respBody, nil
}

// resolve appends an escaped absolute path to the base URL's path and
// refuses anything that could change the origin.
func (h *httpConn) resolve(path string) (*url.URL, error) {
	if !strings.HasPrefix(path, "/") || strings.HasPrefix(path, "//") || strings.ContainsAny(path, "?#") {
		return nil, fmt.Errorf("connection: path %q must be an absolute path with no query or fragment", path)
	}
	basePath := strings.TrimSuffix(h.base.EscapedPath(), "/")
	u, err := url.Parse(h.base.Scheme + "://" + h.base.Host + basePath + path)
	if err != nil {
		return nil, fmt.Errorf("connection: invalid path %q: %w", path, err)
	}
	if u.Host != h.base.Host || u.Scheme != h.base.Scheme {
		return nil, fmt.Errorf("connection: path %q changes the origin", path)
	}
	return u, nil
}
