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
	"testing"

	"github.com/HaikeiLabs/kei-connector-contracts/contract"
)

func TestConnectorContractFreezeResourceAndResultShapes(t *testing.T) {
	ctx := context.Background()
	github := NewGitHub(githubStore())
	gm := metaFor(t, contract.ProviderGitHub, []string{"repos/acme/kei"}, nil)
	githubCases := []struct {
		resource string
		payload  Payload
	}{
		{"repos/acme/kei", RepositoryReadPayload{}},
		{"repos/acme/kei/issues/7", GitHubIssueReadPayload{}},
		{"repos/acme/kei/pulls/8", PullRequestReadPayload{}},
		{"repos/acme/kei/checks/chk-1", CheckReadPayload{}},
		{"repos/acme/kei/workflows/wf-1", WorkflowReadPayload{}},
	}
	for _, tc := range githubCases {
		inv := invocation(tc.payload.Capability(), contract.ActionRead, tc.resource)
		if _, err := github.Invoke(ctx, gm, inv, tc.payload); err != nil {
			t.Errorf("GitHub resource %q rejected: %v", tc.resource, err)
		}
	}

	google := NewGoogle(googleStore())
	drive := metaFor(t, contract.ProviderGoogle, []string{"drive/d-1"}, nil)
	for _, tc := range []struct {
		resource string
		payload  Payload
	}{
		{"drive/d-1", DriveSearchPayload{}},
		{"drive/d-1/files/f-1", DriveMetadataPayload{}},
		{"drive/d-1/files/f-1", DocsReadPayload{}},
	} {
		inv := invocation(tc.payload.Capability(), contract.ActionRead, tc.resource)
		if _, err := google.Invoke(ctx, drive, inv, tc.payload); err != nil {
			t.Errorf("Google Drive resource %q rejected: %v", tc.resource, err)
		}
	}

	linear := NewLinear(linearStore())
	lm := metaFor(t, contract.ProviderLinear, []string{"linear/team/KEI"}, nil)
	for _, tc := range []struct {
		resource string
		payload  Payload
	}{
		{"linear/team/KEI", TeamReadPayload{}},
		{"linear/project/p-1", ProjectReadPayload{}},
		{"linear/cycle/cy-1", CycleReadPayload{}},
		{"linear/issue/i-1", LinearIssueReadPayload{}},
	} {
		inv := invocation(tc.payload.Capability(), contract.ActionRead, tc.resource)
		if _, err := linear.Invoke(ctx, lm, inv, tc.payload); err != nil {
			t.Errorf("Linear resource %q rejected: %v", tc.resource, err)
		}
	}

	notion := NewNotion(notionStore())
	nm := metaFor(t, contract.ProviderNotion, []string{"pages/p-1", "databases/d-1"}, nil)
	notionCases := []struct {
		resource string
		payload  Payload
	}{
		{"pages/p-1", PageReadPayload{}},
		{"databases/d-1", DatabaseQueryPayload{}},
		{"search", SearchPayload{Query: "test"}},
	}
	for _, tc := range notionCases {
		inv := invocation(tc.payload.Capability(), contract.ActionRead, tc.resource)
		if _, err := notion.Invoke(ctx, nm, inv, tc.payload); err != nil {
			t.Errorf("Notion resource %q rejected: %v", tc.resource, err)
		}
	}

	freshbooks := NewFreshBooks(freshBooksStore())
	fm := metaFor(t, contract.ProviderFreshBooks, []string{"accounts/acct-TEST-1"}, nil)
	for _, tc := range []struct {
		resource string
		payload  Payload
	}{
		{"accounts/acct-TEST-1/invoices/inv-TEST-1", InvoiceReadPayload{}},
		{"accounts/acct-TEST-1/expenses/exp-TEST-1", ExpenseReadPayload{}},
		{"accounts/acct-TEST-1/payments/pay-TEST-1", PaymentReadPayload{}},
		{"accounts/acct-TEST-1/clients/cli-TEST-1", ClientReadPayload{}},
	} {
		inv := invocation(tc.payload.Capability(), contract.ActionRead, tc.resource)
		if _, err := freshbooks.Invoke(ctx, fm, inv, tc.payload); err != nil {
			t.Errorf("FreshBooks resource %q rejected: %v", tc.resource, err)
		}
	}

	mercury := NewMercury(mercuryStore())
	mm := metaFor(t, contract.ProviderMercury, []string{"accounts/acct-TEST-1"}, nil)
	for _, tc := range []struct {
		resource string
		payload  Payload
	}{
		{"accounts/acct-TEST-1", AccountReadPayload{}},
		{"accounts/acct-TEST-1", BalanceReadPayload{}},
		{"accounts/acct-TEST-1/transactions/txn-TEST-1", TransactionReadPayload{}},
	} {
		inv := invocation(tc.payload.Capability(), contract.ActionRead, tc.resource)
		if _, err := mercury.Invoke(ctx, mm, inv, tc.payload); err != nil {
			t.Errorf("Mercury resource %q rejected: %v", tc.resource, err)
		}
	}

	// HAI-210: investor resources are investors (list) and investors/<id>
	// (read).
	crm := NewCRM(crmStore())
	cm := metaFor(t, contract.ProviderCRM, []string{"investors"}, nil)
	for _, tc := range []struct {
		resource string
		payload  Payload
	}{
		{"investors", InvestorListPayload{}},
		{"investors/i-1", InvestorReadPayload{}},
	} {
		inv := invocation(tc.payload.Capability(), contract.ActionRead, tc.resource)
		if _, err := crm.Invoke(ctx, cm, inv, tc.payload); err != nil {
			t.Errorf("CRM resource %q rejected: %v", tc.resource, err)
		}
	}

	// HAI-202: Gmail resources are gmail/messages (search) and
	// gmail/messages/<message_id> (get).
	gmail := NewGmail(gmailStore())
	gmm := metaFor(t, contract.ProviderGmail, []string{"gmail/messages"}, nil)
	for _, tc := range []struct {
		resource string
		payload  Payload
	}{
		{"gmail/messages", MessageSearchPayload{}},
		{"gmail/messages/m-1", MessageGetPayload{}},
	} {
		inv := invocation(tc.payload.Capability(), contract.ActionRead, tc.resource)
		if _, err := gmail.Invoke(ctx, gmm, inv, tc.payload); err != nil {
			t.Errorf("Gmail resource %q rejected: %v", tc.resource, err)
		}
	}

	// HAI-203: Tito resources are tito/<account>/events and the event,
	// releases, and ticket-summary resources beneath one event.
	tito := NewTito(titoStore())
	tm := metaFor(t, contract.ProviderTito, []string{"tito/acme/events"}, nil)
	for _, tc := range []struct {
		resource string
		payload  Payload
	}{
		{"tito/acme/events", EventListPayload{}},
		{"tito/acme/events/conf-2026", EventGetPayload{}},
		{"tito/acme/events/conf-2026/releases", ReleaseListPayload{}},
		{"tito/acme/events/conf-2026/ticket-summary", TicketSummaryPayload{}},
	} {
		inv := invocation(tc.payload.Capability(), contract.ActionRead, tc.resource)
		if _, err := tito.Invoke(ctx, tm, inv, tc.payload); err != nil {
			t.Errorf("Tito resource %q rejected: %v", tc.resource, err)
		}
	}
}
