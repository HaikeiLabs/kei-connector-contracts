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

package audit

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/HaikeiLabs/kei-connector-contracts/contract"
	"github.com/HaikeiLabs/kei-connector-contracts/governance"
)

func TestCredentialKindForSource(t *testing.T) {
	cases := []struct {
		source contract.CredentialSource
		want   CredentialKind
		ok     bool
	}{
		{source: contract.CredentialSourceOpaqueRef, want: CredentialKindSecretStore, ok: true},
		{source: "", want: CredentialKindSecretStore, ok: true},
		{source: contract.CredentialSourceOAuth, want: CredentialKindOAuth, ok: true},
		{source: "kerberos", ok: false},
	}
	for _, tc := range cases {
		got, ok := CredentialKindFor(tc.source)
		if got != tc.want || ok != tc.ok {
			t.Errorf("CredentialKindFor(%q) = %q, %v; want %q, %v", tc.source, got, ok, tc.want, tc.ok)
		}
	}
}

func TestCredentialKindValid(t *testing.T) {
	for kind, want := range map[CredentialKind]bool{
		CredentialKindSecretStore: true,
		CredentialKindOAuth:       true,
		"":                        false,
		"opaque_ref":              false,
		"OAuth":                   false,
	} {
		if got := kind.Valid(); got != want {
			t.Errorf("CredentialKind(%q).Valid() = %v, want %v", kind, got, want)
		}
	}
}

func TestAuditRecordCredentialKindWireShape(t *testing.T) {
	cases := []struct {
		name string
		kind CredentialKind
		want string
	}{
		{name: "secret store", kind: CredentialKindSecretStore, want: `"secret_store"`},
		{name: "oauth", kind: CredentialKindOAuth, want: `"oauth"`},
		{name: "no credential used", kind: "", want: ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := BuildAuditRecord(decisionEnvelope(), governance.PolicyDecision{Decision: governance.DecisionAllow}, time.Now())
			rec.CredentialKind = tc.kind
			data, err := json.Marshal(rec)
			if err != nil {
				t.Fatal(err)
			}
			var fields map[string]json.RawMessage
			if err := json.Unmarshal(data, &fields); err != nil {
				t.Fatal(err)
			}
			if got := string(fields["credential_kind"]); got != tc.want {
				t.Fatalf("credential_kind = %q, want %q (%s)", got, tc.want, data)
			}
			if _, ok := fields["credential_released"]; ok {
				t.Fatalf("record carries credential_released; release status is not part of the audit record (%s)", data)
			}
			var back AuditRecord
			if err := json.Unmarshal(data, &back); err != nil {
				t.Fatal(err)
			}
			if back.CredentialKind != tc.kind {
				t.Fatalf("round trip credential_kind = %q, want %q", back.CredentialKind, tc.kind)
			}
		})
	}
}

func TestBuildAuditRecordLeavesCredentialKindUnset(t *testing.T) {
	rec := BuildAuditRecord(decisionEnvelope(), governance.PolicyDecision{Decision: governance.DecisionAllow}, time.Now())
	if rec.CredentialKind != "" {
		t.Fatalf("credential_kind = %q; only a call that used a credential sets it", rec.CredentialKind)
	}
}
