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

// Copyright 2026 Haikei Labs
// SPDX-License-Identifier: Apache-2.0

package policybundlehealth

import (
	"os"
	"strings"
	"testing"
	"time"
)

func ptr[T any](v T) *T { return &v }

func TestCanonicalFixtures(t *testing.T) {
	for _, tc := range []struct {
		path string
		read bool
	}{
		{"../schemas/examples/runtime-policy-bundle-health/heartbeat-active.v1.json", false},
		{"../schemas/examples/runtime-policy-bundle-health/read-active.v1.json", true},
	} {
		raw, err := os.ReadFile(tc.path)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := Decode(raw, tc.read); err != nil {
			t.Errorf("Decode(%s): %v", tc.path, err)
		}
	}
}

func TestHealthValidation(t *testing.T) {
	checked := time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)
	accepted := checked.Add(-time.Hour)
	expires := checked.Add(24 * time.Hour)
	version, revision := int64(42), int64(7)
	digest := "sha256:" + strings.Repeat("a", 64)
	tests := []struct {
		name  string
		h     Health
		valid bool
	}{
		{"cold before first check", Health{SchemaVersion: 1, State: StateCold}, true},
		{"cold checked", Health{SchemaVersion: 1, State: StateCold, CheckedAt: &checked, ReasonCode: ptr(ReasonBundleMissing)}, true},
		{"active", Health{SchemaVersion: 1, State: StateActive, BundleVersion: &version, PolicyRevision: &revision, BundleDigest: &digest, CheckedAt: &checked, AcceptedAt: &accepted, ExpiresAt: &expires}, true},
		{"stale", Health{SchemaVersion: 1, State: StateStale, BundleVersion: &version, PolicyRevision: &revision, BundleDigest: &digest, CheckedAt: &checked, AcceptedAt: &accepted, ExpiresAt: &expires, ReasonCode: ptr(ReasonBundleStale)}, true},
		{"unknown state", Health{SchemaVersion: 1, State: "healthy"}, false},
		{"partial accepted metadata", Health{SchemaVersion: 1, State: StateInvalid, CheckedAt: &checked, ReasonCode: ptr(ReasonRefreshFailed), BundleVersion: &version}, false},
		{"active with reason", Health{SchemaVersion: 1, State: StateActive, BundleVersion: &version, PolicyRevision: &revision, BundleDigest: &digest, CheckedAt: &checked, AcceptedAt: &accepted, ExpiresAt: &expires, ReasonCode: ptr(ReasonRefreshFailed)}, false},
		{"bad digest", Health{SchemaVersion: 1, State: StateActive, BundleVersion: &version, PolicyRevision: &revision, BundleDigest: ptr("secret"), CheckedAt: &checked, AcceptedAt: &accepted, ExpiresAt: &expires}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.h.Validate(false)
			if (err == nil) != tt.valid {
				t.Fatalf("Validate() error = %v, valid = %v", err, tt.valid)
			}
		})
	}
}
