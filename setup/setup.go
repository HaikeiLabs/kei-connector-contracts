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

package setup

import (
	"fmt"
	"net/mail"
	"net/url"
	"regexp"

	"github.com/HaikeiLabs/kei-connector-contracts/contract"
)

// SetupSchemaVersion identifies the connector setup schema that the CLI and the
// web UI render. schemas/connector-setup.v1.json is its JSON export.
const SetupSchemaVersion = "kei.connector-setup/v1"

// SetupFieldType is how a setup field is entered and validated.
type SetupFieldType string

const (
	SetupFieldString     SetupFieldType = "string"
	SetupFieldEmail      SetupFieldType = "email"
	SetupFieldHTTPSURL   SetupFieldType = "https_url"
	SetupFieldBool       SetupFieldType = "bool"
	SetupFieldStringList SetupFieldType = "string_list"
	// SetupFieldServiceAccountJSON is a Google service-account key file. The
	// client checks that it is JSON with type=service_account, client_email,
	// and private_key before sealing it.
	SetupFieldServiceAccountJSON SetupFieldType = "service_account_json"
	// SetupFieldHTTPAPI is the http_api block, stored in contract.Metadata.HTTPAPI and
	// validated by contract.HTTPAPI.Validate.
	SetupFieldHTTPAPI SetupFieldType = "http_api"
)

// SetupFieldLocation is where a field's value is stored.
type SetupFieldLocation string

const (
	// SetupLocationConfig: non-secret, stored in contract.Metadata.Config.
	SetupLocationConfig SetupFieldLocation = "config"
	// SetupLocationCredential: secret. It is sealed to the runtime
	// installation and written to the tenant secret manager; the connector
	// keeps only CredentialRef. The control plane never stores the value.
	SetupLocationCredential SetupFieldLocation = "credential"
	// SetupLocationHTTPAPI: stored in contract.Metadata.HTTPAPI.
	SetupLocationHTTPAPI SetupFieldLocation = "http_api"
)

// SetupField is one input on a connector setup screen.
type SetupField struct {
	Name     string             `json:"name"`
	Label    string             `json:"label"`
	Type     SetupFieldType     `json:"type"`
	Location SetupFieldLocation `json:"location"`
	Required bool               `json:"required"`
	// Secret is true exactly when Location is credential.
	Secret    bool   `json:"secret"`
	Pattern   string `json:"pattern,omitempty"`
	MinLength int    `json:"min_length,omitempty"`
	MaxLength int    `json:"max_length,omitempty"`
	Default   any    `json:"default,omitempty"`
	// AccountModels limits the field to these models; empty means every model.
	AccountModels []contract.AccountModel `json:"account_models,omitempty"`
}

// SetupAuth is one way a connector can authenticate.
type SetupAuth struct {
	CredentialSource    contract.CredentialSource `json:"credential_source"`
	AccountModels       []contract.AccountModel   `json:"account_models,omitempty"`
	DefaultAccountModel contract.AccountModel     `json:"default_account_model,omitempty"`
}

// SetupSchema is everything a setup screen needs for one provider.
type SetupSchema struct {
	Schema   string            `json:"schema"`
	Provider contract.Provider `json:"provider"`
	Auth     []SetupAuth       `json:"auth"`
	Fields   []SetupField      `json:"fields"`
}

// FieldsFor returns the fields that apply to an account model.
func (s SetupSchema) FieldsFor(model contract.AccountModel) []SetupField {
	var fields []SetupField
	for _, f := range s.Fields {
		if len(f.AccountModels) == 0 || containsAccountModel(f.AccountModels, model) {
			fields = append(fields, f)
		}
	}
	return fields
}

var delegationOnly = []contract.AccountModel{contract.AccountModelDomainDelegation}

func googleAuth() []SetupAuth {
	return []SetupAuth{
		{CredentialSource: contract.CredentialSourceOAuth, AccountModels: []contract.AccountModel{contract.AccountModelPerUser, contract.AccountModelShared}, DefaultAccountModel: contract.AccountModelPerUser},
		{CredentialSource: contract.CredentialSourceOpaqueRef, AccountModels: delegationOnly},
	}
}

func oauthOnlyAuth() []SetupAuth {
	return []SetupAuth{{CredentialSource: contract.CredentialSourceOAuth, AccountModels: []contract.AccountModel{contract.AccountModelPerUser, contract.AccountModelShared}, DefaultAccountModel: contract.AccountModelPerUser}}
}

func sharedSecretAuth() []SetupAuth {
	return []SetupAuth{{CredentialSource: contract.CredentialSourceOpaqueRef}}
}

func delegationFields() []SetupField {
	return []SetupField{
		{Name: "impersonate_email", Label: "Mailbox to impersonate", Type: SetupFieldEmail, Location: SetupLocationConfig, Required: true, MaxLength: 254, AccountModels: delegationOnly},
		{Name: "service_account_key", Label: "Service account key (JSON)", Type: SetupFieldServiceAccountJSON, Location: SetupLocationCredential, Required: true, Secret: true, MaxLength: 16384, AccountModels: delegationOnly},
	}
}

// setupSchemas lists the connectors Kei builds, in display order. Providers
// absent here (s3 and the finance providers) cannot be set up.
var setupSchemas = []SetupSchema{
	{Schema: SetupSchemaVersion, Provider: contract.ProviderGmail, Auth: googleAuth(), Fields: append([]SetupField{
		{Name: "include_body", Label: "Return message bodies from message.get", Type: SetupFieldBool, Location: SetupLocationConfig, Default: false},
	}, delegationFields()...)},
	{Schema: SetupSchemaVersion, Provider: contract.ProviderGoogle, Auth: googleAuth(), Fields: append([]SetupField{
		{Name: "drive_id", Label: "Drive ID", Type: SetupFieldString, Location: SetupLocationConfig, Pattern: `^[A-Za-z0-9_-]{1,128}$`},
	}, delegationFields()...)},
	{Schema: SetupSchemaVersion, Provider: contract.ProviderLinear, Auth: oauthOnlyAuth(), Fields: []SetupField{}},
	{Schema: SetupSchemaVersion, Provider: contract.ProviderGitHub, Auth: oauthOnlyAuth(), Fields: []SetupField{}},
	{Schema: SetupSchemaVersion, Provider: contract.ProviderTito, Auth: sharedSecretAuth(), Fields: []SetupField{
		{Name: "account_slug", Label: "Tito account slug", Type: SetupFieldString, Location: SetupLocationConfig, Required: true, Pattern: `^[A-Za-z0-9][A-Za-z0-9_-]{0,127}$`},
		{Name: "api_token", Label: "Tito API token", Type: SetupFieldString, Location: SetupLocationCredential, Required: true, Secret: true, Pattern: `^\S+$`, MinLength: 20, MaxLength: 256},
	}},
	// Notion: an internal integration token (ntn_…, or secret_… for older
	// integrations). What it can read is the pages shared with the
	// integration in Notion, so there is no non-secret config.
	{Schema: SetupSchemaVersion, Provider: contract.ProviderNotion, Auth: sharedSecretAuth(), Fields: []SetupField{
		{Name: "api_token", Label: "Notion internal integration token", Type: SetupFieldString, Location: SetupLocationCredential, Required: true, Secret: true, Pattern: `^(ntn|secret)_[A-Za-z0-9]+$`, MinLength: 20, MaxLength: 256},
	}},
	{Schema: SetupSchemaVersion, Provider: contract.ProviderCRM, Auth: sharedSecretAuth(), Fields: []SetupField{
		{Name: "base_url", Label: "CRM Worker origin", Type: SetupFieldHTTPSURL, Location: SetupLocationConfig, Required: true, MaxLength: 2048},
		{Name: "assertion_audience", Label: "Service assertion audience", Type: SetupFieldString, Location: SetupLocationConfig, Required: true, Pattern: `^[A-Za-z0-9][A-Za-z0-9._:/-]{0,254}$`},
		{Name: "allowed_actions", Label: "Allowed actions", Type: SetupFieldStringList, Location: SetupLocationConfig, Pattern: `^(read|create|update)$`, Default: []string{"read"}},
		{Name: "allowed_resources", Label: "Allowed resources", Type: SetupFieldStringList, Location: SetupLocationConfig, Pattern: `^(leads|customers|investors)(/\*)?$`, Default: []string{"leads", "customers"}},
		{Name: "api_key", Label: "CRM connector API key", Type: SetupFieldString, Location: SetupLocationCredential, Required: true, Secret: true, Pattern: `^\S+$`, MinLength: 32, MaxLength: 512},
	}},
	{Schema: SetupSchemaVersion, Provider: contract.ProviderHTTPAPI, Auth: sharedSecretAuth(), Fields: []SetupField{
		{Name: "http_api", Label: "HTTP API registration", Type: SetupFieldHTTPAPI, Location: SetupLocationHTTPAPI, Required: true},
		{Name: "api_key", Label: "API key", Type: SetupFieldString, Location: SetupLocationCredential, Required: true, Secret: true, Pattern: `^\S+$`, MinLength: 1, MaxLength: 4096},
	}},
}

// SetupSchemas returns every provider's setup schema in display order.
func SetupSchemas() []SetupSchema {
	result := make([]SetupSchema, len(setupSchemas))
	copy(result, setupSchemas)
	return result
}

// SetupSchemaFor returns a provider's setup schema, or false when the provider
// cannot be set up.
func SetupSchemaFor(provider contract.Provider) (SetupSchema, bool) {
	for _, s := range setupSchemas {
		if s.Provider == provider {
			return s, true
		}
	}
	return SetupSchema{}, false
}

// ValidateConfig checks a connector's non-secret config against its setup
// schema for the given account model: unknown fields, secret or http_api
// fields placed in config, fields for another account model, missing required
// fields, and type or pattern violations are all rejected. Errors name the
// field and never echo a submitted value.
func ValidateConfig(provider contract.Provider, model contract.AccountModel, config map[string]any) error {
	schema, ok := SetupSchemaFor(provider)
	if !ok {
		return fmt.Errorf("provider %q has no setup schema", provider)
	}
	allowed := contract.AccountModelsFor(provider)
	if model != "" && !containsAccountModel(allowed, model) {
		return fmt.Errorf("account_model %q is not allowed for provider %q", model, provider)
	}
	applicable := map[string]SetupField{}
	for _, f := range schema.FieldsFor(model) {
		applicable[f.Name] = f
	}
	for name, value := range config {
		f, ok := applicable[name]
		if !ok {
			if isSchemaField(schema, name) {
				return fmt.Errorf("config field %q does not apply to account model %q", name, model)
			}
			return fmt.Errorf("config field %q is not defined for provider %q", name, provider)
		}
		if f.Location != SetupLocationConfig {
			return fmt.Errorf("field %q must not be stored in config", name)
		}
		if err := validateSetupValue(f, value); err != nil {
			return err
		}
	}
	for _, f := range applicable {
		if _, present := config[f.Name]; f.Required && f.Location == SetupLocationConfig && !present {
			return fmt.Errorf("config field %q is required", f.Name)
		}
	}
	return nil
}

// ValidateMetadata is the complete connector check the tenant runtime runs
// before executing: the contract's structural, credential-source, and
// account-model rules (Metadata.Validate), the executable check that every
// declared capability is one the runtime has code for
// (Metadata.ValidateExecutable), then the provider's setup-schema rules for
// Metadata.Config. A control plane that must only store connectors the runtime
// can execute runs it on create. A legacy connector (no account model)
// without config has nothing to check beyond the capabilities.
func ValidateMetadata(m contract.Metadata) error {
	if err := m.Validate(); err != nil {
		return err
	}
	if err := m.ValidateExecutable(); err != nil {
		return err
	}
	if m.AccountModel == "" && m.Config == nil {
		return nil
	}
	return ValidateConfig(m.Provider, m.AccountModel, m.Config)
}

func isSchemaField(schema SetupSchema, name string) bool {
	for _, f := range schema.Fields {
		if f.Name == name {
			return true
		}
	}
	return false
}

func validateSetupValue(f SetupField, value any) error {
	invalid := fmt.Errorf("config field %q is invalid", f.Name)
	switch f.Type {
	case SetupFieldBool:
		if _, ok := value.(bool); !ok {
			return invalid
		}
		return nil
	case SetupFieldStringList:
		items, ok := stringList(value)
		if !ok {
			return invalid
		}
		for _, item := range items {
			if !validSetupString(f, item) {
				return invalid
			}
		}
		return nil
	}
	s, ok := value.(string)
	if !ok || !validSetupString(f, s) {
		return invalid
	}
	switch f.Type {
	case SetupFieldEmail:
		addr, err := mail.ParseAddress(s)
		if err != nil || addr.Name != "" || addr.Address != s {
			return invalid
		}
	case SetupFieldHTTPSURL:
		u, err := url.Parse(s)
		if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || u.Opaque != "" {
			return invalid
		}
	}
	return nil
}

func validSetupString(f SetupField, s string) bool {
	if s == "" || len(s) < f.MinLength || (f.MaxLength > 0 && len(s) > f.MaxLength) {
		return false
	}
	return f.Pattern == "" || regexp.MustCompile(f.Pattern).MatchString(s)
}

func stringList(value any) ([]string, bool) {
	switch v := value.(type) {
	case []string:
		return v, len(v) > 0
	case []any:
		items := make([]string, 0, len(v))
		for _, item := range v {
			s, ok := item.(string)
			if !ok {
				return nil, false
			}
			items = append(items, s)
		}
		return items, len(items) > 0
	}
	return nil, false
}

func containsAccountModel(models []contract.AccountModel, want contract.AccountModel) bool {
	for _, model := range models {
		if model == want {
			return true
		}
	}
	return false
}
