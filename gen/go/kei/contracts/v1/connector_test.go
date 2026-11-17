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

package contracts_test

import (
	"slices"
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"

	contracts "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/contracts/v1"
)

// A connector declared in proto must survive the wire unchanged: the
// manifest and catalog are built from it, and a lost secret flag would turn
// a credential field into a plain one.
func TestConnectorOptionRoundTrip(t *testing.T) {
	want := &contracts.Connector{
		Id:      "github",
		BaseUrl: "https://api.github.com",
		Connections: []*contracts.Connection{{
			Name:         "pat",
			Kind:         contracts.ConnectionKind_CONNECTION_KIND_HTTP_BEARER,
			AuthMethod:   "secret_ref",
			AccountModel: contracts.AccountModel_ACCOUNT_MODEL_SHARED,
			SetupFields:  []*contracts.SetupField{{Name: "token", Type: "string", Secret: true}},
		}},
	}

	b, err := proto.Marshal(want)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got := &contracts.Connector{}
	if err := proto.Unmarshal(b, got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !proto.Equal(want, got) {
		t.Fatalf("round trip changed the connector:\nwant %v\ngot  %v", want, got)
	}
	if !got.GetConnections()[0].GetSetupFields()[0].GetSecret() {
		t.Error("setup field token lost its secret flag")
	}
}

// The operation option gains the fields the generators and the manifest
// read. Every one must round-trip, including the nested semantic mapping and
// the error table.
func TestOperationNewFieldsRoundTrip(t *testing.T) {
	want := &contracts.Operation{
		Id:               "github.issues.get",
		Connector:        "github",
		Class:            contracts.OperationClass_OPERATION_CLASS_READ,
		HttpMethod:       "GET",
		PathTemplate:     "/repos/{owner}/{repository}/issues/{issue_number}",
		Scopes:           []string{"issues:read"},
		MethodName:       "GetIssue",
		Connections:      []string{"user_oauth", "pat"},
		ResourceTemplate: "repos/{owner}/{repository}",
		Errors: []*contracts.ErrorMapping{
			{Status: 404, Kind: contracts.ErrorKind_ERROR_KIND_NOT_FOUND},
			{Status: 429, Kind: contracts.ErrorKind_ERROR_KIND_RATE_LIMITED, Retryable: true},
		},
		Deprecated:       true,
		ReplacedBy:       "github.issues.read",
		GraphqlDocument:  "query { viewer { id } }",
		PositionalHelper: true,
		Semantics: &contracts.Semantics{
			Action:       contracts.SemanticAction_SEMANTIC_ACTION_READ,
			ResourceType: "issue",
			Version:      1,
		},
		CursorField:     "page",
		NextCursorField: "next_page",
	}

	b, err := proto.Marshal(want)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got := &contracts.Operation{}
	if err := proto.Unmarshal(b, got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !proto.Equal(want, got) {
		t.Fatalf("round trip changed the operation:\nwant %v\ngot  %v", want, got)
	}
}

// The option extensions are part of the public contract: a renumbered
// extension silently drops every annotation compiled against the old number.
// 51234 is the existing operation option; 51235 and 51236 are new.
func TestExtensionNumbersAreStable(t *testing.T) {
	cases := []struct {
		name string
		got  int32
		want int32
	}{
		{"operation", int32(contracts.E_Operation.TypeDescriptor().Number()), 51234},
		{"connector", int32(contracts.E_Connector.TypeDescriptor().Number()), 51235},
		{"param", int32(contracts.E_Param.TypeDescriptor().Number()), 51236},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("extension %s = %d, want %d", c.name, c.got, c.want)
		}
	}
	if got := contracts.E_Connector.TypeDescriptor().ContainingMessage().FullName(); got != "google.protobuf.FileOptions" {
		t.Errorf("connector extends %s, want google.protobuf.FileOptions", got)
	}
	if got := contracts.E_Param.TypeDescriptor().ContainingMessage().FullName(); got != "google.protobuf.FieldOptions" {
		t.Errorf("param extends %s, want google.protobuf.FieldOptions", got)
	}
}

// Connectors() is sorted and unique by id, and LookupConnector agrees with
// it; an unknown id is not found rather than a zero-value connector.
func TestConnectorsSortedAndLookupAgrees(t *testing.T) {
	cs := contracts.Connectors()
	if !slices.IsSortedFunc(cs, func(a, b *contracts.Connector) int {
		return strings.Compare(a.GetId(), b.GetId())
	}) {
		t.Error("Connectors() is not sorted by id")
	}
	seen := map[string]bool{}
	for _, c := range cs {
		if seen[c.GetId()] {
			t.Errorf("duplicate connector id %q", c.GetId())
		}
		seen[c.GetId()] = true
		got, ok := contracts.LookupConnector(c.GetId())
		if !ok || !proto.Equal(got, c) {
			t.Errorf("LookupConnector(%q) = %v, %v; want the Connectors() entry", c.GetId(), got, ok)
		}
	}
	if c, ok := contracts.LookupConnector("no.such.connector"); ok || c != nil {
		t.Errorf("LookupConnector(unknown) = %v, %v; want nil, false", c, ok)
	}
}
