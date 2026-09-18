package providers

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/HaikeiLabs/kei-connector-contracts"
)

func metaFor(t *testing.T, provider connectors.Provider, resources, prefixes []string, caps ...connectors.Capability) connectors.Metadata {
	t.Helper()
	if len(caps) == 0 {
		caps = connectors.CapabilitiesFor(provider)
	}
	m := connectors.Metadata{
		ID:            "c-1",
		TenantID:      "t-1",
		WorkspaceID:   "w-1",
		Name:          "test",
		Provider:      provider,
		Status:        connectors.StatusActive,
		CredentialRef: "vault/tenant/t-1/test",
		Scopes:        []string{"read"},
		Resources:     resources,
		Policy:        connectors.PolicyAttributes{AllowedActions: []connectors.Action{connectors.ActionRead}, AllowedResources: resources, AllowedPrefixes: prefixes},
		Capabilities:  caps,
		CreatedBy:     "u-1",
	}
	if err := m.Validate(); err != nil {
		t.Fatalf("fixture metadata is invalid: %v", err)
	}
	return m
}

func invocation(capability string, action connectors.Action, resource string) connectors.Invocation {
	return connectors.Invocation{TenantID: "t-1", WorkspaceID: "w-1", Subject: "u-1", AgentID: "a-1", ConnectorID: "c-1", Capability: capability, Action: action, Resource: resource, TraceID: "trace-1"}
}

func TestGuardEnforcesWorkspaceAndTenant(t *testing.T) {
	m := metaFor(t, connectors.ProviderS3, []string{"s3://bucket-a"}, nil)
	inv := invocation("object.list", connectors.ActionRead, "s3://bucket-a")
	if err := Guard(m, inv); err != nil {
		t.Fatalf("valid call rejected: %v", err)
	}
	inv.WorkspaceID = "w-2"
	if err := Guard(m, inv); err == nil || err.Error() != "connector not found" {
		t.Fatalf("cross-workspace call error = %v", err)
	}
	inv = invocation("object.list", connectors.ActionRead, "s3://bucket-a")
	inv.TenantID = "t-2"
	if err := Guard(m, inv); err == nil || err.Error() != "connector not found" {
		t.Fatalf("cross-tenant call error = %v", err)
	}
	m.Status = connectors.StatusSuspended
	if err := Guard(m, invocation("object.list", connectors.ActionRead, "s3://bucket-a")); err == nil || err.Error() != "connector is not active" {
		t.Fatalf("suspended connector error = %v", err)
	}
}

func TestGuardRejectsPathTraversal(t *testing.T) {
	m := metaFor(t, connectors.ProviderCRM, []string{"leads"}, nil)
	// "leads/../../customers" passes the contract's raw prefix check but
	// normalizes outside the leads boundary; the provider guard rejects it.
	inv := invocation("lead.read", connectors.ActionRead, "leads/../../customers")
	if err := Guard(m, inv); err == nil || err.Error() != "resource contains a path traversal segment" {
		t.Fatalf("traversal resource error = %v", err)
	}
	// A trailing traversal segment is rejected too.
	if err := Guard(m, invocation("lead.read", connectors.ActionRead, "leads/..")); err == nil {
		t.Fatal("trailing traversal segment was allowed")
	}
	// A bare ".." resource is rejected.
	if err := Guard(m, invocation("lead.read", connectors.ActionRead, "..")); err == nil {
		t.Fatal("bare .. resource was allowed")
	}
	// A name that merely contains ".." is not a traversal segment.
	if err := Guard(m, invocation("lead.read", connectors.ActionRead, "leads/foo..bar")); err != nil {
		t.Fatalf("benign resource rejected: %v", err)
	}
}

func TestGuardRejectsEncodedPathTraversal(t *testing.T) {
	m := metaFor(t, connectors.ProviderCRM, []string{"leads"}, nil)
	for _, resource := range []string{"leads/%2e%2e/customers", "leads/%2E%2E/customers", "leads/..%2fcustomers"} {
		inv := invocation("lead.read", connectors.ActionRead, resource)
		if err := Guard(m, inv); err == nil || err.Error() != "resource contains a path traversal segment" {
			t.Fatalf("encoded traversal %q error = %v", resource, err)
		}
	}
}

func TestS3ClientRejectsPathTraversal(t *testing.T) {
	c := NewS3(s3Store())
	m := metaFor(t, connectors.ProviderS3, []string{"s3://bucket-a"}, nil)
	// A key that traverses out of the sanctioned scope is rejected at the
	// guard even though it passes the raw prefix check.
	inv := invocation("object.read", connectors.ActionRead, "s3://bucket-a/exports/../../secrets/keys.pem")
	if _, err := c.Invoke(context.Background(), m, inv, ObjectReadPayload{Key: "exports/../../secrets/keys.pem"}); err == nil || err.Error() != "resource contains a path traversal segment" {
		t.Fatalf("s3 traversal error = %v", err)
	}
}

func TestDestructiveOperationsDisabledWithoutApproval(t *testing.T) {
	m := metaFor(t, connectors.ProviderS3, []string{"s3://bucket-a"}, nil)
	del := invocation("object.list", connectors.ActionDelete, "s3://bucket-a")
	if err := DestructiveAllowed(m, del); err == nil {
		t.Fatal("delete allowed without destructive_enabled")
	}
	m.Policy.DestructiveEnabled = true
	if err := DestructiveAllowed(m, del); err == nil {
		t.Fatal("delete allowed without approval")
	}
	del.ApprovalID = "approval-1"
	if err := DestructiveAllowed(m, del); err != nil {
		t.Fatalf("approved delete rejected: %v", err)
	}
	if err := DestructiveAllowed(m, invocation("object.list", connectors.ActionRead, "s3://bucket-a")); err != nil {
		t.Fatalf("read gated by destructive check: %v", err)
	}
}

// TestClientAcceptsMutationWithoutContractApproval documents the decoupling:
// the provider client no longer demands an approval reference for mutating
// capabilities; approval is decided by the ABAC policy layer. The client still
// fails closed on a payload/capability mismatch.
func TestClientAcceptsMutationWithoutContractApproval(t *testing.T) {
	m := metaFor(t, connectors.ProviderGitHub, []string{"repos/acme/kei"}, nil)
	m.Policy.AllowedActions = []connectors.Action{connectors.ActionRead, connectors.ActionCreate}
	c := NewGitHub(githubStore())
	inv := invocation("issue.create", connectors.ActionCreate, "repos/acme/kei")
	if _, err := c.Invoke(context.Background(), m, inv, RepositoryReadPayload{}); err == nil || !strings.Contains(err.Error(), "does not match invocation capability") {
		t.Fatalf("mutation without approval should fail only on payload mismatch, got %v", err)
	}
}

func googleStore() MemoryGoogle {
	now := time.Now().UTC()
	return MemoryGoogle{Files: map[string]GoogleFile{
		"f-1": {ID: "f-1", DriveID: "d-1", Name: "Q3 Plan", MimeType: "application/vnd.google-apps.document", ModifiedAt: now, Content: "plan body"},
		"f-2": {ID: "f-2", DriveID: "d-1", Name: "Budget", MimeType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", ModifiedAt: now},
		"f-3": {ID: "f-3", DriveID: "d-2", Name: "Other Drive Doc", MimeType: "text/plain", ModifiedAt: now, Content: "other drive"},
	}}
}

func TestGoogleDriveSearchMetadataAndDocsRead(t *testing.T) {
	c := NewGoogle(googleStore())
	m := metaFor(t, connectors.ProviderGoogle, []string{"drive/d-1"}, nil)
	ctx := context.Background()

	files, err := c.Invoke(ctx, m, invocation("drive.search", connectors.ActionRead, "drive/d-1"), DriveSearchPayload{Query: "plan"})
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if got := files.Data.([]GoogleFile); len(got) != 1 || got[0].ID != "f-1" {
		t.Fatalf("search results = %+v", got)
	}

	meta, err := c.Invoke(ctx, m, invocation("drive.metadata.read", connectors.ActionRead, "drive/d-1/files/f-1"), DriveMetadataPayload{})
	if err != nil {
		t.Fatalf("metadata failed: %v", err)
	}
	if got := meta.Data.(GoogleFile); got.ID != "f-1" || got.DriveID != "d-1" {
		t.Fatalf("metadata = %+v", got)
	}

	doc, err := c.Invoke(ctx, m, invocation("docs.read", connectors.ActionRead, "drive/d-1/files/f-1"), DocsReadPayload{})
	if err != nil {
		t.Fatalf("docs read failed: %v", err)
	}
	if doc.Data.(string) != "plan body" {
		t.Fatalf("doc content = %q", doc.Data)
	}

	if _, err := c.Invoke(ctx, m, invocation("docs.read", connectors.ActionRead, "drive/d-1/files/f-2"), DocsReadPayload{}); err == nil {
		t.Fatal("non-document read was allowed")
	}
	if _, err := c.Invoke(ctx, m, invocation("drive.metadata.read", connectors.ActionRead, "drive/d-1/files"), DriveMetadataPayload{}); err == nil {
		t.Fatal("malformed file resource was accepted")
	}
}

func TestGoogleDriveIsolationAndPrefixBoundary(t *testing.T) {
	c := NewGoogle(googleStore())
	ctx := context.Background()

	// A file that lives in another drive cannot be read through this connector.
	m := metaFor(t, connectors.ProviderGoogle, []string{"drive/d-1"}, nil)
	if _, err := c.Invoke(ctx, m, invocation("drive.metadata.read", connectors.ActionRead, "drive/d-1/files/f-3"), DriveMetadataPayload{}); err == nil {
		t.Fatal("cross-drive file read was allowed")
	}

	// Prefix policy is dormant on the provider client: boundary decisions live
	// in the ABAC policy layer, so a drive outside the declared prefixes is
	// still reachable structurally.
	scoped := metaFor(t, connectors.ProviderGoogle, []string{"drive/d-1", "drive/d-2"}, []string{"drive/d-1"})
	if _, err := c.Invoke(ctx, scoped, invocation("drive.search", connectors.ActionRead, "drive/d-2"), DriveSearchPayload{}); err != nil {
		t.Fatalf("drive outside the dormant prefix policy rejected: %v", err)
	}
	if _, err := c.Invoke(ctx, scoped, invocation("drive.search", connectors.ActionRead, "drive/d-1"), DriveSearchPayload{}); err != nil {
		t.Fatalf("in-prefix drive rejected: %v", err)
	}
}

func linearStore() MemoryLinear {
	now := time.Now().UTC()
	return MemoryLinear{
		Teams:    map[string]LinearTeam{"KEI": {Key: "KEI", Name: "Kei Engineering"}},
		Projects: map[string]LinearProject{"p-1": {ID: "p-1", TeamKey: "KEI", Name: "Connectors", State: "started"}},
		Cycles:   map[string]LinearCycle{"cy-1": {ID: "cy-1", TeamKey: "KEI", Name: "Cycle 12", StartsAt: now, EndsAt: now.AddDate(0, 1, 0)}},
		Issues:   map[string]LinearIssue{"i-1": {ID: "i-1", TeamKey: "KEI", ProjectID: "p-1", Identifier: "KEI-42", Title: "Provider seams", State: "in_progress"}},
	}
}

func TestLinearReads(t *testing.T) {
	c := NewLinear(linearStore())
	m := metaFor(t, connectors.ProviderLinear, []string{"linear/team/KEI", "linear/project/p-1", "linear/cycle/cy-1", "linear/issue/i-1"}, nil)
	ctx := context.Background()

	team, err := c.Invoke(ctx, m, invocation("team.read", connectors.ActionRead, "linear/team/KEI"), TeamReadPayload{})
	if err != nil || team.Data.(LinearTeam).Key != "KEI" {
		t.Fatalf("team read = %+v, %v", team.Data, err)
	}
	project, err := c.Invoke(ctx, m, invocation("project.read", connectors.ActionRead, "linear/project/p-1"), ProjectReadPayload{})
	if err != nil || project.Data.(LinearProject).Name != "Connectors" {
		t.Fatalf("project read = %+v, %v", project.Data, err)
	}
	cycle, err := c.Invoke(ctx, m, invocation("cycle.read", connectors.ActionRead, "linear/cycle/cy-1"), CycleReadPayload{})
	if err != nil || cycle.Data.(LinearCycle).Name != "Cycle 12" {
		t.Fatalf("cycle read = %+v, %v", cycle.Data, err)
	}
	issue, err := c.Invoke(ctx, m, invocation("issue.read", connectors.ActionRead, "linear/issue/i-1"), LinearIssueReadPayload{})
	if err != nil || issue.Data.(LinearIssue).Identifier != "KEI-42" {
		t.Fatalf("issue read = %+v, %v", issue.Data, err)
	}

	if _, err := c.Invoke(ctx, m, invocation("issue.read", connectors.ActionRead, "linear/project/p-1"), LinearIssueReadPayload{}); err == nil {
		t.Fatal("issue read against a project resource was allowed")
	}
	if _, err := c.Invoke(ctx, m, invocation("team.read", connectors.ActionRead, "linear/team/KEI"), LinearIssueReadPayload{}); err == nil {
		t.Fatal("payload/capability mismatch was accepted")
	}
}

func githubStore() MemoryGitHub {
	return MemoryGitHub{
		Repositories: map[string]GitHubRepository{"acme/kei": {Owner: "acme", Name: "kei", Description: "governed agents", DefaultBranch: "main", Private: true}},
		Issues:       map[string]GitHubIssue{"acme/kei#7": {Owner: "acme", Repo: "kei", Number: 7, Title: "connector bug", State: "open", Labels: []string{"bug"}}},
		PullRequests: map[string]GitHubPullRequest{"acme/kei#8": {Owner: "acme", Repo: "kei", Number: 8, Title: "provider seams", State: "open", HeadBranch: "feat/seams", BaseBranch: "main"}},
		Checks:       map[string]GitHubCheck{"acme/kei/chk-1": {Owner: "acme", Repo: "kei", ID: "chk-1", Name: "ci", Status: "completed", Conclusion: "success"}},
		Workflows:    map[string]GitHubWorkflow{"acme/kei/wf-1": {Owner: "acme", Repo: "kei", ID: "wf-1", Name: "build", Path: ".github/workflows/build.yml", State: "active"}},
	}
}

func TestGitHubReads(t *testing.T) {
	c := NewGitHub(githubStore())
	m := metaFor(t, connectors.ProviderGitHub, []string{"repos/acme/kei", "repos/acme/kei/issues/7", "repos/acme/kei/pulls/8", "repos/acme/kei/checks/chk-1", "repos/acme/kei/workflows/wf-1"}, nil)
	ctx := context.Background()

	repo, err := c.Invoke(ctx, m, invocation("repository.read", connectors.ActionRead, "repos/acme/kei"), RepositoryReadPayload{})
	if err != nil || repo.Data.(GitHubRepository).DefaultBranch != "main" {
		t.Fatalf("repository read = %+v, %v", repo.Data, err)
	}
	issue, err := c.Invoke(ctx, m, invocation("issue.read", connectors.ActionRead, "repos/acme/kei/issues/7"), GitHubIssueReadPayload{})
	if err != nil || issue.Data.(GitHubIssue).Number != 7 {
		t.Fatalf("issue read = %+v, %v", issue.Data, err)
	}
	pr, err := c.Invoke(ctx, m, invocation("pull_request.read", connectors.ActionRead, "repos/acme/kei/pulls/8"), PullRequestReadPayload{})
	if err != nil || pr.Data.(GitHubPullRequest).HeadBranch != "feat/seams" {
		t.Fatalf("pull request read = %+v, %v", pr.Data, err)
	}
	check, err := c.Invoke(ctx, m, invocation("check.read", connectors.ActionRead, "repos/acme/kei/checks/chk-1"), CheckReadPayload{})
	if err != nil || check.Data.(GitHubCheck).Conclusion != "success" {
		t.Fatalf("check read = %+v, %v", check.Data, err)
	}
	workflow, err := c.Invoke(ctx, m, invocation("workflow.read", connectors.ActionRead, "repos/acme/kei/workflows/wf-1"), WorkflowReadPayload{})
	if err != nil || workflow.Data.(GitHubWorkflow).Name != "build" {
		t.Fatalf("workflow read = %+v, %v", workflow.Data, err)
	}

	// Issue numbers under a sanctioned repository are in scope; a missing
	// issue fails at the backend.
	if _, err := c.Invoke(ctx, m, invocation("issue.read", connectors.ActionRead, "repos/acme/kei/issues/9"), GitHubIssueReadPayload{}); err == nil || err.Error() != "issue not found" {
		t.Fatalf("missing issue error = %v", err)
	}
	if _, err := c.Invoke(ctx, m, invocation("repository.read", connectors.ActionRead, "repos/evil/kei"), RepositoryReadPayload{}); err == nil {
		t.Fatal("out-of-scope repository was read")
	}
}

func s3Store() MemoryS3 {
	now := time.Now().UTC()
	return MemoryS3{
		Objects: map[string]S3Object{
			"bucket-a/exports/2026/report.json": {Key: "exports/2026/report.json", Size: 12, ETag: "e1", LastModified: now},
			"bucket-a/exports/2026/notes.txt":   {Key: "exports/2026/notes.txt", Size: 5, ETag: "e2", LastModified: now},
			"bucket-a/secrets/keys.pem":         {Key: "secrets/keys.pem", Size: 9, ETag: "e3", LastModified: now},
		},
		Contents: map[string][]byte{
			"bucket-a/exports/2026/report.json": []byte(`{"ok":true}`),
			"bucket-a/secrets/keys.pem":         []byte("pem-data"),
		},
	}
}

func TestS3ScopedListAndGet(t *testing.T) {
	c := NewS3(s3Store())
	m := metaFor(t, connectors.ProviderS3, []string{"s3://bucket-a"}, []string{"s3://bucket-a/exports"})
	ctx := context.Background()

	list, err := c.Invoke(ctx, m, invocation("object.list", connectors.ActionRead, "s3://bucket-a/exports"), ObjectListPayload{Prefix: "exports"})
	if err != nil {
		t.Fatalf("scoped list failed: %v", err)
	}
	if got := list.Data.([]S3Object); len(got) != 2 {
		t.Fatalf("scoped list returned %d objects: %+v", len(got), got)
	}

	get, err := c.Invoke(ctx, m, invocation("object.read", connectors.ActionRead, "s3://bucket-a/exports/2026/report.json"), ObjectReadPayload{Key: "exports/2026/report.json"})
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if string(get.Data.([]byte)) != `{"ok":true}` {
		t.Fatalf("object content = %q", get.Data)
	}

	// Resource and prefix boundaries are decided by the ABAC policy layer, not
	// the provider client, so structurally valid reads outside the declared
	// prefix are no longer denied here.
	if _, err := c.Invoke(ctx, m, invocation("object.read", connectors.ActionRead, "s3://bucket-a/secrets/keys.pem"), ObjectReadPayload{Key: "secrets/keys.pem"}); err != nil {
		t.Fatalf("structurally valid read outside the dormant prefix rejected: %v", err)
	}
	// A whole-bucket list is still structurally valid.
	if _, err := c.Invoke(ctx, m, invocation("object.list", connectors.ActionRead, "s3://bucket-a"), ObjectListPayload{}); err != nil {
		t.Fatalf("whole-bucket list rejected: %v", err)
	}
	// The payload cannot point at a different target than the resource.
	if _, err := c.Invoke(ctx, m, invocation("object.read", connectors.ActionRead, "s3://bucket-a/exports/2026/report.json"), ObjectReadPayload{Key: "secrets/keys.pem"}); err == nil {
		t.Fatal("payload/resource key mismatch was accepted")
	}
	// Other buckets are structurally valid too.
	if _, err := c.Invoke(ctx, m, invocation("object.list", connectors.ActionRead, "s3://bucket-b"), ObjectListPayload{}); err != nil {
		t.Fatalf("list in another bucket rejected: %v", err)
	}
}

func crmStore() MemoryCRM {
	now := time.Now().UTC()
	return MemoryCRM{
		AuthRef: "vault/tenant/t-1/test",
		Leads: map[string]CRMLead{
			"l-1": {ID: "l-1", Name: "Ada", Email: "ada@example.com", Company: "Acme", Stage: "qualified", CreatedAt: now},
			"l-2": {ID: "l-2", Name: "Grace", Email: "grace@example.com", Company: "Bolt", Stage: "contacted", CreatedAt: now},
		},
	}
}

func TestCRMLeadReadsRequireAuthentication(t *testing.T) {
	c := NewCRM(crmStore())
	m := metaFor(t, connectors.ProviderCRM, []string{"leads"}, nil)
	ctx := context.Background()

	list, err := c.Invoke(ctx, m, invocation("lead.read", connectors.ActionRead, "leads"), LeadReadPayload{})
	if err != nil {
		t.Fatalf("lead list failed: %v", err)
	}
	if got := list.Data.([]CRMLead); len(got) != 2 {
		t.Fatalf("lead list returned %d leads: %+v", len(got), got)
	}

	lead, err := c.Invoke(ctx, m, invocation("lead.read", connectors.ActionRead, "leads/l-1"), LeadReadPayload{ID: "l-1"})
	if err != nil || lead.Data.(CRMLead).Name != "Ada" {
		t.Fatalf("lead read = %+v, %v", lead.Data, err)
	}

	// A connector whose credential reference does not resolve to an
	// authenticated session is rejected by the backend.
	unauth := metaFor(t, connectors.ProviderCRM, []string{"leads"}, nil)
	unauth.CredentialRef = "vault/tenant/other/crm"
	if _, err := c.Invoke(ctx, unauth, invocation("lead.read", connectors.ActionRead, "leads"), LeadReadPayload{}); err == nil || err.Error() != "crm session is not authenticated" {
		t.Fatalf("unauthenticated lead list error = %v", err)
	}

	// The payload cannot point at a different lead than the resource.
	if _, err := c.Invoke(ctx, m, invocation("lead.read", connectors.ActionRead, "leads/l-1"), LeadReadPayload{ID: "l-2"}); err == nil {
		t.Fatal("payload/resource id mismatch was accepted")
	}
	// Customer reads are not part of the contract surface: the provider client
	// enforces resource shape (structural), not policy boundaries.
	if _, err := c.Invoke(ctx, m, invocation("lead.read", connectors.ActionRead, "customers"), LeadReadPayload{}); err == nil || err.Error() != "resource must be leads or leads/<id>" {
		t.Fatalf("customer resource error = %v", err)
	}
}

func notionStore() MemoryNotion {
	now := time.Now().UTC()
	return MemoryNotion{
		Pages: map[string]NotionPage{
			"p-1": {ID: "p-1", Title: "Onboarding Guide", ParentID: "d-1", ParentType: "database_id", URL: "https://notion.so/p-1", CreatedAt: now, UpdatedAt: now, Properties: map[string]any{"status": "Published"}},
			"p-2": {ID: "p-2", Title: "API Reference", ParentID: "d-1", ParentType: "database_id", URL: "https://notion.so/p-2", CreatedAt: now, UpdatedAt: now},
			"p-3": {ID: "p-3", Title: "Meeting Notes", ParentID: "d-2", ParentType: "database_id", URL: "https://notion.so/p-3", CreatedAt: now, UpdatedAt: now},
		},
		Databases: map[string]NotionDatabase{
			"d-1": {ID: "d-1", Title: "Engineering Wiki", URL: "https://notion.so/d-1", CreatedAt: now, UpdatedAt: now},
			"d-2": {ID: "d-2", Title: "Meetings", URL: "https://notion.so/d-2", CreatedAt: now, UpdatedAt: now},
		},
	}
}

func TestNotionPageRead(t *testing.T) {
	c := NewNotion(notionStore())
	m := metaFor(t, connectors.ProviderNotion, []string{"pages/p-1", "pages/p-2", "databases/d-1"}, nil)
	ctx := context.Background()

	page, err := c.Invoke(ctx, m, invocation("page.read", connectors.ActionRead, "pages/p-1"), PageReadPayload{})
	if err != nil || page.Data.(NotionPage).Title != "Onboarding Guide" {
		t.Fatalf("page read = %+v, %v", page.Data, err)
	}

	// Missing page is rejected.
	if _, err := c.Invoke(ctx, m, invocation("page.read", connectors.ActionRead, "pages/p-99"), PageReadPayload{}); err == nil || err.Error() != "page not found" {
		t.Fatalf("missing page error = %v", err)
	}
}

func TestNotionDatabaseQuery(t *testing.T) {
	c := NewNotion(notionStore())
	m := metaFor(t, connectors.ProviderNotion, []string{"databases/d-1"}, nil)
	ctx := context.Background()

	pages, err := c.Invoke(ctx, m, invocation("database.query", connectors.ActionRead, "databases/d-1"), DatabaseQueryPayload{})
	if err != nil {
		t.Fatalf("database query failed: %v", err)
	}
	got := pages.Data.([]NotionPage)
	if len(got) != 2 {
		t.Fatalf("database query returned %d pages, want 2", len(got))
	}

	// Missing database is rejected.
	if _, err := c.Invoke(ctx, m, invocation("database.query", connectors.ActionRead, "databases/d-99"), DatabaseQueryPayload{}); err == nil || err.Error() != "database not found" {
		t.Fatalf("missing database error = %v", err)
	}
}

func TestNotionSearch(t *testing.T) {
	c := NewNotion(notionStore())
	m := metaFor(t, connectors.ProviderNotion, []string{"search"}, nil)
	ctx := context.Background()

	// Search across all content.
	result, err := c.Invoke(ctx, m, invocation("search", connectors.ActionRead, "search"), SearchPayload{Query: "Onboarding"})
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	sr := result.Data.(NotionSearchResult)
	if len(sr.Results) != 1 {
		t.Fatalf("search returned %d results, want 1", len(sr.Results))
	}

	// Empty query returns all.
	result, err = c.Invoke(ctx, m, invocation("search", connectors.ActionRead, "search"), SearchPayload{})
	if err != nil {
		t.Fatalf("empty search failed: %v", err)
	}
	sr = result.Data.(NotionSearchResult)
	if len(sr.Results) != 5 {
		t.Fatalf("empty search returned %d results, want 5 (3 pages + 2 databases)", len(sr.Results))
	}

	// Payload capability mismatch is rejected.
	if _, err := c.Invoke(ctx, m, invocation("search", connectors.ActionRead, "search"), PageReadPayload{}); err == nil {
		t.Fatal("payload mismatch was accepted")
	}
}

func TestNotionStructuralBoundary(t *testing.T) {
	c := NewNotion(notionStore())
	m := metaFor(t, connectors.ProviderNotion, []string{"pages/p-1", "databases/d-1"}, nil)
	ctx := context.Background()

	// Malformed resource paths are rejected at the structural level.
	if _, err := c.Invoke(ctx, m, invocation("page.read", connectors.ActionRead, "pages"), PageReadPayload{}); err == nil {
		t.Fatal("bare pages resource was accepted")
	}
	if _, err := c.Invoke(ctx, m, invocation("database.query", connectors.ActionRead, "databases"), DatabaseQueryPayload{}); err == nil {
		t.Fatal("bare databases resource was accepted")
	}
	if _, err := c.Invoke(ctx, m, invocation("page.read", connectors.ActionRead, "databases/d-1"), PageReadPayload{}); err == nil {
		t.Fatal("cross-capability resource format was accepted")
	}

	// A missing page returns a backend error.
	if _, err := c.Invoke(ctx, m, invocation("page.read", connectors.ActionRead, "pages/p-99"), PageReadPayload{}); err == nil || err.Error() != "page not found" {
		t.Fatalf("missing page error = %v", err)
	}

	// Payload/capability mismatch is rejected.
	if _, err := c.Invoke(ctx, m, invocation("page.read", connectors.ActionRead, "pages/p-1"), SearchPayload{}); err == nil {
		t.Fatal("payload mismatch was accepted")
	}
}

func TestRegistryAndProviderMismatch(t *testing.T) {
	for _, p := range []connectors.Provider{connectors.ProviderCRM, connectors.ProviderLinear, connectors.ProviderGitHub, connectors.ProviderGoogle, connectors.ProviderNotion, connectors.ProviderS3} {
		c, err := New(p)
		if err != nil {
			t.Fatalf("New(%s) failed: %v", p, err)
		}
		if c.Provider() != p {
			t.Fatalf("New(%s).Provider() = %s", p, c.Provider())
		}
	}
	if _, err := New(connectors.Provider("slack")); err == nil {
		t.Fatal("unknown provider was accepted")
	}

	c := NewS3(s3Store())
	foreign := metaFor(t, connectors.ProviderGitHub, []string{"repos/acme/kei"}, nil)
	if _, err := c.Invoke(context.Background(), foreign, invocation("repository.read", connectors.ActionRead, "repos/acme/kei"), RepositoryReadPayload{}); err == nil || !strings.Contains(err.Error(), "does not match client") {
		t.Fatalf("foreign connector error = %v", err)
	}
}

func TestPayloadCapabilityMismatchRejected(t *testing.T) {
	c := NewGoogle(googleStore())
	m := metaFor(t, connectors.ProviderGoogle, []string{"drive/d-1"}, nil)
	inv := invocation("docs.read", connectors.ActionRead, "drive/d-1/files/f-1")
	if _, err := c.Invoke(context.Background(), m, inv, DriveSearchPayload{}); err == nil || !strings.Contains(err.Error(), "does not match invocation capability") {
		t.Fatalf("payload mismatch error = %v", err)
	}
}
