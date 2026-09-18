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
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/HaikeiLabs/kei-connector-contracts"
)

// FreshBooks is the bookkeeping provider. Its contract surface is read-only:
// the capability catalog names invoice.read, expense.read, payment.read and
// client.read, and nothing else. Creating an invoice, recording a payment or
// altering an expense is not expressible here, so a bookkeeping connector
// cannot be talked into mutating a ledger.
//
// Money amounts are modelled as integer minor units (cents) with an explicit
// currency, never as floats: a float cannot represent 0.10 exactly, and
// silent rounding in a ledger is a correctness bug, not a display concern.
//
// Only the shapes FreshBooks returns are modelled here, per CONTRIBUTING.md —
// this package models providers, it does not call them. The backend interface
// is the seam a tenant-side runtime implements; the bundled memory backend is
// the only implementation in this module, so no network call can occur.

// FreshBooksInvoice is a bookkeeping invoice as returned by the provider.
type FreshBooksInvoice struct {
	ID               string    `json:"id"`
	Number           string    `json:"number"`
	ClientID         string    `json:"client_id"`
	Status           string    `json:"status"`
	AmountMinor      int64     `json:"amount_minor"`
	OutstandingMinor int64     `json:"outstanding_minor"`
	Currency         string    `json:"currency"`
	IssuedAt         time.Time `json:"issued_at"`
	DueAt            time.Time `json:"due_at"`
}

// FreshBooksExpense is a recorded expense.
type FreshBooksExpense struct {
	ID          string    `json:"id"`
	VendorName  string    `json:"vendor_name"`
	CategoryID  string    `json:"category_id"`
	AmountMinor int64     `json:"amount_minor"`
	Currency    string    `json:"currency"`
	IncurredAt  time.Time `json:"incurred_at"`
	Notes       string    `json:"notes,omitempty"`
}

// FreshBooksPayment is a payment applied against an invoice.
type FreshBooksPayment struct {
	ID          string    `json:"id"`
	InvoiceID   string    `json:"invoice_id"`
	AmountMinor int64     `json:"amount_minor"`
	Currency    string    `json:"currency"`
	Method      string    `json:"method"`
	ReceivedAt  time.Time `json:"received_at"`
}

// FreshBooksClient is a billing client record.
type FreshBooksClient struct {
	ID           string `json:"id"`
	Organization string `json:"organization"`
	Email        string `json:"email"`
	Currency     string `json:"currency"`
}

// InvoiceReadPayload requests invoice.read over an invoice resource.
type InvoiceReadPayload struct{}

func (InvoiceReadPayload) Capability() string { return "invoice.read" }

// ExpenseReadPayload requests expense.read over an expense resource.
type ExpenseReadPayload struct{}

func (ExpenseReadPayload) Capability() string { return "expense.read" }

// PaymentReadPayload requests payment.read over a payment resource.
type PaymentReadPayload struct{}

func (PaymentReadPayload) Capability() string { return "payment.read" }

// ClientReadPayload requests client.read over a client resource.
type ClientReadPayload struct{}

func (ClientReadPayload) Capability() string { return "client.read" }

// FreshBooksBackend is the seam for FreshBooks reads. A tenant-side runtime
// implements it against the provider's API using only `user:<resource>:read`
// OAuth scopes; requesting a `:write` scope would exceed every capability
// this contract defines.
type FreshBooksBackend interface {
	Invoice(ctx context.Context, accountID, invoiceID string) (FreshBooksInvoice, error)
	Expense(ctx context.Context, accountID, expenseID string) (FreshBooksExpense, error)
	Payment(ctx context.Context, accountID, paymentID string) (FreshBooksPayment, error)
	Client(ctx context.Context, accountID, clientID string) (FreshBooksClient, error)
}

// MemoryFreshBooks is the in-memory bookkeeping backend used until a real one
// exists. Keys are "<account_id>/<record_id>" so fixtures cannot be read
// across accounts by accident.
type MemoryFreshBooks struct {
	Invoices map[string]FreshBooksInvoice
	Expenses map[string]FreshBooksExpense
	Payments map[string]FreshBooksPayment
	Clients  map[string]FreshBooksClient
}

func (m MemoryFreshBooks) Invoice(_ context.Context, accountID, invoiceID string) (FreshBooksInvoice, error) {
	v, ok := m.Invoices[accountID+"/"+invoiceID]
	if !ok {
		return FreshBooksInvoice{}, errors.New("invoice not found")
	}
	return v, nil
}

func (m MemoryFreshBooks) Expense(_ context.Context, accountID, expenseID string) (FreshBooksExpense, error) {
	v, ok := m.Expenses[accountID+"/"+expenseID]
	if !ok {
		return FreshBooksExpense{}, errors.New("expense not found")
	}
	return v, nil
}

func (m MemoryFreshBooks) Payment(_ context.Context, accountID, paymentID string) (FreshBooksPayment, error) {
	v, ok := m.Payments[accountID+"/"+paymentID]
	if !ok {
		return FreshBooksPayment{}, errors.New("payment not found")
	}
	return v, nil
}

func (m MemoryFreshBooks) Client(_ context.Context, accountID, clientID string) (FreshBooksClient, error) {
	v, ok := m.Clients[accountID+"/"+clientID]
	if !ok {
		return FreshBooksClient{}, errors.New("client not found")
	}
	return v, nil
}

// FreshBooksClientConn serves invoice.read, expense.read, payment.read and
// client.read. Resources are accounts/<account_id>/<collection>/<record_id>,
// so every read is scoped to one bookkeeping account.
type FreshBooksClientConn struct {
	backend FreshBooksBackend
}

func NewFreshBooks(backend FreshBooksBackend) *FreshBooksClientConn {
	return &FreshBooksClientConn{backend: backend}
}

func (c *FreshBooksClientConn) Provider() connectors.Provider {
	return connectors.ProviderFreshBooks
}

func (c *FreshBooksClientConn) Invoke(ctx context.Context, meta connectors.Metadata, inv connectors.Invocation, payload Payload) (Result, error) {
	if err := checkProvider(meta, c.Provider()); err != nil {
		return Result{}, err
	}
	if err := Guard(meta, inv); err != nil {
		return Result{}, err
	}
	if err := matchCapability(inv, payload); err != nil {
		return Result{}, err
	}
	ctx = withInvocationScope(ctx, meta, inv)

	switch payload.(type) {
	case InvoiceReadPayload:
		accountID, recordID, ok := financeResource(inv.Resource, "invoices")
		if !ok {
			return Result{}, errors.New("resource must be accounts/<account_id>/invoices/<invoice_id>")
		}
		v, err := c.backend.Invoice(ctx, accountID, recordID)
		if err != nil {
			return Result{}, err
		}
		return Result{Capability: inv.Capability, Data: v}, nil
	case ExpenseReadPayload:
		accountID, recordID, ok := financeResource(inv.Resource, "expenses")
		if !ok {
			return Result{}, errors.New("resource must be accounts/<account_id>/expenses/<expense_id>")
		}
		v, err := c.backend.Expense(ctx, accountID, recordID)
		if err != nil {
			return Result{}, err
		}
		return Result{Capability: inv.Capability, Data: v}, nil
	case PaymentReadPayload:
		accountID, recordID, ok := financeResource(inv.Resource, "payments")
		if !ok {
			return Result{}, errors.New("resource must be accounts/<account_id>/payments/<payment_id>")
		}
		v, err := c.backend.Payment(ctx, accountID, recordID)
		if err != nil {
			return Result{}, err
		}
		return Result{Capability: inv.Capability, Data: v}, nil
	case ClientReadPayload:
		accountID, recordID, ok := financeResource(inv.Resource, "clients")
		if !ok {
			return Result{}, errors.New("resource must be accounts/<account_id>/clients/<client_id>")
		}
		v, err := c.backend.Client(ctx, accountID, recordID)
		if err != nil {
			return Result{}, err
		}
		return Result{Capability: inv.Capability, Data: v}, nil
	default:
		return Result{}, fmt.Errorf("unsupported payload %T", payload)
	}
}

// financeResource parses accounts/<account_id>/<collection>/<record_id> and
// reports whether it matched the expected collection. Requiring the exact
// four-segment shape keeps a read for one collection from resolving against
// another, and keeps an account id from being omitted.
func financeResource(resource, collection string) (accountID, recordID string, ok bool) {
	parts := strings.Split(resource, "/")
	if len(parts) != 4 || parts[0] != "accounts" || parts[2] != collection {
		return "", "", false
	}
	if parts[1] == "" || parts[3] == "" {
		return "", "", false
	}
	return parts[1], parts[3], true
}
