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
	"sort"
	"strings"
	"time"

	"github.com/HaikeiLabs/kei-connector-contracts/harnessmatch"
)

// HarnessCall is one tool call of a pre-built harness, described by the
// inputs kei.harness-match/v1 matches on (ADR-029 §1.2, §1.3).
type HarnessCall = harnessmatch.Call

// Harness outcomes (ADR-029 §1.5). Every non-permit, non-deny outcome carries a
// reason code (HarnessReason*).
const (
	HarnessPermit        = string(harnessmatch.OutcomePermit)
	HarnessDeny          = string(harnessmatch.OutcomeDeny)
	HarnessUnmatched     = string(harnessmatch.OutcomeUnmatched)
	HarnessNotRenderable = string(harnessmatch.OutcomeNotRenderable)
)

// Harness reason codes (ADR-029 §1.5 extension). Every non-permit, non-deny
// outcome carries one; permit and deny have no reason code.
const (
	HarnessReasonUnsupportedKind  = harnessmatch.ReasonUnsupportedKind
	HarnessReasonPathNotAbsolute  = harnessmatch.ReasonPathNotAbsolute
	HarnessReasonBundleInvalid    = harnessmatch.ReasonBundleInvalid
	HarnessReasonHumanSource      = harnessmatch.ReasonHumanSource
	HarnessReasonDstNotRenderable = harnessmatch.ReasonDstNotRenderable
	HarnessReasonEvaluateViaPDP   = harnessmatch.ReasonEvaluateViaPDP
)

// EvaluateHarness applies kei.harness-match/v1 to a harness tool call. It
// decodes the bundle's policy set and delegates to harnessmatch.Evaluate.
// EvaluateV2Shell applies ordered v2 harness policies to a shell command
// using the authenticated session subject and the verified membership snapshot.
func (v *Verified) EvaluateV2Shell(argv []string, agentID string, subject *SessionSubject, now time.Time) Decision {
	return v.EvaluateV2ShellForKind("pi", argv, agentID, subject, now)
}

// EvaluateV2ShellForKind is EvaluateV2Shell for any harness kind whose shell
// the local PDP decides: harness:<kind> sources match the given kind.
func (v *Verified) EvaluateV2ShellForKind(kind string, argv []string, agentID string, subject *SessionSubject, now time.Time) Decision {
	if agentID == "" {
		return Decision{Outcome: "deny", ReasonCode: "agent_not_assigned", Reason: "the trusted default agent is unavailable"}
	}
	if decision, denied := v.v2SubjectDenial(subject, now); denied {
		return decision
	}
	members := BundleMembershipsAt(v, now)
	if members == nil {
		return Decision{Outcome: "deny", ReasonCode: "subject_pending", Reason: "the verified membership snapshot is unavailable"}
	}
	groups, roles, granted, known := members.Lookup(subject.UserID)
	if !known || !granted {
		return Decision{Outcome: "deny", ReasonCode: "subject_pending", Reason: "the verified membership snapshot does not grant the session subject"}
	}
	var set v2PolicySet
	if err := strictJSON(v.Payload.PolicySet, &set); err != nil || set.HarnessMatchSemantics != "kei.harness-match/v1" {
		return invalidDecision()
	}
	call := harnessmatch.Call{Kind: kind, AgentID: agentID, Argv: argv}
	matching := make([]v2HarnessPolicy, 0, len(set.HarnessPolicies))
	for _, policy := range set.HarnessPolicies {
		if !v2ScopeMatches(policy.Scope, agentID, v.Payload.Audience.OrgID, v.Payload.Audience.WorkspaceID) {
			continue
		}
		matched, renderable := harnessmatch.MatchDst(policy.Dst, call)
		if !renderable || !matched {
			continue
		}
		// Shell calls can only be permitted by a literal shell: argv prefix.
		// A registered tool grant (including tool:pi.bash) is never a shell grant.
		if policy.Action == policyPermitAction && (!strings.HasPrefix(policy.Dst, "shell:") || policy.Dst == "shell:*") {
			continue
		}
		if sourceMatches(policy.Src, call, subject, groups, roles) {
			matching = append(matching, policy)
		}
	}
	if len(matching) == 0 {
		return Decision{Outcome: "deny", ReasonCode: "no_matching_policy", Reason: "no shell policy matched"}
	}
	sort.Slice(matching, func(i, j int) bool {
		if !set.denyOverrides() && matching[i].Priority != matching[j].Priority {
			return matching[i].Priority > matching[j].Priority
		}
		if matching[i].Action != matching[j].Action {
			return matching[i].Action == "deny"
		}
		return matching[i].PolicyID < matching[j].PolicyID
	})
	winner := matching[0]
	if winner.Action == "deny" {
		return Decision{Outcome: "deny", PolicyID: winner.PolicyID, ReasonCode: "policy_deny", Reason: "policy denied"}
	}
	return Decision{Outcome: "allow", PolicyID: winner.PolicyID, ReasonCode: "policy_permit"}
}

func sourceMatches(pattern string, call harnessmatch.Call, subject *SessionSubject, groups, roles []string) bool {
	if pattern == "harness:*" {
		return true
	}
	if strings.HasPrefix(pattern, "harness:") {
		return strings.TrimPrefix(pattern, "harness:") == call.Kind
	}
	// Non-harness patterns route through srcMatch, which delegates agent: and
	// org: matching to its own logic and routes everything else to
	// pol.MatchPrincipal.
	return srcMatch(pattern, Call{AgentID: call.AgentID, UserID: subject.UserID, Email: subject.Email}, groups, roles)
}

// EvaluateV2MCP applies ordered v2 harness policies to an MCP tool call
// using the authenticated session subject, the verified membership snapshot,
// and the calling harness kind (ADR-036 §4.6).
func (v *Verified) EvaluateV2MCP(slug, tool, agentID, harnessKind string, subject *SessionSubject, now time.Time) Decision {
	if agentID == "" {
		return Decision{Outcome: "deny", ReasonCode: "agent_not_assigned", Reason: "the trusted default agent is unavailable"}
	}
	if decision, denied := v.v2SubjectDenial(subject, now); denied {
		return decision
	}
	members := BundleMembershipsAt(v, now)
	if members == nil {
		return Decision{Outcome: "deny", ReasonCode: "subject_pending", Reason: "the verified membership snapshot is unavailable"}
	}
	groups, roles, granted, known := members.Lookup(subject.UserID)
	if !known || !granted {
		return Decision{Outcome: "deny", ReasonCode: "subject_pending", Reason: "the verified membership snapshot does not grant the session subject"}
	}
	var set v2PolicySet
	if err := strictJSON(v.Payload.PolicySet, &set); err != nil || set.HarnessMatchSemantics != "kei.harness-match/v1" {
		return invalidDecision()
	}
	call := harnessmatch.Call{Kind: harnessKind, AgentID: agentID, MCPServer: slug, MCPTool: tool}
	matching := make([]v2HarnessPolicy, 0, len(set.HarnessPolicies))
	for _, policy := range set.HarnessPolicies {
		if !v2ScopeMatches(policy.Scope, agentID, v.Payload.Audience.OrgID, v.Payload.Audience.WorkspaceID) {
			continue
		}
		matched, renderable := harnessmatch.MatchDst(policy.Dst, call)
		if !renderable || !matched {
			continue
		}
		if sourceMatches(policy.Src, call, subject, groups, roles) {
			matching = append(matching, policy)
		}
	}
	if len(matching) == 0 {
		return Decision{Outcome: "deny", ReasonCode: "no_matching_policy", Reason: "no MCP policy matched"}
	}
	sort.Slice(matching, func(i, j int) bool {
		if !set.denyOverrides() && matching[i].Priority != matching[j].Priority {
			return matching[i].Priority > matching[j].Priority
		}
		if matching[i].Action != matching[j].Action {
			return matching[i].Action == "deny"
		}
		return matching[i].PolicyID < matching[j].PolicyID
	})
	winner := matching[0]
	if winner.Action == "deny" {
		return Decision{Outcome: "deny", PolicyID: winner.PolicyID, ReasonCode: "policy_deny", Reason: "policy denied"}
	}
	return Decision{Outcome: "allow", PolicyID: winner.PolicyID, ReasonCode: "policy_permit"}
}

func (v *Verified) EvaluateHarness(c HarnessCall) (outcome, policyID, reason string) {
	if v != nil && v.Payload.Schema == "kei.policy-bundle/v2" && !v.SubjectSnapshotFresh(time.Now()) {
		return HarnessNotRenderable, "", HarnessReasonBundleInvalid
	}
	var policies []harnessmatch.Policy
	denyFirst := false
	if v.Payload.Schema == "kei.policy-bundle/v2" {
		// v2 carries harness policies under policy_set (kei#852).
		var set v2PolicySet
		if err := json.Unmarshal(v.Payload.PolicySet, &set); err != nil || set.HarnessMatchSemantics != "kei.harness-match/v1" {
			return HarnessNotRenderable, "", HarnessReasonBundleInvalid
		}
		denyFirst = set.denyOverrides()
		policies = make([]harnessmatch.Policy, len(set.HarnessPolicies))
		for i, p := range set.HarnessPolicies {
			policies[i] = harnessmatch.Policy{ID: p.PolicyID, Src: p.Src, Dst: p.Dst, Action: p.Action, Enabled: true, Scope: p.Scope.AgentID}
		}
	} else {
		set, err := v.decodePolicySet()
		if err != nil {
			return HarnessNotRenderable, "", HarnessReasonBundleInvalid
		}
		denyFirst = set.Evaluation.Order == evalDenyOverrides
		policies = make([]harnessmatch.Policy, len(set.Policies))
		for i, p := range set.Policies {
			policies[i] = harnessmatch.Policy{ID: p.ID, Src: p.Src, Dst: p.Dst, Action: p.Action, Enabled: true, Scope: p.Scope.AgentID}
		}
	}
	if denyFirst {
		policies = denyBeforePermit(policies)
	}
	result := harnessmatch.EvaluateWith(c, policies, harnessmatch.KindSet{Known: v.Payload.KnownHarnessKind, ViaPDP: v.Payload.HarnessKindDecidesViaPDP})
	return string(result.Outcome), result.PolicyID, result.Reason
}
