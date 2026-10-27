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

// Package policy defines the shared policy model: principals (who a policy
// applies to), targets (what a policy acts on), and pattern matching. The
// package is stdlib-only and dependency-free.
package policy

import (
	"fmt"
	"strings"
)

// PrincipalKind is the scheme of a principal pattern:
// "any", "user", "email", "group", "autogroup", "org", "agent", "harness".
type PrincipalKind string

// Principal is a parsed principal pattern.
type Principal struct {
	Kind  PrincipalKind
	Value string
}

var principalSchemes = map[string]PrincipalKind{
	"user":      "user",
	"email":     "email",
	"group":     "group",
	"autogroup": "autogroup",
	"org":       "org",
	"agent":     "agent",
	"harness":   "harness",
}

// autogroups are the built-in groups the control plane resolves from a
// caller's role instead of storing: every workspace member, and workspace
// admins (owners count as admins for policy). "autogroup:<name>" parses only
// for these names.
var autogroups = map[string]bool{"member": true, "admin": true}

// ParsePrincipal parses a principal pattern. "*" is the any principal, a
// bare name is a group, and any other pattern is "<scheme>:<value>" with a
// known principal scheme and a non-empty value. An autogroup value must be
// one of member or admin. An unknown scheme, an unknown
// autogroup name or an empty value is an error.
func ParsePrincipal(pattern string) (Principal, error) {
	if pattern == "*" {
		return Principal{Kind: "any"}, nil
	}
	scheme, value, ok := strings.Cut(pattern, ":")
	if !ok {
		return Principal{Kind: "group", Value: pattern}, nil
	}
	kind, known := principalSchemes[scheme]
	if !known {
		return Principal{}, fmt.Errorf("unknown principal scheme %q", scheme)
	}
	if value == "" {
		return Principal{}, fmt.Errorf("principal pattern %q has an empty value", pattern)
	}
	if kind == "autogroup" && !autogroups[value] {
		return Principal{}, fmt.Errorf("unknown autogroup %q; valid names: member, admin", value)
	}
	return Principal{Kind: kind, Value: value}, nil
}

// TargetKind is the scheme of a target pattern:
// "any", "shell", "skill", "path", "capability", "connector", "tool", "mcp".
type TargetKind string

// Target is a parsed target pattern.
type Target struct {
	Kind  TargetKind
	Value string
}

var targetSchemes = map[string]TargetKind{
	"shell":      "shell",
	"skill":      "skill",
	"path":       "path",
	"capability": "capability",
	"connector":  "connector",
	"tool":       "tool",
	"mcp":        "mcp",
}

// ParseTarget parses a target pattern. "*" is the any target; any other
// pattern must be "<scheme>:<value>" with a known target scheme and a
// non-empty value. A bare value without a scheme, an unknown scheme, and an
// empty value are errors.
func ParseTarget(pattern string) (Target, error) {
	if pattern == "*" {
		return Target{Kind: "any"}, nil
	}
	scheme, value, ok := strings.Cut(pattern, ":")
	if !ok {
		return Target{}, fmt.Errorf("target pattern %q has no scheme", pattern)
	}
	kind, known := targetSchemes[scheme]
	if !known {
		return Target{}, fmt.Errorf("unknown target scheme %q", scheme)
	}
	if value == "" {
		return Target{}, fmt.Errorf("target pattern %q has an empty value", pattern)
	}
	return Target{Kind: kind, Value: value}, nil
}

// MatchPrincipal reports whether pattern matches the given principal.
// input1 is the user id, input2 the email. Semantics:
//
//   - any principal
//     group:<name>      one of the supplied groups
//     <name>            one of the supplied groups (bare)
//     user:<id>         the user id
//     email:<address>   the user id or the email
//     <value>           the user id or the email, whole
//     <scheme>:*        any non-empty value in exactly that scheme, where the
//     scheme is one segment (no ":" or "*")
//     *@<domain>        an email whose part after the last "@" equals
//     email:*@<domain>  <domain>, ignoring case
//
// A star anywhere else is a literal, and no form matches a prefix or suffix
// of a value. A pattern that matches nothing returns false; there is no
// implicit allow.
func MatchPrincipal(pattern string, user, email string, groups []string) bool {
	if pattern == "*" {
		return true
	}

	for _, g := range groups {
		if "group:"+g == pattern || g == pattern {
			return true
		}
	}

	if user != "" && "user:"+user == pattern {
		return true
	}

	for _, inp := range []string{user, email} {
		if inp == "" {
			continue
		}
		if inp == pattern || "email:"+inp == pattern {
			return true
		}
		if matchesSchemeWildcard(inp, pattern) || matchesEmailDomain(inp, pattern) {
			return true
		}
	}

	return false
}

// matchesSchemeWildcard reports whether pattern is "<scheme>:*" and input is
// a non-empty value in that scheme.
func matchesSchemeWildcard(input, pattern string) bool {
	scheme, ok := strings.CutSuffix(pattern, ":*")
	if !ok || scheme == "" || strings.ContainsAny(scheme, ":*") {
		return false
	}
	value, ok := strings.CutPrefix(input, scheme+":")
	return ok && value != ""
}

// matchesEmailDomain reports whether pattern is "*@<domain>" or
// "email:*@<domain>" and input is an email address, bare or email:-qualified,
// in exactly that domain.
func matchesEmailDomain(input, pattern string) bool {
	domain, ok := strings.CutPrefix(strings.TrimPrefix(pattern, "email:"), "*@")
	if !ok || domain == "" || strings.ContainsAny(domain, "@*") {
		return false
	}
	at := strings.LastIndex(input, "@")
	if at <= 0 {
		return false
	}
	local := strings.TrimPrefix(input[:at], "email:")
	return local != "" && !strings.Contains(local, ":") && strings.EqualFold(input[at+1:], domain)
}
