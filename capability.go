package connectors

import (
	"errors"
	"fmt"
	"sort"
)

// LookupCapability returns the definition for a named capability on a
// provider, plus whether it is defined. Capability lookup is the gate a
// connector client runs before any call: an operation that is not in the
// provider's deliberate surface is refused.
func LookupCapability(provider Provider, name string) (Capability, bool) {
	for _, c := range definitions[provider] {
		if c.Name == name {
			return c, true
		}
	}
	return Capability{}, false
}

// CapabilityFor returns the defined capability for provider and name, or an
// error naming the unsupported capability.
func CapabilityFor(provider Provider, name string) (Capability, error) {
	c, ok := LookupCapability(provider, name)
	if !ok {
		return Capability{}, fmt.Errorf("capability %q is not defined for provider %q", name, provider)
	}
	return c, nil
}

// CapabilitiesForProvider returns a copy of the provider's defined
// capabilities. An unknown provider is an error (fail-closed) rather than an
// empty set.
func CapabilitiesForProvider(provider Provider) ([]Capability, error) {
	if !validProvider(provider) {
		return nil, errors.New("unsupported provider")
	}
	return CapabilitiesFor(provider), nil
}

// Providers returns all provider keys in a stable order for discovery.
func Providers() []Provider {
	out := make([]Provider, 0, len(definitions))
	for p := range definitions {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}
