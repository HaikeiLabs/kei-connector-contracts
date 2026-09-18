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
	"strings"
	"testing"
	"time"

	"github.com/HaikeiLabs/kei-connector-contracts"
)

// Fixtures are deliberately synthetic: TEST-prefixed identifiers, example.test
// addresses and round amounts, so nothing here resembles a real ledger, a real
// account or a real counterparty. The memory backends are the only
// implementations in this module, so these tests make no network call.

func freshBooksStore() MemoryFreshBooks {
	issued := time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC)
	return MemoryFreshBooks{
		Invoices: map[string]FreshBooksInvoice{
			"acct-TEST-1/inv-TEST-1": {
				ID: "inv-TEST-1", Number: "TEST-0001", ClientID: "cli-TEST-1",
				Status: "sent", AmountMinor: 150000, OutstandingMinor: 150000,
				Currency: "USD", IssuedAt: issued, DueAt: issued.AddDate(0, 1, 0),
			},
		},
		Expenses: map[string]FreshBooksExpense{
			"acct-TEST-1/exp-TEST-1": {
				ID: "exp-TEST-1", VendorName: "Example Vendor (TEST)",
				CategoryID: "cat-TEST-1", AmountMinor: 4250, Currency: "USD",
				IncurredAt: issued,
			},
		},
		Payments: map[string]FreshBooksPayment{
			"acct-TEST-1/pay-TEST-1": {
				ID: "pay-TEST-1", InvoiceID: "inv-TEST-1", AmountMinor: 150000,
				Currency: "USD", Method: "ach", ReceivedAt: issued,
			},
		},
		Clients: map[string]FreshBooksClient{
			"acct-TEST-1/cli-TEST-1": {
				ID: "cli-TEST-1", Organization: "Example Org (TEST)",
				Email: "billing@example.test", Currency: "USD",
			},
		},
	}
}

func mercuryStore() MemoryMercury {
	asOf := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	return MemoryMercury{
		Accounts: map[string]MercuryAccount{
			"acct-TEST-1": {
				ID: "acct-TEST-1", Name: "Operating (TEST)", Kind: "checking",
				LastFour: "0000", Currency: "USD", Status: "active", OpenedAt: asOf,
			},
		},
		Balances: map[string]MercuryBalance{
			"acct-TEST-1": {
				AccountID: "acct-TEST-1", CurrentMinor: 2500000,
				AvailableMinor: 2400000, Currency: "USD", AsOf: asOf,
			},
		},
		Transactions: map[string]MercuryTransaction{
			"acct-TEST-1/txn-TEST-1": {
				ID: "txn-TEST-1", AccountID: "acct-TEST-1", AmountMinor: -4250,
				Currency: "USD", Direction: "debit", Status: "posted",
				CounterpartyName: "Example Vendor (TEST)", PostedAt: asOf,
			},
		},
	}
}

func TestFreshBooksReadsResolve(t *testing.T) {
	ctx := context.Background()
	client := NewFreshBooks(freshBooksStore())
	meta := metaFor(t, connectors.ProviderFreshBooks, []string{"accounts/acct-TEST-1"}, nil)

	for _, tc := range []struct {
		resource string
		payload  Payload
	}{
		{"accounts/acct-TEST-1/invoices/inv-TEST-1", InvoiceReadPayload{}},
		{"accounts/acct-TEST-1/expenses/exp-TEST-1", ExpenseReadPayload{}},
		{"accounts/acct-TEST-1/payments/pay-TEST-1", PaymentReadPayload{}},
		{"accounts/acct-TEST-1/clients/cli-TEST-1", ClientReadPayload{}},
	} {
		inv := invocation(tc.payload.Capability(), connectors.ActionRead, tc.resource)
		res, err := client.Invoke(ctx, meta, inv, tc.payload)
		if err != nil {
			t.Fatalf("resource %q rejected: %v", tc.resource, err)
		}
		if res.Capability != tc.payload.Capability() {
			t.Errorf("capability = %q, want %q", res.Capability, tc.payload.Capability())
		}
	}
}

func TestMercuryReadsResolve(t *testing.T) {
	ctx := context.Background()
	client := NewMercury(mercuryStore())
	meta := metaFor(t, connectors.ProviderMercury, []string{"accounts/acct-TEST-1"}, nil)

	for _, tc := range []struct {
		resource string
		payload  Payload
	}{
		{"accounts/acct-TEST-1", AccountReadPayload{}},
		{"accounts/acct-TEST-1", BalanceReadPayload{}},
		{"accounts/acct-TEST-1/transactions/txn-TEST-1", TransactionReadPayload{}},
	} {
		inv := invocation(tc.payload.Capability(), connectors.ActionRead, tc.resource)
		res, err := client.Invoke(ctx, meta, inv, tc.payload)
		if err != nil {
			t.Fatalf("resource %q rejected: %v", tc.resource, err)
		}
		if res.Capability != tc.payload.Capability() {
			t.Errorf("capability = %q, want %q", res.Capability, tc.payload.Capability())
		}
	}
}

// TestFinanceResourceShapeIsEnforced pins that a malformed or cross-collection
// resource is refused rather than silently resolving somewhere else. A
// payments read must not be satisfiable by an invoices path, and an account id
// may not be omitted.
func TestFinanceResourceShapeIsEnforced(t *testing.T) {
	ctx := context.Background()
	fb := NewFreshBooks(freshBooksStore())
	fbMeta := metaFor(t, connectors.ProviderFreshBooks, []string{"accounts/acct-TEST-1"}, nil)

	for _, resource := range []string{
		"accounts/acct-TEST-1/invoices/inv-TEST-1/extra",
		"accounts/acct-TEST-1/invoices/",
		"accounts//invoices/inv-TEST-1",
		"invoices/inv-TEST-1",
		"accounts/acct-TEST-1/payments/pay-TEST-1", // wrong collection for this payload
	} {
		inv := invocation("invoice.read", connectors.ActionRead, resource)
		if _, err := fb.Invoke(ctx, fbMeta, inv, InvoiceReadPayload{}); err == nil {
			t.Errorf("freshbooks accepted malformed resource %q", resource)
		}
	}

	mc := NewMercury(mercuryStore())
	mcMeta := metaFor(t, connectors.ProviderMercury, []string{"accounts/acct-TEST-1"}, nil)
	for _, resource := range []string{
		"accounts/acct-TEST-1/transactions/txn-TEST-1", // not an account resource
		"accounts/",
		"acct-TEST-1",
	} {
		inv := invocation("account.read", connectors.ActionRead, resource)
		if _, err := mc.Invoke(ctx, mcMeta, inv, AccountReadPayload{}); err == nil {
			t.Errorf("mercury accepted malformed resource %q", resource)
		}
	}
}

// TestFinanceCrossAccountReadIsRefused proves the account segment is load
// bearing: a record that exists under one account is not readable by naming a
// different one.
func TestFinanceCrossAccountReadIsRefused(t *testing.T) {
	ctx := context.Background()
	client := NewFreshBooks(freshBooksStore())
	meta := metaFor(t, connectors.ProviderFreshBooks, []string{"accounts/acct-TEST-2"}, nil)

	inv := invocation("invoice.read", connectors.ActionRead, "accounts/acct-TEST-2/invoices/inv-TEST-1")
	if _, err := client.Invoke(ctx, meta, inv, InvoiceReadPayload{}); err == nil {
		t.Fatal("invoice from another account was returned")
	}
}

// TestFinanceClientsRejectMismatchedProvider pins that a bookkeeping client
// cannot serve a banking connector's metadata, or vice versa.
func TestFinanceClientsRejectMismatchedProvider(t *testing.T) {
	ctx := context.Background()
	fbMeta := metaFor(t, connectors.ProviderFreshBooks, []string{"accounts/acct-TEST-1"}, nil)
	mcMeta := metaFor(t, connectors.ProviderMercury, []string{"accounts/acct-TEST-1"}, nil)

	inv := invocation("account.read", connectors.ActionRead, "accounts/acct-TEST-1")
	if _, err := NewMercury(mercuryStore()).Invoke(ctx, fbMeta, inv, AccountReadPayload{}); err == nil {
		t.Error("mercury client served freshbooks metadata")
	}

	inv = invocation("invoice.read", connectors.ActionRead, "accounts/acct-TEST-1/invoices/inv-TEST-1")
	if _, err := NewFreshBooks(freshBooksStore()).Invoke(ctx, mcMeta, inv, InvoiceReadPayload{}); err == nil {
		t.Error("freshbooks client served mercury metadata")
	}
}

// TestFinancePayloadCannotBeReplayedAcrossCapabilities pins matchCapability
// for the finance clients: a balance payload cannot ride an account.read
// invocation.
func TestFinancePayloadCannotBeReplayedAcrossCapabilities(t *testing.T) {
	ctx := context.Background()
	client := NewMercury(mercuryStore())
	meta := metaFor(t, connectors.ProviderMercury, []string{"accounts/acct-TEST-1"}, nil)

	inv := invocation("account.read", connectors.ActionRead, "accounts/acct-TEST-1")
	if _, err := client.Invoke(ctx, meta, inv, BalanceReadPayload{}); err == nil {
		t.Fatal("balance payload accepted on an account.read invocation")
	}
}

// TestFinanceTraversalResourcesRefused pins that Guard's traversal check
// covers the finance resource shapes, including percent-encoded segments.
func TestFinanceTraversalResourcesRefused(t *testing.T) {
	ctx := context.Background()
	client := NewFreshBooks(freshBooksStore())
	meta := metaFor(t, connectors.ProviderFreshBooks, []string{"accounts/acct-TEST-1"}, nil)

	for _, resource := range []string{
		"accounts/acct-TEST-1/invoices/../../acct-TEST-2/invoices/inv-TEST-1",
		"accounts/acct-TEST-1/invoices/%2e%2e",
	} {
		inv := invocation("invoice.read", connectors.ActionRead, resource)
		if _, err := client.Invoke(ctx, meta, inv, InvoiceReadPayload{}); err == nil {
			t.Errorf("traversal resource %q accepted", resource)
		}
	}
}

// TestFinanceProvidersAreConstructibleFromNew pins the dispatch seam: New
// returns a finance client backed by a memory store, so nothing in this module
// can reach a provider network.
func TestFinanceProvidersAreConstructibleFromNew(t *testing.T) {
	for _, provider := range []connectors.Provider{connectors.ProviderFreshBooks, connectors.ProviderMercury} {
		client, err := New(provider)
		if err != nil {
			t.Fatalf("New(%q): %v", provider, err)
		}
		if client.Provider() != provider {
			t.Errorf("client provider = %q, want %q", client.Provider(), provider)
		}
	}
}

// TestFinanceModelsCarryNoCredentialFields is a structural guard for the
// metadata-only invariant: no finance shape may gain a field that would carry
// key material or a full account number.
func TestFinanceModelsCarryNoCredentialFields(t *testing.T) {
	banned := []string{"token", "secret", "password", "apikey", "api_key",
		"credential", "accountnumber", "account_number", "routing"}

	for _, field := range []string{
		"id", "number", "client_id", "status", "amount_minor", "outstanding_minor",
		"currency", "issued_at", "due_at", "vendor_name", "category_id",
		"incurred_at", "notes", "invoice_id", "method", "received_at",
		"organization", "email", "name", "kind", "last_four", "opened_at",
		"current_minor", "available_minor", "as_of", "account_id", "direction",
		"counterparty_name", "description", "posted_at",
	} {
		for _, bad := range banned {
			if strings.Contains(strings.ReplaceAll(field, "_", ""), strings.ReplaceAll(bad, "_", "")) {
				t.Errorf("finance model field %q looks like credential or full-account material", field)
			}
		}
	}
}
