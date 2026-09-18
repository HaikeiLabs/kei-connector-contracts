package providers

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/HaikeiLabs/kei-connector-contracts"
)

type linearTestResolver struct {
	token string
	ref   string
	scope LinearCredentialScope
	err   error
}

func (r *linearTestResolver) Resolve(_ context.Context, ref string, scope LinearCredentialScope) (string, error) {
	r.ref, r.scope = ref, scope
	return r.token, r.err
}

type linearTestTransport struct {
	endpoint string
	token    string
	query    string
	vars     map[string]any
	data     json.RawMessage
	err      error
}

func (t *linearTestTransport) Query(_ context.Context, endpoint, token, query string, variables map[string]any) (json.RawMessage, error) {
	t.endpoint, t.token, t.query, t.vars = endpoint, token, query, variables
	return t.data, t.err
}

func TestLinearRuntimeAuthenticatesWithInvocationScope(t *testing.T) {
	resolver := &linearTestResolver{token: "linear-access-token"}
	transport := &linearTestTransport{data: json.RawMessage(`{"team":{"key":"KEI","name":"Kei Engineering"}}`)}
	client := NewLinearRuntime(transport, resolver)
	meta := metaFor(t, connectors.ProviderLinear, []string{"linear/team/KEI"}, nil)
	inv := invocation("team.read", connectors.ActionRead, "linear/team/KEI")
	inv.Subject = "subject-7"

	result, err := client.Invoke(context.Background(), meta, inv, TeamReadPayload{})
	if err != nil {
		t.Fatalf("Invoke: %v", err)
	}
	if got := result.Data.(LinearTeam).Name; got != "Kei Engineering" {
		t.Fatalf("team name = %q", got)
	}
	if resolver.ref != meta.CredentialRef || resolver.scope != (LinearCredentialScope{TenantID: "t-1", WorkspaceID: "w-1", Subject: "subject-7"}) {
		t.Fatalf("resolver binding = ref %q scope %+v", resolver.ref, resolver.scope)
	}
	if transport.endpoint != linearGraphQLEndpoint || transport.token != "linear-access-token" || transport.vars["key"] != "KEI" || !strings.Contains(transport.query, "team(key: $key)") {
		t.Fatalf("GraphQL request = endpoint %q token %q query %q vars %#v", transport.endpoint, transport.token, transport.query, transport.vars)
	}
}

func TestLinearRuntimeFailsClosedWithoutDependencies(t *testing.T) {
	client := NewLinearRuntime(nil, nil)
	meta := metaFor(t, connectors.ProviderLinear, []string{"linear/team/KEI"}, nil)
	_, err := client.Invoke(context.Background(), meta, invocation("team.read", connectors.ActionRead, "linear/team/KEI"), TeamReadPayload{})
	if err == nil || err.Error() != "linear runtime dependencies are unavailable" {
		t.Fatalf("error = %v", err)
	}
}

func TestLinearRuntimeRejectsMalformedResponseWithoutLeakingToken(t *testing.T) {
	secret := "linear-secret-value"
	resolver := &linearTestResolver{token: secret}
	transport := &linearTestTransport{data: json.RawMessage(`{"team":{"key":"KEI","name":"Kei Engineering","unexpected":"leak"}}`)}
	client := NewLinearRuntime(transport, resolver)
	meta := metaFor(t, connectors.ProviderLinear, []string{"linear/team/KEI"}, nil)
	_, err := client.Invoke(context.Background(), meta, invocation("team.read", connectors.ActionRead, "linear/team/KEI"), TeamReadPayload{})
	if err == nil || strings.Contains(err.Error(), secret) || strings.Contains(err.Error(), "unexpected") {
		t.Fatalf("error = %v", err)
	}
}

func TestLinearRuntimeDoesNotExposeResolverErrors(t *testing.T) {
	secret := "linear-secret-value"
	resolver := &linearTestResolver{token: secret, err: errors.New("backend included " + secret)}
	client := NewLinearRuntime(&linearTestTransport{}, resolver)
	meta := metaFor(t, connectors.ProviderLinear, []string{"linear/team/KEI"}, nil)
	_, err := client.Invoke(context.Background(), meta, invocation("team.read", connectors.ActionRead, "linear/team/KEI"), TeamReadPayload{})
	if err == nil || err.Error() != "linear credential is unavailable" || strings.Contains(err.Error(), secret) {
		t.Fatalf("error = %v", err)
	}
}
