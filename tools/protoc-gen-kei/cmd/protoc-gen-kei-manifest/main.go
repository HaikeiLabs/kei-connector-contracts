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

// Command protoc-gen-kei-manifest generates kei.connector-manifest/v1 from
// the kei.contracts.v1 connector and operation options (proposal 0002,
// section 4.6), with the legacy capability mapping of section 8.1.
//
// Parameters: contracts_version=<tag> (required) and legacy=<path to
// legacy_capabilities.yaml, relative to the plugin's working directory>.
// It runs the same validation pass as protoc-gen-kei-ops first.
package main

import (
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/HaikeiLabs/kei-connector-contracts/tools/protoc-gen-kei/internal/manifest"
)

func main() {
	var opts manifest.Options
	protogen.Options{ParamFunc: opts.Set}.Run(func(gen *protogen.Plugin) error {
		return manifest.Generate(gen, opts)
	})
}
