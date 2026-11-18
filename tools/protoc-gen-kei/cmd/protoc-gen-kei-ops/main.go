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

// Command protoc-gen-kei-ops generates typed connector operations from the
// kei.contracts.v1 connector and operation options (proposal 0002).
//
// Parameters: lang=go|python|ts|rust (default go) and
// signature=message|positional (default message). The plugin validates every
// operation first and fails the build on any rule violation.
package main

import (
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/HaikeiLabs/kei-connector-contracts/tools/protoc-gen-kei/internal/render"
)

func main() {
	var opts render.Options
	protogen.Options{ParamFunc: opts.Set}.Run(func(gen *protogen.Plugin) error {
		return render.Generate(gen, opts)
	})
}
