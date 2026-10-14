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

package canonicaltest

import "testing"

// TestResourceCasesWellFormed: every case is named once and either expects a
// result or an error, never both, and the table covers the rules ADR-036
// section 3.4 names.
func TestResourceCasesWellFormed(t *testing.T) {
	cases := ResourceCases()
	if len(cases) == 0 {
		t.Fatal("no resource cases")
	}
	names := map[string]bool{}
	for _, c := range cases {
		if c.Name == "" || names[c.Name] {
			t.Errorf("case name %q is empty or repeated", c.Name)
		}
		names[c.Name] = true
		if c.Error == (c.Want != "") {
			t.Errorf("%s: want %q with error %v; a case expects exactly one", c.Name, c.Want, c.Error)
		}
	}
	for _, name := range []string{"host is lowercased", "every trailing slash is removed", "root path is kept", "default port 443 is removed", "port 0443 is kept as written", "IPv6 literal is kept as written", "user info is rejected", "fragment is rejected", "query is rejected", "http is rejected"} {
		if !names[name] {
			t.Errorf("golden table lacks %q", name)
		}
	}
	// ResourceCases returns a copy, so a caller cannot change the table.
	cases[0].Want = "changed"
	if ResourceCases()[0].Want == "changed" {
		t.Error("ResourceCases shares its backing array")
	}
}
