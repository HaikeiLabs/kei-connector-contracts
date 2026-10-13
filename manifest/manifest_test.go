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

package manifest

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// goldenTools is the HAI-484 golden vector (ADR-036 section 5.2).
var goldenTools = []json.RawMessage{
	[]byte(`{"name": "search_dashboards", "title": "Search", "description": "Search dashboards by title",
		"inputSchema": {"type": "object", "properties": {"query": {"type": "string"}}, "required": ["query"]},
		"annotations": {"readOnlyHint": true}}`),
	[]byte(`{"name": "get_dashboard",
		"inputSchema": {"type": "object", "properties": {"uid": {"type": "string"}}, "required": ["uid"], "additionalProperties": false}}`),
}

func TestGoldenVectorCanonical(t *testing.T) {
	canon, err := Canonicalize(goldenTools)
	if err != nil {
		t.Fatal(err)
	}
	var gd *Tool
	for i := range canon {
		if canon[i].Name == "get_dashboard" {
			gd = &canon[i]
		}
	}
	if gd == nil {
		t.Fatal("get_dashboard not found in canonical tools")
	}
	b, err := jcs(gd.canonicalObject())
	if err != nil {
		t.Fatal(err)
	}
	want := `{"description":"","inputSchema":{"additionalProperties":false,"properties":{"uid":{"type":"string"}},"required":["uid"],"type":"object"},"name":"get_dashboard"}`
	if string(b) != want {
		t.Errorf("canonical get_dashboard =\n%s\nwant\n%s", b, want)
	}
}

func TestGoldenVectorToolHashes(t *testing.T) {
	got, err := ToolHash(goldenTools[1]) // get_dashboard
	if err != nil {
		t.Fatal(err)
	}
	if want := "sha256:9e84a27116887ce94c79f25184ee07da9b9dd0288847be9e0ee885fdc61ec537"; got != want {
		t.Errorf("tool_hash(get_dashboard) = %s, want %s", got, want)
	}
	got, err = ToolHash(goldenTools[0]) // search_dashboards
	if err != nil {
		t.Fatal(err)
	}
	if want := "sha256:364d2fd75a2d8581c118e66eae51057468a5848031e1586d0f8566b34e531499"; got != want {
		t.Errorf("tool_hash(search_dashboards) = %s, want %s", got, want)
	}
}

func TestGoldenVectorManifestHash(t *testing.T) {
	got, err := Hash(goldenTools)
	if err != nil {
		t.Fatal(err)
	}
	if want := "sha256:a94748f3bd7e1bc1e46ff952e802378a98249adfd5fcd2677ff48bd0df848dbb"; got != want {
		t.Errorf("manifest_hash = %s, want %s", got, want)
	}
}

func TestHashIsOrderIndependent(t *testing.T) {
	a, err := Hash(goldenTools)
	if err != nil {
		t.Fatal(err)
	}
	reversed := []json.RawMessage{goldenTools[1], goldenTools[0]}
	b, err := Hash(reversed)
	if err != nil {
		t.Fatal(err)
	}
	if a != b {
		t.Errorf("Hash depends on input order: %s != %s", a, b)
	}
}

func TestDescriptionChangeChangesHash(t *testing.T) {
	base, err := Hash(goldenTools)
	if err != nil {
		t.Fatal(err)
	}
	changed := []json.RawMessage{
		[]byte(`{"name": "get_dashboard", "description": "changed", "inputSchema": {"type": "object"}}`),
		goldenTools[0],
	}
	got, err := Hash(changed)
	if err != nil {
		t.Fatal(err)
	}
	if got == base {
		t.Error("changing description did not change the manifest hash")
	}
}

func TestInputSchemaChangeChangesHash(t *testing.T) {
	base, err := Hash(goldenTools)
	if err != nil {
		t.Fatal(err)
	}
	changed := []json.RawMessage{
		[]byte(`{"name": "get_dashboard", "inputSchema": {"type": "object", "properties": {"uid": {"type": "integer"}}}}`),
		goldenTools[0],
	}
	got, err := Hash(changed)
	if err != nil {
		t.Fatal(err)
	}
	if got == base {
		t.Error("changing inputSchema did not change the manifest hash")
	}
}

func TestInvalidManifestsFailClosed(t *testing.T) {
	cases := map[string][]json.RawMessage{
		"duplicate name": {
			[]byte(`{"name": "a", "inputSchema": {}}`),
			[]byte(`{"name": "a", "inputSchema": {}}`),
		},
		"name with space": {
			[]byte(`{"name": "has space", "inputSchema": {}}`),
		},
		"name too long": {
			[]byte(`{"name": "` + strings.Repeat("a", 129) + `", "inputSchema": {}}`),
		},
		"missing name": {
			[]byte(`{"inputSchema": {}}`),
		},
		"non-string name": {
			[]byte(`{"name": 1, "inputSchema": {}}`),
		},
		"missing inputSchema": {
			[]byte(`{"name": "a"}`),
		},
		"null inputSchema": {
			[]byte(`{"name": "a", "inputSchema": null}`),
		},
		"non-object inputSchema": {
			[]byte(`{"name": "a", "inputSchema": "x"}`),
		},
		"non-string description": {
			[]byte(`{"name": "a", "description": 1, "inputSchema": {}}`),
		},
		"tool not an object": {
			[]byte(`"just a string"`),
		},
	}
	for name, tools := range cases {
		if _, err := Canonicalize(tools); err == nil {
			t.Errorf("%s: accepted, want error", name)
		}
	}
	// more than 512 tools
	many := make([]json.RawMessage, maxTools+1)
	for i := range many {
		many[i] = json.RawMessage(fmt.Sprintf(`{"name": "t%d", "inputSchema": {}}`, i))
	}
	if _, err := Canonicalize(many); err == nil {
		t.Error("513 tools: accepted, want error")
	}
	// more than 4 MiB
	big := json.RawMessage(`{"name": "big", "inputSchema": {"pad": "` + strings.Repeat("x", maxManifestBytes) + `"}}`)
	if _, err := Canonicalize([]json.RawMessage{big}); err == nil {
		t.Error("over-4MiB tool: accepted, want error")
	}
}

// The number serializer must match RFC 8785 (ECMAScript) exactly, because the
// hash is computed in more than one language.
func TestJCSNumberSerialization(t *testing.T) {
	cases := []struct {
		in   float64
		want string
	}{
		{0, "0"},
		{math.Copysign(0, -1), "0"},
		{1, "1"},
		{123.45, "123.45"},
		{0.000001, "0.000001"},
		{1e-7, "1e-7"},
		{1e20, "100000000000000000000"},
		{1e21, "1e+21"},
		{9007199254740992, "9007199254740992"},
		{5e-324, "5e-324"},
		{1.7976931348623157e308, "1.7976931348623157e+308"},
		{9.999999999999997e-7, "9.999999999999997e-7"},
		{333333333.3333332, "333333333.3333332"},
		{-123.45, "-123.45"},
	}
	for _, c := range cases {
		got, err := jcsNumber(c.in)
		if err != nil {
			t.Errorf("jcsNumber(%v) error: %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("jcsNumber(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestJCSNumberRejectsNaNInf(t *testing.T) {
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := jcsNumber(v); err == nil {
			t.Errorf("jcsNumber(%v): accepted, want error", v)
		}
	}
}

// Object members are sorted by UTF-16 code units (RFC 8785), which differs
// from UTF-8 byte order for non-BMP characters.
func TestJCSKeySortingIsUTF16(t *testing.T) {
	keys := []string{"\u20ac", "\r", "\ufb33", "1", "\U0001F600", "\u0080", "\u00f6"}
	sort.Slice(keys, func(i, j int) bool { return jcsLess(keys[i], keys[j]) })
	want := []string{"\r", "1", "\u0080", "\u00f6", "\u20ac", "\U0001F600", "\ufb33"}
	if !reflect.DeepEqual(keys, want) {
		t.Errorf("sorted keys = %q, want %q", keys, want)
	}
}

// Strings are escaped per RFC 8785: only ASCII control characters (except the
// five predefined), backslash, and quote are escaped; everything else is
// emitted as-is (no HTML escaping, no \uXXXX for non-ASCII).
func TestJCSStringEscaping(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{`a"b`, `"a\"b"`},
		{`a\b`, `"a\\b"`},
		{"a\nb", `"a\nb"`},
		{"a\tb", `"a\tb"`},
		{"a\u0000b", `"a\u0000b"`},
		{"a\u001fb", `"a\u001fb"`},
		{"a\u000bb", `"a\u000bb"`},
		{"a&b<c>", `"a&b<c>"`},
		{"a€b", `"a€b"`},
	}
	for _, c := range cases {
		if got := jcsString(c.in); got != c.want {
			t.Errorf("jcsString(%q) = %s, want %s", c.in, got, c.want)
		}
	}
}
