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

package bundle

import (
	"testing"

	"github.com/HaikeiLabs/kei-connector-contracts/harnessmatch"
)

func TestSourceMatches(t *testing.T) {
	cases := []struct {
		name    string
		pattern string
		kind    string // harness kind
		userID  string
		email   string
		groups  []string
		roles   []string
		want    bool
	}{
		// harness:* matches any harness kind
		{"harness wildcard matches pi", "harness:*", "pi", "", "", nil, nil, true},
		{"harness wildcard matches codex", "harness:*", "codex", "", "", nil, nil, true},
		{"harness wildcard matches unknown", "harness:*", "unknown", "", "", nil, nil, true},

		// harness:<kind> matches when kind matches
		{"harness:pi matches pi kind", "harness:pi", "pi", "", "", nil, nil, true},
		{"harness:codex matches codex kind", "harness:codex", "codex", "", "", nil, nil, true},

		// harness:<kind> does NOT match when kind differs
		{"harness:pi does not match codex", "harness:pi", "codex", "", "", nil, nil, false},
		{"harness:codex does not match pi", "harness:codex", "pi", "", "", nil, nil, false},

		// harness:<unknown> does not match anything
		{"harness:unknown does not match pi", "harness:unknown", "pi", "", "", nil, nil, false},
		{"harness:unknown does not match codex", "harness:unknown", "codex", "", "", nil, nil, false},
		{"harness:unknown does not match empty kind", "harness:unknown", "", "", "", nil, nil, false},

		// Non-harness patterns delegate to MatchPrincipal
		{"star matches any", "*", "", "user-a", "", nil, nil, true},
		{"star matches empty user", "*", "", "", "", nil, nil, true},
		{"user matches", "user:user-a", "", "user-a", "", nil, nil, true},
		{"user mismatch", "user:user-b", "", "user-a", "", nil, nil, false},
		{"email matches", "email:a@b.com", "", "", "a@b.com", nil, nil, true},
		{"group matches", "group:admins", "", "", "", []string{"admins"}, nil, true},
		{"group mismatch", "group:admins", "", "", "", []string{"users"}, nil, false},

		// Non-harness with non-empty kind (kind is irrelevant for non-harness patterns)
		{"star with non-empty kind", "*", "pi", "user-a", "", nil, nil, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			call := harnessmatch.Call{Kind: tc.kind}
			subject := &SessionSubject{UserID: tc.userID, Email: tc.email}
			got := sourceMatches(tc.pattern, call, subject, tc.groups, tc.roles)
			if got != tc.want {
				t.Errorf("sourceMatches(%q, call{Kind:%q}, subject{%q,%q}, %v, %v) = %v, want %v",
					tc.pattern, tc.kind, tc.userID, tc.email, tc.groups, tc.roles, got, tc.want)
			}
		})
	}
}
