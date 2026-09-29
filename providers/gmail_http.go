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
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const gmailAPIBase = "https://gmail.googleapis.com/gmail/v1/users/me"

// AuthenticatedGmail is the production Gmail backend. It calls only GET
// endpoints on a fixed base, and asks Gmail for format=full (the only format
// that carries a body) only when the connector policy opted in.
type AuthenticatedGmail struct {
	httpClient  *http.Client
	credentials CredentialResolver
}

func NewAuthenticatedGmail(config RuntimeConfig) *GmailClient {
	config = normalizeRuntimeConfig(config)
	return NewGmail(&AuthenticatedGmail{httpClient: config.HTTPClient, credentials: config.Credentials})
}

func NewGmailHTTP(httpClient *http.Client, credentials CredentialResolver) *GmailClient {
	return NewAuthenticatedGmail(RuntimeConfig{HTTPClient: httpClient, Credentials: credentials})
}

func (b *AuthenticatedGmail) SearchMessages(ctx context.Context, query string, pageSize int, pageToken string) ([]GmailMessage, string, error) {
	params := url.Values{"maxResults": {strconv.Itoa(pageSize)}}
	if query != "" {
		params.Set("q", query)
	}
	if pageToken != "" {
		params.Set("pageToken", pageToken)
	}
	var response struct {
		Messages []struct {
			ID string `json:"id"`
		} `json:"messages"`
		NextPageToken string `json:"nextPageToken"`
	}
	if err := b.getJSON(ctx, gmailAPIBase+"/messages?"+params.Encode(), &response); err != nil {
		return nil, "", err
	}
	messages := make([]GmailMessage, 0, len(response.Messages))
	for _, ref := range response.Messages {
		message, err := b.Message(ctx, ref.ID, false)
		if err != nil {
			return nil, "", err
		}
		messages = append(messages, message)
	}
	return messages, response.NextPageToken, nil
}

func (b *AuthenticatedGmail) Message(ctx context.Context, messageID string, includeBody bool) (GmailMessage, error) {
	params := url.Values{"format": {"metadata"}, "metadataHeaders": {"From", "To", "Subject"}}
	if includeBody {
		params = url.Values{"format": {"full"}}
	}
	var response gmailMessageResponse
	if err := b.getJSON(ctx, gmailAPIBase+"/messages/"+url.PathEscape(messageID)+"?"+params.Encode(), &response); err != nil {
		return GmailMessage{}, err
	}
	message := response.model()
	if includeBody {
		message.Body = response.Payload.text()
	}
	return message, nil
}

type gmailMessageResponse struct {
	ID           string        `json:"id"`
	ThreadID     string        `json:"threadId"`
	LabelIDs     []string      `json:"labelIds"`
	Snippet      string        `json:"snippet"`
	InternalDate string        `json:"internalDate"`
	Payload      gmailMIMEPart `json:"payload"`
}

type gmailMIMEPart struct {
	MIMEType string `json:"mimeType"`
	Headers  []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"headers"`
	Body struct {
		Data string `json:"data"`
	} `json:"body"`
	Parts []gmailMIMEPart `json:"parts"`
}

func (r gmailMessageResponse) model() GmailMessage {
	message := GmailMessage{ID: r.ID, ThreadID: r.ThreadID, Labels: r.LabelIDs, Snippet: r.Snippet}
	if millis, err := strconv.ParseInt(r.InternalDate, 10, 64); err == nil {
		message.InternalDate = time.UnixMilli(millis).UTC()
	}
	for _, header := range r.Payload.Headers {
		switch strings.ToLower(header.Name) {
		case "from":
			message.From = header.Value
		case "to":
			message.To = header.Value
		case "subject":
			message.Subject = header.Value
		}
	}
	return message
}

// text returns the first text/plain part, falling back to text/html.
func (p gmailMIMEPart) text() string {
	if body := p.find("text/plain"); body != "" {
		return body
	}
	return p.find("text/html")
}

func (p gmailMIMEPart) find(mimeType string) string {
	if strings.EqualFold(p.MIMEType, mimeType) && p.Body.Data != "" {
		// Gmail uses base64url; accept both padded and unpadded data.
		if decoded, err := base64.RawURLEncoding.DecodeString(strings.TrimRight(p.Body.Data, "=")); err == nil {
			return string(decoded)
		}
	}
	for _, part := range p.Parts {
		if body := part.find(mimeType); body != "" {
			return body
		}
	}
	return ""
}

func (b *AuthenticatedGmail) getJSON(ctx context.Context, endpoint string, out any) error {
	scope, err := scopeFromContext(ctx)
	if err != nil {
		return err
	}
	token, err := b.credentials.Resolve(ctx, scope.tenantID, scope.workspaceID, scope.subject, scope.credential)
	if err != nil || token == "" {
		return errors.New("credential resolution failed")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return errors.New("provider request could not be created")
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := b.httpClient.Do(req)
	if err != nil {
		return errors.New("provider request failed")
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("provider request was rejected")
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return errors.New("provider response was invalid")
	}
	return nil
}

var _ GmailBackend = (*AuthenticatedGmail)(nil)
