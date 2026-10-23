package upgradeconfigsyncfailureover4hr

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	amv1 "github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	"github.com/openshift/configuration-anomaly-detection/pkg/executor"
	"github.com/openshift/configuration-anomaly-detection/pkg/investigations/investigation"
	"github.com/openshift/configuration-anomaly-detection/pkg/logging"
	"github.com/openshift/configuration-anomaly-detection/pkg/notewriter"
	ocmmock "github.com/openshift/configuration-anomaly-detection/pkg/ocm/mock"
	"github.com/openshift/configuration-anomaly-detection/pkg/testutil"
	"github.com/openshift/configuration-anomaly-detection/pkg/types"
	"go.uber.org/mock/gomock"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/openshift/configuration-anomaly-detection/pkg/pagerduty"
)

func testScheme() *runtime.Scheme {
	s := runtime.NewScheme()
	utilruntime.Must(corev1.AddToScheme(s))
	return s
}

func hasSilenceAction(actions []types.Action) bool {
	return testutil.HasActionType(actions, string(executor.ActionTypeSilenceIncident))
}

func hasEscalateAction(actions []types.Action) bool {
	return testutil.HasActionType(actions, string(executor.ActionTypeEscalateIncident))
}

func hasServiceLogAction(actions []types.Action) bool {
	return testutil.HasActionType(actions, string(executor.ActionTypeServiceLog))
}

// k8sFailBuilder wraps ResourceBuilderMock to fail on the second Build call
// (after WithK8sClient), simulating cluster access errors.
type k8sFailBuilder struct {
	inner    *investigation.ResourceBuilderMock
	k8sErr   error
	buildK8s bool
}

func (b *k8sFailBuilder) WithCluster() investigation.ResourceBuilder { b.inner.WithCluster(); return b }

func (b *k8sFailBuilder) WithClusterDeployment() investigation.ResourceBuilder {
	b.inner.WithClusterDeployment()
	return b
}

func (b *k8sFailBuilder) WithAwsClient() investigation.ResourceBuilder {
	b.inner.WithAwsClient()
	return b
}

func (b *k8sFailBuilder) WithRestConfig() investigation.ResourceBuilder {
	b.inner.WithRestConfig()
	return b
}
func (b *k8sFailBuilder) WithOC() investigation.ResourceBuilder    { b.inner.WithOC(); return b }
func (b *k8sFailBuilder) WithNotes() investigation.ResourceBuilder { b.inner.WithNotes(); return b }
func (b *k8sFailBuilder) WithManagementRestConfig() investigation.ResourceBuilder {
	b.inner.WithManagementRestConfig()
	return b
}

func (b *k8sFailBuilder) WithManagementK8sClient() investigation.ResourceBuilder {
	b.inner.WithManagementK8sClient()
	return b
}

func (b *k8sFailBuilder) WithManagementOCClient() investigation.ResourceBuilder {
	b.inner.WithManagementOCClient()
	return b
}
func (b *k8sFailBuilder) WithPdClient(_ pagerduty.Client) investigation.ResourceBuilder { return b }
func (b *k8sFailBuilder) WithK8sClient() investigation.ResourceBuilder                  { b.buildK8s = true; return b }

func (b *k8sFailBuilder) Build() (*investigation.Resources, error) {
	if b.buildK8s {
		return b.inner.Resources, b.k8sErr
	}
	return b.inner.Build()
}

var _ = Describe("parseMUOLogs", func() {
	var notes *notewriter.NoteWriter

	BeforeEach(func() {
		logging.InitLogger("fatal", "", "")
		notes = notewriter.New("test", logging.RawLogger)
	})

	It("should detect leader election deadlock", func() {
		logs := `2026-10-19T10:00:00Z INFO starting operator
2026-10-19T10:01:00Z ERROR Leader pod has been deleted, waiting for garbage collection
2026-10-19T10:02:00Z INFO retrying`
		parseMUOLogs(logs, "test-cluster-id", notes)
		Expect(notes.String()).To(ContainSubstring("MUO deadlock detected"))
	})

	It("should detect TLS handshake timeout", func() {
		logs := `2026-10-19T10:00:00Z INFO starting sync
2026-10-19T10:01:00Z ERROR net/http: TLS handshake timeout
2026-10-19T10:02:00Z INFO retrying`
		parseMUOLogs(logs, "test-cluster-id", notes)
		Expect(notes.String()).To(ContainSubstring("TLS handshake timeout"))
	})

	It("should detect ocm-config-getter errors", func() {
		logs := `2026-10-19T10:00:00Z ERROR ocm-config-getter failed to retrieve config from OCM`
		parseMUOLogs(logs, "test-cluster-id", notes)
		Expect(notes.String()).To(ContainSubstring("ocm-config-getter"))
	})

	It("should detect multiple patterns", func() {
		logs := `Leader pod has been deleted, waiting for garbage collection
net/http: TLS handshake timeout`
		parseMUOLogs(logs, "test-cluster-id", notes)
		output := notes.String()
		Expect(output).To(ContainSubstring("deadlock"))
		Expect(output).To(ContainSubstring("TLS handshake"))
	})

	It("should report success when no patterns found", func() {
		logs := `2026-10-19T10:00:00Z INFO everything is fine
2026-10-19T10:01:00Z INFO sync complete`
		parseMUOLogs(logs, "test-cluster-id", notes)
		Expect(notes.String()).To(ContainSubstring("No known error patterns"))
	})

	It("should handle empty logs", func() {
		parseMUOLogs("", "test-cluster-id", notes)
		Expect(notes.String()).To(ContainSubstring("No known error patterns"))
	})

	It("should return OHSS action for TLS handshake timeout", func() {
		logs := `2026-10-19T10:01:00Z ERROR net/http: TLS handshake timeout`
		actions := parseMUOLogs(logs, "test-cluster-id", notes)
		Expect(actions).To(HaveLen(1))
		ohss, ok := actions[0].(*executor.OHSSCaseAction)
		Expect(ok).To(BeTrue())
		Expect(ohss.ClusterID).To(Equal("test-cluster-id"))
		Expect(ohss.Summary).To(ContainSubstring("TLS handshake timeout"))
		Expect(ohss.Priority).To(Equal("Major"))
	})

	It("should return no actions when no OHSS-worthy patterns found", func() {
		logs := `2026-10-19T10:00:00Z ERROR ocm-config-getter failed`
		actions := parseMUOLogs(logs, "test-cluster-id", notes)
		Expect(actions).To(BeEmpty())
	})
})

var _ = Describe("checkMUOConfig", func() {
	var notes *notewriter.NoteWriter

	BeforeEach(func() {
		logging.InitLogger("fatal", "", "")
		notes = notewriter.New("test", logging.RawLogger)
	})

	It("should succeed with correct ocmBaseUrl", func() {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      muoConfigMapName,
				Namespace: muoNamespace,
			},
			Data: map[string]string{
				muoConfigKey: `configManager:
  source: OCM
  ocmBaseUrl: https://api.openshift.com`,
			},
		}
		k8sClient := fake.NewClientBuilder().WithScheme(testScheme()).WithObjects(cm).Build()
		err := checkMUOConfig(context.Background(), k8sClient, notes)
		Expect(err).NotTo(HaveOccurred())
		Expect(notes.String()).To(ContainSubstring("ocmBaseUrl is correct"))
	})

	It("should warn with wrong ocmBaseUrl", func() {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      muoConfigMapName,
				Namespace: muoNamespace,
			},
			Data: map[string]string{
				muoConfigKey: `configManager:
  source: OCM
  ocmBaseUrl: https://wrong.example.com`,
			},
		}
		k8sClient := fake.NewClientBuilder().WithScheme(testScheme()).WithObjects(cm).Build()
		err := checkMUOConfig(context.Background(), k8sClient, notes)
		Expect(err).NotTo(HaveOccurred())
		Expect(notes.String()).To(ContainSubstring("wrong.example.com"))
		Expect(notes.String()).To(ContainSubstring(expectedOCMBase))
	})

	It("should error when configmap not found", func() {
		k8sClient := fake.NewClientBuilder().WithScheme(testScheme()).Build()
		err := checkMUOConfig(context.Background(), k8sClient, notes)
		Expect(err).To(HaveOccurred())
		Expect(notes.String()).To(ContainSubstring("Failed to read MUO configmap"))
	})

	It("should error when config.yaml key missing", func() {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      muoConfigMapName,
				Namespace: muoNamespace,
			},
			Data: map[string]string{
				"other-key": "value",
			},
		}
		k8sClient := fake.NewClientBuilder().WithScheme(testScheme()).WithObjects(cm).Build()
		err := checkMUOConfig(context.Background(), k8sClient, notes)
		Expect(err).To(HaveOccurred())
		Expect(notes.String()).To(ContainSubstring("missing 'config.yaml' key"))
	})

	It("should error on invalid YAML", func() {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      muoConfigMapName,
				Namespace: muoNamespace,
			},
			Data: map[string]string{
				muoConfigKey: `: invalid yaml [[[`,
			},
		}
		k8sClient := fake.NewClientBuilder().WithScheme(testScheme()).WithObjects(cm).Build()
		err := checkMUOConfig(context.Background(), k8sClient, notes)
		Expect(err).To(HaveOccurred())
		Expect(notes.String()).To(ContainSubstring("Failed to parse MUO config"))
	})
})

var _ = Describe("Run", func() {
	var (
		mockCtrl *gomock.Controller
		r        *investigation.ResourceBuilderMock
		cluster  *cmv1.Cluster
		inv      Investigation
	)

	BeforeEach(func() {
		logging.InitLogger("fatal", "", "")
		mockCtrl = gomock.NewController(GinkgoT())
		var err error
		cluster, err = cmv1.NewCluster().ID("test-cluster-id").Build()
		Expect(err).NotTo(HaveOccurred())

		r = &investigation.ResourceBuilderMock{
			Resources: &investigation.Resources{
				Cluster:   cluster,
				OcmClient: ocmmock.NewMockClient(mockCtrl),
			},
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	When("user is banned due to export control", func() {
		It("should escalate", func() {
			account, _ := amv1.NewAccount().ID("user-1").
				Banned(true).BanCode("export_control_compliance").
				BanDescription("Export control").Build()
			r.Resources.OcmClient.(*ocmmock.MockClient).EXPECT().
				GetCreatorFromCluster(gomock.Any()).
				Return(account, nil)

			result, err := inv.Run(context.Background(), r)
			Expect(err).NotTo(HaveOccurred())
			Expect(hasEscalateAction(result.Actions)).To(BeTrue())
			Expect(hasSilenceAction(result.Actions)).To(BeFalse())
		})
	})

	When("user is banned for other reasons", func() {
		It("should send service log and escalate", func() {
			account, _ := amv1.NewAccount().ID("user-1").
				Banned(true).BanCode("terms_and_conditions").
				BanDescription("Terms violation").Build()
			r.Resources.OcmClient.(*ocmmock.MockClient).EXPECT().
				GetCreatorFromCluster(gomock.Any()).
				Return(account, nil)

			result, err := inv.Run(context.Background(), r)
			Expect(err).NotTo(HaveOccurred())
			Expect(hasServiceLogAction(result.Actions)).To(BeTrue())
			Expect(hasEscalateAction(result.Actions)).To(BeTrue())
			Expect(hasSilenceAction(result.Actions)).To(BeFalse())
		})
	})

	When("GetCreatorFromCluster fails", func() {
		It("should escalate", func() {
			r.Resources.OcmClient.(*ocmmock.MockClient).EXPECT().
				GetCreatorFromCluster(gomock.Any()).
				Return(nil, fmt.Errorf("connection timeout"))

			result, err := inv.Run(context.Background(), r)
			Expect(err).NotTo(HaveOccurred())
			Expect(hasEscalateAction(result.Actions)).To(BeTrue())
			Expect(hasSilenceAction(result.Actions)).To(BeFalse())
		})
	})

	When("initial resource build fails", func() {
		It("should propagate the error", func() {
			r.BuildError = fmt.Errorf("infrastructure error: resource build failed")

			_, err := inv.Run(context.Background(), r)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("infrastructure error"))
		})
	})

	When("second build fails with cluster access error", func() {
		It("should escalate with cluster access message", func() {
			account, _ := amv1.NewAccount().ID("user-1").Banned(false).Build()
			r.Resources.OcmClient.(*ocmmock.MockClient).EXPECT().
				GetCreatorFromCluster(gomock.Any()).
				Return(account, nil)

			fb := &k8sFailBuilder{
				inner:  r,
				k8sErr: investigation.RestConfigError{ClusterID: "test-cluster-id", Err: fmt.Errorf("unable to get credentials")},
			}

			result, err := inv.Run(context.Background(), fb)
			Expect(err).NotTo(HaveOccurred())
			Expect(hasEscalateAction(result.Actions)).To(BeTrue())
			Expect(hasSilenceAction(result.Actions)).To(BeFalse())
		})
	})

	When("second build fails with non-cluster-access error", func() {
		It("should propagate as infrastructure error", func() {
			account, _ := amv1.NewAccount().ID("user-1").Banned(false).Build()
			r.Resources.OcmClient.(*ocmmock.MockClient).EXPECT().
				GetCreatorFromCluster(gomock.Any()).
				Return(account, nil)

			fb := &k8sFailBuilder{
				inner:  r,
				k8sErr: fmt.Errorf("unexpected failure"),
			}

			_, err := inv.Run(context.Background(), fb)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Resource build error"))
		})
	})
})
