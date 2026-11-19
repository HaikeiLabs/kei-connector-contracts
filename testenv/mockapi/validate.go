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

package mockapi

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connection"
	contracts "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/contracts/v1"
)

// validateCredential checks that the request carries the credential where
// the connector's connections declare it. A request that places it anywhere
// else — or nowhere — is a 400 with a reason, because a runtime that places
// credentials wrong is exactly the bug the mock exists to catch.
func validateCredential(route *route, r *http.Request) string {
	if route.connector == nil {
		return ""
	}
	var accepted []string
	for _, conn := range route.connector.GetConnections() {
		switch conn.GetKind() {
		case contracts.ConnectionKind_CONNECTION_KIND_HTTP_BEARER, contracts.ConnectionKind_CONNECTION_KIND_TOKEN_EXCHANGE:
			if strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") && len(r.Header.Get("Authorization")) > len("Bearer ") {
				return ""
			}
			accepted = append(accepted, "Authorization: Bearer <credential>")
		case contracts.ConnectionKind_CONNECTION_KIND_HTTP_HEADER:
			name := conn.GetHeaderName()
			if name == "" {
				continue
			}
			if r.Header.Get(name) != "" {
				return ""
			}
			accepted = append(accepted, "the "+name+" header")
		case contracts.ConnectionKind_CONNECTION_KIND_HTTP_BASIC:
			if strings.HasPrefix(r.Header.Get("Authorization"), "Basic ") && len(r.Header.Get("Authorization")) > len("Basic ") {
				return ""
			}
			accepted = append(accepted, "Authorization: Basic <base64 user:credential>")
		default:
			// SQL DSN, key-pair JWT and other non-HTTP placements: the
			// mock cannot check them on an HTTP request.
		}
	}
	if len(accepted) == 0 {
		return ""
	}
	return "mockapi: credential placement wrong: want " + strings.Join(accepted, " or ")
}

// validatePath checks that the incoming path's variable values are escaped
// the way the connection runtime escapes them. The mock re-renders the path
// from the decoded values; any difference means the client skipped the
// escaping rules (raw reserved characters, a raw "/" inside a single
// segment), which real providers would misroute.
func validatePath(route *route, r *http.Request) string {
	raw := r.URL.EscapedPath()
	match := route.re.FindStringSubmatch(raw)
	if match == nil {
		return "mockapi: path does not match " + route.Template
	}
	values := map[string]string{}
	for i, name := range route.re.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}
		decoded, err := url.PathUnescape(match[i])
		if err != nil {
			return "mockapi: path value for " + name + " is not valid percent-encoding"
		}
		values[name] = decoded
	}
	rendered, err := connection.RenderPath(route.Template, values, route.multiSegment...)
	if err != nil || rendered != raw {
		return "mockapi: path values must be escaped as the connection runtime escapes them (reserved characters like ( ) / must be percent-encoded)"
	}
	return ""
}
