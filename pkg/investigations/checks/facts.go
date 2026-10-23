package checks

import (
	"fmt"
	"sync"
)

// FactID identifies a fact in the store. Used for dependency declaration.
type FactID string

// Key is a typed handle for getting and setting facts.
// The type parameter ensures compile-time type safety at the call site.
type Key[T any] struct {
	id FactID
}

// NewKey creates a typed fact key with the given name.
func NewKey[T any](name string) Key[T] {
	return Key[T]{id: FactID(name)}
}

// ID returns the untyped FactID for use in dependency declarations.
func (k Key[T]) ID() FactID {
	return k.id
}

// Facts is a type-safe store that checks read from and write to.
type Facts struct {
	mu    sync.RWMutex
	store map[FactID]any
}

// NewFacts creates an empty fact store.
func NewFacts() *Facts {
	return &Facts{store: make(map[FactID]any)}
}

// Set stores a typed value under the given key.
func Set[T any](f *Facts, k Key[T], val T) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.store[k.id] = val
}

// Get retrieves a typed value. Returns the zero value and false if not set.
func Get[T any](f *Facts, k Key[T]) (T, bool) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	v, ok := f.store[k.id]
	if !ok {
		var zero T
		return zero, false
	}
	return v.(T), true
}

// MustGet retrieves a typed value or panics.
// Use only when the check has declared a hard Requires dependency on the producing check.
func MustGet[T any](f *Facts, k Key[T]) T {
	v, ok := Get(f, k)
	if !ok {
		panic(fmt.Sprintf("checks: required fact %q not set", k.id))
	}
	return v
}

// Has returns whether a fact has been set.
func (f *Facts) Has(id FactID) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	_, ok := f.store[id]
	return ok
}

// IDs returns all fact IDs currently in the store.
func (f *Facts) IDs() []FactID {
	f.mu.RLock()
	defer f.mu.RUnlock()
	ids := make([]FactID, 0, len(f.store))
	for id := range f.store {
		ids = append(ids, id)
	}
	return ids
}
