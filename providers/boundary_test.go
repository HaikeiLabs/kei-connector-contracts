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

package providers

import (
	"os/exec"
	"strings"
	"testing"
)

// TestProvidersImportBoundary asserts that the providers packages import
// nothing from kei-connector-runtime and no HaikeiLabs internal package, so
// the package stays a harness-consumable library (ADR-036 section 9.2). It
// inspects the full transitive import set of every providers package. Go
// standard-library internal packages are not a boundary violation.
func TestProvidersImportBoundary(t *testing.T) {
	out, err := exec.Command("go", "list", "-deps", "-f", "{{range .Imports}}{{.}}\n{{end}}", "./...").Output()
	if err != nil {
		t.Fatalf("go list: %v", err)
	}
	var forbidden []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		importPath := strings.TrimSpace(line)
		if importPath == "" {
			continue
		}
		if strings.HasPrefix(importPath, "github.com/HaikeiLabs/kei-connector-runtime") {
			forbidden = append(forbidden, importPath)
			continue
		}
		if strings.HasPrefix(importPath, "github.com/HaikeiLabs/") && strings.Contains(importPath, "/internal/") {
			forbidden = append(forbidden, importPath)
		}
	}
	if len(forbidden) > 0 {
		t.Fatalf("providers import forbidden packages: %v", forbidden)
	}
}
