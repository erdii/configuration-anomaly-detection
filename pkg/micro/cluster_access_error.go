package micro

import (
	"github.com/openshift/configuration-anomaly-detection/pkg/executor"
	"github.com/openshift/configuration-anomaly-detection/pkg/investigations/investigation"
	"github.com/openshift/configuration-anomaly-detection/pkg/notewriter"
	"github.com/openshift/configuration-anomaly-detection/pkg/types"
)

// ClusterAccessError checks whether a Build error is a known cluster access
// issue and returns an escalation result if so. Pure — no external calls.
func ClusterAccessError(err error, notes *notewriter.NoteWriter) DataResult[string] {
	msg, ok := investigation.ClusterAccessErrorMessage(err)
	if !ok {
		return DataResult[string]{}
	}

	if notes != nil {
		notes.AppendWarning("%s", msg)
	}
	return DataResult[string]{
		Result: Result{
			Done: true,
			Actions: []types.Action{
				executor.Escalate(msg),
			},
		},
		Value: msg,
	}
}
