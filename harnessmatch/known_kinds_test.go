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
	known := func(kind string) bool { return kind == "example" }
	if got := EvaluateWith(call, nil, known); got.Reason == ReasonUnsupportedKind || got.Outcome != OutcomeUnmatched {
		t.Fatalf("EvaluateWith a set that lists it = %+v", got)
	}
	if got := EvaluateWith(Call{Kind: "codex"}, nil, known); got.Reason != ReasonUnsupportedKind {
		t.Fatalf("EvaluateWith a set that omits codex = %+v", got)
	}
}
