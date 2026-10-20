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

package bundle_test

import (
	"strings"
	"testing"
	"time"

	"github.com/HaikeiLabs/kei-connector-contracts/bundle"
)

func TestValidate_ExternalCompileCheck(t *testing.T) {
	// This test verifies the public API compiles from an external package.
	// It sends a minimal v2 bundle; without a registered schema validator
	// ValidateV2 should fail with a clear error.
	input := []byte(`{"schema":"kei.policy-bundle/v2"}`)
	_, err := bundle.Validate(input, "inst", "org", "ws", 0, time.Now())
	if err == nil {
		t.Fatal("Validate on a bare v2 bundle without schema validator must fail")
	}
	if !strings.Contains(err.Error(), "schema validator") {
		t.Fatalf("expected error mentioning schema validator, got: %v", err)
	}
}
