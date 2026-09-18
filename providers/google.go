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

// GoogleFile is the read-only Drive/Docs record exposed by the client.
type GoogleFile struct {
	ID         string    `json:"id"`
	DriveID    string    `json:"drive_id"`
	Name       string    `json:"name"`
	MimeType   string    `json:"mime_type"`
	ModifiedAt time.Time `json:"modified_at"`
	Content    string    `json:"-"`
}

// DriveSearchPayload requests drive.search over a drive resource.
type DriveSearchPayload struct {
	Query string `json:"query"`
}

func (DriveSearchPayload) Capability() string { return "drive.search" }

// DriveMetadataPayload requests drive.metadata.read over a file resource.
type DriveMetadataPayload struct{}

func (DriveMetadataPayload) Capability() string { return "drive.metadata.read" }

// DocsReadPayload requests docs.read over a file resource.
type DocsReadPayload struct{}

func (DocsReadPayload) Capability() string { return "docs.read" }

// GoogleBackend is the seam for Google Drive/Docs. Implementations receive
// opaque drive and file identifiers only.
type GoogleBackend interface {
	SearchFiles(ctx context.Context, driveID, query string) ([]GoogleFile, error)
	File(ctx context.Context, driveID, fileID string) (GoogleFile, error)
	Document(ctx context.Context, driveID, fileID string) (string, error)
}

// MemoryGoogle is the in-memory Google backend used until a real one exists.
type MemoryGoogle struct {
	Files map[string]GoogleFile
}

func (m MemoryGoogle) SearchFiles(_ context.Context, driveID, query string) ([]GoogleFile, error) {
	var out []GoogleFile
	for _, f := range m.Files {
		if f.DriveID != driveID {
			continue
		}
		if query == "" || strings.Contains(strings.ToLower(f.Name), strings.ToLower(query)) {
			out = append(out, f)
		}
	}
	return out, nil
}

func (m MemoryGoogle) File(_ context.Context, driveID, fileID string) (GoogleFile, error) {
	f, ok := m.Files[fileID]
	if !ok || f.DriveID != driveID {
		return GoogleFile{}, errors.New("file not found")
	}
	return f, nil
}

func (m MemoryGoogle) Document(ctx context.Context, driveID, fileID string) (string, error) {
	f, err := m.File(ctx, driveID, fileID)
	if err != nil {
		return "", err
	}
	if f.MimeType != "application/vnd.google-apps.document" && !strings.HasPrefix(f.MimeType, "text/") {
		return "", errors.New("file is not a readable document")
	}
	return f.Content, nil
}

// GoogleClient serves drive.search, drive.metadata.read, and docs.read.
// Resources are drive/<drive_id> for search and
// drive/<drive_id>/files/<file_id> for file reads.
type GoogleClient struct {
	backend GoogleBackend
}

func NewGoogle(backend GoogleBackend) *GoogleClient { return &GoogleClient{backend: backend} }

func (c *GoogleClient) Provider() connectors.Provider { return connectors.ProviderGoogle }

func (c *GoogleClient) Invoke(ctx context.Context, meta connectors.Metadata, inv connectors.Invocation, payload Payload) (Result, error) {
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
	case DriveSearchPayload:
		driveID, ok := driveResource(inv.Resource)
		if !ok {
			return Result{}, errors.New("resource must be drive/<drive_id>")
		}
		files, err := c.backend.SearchFiles(ctx, driveID, p.Query)
		if err != nil {
			return Result{}, err
		}
		return Result{Capability: inv.Capability, Data: files}, nil
	case DriveMetadataPayload:
		driveID, fileID, ok := fileResource(inv.Resource)
		if !ok {
			return Result{}, errors.New("resource must be drive/<drive_id>/files/<file_id>")
		}
		f, err := c.backend.File(ctx, driveID, fileID)
		if err != nil {
			return Result{}, err
		}
		return Result{Capability: inv.Capability, Data: f}, nil
	case DocsReadPayload:
		driveID, fileID, ok := fileResource(inv.Resource)
		if !ok {
			return Result{}, errors.New("resource must be drive/<drive_id>/files/<file_id>")
		}
		content, err := c.backend.Document(ctx, driveID, fileID)
		if err != nil {
			return Result{}, err
		}
		return Result{Capability: inv.Capability, Data: content}, nil
	default:
		return Result{}, fmt.Errorf("unsupported payload %T", payload)
	}
}

func driveResource(resource string) (string, bool) {
	parts := strings.Split(resource, "/")
	if len(parts) != 2 || parts[0] != "drive" || parts[1] == "" {
		return "", false
	}
	return parts[1], true
}

func fileResource(resource string) (string, string, bool) {
	parts := strings.Split(resource, "/")
	if len(parts) != 4 || parts[0] != "drive" || parts[2] != "files" || parts[1] == "" || parts[3] == "" {
		return "", "", false
	}
	return parts[1], parts[3], true
}
