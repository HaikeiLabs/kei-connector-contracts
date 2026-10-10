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

package providers

import (
	"fmt"
	"sort"

	"github.com/HaikeiLabs/kei-connector-contracts/contract"
)

// CapabilityOAuthScopes is the reviewed provider contract linking executable
// capabilities to OAuth scopes. An absent provider or capability has no mapping;
// callers must not infer one from a capability name or action.
var CapabilityOAuthScopes = map[contract.Provider]map[string][]string{
	contract.ProviderGoogle: {
		"drive.search":        {"https://www.googleapis.com/auth/drive.readonly"},
		"drive.metadata.read": {"https://www.googleapis.com/auth/drive.readonly"},
		"docs.read":           {"https://www.googleapis.com/auth/drive.readonly"},
	},
	contract.ProviderGmail: {
		"message.search": {GmailReadonlyScope},
		"message.get":    {GmailReadonlyScope},
	},
	// Linear scopes are its OAuth app scopes, requested comma-separated.
	// Writes map to the narrowest scope that covers them: issues:create and
	// comments:create rather than the broad write scope, which only
	// issue.update needs.
	contract.ProviderLinear: {
		"team.read":      {LinearReadScope},
		"project.read":   {LinearReadScope},
		"cycle.read":     {LinearReadScope},
		"issue.read":     {LinearReadScope},
		"issue.create":   {LinearIssuesCreateScope},
		"comment.create": {LinearCommentsCreateScope},
		"issue.update":   {LinearWriteScope},
	},
}

// Linear OAuth scopes (https://linear.app/developers/oauth-2-0-authentication).
const (
	LinearReadScope           = "read"
	LinearWriteScope          = "write"
	LinearIssuesCreateScope   = "issues:create"
	LinearCommentsCreateScope = "comments:create"
)

// DeclaresOAuthScopes reports whether the OAuth scope preflight applies to a
// connector: its provider has a scope map, and its Scopes name at least one of
// that provider's mapped OAuth scopes. Scopes that are only descriptive, such
// as capability names or the runtime's session placeholder, and providers
// without a map (GitHub) skip the preflight. This is the interim rule (owner
// decision 2026-10-10, option C); the target is a preflight against the
// user's actual grant, not connector metadata.
func DeclaresOAuthScopes(meta contract.Metadata) bool {
	mappings, ok := CapabilityOAuthScopes[meta.Provider]
	if !ok {
		return false
	}
	known := map[string]struct{}{}
	for _, scopes := range mappings {
		for _, scope := range scopes {
			known[scope] = struct{}{}
		}
	}
	for _, scope := range meta.Scopes {
		if _, ok := known[scope]; ok {
			return true
		}
	}
	return false
}

// RequiredOAuthScopes returns the sorted, deduplicated least-privilege union
// for explicitly mapped capabilities. Unknown capabilities fail closed.
func RequiredOAuthScopes(provider contract.Provider, capabilities []contract.Capability) ([]string, error) {
	mappings, ok := CapabilityOAuthScopes[provider]
	if !ok {
		if len(capabilities) == 0 {
			return []string{}, nil
		}
		return nil, fmt.Errorf("OAuth scope mappings are not declared for provider %q", provider)
	}
	set := map[string]struct{}{}
	for _, capability := range capabilities {
		scopes, ok := mappings[capability.Name]
		if !ok {
			return nil, fmt.Errorf("OAuth scopes are not declared for %s capability %q", provider, capability.Name)
		}
		for _, scope := range scopes {
			set[scope] = struct{}{}
		}
	}
	return sortedScopeSet(set), nil
}

// MissingOAuthScopes compares required scopes with an existing OAuth grant.
func MissingOAuthScopes(required, granted []string) []string {
	grants := make(map[string]struct{}, len(granted))
	for _, scope := range granted {
		grants[scope] = struct{}{}
	}
	missing := make([]string, 0)
	for _, scope := range required {
		if _, ok := grants[scope]; !ok {
			missing = append(missing, scope)
		}
	}
	sort.Strings(missing)
	return missing
}

// ValidateGrantedOAuthScopes ensures a connector's stored grant covers its
// enabled capabilities. It is intended for consent and runtime preflight,
// before provider execution; it does not perform policy evaluation or I/O.
func ValidateGrantedOAuthScopes(provider contract.Provider, capabilities []contract.Capability, granted []string) error {
	required, err := RequiredOAuthScopes(provider, capabilities)
	if err != nil {
		return err
	}
	missing := MissingOAuthScopes(required, granted)
	if len(missing) != 0 {
		return fmt.Errorf("OAuth grant is missing required scopes: %v", missing)
	}
	return nil
}

// OAuthReconsentScopes returns only scopes not already granted for the new
// enabled-capability set, supporting incremental re-consent.
func OAuthReconsentScopes(provider contract.Provider, capabilities []contract.Capability, granted []string) ([]string, error) {
	required, err := RequiredOAuthScopes(provider, capabilities)
	if err != nil {
		return nil, err
	}
	return MissingOAuthScopes(required, granted), nil
}

func sortedScopeSet(set map[string]struct{}) []string {
	out := make([]string, 0, len(set))
	for scope := range set {
		out = append(out, scope)
	}
	sort.Strings(out)
	return out
}
