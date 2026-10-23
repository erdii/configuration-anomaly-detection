package checks

import (
	"context"
	"errors"
	"testing"

	"github.com/openshift/configuration-anomaly-detection/pkg/investigations/investigation"
)

func TestComposeRunsInDependencyOrder(t *testing.T) {
	type A struct{ V string }
	type B struct{ V string }
	keyA := NewKey[A]("a")
	keyB := NewKey[B]("b")

	var order []string

	checkA := Define("produce_a").
		Produces(keyA.ID()).
		Run(func(_ context.Context, _ *investigation.Resources, f *Facts) error {
			order = append(order, "a")
			Set(f, keyA, A{V: "from_a"})
			return nil
		})

	checkB := Define("produce_b").
		Requires(keyA.ID()).
		Produces(keyB.ID()).
		Run(func(_ context.Context, _ *investigation.Resources, f *Facts) error {
			order = append(order, "b")
			a := MustGet(f, keyA)
			Set(f, keyB, B{V: "from_b_saw_" + a.V})
			return nil
		})

	// Pass B before A — Compose should still run A first.
	f, err := Compose(context.Background(), nil, checkB, checkA)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(order) != 2 || order[0] != "a" || order[1] != "b" {
		t.Fatalf("expected order [a, b], got %v", order)
	}

	b := MustGet(f, keyB)
	if b.V != "from_b_saw_from_a" {
		t.Fatalf("expected B to see A's value, got %q", b.V)
	}
}

func TestComposeOptionalDependency(t *testing.T) {
	type A struct{ V string }
	type B struct{ V string }
	keyA := NewKey[A]("a")
	keyB := NewKey[B]("b")

	var order []string

	checkA := Define("produce_a").
		Produces(keyA.ID()).
		Run(func(_ context.Context, _ *investigation.Resources, f *Facts) error {
			order = append(order, "a")
			Set(f, keyA, A{V: "hello"})
			return nil
		})

	checkB := Define("produce_b").
		Optional(keyA.ID()).
		Produces(keyB.ID()).
		Run(func(_ context.Context, _ *investigation.Resources, f *Facts) error {
			order = append(order, "b")
			v := "no_a"
			if a, ok := Get(f, keyA); ok {
				v = "saw_" + a.V
			}
			Set(f, keyB, B{V: v})
			return nil
		})

	// With both: A runs first due to optional ordering.
	f, err := Compose(context.Background(), nil, checkB, checkA)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if order[0] != "a" {
		t.Fatalf("expected A first, got %v", order)
	}
	b := MustGet(f, keyB)
	if b.V != "saw_hello" {
		t.Fatalf("expected B to see A, got %q", b.V)
	}

	// Without A: B should still run fine.
	order = nil
	f, err = Compose(context.Background(), nil, checkB)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order) != 1 || order[0] != "b" {
		t.Fatalf("expected only b, got %v", order)
	}
	b = MustGet(f, keyB)
	if b.V != "no_a" {
		t.Fatalf("expected B without A, got %q", b.V)
	}
}

func TestComposeMissingRequiredDependency(t *testing.T) {
	keyMissing := NewKey[struct{}]("missing")

	check := Define("needs_missing").
		Requires(keyMissing.ID()).
		Run(func(_ context.Context, _ *investigation.Resources, _ *Facts) error {
			t.Fatal("should not run")
			return nil
		})

	_, err := Compose(context.Background(), nil, check)
	if err == nil {
		t.Fatal("expected error for missing dependency")
	}

	var depErr *MissingDependencyError
	if !errors.As(err, &depErr) {
		t.Fatalf("expected MissingDependencyError, got %T: %v", err, err)
	}
	if depErr.Check != "needs_missing" || depErr.FactID != "missing" {
		t.Fatalf("unexpected error details: %+v", depErr)
	}
}

func TestComposeDuplicateProducer(t *testing.T) {
	key := NewKey[struct{}]("shared")

	a := Define("a").Produces(key.ID()).Run(func(_ context.Context, _ *investigation.Resources, _ *Facts) error { return nil })
	b := Define("b").Produces(key.ID()).Run(func(_ context.Context, _ *investigation.Resources, _ *Facts) error { return nil })

	_, err := Compose(context.Background(), nil, a, b)
	if err == nil {
		t.Fatal("expected error for duplicate producer")
	}
	if !containsSubstring(err.Error(), "produced by both") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestComposeCycleDetection(t *testing.T) {
	keyA := NewKey[struct{}]("a")
	keyB := NewKey[struct{}]("b")

	a := Define("a").Requires(keyB.ID()).Produces(keyA.ID()).Run(func(_ context.Context, _ *investigation.Resources, _ *Facts) error { return nil })
	b := Define("b").Requires(keyA.ID()).Produces(keyB.ID()).Run(func(_ context.Context, _ *investigation.Resources, _ *Facts) error { return nil })

	_, err := Compose(context.Background(), nil, a, b)
	if err == nil {
		t.Fatal("expected error for cycle")
	}
	if !containsSubstring(err.Error(), "cycle") {
		t.Fatalf("expected cycle error, got: %v", err)
	}
}

func TestComposeCheckError(t *testing.T) {
	checkOK := Define("ok").Run(func(_ context.Context, _ *investigation.Resources, _ *Facts) error {
		return nil
	})
	checkFail := Define("fail").Run(func(_ context.Context, _ *investigation.Resources, _ *Facts) error {
		return errors.New("boom")
	})

	_, err := Compose(context.Background(), nil, checkOK, checkFail)
	if err == nil {
		t.Fatal("expected error from failing check")
	}
	if !containsSubstring(err.Error(), "check fail") {
		t.Fatalf("expected error to mention check name, got: %v", err)
	}
}

func TestComposeDiamondDependency(t *testing.T) {
	type Root struct{}
	type Left struct{}
	type Right struct{}
	type Sink struct{}
	rootKey := NewKey[Root]("root")
	leftKey := NewKey[Left]("left")
	rightKey := NewKey[Right]("right")
	sinkKey := NewKey[Sink]("sink")

	var order []string

	root := Define("root").Produces(rootKey.ID()).Run(func(_ context.Context, _ *investigation.Resources, f *Facts) error {
		order = append(order, "root")
		Set(f, rootKey, Root{})
		return nil
	})
	left := Define("left").Requires(rootKey.ID()).Produces(leftKey.ID()).Run(func(_ context.Context, _ *investigation.Resources, f *Facts) error {
		order = append(order, "left")
		Set(f, leftKey, Left{})
		return nil
	})
	right := Define("right").Requires(rootKey.ID()).Produces(rightKey.ID()).Run(func(_ context.Context, _ *investigation.Resources, f *Facts) error {
		order = append(order, "right")
		Set(f, rightKey, Right{})
		return nil
	})
	sink := Define("sink").Requires(leftKey.ID(), rightKey.ID()).Produces(sinkKey.ID()).Run(func(_ context.Context, _ *investigation.Resources, f *Facts) error {
		order = append(order, "sink")
		Set(f, sinkKey, Sink{})
		return nil
	})

	_, err := Compose(context.Background(), nil, sink, right, left, root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if order[0] != "root" {
		t.Fatalf("expected root first, got %v", order)
	}
	if order[3] != "sink" {
		t.Fatalf("expected sink last, got %v", order)
	}
}

func TestComposeIndependentChecks(t *testing.T) {
	var ran []string

	a := Define("a").Run(func(_ context.Context, _ *investigation.Resources, _ *Facts) error {
		ran = append(ran, "a")
		return nil
	})
	b := Define("b").Run(func(_ context.Context, _ *investigation.Resources, _ *Facts) error {
		ran = append(ran, "b")
		return nil
	})

	_, err := Compose(context.Background(), nil, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ran) != 2 {
		t.Fatalf("expected 2 checks to run, got %d", len(ran))
	}
}

func TestValidateChecks(t *testing.T) {
	keyA := NewKey[struct{}]("a")
	keyB := NewKey[struct{}]("b")

	good := []*CheckDef{
		Define("a").Produces(keyA.ID()).Run(func(_ context.Context, _ *investigation.Resources, _ *Facts) error { return nil }),
		Define("b").Requires(keyA.ID()).Produces(keyB.ID()).Run(func(_ context.Context, _ *investigation.Resources, _ *Facts) error { return nil }),
	}
	if errs := ValidateChecks(good); len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}

	bad := []*CheckDef{
		Define("orphan").Requires(FactID("nonexistent")).Run(func(_ context.Context, _ *investigation.Resources, _ *Facts) error { return nil }),
	}
	if errs := ValidateChecks(bad); len(errs) != 1 {
		t.Fatalf("expected 1 error, got %v", errs)
	}
}

func containsSubstring(s, sub string) bool {
	return len(s) >= len(sub) && searchSubstring(s, sub)
}

func searchSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
