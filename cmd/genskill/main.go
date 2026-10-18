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

// Command genskill generates a per-connector SKILL.md from a v1 tool-definition
// JSON file. It reads the definition from a file argument and writes the
// generated SKILL.md to stdout.
//
// Usage:
//
//	genskill -connector linear schemas/linear.tool-definition.v1.json > linear/SKILL.md
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/HaikeiLabs/kei-connector-contracts/generate/skill"
)

func main() {
	connector := flag.String("connector", "", "Connector ID (e.g. linear, github, google_drive). Required when the filename does not match '*.tool-definition.v1.json'.")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: genskill [-connector <id>] <tool-definition-file>")
		os.Exit(2)
	}

	defPath := args[0]
	f, err := os.Open(defPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	opts := skill.Options{ConnectorID: *connector}
	if opts.ConnectorID == "" {
		// Derive from filename: linear.tool-definition.v1.json → linear
		base := filepath.Base(defPath)
		opts.ConnectorID = strings.TrimSuffix(base, ".tool-definition.v1.json")
	}

	out, err := skill.Generate(f, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	os.Stdout.Write(out)
}
