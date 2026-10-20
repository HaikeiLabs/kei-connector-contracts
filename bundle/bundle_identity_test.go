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

package bundle

import (
	"encoding/json"
	"testing"
	"time"
)

func TestBundleMembershipsReturnsNilForModeOnline(t *testing.T) {
	v := &Verified{Payload: Payload{SubjectResolution: struct {
		Mode            string `json:"mode"`
		MaxCacheSeconds int    `json:"max_cache_seconds"`
	}{Mode: "online"}}}
	if m := BundleMemberships(v); m != nil {
		t.Fatal("expected nil memberships for mode=online bundle")
	}
}

func TestBundleMembershipsReturnsNilForEmptyIdentities(t *testing.T) {
	v := &Verified{Payload: Payload{SubjectResolution: struct {
		Mode            string `json:"mode"`
		MaxCacheSeconds int    `json:"max_cache_seconds"`
	}{Mode: "bundle"}}}
	if m := BundleMemberships(v); m != nil {
		t.Fatal("expected nil memberships when identities slice is empty")
	}
}

func TestBundleMembershipsLookup(t *testing.T) {
	now := time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)

	cases := []struct {
		name     string
		userID   string
		expected struct {
			groups  []string
			roles   []string
			granted bool
			known   bool
		}
	}{
		{
			name:   "linked identity grants groups and roles",
			userID: "u-alice",
			expected: struct {
				groups  []string
				roles   []string
				granted bool
				known   bool
			}{groups: []string{"members", "admins"}, roles: []string{"org:admin"}, granted: true, known: true},
		},
		{
			name:   "pending identity is known but not granted",
			userID: "u-bob",
			expected: struct {
				groups  []string
				roles   []string
				granted bool
				known   bool
			}{groups: []string{"members"}, roles: nil, granted: false, known: true},
		},
		{
			name:   "revoked identity is known but not granted",
			userID: "u-carol",
			expected: struct {
				groups  []string
				roles   []string
				granted bool
				known   bool
			}{groups: []string{"former"}, roles: nil, granted: false, known: true},
		},
		{
			name:   "cross-workspace identity is known but not granted",
			userID: "u-dave",
			expected: struct {
				groups  []string
				roles   []string
				granted bool
				known   bool
			}{groups: nil, roles: nil, granted: false, known: true},
		},
		{
			name:   "expired identity is not known",
			userID: "u-expired",
			expected: struct {
				groups  []string
				roles   []string
				granted bool
				known   bool
			}{groups: nil, roles: nil, granted: false, known: false},
		},
		{
			name:   "unknown userID is not known",
			userID: "u-unknown",
			expected: struct {
				groups  []string
				roles   []string
				granted bool
				known   bool
			}{groups: nil, roles: nil, granted: false, known: false},
		},
	}

	identities := []Identity{
		{
			UserID:     "u-alice",
			Status:     SubjectLinked,
			Groups:     []string{"members", "admins"},
			OrgRoles:   []string{"org:admin"},
			ExpireTime: time.Time{},
		},
		{
			UserID: "u-bob",
			Status: SubjectPending,
			Groups: []string{"members"},
		},
		{
			UserID: "u-carol",
			Status: SubjectRevoked,
			Groups: []string{"former"},
		},
		{
			UserID: "u-dave",
			Status: SubjectCrossWorkspace,
		},
		{
			UserID:     "u-expired",
			Status:     SubjectLinked,
			Groups:     []string{"members"},
			ExpireTime: now.Add(-1 * time.Hour),
		},
	}

	v := &Verified{Payload: Payload{
		SubjectResolution: struct {
			Mode            string `json:"mode"`
			MaxCacheSeconds int    `json:"max_cache_seconds"`
		}{Mode: "bundle"},
		Identities: identities,
	}}

	m := BundleMembershipsAt(v, now)
	if m == nil {
		t.Fatal("expected non-nil memberships for mode=bundle with identities")
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			groups, roles, granted, known := m.Lookup(tc.userID)
			if !stringSliceEqual(groups, tc.expected.groups) {
				t.Errorf("groups = %v, want %v", groups, tc.expected.groups)
			}
			if !stringSliceEqual(roles, tc.expected.roles) {
				t.Errorf("roles = %v, want %v", roles, tc.expected.roles)
			}
			if granted != tc.expected.granted {
				t.Errorf("granted = %v, want %v", granted, tc.expected.granted)
			}
			if known != tc.expected.known {
				t.Errorf("known = %v, want %v", known, tc.expected.known)
			}
		})
	}
}

func TestLookupIdentityUsesBundleKeysAndFailsClosed(t *testing.T) {
	now := time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)
	v := &Verified{Payload: Payload{
		SubjectResolution: struct {
			Mode            string `json:"mode"`
			MaxCacheSeconds int    `json:"max_cache_seconds"`
		}{Mode: "bundle", MaxCacheSeconds: 60},
		Identities: []Identity{
			{Source: "discord", ExternalID: "alice", UserID: "u-alice", Status: SubjectLinked},
			{Source: "discord", ExternalID: "pending", UserID: "u-pending", Status: SubjectPending},
			{Source: "discord", ExternalID: "expired", UserID: "u-expired", Status: SubjectLinked, ExpireTime: now},
		},
	}}
	if got, ok := v.LookupIdentity("discord", "alice", now); !ok || got.UserID != "u-alice" || got.Status != SubjectLinked {
		t.Fatalf("linked identity = %+v, %v", got, ok)
	}
	if got, ok := v.LookupIdentity("discord", "pending", now); !ok || got.Status != SubjectPending {
		t.Fatalf("pending identity = %+v, %v", got, ok)
	}
	for _, key := range [][2]string{{"discord", "missing"}, {"slack", "alice"}, {"discord", "expired"}, {"", "alice"}, {"discord", ""}} {
		if got, ok := v.LookupIdentity(key[0], key[1], now); ok {
			t.Errorf("LookupIdentity(%q,%q) = %+v, true; want fail closed", key[0], key[1], got)
		}
	}
	v.Payload.Identities = append(v.Payload.Identities, Identity{Source: "discord", ExternalID: "alice", UserID: "u-substitution", Status: SubjectLinked})
	if got, ok := v.LookupIdentity("discord", "alice", now); ok {
		t.Fatalf("duplicate identity mapping selected %+v", got)
	}
	v.Payload.SubjectResolution.Mode = "online"
	if got, ok := v.LookupIdentity("discord", "pending", now); ok {
		t.Fatalf("online mode returned identity %+v", got)
	}
}

func TestBundleMembershipsRegisteredWithPDPEvaluator(t *testing.T) {
	v := &Verified{Payload: Payload{
		SubjectResolution: struct {
			Mode            string `json:"mode"`
			MaxCacheSeconds int    `json:"max_cache_seconds"`
		}{Mode: "bundle"},
		Identities: []Identity{
			{UserID: "u-alice", Status: SubjectLinked, Groups: []string{"members"}},
		},
		PolicySet: mustRawMessage(t, `{"groups":[],"agents":[],"tools":[{"name":"web_search","service":"","action":"tool:web_search"}],"policies":[{"id":"p","src_pattern":"group:members","dst_pattern":"tool:web_search","action":"permit","scope":{"agent_id":null}}]}`),
	}}

	call := Call{Tool: "web_search", Action: "tool:web_search"}
	subject := &SessionSubject{UserID: "u-alice", Status: SubjectLinked}
	m := BundleMemberships(v)
	d := v.EvaluateSession(call, subject, m)
	if d.Outcome != "allow" {
		t.Fatalf("EvaluateSession outcome = %s, want allow (%+v)", d.Outcome, d)
	}
}

func TestBundleMembershipsLinkedViaEvaluateSessionWithBundleMode(t *testing.T) {
	v := &Verified{Payload: Payload{
		SubjectResolution: struct {
			Mode            string `json:"mode"`
			MaxCacheSeconds int    `json:"max_cache_seconds"`
		}{Mode: "bundle"},
		Identities: []Identity{
			{UserID: "u-alice", Status: SubjectLinked, Groups: []string{"members"}},
		},
		PolicySet: mustRawMessage(t, `{"groups":[],"agents":[{"id":"a1","is_default":true}],"tools":[{"name":"web_search","service":"","action":"tool:web_search"}],"policies":[{"id":"p","src_pattern":"group:members","dst_pattern":"tool:web_search","action":"permit","scope":{"agent_id":null}}]}`),
	}}

	t.Run("linked user gets allow", func(t *testing.T) {
		d := v.EvaluateSession(Call{Tool: "web_search", Action: "tool:web_search"}, &SessionSubject{UserID: "u-alice", Status: SubjectLinked}, BundleMemberships(v))
		if d.Outcome != "allow" {
			t.Fatalf("outcome = %s, want allow (%+v)", d.Outcome, d)
		}
	})

	t.Run("revoked user gets subject_revoked", func(t *testing.T) {
		v.Payload.Identities = []Identity{
			{UserID: "u-alice", Status: SubjectRevoked},
		}
		d := v.EvaluateSession(Call{Tool: "web_search", Action: "tool:web_search"}, &SessionSubject{UserID: "u-alice", Status: SubjectLinked}, BundleMemberships(v))
		if d.Outcome != "deny" || d.ReasonCode != "subject_revoked" {
			t.Fatalf("outcome = %s reason_code = %s, want deny subject_revoked (%+v)", d.Outcome, d.ReasonCode, d)
		}
	})

	t.Run("user not in bundle still gets subject_pending for group policy", func(t *testing.T) {
		v.Payload.Identities = []Identity{
			{UserID: "u-alice", Status: SubjectLinked, Groups: []string{"members"}},
		}
		d := v.EvaluateSession(Call{Tool: "web_search", Action: "tool:web_search"}, &SessionSubject{UserID: "u-other", Status: SubjectLinked}, BundleMemberships(v))
		if d.Outcome != "deny" || d.ReasonCode != "subject_pending" {
			t.Fatalf("outcome = %s reason_code = %s, want deny subject_pending (%+v)", d.Outcome, d.ReasonCode, d)
		}
	})
}

func TestBundleValidateAcceptsBundleMode(t *testing.T) {
	raw := `{
		"schema": "kei.policy-bundle/v1",
		"bundle_id": "bundle-001",
		"bundle_version": 42,
		"policy_revision": 1,
		"issuer": "kei-policy-catalog",
		"audience": {"installation_id": "inst-001", "org_id": "org-001", "workspace_id": "ws-001"},
		"issued_at": "2026-10-01T11:00:00Z",
		"not_after": "2026-10-01T23:00:00Z",
		"refresh": {"poll_interval_seconds": 300, "stale_after_seconds": 600},
		"subject_resolution": {"mode": "bundle"},
		"policy_set": {"match_semantics":"kei.match/v1","evaluation":{"order":"array_order_first_match","default":"deny"},"groups":[],"agents":[],"tools":[],"policies":[]},
		"agent_connector_bindings": []
	}`
	SetSchemaValidator(func([]byte) error { return nil })
	v, err := Validate([]byte(raw), "inst-001", "org-001", "ws-001", 41, time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if v.Payload.SubjectResolution.Mode != "bundle" {
		t.Fatalf("mode = %s, want bundle", v.Payload.SubjectResolution.Mode)
	}
}

func TestBundleValidateRejectsUnsupportedMode(t *testing.T) {
	raw := `{
		"schema": "kei.policy-bundle/v1",
		"bundle_id": "bundle-001",
		"bundle_version": 42,
		"policy_revision": 1,
		"issuer": "kei-policy-catalog",
		"audience": {"installation_id": "inst-001", "org_id": "org-001", "workspace_id": "ws-001"},
		"issued_at": "2026-10-01T11:00:00Z",
		"not_after": "2026-10-01T23:00:00Z",
		"refresh": {"poll_interval_seconds": 300, "stale_after_seconds": 600},
		"subject_resolution": {"mode": "hybrid"},
		"policy_set": {"match_semantics":"kei.match/v1","evaluation":{"order":"array_order_first_match","default":"deny"},"groups":[],"agents":[],"tools":[],"policies":[]},
		"agent_connector_bindings": []
	}`
	SetSchemaValidator(func([]byte) error { return nil })
	if _, err := Validate([]byte(raw), "inst-001", "org-001", "ws-001", 41, time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)); err != ErrBundleUnsupported {
		t.Fatalf("expected ErrBundleUnsupported, got %v", err)
	}
}

func TestBundleValidateAcceptsOnlineWithHighMaxCache(t *testing.T) {
	raw := `{
		"schema": "kei.policy-bundle/v1",
		"bundle_id": "bundle-001",
		"bundle_version": 42,
		"policy_revision": 1,
		"issuer": "kei-policy-catalog",
		"audience": {"installation_id": "inst-001", "org_id": "org-001", "workspace_id": "ws-001"},
		"issued_at": "2026-10-01T11:00:00Z",
		"not_after": "2026-10-01T23:00:00Z",
		"refresh": {"poll_interval_seconds": 300, "stale_after_seconds": 600},
		"subject_resolution": {"mode": "online", "max_cache_seconds": 120},
		"policy_set": {"match_semantics":"kei.match/v1","evaluation":{"order":"array_order_first_match","default":"deny"},"groups":[],"agents":[],"tools":[],"policies":[]},
		"agent_connector_bindings": []
	}`
	SetSchemaValidator(func([]byte) error { return nil })
	v, err := Validate([]byte(raw), "inst-001", "org-001", "ws-001", 41, time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if v.Payload.SubjectResolution.MaxCacheSeconds != 120 {
		t.Fatalf("unexpected max_cache_seconds: %d", v.Payload.SubjectResolution.MaxCacheSeconds)
	}
}

func TestBundleValidateAcceptsOnlineWithLowMaxCache(t *testing.T) {
	raw := `{
		"schema": "kei.policy-bundle/v1",
		"bundle_id": "bundle-001",
		"bundle_version": 42,
		"policy_revision": 1,
		"issuer": "kei-policy-catalog",
		"audience": {"installation_id": "inst-001", "org_id": "org-001", "workspace_id": "ws-001"},
		"issued_at": "2026-10-01T11:00:00Z",
		"not_after": "2026-10-01T23:00:00Z",
		"refresh": {"poll_interval_seconds": 300, "stale_after_seconds": 600},
		"subject_resolution": {"mode": "online", "max_cache_seconds": 60},
		"policy_set": {"match_semantics":"kei.match/v1","evaluation":{"order":"array_order_first_match","default":"deny"},"groups":[],"agents":[],"tools":[],"policies":[]},
		"agent_connector_bindings": []
	}`
	SetSchemaValidator(func([]byte) error { return nil })
	v, err := Validate([]byte(raw), "inst-001", "org-001", "ws-001", 41, time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if v.Payload.SubjectResolution.Mode != "online" || v.Payload.SubjectResolution.MaxCacheSeconds != 60 {
		t.Fatalf("unexpected payload: mode=%s max_cache=%d", v.Payload.SubjectResolution.Mode, v.Payload.SubjectResolution.MaxCacheSeconds)
	}
}

func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func mustRawMessage(t *testing.T, s string) json.RawMessage {
	t.Helper()
	return json.RawMessage(s)
}
