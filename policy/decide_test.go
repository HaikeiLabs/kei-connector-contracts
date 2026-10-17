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

// Copyright 2026 Haikei Labs

package policy

import "testing"

func TestMatchTarget(t *testing.T) {
	cases := []struct {
		rule, req string
		want      bool
	}{
		{"connector:jira", "connector:jira", true},
		{"connector:jira", "connector:myjira", false},
		{"*", "capability:read", true},
		{"connector:*", "connector:jira", true},
		{"connector:*", "capability:read", false},
		{"tool:*", "tool:", false},
		{"tool:*", "tool:x", true},
	}
	for _, tc := range cases {
		if got := MatchTarget(tc.rule, tc.req); got != tc.want {
			t.Errorf("MatchTarget(%q,%q)=%v want %v", tc.rule, tc.req, got, tc.want)
		}
	}
}
