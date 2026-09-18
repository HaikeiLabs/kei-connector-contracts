package providers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const githubAPIBase = "https://api.github.com"

// AuthenticatedGitHub is the production GitHub backend. Its base endpoint is
// fixed so connector resources cannot redirect requests to another service.
type AuthenticatedGitHub struct {
	httpClient  *http.Client
	credentials CredentialResolver
}

func NewAuthenticatedGitHub(config RuntimeConfig) *GitHubClient {
	config = normalizeRuntimeConfig(config)
	return NewGitHub(&AuthenticatedGitHub{httpClient: config.HTTPClient, credentials: config.Credentials})
}

func NewGitHubHTTP(httpClient *http.Client, credentials CredentialResolver) *GitHubClient {
	return NewAuthenticatedGitHub(RuntimeConfig{HTTPClient: httpClient, Credentials: credentials})
}

func (b *AuthenticatedGitHub) Repository(ctx context.Context, owner, name string) (GitHubRepository, error) {
	var response struct {
		Owner struct {
			Login string `json:"login"`
		} `json:"owner"`
		Name          string  `json:"name"`
		Description   *string `json:"description"`
		DefaultBranch string  `json:"default_branch"`
		Private       bool    `json:"private"`
	}
	if err := b.getJSON(ctx, "/repos/"+url.PathEscape(owner)+"/"+url.PathEscape(name), &response); err != nil {
		return GitHubRepository{}, err
	}
	return GitHubRepository{Owner: response.Owner.Login, Name: response.Name, Description: stringValue(response.Description), DefaultBranch: response.DefaultBranch, Private: response.Private}, nil
}

func (b *AuthenticatedGitHub) Issue(ctx context.Context, owner, name string, number int) (GitHubIssue, error) {
	var response struct {
		User struct {
			Login string `json:"login"`
		} `json:"user"`
		Number int    `json:"number"`
		Title  string `json:"title"`
		State  string `json:"state"`
		Labels []struct {
			Name string `json:"name"`
		} `json:"labels"`
	}
	if err := b.getJSON(ctx, fmt.Sprintf("/repos/%s/%s/issues/%d", url.PathEscape(owner), url.PathEscape(name), number), &response); err != nil {
		return GitHubIssue{}, err
	}
	labels := make([]string, 0, len(response.Labels))
	for _, label := range response.Labels {
		labels = append(labels, label.Name)
	}
	return GitHubIssue{Owner: owner, Repo: name, Number: response.Number, Title: response.Title, State: response.State, Labels: labels}, nil
}

func (b *AuthenticatedGitHub) PullRequest(ctx context.Context, owner, name string, number int) (GitHubPullRequest, error) {
	var response struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
		State  string `json:"state"`
		Head   struct {
			Ref string `json:"ref"`
		} `json:"head"`
		Base struct {
			Ref string `json:"ref"`
		} `json:"base"`
		Merged bool `json:"merged"`
	}
	if err := b.getJSON(ctx, fmt.Sprintf("/repos/%s/%s/pulls/%d", url.PathEscape(owner), url.PathEscape(name), number), &response); err != nil {
		return GitHubPullRequest{}, err
	}
	return GitHubPullRequest{Owner: owner, Repo: name, Number: response.Number, Title: response.Title, State: response.State, HeadBranch: response.Head.Ref, BaseBranch: response.Base.Ref, Merged: response.Merged}, nil
}

func (b *AuthenticatedGitHub) Check(ctx context.Context, owner, name, id string) (GitHubCheck, error) {
	var response struct {
		ID         json.Number `json:"id"`
		Name       string      `json:"name"`
		Status     string      `json:"status"`
		Conclusion *string     `json:"conclusion"`
	}
	if err := b.getJSON(ctx, "/repos/"+url.PathEscape(owner)+"/"+url.PathEscape(name)+"/check-runs/"+url.PathEscape(id), &response); err != nil {
		return GitHubCheck{}, err
	}
	return GitHubCheck{Owner: owner, Repo: name, ID: response.ID.String(), Name: response.Name, Status: response.Status, Conclusion: stringValue(response.Conclusion)}, nil
}

func (b *AuthenticatedGitHub) Workflow(ctx context.Context, owner, name, id string) (GitHubWorkflow, error) {
	var response struct {
		ID    json.Number `json:"id"`
		Name  string      `json:"name"`
		Path  string      `json:"path"`
		State string      `json:"state"`
	}
	if err := b.getJSON(ctx, "/repos/"+url.PathEscape(owner)+"/"+url.PathEscape(name)+"/actions/workflows/"+url.PathEscape(id), &response); err != nil {
		return GitHubWorkflow{}, err
	}
	return GitHubWorkflow{Owner: owner, Repo: name, ID: response.ID.String(), Name: response.Name, Path: response.Path, State: response.State}, nil
}

func (b *AuthenticatedGitHub) getJSON(ctx context.Context, path string, out any) error {
	scope, err := scopeFromContext(ctx)
	if err != nil {
		return err
	}
	token, err := b.credentials.Resolve(ctx, scope.tenantID, scope.workspaceID, scope.subject, scope.credential)
	if err != nil || token == "" {
		return errors.New("credential resolution failed")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, githubAPIBase+path, nil)
	if err != nil {
		return errors.New("provider request could not be created")
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
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

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

var _ GitHubBackend = (*AuthenticatedGitHub)(nil)
