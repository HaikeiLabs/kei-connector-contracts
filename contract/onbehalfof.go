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
	"regexp"
	"strings"
)

var guidPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// OnBehalfOf is the OAuth 2.0 on-behalf-of (client-credentials-for-a-user)
// configuration for an http_api connector. The runtime uses client_id and
// client_secret to obtain an access token for the invoking user's identity
// from the identity provider at TenantID, scoped to Resource.
type OnBehalfOf struct {
	// TenantID is the Azure AD / Entra directory (tenant) the user and
	// application are registered in.
	TenantID string `json:"tenant_id"`
	// ClientID is the application (client) registration's application ID.
	ClientID string `json:"client_id"`
	// Resource is the resource the access token is for (the audience).
	Resource string `json:"resource"`
	// Scopes are the individual permissions the connector requests. Each
	// scope must be a resource-scoped name under Resource (e.g.
	// https://resource.example.com/Scope.Name) and must not use the
	// ".default" shortcut.
	Scopes []string `json:"scopes"`
	// ClientSecretRef is the Kei credential reference holding the client
	// secret. It must be of the form kei://connectors/<connector-id>/client_secret.
	ClientSecretRef string `json:"client_secret_ref"`
}

// Validate checks every field of the OnBehalfOf configuration against the
// contract's structural rules. connectorID is the metadata id of the owning
// http_api connector, used to verify ClientSecretRef.
func (o OnBehalfOf) Validate(connectorID string) error {
	if !guidPattern.MatchString(o.TenantID) {
		return errors.New("on_behalf_of tenant_id must be a valid GUID")
	}
	if !guidPattern.MatchString(o.ClientID) {
		return errors.New("on_behalf_of client_id must be a valid GUID")
	}
	if !strings.HasPrefix(o.Resource, "https://") {
		return errors.New("on_behalf_of resource must be an https URL")
	}
	if len(o.Scopes) == 0 {
		return errors.New("on_behalf_of must declare at least one scope")
	}
	prefix := o.Resource + "/"
	for _, scope := range o.Scopes {
		if !strings.HasPrefix(scope, prefix) {
			return fmt.Errorf("on_behalf_of scope %q is not under resource %q", scope, o.Resource)
		}
		if strings.HasSuffix(scope, "/.default") {
			return errors.New("on_behalf_of scopes must not use .default")
		}
	}
	want := fmt.Sprintf("kei://connectors/%s/client_secret", connectorID)
	if o.ClientSecretRef != want {
		return errors.New("on_behalf_of client_secret_ref must reference the connector's own client_secret in the kei://connectors/ namespace")
	}
	return nil
}
