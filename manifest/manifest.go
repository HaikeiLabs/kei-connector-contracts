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

// Package manifest computes the canonical form and hash of an MCP
// tool manifest, per ADR-036 section 5.2. The hash is computed over the
// JSON Canonicalization Scheme (RFC 8785) bytes so that every language that
// computes it agrees: the catalog, the runtime, the CLI and the console all
// use this package unchanged.
package manifest

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf16"
)

const (
	// maxTools is the maximum number of tools in one manifest.
	maxTools = 512
	// maxManifestBytes is the maximum total size of the tools in one manifest.
	maxManifestBytes = 4 * 1024 * 1024
)

// namePattern is the character set a tool name must match (ADR-036 3.5).
var namePattern = regexp.MustCompile(`^[A-Za-z0-9_.-]{1,128}$`)

// Tool is the canonical form of one upstream tool: name, description, and
// inputSchema. Every other member (title, annotations, outputSchema, icons,
// _meta, and any unknown member) is dropped.
type Tool struct {
	Name        string
	Description string
	InputSchema map[string]interface{}
}

// Canonicalize reduces each received tool to its canonical object and
// validates the manifest. tools is the concatenation of the tools arrays of
// every page of an upstream tools/list.
func Canonicalize(tools []json.RawMessage) ([]Tool, error) {
	if len(tools) > maxTools {
		return nil, fmt.Errorf("manifest: %d tools exceeds the maximum of %d", len(tools), maxTools)
	}
	total := 0
	for _, t := range tools {
		total += len(t)
	}
	if total > maxManifestBytes {
		return nil, fmt.Errorf("manifest: %d bytes exceeds the maximum of %d", total, maxManifestBytes)
	}
	seen := make(map[string]bool, len(tools))
	out := make([]Tool, 0, len(tools))
	for _, raw := range tools {
		tool, err := canonicalTool(raw)
		if err != nil {
			return nil, err
		}
		if seen[tool.Name] {
			return nil, fmt.Errorf("manifest: duplicate tool name %q", tool.Name)
		}
		seen[tool.Name] = true
		out = append(out, tool)
	}
	return out, nil
}

// ToolHash returns the sha256: tool_hash of one received tool.
func ToolHash(raw json.RawMessage) (string, error) {
	tool, err := canonicalTool(raw)
	if err != nil {
		return "", err
	}
	return hashValue(tool.canonicalObject())
}

// Hash returns the sha256: manifest_hash of the tools: the JCS bytes of the
// array of canonical objects, sorted by name in ascending byte order.
func Hash(tools []json.RawMessage) (string, error) {
	canon, err := Canonicalize(tools)
	if err != nil {
		return "", err
	}
	sort.Slice(canon, func(i, j int) bool { return canon[i].Name < canon[j].Name })
	arr := make([]interface{}, len(canon))
	for i, t := range canon {
		arr[i] = t.canonicalObject()
	}
	return hashValue(arr)
}

func canonicalTool(raw json.RawMessage) (Tool, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		return Tool{}, fmt.Errorf("manifest: tool is not a JSON object: %w", err)
	}
	name, ok := m["name"].(string)
	if !ok {
		return Tool{}, errors.New("manifest: tool name is missing or not a string")
	}
	if !namePattern.MatchString(name) {
		return Tool{}, fmt.Errorf("manifest: tool name %q does not match %s", name, namePattern.String())
	}
	description := ""
	if d, present := m["description"]; present {
		ds, ok := d.(string)
		if !ok {
			return Tool{}, errors.New("manifest: tool description is not a string")
		}
		description = ds
	}
	schema, present := m["inputSchema"]
	if !present {
		return Tool{}, errors.New("manifest: tool inputSchema is missing")
	}
	schemaMap, ok := schema.(map[string]interface{})
	if !ok {
		return Tool{}, errors.New("manifest: tool inputSchema is not an object")
	}
	return Tool{Name: name, Description: description, InputSchema: schemaMap}, nil
}

func (t Tool) canonicalObject() map[string]interface{} {
	return map[string]interface{}{
		"name":        t.Name,
		"description": t.Description,
		"inputSchema": t.InputSchema,
	}
}

func hashValue(v interface{}) (string, error) {
	b, err := jcs(v)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

// jcs serializes v as RFC 8785 (JSON Canonicalization Scheme) bytes.
func jcs(v interface{}) ([]byte, error) {
	var b strings.Builder
	if err := jcsWrite(&b, v); err != nil {
		return nil, err
	}
	return []byte(b.String()), nil
}

func jcsWrite(b *strings.Builder, v interface{}) error {
	switch t := v.(type) {
	case nil:
		b.WriteString("null")
	case bool:
		if t {
			b.WriteString("true")
		} else {
			b.WriteString("false")
		}
	case float64:
		s, err := jcsNumber(t)
		if err != nil {
			return err
		}
		b.WriteString(s)
	case string:
		b.WriteString(jcsString(t))
	case []interface{}:
		b.WriteByte('[')
		for i, e := range t {
			if i > 0 {
				b.WriteByte(',')
			}
			if err := jcsWrite(b, e); err != nil {
				return err
			}
		}
		b.WriteByte(']')
	case map[string]interface{}:
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool { return jcsLess(keys[i], keys[j]) })
		b.WriteByte('{')
		for i, k := range keys {
			if i > 0 {
				b.WriteByte(',')
			}
			b.WriteString(jcsString(k))
			b.WriteByte(':')
			if err := jcsWrite(b, t[k]); err != nil {
				return err
			}
		}
		b.WriteByte('}')
	default:
		return fmt.Errorf("manifest: unsupported JCS value type %T", v)
	}
	return nil
}

// jcsNumber serializes f per RFC 8785 section 3.2.2.3 (ECMAScript
// Number.prototype.toString, with the -0 -> 0 enhancement).
func jcsNumber(f float64) (string, error) {
	if math.IsNaN(f) {
		return "", errors.New("manifest: NaN is not a valid JSON number")
	}
	if math.IsInf(f, 0) {
		return "", errors.New("manifest: Infinity is not a valid JSON number")
	}
	if f == 0 {
		return "0", nil
	}
	neg := f < 0
	s := strconv.FormatFloat(math.Abs(f), 'e', -1, 64) // e.g. "1.2345e+02"
	eIdx := strings.IndexByte(s, 'e')
	mantissa := s[:eIdx]
	exp, err := strconv.Atoi(s[eIdx+1:])
	if err != nil {
		return "", err
	}
	dot := strings.IndexByte(mantissa, '.')
	digits := mantissa
	if dot >= 0 {
		digits = mantissa[:dot] + mantissa[dot+1:]
	}
	body := formatJCSNumber(digits, exp+1)
	if neg {
		body = "-" + body
	}
	return body, nil
}

// formatJCSNumber renders the significant digits (no leading zeros, last
// digit non-zero) with the decimal point at position n (the number of digits
// to the left of the point), choosing fixed or scientific notation exactly as
// ECMAScript does.
func formatJCSNumber(digits string, n int) string {
	e := len(digits)
	switch {
	case 1 <= n && n <= e && e <= 21:
		if n == e {
			return digits
		}
		return digits[:n] + "." + digits[n:]
	case 0 < n && n <= 21:
		return digits + strings.Repeat("0", n-e)
	case -6 < n && n <= 0:
		return "0." + strings.Repeat("0", -n) + digits
	default:
		exp := n - 1
		sign := "+"
		if exp < 0 {
			sign = "-"
		}
		absExp := exp
		if exp < 0 {
			absExp = -exp
		}
		frac := ""
		if e > 1 {
			frac = "." + digits[1:]
		}
		return string(digits[0]) + frac + "e" + sign + strconv.Itoa(absExp)
	}
}

// jcsString serializes s per RFC 8785 section 3.2.2.2: only ASCII control
// characters (except the five predefined escapes), backslash, and quote are
// escaped; every other code point is emitted as-is (UTF-8, no HTML escaping).
func jcsString(s string) string {
	var b strings.Builder
	b.WriteByte('"')
	for _, r := range s {
		switch {
		case r == '"':
			b.WriteString(`\"`)
		case r == '\\':
			b.WriteString(`\\`)
		case r == '\b':
			b.WriteString(`\b`)
		case r == '\t':
			b.WriteString(`\t`)
		case r == '\n':
			b.WriteString(`\n`)
		case r == '\f':
			b.WriteString(`\f`)
		case r == '\r':
			b.WriteString(`\r`)
		case r < 0x20:
			b.WriteString(fmt.Sprintf(`\u%04x`, r))
		default:
			b.WriteRune(r)
		}
	}
	b.WriteByte('"')
	return b.String()
}

// jcsLess reports whether a sorts before b in JCS object-member order:
// ascending UTF-16 code units, the shorter prefix first.
func jcsLess(a, b string) bool {
	au := utf16.Encode([]rune(a))
	bu := utf16.Encode([]rune(b))
	for i := 0; i < len(au) && i < len(bu); i++ {
		if au[i] != bu[i] {
			return au[i] < bu[i]
		}
	}
	return len(au) < len(bu)
}
