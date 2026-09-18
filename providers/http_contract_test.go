package providers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/HaikeiLabs/kei-connector-contracts"
)

type recordingResolver struct {
	token                        string
	ten, workspace, subject, ref string
}

func (r *recordingResolver) Resolve(_ context.Context, tenant, workspace, subject, ref string) (string, error) {
	r.ten, r.workspace, r.subject, r.ref = tenant, workspace, subject, ref
	return r.token, nil
}

type rewriteTransport struct {
	serverURL    string
	check        func(*testing.T, *http.Request)
	roundTripper func(*http.Request) (*http.Response, error)
}

func (t rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.check(nil, req)
	copy := req.Clone(req.Context())
	copy.URL.Scheme = "http"
	copy.URL.Host = strings.TrimPrefix(t.serverURL, "http://")
	return t.roundTripper(copy)
}

func testHTTPClient(server *httptest.Server, check func(*testing.T, *http.Request)) *http.Client {
	return &http.Client{Transport: rewriteTransport{serverURL: server.URL, check: check, roundTripper: http.DefaultTransport.RoundTrip}}
}

func TestAuthenticatedGitHubContract(t *testing.T) {
	resolver := &recordingResolver{token: "secret-token"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/acme/kei/issues/7" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer secret-token" {
			t.Errorf("authorization = %q", got)
		}
		if got := r.Header.Get("X-GitHub-Api-Version"); got != "2022-11-28" {
			t.Errorf("api version = %q", got)
		}
		_, _ = io.WriteString(w, `{"number":7,"title":"hello","state":"open","labels":[{"name":"bug"}]}`)
	}))
	defer server.Close()
	client := NewGitHubHTTP(testHTTPClient(server, func(_ *testing.T, req *http.Request) {
		if req.URL.Scheme != "https" || req.URL.Host != "api.github.com" {
			t.Errorf("endpoint = %s", req.URL.String())
		}
	}), resolver)
	meta := metaFor(t, connectors.ProviderGitHub, []string{"repos/acme/kei"}, nil)
	inv := invocation("issue.read", connectors.ActionRead, "repos/acme/kei/issues/7")
	result, err := client.Invoke(context.Background(), meta, inv, GitHubIssueReadPayload{})
	if err != nil {
		t.Fatalf("invoke: %v", err)
	}
	issue := result.Data.(GitHubIssue)
	if issue.Owner != "acme" || issue.Repo != "kei" || issue.Number != 7 || issue.Labels[0] != "bug" {
		t.Fatalf("issue = %+v", issue)
	}
	if resolver.ten != "t-1" || resolver.workspace != "w-1" || resolver.subject != "u-1" || resolver.ref != meta.CredentialRef {
		t.Fatalf("scope = %+v", resolver)
	}
}

func TestAuthenticatedGoogleContract(t *testing.T) {
	resolver := &recordingResolver{token: "google-secret"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/files/file-1" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if got := r.URL.Query().Get("alt"); got != "media" {
			t.Errorf("alt = %q", got)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer google-secret" {
			t.Errorf("authorization = %q", got)
		}
		_, _ = io.WriteString(w, "document text")
	}))
	defer server.Close()
	client := NewGoogleHTTP(testHTTPClient(server, func(_ *testing.T, req *http.Request) {
		if req.URL.Scheme != "https" || req.URL.Host != "www.googleapis.com" {
			t.Errorf("endpoint = %s", req.URL.String())
		}
	}), resolver)
	meta := metaFor(t, connectors.ProviderGoogle, []string{"drive/drive-1/files/file-1"}, nil)
	inv := invocation("docs.read", connectors.ActionRead, "drive/drive-1/files/file-1")
	// The metadata request is intentionally mocked as a separate fixed endpoint
	// response by using the same server and dispatching on the query.
	server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("alt") == "media" {
			_, _ = io.WriteString(w, "document text")
			return
		}
		_, _ = io.WriteString(w, `{"id":"file-1","driveId":"drive-1","name":"notes","mimeType":"text/plain"}`)
	})
	result, err := client.Invoke(context.Background(), meta, inv, DocsReadPayload{})
	if err != nil {
		t.Fatalf("invoke: %v", err)
	}
	if result.Data != "document text" {
		t.Fatalf("content = %v", result.Data)
	}
}

func TestAuthenticatedProviderFailsClosedWithoutCredentialResolver(t *testing.T) {
	client, err := New(connectors.ProviderGitHub)
	if err != nil {
		t.Fatal(err)
	}
	meta := metaFor(t, connectors.ProviderGitHub, []string{"repos/acme/kei"}, nil)
	_, err = client.Invoke(context.Background(), meta, invocation("repository.read", connectors.ActionRead, "repos/acme/kei"), RepositoryReadPayload{})
	if err == nil || !strings.Contains(err.Error(), "credential") {
		t.Fatalf("error = %v", err)
	}
	if strings.Contains(err.Error(), "secret") {
		t.Fatalf("secret material in error: %v", err)
	}
}

func TestAuthenticatedProviderDoesNotExposeCredentialErrors(t *testing.T) {
	resolver := credentialErrorResolver{}
	client := NewGitHubHTTP(&http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) { return nil, fmt.Errorf("must not send") })}, resolver)
	meta := metaFor(t, connectors.ProviderGitHub, []string{"repos/acme/kei"}, nil)
	_, err := client.Invoke(context.Background(), meta, invocation("repository.read", connectors.ActionRead, "repos/acme/kei"), RepositoryReadPayload{})
	if err == nil || err.Error() != "credential resolution failed" {
		t.Fatalf("error = %v", err)
	}
}

type credentialErrorResolver struct{}

func (credentialErrorResolver) Resolve(context.Context, string, string, string, string) (string, error) {
	return "", fmt.Errorf("secret-token leaked")
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
