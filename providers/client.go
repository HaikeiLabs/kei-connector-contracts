// Package providers contains the initial read-only connector clients for the
// governed providers defined by the PR #312 contract. Every client is a seam:
// the backend interface is the integration point, and the bundled memory
// backends are the only implementations, so no provider network call can
// occur. Real backends will replace the memory ones behind the same seams.
//
// CRM customer reads are intentionally absent: the contract only defines
// lead.* capabilities for the CRM provider, and new operations must be added
// to the contract before a client may expose them.
package providers

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/HaikeiLabs/kei-connector-contracts"
)

// Client is the provider-neutral seam invoked by the control plane.
type Client interface {
	Provider() connectors.Provider
	Invoke(ctx context.Context, meta connectors.Metadata, inv connectors.Invocation, payload Payload) (Result, error)
}

// Payload is a provider-specific, read-only request. Each payload declares
// the capability it implements so an invocation cannot be replayed against a
// different operation.
type Payload interface {
	Capability() string
}

// Result is the provider-neutral envelope returned by every client.
type Result struct {
	Capability string `json:"capability"`
	Data       any    `json:"data"`
}

// Guard is the fail-closed gate every client runs before touching a backend.
// It applies the contract's structural boundaries (tenant/workspace/status/
// capability), rejects path-traversal resources that would normalize outside
// the resource boundary, and keeps the destructive-operation lock. Resource,
// prefix, and action policy are decided by the ABAC policy layer, not here.
func Guard(meta connectors.Metadata, inv connectors.Invocation) error {
	if err := connectors.ValidateCall(meta, inv); err != nil {
		return err
	}
	if !traversalSafe(inv.Resource) {
		return errors.New("resource contains a path traversal segment")
	}
	if err := DestructiveAllowed(meta, inv); err != nil {
		return err
	}
	return nil
}

// DestructiveAllowed keeps delete-class operations disabled unless the
// connector policy explicitly enables them and an approval reference is
// present. The initial capability surface defines no delete capabilities, so
// this is defense in depth for the contract's fail-closed default.
func DestructiveAllowed(meta connectors.Metadata, inv connectors.Invocation) error {
	if inv.Action != connectors.ActionDelete {
		return nil
	}
	if !meta.Policy.DestructiveEnabled {
		return errors.New("destructive operations are disabled")
	}
	if inv.ApprovalID == "" {
		return errors.New("approval required for destructive operations")
	}
	return nil
}

// traversalSafe reports whether resource contains no path-traversal segment.
// Reject any ".." segment, literally or percent-encoded, so a resource cannot
// escape its sanctioned scope.
func traversalSafe(resource string) bool {
	if hasDotDotSegment(resource) {
		return false
	}
	if decoded, err := url.PathUnescape(resource); err == nil && decoded != resource && hasDotDotSegment(decoded) {
		return false
	}
	return true
}

func hasDotDotSegment(path string) bool {
	for _, segment := range strings.Split(path, "/") {
		if segment == ".." {
			return true
		}
	}
	return false
}

func checkProvider(meta connectors.Metadata, provider connectors.Provider) error {
	if meta.Provider != provider {
		return fmt.Errorf("connector provider %q does not match client %q", meta.Provider, provider)
	}
	return nil
}

func matchCapability(inv connectors.Invocation, payload Payload) error {
	if payload.Capability() != inv.Capability {
		return fmt.Errorf("payload capability %q does not match invocation capability %q", payload.Capability(), inv.Capability)
	}
	return nil
}

// New returns a read-only client. Linear is wired to its authenticated
// GraphQL runtime; credential resolution must be injected by the runtime
// integration before an invocation can succeed.
func New(provider connectors.Provider) (Client, error) {
	switch provider {
	case connectors.ProviderGoogle:
		return NewAuthenticatedGoogle(RuntimeConfig{}), nil
	case connectors.ProviderLinear:
		return NewLinearRuntime(HTTPLinearGraphQLTransport{}, nil), nil
	case connectors.ProviderGitHub:
		return NewAuthenticatedGitHub(RuntimeConfig{}), nil
	case connectors.ProviderS3:
		return NewS3(MemoryS3{}), nil
	case connectors.ProviderNotion:
		return NewNotion(MemoryNotion{}), nil
	case connectors.ProviderCRM:
		return NewCRM(MemoryCRM{}), nil
	default:
		return nil, fmt.Errorf("unsupported provider %q", provider)
	}
}
