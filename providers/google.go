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
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/HaikeiLabs/kei-connector-contracts/contract"
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

// DriveDocCreatePayload requests drive.doc.create on a Drive resource. The
// resource is the parent drive: drive/<drive_id>.
type DriveDocCreatePayload struct {
	Title string `json:"title"`
	Body  string `json:"body,omitempty"`
}

func (DriveDocCreatePayload) Capability() string { return "drive.doc.create" }

func (p DriveDocCreatePayload) Validate() error {
	if strings.TrimSpace(p.Title) == "" {
		return errors.New("title is required")
	}
	if len(p.Body) > maxDocBodyBytes {
		return errors.New("body exceeds maximum size")
	}
	return nil
}

// DriveDocUpdatePayload requests drive.doc.update on a file resource. The
// resource is the target document: drive/<drive_id>/files/<file_id>.
type DriveDocUpdatePayload struct {
	Title *string `json:"title,omitempty"`
	Body  *string `json:"body,omitempty"`
}

func (DriveDocUpdatePayload) Capability() string { return "drive.doc.update" }

func (p DriveDocUpdatePayload) Validate() error {
	if p.Title == nil && p.Body == nil {
		return errors.New("at least one of title or body must be set")
	}
	if p.Body != nil && len(*p.Body) > maxDocBodyBytes {
		return errors.New("body exceeds maximum size")
	}
	if p.Title != nil && strings.TrimSpace(*p.Title) == "" {
		return errors.New("title must not be empty")
	}
	return nil
}

// DriveDocWriteBounds keep writes inside the runtime's input cap.
const maxDocBodyBytes = 12 << 10

// DriveDocCreated is the drive.doc.create result.
type DriveDocCreated struct {
	ID      string `json:"id"`
	DriveID string `json:"drive_id"`
	Name    string `json:"name"`
	URL     string `json:"url"`
}

// DriveDocUpdated is the drive.doc.update result.
type DriveDocUpdated struct {
	ID      string `json:"id"`
	DriveID string `json:"drive_id"`
	Name    string `json:"name"`
	URL     string `json:"url"`
}

// GoogleBackend is the seam for Google Drive/Docs reads. Implementations
// receive opaque drive and file identifiers only.
type GoogleBackend interface {
	SearchFiles(ctx context.Context, driveID, query string) ([]GoogleFile, error)
	File(ctx context.Context, driveID, fileID string) (GoogleFile, error)
	Document(ctx context.Context, driveID, fileID string) (string, error)
}

// GoogleWriteBackend is the seam for Google Drive/Docs writes.
type GoogleWriteBackend interface {
	CreateDocument(ctx context.Context, clientID string, driveID string, in DriveDocCreatePayload) (DriveDocCreated, error)
	UpdateDocument(ctx context.Context, clientID string, driveID, fileID string, in DriveDocUpdatePayload) (DriveDocUpdated, error)
}

// GoogleWriteID derives a deterministic UUID for one write from the
// invocation identifiers. A replayed write with the same idempotency key
// produces the same client ID, so it cannot create a duplicate resource.
func GoogleWriteID(inv contract.Invocation) string {
	sum := sha256.Sum256([]byte(strings.Join([]string{inv.TenantID, inv.WorkspaceID, inv.ConnectorID, inv.Subject, inv.Capability, inv.IdempotencyKey}, "\x00")))
	b := sum[:16]
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	h := hex.EncodeToString(b)
	return h[0:8] + "-" + h[8:12] + "-" + h[12:16] + "-" + h[16:20] + "-" + h[20:32]
}

// MemoryGoogle is the in-memory Google backend used until a real one exists.
type MemoryGoogle struct {
	Files   map[string]GoogleFile
	Creates map[string]DriveDocCreated
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

func (m MemoryGoogle) CreateDocument(_ context.Context, clientID, driveID string, in DriveDocCreatePayload) (DriveDocCreated, error) {
	if m.Creates == nil {
		return DriveDocCreated{}, errors.New("drive writes are unavailable")
	}
	if existing, ok := m.Creates[clientID]; ok {
		return existing, nil
	}
	doc := DriveDocCreated{
		ID:      clientID,
		DriveID: driveID,
		Name:    in.Title,
		URL:     "https://docs.google.com/document/d/" + clientID,
	}
	if m.Creates == nil {
		m.Creates = map[string]DriveDocCreated{}
	}
	m.Creates[clientID] = doc
	return doc, nil
}

func (m MemoryGoogle) UpdateDocument(_ context.Context, clientID, driveID, fileID string, in DriveDocUpdatePayload) (DriveDocUpdated, error) {
	if m.Creates == nil {
		return DriveDocUpdated{}, errors.New("drive writes are unavailable")
	}
	if existing, ok := m.Creates[clientID]; ok {
		return DriveDocUpdated{ID: existing.ID, DriveID: existing.DriveID, Name: existing.Name, URL: existing.URL}, nil
	}
	doc := DriveDocUpdated{
		ID:      fileID,
		DriveID: driveID,
		URL:     "https://docs.google.com/document/d/" + fileID,
	}
	if in.Title != nil {
		doc.Name = *in.Title
	}
	m.Creates[clientID] = DriveDocCreated{ID: doc.ID, DriveID: doc.DriveID, Name: doc.Name, URL: doc.URL}
	return doc, nil
}

// GoogleClient serves drive.search, drive.metadata.read, docs.read,
// drive.doc.create, and drive.doc.update. Resources are drive/<drive_id> for
// search and create, and drive/<drive_id>/files/<file_id> for file reads and
// update.
type GoogleClient struct {
	backend GoogleBackend
}

func NewGoogle(backend GoogleBackend) *GoogleClient { return &GoogleClient{backend: backend} }

func (c *GoogleClient) Provider() contract.Provider { return contract.ProviderGoogle }

func (c *GoogleClient) Invoke(ctx context.Context, meta contract.Metadata, inv contract.Invocation, payload Payload) (Result, error) {
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
	case DriveDocCreatePayload:
		return c.invokeDocCreate(ctx, inv, p)
	case DriveDocUpdatePayload:
		return c.invokeDocUpdate(ctx, inv, p)
	default:
		return Result{}, fmt.Errorf("unsupported payload %T", payload)
	}
}

func (c *GoogleClient) invokeDocCreate(ctx context.Context, inv contract.Invocation, p DriveDocCreatePayload) (Result, error) {
	writer, ok := c.backend.(GoogleWriteBackend)
	if !ok {
		return Result{}, errors.New("drive writes are unavailable")
	}
	if inv.IdempotencyKey == "" {
		return Result{}, errors.New("drive writes require an idempotency key")
	}
	if err := p.Validate(); err != nil {
		return Result{}, err
	}
	driveID, ok := driveResource(inv.Resource)
	if !ok {
		return Result{}, errors.New("resource must be drive/<drive_id>")
	}
	clientID := GoogleWriteID(inv)
	data, err := writer.CreateDocument(ctx, clientID, driveID, p)
	if err != nil {
		return Result{}, err
	}
	return Result{Capability: inv.Capability, Data: data}, nil
}

func (c *GoogleClient) invokeDocUpdate(ctx context.Context, inv contract.Invocation, p DriveDocUpdatePayload) (Result, error) {
	writer, ok := c.backend.(GoogleWriteBackend)
	if !ok {
		return Result{}, errors.New("drive writes are unavailable")
	}
	if inv.IdempotencyKey == "" {
		return Result{}, errors.New("drive writes require an idempotency key")
	}
	if err := p.Validate(); err != nil {
		return Result{}, err
	}
	driveID, fileID, ok := fileResource(inv.Resource)
	if !ok {
		return Result{}, errors.New("resource must be drive/<drive_id>/files/<file_id>")
	}
	clientID := GoogleWriteID(inv)
	data, err := writer.UpdateDocument(ctx, clientID, driveID, fileID, p)
	if err != nil {
		return Result{}, err
	}
	return Result{Capability: inv.Capability, Data: data}, nil
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
