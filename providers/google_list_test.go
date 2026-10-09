package providers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/HaikeiLabs/kei-connector-contracts/contract"
)

func TestAuthenticatedGoogleListsOnlyRequestedBoundedDriveFiles(t *testing.T) {
	resolver := &recordingResolver{token: "test-token"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/drive/v3/files" {
			t.Errorf("request = %s %s", r.Method, r.URL.Path)
		}
		query := r.URL.Query().Get("q")
		if !strings.Contains(query, "trashed = false") || !strings.Contains(query, "name contains 'plan\\'s'") || !strings.Contains(query, "mimeType = 'text/plain'") {
			t.Errorf("Drive query = %q", query)
		}
		if r.URL.Query().Get("driveId") != "drive-1" || r.URL.Query().Get("corpora") != "drive" ||
			r.URL.Query().Get("pageSize") != "25" || r.URL.Query().Get("supportsAllDrives") != "true" ||
			r.URL.Query().Get("includeItemsFromAllDrives") != "true" {
			t.Errorf("unexpected list parameters: %v", r.URL.Query())
		}
		if r.URL.Query().Get("fields") != "files(id,driveId,name,mimeType,modifiedTime)" {
			t.Errorf("fields = %q", r.URL.Query().Get("fields"))
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("authorization header missing delegated token")
		}
		_, _ = io.WriteString(w, `{"files":[{"id":"f-1","driveId":"drive-1","name":"Plan","mimeType":"text/plain"}]}`)
	}))
	defer server.Close()

	client := NewGoogleHTTP(testHTTPClient(server, func(_ *testing.T, req *http.Request) {
		if req.URL.Scheme != "https" || req.URL.Host != "www.googleapis.com" {
			t.Errorf("provider endpoint = %s", req.URL)
		}
	}), resolver)
	meta := metaFor(t, contract.ProviderGoogle, []string{"drive/drive-1"}, nil)
	result, err := client.Invoke(context.Background(), meta, invocation("drive.search", contract.ActionRead, "drive/drive-1"), DriveSearchPayload{
		Query: "plan's", MimeType: "text/plain", PageSize: 25,
	})
	if err != nil {
		t.Fatalf("list invocation: %v", err)
	}
	files, ok := result.Data.([]GoogleFile)
	if !ok || len(files) != 1 || files[0].ID != "f-1" || files[0].DriveID != "drive-1" {
		t.Fatalf("normalized files = %#v", result.Data)
	}
}

func TestAuthenticatedGoogleRejectsUnboundedPageSize(t *testing.T) {
	backend := &AuthenticatedGoogle{httpClient: http.DefaultClient, credentials: &recordingResolver{token: "test-token"}}
	if _, err := backend.ListFiles(context.Background(), "drive-1", "", "", 101); err == nil {
		t.Fatal("accepted a Drive page size over the runtime bound")
	}
}
