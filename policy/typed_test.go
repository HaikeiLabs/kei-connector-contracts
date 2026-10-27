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

package policy

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

type typedWant struct {
	Error   bool   `json:"error"`
	Pattern string `json:"pattern"`
}

func TestTypedTargetVectors(t *testing.T) {
	var suite struct {
		Cases []struct {
			Name   string      `json:"name"`
			Target TypedTarget `json:"target"`
			Want   typedWant   `json:"want"`
		} `json:"cases"`
	}
	loadVectors(t, "typed-target.json", &suite)
	if len(suite.Cases) == 0 {
		t.Fatal("no cases")
	}
	for _, tc := range suite.Cases {
		t.Run(tc.Name, func(t *testing.T) {
			err := tc.Target.Validate()
			if tc.Want.Error {
				if err == nil {
					t.Fatalf("Validate() = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("Validate() = %v", err)
			}
			if got := tc.Target.Pattern(); got != tc.Want.Pattern {
				t.Errorf("Pattern() = %q, want %q", got, tc.Want.Pattern)
			}
		})
	}
}

func TestTypedPrincipalVectors(t *testing.T) {
	var suite struct {
		Cases []struct {
			Name      string         `json:"name"`
			Context   string         `json:"context"`
			Principal TypedPrincipal `json:"principal"`
			Want      typedWant      `json:"want"`
		} `json:"cases"`
	}
	loadVectors(t, "typed-principal.json", &suite)
	if len(suite.Cases) == 0 {
		t.Fatal("no cases")
	}
	for _, tc := range suite.Cases {
		t.Run(tc.Name, func(t *testing.T) {
			var err error
			switch tc.Context {
			case "policy":
				err = tc.Principal.ValidatePolicy()
			case "grant":
				err = tc.Principal.ValidateGrant()
			default:
				t.Fatalf("unknown context %q", tc.Context)
			}
			if tc.Want.Error {
				if err == nil {
					t.Fatalf("validate = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("validate = %v", err)
			}
			if got := tc.Principal.Pattern(); got != tc.Want.Pattern {
				t.Errorf("Pattern() = %q, want %q", got, tc.Want.Pattern)
			}
		})
	}
}

// The typed forms convert from and back to the patterns the existing parsers
// accept, so a target or principal survives the round trip unchanged.
func TestTypedPatternRoundTrip(t *testing.T) {
	var targets, principals struct {
		Cases []struct {
			Input string `json:"input"`
			Want  struct {
				Error bool `json:"error"`
			} `json:"want"`
		} `json:"cases"`
	}
	loadVectors(t, "parse-target.json", &targets)
	loadVectors(t, "parse-principal.json", &principals)
	for _, tc := range targets.Cases {
		typed, err := TypedTargetFromPattern(tc.Input)
		if tc.Want.Error {
			if err == nil {
				t.Errorf("TypedTargetFromPattern(%q) = nil error", tc.Input)
			}
			continue
		}
		if err != nil || typed.Pattern() != tc.Input {
			t.Errorf("TypedTargetFromPattern(%q) = %+v, %v", tc.Input, typed, err)
		}
	}
	for _, tc := range principals.Cases {
		typed, err := TypedPrincipalFromPattern(tc.Input)
		if tc.Want.Error {
			if err == nil {
				t.Errorf("TypedPrincipalFromPattern(%q) = nil error", tc.Input)
			}
			continue
		}
		// A bare name parses as a group and renders as "group:<name>".
		want := tc.Input
		if typed.Kind == "group" && want == typed.Value {
			want = "group:" + want
		}
		if err != nil || typed.Pattern() != want {
			t.Errorf("TypedPrincipalFromPattern(%q) = %+v, %v", tc.Input, typed, err)
		}
	}
}

// The JSON Schemas list exactly the kinds and autogroup names the parsers know.
func TestTypedSchemasMatchParsers(t *testing.T) {
	var target struct {
		Properties struct {
			Kind struct {
				Enum []string `json:"enum"`
			} `json:"kind"`
		} `json:"properties"`
	}
	readSchema(t, "policy-target.v1.schema.json", &target)
	wantKinds := []string{"any"}
	for scheme := range targetSchemes {
		wantKinds = append(wantKinds, scheme)
	}
	assertSameSet(t, "target kinds", target.Properties.Kind.Enum, wantKinds)

	var principal struct {
		Defs struct {
			Autogroup struct {
				Properties struct {
					Value struct {
						Enum []string `json:"enum"`
					} `json:"value"`
				} `json:"properties"`
			} `json:"autogroup"`
		} `json:"$defs"`
	}
	readSchema(t, "policy-principal.v1.schema.json", &principal)
	var wantGroups []string
	for name := range autogroups {
		wantGroups = append(wantGroups, name)
	}
	assertSameSet(t, "autogroup names", principal.Defs.Autogroup.Properties.Value.Enum, wantGroups)
}

func loadVectors(t *testing.T, name string, into any) {
	t.Helper()
	raw, err := os.ReadFile(vectorPath(name))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(raw, into); err != nil {
		t.Fatalf("%s: %v", name, err)
	}
}

func readSchema(t *testing.T, name string, into any) {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("..", "schemas", name))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(raw, into); err != nil {
		t.Fatalf("%s: %v", name, err)
	}
}

func assertSameSet(t *testing.T, what string, got, want []string) {
	t.Helper()
	sort.Strings(got)
	sort.Strings(want)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s: schema %v, parsers %v", what, got, want)
	}
}
