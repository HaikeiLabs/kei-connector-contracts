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
	"sort"
	"strings"
	"time"

	"github.com/HaikeiLabs/kei-connector-contracts"
)

// GmailReadonlyScope is the only OAuth scope a gmail connector needs. The
// provider never requests a broader Gmail scope.
const GmailReadonlyScope = "https://www.googleapis.com/auth/gmail.readonly"

const (
	gmailDefaultPageSize = 10
	// gmailMaxPageSize bounds a search: each result costs one metadata
	// request against the caller's Gmail quota.
	gmailMaxPageSize = 25
)

// GmailMessage is the read-only message record exposed by the client. Body is
// set only by message.get on a connector whose policy sets
// gmail_include_body; every other result carries metadata and snippet only.
type GmailMessage struct {
	ID           string    `json:"message_id"`
	ThreadID     string    `json:"thread_id"`
	InternalDate time.Time `json:"internal_date"`
	From         string    `json:"from,omitempty"`
	To           string    `json:"to,omitempty"`
	Subject      string    `json:"subject,omitempty"`
	Labels       []string  `json:"labels,omitempty"`
	Snippet      string    `json:"snippet"`
	Body         string    `json:"body,omitempty"`
}

// GmailMessagePage is the message.search result.
type GmailMessagePage struct {
	Messages      []GmailMessage `json:"messages"`
	NextPageToken string         `json:"next_page_token,omitempty"`
}

// MessageSearchPayload requests message.search over gmail/messages. Query
// uses Gmail search syntax; PageSize defaults to 10 and is capped at 25.
type MessageSearchPayload struct {
	Query     string `json:"query,omitempty"`
	PageSize  int    `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
}

func (MessageSearchPayload) Capability() string { return "message.search" }

// MessageGetPayload requests message.get over gmail/messages/<message_id>.
// MessageID is optional; when set it must match the resource.
type MessageGetPayload struct {
	MessageID string `json:"message_id,omitempty"`
}

func (MessageGetPayload) Capability() string { return "message.get" }

// GmailBackend is the seam for Gmail. includeBody is true only when the
// connector policy opted in, so a backend can avoid fetching bodies at all.
type GmailBackend interface {
	SearchMessages(ctx context.Context, query string, pageSize int, pageToken string) ([]GmailMessage, string, error)
	Message(ctx context.Context, messageID string, includeBody bool) (GmailMessage, error)
}

// MemoryGmail is the in-memory Gmail backend used by tests and fixtures.
type MemoryGmail struct {
	Messages map[string]GmailMessage
}

func (m MemoryGmail) SearchMessages(_ context.Context, query string, pageSize int, _ string) ([]GmailMessage, string, error) {
	ids := make([]string, 0, len(m.Messages))
	for id := range m.Messages {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	var out []GmailMessage
	for _, id := range ids {
		msg := m.Messages[id]
		if query == "" || strings.Contains(strings.ToLower(msg.Subject+" "+msg.Snippet), strings.ToLower(query)) {
			msg.Body = ""
			out = append(out, msg)
		}
		if len(out) == pageSize {
			break
		}
	}
	return out, "", nil
}

func (m MemoryGmail) Message(_ context.Context, messageID string, includeBody bool) (GmailMessage, error) {
	msg, ok := m.Messages[messageID]
	if !ok {
		return GmailMessage{}, errors.New("message not found")
	}
	if !includeBody {
		msg.Body = ""
	}
	return msg, nil
}

// GmailClient serves message.search and message.get. Resources are
// gmail/messages for search and gmail/messages/<message_id> for get. The
// client strips bodies itself, so a backend that ignores includeBody still
// cannot leak one.
type GmailClient struct {
	backend GmailBackend
}

func NewGmail(backend GmailBackend) *GmailClient { return &GmailClient{backend: backend} }

func (c *GmailClient) Provider() connectors.Provider { return connectors.ProviderGmail }

func (c *GmailClient) Invoke(ctx context.Context, meta connectors.Metadata, inv connectors.Invocation, payload Payload) (Result, error) {
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
	switch p := payload.(type) {
	case MessageSearchPayload:
		if inv.Resource != "gmail/messages" {
			return Result{}, errors.New("resource must be gmail/messages")
		}
		pageSize := p.PageSize
		if pageSize == 0 {
			pageSize = gmailDefaultPageSize
		}
		if pageSize < 0 || pageSize > gmailMaxPageSize {
			return Result{}, fmt.Errorf("page_size must be between 1 and %d", gmailMaxPageSize)
		}
		messages, next, err := c.backend.SearchMessages(ctx, p.Query, pageSize, p.PageToken)
		if err != nil {
			return Result{}, err
		}
		for i := range messages {
			messages[i].Body = ""
		}
		return Result{Capability: inv.Capability, Data: GmailMessagePage{Messages: messages, NextPageToken: next}}, nil
	case MessageGetPayload:
		messageID, ok := gmailMessageResource(inv.Resource)
		if !ok {
			return Result{}, errors.New("resource must be gmail/messages/<message_id>")
		}
		if p.MessageID != "" && p.MessageID != messageID {
			return Result{}, errors.New("message_id does not match the resource")
		}
		includeBody := meta.Policy.GmailIncludeBody
		message, err := c.backend.Message(ctx, messageID, includeBody)
		if err != nil {
			return Result{}, err
		}
		if !includeBody {
			message.Body = ""
		}
		return Result{Capability: inv.Capability, Data: message}, nil
	default:
		return Result{}, fmt.Errorf("unsupported payload %T", payload)
	}
}

func gmailMessageResource(resource string) (string, bool) {
	parts := strings.Split(resource, "/")
	if len(parts) != 3 || parts[0] != "gmail" || parts[1] != "messages" || parts[2] == "" {
		return "", false
	}
	return parts[2], true
}
