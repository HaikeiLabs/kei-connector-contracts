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

// Package manifestdiff compares two connector manifests and reports the
// breaking changes of proposal 0002, section 6 (and the semantic-mapping
// rule of section 8.3). buf breaking checks the messages; option values are
// only visible in the manifest, so this is the gate for them.
package manifestdiff

import (
	"fmt"
	"slices"
	"sort"

	"github.com/HaikeiLabs/kei-connector-contracts/tools/protoc-gen-kei/internal/manifest"
)

// Finding is one breaking change.
type Finding struct {
	Op     string
	Reason string
}

func (f Finding) String() string { return f.Op + ": " + f.Reason }

// Diff returns the breaking changes from old to new, sorted by operation id
// and then reason. For every operation in old:
//
//   - removing it is breaking;
//   - so is changing its class from read to write, making it destructive,
//     adding a required scope, removing a connection it accepts, changing
//     its resource template or its request or response type, and changing
//     its semantic mapping without raising the mapping version (or lowering
//     the version).
//
// An operation that old already marks deprecated, with a replaced_by that
// exists in new, may change or go away: that is the rename path. A
// replaced_by in new that names no operation of new is breaking.
func Diff(old, new *manifest.Manifest) []Finding {
	oldCaps, newCaps := index(old), index(new)
	var out []Finding
	add := func(op, format string, args ...any) {
		out = append(out, Finding{Op: op, Reason: fmt.Sprintf(format, args...)})
	}
	for id, o := range oldCaps {
		if o.Deprecated && o.ReplacedBy != "" {
			if _, ok := newCaps[o.ReplacedBy]; ok {
				continue
			}
		}
		n, ok := newCaps[id]
		if !ok {
			add(id, "removed")
			continue
		}
		if o.Class == "read" && n.Class == "write" {
			add(id, "class changed from read to write")
		}
		if !o.Destructive && n.Destructive {
			add(id, "became destructive")
		}
		for _, s := range n.Scopes {
			if !slices.Contains(o.Scopes, s) {
				add(id, "adds required scope %s", s)
			}
		}
		for _, c := range o.Connections {
			if !slices.Contains(n.Connections, c) {
				add(id, "removes connection %s", c)
			}
		}
		if o.ResourceTemplate != n.ResourceTemplate {
			add(id, "resource_template changed from %s to %s", o.ResourceTemplate, n.ResourceTemplate)
		}
		if o.Request != n.Request {
			add(id, "request type changed")
		}
		if o.Response != n.Response {
			add(id, "response type changed")
		}
		switch {
		case n.Semantics.Version < o.Semantics.Version:
			add(id, "semantics version lowered from %d to %d", o.Semantics.Version, n.Semantics.Version)
		case n.Semantics.Version == o.Semantics.Version && n.Semantics != o.Semantics:
			add(id, "semantics changed without a version bump")
		}
	}
	for id, n := range newCaps {
		if n.ReplacedBy == "" {
			continue
		}
		if _, ok := newCaps[n.ReplacedBy]; !ok {
			add(id, "replaced_by names unknown operation %s", n.ReplacedBy)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Op != out[j].Op {
			return out[i].Op < out[j].Op
		}
		return out[i].Reason < out[j].Reason
	})
	return out
}

func index(m *manifest.Manifest) map[string]manifest.Capability {
	out := map[string]manifest.Capability{}
	for _, c := range m.Connectors {
		for _, cap := range c.Capabilities {
			out[cap.ID] = cap
		}
	}
	return out
}
