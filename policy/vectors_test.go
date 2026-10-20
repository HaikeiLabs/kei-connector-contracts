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
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type parseWant struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
	Error bool   `json:"error"`
}

type parseCase struct {
	Name  string    `json:"name"`
	Input string    `json:"input"`
	Want  parseWant `json:"want"`
}

type matchCase struct {
	Name    string   `json:"name"`
	Pattern string   `json:"pattern"`
	User    string   `json:"user"`
	Email   string   `json:"email"`
	Groups  []string `json:"groups"`
	Want    bool     `json:"want"`
}

func vectorPath(name string) string {
	return filepath.Join("testdata", "vectors", name)
}

func loadParseCases(name string) ([]parseCase, error) {
	raw, err := os.ReadFile(vectorPath(name))
	if err != nil {
		return nil, err
	}
	var suite struct {
		File    string      `json:"file"`
		Version int         `json:"version"`
		Cases   []parseCase `json:"cases"`
	}
	if err := json.Unmarshal(raw, &suite); err != nil {
		return nil, err
	}
	return suite.Cases, nil
}

func loadMatchCases(name string) ([]matchCase, error) {
	raw, err := os.ReadFile(vectorPath(name))
	if err != nil {
		return nil, err
	}
	var suite struct {
		File    string      `json:"file"`
		Version int         `json:"version"`
		Cases   []matchCase `json:"cases"`
	}
	if err := json.Unmarshal(raw, &suite); err != nil {
		return nil, err
	}
	return suite.Cases, nil
}

func TestParsePrincipal(t *testing.T) {
	cases, err := loadParseCases("parse-principal.json")
	if err != nil {
		t.Fatalf("loading fixtures: %v", err)
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			got, err := ParsePrincipal(tc.Input)
			if tc.Want.Error {
				if err == nil {
					t.Fatalf("expected error, got %+v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(got.Kind) != tc.Want.Kind {
				t.Errorf("kind: got %q want %q", got.Kind, tc.Want.Kind)
			}
			if got.Value != tc.Want.Value {
				t.Errorf("value: got %q want %q", got.Value, tc.Want.Value)
			}
		})
	}
}

func TestParseTarget(t *testing.T) {
	cases, err := loadParseCases("parse-target.json")
	if err != nil {
		t.Fatalf("loading fixtures: %v", err)
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			got, err := ParseTarget(tc.Input)
			if tc.Want.Error {
				if err == nil {
					t.Fatalf("expected error, got %+v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(got.Kind) != tc.Want.Kind {
				t.Errorf("kind: got %q want %q", got.Kind, tc.Want.Kind)
			}
			if got.Value != tc.Want.Value {
				t.Errorf("value: got %q want %q", got.Value, tc.Want.Value)
			}
		})
	}
}

type decideRule struct {
	ID        string `json:"id"`
	Principal string `json:"principal"`
	Target    string `json:"target"`
	Effect    string `json:"effect"`
	Priority  int    `json:"priority"`
	Enabled   bool   `json:"enabled"`
}

type decideGrant struct {
	Principal  string `json:"principal"`
	Target     string `json:"target"`
	ExpireTime string `json:"expire_time"` // RFC 3339
}

type decideCase struct {
	Name    string        `json:"name"`
	Rules   []decideRule  `json:"rules"`
	Request Request       `json:"request"`
	Now     string        `json:"now,omitempty"` // RFC 3339
	Grants  []decideGrant `json:"grants,omitempty"`
	Want    decideWant    `json:"want"`
}

type decideWant struct {
	Decision  string   `json:"decision"`
	Reason    string   `json:"reason,omitempty"`
	DecidedBy string   `json:"decided_by"`
	Outranked []string `json:"outranked"`
}

func loadDecideCases(name string) ([]decideCase, error) {
	raw, err := os.ReadFile(vectorPath(name))
	if err != nil {
		return nil, err
	}
	var suite struct {
		File    string       `json:"file"`
		Version int          `json:"version"`
		Cases   []decideCase `json:"cases"`
	}
	if err := json.Unmarshal(raw, &suite); err != nil {
		return nil, err
	}
	return suite.Cases, nil
}

func TestDecide(t *testing.T) {
	cases, err := loadDecideCases("decide.json")
	if err != nil {
		t.Fatalf("loading fixtures: %v", err)
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			rules := make([]Rule, len(tc.Rules))
			for i, r := range tc.Rules {
				rules[i] = Rule{
					ID:        r.ID,
					Principal: r.Principal,
					Target:    r.Target,
					Effect:    Effect(r.Effect),
					Priority:  r.Priority,
					Enabled:   r.Enabled,
				}
			}
			got := Decide(rules, tc.Request)
			if got.Decision != tc.Want.Decision {
				t.Errorf("decision: got %q want %q", got.Decision, tc.Want.Decision)
			}
			if got.DecidedBy != tc.Want.DecidedBy {
				t.Errorf("decided_by: got %q want %q", got.DecidedBy, tc.Want.DecidedBy)
			}
			if len(got.Outranked) != len(tc.Want.Outranked) {
				t.Errorf("outranked: got %v want %v", got.Outranked, tc.Want.Outranked)
				return
			}
			for i := range got.Outranked {
				if got.Outranked[i] != tc.Want.Outranked[i] {
					t.Errorf("outranked[%d]: got %q want %q", i, got.Outranked[i], tc.Want.Outranked[i])
				}
			}
			if tc.Want.Reason != "" && got.Reason != tc.Want.Reason {
				t.Errorf("reason: got %q want %q", got.Reason, tc.Want.Reason)
			}
		})
	}
}

func loadDecideAtCases(name string) ([]decideCase, error) {
	raw, err := os.ReadFile(vectorPath(name))
	if err != nil {
		return nil, err
	}
	var suite struct {
		File    string       `json:"file"`
		Version int          `json:"version"`
		Cases   []decideCase `json:"cases"`
	}
	if err := json.Unmarshal(raw, &suite); err != nil {
		return nil, err
	}
	return suite.Cases, nil
}

func TestDecideAt(t *testing.T) {
	cases, err := loadDecideAtCases("decide-at.json")
	if err != nil {
		t.Fatalf("loading fixtures: %v", err)
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			rules := make([]Rule, len(tc.Rules))
			for i, r := range tc.Rules {
				rules[i] = Rule{
					ID:        r.ID,
					Principal: r.Principal,
					Target:    r.Target,
					Effect:    Effect(r.Effect),
					Priority:  r.Priority,
					Enabled:   r.Enabled,
				}
			}
			grants := make([]AccessGrant, len(tc.Grants))
			for i, g := range tc.Grants {
				exp, err := time.Parse(time.RFC3339, g.ExpireTime)
				if err != nil {
					t.Fatalf("bad expire_time %q: %v", g.ExpireTime, err)
				}
				grants[i] = AccessGrant{
					Principal:  g.Principal,
					Target:     g.Target,
					ExpireTime: exp,
				}
			}
			var now time.Time
			if tc.Now != "" {
				var err error
				now, err = time.Parse(time.RFC3339, tc.Now)
				if err != nil {
					t.Fatalf("bad now %q: %v", tc.Now, err)
				}
			}
			got := DecideAt(rules, tc.Request, grants, now)
			if got.Decision != tc.Want.Decision {
				t.Errorf("decision: got %q want %q", got.Decision, tc.Want.Decision)
			}
			if got.Reason != tc.Want.Reason {
				t.Errorf("reason: got %q want %q", got.Reason, tc.Want.Reason)
			}
			if got.DecidedBy != tc.Want.DecidedBy {
				t.Errorf("decided_by: got %q want %q", got.DecidedBy, tc.Want.DecidedBy)
			}
			if len(got.Outranked) != len(tc.Want.Outranked) {
				t.Errorf("outranked: got %v want %v", got.Outranked, tc.Want.Outranked)
				return
			}
			for i := range got.Outranked {
				if got.Outranked[i] != tc.Want.Outranked[i] {
					t.Errorf("outranked[%d]: got %q want %q", i, got.Outranked[i], tc.Want.Outranked[i])
				}
			}
		})
	}
}

type renderEntry struct {
	RuleID string `json:"rule_id"`
	Target string `json:"target"`
	Effect string `json:"effect"`
}

type renderDiagnostic struct {
	RuleID string `json:"rule_id"`
	Kind   string `json:"kind"`
	Reason string `json:"reason"`
}

type renderCase struct {
	Name        string       `json:"name"`
	Rules       []decideRule `json:"rules"`
	User        string       `json:"user"`
	Email       string       `json:"email"`
	Groups      []string     `json:"groups"`
	HarnessKind string       `json:"harness_kind"`
	Want        renderWant   `json:"want"`
}

type renderWant struct {
	Entries     []renderEntry      `json:"entries"`
	Diagnostics []renderDiagnostic `json:"diagnostics"`
}

func loadRenderCases(name string) ([]renderCase, error) {
	raw, err := os.ReadFile(vectorPath(name))
	if err != nil {
		return nil, err
	}
	var suite struct {
		File    string       `json:"file"`
		Version int          `json:"version"`
		Cases   []renderCase `json:"cases"`
	}
	if err := json.Unmarshal(raw, &suite); err != nil {
		return nil, err
	}
	return suite.Cases, nil
}

func TestRender(t *testing.T) {
	cases, err := loadRenderCases("render.json")
	if err != nil {
		t.Fatalf("loading fixtures: %v", err)
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			rules := make([]Rule, len(tc.Rules))
			for i, r := range tc.Rules {
				rules[i] = Rule{
					ID:        r.ID,
					Principal: r.Principal,
					Target:    r.Target,
					Effect:    Effect(r.Effect),
					Priority:  r.Priority,
					Enabled:   r.Enabled,
				}
			}
			got := Render(rules, tc.User, tc.Email, tc.Groups, tc.HarnessKind)
			if len(got.Entries) != len(tc.Want.Entries) {
				t.Errorf("entries: got %v want %v", got.Entries, tc.Want.Entries)
			} else {
				for i := range got.Entries {
					w := tc.Want.Entries[i]
					if got.Entries[i].RuleID != w.RuleID || got.Entries[i].Target != w.Target || string(got.Entries[i].Effect) != w.Effect {
						t.Errorf("entries[%d]: got %+v want %+v", i, got.Entries[i], w)
					}
				}
			}
			if len(got.Diagnostics) != len(tc.Want.Diagnostics) {
				t.Errorf("diagnostics: got %v want %v", got.Diagnostics, tc.Want.Diagnostics)
			} else {
				for i := range got.Diagnostics {
					w := tc.Want.Diagnostics[i]
					if got.Diagnostics[i].RuleID != w.RuleID || got.Diagnostics[i].Kind != w.Kind || got.Diagnostics[i].Reason != w.Reason {
						t.Errorf("diagnostics[%d]: got %+v want %+v", i, got.Diagnostics[i], w)
					}
				}
			}
		})
	}
}

func TestMatchPrincipal(t *testing.T) {
	cases, err := loadMatchCases("match-principal.json")
	if err != nil {
		t.Fatalf("loading fixtures: %v", err)
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			got := MatchPrincipal(tc.Pattern, tc.User, tc.Email, tc.Groups)
			if got != tc.Want {
				t.Errorf("got %v want %v", got, tc.Want)
			}
		})
	}
}
