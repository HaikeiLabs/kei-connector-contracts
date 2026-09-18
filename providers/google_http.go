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
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	googleDriveAPIBase = "https://www.googleapis.com/drive/v3"
	googleDocsAPIBase  = "https://docs.googleapis.com/v1"
)

// AuthenticatedGoogle is the production Drive/Docs backend. Drive and Docs
// endpoints are fixed constants; connector resources supply identifiers only.
type AuthenticatedGoogle struct {
	httpClient  *http.Client
	credentials CredentialResolver
}

func NewAuthenticatedGoogle(config RuntimeConfig) *GoogleClient {
	config = normalizeRuntimeConfig(config)
	return NewGoogle(&AuthenticatedGoogle{httpClient: config.HTTPClient, credentials: config.Credentials})
}

func NewGoogleHTTP(httpClient *http.Client, credentials CredentialResolver) *GoogleClient {
	return NewAuthenticatedGoogle(RuntimeConfig{HTTPClient: httpClient, Credentials: credentials})
}

func (b *AuthenticatedGoogle) SearchFiles(ctx context.Context, driveID, query string) ([]GoogleFile, error) {
	params := url.Values{}
	queryParts := []string{"trashed = false"}
	if query != "" {
		queryParts = append(queryParts, "name contains '"+escapeDriveQueryValue(query)+"'")
	}
	params.Set("q", strings.Join(queryParts, " and "))
	params.Set("corpora", "drive")
	params.Set("driveId", driveID)
	params.Set("includeItemsFromAllDrives", "true")
	params.Set("supportsAllDrives", "true")
	params.Set("fields", "files(id,driveId,name,mimeType,modifiedTime)")
	var response struct {
		Files []googleFileResponse `json:"files"`
	}
	if err := b.getJSON(ctx, googleDriveAPIBase+"/files?"+params.Encode(), &response); err != nil {
		return nil, err
	}
	files := make([]GoogleFile, 0, len(response.Files))
	for _, file := range response.Files {
		files = append(files, file.model(driveID))
	}
	return files, nil
}

func (b *AuthenticatedGoogle) File(ctx context.Context, driveID, fileID string) (GoogleFile, error) {
	params := url.Values{"supportsAllDrives": {"true"}, "fields": {"id,driveId,name,mimeType,modifiedTime"}}
	var response googleFileResponse
	if err := b.getJSON(ctx, googleDriveAPIBase+"/files/"+url.PathEscape(fileID)+"?"+params.Encode(), &response); err != nil {
		return GoogleFile{}, err
	}
	if response.DriveID != "" && response.DriveID != driveID {
		return GoogleFile{}, errors.New("provider resource was rejected")
	}
	return response.model(driveID), nil
}

func (b *AuthenticatedGoogle) Document(ctx context.Context, driveID, fileID string) (string, error) {
	file, err := b.File(ctx, driveID, fileID)
	if err != nil {
		return "", err
	}
	if file.MimeType != "application/vnd.google-apps.document" && !strings.HasPrefix(file.MimeType, "text/") {
		return "", errors.New("file is not a readable document")
	}
	if file.MimeType != "application/vnd.google-apps.document" {
		return b.getText(ctx, googleDriveAPIBase+"/files/"+url.PathEscape(fileID)+"?alt=media&supportsAllDrives=true")
	}
	var document struct {
		Body struct {
			Content []struct {
				Paragraph *struct {
					Elements []struct {
						TextRun *struct {
							Content string `json:"content"`
						} `json:"textRun"`
					} `json:"elements"`
				} `json:"paragraph"`
			} `json:"content"`
		} `json:"body"`
	}
	if err := b.getJSON(ctx, googleDocsAPIBase+"/documents/"+url.PathEscape(fileID), &document); err != nil {
		return "", err
	}
	var content strings.Builder
	for _, block := range document.Body.Content {
		if block.Paragraph == nil {
			continue
		}
		for _, element := range block.Paragraph.Elements {
			if element.TextRun != nil {
				content.WriteString(element.TextRun.Content)
			}
		}
	}
	return content.String(), nil
}

type googleFileResponse struct {
	ID         string    `json:"id"`
	DriveID    string    `json:"driveId"`
	Name       string    `json:"name"`
	MimeType   string    `json:"mimeType"`
	ModifiedAt time.Time `json:"modifiedTime"`
}

func (f googleFileResponse) model(driveID string) GoogleFile {
	if f.DriveID == "" {
		f.DriveID = driveID
	}
	return GoogleFile{ID: f.ID, DriveID: f.DriveID, Name: f.Name, MimeType: f.MimeType, ModifiedAt: f.ModifiedAt}
}
func escapeDriveQueryValue(value string) string { return strings.ReplaceAll(value, "'", "\\'") }

func (b *AuthenticatedGoogle) getJSON(ctx context.Context, endpoint string, out any) error {
	req, err := b.authorizedRequest(ctx, http.MethodGet, endpoint)
	if err != nil {
		return err
	}
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

func (b *AuthenticatedGoogle) getText(ctx context.Context, endpoint string) (string, error) {
	req, err := b.authorizedRequest(ctx, http.MethodGet, endpoint)
	if err != nil {
		return "", err
	}
	resp, err := b.httpClient.Do(req)
	if err != nil {
		return "", errors.New("provider request failed")
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", errors.New("provider request was rejected")
	}
	var content strings.Builder
	if _, err := io.Copy(&content, resp.Body); err != nil {
		return "", errors.New("provider response was invalid")
	}
	return content.String(), nil
}

func (b *AuthenticatedGoogle) authorizedRequest(ctx context.Context, method, endpoint string) (*http.Request, error) {
	scope, err := scopeFromContext(ctx)
	if err != nil {
		return nil, err
	}
	token, err := b.credentials.Resolve(ctx, scope.tenantID, scope.workspaceID, scope.subject, scope.credential)
	if err != nil || token == "" {
		return nil, errors.New("credential resolution failed")
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint, nil)
	if err != nil {
		return nil, errors.New("provider request could not be created")
	}
	req.Header.Set("Authorization", "Bearer "+token)
	return req, nil
}

var _ GoogleBackend = (*AuthenticatedGoogle)(nil)
