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

package skill

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateLinearGolden(t *testing.T) {
	// Open the linear tool-definition file.
	defPath := filepath.Join("..", "..", "schemas", "linear.tool-definition.v1.json")
	f, err := os.Open(defPath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	got, err := Generate(f, Options{ConnectorID: "linear"})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if !strings.Contains(string(got), `"url_template": "https://api.linear.app/graphql"`) {
		t.Fatal("generated skill must preserve binding.transport for harness calls")
	}

	// Write actual output for inspection.
	actualPath := filepath.Join("testdata", "linear.SKILL.md.actual")
	if err := os.WriteFile(actualPath, got, 0644); err != nil {
		t.Fatal(err)
	}

	// Compare with golden file.
	goldenPath := filepath.Join("testdata", "linear.SKILL.md.golden")
	golden, err := os.ReadFile(goldenPath)
	if err != nil {
		// Golden doesn't exist yet: first run. Copy actual to golden so the
		// test passes after review.
		t.Logf("golden file %s not found; writing actual as golden", goldenPath)
		if err := os.WriteFile(goldenPath, got, 0644); err != nil {
			t.Fatal(err)
		}
		// Fail so the developer reviews the output.
		t.Fatal("golden file did not exist; created from actual. Review and re-run.")
	}

	if string(got) != string(golden) {
		t.Errorf("generated SKILL.md differs from golden.\n--- want (golden) ---\n%s\n--- got ---\n%s",
			string(golden), string(got))
	}
}

func TestGenerateAllDefinitions(t *testing.T) {
	// Verify every .tool-definition.v1.json in schemas/ can be processed.
	pattern := filepath.Join("..", "..", "schemas", "*.tool-definition.v1.json")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) == 0 {
		t.Fatal("no .tool-definition.v1.json files found")
	}

	for _, path := range matches {
		name := strings.TrimSuffix(filepath.Base(path), ".tool-definition.v1.json")
		t.Run(name, func(t *testing.T) {
			f, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			_, err = Generate(f, Options{ConnectorID: name})
			if err != nil {
				t.Fatalf("Generate(%s): %v", name, err)
			}
		})
	}
}
