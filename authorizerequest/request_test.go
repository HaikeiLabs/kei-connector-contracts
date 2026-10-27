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

package authorizerequest

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestCanonicalExamples(t *testing.T) {
	for _, tc := range []struct {
		path       string
		wantCaller bool
	}{
		{"../schemas/examples/authorize-request/tool.v1.json", false},
		{"../schemas/examples/authorize-request/tool-with-caller.v1.json", true},
		{"../schemas/examples/authorize-request/connector-with-caller.v1.json", true},
	} {
		raw, err := os.ReadFile(tc.path)
		if err != nil {
			t.Fatal(err)
		}
		req, err := Decode(raw)
		if err != nil {
			t.Errorf("Decode(%s): %v", tc.path, err)
			continue
		}
		if (req.Caller != nil) != tc.wantCaller {
			t.Errorf("%s: caller present = %v, want %v", tc.path, req.Caller != nil, tc.wantCaller)
		}
	}
}

// A request without caller metadata (every request before this field existed)
// stays valid: the caller is simply unknown.
func TestCallerIsOptional(t *testing.T) {
	req, err := Decode([]byte(`{"schema":"kei.authorize-request/v1","tool":"web_search"}`))
	if err != nil {
		t.Fatal(err)
	}
	if req.Caller != nil {
		t.Errorf("caller = %+v, want nil", req.Caller)
	}
}

func TestDecodeRejectsInvalidRequests(t *testing.T) {
	base := `"schema":"kei.authorize-request/v1","tool":"web_search"`
	for name, input := range map[string]string{
		"wrong schema":                   `{"schema":"kei.authorize-request/v2","tool":"x"}`,
		"unknown top-level field":        `{` + base + `,"runtime_token":"fake-secret-marker"}`,
		"neither tool nor connector":     `{"schema":"kei.authorize-request/v1"}`,
		"both tool and connector":        `{` + base + `,"connector":{"connector_id":"github","capability":"read"}}`,
		"connector without capability":   `{"schema":"kei.authorize-request/v1","connector":{"connector_id":"github"}}`,
		"unknown caller field":           `{` + base + `,"caller":{"origin":"interactive","message":"hello"}}`,
		"origin outside the closed set":  `{` + base + `,"caller":{"origin":"batch"}}`,
		"empty origin":                   `{` + base + `,"caller":{"origin":""}}`,
		"harness with uppercase":         `{` + base + `,"caller":{"harness":"OpenWebUI"}}`,
		"harness with a space":           `{` + base + `,"caller":{"harness":"open webui"}}`,
		"harness too long":               `{` + base + `,"caller":{"harness":"` + strings.Repeat("a", 65) + `"}}`,
		"session ref with whitespace":    `{` + base + `,"caller":{"session_ref":"what is the plan"}}`,
		"session ref too long":           `{` + base + `,"caller":{"session_ref":"` + strings.Repeat("a", 129) + `"}}`,
		"trailing data after the object": `{` + base + `}{}`,
	} {
		t.Run(name, func(t *testing.T) {
			_, err := Decode([]byte(input))
			if err == nil {
				t.Fatal("Decode accepted an invalid request")
			}
			if strings.Contains(err.Error(), "fake-secret-marker") {
				t.Errorf("error echoes the submitted value: %v", err)
			}
		})
	}
}

func TestCallerAtLimits(t *testing.T) {
	ref := strings.Repeat("a", 128)
	harness := "a" + strings.Repeat("b", 63)
	if _, err := Decode([]byte(`{"schema":"kei.authorize-request/v1","tool":"x","caller":{"harness":"` + harness + `","origin":"automation","session_ref":"` + ref + `"}}`)); err != nil {
		t.Fatalf("Decode at the limits: %v", err)
	}
}

// The JSON Schema and the Go validation are two statements of one contract;
// this keeps their closed sets from drifting.
func TestSchemaMatchesValidation(t *testing.T) {
	raw, err := os.ReadFile("../schemas/authorize-request.v1.schema.json")
	if err != nil {
		t.Fatal(err)
	}
	var schema struct {
		ID         string `json:"$id"`
		Properties map[string]json.RawMessage
		Defs       struct {
			Caller struct {
				AdditionalProperties bool `json:"additionalProperties"`
				Properties           struct {
					Origin struct {
						Enum []string `json:"enum"`
					} `json:"origin"`
					Harness    struct{ Pattern string } `json:"harness"`
					SessionRef struct{ Pattern string } `json:"session_ref"`
				} `json:"properties"`
			} `json:"caller"`
		} `json:"$defs"`
	}
	if err := json.Unmarshal(raw, &schema); err != nil {
		t.Fatal(err)
	}
	if _, ok := schema.Properties["caller"]; !ok {
		t.Error("schema has no caller property")
	}
	c := schema.Defs.Caller
	if c.AdditionalProperties {
		t.Error("caller must be closed (additionalProperties false)")
	}
	if got := strings.Join(c.Properties.Origin.Enum, ","); got != OriginInteractive+","+OriginAutomation {
		t.Errorf("origin enum = %q", got)
	}
	if c.Properties.Harness.Pattern != harnessPattern.String() {
		t.Errorf("harness pattern = %q, want %q", c.Properties.Harness.Pattern, harnessPattern.String())
	}
	if c.Properties.SessionRef.Pattern != sessionRefPattern.String() {
		t.Errorf("session_ref pattern = %q, want %q", c.Properties.SessionRef.Pattern, sessionRefPattern.String())
	}
}
