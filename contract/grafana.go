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
	"errors"
	"net"
	"net/url"
	"strings"
)

// GrafanaConfig is the grafana provider's instance: the origin of a Grafana
// Cloud stack or a self-hosted server. It is non-secret; the service-account
// token is behind the connector's CredentialRef.
type GrafanaConfig struct {
	// BaseURL is an https origin: scheme and host, an optional port, and no
	// path, query, fragment, or userinfo. The runtime appends the API paths.
	BaseURL string `json:"base_url"`
}

// Validate checks the Grafana origin. Unlike http_api, private and internal
// hosts are allowed, because Grafana is commonly self-hosted inside the
// customer's network next to the tenant runtime. Link-local addresses
// (169.254.0.0/16, fe80::/10), the AWS IPv6 metadata address fd00:ec2::254,
// unspecified addresses, and cloud-metadata hostnames are always refused: they
// are credential endpoints, never a Grafana.
func (g GrafanaConfig) Validate() error {
	invalid := errors.New("grafana base_url must be an https origin")
	if g.BaseURL == "" || len(g.BaseURL) > 2048 {
		return invalid
	}
	u, err := url.Parse(g.BaseURL)
	if err != nil || u.Scheme != "https" || u.Host == "" || u.Hostname() == "" || u.User != nil || u.Opaque != "" ||
		u.RawQuery != "" || u.ForceQuery || u.Fragment != "" || strings.Contains(g.BaseURL, "#") || (u.Path != "" && u.Path != "/") {
		return invalid
	}
	if metadataDestination(u.Hostname()) {
		return errors.New("grafana base_url must not be a link-local or cloud metadata address")
	}
	return nil
}

var awsIPv6Metadata = net.ParseIP("fd00:ec2::254")

// metadataDestination reports whether host is a link-local address, an
// unspecified address, or a cloud-metadata hostname. A host whose last label
// is numeric but which is not a dotted-quad IP (2852039166, 0xa9fea9fe,
// 0251.0376.0251.0376) is refused too: resolvers read those as IPv4 literals,
// so they would bypass the address check.
func metadataDestination(host string) bool {
	host = strings.TrimSuffix(strings.ToLower(host), ".")
	if ip := net.ParseIP(host); ip != nil {
		return ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified() || ip.Equal(awsIPv6Metadata)
	}
	if numericLabel(host[strings.LastIndex(host, ".")+1:]) {
		return true
	}
	for _, name := range []string{"metadata", "instance-data"} {
		if host == name || strings.HasPrefix(host, name+".") {
			return true
		}
	}
	return false
}

func numericLabel(label string) bool {
	digits := "0123456789"
	if strings.HasPrefix(label, "0x") {
		label, digits = label[2:], "0123456789abcdef"
	}
	if label == "" {
		return false
	}
	for _, r := range label {
		if !strings.ContainsRune(digits, r) {
			return false
		}
	}
	return true
}
