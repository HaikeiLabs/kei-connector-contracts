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

// Package canonicaltest holds the shared golden cases for the canonical
// package, so every service that computes or compares a canonical value
// tests against the same table. The cases are also published as
// resource.json for implementations in other languages.
package canonicaltest

import (
	_ "embed"
	"encoding/json"
)

//go:embed resource.json
var resourceJSON []byte

// ResourceCase is one golden case for canonical.Resource. Error cases have
// no Want. Note explains a case that pins surprising behaviour.
type ResourceCase struct {
	Name  string `json:"name"`
	Input string `json:"input"`
	Want  string `json:"want,omitempty"`
	Error bool   `json:"error,omitempty"`
	Note  string `json:"note,omitempty"`
}

var resourceCases = func() []ResourceCase {
	var golden struct {
		Cases []ResourceCase `json:"cases"`
	}
	if err := json.Unmarshal(resourceJSON, &golden); err != nil {
		panic("canonicaltest: resource.json: " + err.Error())
	}
	return golden.Cases
}()

// ResourceCases returns a copy of the golden cases for canonical.Resource,
// in the order resource.json lists them.
func ResourceCases() []ResourceCase {
	return append([]ResourceCase(nil), resourceCases...)
}
