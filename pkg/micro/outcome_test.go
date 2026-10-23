package micro_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift/configuration-anomaly-detection/pkg/executor"
	"github.com/openshift/configuration-anomaly-detection/pkg/logging"
	"github.com/openshift/configuration-anomaly-detection/pkg/micro"
	"github.com/openshift/configuration-anomaly-detection/pkg/notewriter"
	"github.com/openshift/configuration-anomaly-detection/pkg/testutil"
)

var _ = Describe("Outcome", func() {
	var (
		notes *notewriter.NoteWriter
		o     *micro.Outcome
	)

	BeforeEach(func() {
		logging.InitLogger("fatal", "", "")
		notes = notewriter.New("test", logging.RawLogger)
		o = micro.NewOutcome(notes, "cluster-1", "test-investigation")
	})

	Describe("Stop", func() {
		It("should return actions passed to it", func() {
			result := o.Stop(executor.Escalate("something broke"))
			Expect(testutil.HasActionType(result.Actions, string(executor.ActionTypeEscalateIncident))).To(BeTrue())
		})

		It("should flush previously collected actions", func() {
			ohss := executor.NewOHSSCaseAction("TLS issue").
				ForCluster("cluster-1").
				WithPriority("Major").
				Build()
			o.Collect(ohss)

			result := o.Stop(executor.Escalate("early exit"))
			Expect(testutil.HasActionType(result.Actions, string(executor.ActionTypeOHSSCase))).To(BeTrue())
			Expect(testutil.HasActionType(result.Actions, string(executor.ActionTypeEscalateIncident))).To(BeTrue())
		})

		It("should work with no collected actions", func() {
			result := o.Stop(executor.Silence("all clear"))
			Expect(testutil.HasActionType(result.Actions, string(executor.ActionTypeSilenceIncident))).To(BeTrue())
		})
	})

	Describe("Finalize", func() {
		It("should combine collected and final actions", func() {
			ohss := executor.NewOHSSCaseAction("TLS issue").
				ForCluster("cluster-1").
				WithPriority("Major").
				Build()
			o.Collect(ohss)

			result := o.Finalize(executor.Escalate("done"))
			Expect(testutil.HasActionType(result.Actions, string(executor.ActionTypeOHSSCase))).To(BeTrue())
			Expect(testutil.HasActionType(result.Actions, string(executor.ActionTypeEscalateIncident))).To(BeTrue())
		})

		It("should work with no collected actions", func() {
			result := o.Finalize(executor.Escalate("done"))
			Expect(testutil.HasActionType(result.Actions, string(executor.ActionTypeEscalateIncident))).To(BeTrue())
			Expect(testutil.HasActionType(result.Actions, string(executor.ActionTypeOHSSCase))).To(BeFalse())
		})
	})

	Describe("Collect", func() {
		It("should accumulate multiple calls", func() {
			ohss1 := executor.NewOHSSCaseAction("issue 1").ForCluster("cluster-1").Build()
			ohss2 := executor.NewOHSSCaseAction("issue 2").ForCluster("cluster-1").Build()
			o.Collect(ohss1)
			o.Collect(ohss2)

			result := o.Finalize(executor.Escalate("done"))
			ohssCount := 0
			for _, a := range result.Actions {
				if a.Type() == string(executor.ActionTypeOHSSCase) {
					ohssCount++
				}
			}
			Expect(ohssCount).To(Equal(2))
		})
	})
})
