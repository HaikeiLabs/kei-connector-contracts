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

package mockapi

import (
	"fmt"
	"strings"
	"time"
)

func sleepFor(seconds int) { time.Sleep(time.Duration(seconds) * time.Second) }

// vary rewrites a decoded fixture with a deterministic synthetic offset:
// numeric ids gain n, varying strings gain a "-m<n>" suffix, recursively.
// Only fields a decoder treats as free-form vary — enum-shaped keys (state,
// status, …) are left alone so the generated clients still decode the page.
//
// The offset n must be derived deterministically (seed, page, item index):
// the same request always yields the same page.
func vary(v any, n int) any {
	switch t := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(t))
		for k, val := range t {
			if varyableKey(k) {
				out[k] = varyValue(val, n)
			} else {
				out[k] = vary(val, n)
			}
		}
		return out
	case []any:
		out := make([]any, len(t))
		for i, item := range t {
			out[i] = vary(item, n+i*7)
		}
		return out
	default:
		return v
	}
}

// varyValue applies the offset to one scalar.
func varyValue(v any, n int) any {
	switch t := v.(type) {
	case string:
		return t + "-m" + fmt.Sprint(n)
	case float64:
		return t + float64(n)
	default:
		return v
	}
}

// varyableKey reports whether the field is free-form enough to vary: ids,
// urls and display strings. Enum-shaped keys (state, status, type, …) pass
// through untouched, so the generated clients still decode the page.
func varyableKey(key string) bool {
	lk := strings.ToLower(key)
	switch lk {
	case "id", "url", "title", "name", "sha", "text", "body", "description", "email", "slug", "subject", "summary", "key":
		return true
	}
	return strings.Contains(lk, "_id") || strings.Contains(lk, "_url")
}

// growList lengthens the first top-level list by one varied item, so a
// generated page's length differs from the fixture's.
func growList(doc any) any {
	obj, ok := doc.(map[string]any)
	if !ok {
		return doc
	}
	for k, v := range obj {
		list, ok := v.([]any)
		if !ok || len(list) == 0 {
			continue
		}
		grown := append(append([]any{}, list...), vary(list[0], 991))
		obj[k] = grown
		return obj
	}
	return obj
}
