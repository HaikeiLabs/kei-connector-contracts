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

package conformance

import "testing"

// TestParseCSVSkipsBlankAndFooter feeds parseCSV a sqlcmd-shaped output that
// ends with a blank line and the "(2 rows affected)" footer. sqlcmd prints
// such junk when NOCOUNT is off; the mssql queries suppress it via
// "SET NOCOUNT ON;", but parseCSV drops both kinds of lines defensively so a
// stray footer can never break parsing.
func TestParseCSVSkipsBlankAndFooter(t *testing.T) {
	sample := "10,__NULL__\n20,__NULL__\n\n(2 rows affected)\n"
	got := parseCSV(t, sample, ',')
	want := [][]string{{"10", "__NULL__"}, {"20", "__NULL__"}}
	if len(got) != len(want) {
		t.Fatalf("parsed %d rows, want %d: %v", len(got), len(want), got)
	}
	for i, w := range want {
		if got[i][0] != w[0] || got[i][1] != w[1] {
			t.Errorf("row %d = %v, want %v", i, got[i], w)
		}
	}
}
