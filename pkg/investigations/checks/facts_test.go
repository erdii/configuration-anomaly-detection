package checks

import (
	"testing"
)

func TestSetAndGet(t *testing.T) {
	type MyFact struct{ Value int }
	key := NewKey[MyFact]("my_fact")
	f := NewFacts()

	if _, ok := Get(f, key); ok {
		t.Fatal("expected fact to be absent before Set")
	}

	Set(f, key, MyFact{Value: 42})

	got, ok := Get(f, key)
	if !ok {
		t.Fatal("expected fact to be present after Set")
	}
	if got.Value != 42 {
		t.Fatalf("expected Value=42, got %d", got.Value)
	}
}

func TestMustGetPanicsWhenMissing(t *testing.T) {
	type Missing struct{}
	key := NewKey[Missing]("missing")
	f := NewFacts()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected MustGet to panic on missing fact")
		}
	}()
	MustGet(f, key)
}

func TestMustGetReturnsValue(t *testing.T) {
	type Present struct{ X string }
	key := NewKey[Present]("present")
	f := NewFacts()
	Set(f, key, Present{X: "hello"})

	got := MustGet(f, key)
	if got.X != "hello" {
		t.Fatalf("expected X=hello, got %s", got.X)
	}
}

func TestHas(t *testing.T) {
	type SomeFact struct{}
	key := NewKey[SomeFact]("some")
	f := NewFacts()

	if f.Has(key.ID()) {
		t.Fatal("expected Has=false before Set")
	}

	Set(f, key, SomeFact{})

	if !f.Has(key.ID()) {
		t.Fatal("expected Has=true after Set")
	}
}

func TestIDs(t *testing.T) {
	type A struct{}
	type B struct{}
	keyA := NewKey[A]("a")
	keyB := NewKey[B]("b")
	f := NewFacts()

	Set(f, keyA, A{})
	Set(f, keyB, B{})

	ids := f.IDs()
	if len(ids) != 2 {
		t.Fatalf("expected 2 IDs, got %d", len(ids))
	}

	found := map[FactID]bool{}
	for _, id := range ids {
		found[id] = true
	}
	if !found["a"] || !found["b"] {
		t.Fatalf("expected IDs a and b, got %v", ids)
	}
}

func TestDistinctTypesWithSameStructure(t *testing.T) {
	type NodeList struct{ Count int }
	type PodList struct{ Count int }
	nodeKey := NewKey[NodeList]("nodes")
	podKey := NewKey[PodList]("pods")
	f := NewFacts()

	Set(f, nodeKey, NodeList{Count: 3})
	Set(f, podKey, PodList{Count: 7})

	nodes := MustGet(f, nodeKey)
	pods := MustGet(f, podKey)
	if nodes.Count != 3 {
		t.Fatalf("expected nodes.Count=3, got %d", nodes.Count)
	}
	if pods.Count != 7 {
		t.Fatalf("expected pods.Count=7, got %d", pods.Count)
	}
}
