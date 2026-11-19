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
	"strings"
	"testing"

	"github.com/HaikeiLabs/kei-connector-contracts/contract"
)

// A postgres connector (Supabase included) takes one secret, the DSN
// of a read-only role, sealed to the credential location. It has no
// non-secret config and no OAuth.
func TestPostgresSetupIsOneSecretDSN(t *testing.T) {
	s, ok := SetupSchemaFor(contract.ProviderPostgres)
	if !ok {
		t.Fatal("postgres has no setup schema")
	}
	if len(s.Auth) != 1 || s.Auth[0].CredentialSource != contract.CredentialSourceOpaqueRef || len(s.Auth[0].AccountModels) != 0 {
		t.Fatalf("postgres auth = %+v, want opaque_ref only", s.Auth)
	}
	if len(s.Fields) != 1 {
		t.Fatalf("postgres fields = %+v, want one", s.Fields)
	}
	dsn := s.Fields[0]
	if dsn.Name != "dsn" || dsn.Type != SetupFieldString || dsn.Location != SetupLocationCredential || !dsn.Secret || !dsn.Required {
		t.Errorf("postgres dsn field = %+v", dsn)
	}
	if !strings.Contains(strings.ToLower(dsn.Label), "read-only") || !strings.Contains(dsn.Label, "Supabase") {
		t.Errorf("postgres dsn label %q must name the read-only role and Supabase", dsn.Label)
	}
	if err := ValidateConfig(contract.ProviderPostgres, "", map[string]any{"dsn": "host=VALUE"}); err == nil {
		t.Error("postgres dsn in config was accepted")
	}
	m := connectorMetadata(contract.ProviderPostgres, contract.CredentialSourceOpaqueRef, "", nil)
	m.Capabilities = contract.CapabilitiesFor(contract.ProviderPostgres)
	if err := ValidateMetadata(m); err != nil {
		t.Errorf("postgres connector rejected: %v", err)
	}
}
