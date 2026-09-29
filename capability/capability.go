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

package capability

import (
	"errors"
	"fmt"
	"sort"

	"github.com/HaikeiLabs/kei-connector-contracts/contract"
)

// LookupCapability returns the definition for a named capability on a
// provider, plus whether it is defined. contract.Capability lookup is the gate a
// connector client runs before any call: an operation that is not in the
// provider's deliberate surface is refused.
func LookupCapability(provider contract.Provider, name string) (contract.Capability, bool) {
	for _, c := range contract.CapabilitiesFor(provider) {
		if c.Name == name {
			return c, true
		}
	}
	return contract.Capability{}, false
}

// CapabilityFor returns the defined capability for provider and name, or an
// error naming the unsupported capability.
func CapabilityFor(provider contract.Provider, name string) (contract.Capability, error) {
	c, ok := LookupCapability(provider, name)
	if !ok {
		return contract.Capability{}, fmt.Errorf("capability %q is not defined for provider %q", name, provider)
	}
	return c, nil
}

// CapabilitiesForProvider returns a copy of the provider's defined
// capabilities. An unknown provider is an error (fail-closed) rather than an
// empty set.
func CapabilitiesForProvider(provider contract.Provider) ([]contract.Capability, error) {
	if !contract.ValidProvider(provider) {
		return nil, errors.New("unsupported provider")
	}
	return contract.CapabilitiesFor(provider), nil
}

// Providers returns all provider keys in a stable order for discovery.
func Providers() []contract.Provider {
	out := contract.DefinedProviders()
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}
