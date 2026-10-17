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
	"sort"
	"strings"
)

// Effect is the decision a rule makes: "permit" or "deny".
type Effect string

// Rule is a single policy rule.
type Rule struct {
	ID        string
	Principal string
	Target    string
	Effect    Effect
	Priority  int
	Enabled   bool
}

// Request is the principal and target a decision is made for.
type Request struct {
	User   string
	Email  string
	Groups []string
	Target string
}

// Outcome is the result of a decision.
type Outcome struct {
	Decision  string
	DecidedBy string
	Outranked []string
}

// Precedence compares two rules: priority DESC, then deny before permit,
// then ID ASC. It returns a negative number if a ranks before b, a positive
// number if b ranks before a, and zero if they are equivalent.
var Precedence func(a, b Rule) int = func(a, b Rule) int {
	if a.Priority != b.Priority {
		if a.Priority > b.Priority {
			return -1
		}
		return 1
	}
	aDeny, bDeny := a.Effect == "deny", b.Effect == "deny"
	if aDeny != bDeny {
		if aDeny {
			return -1
		}
		return 1
	}
	if a.ID != b.ID {
		if a.ID < b.ID {
			return -1
		}
		return 1
	}
	return 0
}

// Decide decides the request against the rules. It considers the enabled
// rules whose principal matches the request and whose target equals the
// request target exactly, or is "*" or "<scheme>:*" with the same scheme,
// orders them by Precedence, and lets the first one decide. DecidedBy is
// the deciding rule's ID and Outranked the IDs of the other matching rules
// in order. If no rule matches, the decision is "no_match".
func Decide(rules []Rule, req Request) Outcome {
	matched := make([]Rule, 0, len(rules))
	for _, r := range rules {
		if !r.Enabled {
			continue
		}
		if !MatchPrincipal(r.Principal, req.User, req.Email, req.Groups) {
			continue
		}
		if !targetMatches(r.Target, req.Target) {
			continue
		}
		matched = append(matched, r)
	}
	if len(matched) == 0 {
		return Outcome{Decision: "no_match"}
	}
	sort.SliceStable(matched, func(i, j int) bool {
		return Precedence(matched[i], matched[j]) < 0
	})
	decided := matched[0]
	out := Outcome{
		Decision:  string(decided.Effect),
		DecidedBy: decided.ID,
		Outranked: make([]string, 0, len(matched)-1),
	}
	for _, r := range matched[1:] {
		out.Outranked = append(out.Outranked, r.ID)
	}
	return out
}

// targetMatches reports whether ruleTarget matches reqTarget: an exact
// equality, the any target "*", or a "<scheme>:*" wildcard with the same
// scheme and a non-empty value in reqTarget.
func targetMatches(ruleTarget, reqTarget string) bool {
	if ruleTarget == reqTarget || ruleTarget == "*" {
		return true
	}
	scheme, ok := strings.CutSuffix(ruleTarget, ":*")
	if !ok || scheme == "" || strings.ContainsAny(scheme, ":*") {
		return false
	}
	reqScheme, reqValue, hasScheme := strings.Cut(reqTarget, ":")
	return hasScheme && reqScheme == scheme && reqValue != ""
}

// MatchTarget reports whether ruleTarget matches reqTarget: an exact
// equality, the any target "*", or a "<scheme>:*" wildcard with the same
// scheme and a non-empty value in reqTarget.
func MatchTarget(ruleTarget, reqTarget string) bool {
	return targetMatches(ruleTarget, reqTarget)
}
