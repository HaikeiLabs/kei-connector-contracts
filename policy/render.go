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

package policy

import (
	"slices"
	"strings"

	"github.com/HaikeiLabs/kei-connector-contracts/harnessmatch"
)

// RenderEntry is a rule rendered for a harness.
type RenderEntry struct {
	RuleID string
	Target string
	Effect Effect
}

// Diagnostic is a rule that was not rendered and why.
type Diagnostic struct {
	RuleID string
	Kind   string // "withheld" or "unrenderable"
	Reason string
}

// RenderResult is the output of Render.
type RenderResult struct {
	Entries     []RenderEntry
	Diagnostics []Diagnostic
}

// Render renders the rules for one harness principal and kind. A rule whose
// principal matches the given user/email/groups is rendered as an entry,
// unless the harness kind is unknown (unrenderable) or, for a permit, a
// deny that is not rendered for this harness overlaps its target (withheld,
// fail closed). A rule whose principal does not match is withheld. harnessmatch
// supplies the known-kind check.
//
// MatchPrincipal only matches the supported principal kinds (user:, email:,
// group:, org:, and "*"). harness:/agent: sources are not matched and return
// false. The caller must pre-filter those before calling Render.
func Render(rules []Rule, user, email string, groups []string, harnessKind string) RenderResult {
	res := RenderResult{}
	kindOK := harnessmatch.ValidKind(harnessKind)

	var candidates []Rule
	var withheldDenyTargets []string

	for _, r := range rules {
		if !r.Enabled {
			if r.Effect == "deny" {
				withheldDenyTargets = append(withheldDenyTargets, r.Target)
			}
			continue
		}
		if !MatchPrincipal(r.Principal, user, email, groups) {
			res.Diagnostics = append(res.Diagnostics, Diagnostic{
				RuleID: r.ID,
				Kind:   "withheld",
				Reason: "principal does not match",
			})
			if r.Effect == "deny" {
				withheldDenyTargets = append(withheldDenyTargets, r.Target)
			}
			continue
		}
		candidates = append(candidates, r)
	}

	for _, r := range candidates {
		if !kindOK {
			res.Diagnostics = append(res.Diagnostics, Diagnostic{
				RuleID: r.ID,
				Kind:   "unrenderable",
				Reason: harnessmatch.ReasonUnsupportedKind,
			})
			continue
		}
		if r.Effect == "permit" && overlapsAny(r.Target, withheldDenyTargets) {
			res.Diagnostics = append(res.Diagnostics, Diagnostic{
				RuleID: r.ID,
				Kind:   "withheld",
				Reason: "overlapping deny is not rendered",
			})
			continue
		}
		res.Entries = append(res.Entries, RenderEntry{
			RuleID: r.ID,
			Target: r.Target,
			Effect: r.Effect,
		})
	}
	return res
}

// overlapsAny reports whether target overlaps any of the given targets.
func overlapsAny(target string, others []string) bool {
	for _, o := range others {
		if targetsOverlap(target, o) {
			return true
		}
	}
	return false
}

// targetsOverlap reports whether two target patterns can match the same
// concrete target: exact equality, either is "*", or one is a "<scheme>:*"
// wildcard covering the other's scheme, one shell: pattern's whitespace-split
// argv is a prefix of the other's (either direction), or two path: patterns'
// glob segments could match a common path.
func targetsOverlap(a, b string) bool {
	if a == b || a == "*" || b == "*" {
		return true
	}
	if s, ok := schemeWildcard(a); ok {
		if bs, _, has := strings.Cut(b, ":"); has && bs == s {
			return true
		}
	}
	if s, ok := schemeWildcard(b); ok {
		if as, _, has := strings.Cut(a, ":"); has && as == s {
			return true
		}
	}
	aScheme, aValue, aOK := strings.Cut(a, ":")
	bScheme, bValue, bOK := strings.Cut(b, ":")
	if aOK && bOK && aScheme == bScheme {
		switch aScheme {
		case "shell":
			da, pa := strings.Fields(aValue), strings.Fields(bValue)
			if len(da) == 0 || da[0] == "*" || len(pa) == 0 || pa[0] == "*" {
				return true
			}
			n := min(len(da), len(pa))
			return slices.Equal(da[:n], pa[:n])
		case "path":
			return globsOverlap(strings.Split(aValue, "/"), strings.Split(bValue, "/"))
		}
	}
	return false
}

// schemeWildcard reports whether t is a "<scheme>:*" wildcard and returns
// its scheme.
func schemeWildcard(t string) (string, bool) {
	s, ok := strings.CutSuffix(t, ":*")
	if !ok || s == "" || strings.ContainsAny(s, ":*") {
		return "", false
	}
	return s, true
}

// globsOverlap reports whether two path globs could match a common path. A
// ** or ~ segment overlaps anything that follows; a segment holding * is
// taken to overlap any segment.
func globsOverlap(a, b []string) bool {
	if len(a) == 0 || len(b) == 0 {
		return len(a) == len(b)
	}
	if a[0] == "**" || b[0] == "**" || a[0] == "*" && len(a) == 1 || b[0] == "*" && len(b) == 1 || a[0] == "~" || b[0] == "~" {
		return true
	}
	if a[0] != b[0] && !strings.Contains(a[0], "*") && !strings.Contains(b[0], "*") {
		return false
	}
	return globsOverlap(a[1:], b[1:])
}
