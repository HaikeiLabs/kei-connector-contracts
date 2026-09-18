package providers

import (
	"context"
	"errors"
	"net/http"

	"github.com/HaikeiLabs/kei-connector-contracts"
)

// CredentialResolver resolves an opaque reference for one governed
// invocation. Implementations must not log or return the resolved secret.
type CredentialResolver interface {
	Resolve(ctx context.Context, tenantID, workspaceID, subject, ref string) (string, error)
}

// RuntimeConfig contains the only injectable seams of authenticated provider
// backends. The resolver is called after the invocation has passed Guard.
type RuntimeConfig struct {
	HTTPClient  *http.Client
	Credentials CredentialResolver
}

type rejectingCredentialResolver struct{}

func (rejectingCredentialResolver) Resolve(context.Context, string, string, string, string) (string, error) {
	return "", errors.New("credential resolver is not configured")
}

func normalizeRuntimeConfig(config RuntimeConfig) RuntimeConfig {
	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{}
	}
	if config.Credentials == nil {
		config.Credentials = rejectingCredentialResolver{}
	}
	return config
}

type invocationScope struct {
	tenantID    string
	workspaceID string
	subject     string
	credential  string
}

type invocationScopeKey struct{}

func withInvocationScope(ctx context.Context, meta connectors.Metadata, inv connectors.Invocation) context.Context {
	return context.WithValue(ctx, invocationScopeKey{}, invocationScope{
		tenantID:    inv.TenantID,
		workspaceID: inv.WorkspaceID,
		subject:     inv.Subject,
		credential:  meta.CredentialRef,
	})
}

func scopeFromContext(ctx context.Context) (invocationScope, error) {
	scope, ok := ctx.Value(invocationScopeKey{}).(invocationScope)
	if !ok || scope.tenantID == "" || scope.workspaceID == "" || scope.subject == "" || scope.credential == "" {
		return invocationScope{}, errors.New("invocation scope is required")
	}
	if err := connectors.ValidateCredentialRef(scope.credential); err != nil {
		return invocationScope{}, errors.New("credential reference is invalid")
	}
	return scope, nil
}
