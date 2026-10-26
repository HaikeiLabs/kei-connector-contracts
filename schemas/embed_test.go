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

import (
	"io/fs"
	"testing"

	"github.com/HaikeiLabs/kei-connector-contracts/calltemplate"
)

func TestEmbeddedToolDefinitionsAreValidCallTemplates(t *testing.T) {
	// Find all shipped template files in the embedded FS.
	matches, err := fs.Glob(FS, "*.tool-definition.v1.json")
	if err != nil {
		t.Fatalf("fs.Glob: %v", err)
	}
	if len(matches) == 0 {
		t.Fatal("embedded FS contains no tool-definition files")
	}

	for _, name := range matches {
		t.Run(name, func(t *testing.T) {
			data, err := FS.ReadFile(name)
			if err != nil {
				t.Fatalf("ReadFile(%s): %v", name, err)
			}
			if _, err := calltemplate.Parse(data); err != nil {
				t.Fatalf("calltemplate.Parse(%s): %v", name, err)
			}
		})
	}
}
