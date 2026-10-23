package micro

import (
	"fmt"

	amv1 "github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	"github.com/openshift/configuration-anomaly-detection/pkg/executor"
	"github.com/openshift/configuration-anomaly-detection/pkg/notewriter"
	"github.com/openshift/configuration-anomaly-detection/pkg/ocm"
	"github.com/openshift/configuration-anomaly-detection/pkg/types"
)

// BanInfo carries ban details for callers that need to branch on ban code.
type BanInfo struct {
	Code        string
	Description string
}

// CheckBan inspects an already-fetched account for ban status.
// Pure — no external calls.
func CheckBan(account *amv1.Account, notes *notewriter.NoteWriter) DataResult[BanInfo] {
	if !account.Banned() {
		notes.AppendSuccess("User is not banned.")
		return DataResult[BanInfo]{Result: Result{Done: false}}
	}

	code := account.BanCode()
	desc := account.BanDescription()
	info := BanInfo{Code: code, Description: desc}

	if code == "export_control_compliance" {
		notes.AppendWarning("User is banned (export control): %s", desc)
		return DataResult[BanInfo]{
			Result: Result{
				Done: true,
				Actions: []types.Action{
					executor.Escalate("Export Control Compliance ban detected, please refer to the SOP."),
				},
			},
			Value: info,
		}
	}

	sl := ocm.NewOCMBannedUserServiceLog()
	notes.AppendWarning("User is banned (%s): %s\nPlease open a proactive case, so that MCS can resolve the ban or organize an ownership transfer.", code, desc)
	notes.AppendWarning("Sending out Service Log (%s)", sl.Summary)
	return DataResult[BanInfo]{
		Result: Result{
			Done: true,
			Actions: []types.Action{
				executor.NewServiceLogAction(sl.Severity, sl.Summary).
					WithDescription(sl.Description).
					WithServiceName(sl.ServiceName).
					Build(),
				executor.Escalate("Banned OCM user, please open a proactive ticket."),
			},
		},
		Value: info,
	}
}

// BannedUser fetches the cluster creator and checks ban status.
// For callers that need the ban code to branch on (e.g., export control),
// use ClusterCreator + CheckBan instead.
func BannedUser(ocmClient ocm.Client, cluster *cmv1.Cluster, notes *notewriter.NoteWriter) (Result, error) {
	user, err := ocmClient.GetCreatorFromCluster(cluster)
	if err != nil {
		return Result{}, fmt.Errorf("while checking if the cluster owner is banned: %w", err)
	}
	result := CheckBan(user, notes)
	return result.Result, nil
}
