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
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/HaikeiLabs/kei-connector-contracts/contract"
)

func investorsMeta(t *testing.T) contract.Metadata {
	return metaFor(t, contract.ProviderCRM, []string{"investors"}, nil)
}

func TestCRMInvestorCapabilitiesAreDefined(t *testing.T) {
	for _, name := range []string{"investor.list", "investor.read", "investor.update"} {
		found := false
		for _, c := range contract.CapabilitiesFor(contract.ProviderCRM) {
			if c.Name == name {
				found = true
				switch name {
				case "investor.list", "investor.read":
					if c.Action != contract.ActionRead {
						t.Errorf("%s action = %q, want read", name, c.Action)
					}
				case "investor.update":
					if c.Action != contract.ActionUpdate {
						t.Errorf("%s action = %q, want update", name, c.Action)
					}
				}
			}
		}
		if !found {
			t.Errorf("crm does not declare %s", name)
		}
	}
}

func TestCRMInvestorList(t *testing.T) {
	c := NewCRM(crmStore())
	ctx := context.Background()
	inv := invocation("investor.list", contract.ActionRead, "investors")

	all, err := c.Invoke(ctx, investorsMeta(t), inv, InvestorListPayload{})
	if err != nil {
		t.Fatalf("investor list failed: %v", err)
	}
	page := all.Data.(InvestorPage)
	if len(page.Records) != 3 || page.HasMore || page.NextCursor != nil {
		t.Fatalf("investor list = %+v", page)
	}

	diligence, err := c.Invoke(ctx, investorsMeta(t), inv, InvestorListPayload{Stage: "diligence"})
	if err != nil {
		t.Fatalf("stage filter failed: %v", err)
	}
	if got := diligence.Data.(InvestorPage).Records; len(got) != 1 || got[0].ID != "i-2" {
		t.Fatalf("stage filter returned %+v", got)
	}

	byFirm, err := c.Invoke(ctx, investorsMeta(t), inv, InvestorListPayload{Query: "north"})
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if got := byFirm.Data.(InvestorPage).Records; len(got) != 1 || got[0].ID != "i-3" {
		t.Fatalf("query returned %+v", got)
	}
}

func TestCRMInvestorListPaginates(t *testing.T) {
	c := NewCRM(crmStore())
	ctx := context.Background()
	inv := invocation("investor.list", contract.ActionRead, "investors")

	first, err := c.Invoke(ctx, investorsMeta(t), inv, InvestorListPayload{Limit: 2})
	if err != nil {
		t.Fatalf("first page failed: %v", err)
	}
	p1 := first.Data.(InvestorPage)
	if len(p1.Records) != 2 || !p1.HasMore || p1.NextCursor == nil {
		t.Fatalf("first page = %+v", p1)
	}
	second, err := c.Invoke(ctx, investorsMeta(t), inv, InvestorListPayload{Limit: 2, Cursor: *p1.NextCursor})
	if err != nil {
		t.Fatalf("second page failed: %v", err)
	}
	p2 := second.Data.(InvestorPage)
	if len(p2.Records) != 1 || p2.HasMore || p2.NextCursor != nil {
		t.Fatalf("second page = %+v", p2)
	}
	seen := map[string]bool{}
	for _, r := range append(p1.Records, p2.Records...) {
		if seen[r.ID] {
			t.Fatalf("investor %s returned twice", r.ID)
		}
		seen[r.ID] = true
	}
}

func TestCRMInvestorListRejectsInvalidInput(t *testing.T) {
	c := NewCRM(crmStore())
	inv := invocation("investor.list", contract.ActionRead, "investors")
	for name, p := range map[string]InvestorListPayload{
		"unknown stage":  {Stage: "won"},
		"long query":     {Query: strings.Repeat("q", 201)},
		"negative limit": {Limit: -1},
		"large limit":    {Limit: 101},
		"long cursor":    {Cursor: strings.Repeat("c", 129)},
	} {
		if _, err := c.Invoke(context.Background(), investorsMeta(t), inv, p); err == nil {
			t.Errorf("%s was accepted", name)
		}
	}
	// investor.list lists the collection only.
	if _, err := c.Invoke(context.Background(), investorsMeta(t), invocation("investor.list", contract.ActionRead, "investors/i-1"), InvestorListPayload{}); err == nil {
		t.Error("investor.list accepted a single-investor resource")
	}
}

func TestCRMInvestorRead(t *testing.T) {
	c := NewCRM(crmStore())
	ctx := context.Background()
	got, err := c.Invoke(ctx, investorsMeta(t), invocation("investor.read", contract.ActionRead, "investors/i-2"), InvestorReadPayload{})
	if err != nil {
		t.Fatalf("investor read failed: %v", err)
	}
	if inv := got.Data.(CRMInvestor); inv.Name != "Grace Hopper" || inv.Stage != "diligence" {
		t.Fatalf("investor read = %+v", inv)
	}
	for _, resource := range []string{"investors", "investors/", "investors/i-1/notes", "leads/l-1"} {
		if _, err := c.Invoke(ctx, investorsMeta(t), invocation("investor.read", contract.ActionRead, resource), InvestorReadPayload{}); err == nil {
			t.Errorf("investor.read accepted resource %q", resource)
		}
	}
	if _, err := c.Invoke(ctx, investorsMeta(t), invocation("investor.read", contract.ActionRead, "investors/missing"), InvestorReadPayload{}); err == nil {
		t.Error("missing investor was returned")
	}
}

func TestCRMInvestorUpdateStage(t *testing.T) {
	c := NewCRM(crmStore())
	ctx := context.Background()
	stage := "committed"
	got, err := c.Invoke(ctx, investorsMeta(t), invocation("investor.update", contract.ActionUpdate, "investors/i-1"), InvestorUpdatePayload{Stage: &stage})
	if err != nil {
		t.Fatalf("investor update failed: %v", err)
	}
	inv := got.Data.(CRMInvestor)
	if inv.ID != "i-1" || inv.Stage != "committed" {
		t.Fatalf("investor after update = %+v", inv)
	}
	// Verify the change persisted by reading it back.
	read, err := c.Invoke(ctx, investorsMeta(t), invocation("investor.read", contract.ActionRead, "investors/i-1"), InvestorReadPayload{})
	if err != nil {
		t.Fatalf("investor read after update failed: %v", err)
	}
	if read.Data.(CRMInvestor).Stage != "committed" {
		t.Fatalf("investor stage after update = %q, want committed", read.Data.(CRMInvestor).Stage)
	}
}

func TestCRMInvestorUpdateValidatesStage(t *testing.T) {
	c := NewCRM(crmStore())
	inv := invocation("investor.update", contract.ActionUpdate, "investors/i-1")
	invalid := "won"
	for name, p := range map[string]InvestorUpdatePayload{
		"unknown stage": {Stage: &invalid},
	} {
		if _, err := c.Invoke(context.Background(), investorsMeta(t), inv, p); err == nil {
			t.Errorf("%s was accepted", name)
		}
	}
}

func TestCRMInvestorUpdateRejectsMissingResource(t *testing.T) {
	c := NewCRM(crmStore())
	stage := "committed"
	for _, resource := range []string{"investors", "investors/", "investors/i-1/notes", "leads/l-1"} {
		if _, err := c.Invoke(context.Background(), investorsMeta(t), invocation("investor.update", contract.ActionUpdate, resource), InvestorUpdatePayload{Stage: &stage}); err == nil {
			t.Errorf("investor.update accepted resource %q", resource)
		}
	}
}

func TestCRMInvestorUpdateRejectsMissingInvestor(t *testing.T) {
	c := NewCRM(crmStore())
	stage := "committed"
	if _, err := c.Invoke(context.Background(), investorsMeta(t), invocation("investor.update", contract.ActionUpdate, "investors/missing"), InvestorUpdatePayload{Stage: &stage}); err == nil {
		t.Error("missing investor was updated")
	}
}

func TestCRMInvestorUpdateFailsClosedWithoutInvestorBackend(t *testing.T) {
	c := NewCRM(leadOnlyBackend{crmStore()})
	stage := "committed"
	_, err := c.Invoke(context.Background(), investorsMeta(t), invocation("investor.update", contract.ActionUpdate, "investors/i-1"), InvestorUpdatePayload{Stage: &stage})
	if err == nil || err.Error() != "crm backend does not serve investors" {
		t.Fatalf("lead-only backend error = %v", err)
	}
}

func TestCRMInvestorReadsRequireAuthentication(t *testing.T) {
	c := NewCRM(crmStore())
	m := investorsMeta(t)
	m.CredentialRef = "vault/tenant/other/crm"
	ctx := context.Background()
	if _, err := c.Invoke(ctx, m, invocation("investor.list", contract.ActionRead, "investors"), InvestorListPayload{}); err == nil || err.Error() != "crm session is not authenticated" {
		t.Fatalf("unauthenticated investor list error = %v", err)
	}
	if _, err := c.Invoke(ctx, m, invocation("investor.read", contract.ActionRead, "investors/i-1"), InvestorReadPayload{}); err == nil || err.Error() != "crm session is not authenticated" {
		t.Fatalf("unauthenticated investor read error = %v", err)
	}
	stage := "committed"
	if _, err := c.Invoke(ctx, m, invocation("investor.update", contract.ActionUpdate, "investors/i-1"), InvestorUpdatePayload{Stage: &stage}); err == nil || err.Error() != "crm session is not authenticated" {
		t.Fatalf("unauthenticated investor update error = %v", err)
	}
}

// leadOnlyBackend implements CRMBackend without investor support.
type leadOnlyBackend struct{ leads MemoryCRM }

func (b leadOnlyBackend) ListLeads(ctx context.Context, authRef string) ([]CRMLead, error) {
	return b.leads.ListLeads(ctx, authRef)
}

func (b leadOnlyBackend) Lead(ctx context.Context, authRef, id string) (CRMLead, error) {
	return b.leads.Lead(ctx, authRef, id)
}

func TestCRMInvestorReadsFailClosedWithoutInvestorBackend(t *testing.T) {
	c := NewCRM(leadOnlyBackend{crmStore()})
	_, err := c.Invoke(context.Background(), investorsMeta(t), invocation("investor.list", contract.ActionRead, "investors"), InvestorListPayload{})
	if err == nil || err.Error() != "crm backend does not serve investors" {
		t.Fatalf("lead-only backend error = %v", err)
	}
}

func TestCRMInvestorWireShape(t *testing.T) {
	page := InvestorPage{Records: []CRMInvestor{{ID: "i-1", Name: "Ada", Firm: "Acme", Email: "ada@example.com", Stage: "prospect"}}}
	raw, err := json.Marshal(page)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"records", "next_cursor", "has_more"} {
		if _, ok := got[key]; !ok {
			t.Errorf("investor page is missing %q: %s", key, raw)
		}
	}
	record := got["records"].([]any)[0].(map[string]any)
	for _, key := range []string{"id", "name", "firm", "email", "stage", "created_at", "updated_at"} {
		if _, ok := record[key]; !ok {
			t.Errorf("investor record is missing %q: %s", key, raw)
		}
	}
	if len(record) != 7 {
		t.Errorf("investor record has %d fields, want the 7-field allowlist: %s", len(record), raw)
	}
	in, _ := json.Marshal(InvestorListPayload{Stage: "meeting", Query: "a", Limit: 5, Cursor: "c"})
	if string(in) != `{"stage":"meeting","q":"a","limit":5,"cursor":"c"}` {
		t.Errorf("investor.list input = %s", in)
	}
}
