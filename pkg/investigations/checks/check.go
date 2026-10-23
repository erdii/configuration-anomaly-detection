package checks

import (
	"context"
	"fmt"
	"strings"

	"github.com/openshift/configuration-anomaly-detection/pkg/investigations/investigation"
)

// CheckFunc is the function signature for a check's implementation.
type CheckFunc func(ctx context.Context, r *investigation.Resources, f *Facts) error

// CheckDef defines a check with its dependencies and the function to run.
type CheckDef struct {
	name     string
	requires []FactID // hard deps: must be produced before this check runs
	optional []FactID // soft deps: ordered before this check if present, ignored if absent
	produces []FactID
	run      CheckFunc
}

// Name returns the check's name.
func (c *CheckDef) Name() string { return c.name }

// Requires returns the check's hard dependencies.
func (c *CheckDef) Requires() []FactID { return c.requires }

// Optional returns the check's soft dependencies.
func (c *CheckDef) Optional() []FactID { return c.optional }

// Produces returns the facts this check produces.
func (c *CheckDef) Produces() []FactID { return c.produces }

// Builder constructs a CheckDef with a fluent API.
type Builder struct {
	def CheckDef
}

// Define starts building a new check with the given name.
func Define(name string) *Builder {
	return &Builder{def: CheckDef{name: name}}
}

// Requires declares hard dependencies: facts that must be produced by another
// check before this one runs. Compose will return an error if no check in the
// set produces a required fact.
func (b *Builder) Requires(ids ...FactID) *Builder {
	b.def.requires = append(b.def.requires, ids...)
	return b
}

// Optional declares soft dependencies: if another check in the set produces
// this fact, it will be ordered before this check. If no check produces it,
// this check runs anyway.
func (b *Builder) Optional(ids ...FactID) *Builder {
	b.def.optional = append(b.def.optional, ids...)
	return b
}

// Produces declares the facts this check writes to the store.
func (b *Builder) Produces(ids ...FactID) *Builder {
	b.def.produces = append(b.def.produces, ids...)
	return b
}

// Run sets the check function and returns the completed CheckDef.
func (b *Builder) Run(fn CheckFunc) *CheckDef {
	b.def.run = fn
	return &b.def
}

// registry holds all registered checks for static validation.
var registry []*CheckDef

// Register adds a check to the global registry and returns it.
// Use at package init time so that TestRegistryDependencies can validate the full graph.
func Register(def *CheckDef) *CheckDef {
	registry = append(registry, def)
	return def
}

// Registry returns a copy of all registered checks.
func Registry() []*CheckDef {
	out := make([]*CheckDef, len(registry))
	copy(out, registry)
	return out
}

// Compose resolves dependencies, topologically sorts the given checks, and
// runs them in order. It validates that all hard dependencies are satisfiable
// before any check executes.
func Compose(ctx context.Context, r *investigation.Resources, defs ...*CheckDef) (*Facts, error) {
	sorted, err := topoSort(defs)
	if err != nil {
		return nil, err
	}

	f := NewFacts()
	for _, def := range sorted {
		if err := def.run(ctx, r, f); err != nil {
			return f, fmt.Errorf("check %s: %w", def.name, err)
		}
	}
	return f, nil
}

// topoSort resolves the execution order for the given checks based on their
// Requires and Optional dependencies. Returns an error if:
//   - a required fact is not produced by any check in the set
//   - there is a dependency cycle
func topoSort(defs []*CheckDef) ([]*CheckDef, error) {
	producers := make(map[FactID]*CheckDef)
	for _, d := range defs {
		for _, p := range d.produces {
			if existing, ok := producers[p]; ok {
				return nil, fmt.Errorf("fact %q produced by both %q and %q", p, existing.name, d.name)
			}
			producers[p] = d
		}
	}

	// Validate hard dependencies are satisfiable.
	for _, d := range defs {
		for _, req := range d.requires {
			if _, ok := producers[req]; !ok {
				return nil, &MissingDependencyError{
					Check:    d.name,
					FactID:   req,
					Required: true,
				}
			}
		}
	}

	// Build adjacency list: edge from dependency → dependent.
	type node struct {
		def     *CheckDef
		inCount int
	}
	nodes := make(map[string]*node, len(defs))
	edges := make(map[string][]string)
	for _, d := range defs {
		nodes[d.name] = &node{def: d}
	}

	addEdge := func(from, to string) {
		if from == to {
			return
		}
		edges[from] = append(edges[from], to)
		nodes[to].inCount++
	}

	for _, d := range defs {
		for _, req := range d.requires {
			addEdge(producers[req].name, d.name)
		}
		for _, opt := range d.optional {
			if p, ok := producers[opt]; ok {
				addEdge(p.name, d.name)
			}
		}
	}

	// Kahn's algorithm.
	var queue []string
	for name, n := range nodes {
		if n.inCount == 0 {
			queue = append(queue, name)
		}
	}

	sorted := make([]*CheckDef, 0, len(defs))
	for len(queue) > 0 {
		name := queue[0]
		queue = queue[1:]
		sorted = append(sorted, nodes[name].def)

		for _, dep := range edges[name] {
			nodes[dep].inCount--
			if nodes[dep].inCount == 0 {
				queue = append(queue, dep)
			}
		}
	}

	if len(sorted) != len(defs) {
		// Find the cycle participants for a useful error message.
		var cycled []string
		for name, n := range nodes {
			if n.inCount > 0 {
				cycled = append(cycled, name)
			}
		}
		return nil, fmt.Errorf("dependency cycle involving checks: %s", strings.Join(cycled, ", "))
	}

	return sorted, nil
}

// MissingDependencyError is returned when a check requires a fact that no
// check in the Compose set produces.
type MissingDependencyError struct {
	Check    string
	FactID   FactID
	Required bool
}

func (e *MissingDependencyError) Error() string {
	return fmt.Sprintf("check %q requires fact %q but no check in the set produces it", e.Check, e.FactID)
}

// ValidateRegistry checks that all hard dependencies across all registered
// checks are satisfiable and that no two checks produce the same fact.
// Returns a list of validation errors (empty means valid).
func ValidateRegistry() []error {
	return ValidateChecks(registry)
}

// ValidateChecks checks dependency integrity for an arbitrary set of checks.
func ValidateChecks(defs []*CheckDef) []error {
	var errs []error
	produced := make(map[FactID]string)

	for _, d := range defs {
		for _, p := range d.produces {
			if existing, ok := produced[p]; ok {
				errs = append(errs, fmt.Errorf("fact %q produced by both %q and %q", p, existing, d.name))
			}
			produced[p] = d.name
		}
	}

	for _, d := range defs {
		for _, req := range d.requires {
			if _, ok := produced[req]; !ok {
				errs = append(errs, &MissingDependencyError{
					Check:    d.name,
					FactID:   req,
					Required: true,
				})
			}
		}
	}

	return errs
}
