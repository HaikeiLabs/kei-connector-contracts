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

import (
	"encoding/csv"
	"regexp"
	"strings"
	"testing"
)

// affectedFooterRe matches sqlcmd's "(N rows affected)" footer. The mssql
// queries suppress it via SET NOCOUNT ON, but parseCSV drops a stray footer
// defensively so it can never break parsing.
var affectedFooterRe = regexp.MustCompile(`^\(\d+ rows? affected\)$`)

// parseCSV splits client output into rows, header first (sqlcmd -h -1 prints
// no header). Blank lines — whitespace-only ones included, as sqlcmd emits
// around batches — and the "(N rows affected)" footer are dropped before
// parsing.
func parseCSV(t *testing.T, out string, comma rune) [][]string {
	t.Helper()
	var kept []string
	for _, line := range strings.Split(out, "\n") {
		if trimmed := strings.TrimSpace(line); trimmed == "" || affectedFooterRe.MatchString(trimmed) {
			continue
		}
		kept = append(kept, line)
	}
	if len(kept) == 0 {
		t.Fatal("no rows in client output")
	}
	reader := csv.NewReader(strings.NewReader(strings.Join(kept, "\n")))
	reader.Comma = comma
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("parsing client output: %v", err)
	}
	return records
}
