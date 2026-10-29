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

import "testing"

// A rule's When condition narrows it to calls made through one harness kind
// (ADR-038 D-6). It filters before Precedence orders the survivors, so
// deny-overrides is unchanged.
func TestDecideWhenHarness(t *testing.T) {
	permit := Rule{ID: "p", Principal: "group:engineering", Target: "shell:git", Effect: "permit", Enabled: true}
	denyPush := Rule{ID: "d", Principal: "group:engineering", Target: "shell:git push", Effect: "deny", Enabled: true}
	onlyClaude := func(r Rule) Rule { r.When = When{Harness: "claude_code"}; return r }
	req := func(harness, target string) Request {
		return Request{User: "u1", Groups: []string{"engineering"}, Target: target, Harness: harness}
	}

	cases := []struct {
		name        string
		rules       []Rule
		req         Request
		wantDecide  string
		wantDecider string
	}{
		{"a rule with a matching when applies", []Rule{onlyClaude(permit)}, req("claude_code", "shell:git status"), "permit", "p"},
		{"a rule with a different when does not apply", []Rule{onlyClaude(permit)}, req("codex", "shell:git status"), "no_match", ""},
		{"an unknown request harness never matches a when", []Rule{onlyClaude(permit)}, req("", "shell:git status"), "no_match", ""},
		{"a rule without when applies to every harness", []Rule{permit}, req("codex", "shell:git status"), "permit", "p"},
		{"a rule without when applies when the harness is unknown", []Rule{permit}, req("", "shell:git status"), "permit", "p"},
		{"a deny with a matching when still beats a permit", []Rule{permit, onlyClaude(denyPush)}, req("claude_code", "shell:git push"), "deny", "d"},
		{"a deny with a matching when wins whatever the rule order", []Rule{onlyClaude(denyPush), permit}, req("claude_code", "shell:git push"), "deny", "d"},
		{"a deny whose when does not match cannot deny", []Rule{permit, onlyClaude(denyPush)}, req("codex", "shell:git push"), "permit", "p"},
		{"an unconditional deny beats a permit scoped by when", []Rule{onlyClaude(permit), denyPush}, req("claude_code", "shell:git push"), "deny", "d"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Decide(tc.rules, tc.req)
			if got.Decision != tc.wantDecide || got.DecidedBy != tc.wantDecider {
				t.Fatalf("Decide = %+v, want %s by %q", got, tc.wantDecide, tc.wantDecider)
			}
		})
	}
}

func TestDecideWhenOutrankedListsOnlyMatchingRules(t *testing.T) {
	rules := []Rule{
		{ID: "d", Principal: "group:engineering", Target: "shell:git", Effect: "deny", Enabled: true, When: When{Harness: "claude_code"}},
		{ID: "p-claude", Principal: "group:engineering", Target: "shell:git", Effect: "permit", Enabled: true, When: When{Harness: "claude_code"}},
		{ID: "p-codex", Principal: "group:engineering", Target: "shell:git", Effect: "permit", Enabled: true, When: When{Harness: "codex"}},
	}
	got := Decide(rules, Request{User: "u1", Groups: []string{"engineering"}, Target: "shell:git push", Harness: "claude_code"})
	if got.Decision != "deny" || got.DecidedBy != "d" {
		t.Fatalf("Decide = %+v", got)
	}
	if len(got.Outranked) != 1 || got.Outranked[0] != "p-claude" {
		t.Fatalf("Outranked = %v, want only the permit whose when matched", got.Outranked)
	}
}
