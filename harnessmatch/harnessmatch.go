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

// Package harnessmatch implements the kei.harness-match/v1 dialect (ADR-029).
// It matches a harness tool call against a set of policies and returns an
// outcome: permit, deny, unmatched, or not_renderable.  The package is
// stdlib-only and dependency-free.
package harnessmatch

import (
	"path"
	"strings"
)

// Kinds is the set of known harness kinds.
var Kinds = map[string]bool{"claude_code": true, "codex": true, "opencode": true, "custom": true}

// Outcome is the result of matching a harness call against policies.
type Outcome string

const (
	OutcomePermit        Outcome = "permit"
	OutcomeDeny          Outcome = "deny"
	OutcomeUnmatched     Outcome = "unmatched"
	OutcomeNotRenderable Outcome = "not_renderable"
)

// Reason codes for non-permit, non-deny outcomes (ADR-029 §1.5 extension).
const (
	ReasonUnsupportedKind  = "unsupported_kind"
	ReasonPathNotAbsolute  = "path_not_absolute"
	ReasonBundleInvalid    = "bundle_invalid"
	ReasonHumanSource      = "human_source"
	ReasonDstNotRenderable = "dst_not_renderable"
	ReasonEvaluateViaPDP   = "evaluate_via_pdp"
)

// Call is one tool call of a pre-built harness, described by the inputs
// kei.harness-match/v1 matches on (ADR-029 §1.2, §1.3).
type Call struct {
	HarnessID string   // registered harness id (UUID), if any
	Kind      string   // claude_code | codex | opencode | custom
	AgentID   string   // agent the harness runs as, if any
	Tool      string   // registry tool name, e.g. claude_code.edit
	Argv      []string // shell commands only
	Skill     string   // skill invocations only
	Path      string   // absolute path outside the working directory
	Home      string   // home directory that a leading ~/ in path: globs expands to
	MCPServer string
	MCPTool   string
}

// Policy is a single kei.harness-match/v1 policy entry.
type Policy struct {
	ID      string  // policy identifier (for audit and render blocking)
	Src     string  // source pattern: *, harness:<kind|id>, agent:<id>
	Dst     string  // destination pattern: *, shell:<prefix>, skill:<name>, path:<glob>, mcp:<server>/<tool>, tool:<name>
	Action  string  // "permit" or "deny"
	Enabled bool    // false skips the policy
	Scope   *string // non-nil means the policy is scoped to this agent ID only
}

// Result is the outcome of Evaluate.
type Result struct {
	Outcome  Outcome // permit, deny, unmatched, or not_renderable
	PolicyID string  // the matched policy's ID (empty for unmatched)
	Reason   string  // reason code for not_renderable outcomes
}

// Evaluate matches a harness call against a set of policies and returns the
// first matching outcome.  Policies are evaluated in order; the first policy
// whose source and destination both match wins.  A human-only source pattern
// gives not_renderable.  For a shell command, a permit counts only when its
// dst is a shell: argv prefix (ADR-029 §9.4).
func Evaluate(call Call, policies []Policy) Result {
	if !Kinds[call.Kind] {
		return Result{Outcome: OutcomeNotRenderable, Reason: ReasonUnsupportedKind}
	}
	if call.Kind == "custom" {
		return Result{Outcome: OutcomeNotRenderable, Reason: ReasonEvaluateViaPDP}
	}
	if call.Path != "" && (!path.IsAbs(call.Path) || path.Clean(call.Path) != call.Path) {
		return Result{Outcome: OutcomeNotRenderable, Reason: ReasonPathNotAbsolute}
	}

	shell := len(call.Argv) > 0
	for _, p := range policies {
		if !p.Enabled {
			continue
		}
		if p.Scope != nil && *p.Scope != call.AgentID {
			continue
		}
		if p.Action != "permit" && p.Action != "deny" {
			continue
		}

		dstMatch, dstRenderable := MatchDst(p.Dst, call)
		if !dstMatch {
			continue
		}

		// Shell-permit rule: a permit on a shell call counts only with
		// a shell: argv-prefix dst.  shell:* and non-shell permits on
		// shell calls are rejected (ADR-029 §9.4).
		if shell && p.Action == "permit" && (!strings.HasPrefix(p.Dst, "shell:") || p.Dst == "shell:*") {
			continue
		}

		srcMatch, srcHuman := MatchSrc(p.Src, call)
		if !srcMatch && !srcHuman {
			continue
		}
		if srcHuman {
			return Result{Outcome: OutcomeNotRenderable, PolicyID: p.ID, Reason: ReasonHumanSource}
		}
		if !dstRenderable {
			return Result{Outcome: OutcomeNotRenderable, PolicyID: p.ID, Reason: ReasonDstNotRenderable}
		}
		if p.Action == "deny" {
			return Result{Outcome: OutcomeDeny, PolicyID: p.ID}
		}
		return Result{Outcome: OutcomePermit, PolicyID: p.ID}
	}
	return Result{Outcome: OutcomeUnmatched}
}

// MatchSrc reports whether a source pattern applies to a harness call.
// Returns (applies, human) where human indicates the pattern is a
// human-only source pattern that cannot be resolved by a harness renderer.
func MatchSrc(src string, call Call) (applies, human bool) {
	switch {
	case src == "*", src == "harness:*":
		return true, false
	case strings.HasPrefix(src, "harness:"):
		v := strings.TrimPrefix(src, "harness:")
		return v == call.Kind || (call.HarnessID != "" && v == call.HarnessID), false
	case strings.HasPrefix(src, "agent:"):
		v := strings.TrimPrefix(src, "agent:")
		return call.AgentID != "" && (v == "*" || v == call.AgentID), false
	}
	// Any other pattern (user:, email:, group:, org:, or unknown) is a
	// human-only source that cannot be resolved by a harness renderer.
	return false, true
}

// MatchDst reports whether a destination pattern matches a harness call.
// Returns (matched, renderable).  renderable is false when the pattern
// cannot be evaluated for this call (e.g. a ~/ glob with no home set).
func MatchDst(dst string, call Call) (matched, renderable bool) {
	if dst == "*" {
		return true, true
	}
	scheme, value, ok := strings.Cut(dst, ":")
	if !ok {
		return false, true
	}
	switch scheme {
	case "shell":
		if len(call.Argv) == 0 {
			return false, true
		}
		if value == "*" {
			return true, true
		}
		tokens := strings.Split(value, " ")
		if len(tokens) > len(call.Argv) {
			return false, true
		}
		for i, t := range tokens {
			if t == "" || t != call.Argv[i] {
				return false, true
			}
		}
		return true, true
	case "skill":
		return call.Skill != "" && (value == "*" || value == call.Skill), true
	case "path":
		if call.Path == "" {
			return false, true
		}
		if value == "*" {
			return true, true
		}
		if strings.HasPrefix(value, "~/") {
			if call.Home == "" {
				return true, false
			}
			value = strings.TrimSuffix(call.Home, "/") + value[1:]
		}
		return globMatch(strings.Split(value, "/"), strings.Split(call.Path, "/")), true
	case "mcp":
		if call.MCPServer == "" {
			return false, true
		}
		return value == "*" || value == call.MCPServer || (call.MCPTool != "" && value == call.MCPServer+"."+call.MCPTool), true
	case "tool":
		return call.Tool != "" && (value == "*" || value == call.Tool), true
	}
	return false, true
}

// globMatch matches path segments: ** as a whole segment matches zero or
// more segments, * matches any characters within one segment, nothing else
// is special.  This mirrors the runtime's internal/policybundle
// implementation.
func globMatch(pattern, segs []string) bool {
	if len(pattern) == 0 {
		return len(segs) == 0
	}
	if pattern[0] == "**" {
		for i := 0; i <= len(segs); i++ {
			if globMatch(pattern[1:], segs[i:]) {
				return true
			}
		}
		return false
	}
	if len(segs) == 0 || !segmentMatch(pattern[0], segs[0]) {
		return false
	}
	return globMatch(pattern[1:], segs[1:])
}

// segmentMatch matches a single path segment against a pattern segment.
// The pattern may contain * which matches any sequence of characters.
func segmentMatch(pattern, seg string) bool {
	parts := strings.Split(pattern, "*")
	if len(parts) == 1 {
		return pattern == seg
	}
	if !strings.HasPrefix(seg, parts[0]) {
		return false
	}
	seg = seg[len(parts[0]):]
	for _, part := range parts[1 : len(parts)-1] {
		i := strings.Index(seg, part)
		if i < 0 {
			return false
		}
		seg = seg[i+len(part):]
	}
	return strings.HasSuffix(seg, parts[len(parts)-1])
}

// ValidKind reports whether kind is a known harness kind.
func ValidKind(kind string) bool { return Kinds[kind] }
