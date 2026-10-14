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

package canonical_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/HaikeiLabs/kei-connector-contracts/canonical"
	"github.com/HaikeiLabs/kei-connector-contracts/canonical/canonicaltest"
)

// TestResourceGolden runs the shared golden table: every accepted input
// returns exactly its canonical form, and every rejected input errors with
// ErrInvalidResource.
func TestResourceGolden(t *testing.T) {
	for _, c := range canonicaltest.ResourceCases() {
		got, err := canonical.Resource(c.Input)
		if c.Error {
			if err == nil || !errors.Is(err, canonical.ErrInvalidResource) {
				t.Errorf("%s: Resource(%q) = %q, %v; want ErrInvalidResource", c.Name, c.Input, got, err)
			}
			continue
		}
		if err != nil || got != c.Want {
			t.Errorf("%s: Resource(%q) = %q, %v; want %q", c.Name, c.Input, got, err, c.Want)
		}
	}
}

// TestResourceErrorsNeverEchoInput keeps the submitted URL, which may carry
// user info credentials, out of the error.
func TestResourceErrorsNeverEchoInput(t *testing.T) {
	const marker = "FAKE-secret-marker"
	for _, in := range []string{
		"https://user:" + marker + "@mcp.example.com/mcp",
		"http://" + marker + ".example.com/mcp",
		"https://mcp.example.com/mcp?token=" + marker,
		"https://mcp.example.com/mcp#" + marker,
		"https://" + marker + " .example.com/mcp",
	} {
		_, err := canonical.Resource(in)
		if err == nil {
			t.Fatalf("Resource(%q) accepted", in)
		}
		if strings.Contains(err.Error(), marker) {
			t.Errorf("error %q echoes the submitted value", err)
		}
	}
}
