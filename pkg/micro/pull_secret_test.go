package micro_test

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	sdk "github.com/openshift-online/ocm-sdk-go"
	"github.com/openshift/configuration-anomaly-detection/pkg/executor"
	"github.com/openshift/configuration-anomaly-detection/pkg/logging"
	"github.com/openshift/configuration-anomaly-detection/pkg/micro"
	"github.com/openshift/configuration-anomaly-detection/pkg/notewriter"
	"github.com/openshift/configuration-anomaly-detection/pkg/pullsecret"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func pullSecretScheme() *runtime.Scheme {
	s := runtime.NewScheme()
	utilruntime.Must(corev1.AddToScheme(s))
	return s
}

func newPullSecret(email string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pullsecret.PullSecretName,
			Namespace: pullsecret.PullSecretNamespace,
		},
		Type: corev1.SecretTypeDockerConfigJson,
		Data: map[string][]byte{
			".dockerconfigjson": []byte(fmt.Sprintf(
				`{"auths":{"cloud.openshift.com":{"auth":"dGVzdDp0ZXN0","email":"%s"}}}`,
				email,
			)),
		},
	}
}

func fakeAccessToken() string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"RS256","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"exp":9999999999,"typ":"Bearer"}`))
	return header + "." + payload + ".dGVzdA"
}

func newTestOCMConn(handler http.Handler) (*sdk.Connection, *httptest.Server) {
	server := httptest.NewServer(handler)
	conn, err := sdk.NewConnectionBuilder().
		URL(server.URL).
		Tokens(fakeAccessToken()).
		Build()
	Expect(err).NotTo(HaveOccurred())
	return conn, server
}

func emptyRegistryCredentialsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"items":[],"size":0,"total":0,"page":1}`)
	})
}

func errorHandler(code int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		fmt.Fprint(w, `{"reason":"test error"}`)
	})
}

var _ = Describe("PullSecret", func() {
	var notes *notewriter.NoteWriter

	BeforeEach(func() {
		logging.InitLogger("fatal", "", "")
		notes = notewriter.New("test", logging.RawLogger)
	})

	When("pull secret does not exist on cluster", func() {
		It("should return an infrastructure error", func() {
			k8sClient := fake.NewClientBuilder().WithScheme(pullSecretScheme()).Build()

			_, err := micro.PullSecret(k8sClient, nil, "user-id", "user@test.com", notes)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to get pull secret"))
		})
	})

	When("pull secret email does not match OCM email", func() {
		It("should escalate", func() {
			secret := newPullSecret("cluster@test.com")
			k8sClient := fake.NewClientBuilder().WithScheme(pullSecretScheme()).WithObjects(secret).Build()

			conn, server := newTestOCMConn(emptyRegistryCredentialsHandler())
			defer server.Close()
			defer conn.Close()

			result, err := micro.PullSecret(k8sClient, conn, "user-id", "ocm@test.com", notes)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Done).To(BeTrue())
			Expect(hasActionType(result.Actions, string(executor.ActionTypeEscalateIncident))).To(BeTrue())
			Expect(hasActionType(result.Actions, string(executor.ActionTypeSilenceIncident))).To(BeFalse())
			Expect(notes.String()).To(ContainSubstring("does not match"))
		})
	})

	When("pull secret email matches and no registry credentials exist", func() {
		It("should silence", func() {
			secret := newPullSecret("user@test.com")
			k8sClient := fake.NewClientBuilder().WithScheme(pullSecretScheme()).WithObjects(secret).Build()

			conn, server := newTestOCMConn(emptyRegistryCredentialsHandler())
			defer server.Close()
			defer conn.Close()

			result, err := micro.PullSecret(k8sClient, conn, "user-id", "user@test.com", notes)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Done).To(BeTrue())
			Expect(hasActionType(result.Actions, string(executor.ActionTypeSilenceIncident))).To(BeTrue())
			Expect(hasActionType(result.Actions, string(executor.ActionTypeEscalateIncident))).To(BeFalse())
		})
	})

	When("pull secret email matches but OCM registry credentials request fails", func() {
		It("should escalate", func() {
			secret := newPullSecret("user@test.com")
			k8sClient := fake.NewClientBuilder().WithScheme(pullSecretScheme()).WithObjects(secret).Build()

			conn, server := newTestOCMConn(errorHandler(http.StatusInternalServerError))
			defer server.Close()
			defer conn.Close()

			result, err := micro.PullSecret(k8sClient, conn, "user-id", "user@test.com", notes)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Done).To(BeTrue())
			Expect(hasActionType(result.Actions, string(executor.ActionTypeEscalateIncident))).To(BeTrue())
			Expect(hasActionType(result.Actions, string(executor.ActionTypeSilenceIncident))).To(BeFalse())
		})
	})

	When("email matches and notes reflect success", func() {
		It("should append success note for matching email", func() {
			secret := newPullSecret("user@test.com")
			k8sClient := fake.NewClientBuilder().WithScheme(pullSecretScheme()).WithObjects(secret).Build()

			conn, server := newTestOCMConn(emptyRegistryCredentialsHandler())
			defer server.Close()
			defer conn.Close()

			_, err := micro.PullSecret(k8sClient, conn, "user-id", "user@test.com", notes)
			Expect(err).NotTo(HaveOccurred())
			Expect(notes.String()).To(ContainSubstring("email matches"))
		})
	})
})
