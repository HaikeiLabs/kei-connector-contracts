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

package schemas

import "embed"

// FS is the embedded filesystem containing the shipped call templates
// (*.tool-definition.v1.json) and the JSON Schemas they validate against
// (tool-definition.v1.schema.json, tool-call-template.v1.schema.json), and the
// policy target and principal schemas (policy-target.v1.schema.json,
// policy-principal.v1.schema.json).
//
// Use it from a harness to load templates without a filesystem path:
//
//	data, _ := schemas.FS.ReadFile("github.tool-definition.v1.json")
//	templates, err := calltemplate.Parse(data)
//
//go:embed *.tool-definition.v1.json tool-definition.v1.schema.json tool-call-template.v1.schema.json policy-target.v1.schema.json policy-principal.v1.schema.json
var FS embed.FS
