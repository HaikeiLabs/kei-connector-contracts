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

package contract

import (
	"errors"
	"fmt"
)

// AccountModel says whose account an OAuth-style connector acts as. It is an
// explicit, per-connector choice.
type AccountModel string

const (
	// AccountModelPerUser: each Kei user connects their own account, and an
	// invocation uses the invoking user's token. The connector declares no
	// subject.
	AccountModelPerUser AccountModel = "per_user"
	// AccountModelShared: an admin connects one account for the connector
	// instance, and every user the policy allows reads it. The connector's
	// subject is SharedSubject(connector id).
	AccountModelShared AccountModel = "shared"
	// AccountModelDomainDelegation: Google Workspace domain-wide delegation.
	// A service-account key is the connector's opaque_ref credential and the
	// impersonated mailbox is the impersonate_email config field.
	AccountModelDomainDelegation AccountModel = "domain_delegation"
)

var accountModels = map[Provider][]AccountModel{
	ProviderGmail:  {AccountModelPerUser, AccountModelShared, AccountModelDomainDelegation},
	ProviderGoogle: {AccountModelPerUser, AccountModelShared, AccountModelDomainDelegation},
	ProviderLinear: {AccountModelPerUser, AccountModelShared},
	ProviderGitHub: {AccountModelPerUser, AccountModelShared},
}

// AccountModelsFor returns the account models a provider allows, in the order
// a setup screen offers them; the first is the default. Providers without
// user-delegated accounts return nil.
func AccountModelsFor(provider Provider) []AccountModel {
	models := accountModels[provider]
	if len(models) == 0 {
		return nil
	}
	result := make([]AccountModel, len(models))
	copy(result, models)
	return result
}

// SharedSubject is the service subject a shared or delegated connector's
// credential is bound to: the connector itself, not the admin who set it up.
func SharedSubject(connectorID string) string {
	return "connector:" + connectorID
}

type Status string

const (
	StatusPending   Status = "pending"
	StatusActive    Status = "active"
	StatusSuspended Status = "suspended"
	StatusRevoked   Status = "revoked"
	StatusFailed    Status = "failed"
)

type Action string

const (
	ActionRead    Action = "read"
	ActionCreate  Action = "create"
	ActionUpdate  Action = "update"
	ActionDelete  Action = "delete"
	ActionComment Action = "comment"
)

// Capability is the only provider operation a connector client may expose.
// It describes the credential surface: which operations a connector's
// credential can perform. Data-access and tool-call authorization are decided
// by ABAC policy in the control plane before the runtime, not by capability
// flags on the connector.
type Capability struct {
	Name        string `json:"name"`
	Action      Action `json:"action"`
	Description string `json:"description"`
}

// CredentialSubject is whose credential an invocation uses, for audit: the
// invoking user for per_user, the connector's service subject for shared,
// delegated, and opaque_ref connectors, and the declared subject for a legacy
// oauth connector.
func CredentialSubject(m Metadata, in Invocation) string {
	switch {
	case m.AccountModel == AccountModelPerUser:
		return in.Subject
	case m.AccountModel == "" && m.EffectiveCredentialSource() == CredentialSourceOAuth:
		return m.Subject
	default:
		return SharedSubject(m.ID)
	}
}

// validateCredentialBinding enforces the typed credential source, the subject
// binding, and the account model. It never derives any of them from
// credential_ref.
func (m Metadata) validateCredentialBinding() error {
	source := m.EffectiveCredentialSource()
	if source != CredentialSourceOAuth && source != CredentialSourceOpaqueRef {
		return fmt.Errorf("unsupported credential_source %q", m.CredentialSource)
	}
	if m.Subject != "" {
		if err := validID("subject", m.Subject); err != nil {
			return err
		}
	}
	if (m.Provider == ProviderTito || m.Provider == ProviderDiscord || m.Provider == ProviderGrafana) && source != CredentialSourceOpaqueRef {
		return fmt.Errorf("%s connectors must use an opaque_ref credential source", m.Provider)
	}
	if m.AccountModel == "" {
		if source == CredentialSourceOAuth && m.Subject == "" {
			return errors.New("oauth connectors must declare the subject they act for")
		}
		return nil
	}
	if !containsAccountModel(accountModels[m.Provider], m.AccountModel) {
		return fmt.Errorf("account_model %q is not allowed for provider %q", m.AccountModel, m.Provider)
	}
	switch m.AccountModel {
	case AccountModelPerUser:
		if source != CredentialSourceOAuth {
			return errors.New("per_user connectors must use an oauth credential source")
		}
		if m.Subject != "" {
			return errors.New("per_user connectors act for the invoking user and must not declare a subject")
		}
	case AccountModelShared:
		if source != CredentialSourceOAuth {
			return errors.New("shared connectors must use an oauth credential source")
		}
		if m.Subject != SharedSubject(m.ID) {
			return errors.New("shared connectors must be bound to their connector subject")
		}
	case AccountModelDomainDelegation:
		if source != CredentialSourceOpaqueRef {
			return errors.New("domain_delegation connectors must use an opaque_ref credential source")
		}
		if m.Subject != "" {
			return errors.New("domain_delegation connectors must not declare a subject")
		}
	}
	return nil
}

func containsAccountModel(models []AccountModel, want AccountModel) bool {
	for _, model := range models {
		if model == want {
			return true
		}
	}
	return false
}
