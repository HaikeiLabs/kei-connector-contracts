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
	"errors"
	"strings"
	"testing"
	"time"
)

// descriptorDoc is a v2 bundle whose policy set carries one connector
// descriptor, edited by mutate.
func descriptorDoc(now time.Time, mutate func(c map[string]any)) map[string]any {
	doc := v2BundleDoc(now)
	c := map[string]any{
		"id": "conn-github", "name": "GitHub", "provider": "github",
		"capabilities":      []any{map[string]any{"name": "issues.read", "mutating": false}},
		"credential_source": "opaque_ref", "account_model": "shared",
		"credential_ref": "kei://connectors/conn-github/token", "credential_backend": "vault",
	}
	if mutate != nil {
		mutate(c)
	}
	doc["policy_set"].(map[string]any)["connectors"] = []any{c}
	return doc
}

func descriptorValidate(t *testing.T, doc map[string]any, now time.Time) (*Verified, error) {
	t.Helper()
	SetV2SchemaValidator(func([]byte) error { return nil })
	t.Cleanup(func() { SetV2SchemaValidator(nil) })
	raw, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	return ValidateV2(raw, v2Installation, v2Org, v2Workspace, 1, now)
}

// The descriptor may name the contract connection the installation uses and
// its resolved provider origin; both are non-secret and optional.
func TestConnectorDescriptorCarriesConnectionAndBaseURL(t *testing.T) {
	now := time.Now().UTC()
	doc := descriptorDoc(now, func(c map[string]any) {
		c["connection"] = "pat"
		c["base_url"] = "https://ghe.example.com/api/v3"
	})
	v, err := descriptorValidate(t, doc, now)
	if err != nil {
		t.Fatalf("descriptor with connection and base_url must validate, got %v", err)
	}
	c, ok := v.Connector("conn-github")
	if !ok {
		t.Fatal("connector descriptor not found")
	}
	if c.Connection != "pat" || c.BaseURL != "https://ghe.example.com/api/v3" {
		t.Fatalf("connection/base_url = %q/%q", c.Connection, c.BaseURL)
	}
}

// Both fields are optional: a descriptor from an older catalog still validates
// and decodes with them empty.
func TestConnectorDescriptorWithoutConnectionStillValidates(t *testing.T) {
	now := time.Now().UTC()
	v, err := descriptorValidate(t, descriptorDoc(now, nil), now)
	if err != nil {
		t.Fatalf("descriptor without the new fields must validate, got %v", err)
	}
	c, _ := v.Connector("conn-github")
	if c.Connection != "" || c.BaseURL != "" {
		t.Fatalf("absent fields decoded as %q/%q", c.Connection, c.BaseURL)
	}
	out, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	var keys map[string]any
	_ = json.Unmarshal(out, &keys)
	for _, k := range []string{"connection", "base_url"} {
		if _, has := keys[k]; has {
			t.Errorf("empty %s must be omitted on encode", k)
		}
	}
}

func TestConnectorDescriptorRejectsMalformedConnectionOrBaseURL(t *testing.T) {
	cases := map[string]func(c map[string]any){
		"connection not snake_case":   func(c map[string]any) { c["connection"] = "User OAuth" },
		"connection with a dot":       func(c map[string]any) { c["connection"] = "user.oauth" },
		"connection too long":         func(c map[string]any) { c["connection"] = strings.Repeat("a", 65) },
		"base_url http":               func(c map[string]any) { c["base_url"] = "http://ghe.example.com" },
		"base_url relative":           func(c map[string]any) { c["base_url"] = "/api/v3" },
		"base_url with userinfo":      func(c map[string]any) { c["base_url"] = "https://user:tok@ghe.example.com" },
		"base_url with query":         func(c map[string]any) { c["base_url"] = "https://ghe.example.com/?access_token=x" },
		"base_url with fragment":      func(c map[string]any) { c["base_url"] = "https://ghe.example.com/#x" },
		"base_url without host":       func(c map[string]any) { c["base_url"] = "https:///api" },
		"base_url link-local address": func(c map[string]any) { c["base_url"] = "https://169.254.169.254" },
		"base_url not a string":       func(c map[string]any) { c["base_url"] = 7 },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			now := time.Now().UTC()
			if _, err := descriptorValidate(t, descriptorDoc(now, mutate), now); !errors.Is(err, ErrBundleUnsupported) {
				t.Fatalf("want ErrBundleUnsupported, got %v", err)
			}
		})
	}
}

// The descriptor is non-secret: the closed parser refuses any field that
// would carry a credential value, so a catalog can never ship one.
func TestConnectorDescriptorRejectsCredentialFields(t *testing.T) {
	for _, key := range []string{"credential", "credentials", "api_key", "token", "access_token", "secret", "password", "client_secret", "private_key"} {
		t.Run(key, func(t *testing.T) {
			now := time.Now().UTC()
			doc := descriptorDoc(now, func(c map[string]any) {
				c["connection"] = "pat"
				c[key] = "value"
			})
			if _, err := descriptorValidate(t, doc, now); !errors.Is(err, ErrBundleUnsupported) {
				t.Fatalf("descriptor with %q must be rejected, got %v", key, err)
			}
		})
	}
}
