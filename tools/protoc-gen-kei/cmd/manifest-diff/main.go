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

// Command manifest-diff compares two connector manifests and prints one
// "breaking: <op id>: <reason>" line per breaking change (proposal 0002,
// section 6). Exit status: 0 no breaking change, 1 breaking changes, 2 usage
// or input error. Run it through scripts/manifest-diff.sh.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/HaikeiLabs/kei-connector-contracts/tools/protoc-gen-kei/internal/manifest"
	"github.com/HaikeiLabs/kei-connector-contracts/tools/protoc-gen-kei/internal/manifestdiff"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: manifest-diff <old manifest.json> <new manifest.json>")
		os.Exit(2)
	}
	old, err := load(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	cur, err := load(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	findings := manifestdiff.Diff(old, cur)
	for _, f := range findings {
		fmt.Printf("breaking: %s\n", f)
	}
	if len(findings) > 0 {
		fmt.Printf("%d breaking change(s) from %s to %s; rename with deprecated and replaced_by, or ship a v2 package\n", len(findings), old.ContractsVersion, cur.ContractsVersion)
		os.Exit(1)
	}
	fmt.Printf("OK: no breaking changes from %s to %s\n", old.ContractsVersion, cur.ContractsVersion)
}

func load(path string) (*manifest.Manifest, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("manifest-diff: %w", err)
	}
	var m manifest.Manifest
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, fmt.Errorf("manifest-diff: %s: %w", path, err)
	}
	if m.Schema != manifest.Schema {
		return nil, fmt.Errorf("manifest-diff: %s: schema %q, want %q", path, m.Schema, manifest.Schema)
	}
	return &m, nil
}
