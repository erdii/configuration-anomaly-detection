package micro_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	amv1 "github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	"github.com/openshift/configuration-anomaly-detection/pkg/executor"
	"github.com/openshift/configuration-anomaly-detection/pkg/logging"
	"github.com/openshift/configuration-anomaly-detection/pkg/micro"
	"github.com/openshift/configuration-anomaly-detection/pkg/notewriter"
	"github.com/openshift/configuration-anomaly-detection/pkg/types"
	ocmmock "github.com/openshift/configuration-anomaly-detection/pkg/ocm/mock"
	"go.uber.org/mock/gomock"
)

func hasActionType(actions []types.Action, actionType string) bool {
	for _, action := range actions {
		if action.Type() == actionType {
			return true
		}
	}
	return false
}

var _ = Describe("CheckBan", func() {
	var notes *notewriter.NoteWriter

	BeforeEach(func() {
		logging.InitLogger("fatal", "", "")
		notes = notewriter.New("test", logging.RawLogger)
	})

	When("user is not banned", func() {
		It("should return not Done", func() {
			account, err := amv1.NewAccount().ID("user-1").Banned(false).Build()
			Expect(err).NotTo(HaveOccurred())

			result := micro.CheckBan(account, notes)
			Expect(result.Done).To(BeFalse())
			Expect(result.Actions).To(BeEmpty())
			Expect(notes.String()).To(ContainSubstring("not banned"))
		})
	})

	When("user is banned due to export control", func() {
		It("should return Done with Escalate and ban info", func() {
			account, err := amv1.NewAccount().
				ID("user-1").
				Banned(true).
				BanCode("export_control_compliance").
				BanDescription("Export control violation").
				Build()
			Expect(err).NotTo(HaveOccurred())

			result := micro.CheckBan(account, notes)
			Expect(result.Done).To(BeTrue())
			Expect(result.Value.Code).To(Equal("export_control_compliance"))
			Expect(result.Value.Description).To(Equal("Export control violation"))
			Expect(hasActionType(result.Actions, string(executor.ActionTypeEscalateIncident))).To(BeTrue())
			Expect(hasActionType(result.Actions, string(executor.ActionTypeSilenceIncident))).To(BeFalse())
			Expect(hasActionType(result.Actions, string(executor.ActionTypeServiceLog))).To(BeFalse())
			Expect(notes.String()).To(ContainSubstring("export control"))
		})
	})

	When("user is banned for other reasons", func() {
		It("should return Done with ServiceLog and Silence", func() {
			account, err := amv1.NewAccount().
				ID("user-1").
				Banned(true).
				BanCode("terms_and_conditions").
				BanDescription("Terms violation").
				Build()
			Expect(err).NotTo(HaveOccurred())

			result := micro.CheckBan(account, notes)
			Expect(result.Done).To(BeTrue())
			Expect(result.Value.Code).To(Equal("terms_and_conditions"))
			Expect(result.Value.Description).To(Equal("Terms violation"))
			Expect(hasActionType(result.Actions, string(executor.ActionTypeServiceLog))).To(BeTrue())
			Expect(hasActionType(result.Actions, string(executor.ActionTypeSilenceIncident))).To(BeTrue())
			Expect(hasActionType(result.Actions, string(executor.ActionTypeEscalateIncident))).To(BeFalse())
			Expect(notes.String()).To(ContainSubstring("terms_and_conditions"))
			Expect(notes.String()).To(ContainSubstring("Service Log"))
		})
	})
})

var _ = Describe("BannedUser", func() {
	var (
		mockCtrl  *gomock.Controller
		ocmClient *ocmmock.MockClient
		cluster   *cmv1.Cluster
		notes     *notewriter.NoteWriter
	)

	BeforeEach(func() {
		logging.InitLogger("fatal", "", "")
		mockCtrl = gomock.NewController(GinkgoT())
		ocmClient = ocmmock.NewMockClient(mockCtrl)
		var err error
		cluster, err = cmv1.NewCluster().ID("test-id").Build()
		Expect(err).NotTo(HaveOccurred())
		notes = notewriter.New("test", logging.RawLogger)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	When("user is banned due to export control", func() {
		It("should return Done with Escalate", func() {
			account, _ := amv1.NewAccount().ID("user-1").
				Banned(true).BanCode("export_control_compliance").
				BanDescription("Export control").Build()
			ocmClient.EXPECT().GetCreatorFromCluster(gomock.Any()).Return(account, nil)

			result, err := micro.BannedUser(ocmClient, cluster, notes)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Done).To(BeTrue())
			Expect(hasActionType(result.Actions, string(executor.ActionTypeEscalateIncident))).To(BeTrue())
			Expect(hasActionType(result.Actions, string(executor.ActionTypeSilenceIncident))).To(BeFalse())
		})
	})

	When("user is banned for other reasons", func() {
		It("should return Done with ServiceLog and Silence", func() {
			account, _ := amv1.NewAccount().ID("user-1").
				Banned(true).BanCode("terms_and_conditions").
				BanDescription("Terms violation").Build()
			ocmClient.EXPECT().GetCreatorFromCluster(gomock.Any()).Return(account, nil)

			result, err := micro.BannedUser(ocmClient, cluster, notes)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Done).To(BeTrue())
			Expect(hasActionType(result.Actions, string(executor.ActionTypeServiceLog))).To(BeTrue())
			Expect(hasActionType(result.Actions, string(executor.ActionTypeSilenceIncident))).To(BeTrue())
		})
	})

	When("GetCreatorFromCluster fails", func() {
		It("should return error", func() {
			ocmClient.EXPECT().GetCreatorFromCluster(gomock.Any()).Return(nil, fmt.Errorf("connection timeout"))

			_, err := micro.BannedUser(ocmClient, cluster, notes)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("connection timeout"))
		})
	})

	When("user is not banned", func() {
		It("should return not Done", func() {
			account, _ := amv1.NewAccount().ID("user-1").Banned(false).Build()
			ocmClient.EXPECT().GetCreatorFromCluster(gomock.Any()).Return(account, nil)

			result, err := micro.BannedUser(ocmClient, cluster, notes)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Done).To(BeFalse())
			Expect(notes.String()).To(ContainSubstring("not banned"))
		})
	})
})
