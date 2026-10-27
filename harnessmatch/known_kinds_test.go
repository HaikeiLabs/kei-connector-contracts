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

package harnessmatch

import "testing"

func TestEvaluateWithTakesTheKindSet(t *testing.T) {
	call := Call{Kind: "example", Tool: "example.read"}
	if got := Evaluate(call, nil); got.Reason != ReasonUnsupportedKind {
		t.Fatalf("Evaluate with the built-in kinds = %+v", got)
	}
	set := KindSet{Known: func(kind string) bool { return kind == "example" || kind == "viapdp" }, ViaPDP: func(kind string) bool { return kind == "viapdp" }}
	if got := EvaluateWith(call, nil, set); got.Reason == ReasonUnsupportedKind || got.Outcome != OutcomeUnmatched {
		t.Fatalf("EvaluateWith a set that lists it = %+v", got)
	}
	if got := EvaluateWith(Call{Kind: "viapdp"}, nil, set); got.Reason != ReasonEvaluateViaPDP {
		t.Fatalf("EvaluateWith a kind the set decides via the PDP = %+v", got)
	}
	if got := EvaluateWith(Call{Kind: "codex"}, nil, set); got.Reason != ReasonUnsupportedKind {
		t.Fatalf("EvaluateWith a set that omits codex = %+v", got)
	}
	// A set with no ViaPDP keeps the built-in meaning (custom and pi via the PDP).
	builtin := KindSet{Known: ValidKind}
	if got := EvaluateWith(Call{Kind: "pi"}, nil, builtin); got.Reason != ReasonEvaluateViaPDP {
		t.Fatalf("built-in decides_via for pi = %+v", got)
	}
}
