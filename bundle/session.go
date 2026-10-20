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
	"strings"
	"time"
)

// DecisionPermit is the ADR-028 permit decision word, for callers that report
// a decision in that vocabulary.
const DecisionPermit = policyPermitAction

// Session subject statuses (ADR-026 G0, kei#681). The subject is established
// for a worker session by machine login plus kei or harness login and is held
// only in the enforcing service's memory; it is never read from a per-call
// request.
const (
	SubjectLinked         = "linked"
	SubjectUnlinked       = "unlinked"
	SubjectPending        = "pending"
	SubjectRevoked        = "revoked"
	SubjectCrossWorkspace = "cross_workspace"
)

// SessionSubject is the acting user of the current worker session.
type SessionSubject struct {
	UserID     string
	Email      string
	Source     string
	ExternalID string
	Groups     []string
	Status     string
}

// Memberships answers which groups, org roles and workspace grant the bundle
// holds for a user. G0 puts memberships and grants in the bundle, but the
// kei.policy-bundle/v1 schema has no field for them yet, so a v1 bundle
// yields no Memberships and every membership-dependent source fails closed.
type Memberships interface {
	Lookup(userID string) (groups, roles []string, granted, known bool)
}

// BundleMemberships returns the memberships carried by the bundle, or nil
// when the bundle does not carry them (mode=online or no identities).
func BundleMemberships(v *Verified) Memberships {
	return BundleMembershipsAt(v, time.Now())
}

// BundleMembershipsAt binds identity expiry checks to the evaluator's clock so
// snapshot freshness and expiry can be tested deterministically.
func BundleMembershipsAt(v *Verified, now time.Time) Memberships {
	if v == nil || v.Payload.SubjectResolution.Mode != "bundle" || len(v.Payload.Identities) == 0 || (v.Payload.Schema == "kei.policy-bundle/v2" && !v.SubjectSnapshotFresh(now)) {
		return nil
	}
	return &bundleMemberships{identities: v.Payload.Identities, now: now, aggregate: v.Payload.Schema == "kei.policy-bundle/v2"}
}

type bundleMemberships struct {
	identities []Identity
	now        time.Time
	aggregate  bool
}

func (m *bundleMemberships) Lookup(userID string) (groups, roles []string, granted, known bool) {
	if m.aggregate {
		groupSet, roleSet := map[string]bool{}, map[string]bool{}
		for _, id := range m.identities {
			if id.UserID != userID || (!id.ExpireTime.IsZero() && !m.now.Before(id.ExpireTime)) {
				continue
			}
			known = true
			granted = granted || id.Status == SubjectLinked
			for _, group := range id.Groups {
				groupSet[group] = true
			}
			for _, role := range id.OrgRoles {
				roleSet[role] = true
			}
		}
		for group := range groupSet {
			groups = append(groups, group)
		}
		for role := range roleSet {
			roles = append(roles, role)
		}
		return groups, roles, granted, known
	}
	for _, id := range m.identities {
		if id.UserID != userID {
			continue
		}
		if !id.ExpireTime.IsZero() && !m.now.Before(id.ExpireTime) {
			return nil, nil, false, false
		}
		known = true
		granted = id.Status == SubjectLinked
		return id.Groups, id.OrgRoles, granted, known
	}
	return nil, nil, false, false
}

// LookupLinkedUser resolves a CLI session's user_id against bundle identities.
// Remote sessions must use the exact source/external_id key instead.
func (v *Verified) LookupLinkedUser(userID string, now time.Time) (Identity, bool) {
	if v == nil || v.Payload.Schema != "kei.policy-bundle/v2" || v.Payload.SubjectResolution.Mode != "bundle" || !v.SubjectSnapshotFresh(now) || userID == "" {
		return Identity{}, false
	}
	for _, id := range v.Payload.Identities {
		if id.UserID == userID && id.Status == SubjectLinked && (id.ExpireTime.IsZero() || now.Before(id.ExpireTime)) {
			return id, true
		}
	}
	return Identity{}, false
}

func statusReasonCode(status string) string {
	switch status {
	case SubjectPending:
		return "subject_pending"
	case SubjectRevoked:
		return "subject_revoked"
	case SubjectCrossWorkspace:
		return "subject_cross_workspace"
	default:
		return "subject_unlinked"
	}
}

// SubjectSnapshotFresh reports whether the bundle subject snapshot is usable.
// HAI-452: the snapshot is trusted for the session lifetime — no
// max_cache_seconds window is enforced. Returns true when the bundle
// itself has been issued (IssuedAt ≤ now, allowing MaxIssuedAtClockSkew so a
// bundle Validate accepted is usable at once).
func (v *Verified) SubjectSnapshotFresh(now time.Time) bool {
	if v == nil || v.Payload.Schema != "kei.policy-bundle/v2" || v.Payload.SubjectResolution.Mode != "bundle" || v.Payload.Identities == nil {
		return false
	}
	return !v.Payload.IssuedAt.After(now.Add(MaxIssuedAtClockSkew))
}

// LookupIdentity resolves an external session key against the verified
// bundle's ADR-034 identities table. Duplicate matches and expired entries are
// rejected rather than allowing array order to select an identity.
func (v *Verified) LookupIdentity(source, externalID string, now time.Time) (Identity, bool) {
	if v == nil || v.Payload.SubjectResolution.Mode != "bundle" || source == "" || externalID == "" || (v.Payload.Schema == "kei.policy-bundle/v2" && !v.SubjectSnapshotFresh(now)) {
		return Identity{}, false
	}
	var found Identity
	matched := false
	for _, id := range v.Payload.Identities {
		if id.Source != source || id.ExternalID != externalID {
			continue
		}
		if matched || (!id.ExpireTime.IsZero() && !now.Before(id.ExpireTime)) {
			return Identity{}, false
		}
		found, matched = id, true
	}
	return found, matched
}

// EvaluateSession decides a call with the in-memory session subject and no
// catalog round trip. Registry and agent checks come first, then the subject
// gate (a non-linked session denies before any policy), then the
// catalog-ordered policy walk. A human-source policy whose dst matches never
// falls through: without a session it denies subject_unlinked, and without
// bundle memberships a group or org source denies subject_pending.
func (v *Verified) EvaluateSession(call Call, subject *SessionSubject, members Memberships) Decision {
	return v.EvaluateSessionAt(call, subject, members, time.Now())
}

// EvaluateSessionAt adds the v2 subject snapshot freshness gate while keeping
// v1 decision semantics unchanged.
func (v *Verified) EvaluateSessionAt(call Call, subject *SessionSubject, members Memberships, now time.Time) Decision {
	if v.Payload.Schema == "kei.policy-bundle/v2" {
		if !v.SubjectSnapshotFresh(now) {
			return Decision{Outcome: "deny", ReasonCode: "subject_snapshot_stale", Reason: "the local subject snapshot is missing or stale"}
		}
		return Decision{Outcome: "deny", ReasonCode: "preflight_required", Reason: "v2 connector operations require complete-plan preflight"}
	}
	call.UserID, call.Email = "", ""
	var pre Decision
	if call.ConnectorInvocation {
		pre = v.EvaluateConnector(call, nil)
	} else {
		pre = v.Evaluate(call, nil)
	}
	switch pre.ReasonCode {
	case "policy_bundle_invalid", "tool_not_registered", "tool_action_mismatch", "agent_not_assigned", "agent_not_bound_to_connector":
		return pre
	}
	set, err := v.decodePolicySet()
	if err != nil {
		return invalidDecision()
	}
	if call.AgentID == "" {
		call.AgentID = defaultAgent(set)
	}
	var reg tool
	for _, t := range set.Tools {
		if (call.ConnectorInvocation && t.ConnectorID == call.ConnectorID && t.Capability == call.Capability) ||
			(!call.ConnectorInvocation && t.Name == call.Tool) {
			reg = t
			break
		}
	}
	if call.ConnectorInvocation {
		call.Tool = reg.Name
	} else {
		if call.ConnectorID == "" {
			call.ConnectorID = reg.ConnectorID
		}
		if call.Capability == "" {
			call.Capability = reg.Capability
		}
	}
	if subject != nil {
		if d, denied := subjectGate(subject); denied {
			return d
		}
		if v.Payload.Schema == "kei.policy-bundle/v2" {
			var id Identity
			var found bool
			if subject.Source == "cli" && subject.ExternalID == "" {
				id, found = v.LookupLinkedUser(subject.UserID, now)
			} else {
				id, found = v.LookupIdentity(subject.Source, subject.ExternalID, now)
			}
			if !found || id.UserID != subject.UserID {
				return Decision{Outcome: "deny", ReasonCode: "subject_unlinked", Reason: "the session identity is not present in the fresh bundle snapshot"}
			}
			if id.Status != SubjectLinked {
				return Decision{Outcome: "deny", ReasonCode: statusReasonCode(id.Status), Reason: "the bundle snapshot does not grant the session identity"}
			}
		}
		call.UserID, call.Email = subject.UserID, subject.Email
	}
	var groups, roles []string
	known := false
	if subject != nil && members != nil {
		var granted bool
		groups, roles, granted, known = members.Lookup(subject.UserID)
		if known && !granted {
			return Decision{Outcome: "deny", ReasonCode: "subject_revoked", Reason: "the bundle no longer grants the session subject"}
		}
	}
	for _, p := range set.Policies {
		if p.Scope.AgentID != nil && *p.Scope.AgentID != call.AgentID {
			continue
		}
		if !dstMatch(p.Dst, call, reg) {
			continue
		}
		if requiresSubject(p.Src) && !srcMatch(p.Src, Call{AgentID: call.AgentID}, nil, nil) {
			if subject == nil {
				return Decision{Outcome: "deny", ReasonCode: "subject_unlinked", Reason: "a matching policy names a person and there is no session subject"}
			}
			if needsMemberships(p.Src) && !known {
				return Decision{Outcome: "deny", ReasonCode: "subject_pending", Reason: "the bundle carries no memberships or grants for the session subject"}
			}
		}
		if !srcMatch(p.Src, call, groups, roles) {
			continue
		}
		if p.Action == "deny" {
			return Decision{Outcome: "deny", PolicyID: p.ID, ReasonCode: "policy_deny", Reason: "policy denied"}
		}
		if p.Action == policyPermitAction {
			return Decision{Outcome: "allow", PolicyID: p.ID}
		}
	}
	return Decision{Outcome: "deny", ReasonCode: "no_matching_policy", Reason: "no policy matched"}
}

func subjectGate(s *SessionSubject) (Decision, bool) {
	switch s.Status {
	case SubjectLinked:
		if s.UserID != "" {
			return Decision{}, false
		}
		return Decision{Outcome: "deny", ReasonCode: "subject_unlinked", Reason: "the session subject has no user id"}, true
	case SubjectPending:
		return Decision{Outcome: "deny", ReasonCode: "subject_pending", Reason: "the session subject's workspace access is pending"}, true
	case SubjectRevoked:
		return Decision{Outcome: "deny", ReasonCode: "subject_revoked", Reason: "the session subject was revoked"}, true
	case SubjectCrossWorkspace:
		return Decision{Outcome: "deny", ReasonCode: "subject_cross_workspace", Reason: "the session subject belongs to another workspace"}, true
	default:
		return Decision{Outcome: "deny", ReasonCode: "subject_unlinked", Reason: "the session subject is not linked"}, true
	}
}

// needsMemberships reports whether a source pattern can only be decided from
// group memberships or org roles (srcMatch also treats a bare name as a group).
func needsMemberships(pattern string) bool {
	if strings.HasPrefix(pattern, "group:") || strings.HasPrefix(pattern, "org:") {
		return true
	}
	return !strings.Contains(pattern, ":") && !strings.HasPrefix(pattern, "*")
}
