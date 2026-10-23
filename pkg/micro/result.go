// Package micro provides reusable building blocks for investigations.
//
// Naming convention:
//   - Pure checks (no external calls): Check* prefix, return without error (e.g. CheckBan, ClusterAccessError)
//   - Fetching functions (call OCM/K8s/etc): noun phrases, return with error (e.g. BannedUser, ClusterCreator, PullSecret)
package micro

import "github.com/openshift/configuration-anomaly-detection/pkg/types"

// Result is the return type for micro functions that produce a terminal
// finding. Done indicates the caller should treat this as a complete
// finding and stop further checks in the current branch.
type Result struct {
	Done    bool
	Actions []types.Action
}

// DataResult extends Result with a typed value for callers that need
// to inspect the fetched data (e.g., ban code, error message).
type DataResult[T any] struct {
	Result
	Value T
}
