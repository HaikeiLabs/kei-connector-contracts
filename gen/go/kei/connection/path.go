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

package connection

import (
	"fmt"
	"net/url"
	"slices"
	"strings"
)

// RenderPath fills a path template such as
// "/repos/{owner}/{repository}" from values.
//
// Every value is percent-encoded as a single path segment, so no value can
// change the shape of the path: "/" becomes %2F, "%" becomes %25, "?" and
// "#" are encoded, and a value made only of dots ("." or "..") is encoded as
// %2E so it is never read as a relative segment. A missing or empty value,
// an unterminated or empty variable, or a template that does not start with
// "/" is an error.
//
// A variable named in multiSegment (a param with multi_segment, such as a
// file path or a Git ref) keeps "/" between its segments and escapes each
// segment; an empty, "." or ".." segment is an error.
func RenderPath(template string, values map[string]string, multiSegment ...string) (string, error) {
	if !strings.HasPrefix(template, "/") {
		return "", fmt.Errorf("connection: path template %q must start with /", template)
	}
	var b strings.Builder
	rest := template
	for {
		open := strings.IndexByte(rest, '{')
		close := strings.IndexByte(rest, '}')
		if open < 0 {
			if close >= 0 {
				return "", fmt.Errorf("connection: path template %q has an unmatched }", template)
			}
			b.WriteString(rest)
			return b.String(), nil
		}
		if close >= 0 && close < open {
			return "", fmt.Errorf("connection: path template %q has an unmatched }", template)
		}
		end := strings.IndexByte(rest[open:], '}')
		if end < 0 {
			return "", fmt.Errorf("connection: path template %q has an unterminated variable", template)
		}
		name := rest[open+1 : open+end]
		if name == "" || strings.ContainsAny(name, "{/") {
			return "", fmt.Errorf("connection: path template %q has an invalid variable", template)
		}
		v, ok := values[name]
		if !ok {
			return "", fmt.Errorf("connection: path variable %q has no value", name)
		}
		if v == "" {
			return "", fmt.Errorf("connection: path variable %q is empty", name)
		}
		b.WriteString(rest[:open])
		if slices.Contains(multiSegment, name) {
			segs := strings.Split(v, "/")
			for i, s := range segs {
				if s == "" || s == "." || s == ".." {
					return "", fmt.Errorf("connection: path variable %q has an empty, . or .. segment", name)
				}
				segs[i] = url.PathEscape(s)
			}
			b.WriteString(strings.Join(segs, "/"))
		} else {
			b.WriteString(escapeSegment(v))
		}
		rest = rest[open+end+1:]
	}
}

// escapeSegment percent-encodes v as one path segment.
func escapeSegment(v string) string {
	if strings.Trim(v, ".") == "" {
		return strings.Repeat("%2E", len(v))
	}
	return url.PathEscape(v)
}
