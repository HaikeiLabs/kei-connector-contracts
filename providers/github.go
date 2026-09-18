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
	"strconv"
	"strings"

	"github.com/HaikeiLabs/kei-connector-contracts"
)

type GitHubRepository struct {
	Owner         string `json:"owner"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	DefaultBranch string `json:"default_branch"`
	Private       bool   `json:"private"`
}

type GitHubIssue struct {
	Owner  string   `json:"owner"`
	Repo   string   `json:"repo"`
	Number int      `json:"number"`
	Title  string   `json:"title"`
	State  string   `json:"state"`
	Labels []string `json:"labels"`
}

type GitHubPullRequest struct {
	Owner      string `json:"owner"`
	Repo       string `json:"repo"`
	Number     int    `json:"number"`
	Title      string `json:"title"`
	State      string `json:"state"`
	HeadBranch string `json:"head_branch"`
	BaseBranch string `json:"base_branch"`
	Merged     bool   `json:"merged"`
}

type GitHubCheck struct {
	Owner      string `json:"owner"`
	Repo       string `json:"repo"`
	ID         string `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
}

type GitHubWorkflow struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
	ID    string `json:"id"`
	Name  string `json:"name"`
	Path  string `json:"path"`
	State string `json:"state"`
}

type RepositoryReadPayload struct{}

func (RepositoryReadPayload) Capability() string { return "repository.read" }

type GitHubIssueReadPayload struct{}

func (GitHubIssueReadPayload) Capability() string { return "issue.read" }

type PullRequestReadPayload struct{}

func (PullRequestReadPayload) Capability() string { return "pull_request.read" }

type CheckReadPayload struct{}

func (CheckReadPayload) Capability() string { return "check.read" }

type WorkflowReadPayload struct{}

func (WorkflowReadPayload) Capability() string { return "workflow.read" }

// GitHubBackend is the seam for GitHub reads.
type GitHubBackend interface {
	Repository(ctx context.Context, owner, name string) (GitHubRepository, error)
	Issue(ctx context.Context, owner, name string, number int) (GitHubIssue, error)
	PullRequest(ctx context.Context, owner, name string, number int) (GitHubPullRequest, error)
	Check(ctx context.Context, owner, name, id string) (GitHubCheck, error)
	Workflow(ctx context.Context, owner, name, id string) (GitHubWorkflow, error)
}

// MemoryGitHub is the in-memory GitHub backend used until a real one exists.
type MemoryGitHub struct {
	Repositories map[string]GitHubRepository
	Issues       map[string]GitHubIssue
	PullRequests map[string]GitHubPullRequest
	Checks       map[string]GitHubCheck
	Workflows    map[string]GitHubWorkflow
}

func repoKey(owner, name string) string { return owner + "/" + name }

func (m MemoryGitHub) Repository(_ context.Context, owner, name string) (GitHubRepository, error) {
	r, ok := m.Repositories[repoKey(owner, name)]
	if !ok {
		return GitHubRepository{}, errors.New("repository not found")
	}
	return r, nil
}

func (m MemoryGitHub) Issue(_ context.Context, owner, name string, number int) (GitHubIssue, error) {
	i, ok := m.Issues[fmt.Sprintf("%s#%d", repoKey(owner, name), number)]
	if !ok {
		return GitHubIssue{}, errors.New("issue not found")
	}
	return i, nil
}

func (m MemoryGitHub) PullRequest(_ context.Context, owner, name string, number int) (GitHubPullRequest, error) {
	p, ok := m.PullRequests[fmt.Sprintf("%s#%d", repoKey(owner, name), number)]
	if !ok {
		return GitHubPullRequest{}, errors.New("pull request not found")
	}
	return p, nil
}

func (m MemoryGitHub) Check(_ context.Context, owner, name, id string) (GitHubCheck, error) {
	c, ok := m.Checks[repoKey(owner, name)+"/"+id]
	if !ok {
		return GitHubCheck{}, errors.New("check not found")
	}
	return c, nil
}

func (m MemoryGitHub) Workflow(_ context.Context, owner, name, id string) (GitHubWorkflow, error) {
	w, ok := m.Workflows[repoKey(owner, name)+"/"+id]
	if !ok {
		return GitHubWorkflow{}, errors.New("workflow not found")
	}
	return w, nil
}

// GitHubClient serves repository.read, issue.read, pull_request.read,
// check.read, and workflow.read. Resources are repos/<owner>/<repo> and
// repos/<owner>/<repo>/{issues,pulls,checks,workflows}/<id>.
type GitHubClient struct {
	backend GitHubBackend
}

func NewGitHub(backend GitHubBackend) *GitHubClient { return &GitHubClient{backend: backend} }

func (c *GitHubClient) Provider() connectors.Provider { return connectors.ProviderGitHub }

func (c *GitHubClient) Invoke(ctx context.Context, meta connectors.Metadata, inv connectors.Invocation, payload Payload) (Result, error) {
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
	case RepositoryReadPayload:
		owner, name, ok := repoResource(inv.Resource)
		if !ok {
			return Result{}, errors.New("resource must be repos/<owner>/<repo>")
		}
		data, err := c.backend.Repository(ctx, owner, name)
		if err != nil {
			return Result{}, err
		}
		return Result{Capability: inv.Capability, Data: data}, nil
	case GitHubIssueReadPayload:
		owner, name, tail, ok := subResource(inv.Resource, "issues")
		if !ok {
			return Result{}, errors.New("resource must be repos/<owner>/<repo>/issues/<number>")
		}
		number, err := strconv.Atoi(tail)
		if err != nil {
			return Result{}, errors.New("issue number is invalid")
		}
		data, err := c.backend.Issue(ctx, owner, name, number)
		if err != nil {
			return Result{}, err
		}
		return Result{Capability: inv.Capability, Data: data}, nil
	case PullRequestReadPayload:
		owner, name, tail, ok := subResource(inv.Resource, "pulls")
		if !ok {
			return Result{}, errors.New("resource must be repos/<owner>/<repo>/pulls/<number>")
		}
		number, err := strconv.Atoi(tail)
		if err != nil {
			return Result{}, errors.New("pull request number is invalid")
		}
		data, err := c.backend.PullRequest(ctx, owner, name, number)
		if err != nil {
			return Result{}, err
		}
		return Result{Capability: inv.Capability, Data: data}, nil
	case CheckReadPayload:
		owner, name, tail, ok := subResource(inv.Resource, "checks")
		if !ok {
			return Result{}, errors.New("resource must be repos/<owner>/<repo>/checks/<id>")
		}
		data, err := c.backend.Check(ctx, owner, name, tail)
		if err != nil {
			return Result{}, err
		}
		return Result{Capability: inv.Capability, Data: data}, nil
	case WorkflowReadPayload:
		owner, name, tail, ok := subResource(inv.Resource, "workflows")
		if !ok {
			return Result{}, errors.New("resource must be repos/<owner>/<repo>/workflows/<id>")
		}
		data, err := c.backend.Workflow(ctx, owner, name, tail)
		if err != nil {
			return Result{}, err
		}
		return Result{Capability: inv.Capability, Data: data}, nil
	default:
		return Result{}, fmt.Errorf("unsupported payload %T", p)
	}
}

func repoResource(resource string) (owner, name string, ok bool) {
	parts := strings.Split(resource, "/")
	if len(parts) != 3 || parts[0] != "repos" || parts[1] == "" || parts[2] == "" {
		return "", "", false
	}
	return parts[1], parts[2], true
}

func subResource(resource, kind string) (owner, name, tail string, ok bool) {
	parts := strings.Split(resource, "/")
	if len(parts) != 5 || parts[0] != "repos" || parts[3] != kind || parts[1] == "" || parts[2] == "" || parts[4] == "" {
		return "", "", "", false
	}
	return parts[1], parts[2], parts[4], true
}
