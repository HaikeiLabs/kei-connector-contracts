package providers

import (
	"context"
	"testing"

	"github.com/HaikeiLabs/kei-connector-contracts"
)

func TestConnectorContractFreezeResourceAndResultShapes(t *testing.T) {
	ctx := context.Background()
	github := NewGitHub(githubStore())
	gm := metaFor(t, connectors.ProviderGitHub, []string{"repos/acme/kei"}, nil)
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
		inv := invocation(tc.payload.Capability(), connectors.ActionRead, tc.resource)
		if _, err := github.Invoke(ctx, gm, inv, tc.payload); err != nil {
			t.Errorf("GitHub resource %q rejected: %v", tc.resource, err)
		}
	}

	google := NewGoogle(googleStore())
	drive := metaFor(t, connectors.ProviderGoogle, []string{"drive/d-1"}, nil)
	for _, tc := range []struct {
		resource string
		payload  Payload
	}{
		{"drive/d-1", DriveSearchPayload{}},
		{"drive/d-1/files/f-1", DriveMetadataPayload{}},
		{"drive/d-1/files/f-1", DocsReadPayload{}},
	} {
		inv := invocation(tc.payload.Capability(), connectors.ActionRead, tc.resource)
		if _, err := google.Invoke(ctx, drive, inv, tc.payload); err != nil {
			t.Errorf("Google Drive resource %q rejected: %v", tc.resource, err)
		}
	}

	linear := NewLinear(linearStore())
	lm := metaFor(t, connectors.ProviderLinear, []string{"linear/team/KEI"}, nil)
	for _, tc := range []struct {
		resource string
		payload  Payload
	}{
		{"linear/team/KEI", TeamReadPayload{}},
		{"linear/project/p-1", ProjectReadPayload{}},
		{"linear/cycle/cy-1", CycleReadPayload{}},
		{"linear/issue/i-1", LinearIssueReadPayload{}},
	} {
		inv := invocation(tc.payload.Capability(), connectors.ActionRead, tc.resource)
		if _, err := linear.Invoke(ctx, lm, inv, tc.payload); err != nil {
			t.Errorf("Linear resource %q rejected: %v", tc.resource, err)
		}
	}

	notion := NewNotion(notionStore())
	nm := metaFor(t, connectors.ProviderNotion, []string{"pages/p-1", "databases/d-1"}, nil)
	notionCases := []struct {
		resource string
		payload  Payload
	}{
		{"pages/p-1", PageReadPayload{}},
		{"databases/d-1", DatabaseQueryPayload{}},
		{"search", SearchPayload{Query: "test"}},
	}
	for _, tc := range notionCases {
		inv := invocation(tc.payload.Capability(), connectors.ActionRead, tc.resource)
		if _, err := notion.Invoke(ctx, nm, inv, tc.payload); err != nil {
			t.Errorf("Notion resource %q rejected: %v", tc.resource, err)
		}
	}
}
