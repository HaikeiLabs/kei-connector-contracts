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

// Mercury is the banking provider. Its contract surface is read-only:
// account.read, transaction.read and balance.read, and nothing else. There is
// no transfer, no payment and no recipient mutation, so a banking connector
// built against this contract cannot move money — the guarantee is structural,
// not a policy setting that could be misconfigured.
//
// This matches the provider's own posture: Mercury issues read-only API
// tokens distinct from write tokens, and its hosted MCP surface is documented
// as read-only and unable to initiate transactions. A tenant-side runtime
// should use the read-only token, so the contract and the credential agree.
//
// Balances and transaction amounts are integer minor units with an explicit
// currency, never floats. Rounding error in a banking balance is a
// correctness bug.
//
// Per CONTRIBUTING.md this models the provider's shapes and does not call it.

// MercuryAccount is a bank account record. It deliberately carries only the
// last four digits of the account number: the full number is not needed to
// reconcile, and a contract that never models it cannot leak it.
type MercuryAccount struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Kind     string    `json:"kind"`
	LastFour string    `json:"last_four"`
	Currency string    `json:"currency"`
	Status   string    `json:"status"`
	OpenedAt time.Time `json:"opened_at"`
}

// MercuryBalance is an account balance at a point in time.
type MercuryBalance struct {
	AccountID      string    `json:"account_id"`
	CurrentMinor   int64     `json:"current_minor"`
	AvailableMinor int64     `json:"available_minor"`
	Currency       string    `json:"currency"`
	AsOf           time.Time `json:"as_of"`
}

// MercuryTransaction is a posted or pending transaction.
type MercuryTransaction struct {
	ID               string    `json:"id"`
	AccountID        string    `json:"account_id"`
	AmountMinor      int64     `json:"amount_minor"`
	Currency         string    `json:"currency"`
	Direction        string    `json:"direction"`
	Status           string    `json:"status"`
	CounterpartyName string    `json:"counterparty_name,omitempty"`
	Description      string    `json:"description,omitempty"`
	PostedAt         time.Time `json:"posted_at"`
}

// AccountReadPayload requests account.read over an account resource.
type AccountReadPayload struct{}

func (AccountReadPayload) Capability() string { return "account.read" }

// BalanceReadPayload requests balance.read over an account resource.
type BalanceReadPayload struct{}

func (BalanceReadPayload) Capability() string { return "balance.read" }

// TransactionReadPayload requests transaction.read over a transaction
// resource.
type TransactionReadPayload struct{}

func (TransactionReadPayload) Capability() string { return "transaction.read" }

// MercuryBackend is the seam for Mercury reads. A tenant-side runtime
// implements it against the provider's read-only token or hosted MCP surface.
type MercuryBackend interface {
	Account(ctx context.Context, accountID string) (MercuryAccount, error)
	Balance(ctx context.Context, accountID string) (MercuryBalance, error)
	Transaction(ctx context.Context, accountID, transactionID string) (MercuryTransaction, error)
}

// MemoryMercury is the in-memory banking backend used until a real one
// exists. Transaction keys are "<account_id>/<transaction_id>" so a
// transaction cannot be read outside its account.
type MemoryMercury struct {
	Accounts     map[string]MercuryAccount
	Balances     map[string]MercuryBalance
	Transactions map[string]MercuryTransaction
}

func (m MemoryMercury) Account(_ context.Context, accountID string) (MercuryAccount, error) {
	v, ok := m.Accounts[accountID]
	if !ok {
		return MercuryAccount{}, errors.New("account not found")
	}
	return v, nil
}

func (m MemoryMercury) Balance(_ context.Context, accountID string) (MercuryBalance, error) {
	v, ok := m.Balances[accountID]
	if !ok {
		return MercuryBalance{}, errors.New("balance not found")
	}
	return v, nil
}

func (m MemoryMercury) Transaction(_ context.Context, accountID, transactionID string) (MercuryTransaction, error) {
	v, ok := m.Transactions[accountID+"/"+transactionID]
	if !ok {
		return MercuryTransaction{}, errors.New("transaction not found")
	}
	return v, nil
}

// MercuryClientConn serves account.read, balance.read and transaction.read.
// Resources are accounts/<account_id> for account and balance reads, and
// accounts/<account_id>/transactions/<transaction_id> for transaction reads.
type MercuryClientConn struct {
	backend MercuryBackend
}

func NewMercury(backend MercuryBackend) *MercuryClientConn {
	return &MercuryClientConn{backend: backend}
}

func (c *MercuryClientConn) Provider() connectors.Provider {
	return connectors.ProviderMercury
}

func (c *MercuryClientConn) Invoke(ctx context.Context, meta connectors.Metadata, inv connectors.Invocation, payload Payload) (Result, error) {
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
	case AccountReadPayload:
		accountID, ok := bankAccountResource(inv.Resource)
		if !ok {
			return Result{}, errors.New("resource must be accounts/<account_id>")
		}
		v, err := c.backend.Account(ctx, accountID)
		if err != nil {
			return Result{}, err
		}
		return Result{Capability: inv.Capability, Data: v}, nil
	case BalanceReadPayload:
		accountID, ok := bankAccountResource(inv.Resource)
		if !ok {
			return Result{}, errors.New("resource must be accounts/<account_id>")
		}
		v, err := c.backend.Balance(ctx, accountID)
		if err != nil {
			return Result{}, err
		}
		return Result{Capability: inv.Capability, Data: v}, nil
	case TransactionReadPayload:
		accountID, transactionID, ok := financeResource(inv.Resource, "transactions")
		if !ok {
			return Result{}, errors.New("resource must be accounts/<account_id>/transactions/<transaction_id>")
		}
		v, err := c.backend.Transaction(ctx, accountID, transactionID)
		if err != nil {
			return Result{}, err
		}
		return Result{Capability: inv.Capability, Data: v}, nil
	default:
		return Result{}, fmt.Errorf("unsupported payload %T", payload)
	}
}

// bankAccountResource parses accounts/<account_id>.
func bankAccountResource(resource string) (string, bool) {
	parts := strings.Split(resource, "/")
	if len(parts) != 2 || parts[0] != "accounts" || parts[1] == "" {
		return "", false
	}
	return parts[1], true
}
