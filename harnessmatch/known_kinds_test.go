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
