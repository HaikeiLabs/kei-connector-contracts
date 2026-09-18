package providers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/HaikeiLabs/kei-connector-contracts"
)

type NotionPage struct {
	ID         string         `json:"id"`
	Title      string         `json:"title"`
	ParentID   string         `json:"parent_id"`
	ParentType string         `json:"parent_type"`
	URL        string         `json:"url"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	Archived   bool           `json:"archived"`
	Properties map[string]any `json:"properties,omitempty"`
}

type NotionDatabase struct {
	ID         string         `json:"id"`
	Title      string         `json:"title"`
	URL        string         `json:"url"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	Archived   bool           `json:"archived"`
	Properties map[string]any `json:"properties,omitempty"`
}

type NotionSearchResult struct {
	Results    []any  `json:"results"`
	HasMore    bool   `json:"has_more"`
	NextCursor string `json:"next_cursor,omitempty"`
}

// PageReadPayload requests page.read over a page resource.
type PageReadPayload struct{}

func (PageReadPayload) Capability() string { return "page.read" }

// DatabaseQueryPayload requests database.query over a database resource.
type DatabaseQueryPayload struct {
	Filter   map[string]any   `json:"filter,omitempty"`
	Sorts    []map[string]any `json:"sorts,omitempty"`
	PageSize int              `json:"page_size,omitempty"`
}

func (DatabaseQueryPayload) Capability() string { return "database.query" }

// SearchPayload requests search across all accessible pages and databases.
type SearchPayload struct {
	Query    string `json:"query"`
	PageSize int    `json:"page_size,omitempty"`
}

func (SearchPayload) Capability() string { return "search" }

// NotionBackend is the seam for Notion reads.
type NotionBackend interface {
	Page(ctx context.Context, pageID string) (NotionPage, error)
	QueryDatabase(ctx context.Context, databaseID string, filter map[string]any, sorts []map[string]any, pageSize int) ([]NotionPage, error)
	Search(ctx context.Context, query string, pageSize int) (NotionSearchResult, error)
}

// MemoryNotion is the in-memory Notion backend used until a real one exists.
type MemoryNotion struct {
	Pages     map[string]NotionPage
	Databases map[string]NotionDatabase
}

func (m MemoryNotion) Page(_ context.Context, pageID string) (NotionPage, error) {
	p, ok := m.Pages[pageID]
	if !ok {
		return NotionPage{}, errors.New("page not found")
	}
	return p, nil
}

func (m MemoryNotion) QueryDatabase(_ context.Context, databaseID string, filter map[string]any, sorts []map[string]any, pageSize int) ([]NotionPage, error) {
	db, ok := m.Databases[databaseID]
	if !ok {
		return nil, errors.New("database not found")
	}
	_ = db
	var out []NotionPage
	for _, p := range m.Pages {
		if p.ParentType == "database_id" && p.ParentID == databaseID {
			out = append(out, p)
		}
	}
	if pageSize > 0 && len(out) > pageSize {
		out = out[:pageSize]
	}
	return out, nil
}

func (m MemoryNotion) Search(_ context.Context, query string, pageSize int) (NotionSearchResult, error) {
	var results []any
	lowerQuery := strings.ToLower(query)
	for _, p := range m.Pages {
		if query == "" || strings.Contains(strings.ToLower(p.Title), lowerQuery) {
			results = append(results, p)
		}
	}
	for _, d := range m.Databases {
		if query == "" || strings.Contains(strings.ToLower(d.Title), lowerQuery) {
			results = append(results, d)
		}
	}
	if pageSize > 0 && len(results) > pageSize {
		results = results[:pageSize]
	}
	return NotionSearchResult{Results: results, HasMore: false}, nil
}

// NotionClient serves search, page.read, and database.query.
// Resources are pages/<page_id> for page reads and databases/<database_id> for database queries.
type NotionClient struct {
	backend NotionBackend
}

func NewNotion(backend NotionBackend) *NotionClient { return &NotionClient{backend: backend} }

func (c *NotionClient) Provider() connectors.Provider { return connectors.ProviderNotion }

func (c *NotionClient) Invoke(ctx context.Context, meta connectors.Metadata, inv connectors.Invocation, payload Payload) (Result, error) {
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
	case PageReadPayload:
		pageID, ok := pageResource(inv.Resource)
		if !ok {
			return Result{}, errors.New("resource must be pages/<page_id>")
		}
		page, err := c.backend.Page(ctx, pageID)
		if err != nil {
			return Result{}, err
		}
		return Result{Capability: inv.Capability, Data: page}, nil
	case DatabaseQueryPayload:
		databaseID, ok := databaseResource(inv.Resource)
		if !ok {
			return Result{}, errors.New("resource must be databases/<database_id>")
		}
		pages, err := c.backend.QueryDatabase(ctx, databaseID, p.Filter, p.Sorts, p.PageSize)
		if err != nil {
			return Result{}, err
		}
		return Result{Capability: inv.Capability, Data: pages}, nil
	case SearchPayload:
		// Search is a top-level operation; the resource is "search" or "search/<query>".
		query := p.Query
		if inv.Resource != "search" {
			parts := strings.SplitN(inv.Resource, "/", 2)
			if len(parts) == 2 && parts[0] == "search" && parts[1] != "" {
				query = parts[1]
			} else if inv.Resource != "search" {
				return Result{}, errors.New("resource must be search or search/<query>")
			}
		}
		result, err := c.backend.Search(ctx, query, p.PageSize)
		if err != nil {
			return Result{}, err
		}
		return Result{Capability: inv.Capability, Data: result}, nil
	default:
		return Result{}, fmt.Errorf("unsupported payload %T", payload)
	}
}

func pageResource(resource string) (string, bool) {
	parts := strings.Split(resource, "/")
	if len(parts) != 2 || parts[0] != "pages" || parts[1] == "" {
		return "", false
	}
	return parts[1], true
}

func databaseResource(resource string) (string, bool) {
	parts := strings.Split(resource, "/")
	if len(parts) != 2 || parts[0] != "databases" || parts[1] == "" {
		return "", false
	}
	return parts[1], true
}
