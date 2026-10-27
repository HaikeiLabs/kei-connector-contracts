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

package policy

import (
	"errors"
	"fmt"
	"regexp"
)

// capPattern is the shape of a capability name, the same one the connector
// access schema uses.
var capPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:/-]{0,254}$`)

// TypedTarget is the structured form of a target pattern: the kind and value
// that ParseTarget returns, plus optional capabilities. It is the wire shape
// (schemas/policy-target.v1.schema.json) a policy API carries in place of a
// pattern string. Caps narrow a connector target to the named capabilities and
// are not part of the pattern.
type TypedTarget struct {
	Kind  TargetKind `json:"kind"`
	Value string     `json:"value,omitempty"`
	Caps  []string   `json:"caps,omitempty"`
}

// TypedTargetFromPattern parses a target pattern into its typed form.
func TypedTargetFromPattern(pattern string) (TypedTarget, error) {
	t, err := ParseTarget(pattern)
	if err != nil {
		return TypedTarget{}, err
	}
	return TypedTarget{Kind: t.Kind, Value: t.Value}, nil
}

// Validate reports whether t is a well-formed target: a known kind, a value
// for every kind but "any" (which takes none), and caps only on a connector
// target, each a distinct, well-formed capability name.
func (t TypedTarget) Validate() error {
	if t.Kind == "any" {
		if t.Value != "" || len(t.Caps) > 0 {
			return errors.New(`target kind "any" takes no value or caps`)
		}
		return nil
	}
	if _, known := targetSchemes[string(t.Kind)]; !known {
		return fmt.Errorf("unknown target kind %q", t.Kind)
	}
	if t.Value == "" {
		return fmt.Errorf("target kind %q has an empty value", t.Kind)
	}
	if len(t.Caps) > 0 && t.Kind != "connector" {
		return fmt.Errorf("caps apply only to a connector target, not %q", t.Kind)
	}
	seen := make(map[string]bool, len(t.Caps))
	for _, c := range t.Caps {
		if !capPattern.MatchString(c) {
			return fmt.Errorf("malformed capability %q", c)
		}
		if seen[c] {
			return fmt.Errorf("duplicate capability %q", c)
		}
		seen[c] = true
	}
	return nil
}

// Pattern renders t as the target pattern ParseTarget accepts. Caps are not
// part of the pattern.
func (t TypedTarget) Pattern() string {
	if t.Kind == "any" {
		return "*"
	}
	return string(t.Kind) + ":" + t.Value
}

// TypedPrincipal is the structured form of a principal pattern
// (schemas/policy-principal.v1.schema.json). A policy names a standing source,
// a group or an autogroup. A direct subject (user, email or agent) is valid
// only as the principal of an access grant.
type TypedPrincipal struct {
	Kind  PrincipalKind `json:"kind"`
	Value string        `json:"value"`
}

// TypedPrincipalFromPattern parses a principal pattern into its typed form.
// It does not apply a context: use ValidatePolicy or ValidateGrant for that.
func TypedPrincipalFromPattern(pattern string) (TypedPrincipal, error) {
	p, err := ParsePrincipal(pattern)
	if err != nil {
		return TypedPrincipal{}, err
	}
	return TypedPrincipal{Kind: p.Kind, Value: p.Value}, nil
}

// ValidatePolicy reports whether p may be the principal of a policy: a group,
// or an autogroup with a known name.
func (p TypedPrincipal) ValidatePolicy() error {
	switch p.Kind {
	case "group", "autogroup":
		return p.validateValue()
	}
	return fmt.Errorf("principal kind %q is not a standing source; a policy takes a group or an autogroup", p.Kind)
}

// ValidateGrant reports whether p may be the principal of an access grant: a
// policy principal, or a direct subject (user, email or agent).
func (p TypedPrincipal) ValidateGrant() error {
	switch p.Kind {
	case "group", "autogroup", "user", "email", "agent":
		return p.validateValue()
	}
	return fmt.Errorf("principal kind %q cannot be an access grant principal", p.Kind)
}

func (p TypedPrincipal) validateValue() error {
	if p.Value == "" {
		return fmt.Errorf("principal kind %q has an empty value", p.Kind)
	}
	if p.Kind == "autogroup" && !autogroups[p.Value] {
		return fmt.Errorf("unknown autogroup %q; valid names: member, admin, owner, agent", p.Value)
	}
	return nil
}

// Pattern renders p as the principal pattern ParsePrincipal accepts; "any"
// renders as "*".
func (p TypedPrincipal) Pattern() string {
	if p.Kind == "any" {
		return "*"
	}
	return string(p.Kind) + ":" + p.Value
}
