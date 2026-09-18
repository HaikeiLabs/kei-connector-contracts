package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

const linearGraphQLEndpoint = "https://api.linear.app/graphql"

const (
	linearTeamQuery    = `query Team($key: String!) { team(key: $key) { id key name } }`
	linearProjectQuery = `query Project($id: String!) { project(id: $id) { id name state { name } team { key } } }`
	linearCycleQuery   = `query Cycle($id: String!) { cycle(id: $id) { id name startsAt endsAt team { key } } }`
	linearIssueQuery   = `query Issue($id: String!) { issue(id: $id) { id identifier title state { name } team { key } project { id } } }`
)

// HTTPLinearGraphQLTransport is the production transport. Its endpoint is
// intentionally not configurable: arbitrary destinations are outside the
// governed Linear provider contract.
type HTTPLinearGraphQLTransport struct {
	Client *http.Client
}

func (t HTTPLinearGraphQLTransport) Query(ctx context.Context, endpoint, token, query string, variables map[string]any) (json.RawMessage, error) {
	if endpoint != linearGraphQLEndpoint || strings.TrimSpace(token) == "" || query == "" {
		return nil, errors.New("invalid Linear GraphQL request")
	}
	client := t.Client
	if client == nil {
		client = http.DefaultClient
	}
	body, err := json.Marshal(struct {
		Query     string         `json:"query"`
		Variables map[string]any `json:"variables"`
	}{Query: query, Variables: variables})
	if err != nil {
		return nil, errors.New("encode Linear GraphQL request")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, linearGraphQLEndpoint, bytes.NewReader(body))
	if err != nil {
		return nil, errors.New("create Linear GraphQL request")
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Linear GraphQL request failed")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Linear GraphQL request failed")
	}
	limited := io.LimitReader(resp.Body, 1<<20)
	var envelope struct {
		Data   json.RawMessage `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	decoder := json.NewDecoder(limited)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&envelope); err != nil || len(envelope.Errors) != 0 || len(envelope.Data) == 0 || string(envelope.Data) == "null" {
		return nil, errors.New("invalid Linear GraphQL response")
	}
	return envelope.Data, nil
}

type linearGraphQLBackend struct {
	transport LinearGraphQLTransport
	token     string
}

func (b linearGraphQLBackend) query(ctx context.Context, query string, variables map[string]any, out any) error {
	data, err := b.transport.Query(ctx, linearGraphQLEndpoint, b.token, query, variables)
	if err != nil {
		return errors.New("Linear GraphQL query failed")
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(out); err != nil {
		return errors.New("invalid Linear resource response")
	}
	return nil
}

func (b linearGraphQLBackend) Team(ctx context.Context, key string) (LinearTeam, error) {
	var response struct {
		Team *struct {
			Key  string `json:"key"`
			Name string `json:"name"`
		} `json:"team"`
	}
	if err := b.query(ctx, linearTeamQuery, map[string]any{"key": key}, &response); err != nil || response.Team == nil || response.Team.Key == "" || response.Team.Name == "" {
		return LinearTeam{}, errors.New("invalid Linear team response")
	}
	return LinearTeam{Key: response.Team.Key, Name: response.Team.Name}, nil
}

func (b linearGraphQLBackend) Project(ctx context.Context, id string) (LinearProject, error) {
	var response struct {
		Project *struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			State *struct {
				Name string `json:"name"`
			} `json:"state"`
			Team *struct {
				Key string `json:"key"`
			} `json:"team"`
		} `json:"project"`
	}
	if err := b.query(ctx, linearProjectQuery, map[string]any{"id": id}, &response); err != nil || response.Project == nil || response.Project.ID == "" || response.Project.Name == "" || response.Project.State == nil || response.Project.State.Name == "" || response.Project.Team == nil || response.Project.Team.Key == "" {
		return LinearProject{}, errors.New("invalid Linear project response")
	}
	return LinearProject{ID: response.Project.ID, TeamKey: response.Project.Team.Key, Name: response.Project.Name, State: response.Project.State.Name}, nil
}

func (b linearGraphQLBackend) Cycle(ctx context.Context, id string) (LinearCycle, error) {
	var response struct {
		Cycle *struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			StartsAt string `json:"startsAt"`
			EndsAt   string `json:"endsAt"`
			Team     *struct {
				Key string `json:"key"`
			} `json:"team"`
		} `json:"cycle"`
	}
	if err := b.query(ctx, linearCycleQuery, map[string]any{"id": id}, &response); err != nil || response.Cycle == nil || response.Cycle.ID == "" || response.Cycle.Name == "" || response.Cycle.StartsAt == "" || response.Cycle.EndsAt == "" || response.Cycle.Team == nil || response.Cycle.Team.Key == "" {
		return LinearCycle{}, errors.New("invalid Linear cycle response")
	}
	starts, err := time.Parse(time.RFC3339, response.Cycle.StartsAt)
	if err != nil {
		return LinearCycle{}, errors.New("invalid Linear cycle response")
	}
	ends, err := time.Parse(time.RFC3339, response.Cycle.EndsAt)
	if err != nil {
		return LinearCycle{}, errors.New("invalid Linear cycle response")
	}
	return LinearCycle{ID: response.Cycle.ID, TeamKey: response.Cycle.Team.Key, Name: response.Cycle.Name, StartsAt: starts, EndsAt: ends}, nil
}

func (b linearGraphQLBackend) Issue(ctx context.Context, id string) (LinearIssue, error) {
	var response struct {
		Issue *struct {
			ID         string `json:"id"`
			Identifier string `json:"identifier"`
			Title      string `json:"title"`
			State      *struct {
				Name string `json:"name"`
			} `json:"state"`
			Team *struct {
				Key string `json:"key"`
			} `json:"team"`
			Project *struct {
				ID string `json:"id"`
			} `json:"project"`
		} `json:"issue"`
	}
	if err := b.query(ctx, linearIssueQuery, map[string]any{"id": id}, &response); err != nil || response.Issue == nil || response.Issue.ID == "" || response.Issue.Identifier == "" || response.Issue.Title == "" || response.Issue.State == nil || response.Issue.State.Name == "" || response.Issue.Team == nil || response.Issue.Team.Key == "" {
		return LinearIssue{}, errors.New("invalid Linear issue response")
	}
	result := LinearIssue{ID: response.Issue.ID, TeamKey: response.Issue.Team.Key, Identifier: response.Issue.Identifier, Title: response.Issue.Title, State: response.Issue.State.Name}
	if response.Issue.Project != nil {
		result.ProjectID = response.Issue.Project.ID
	}
	return result, nil
}
