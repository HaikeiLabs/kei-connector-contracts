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

package manifestdiff_test

import (
	"reflect"
	"testing"

	"github.com/HaikeiLabs/kei-connector-contracts/tools/protoc-gen-kei/internal/manifest"
	"github.com/HaikeiLabs/kei-connector-contracts/tools/protoc-gen-kei/internal/manifestdiff"
)

func capability(id string) manifest.Capability {
	return manifest.Capability{
		ID: id, Class: "read", Scopes: []string{"things:read"}, Connections: []string{"pat", "user_oauth"},
		ResourceTemplate: "spaces/{space}",
		Semantics:        manifest.Semantics{Action: "read", ResourceType: "thing", Version: 1},
		Request:          "kei.acme.v1.GetThingRequest", Response: "kei.acme.v1.GetThingResponse",
	}
}

func doc(caps ...manifest.Capability) *manifest.Manifest {
	return &manifest.Manifest{Schema: manifest.Schema, ContractsVersion: "v1.0.0", Connectors: []manifest.Connector{{ID: "acme", Capabilities: caps}}}
}

func TestDiff(t *testing.T) {
	base := capability("acme.things.get")
	with := func(f func(*manifest.Capability)) manifest.Capability {
		c := capability("acme.things.get")
		f(&c)
		return c
	}

	cases := []struct {
		name     string
		old, new *manifest.Manifest
		want     []string
	}{
		{"unchanged", doc(base), doc(base), nil},
		{"added op", doc(base), doc(base, capability("acme.things.list")), nil},
		{"added optional connection", doc(base), doc(with(func(c *manifest.Capability) { c.Connections = append(c.Connections, "api_key") })), nil},
		{"removed op", doc(base), doc(), []string{"acme.things.get: removed"}},
		{"read to write", doc(base), doc(with(func(c *manifest.Capability) { c.Class = "write"; c.Semantics.Action = "update" })),
			[]string{"acme.things.get: class changed from read to write", "acme.things.get: semantics changed without a version bump"}},
		{"write to read is not breaking", doc(with(func(c *manifest.Capability) { c.Class = "write" })), doc(base), nil},
		{"set destructive", doc(base), doc(with(func(c *manifest.Capability) { c.Destructive = true })), []string{"acme.things.get: became destructive"}},
		{"added scope", doc(base), doc(with(func(c *manifest.Capability) { c.Scopes = append(c.Scopes, "admin") })), []string{"acme.things.get: adds required scope admin"}},
		{"removed connection", doc(base), doc(with(func(c *manifest.Capability) { c.Connections = []string{"pat"} })), []string{"acme.things.get: removes connection user_oauth"}},
		{"resource template", doc(base), doc(with(func(c *manifest.Capability) { c.ResourceTemplate = "things" })), []string{"acme.things.get: resource_template changed from spaces/{space} to things"}},
		{"request type", doc(base), doc(with(func(c *manifest.Capability) { c.Request = "kei.acme.v2.GetThingRequest" })), []string{"acme.things.get: request type changed"}},
		{"semantics without bump", doc(base), doc(with(func(c *manifest.Capability) { c.Semantics.ResourceType = "item" })), []string{"acme.things.get: semantics changed without a version bump"}},
		{"semantics with bump", doc(base), doc(with(func(c *manifest.Capability) { c.Semantics.ResourceType = "item"; c.Semantics.Version = 2 })), nil},
		{"semantics version lowered", doc(with(func(c *manifest.Capability) { c.Semantics.Version = 2 })), doc(base), []string{"acme.things.get: semantics version lowered from 2 to 1"}},
		{"rename with deprecated and replaced_by",
			doc(base),
			doc(with(func(c *manifest.Capability) { c.Deprecated = true; c.ReplacedBy = "acme.things.fetch" }), capability("acme.things.fetch")), nil},
		{"replaced_by names an unknown op",
			doc(base),
			doc(with(func(c *manifest.Capability) { c.Deprecated = true; c.ReplacedBy = "acme.things.fetch" })),
			[]string{"acme.things.get: replaced_by names unknown operation acme.things.fetch"}},
		{"removing a deprecated op with a replacement",
			doc(with(func(c *manifest.Capability) { c.Deprecated = true; c.ReplacedBy = "acme.things.fetch" }), capability("acme.things.fetch")),
			doc(capability("acme.things.fetch")), nil},
		{"removing a deprecated op without a replacement",
			doc(with(func(c *manifest.Capability) { c.Deprecated = true })),
			doc(), []string{"acme.things.get: removed"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got []string
			for _, f := range manifestdiff.Diff(tc.old, tc.new) {
				got = append(got, f.String())
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("Diff =\n%q\nwant\n%q", got, tc.want)
			}
		})
	}
}
