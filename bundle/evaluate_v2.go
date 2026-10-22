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
	"time"
)

// Operation is one exact planner-produced connector operation. Resource
// identity remains structured; no field is parsed from or flattened into another.
type Operation struct {
	ID                 string
	ConnectorID        string
	Capability         string
	Action             string
	Resource           string
	InputDigest        string
	ProviderInput      json.RawMessage
	ResourceType       string
	ResourceID         string
	ResourceParent     string
	ResourceParentType string
}

type V2ResourceType struct {
	Type       string `json:"type"`
	ParentType string `json:"parent_type,omitempty"`
}

type V2Tool struct {
	Name           string `json:"name"`
	Source         string `json:"source"`
	OperationClass string `json:"operation_class"`
	Route          struct {
		ConnectorBinding *struct {
			ConnectorID string `json:"connector_id"`
		} `json:"connector_binding,omitempty"`
		HarnessExecutor *json.RawMessage `json:"harness_executor,omitempty"`
	} `json:"route"`
	RequiredCapabilities []string         `json:"required_capabilities,omitempty"`
	ResourceTypes        []V2ResourceType `json:"resource_types,omitempty"`
	Plan                 json.RawMessage  `json:"plan,omitempty"`
}

// V2HarnessTool reports whether name resolves to an explicitly registered
// harness_executor route in a verified v2 bundle.
func (v *Verified) V2HarnessTool(name string, now time.Time) bool {
	if v == nil || v.Payload.Schema != "kei.policy-bundle/v2" || !v.SubjectSnapshotFresh(now) {
		return false
	}
	var set v2PolicySet
	if strictJSON(v.Payload.PolicySet, &set) != nil {
		return false
	}
	for _, tool := range set.Tools {
		if tool.Name == name {
			return tool.Route.HarnessExecutor != nil && tool.Route.ConnectorBinding == nil
		}
	}
	return false
}

// V2ConnectorTool resolves one trusted connector tool from the parsed v2 set.
func (v *Verified) V2ConnectorTool(name string, now time.Time) (V2Tool, bool) {
	if v == nil || v.Payload.Schema != "kei.policy-bundle/v2" || !v.SubjectSnapshotFresh(now) {
		return V2Tool{}, false
	}
	var set v2PolicySet
	if strictJSON(v.Payload.PolicySet, &set) != nil {
		return V2Tool{}, false
	}
	for _, tool := range set.Tools {
		if tool.Name == name && tool.Route.ConnectorBinding != nil && len(tool.Plan) > 0 && planValid(tool) == nil {
			return tool, true
		}
	}
	return V2Tool{}, false
}

type v2Check struct {
	PolicyID string `json:"policy_id"`
	Name     string `json:"name"`
	Effect   string `json:"effect"`
	Priority int    `json:"priority"`
	Scope    struct {
		OrgID       *string `json:"org_id"`
		WorkspaceID *string `json:"workspace_id"`
		AgentID     *string `json:"agent_id"`
	} `json:"scope"`
	Selector map[string]string `json:"selector"`
}

type v2PolicySet struct {
	Schema         string `json:"schema"`
	OrgID          string `json:"org_id"`
	WorkspaceID    string `json:"workspace_id"`
	MatchSemantics string `json:"match_semantics"`
	Evaluation     struct {
		Precedence []string `json:"precedence"`
		Checks     string   `json:"checks"`
		Default    string   `json:"default"`
	} `json:"evaluation"`
	Tools                 []V2Tool          `json:"tools"`
	Checks                []v2Check         `json:"checks"`
	ToolPolicies          []v2Check         `json:"tool_policies"`
	HarnessMatchSemantics string            `json:"harness_match_semantics"`
	HarnessPolicies       []v2HarnessPolicy `json:"harness_policies"`
	// Connectors is the optional descriptor section (HAI-378); see Connector.
	Connectors []Connector `json:"connectors,omitempty"`
}

// v2HarnessPolicy is one ADR-029 harness policy row as the catalog emits it
// under policy_set.harness_policies (kei.harness-match/v1).
type v2HarnessPolicy struct {
	PolicyID string `json:"policy_id"`
	Name     string `json:"name"`
	Src      string `json:"src_pattern"`
	Dst      string `json:"dst_pattern"`
	Action   string `json:"action"`
	Priority int    `json:"priority"`
	Scope    struct {
		OrgID       *string `json:"org_id"`
		WorkspaceID *string `json:"workspace_id"`
		AgentID     *string `json:"agent_id"`
	} `json:"scope"`
}

// v2SubjectDenial is the v2 session-subject gate shared by every local v2
// decision: a fresh identity snapshot must corroborate a linked session
// subject. It reports the denial, or false when the subject passes.
func (v *Verified) v2SubjectDenial(subject *SessionSubject, now time.Time) (Decision, bool) {
	if v == nil || v.Payload.Schema != "kei.policy-bundle/v2" || !v.SubjectSnapshotFresh(now) {
		return Decision{Outcome: "deny", ReasonCode: "subject_snapshot_stale", Reason: "the local subject snapshot is missing or stale"}, true
	}
	if subject == nil {
		return Decision{Outcome: "deny", ReasonCode: "subject_unlinked", Reason: "no session subject"}, true
	}
	if decision, denied := subjectGate(subject); denied {
		return decision, true
	}
	var identity Identity
	var known bool
	if subject.Source == "installation_owner" {
		if v.Payload.Subject != nil && v.Payload.Subject.UserID == subject.UserID {
			identity, known = Identity{UserID: subject.UserID, Status: SubjectLinked}, true
		}
	} else if subject.Source == "cli" && subject.ExternalID == "" {
		identity, known = v.LookupLinkedUser(subject.UserID, now)
	} else {
		identity, known = v.LookupIdentity(subject.Source, subject.ExternalID, now)
	}
	if !known || identity.UserID != subject.UserID {
		return Decision{Outcome: "deny", ReasonCode: "subject_unlinked", Reason: "the session identity is not present in the fresh bundle snapshot"}, true
	}
	if identity.Status != SubjectLinked {
		return Decision{Outcome: "deny", ReasonCode: statusReasonCode(identity.Status), Reason: "the bundle snapshot does not grant the session identity"}, true
	}
	return Decision{}, false
}

// EvaluateV2HarnessTool decides one harness-executed tool against the v2
// tool_policies set. The harness executes the tool; this method only returns
// the local policy decision. A fresh bundle identity must corroborate the
// named session subject before any tool policy can permit it.
func (v *Verified) EvaluateV2HarnessTool(toolName, agentID string, subject *SessionSubject, now time.Time) Decision {
	if decision, denied := v.v2SubjectDenial(subject, now); denied {
		return decision
	}
	var set v2PolicySet
	if err := strictJSON(v.Payload.PolicySet, &set); err != nil {
		return invalidDecision()
	}
	var registered *V2Tool
	for i := range set.Tools {
		if set.Tools[i].Name == toolName {
			registered = &set.Tools[i]
			break
		}
	}
	if registered == nil || registered.Route.HarnessExecutor == nil || registered.Route.ConnectorBinding != nil {
		return Decision{Outcome: "deny", ReasonCode: "tool_not_registered", Reason: "harness tool is not registered in the policy bundle"}
	}
	request := map[string]string{"tool": registered.Name, "source": registered.Source}
	groups, roles := membershipValues(v.Payload.Identities, subject.UserID, now)
	groups = append(groups, subject.Groups...)
	return evaluateV2Check(set.ToolPolicies, request, subject, groups, roles, agentID, v.Payload.Audience.OrgID, v.Payload.Audience.WorkspaceID)
}

// EvaluateV2Plan locally AND-authorizes a complete exact operation plan for a
// registered connector tool. The runtime derives operations from its validated
// cached bundle plan; caller capability declarations are never authority.
func (v *Verified) EvaluateV2Plan(toolName, agentID string, subject *SessionSubject, ops []Operation, now time.Time) Decision {
	if decision, denied := v.v2SubjectDenial(subject, now); denied {
		return decision
	}
	var set v2PolicySet
	if err := strictJSON(v.Payload.PolicySet, &set); err != nil {
		return invalidDecision()
	}
	var tool *V2Tool
	for i := range set.Tools {
		if set.Tools[i].Name == toolName {
			tool = &set.Tools[i]
			break
		}
	}
	if tool == nil || tool.Route.ConnectorBinding == nil || tool.Route.ConnectorBinding.ConnectorID == "" || planValid(*tool) != nil {
		return Decision{Outcome: "deny", ReasonCode: "tool_not_registered", Reason: "connector tool is not registered in the policy bundle"}
	}
	if agentID == "" {
		return Decision{Outcome: "deny", ReasonCode: "agent_not_assigned", Reason: "agent is not assigned to this installation"}
	}
	connectorID := tool.Route.ConnectorBinding.ConnectorID
	if !v.agentBound(agentID, connectorID, tool.RequiredCapabilities) {
		return Decision{Outcome: "deny", ReasonCode: "agent_not_bound_to_connector", Reason: "agent binding does not cover the tool's required capabilities"}
	}
	if len(ops) == 0 || len(tool.RequiredCapabilities) == 0 {
		return Decision{Outcome: "deny", ReasonCode: "operation_plan_invalid", Reason: "the operation plan is empty"}
	}
	covered := make(map[string]bool, len(tool.RequiredCapabilities))
	groups, roles := membershipValues(v.Payload.Identities, subject.UserID, now)
	groups = append(groups, subject.Groups...)
	for _, op := range ops {
		if op.ConnectorID != connectorID || op.Capability == "" || op.Action == "" {
			return Decision{Outcome: "deny", ReasonCode: "operation_plan_invalid", Reason: "the operation plan contains a malformed operation"}
		}
		if !containsExact(tool.RequiredCapabilities, op.Capability) || !v.operationResourceRegistered(tool, op) {
			return Decision{Outcome: "deny", ReasonCode: "operation_plan_invalid", Reason: "the operation does not match the tool registration"}
		}
		covered[op.Capability] = true
		request := map[string]string{"source": tool.Source, "connector_id": connectorID, "tool": tool.Name, "capability": op.Capability}
		if op.ResourceType != "" {
			request["resource_type"] = op.ResourceType
		}
		if op.ResourceID != "" {
			request["resource_id"] = op.ResourceID
		}
		if op.ResourceParent != "" {
			request["resource_parent"] = op.ResourceParent
		}
		decision := evaluateV2Check(set.Checks, request, subject, groups, roles, agentID, v.Payload.Audience.OrgID, v.Payload.Audience.WorkspaceID)
		if decision.Outcome != "allow" {
			return decision
		}
	}
	for _, required := range tool.RequiredCapabilities {
		if !covered[required] {
			return Decision{Outcome: "deny", ReasonCode: "operation_plan_invalid", Reason: "the operation plan omits a required capability"}
		}
	}
	return Decision{Outcome: "allow", ReasonCode: "policy_permit"}
}

func (v *Verified) agentBound(agentID, connectorID string, required []string) bool {
	var bindings []struct {
		AgentID      string   `json:"agent_id"`
		ConnectorID  string   `json:"connector_id"`
		Capabilities []string `json:"capabilities"`
	}
	if json.Unmarshal(v.Payload.AgentConnectorBindings, &bindings) != nil {
		return false
	}
	for _, binding := range bindings {
		if binding.AgentID != agentID || binding.ConnectorID != connectorID {
			continue
		}
		for _, capability := range required {
			if !containsExact(binding.Capabilities, capability) {
				return false
			}
		}
		return true
	}
	return false
}

func (v *Verified) operationResourceRegistered(tool *V2Tool, op Operation) bool {
	if op.ResourceType == "" {
		return op.ResourceID == "" && op.ResourceParent == "" && op.ResourceParentType == ""
	}
	for _, rt := range tool.ResourceTypes {
		if rt.Type == op.ResourceType && rt.ParentType == op.ResourceParentType {
			if rt.ParentType == "" && (op.ResourceParent != "" || op.ResourceParentType != "") {
				continue
			}
			if rt.ParentType != "" && (op.ResourceParent == "" || op.ResourceParentType != rt.ParentType) {
				continue
			}
			return true
		}
	}
	return false
}

func evaluateV2Check(checks []v2Check, request map[string]string, subject *SessionSubject, groups, roles []string, agentID, orgID, workspaceID string) Decision {
	matching := make([]v2Check, 0, len(checks))
	call := Call{AgentID: agentID, UserID: subject.UserID, Email: subject.Email}
	for _, check := range checks {
		// A direct capability permit is a separate compatibility grant and
		// cannot create the missing tool-call permit. Direct capability denies
		// remain effective against the operation.
		if check.Selector["tool"] == "" && check.Effect == policyPermitAction {
			continue
		}
		if check.Effect != policyPermitAction && check.Effect != "deny" || !v2ScopeMatches(check.Scope, agentID, orgID, workspaceID) {
			continue
		}
		if !srcMatch(check.Selector["source_principal"], call, groups, roles) || !v2SelectorMatches(check.Selector, request) {
			continue
		}
		matching = append(matching, check)
	}
	if len(matching) == 0 {
		return Decision{Outcome: "deny", ReasonCode: "no_matching_policy", Reason: "no policy check matched the operation"}
	}
	sort.Slice(matching, func(i, j int) bool {
		if matching[i].Priority != matching[j].Priority {
			return matching[i].Priority > matching[j].Priority
		}
		if matching[i].Effect != matching[j].Effect {
			return matching[i].Effect == "deny"
		}
		return matching[i].PolicyID < matching[j].PolicyID
	})
	winner := matching[0]
	if winner.Effect == "deny" {
		return Decision{Outcome: "deny", PolicyID: winner.PolicyID, ReasonCode: "policy_deny", Reason: "policy denied"}
	}
	return Decision{Outcome: "allow", PolicyID: winner.PolicyID}
}

func v2ScopeMatches(scope struct {
	OrgID       *string `json:"org_id"`
	WorkspaceID *string `json:"workspace_id"`
	AgentID     *string `json:"agent_id"`
}, agentID, orgID, workspaceID string) bool {
	return (scope.OrgID == nil || *scope.OrgID == orgID) && (scope.WorkspaceID == nil || *scope.WorkspaceID == workspaceID) && (scope.AgentID == nil || *scope.AgentID == agentID)
}

func v2SelectorMatches(selector map[string]string, request map[string]string) bool {
	for dimension, expected := range selector {
		if dimension == "source_principal" {
			continue
		}
		actual, ok := request[dimension]
		if !ok || expected == "" || (expected != "*" && expected != actual) {
			return false
		}
	}
	return true
}

func membershipValues(identities []Identity, userID string, now time.Time) (groups, roles []string) {
	groupSet, roleSet := map[string]bool{}, map[string]bool{}
	for _, identity := range identities {
		if identity.UserID != userID || identity.Status != SubjectLinked || (!identity.ExpireTime.IsZero() && !now.Before(identity.ExpireTime)) {
			continue
		}
		for _, group := range identity.Groups {
			groupSet[group] = true
		}
		for _, role := range identity.OrgRoles {
			roleSet[role] = true
		}
	}
	for group := range groupSet {
		groups = append(groups, group)
	}
	for role := range roleSet {
		roles = append(roles, role)
	}
	return groups, roles
}

func containsExact(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
