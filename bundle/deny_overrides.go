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

package bundle

import "github.com/HaikeiLabs/kei-connector-contracts/harnessmatch"

// evalDenyOverrides is the evaluation id for deny-overrides decisions (ADR-038):
// any matching deny wins over every matching permit, priority is not read, and
// no match denies. v1 bundles declare it as evaluation.order; v2 bundles as an
// entry of evaluation.precedence. The older ids keep their meaning.
const evalDenyOverrides = "deny_overrides"

// denyOverrides reports whether the v2 policy set declares deny-overrides.
func (s v2PolicySet) denyOverrides() bool {
	for _, id := range s.Evaluation.Precedence {
		if id == evalDenyOverrides {
			return true
		}
	}
	return false
}

// denyBeforePermit returns the policies with every deny ahead of every permit,
// each group keeping its array order, so harnessmatch.Evaluate's first match
// is the deny-overrides decision.
func denyBeforePermit(policies []harnessmatch.Policy) []harnessmatch.Policy {
	ordered := make([]harnessmatch.Policy, 0, len(policies))
	for _, p := range policies {
		if p.Action == "deny" {
			ordered = append(ordered, p)
		}
	}
	for _, p := range policies {
		if p.Action != "deny" {
			ordered = append(ordered, p)
		}
	}
	return ordered
}
