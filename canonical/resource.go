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

// Package canonical computes canonical identifiers that more than one Kei
// service must compute identically. It is the single implementation; a
// service never keeps its own copy.
package canonical

import (
	"errors"
	"net"
	"net/url"
	"strings"
)

// ErrInvalidResource is returned for a URL that is not an absolute https URL
// without user info, query or fragment. It never includes the URL, which may
// carry credentials in its user info.
var ErrInvalidResource = errors.New("canonical: resource is not an absolute https URL without user info, query or fragment")

// Resource is the canonical resource URI of an MCP server, computed from its
// URL as Kei ADR-036 section 3.4 and MCP "Canonical Server URI" specify: the
// scheme and host lowercased, no default port (:443), no trailing slash
// unless the path is exactly /, and no fragment, user info or query. Any
// other port is kept as written (so :0443 stays), IPv6 literals are kept as
// written apart from case, and the path keeps its case and escaping.
//
// Resource checks shape only. It does not resolve the host or apply SSRF
// rules; callers that connect to the URL still must.
//
// The OAuth resource parameter, a stored connector's mcp_url and a token's
// audience are all compared through this function, so any change to its
// results makes stored resources stop matching their tokens. The golden
// cases in canonicaltest pin them.
func Resource(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil || !strings.EqualFold(u.Scheme, "https") || u.Host == "" || u.User != nil || u.Fragment != "" || u.RawQuery != "" || u.ForceQuery {
		return "", ErrInvalidResource
	}
	host := strings.ToLower(u.Hostname())
	if port := u.Port(); port != "" && port != "443" {
		host = net.JoinHostPort(host, port)
	} else if strings.Contains(host, ":") {
		host = "[" + host + "]"
	}
	path := u.EscapedPath()
	if path != "/" {
		path = strings.TrimRight(path, "/")
	}
	if path == "" {
		path = "/"
	}
	return "https://" + host + path, nil
}
