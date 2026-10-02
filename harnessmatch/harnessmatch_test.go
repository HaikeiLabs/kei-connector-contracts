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

package harnessmatch

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

type testCase struct {
	ID          string       `json:"id"`
	Desc        string       `json:"desc"`
	Call        testCall     `json:"call"`
	Policies    []testPolicy `json:"policies"`
	WantOutcome string       `json:"want_outcome"`
	WantPolicy  string       `json:"want_policy,omitempty"`
	WantReason  string       `json:"want_reason,omitempty"`
}

type testCall struct {
	Kind      string   `json:"kind,omitempty"`
	HarnessID string   `json:"harness_id,omitempty"`
	AgentID   string   `json:"agent_id,omitempty"`
	Tool      string   `json:"tool,omitempty"`
	Argv      []string `json:"argv,omitempty"`
	Skill     string   `json:"skill,omitempty"`
	Path      string   `json:"path,omitempty"`
	Home      string   `json:"home,omitempty"`
	MCPServer string   `json:"mcp_server,omitempty"`
	MCPTool   string   `json:"mcp_tool,omitempty"`
}

type testPolicy struct {
	ID      string `json:"id"`
	Src     string `json:"src"`
	Dst     string `json:"dst"`
	Action  string `json:"action"`
	Enabled bool   `json:"enabled"`
	Scope   string `json:"scope_agent_id,omitempty"`
}

type testSuite struct {
	Description string     `json:"description"`
	Semantics   string     `json:"semantics"`
	Cases       []testCase `json:"cases"`
}

// FixturePath returns the path to the shared test data file.
func FixturePath() string {
	return filepath.Join("testdata", "cases.v1.json")
}

// LoadFixtures loads all test cases from the shared fixture file.
func LoadFixtures() ([]testCase, error) {
	raw, err := os.ReadFile(FixturePath())
	if err != nil {
		return nil, err
	}
	var suite testSuite
	if err := json.Unmarshal(raw, &suite); err != nil {
		return nil, err
	}
	return suite.Cases, nil
}

func toPolicy(tp testPolicy) Policy {
	p := Policy{
		ID:      tp.ID,
		Src:     tp.Src,
		Dst:     tp.Dst,
		Action:  tp.Action,
		Enabled: tp.Enabled,
	}
	if tp.Scope != "" {
		s := tp.Scope
		p.Scope = &s
	}
	return p
}

func toCall(tc testCall) Call {
	return Call{
		Kind:      tc.Kind,
		HarnessID: tc.HarnessID,
		AgentID:   tc.AgentID,
		Tool:      tc.Tool,
		Argv:      tc.Argv,
		Skill:     tc.Skill,
		Path:      tc.Path,
		Home:      tc.Home,
		MCPServer: tc.MCPServer,
		MCPTool:   tc.MCPTool,
	}
}

func TestEvaluate(t *testing.T) {
	cases, err := LoadFixtures()
	if err != nil {
		t.Fatalf("loading fixtures: %v", err)
	}
	for _, tc := range cases {
		t.Run(tc.ID, func(t *testing.T) {
			policies := make([]Policy, len(tc.Policies))
			for i, tp := range tc.Policies {
				policies[i] = toPolicy(tp)
			}
			got := Evaluate(toCall(tc.Call), policies)
			if string(got.Outcome) != tc.WantOutcome {
				t.Errorf("outcome: got %q, want %q", got.Outcome, tc.WantOutcome)
			}
			if tc.WantPolicy != "" && got.PolicyID != tc.WantPolicy {
				t.Errorf("policy: got %q, want %q", got.PolicyID, tc.WantPolicy)
			}
			if tc.WantReason != "" && got.Reason != tc.WantReason {
				t.Errorf("reason: got %q, want %q", got.Reason, tc.WantReason)
			}
		})
	}
}

func TestMatchSrc(t *testing.T) {
	tests := []struct {
		id          string
		src         string
		call        Call
		wantApplies bool
		wantHuman   bool
	}{
		{id: "star", src: "*", call: Call{Kind: "claude_code"}, wantApplies: true, wantHuman: false},
		{id: "harness-star", src: "harness:*", call: Call{Kind: "claude_code"}, wantApplies: true, wantHuman: false},
		{id: "harness-kind-match", src: "harness:claude_code", call: Call{Kind: "claude_code"}, wantApplies: true, wantHuman: false},
		{id: "harness-kind-mismatch", src: "harness:codex", call: Call{Kind: "claude_code"}, wantApplies: false, wantHuman: false},
		{id: "harness-id-match", src: "harness:uuid-123", call: Call{Kind: "claude_code", HarnessID: "uuid-123"}, wantApplies: true, wantHuman: false},
		{id: "agent-star", src: "agent:*", call: Call{AgentID: "agent-1"}, wantApplies: true, wantHuman: false},
		{id: "agent-match", src: "agent:agent-1", call: Call{AgentID: "agent-1"}, wantApplies: true, wantHuman: false},
		{id: "agent-no-agent-id", src: "agent:agent-1", call: Call{Kind: "claude_code"}, wantApplies: false, wantHuman: false},
		{id: "user-human", src: "user:alice", call: Call{Kind: "claude_code"}, wantApplies: false, wantHuman: true},
		{id: "empty-src-not-human", src: "", call: Call{Kind: "claude_code"}, wantApplies: false, wantHuman: true},
	}
	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			applies, human := MatchSrc(tt.src, tt.call)
			if applies != tt.wantApplies || human != tt.wantHuman {
				t.Errorf("MatchSrc(%q) = (%v, %v), want (%v, %v)", tt.src, applies, human, tt.wantApplies, tt.wantHuman)
			}
		})
	}
}

func TestMatchDst(t *testing.T) {
	tests := []struct {
		id             string
		dst            string
		call           Call
		wantMatch      bool
		wantRenderable bool
	}{
		{id: "star", dst: "*", call: Call{}, wantMatch: true, wantRenderable: true},
		{id: "shell-argv-prefix", dst: "shell:git commit", call: Call{Argv: []string{"git", "commit", "-m", "msg"}}, wantMatch: true, wantRenderable: true},
		{id: "shell-argv-too-short", dst: "shell:git commit -m", call: Call{Argv: []string{"git"}}, wantMatch: false, wantRenderable: true},
		{id: "shell-star", dst: "shell:*", call: Call{Argv: []string{"ls"}}, wantMatch: true, wantRenderable: true},
		{id: "shell-no-argv", dst: "shell:*", call: Call{}, wantMatch: false, wantRenderable: true},
		{id: "skill-match", dst: "skill:foo", call: Call{Skill: "foo"}, wantMatch: true, wantRenderable: true},
		{id: "skill-no-match", dst: "skill:foo", call: Call{Skill: "bar"}, wantMatch: false, wantRenderable: true},
		{id: "path-exact", dst: "path:/home/*", call: Call{Path: "/home/user"}, wantMatch: true, wantRenderable: true},
		{id: "path-double-star", dst: "path:/home/**/*.go", call: Call{Path: "/home/user/proj/main.go"}, wantMatch: true, wantRenderable: true},
		{id: "path-tilde-no-home", dst: "path:~/proj", call: Call{Path: "/home/u/proj"}, wantMatch: true, wantRenderable: false},
		{id: "path-tilde-with-home", dst: "path:~/proj/*", call: Call{Path: "/home/u/proj/app.go", Home: "/home/u"}, wantMatch: true, wantRenderable: true},
		{id: "mcp-server", dst: "mcp:my-server", call: Call{MCPServer: "my-server"}, wantMatch: true, wantRenderable: true},
		{id: "mcp-server-tool", dst: "mcp:my-server.read", call: Call{MCPServer: "my-server", MCPTool: "read"}, wantMatch: true, wantRenderable: true},
		{id: "mcp-no-match", dst: "mcp:other", call: Call{MCPServer: "my-server"}, wantMatch: false, wantRenderable: true},
		{id: "tool-match", dst: "tool:my-tool", call: Call{Tool: "my-tool"}, wantMatch: true, wantRenderable: true},
		{id: "tool-no-match", dst: "tool:my-tool", call: Call{Tool: "other"}, wantMatch: false, wantRenderable: true},
		{id: "unknown-scheme", dst: "unknown:val", call: Call{}, wantMatch: false, wantRenderable: true},
	}
	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			gotMatch, gotRenderable := MatchDst(tt.dst, tt.call)
			if gotMatch != tt.wantMatch || gotRenderable != tt.wantRenderable {
				t.Errorf("MatchDst(%q) = (%v, %v), want (%v, %v)", tt.dst, gotMatch, gotRenderable, tt.wantMatch, tt.wantRenderable)
			}
		})
	}
}

func TestFixturesLoad(t *testing.T) {
	cases, err := LoadFixtures()
	if err != nil {
		t.Fatalf("LoadFixtures: %v", err)
	}
	if len(cases) < 40 {
		t.Fatalf("expected at least 40 cases, got %d", len(cases))
	}
}

func TestValidKind(t *testing.T) {
	if !ValidKind("claude_code") {
		t.Error("expected claude_code to be valid")
	}
	if !ValidKind("codex") {
		t.Error("expected codex to be valid")
	}
	if !ValidKind("opencode") {
		t.Error("expected opencode to be valid")
	}
	if !ValidKind("custom") {
		t.Error("expected custom to be valid")
	}
	if ValidKind("vim") {
		t.Error("expected vim to be invalid")
	}
}
