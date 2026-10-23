package jira_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/openshift/configuration-anomaly-detection/pkg/jira"
)

var _ = Describe("Jira Client", func() {
	var (
		mux    *http.ServeMux
		server *httptest.Server
		client *jira.SdkClient
		ctx    context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		mux = http.NewServeMux()
		server = httptest.NewServer(mux)

		var err error
		client, err = jira.New(server.URL, "test@example.com", "test-token")
		Expect(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("New", func() {
		It("should return an error when baseURL is empty", func() {
			_, err := jira.New("", "email", "token")
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("baseURL is required"))
		})

		It("should return an error when email is empty", func() {
			_, err := jira.New("https://jira.example.com", "", "token")
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("email is required"))
		})

		It("should return an error when token is empty", func() {
			_, err := jira.New("https://jira.example.com", "email", "")
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("token is required"))
		})

		It("should succeed with valid inputs", func() {
			c, err := jira.New("https://jira.example.com", "email@test.com", "token123")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(c).ShouldNot(BeNil())
		})
	})

	Describe("CreateIssue", func() {
		It("should create an issue and return it", func() {
			var receivedBody map[string]any
			var receivedAuth string

			mux.HandleFunc("/rest/api/2/issue", func(w http.ResponseWriter, r *http.Request) {
				Expect(r.Method).Should(Equal("POST"))
				Expect(r.Header.Get("Content-Type")).Should(Equal("application/json"))
				Expect(r.Header.Get("Accept")).Should(Equal("application/json"))

				receivedAuth = r.Header.Get("Authorization")

				body, err := io.ReadAll(r.Body)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(json.Unmarshal(body, &receivedBody)).Should(Succeed())

				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`{"id":"10001","key":"OHSS-123","self":"https://jira.example.com/rest/api/2/issue/10001"}`))
			})

			issue, err := client.CreateIssue(ctx, &jira.CreateIssueOptions{
				ProjectKey:  "OHSS",
				IssueType:   "Bug",
				Summary:     "Test issue",
				Description: "Test description",
				Priority:    "High",
				Labels:      []string{"cad", "automated"},
				Components:  []string{"SRE"},
				CustomFields: map[string]any{
					"customfield_123": "custom-value",
				},
			})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(issue.Key).Should(Equal("OHSS-123"))
			Expect(issue.ID).Should(Equal("10001"))
			Expect(issue.Self).Should(Equal("https://jira.example.com/rest/api/2/issue/10001"))
			Expect(issue.WebURL).Should(ContainSubstring("/browse/OHSS-123"))

			Expect(receivedAuth).Should(HavePrefix("Basic "))

			fields, ok := receivedBody["fields"].(map[string]any)
			Expect(ok).Should(BeTrue())
			Expect(fields["summary"]).Should(Equal("Test issue"))
			Expect(fields["description"]).Should(Equal("Test description"))

			project, ok := fields["project"].(map[string]any)
			Expect(ok).Should(BeTrue())
			Expect(project["key"]).Should(Equal("OHSS"))

			issueType, ok := fields["issuetype"].(map[string]any)
			Expect(ok).Should(BeTrue())
			Expect(issueType["name"]).Should(Equal("Bug"))

			priority, ok := fields["priority"].(map[string]any)
			Expect(ok).Should(BeTrue())
			Expect(priority["name"]).Should(Equal("High"))

			labels, ok := fields["labels"].([]any)
			Expect(ok).Should(BeTrue())
			Expect(labels).Should(HaveLen(2))
			Expect(labels[0]).Should(Equal("cad"))
			Expect(labels[1]).Should(Equal("automated"))

			components, ok := fields["components"].([]any)
			Expect(ok).Should(BeTrue())
			Expect(components).Should(HaveLen(1))

			Expect(fields["customfield_123"]).Should(Equal("custom-value"))
		})

		It("should return an error on API failure", func() {
			mux.HandleFunc("/rest/api/2/issue", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(`{"errorMessages":["Summary is required"],"errors":{}}`))
			})

			issue, err := client.CreateIssue(ctx, &jira.CreateIssueOptions{
				ProjectKey: "OHSS",
				IssueType:  "Bug",
			})

			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("Summary is required"))
			Expect(issue).Should(BeNil())
		})

		It("should return an error with field-level errors", func() {
			mux.HandleFunc("/rest/api/2/issue", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(`{"errorMessages":[],"errors":{"priority":"Priority is required"}}`))
			})

			_, err := client.CreateIssue(ctx, &jira.CreateIssueOptions{
				ProjectKey: "OHSS",
				IssueType:  "Bug",
				Summary:    "Test",
			})

			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("priority"))
			Expect(err.Error()).Should(ContainSubstring("Priority is required"))
		})

		It("should send basic auth credentials", func() {
			var receivedAuth string

			mux.HandleFunc("/rest/api/2/issue", func(w http.ResponseWriter, r *http.Request) {
				receivedAuth = r.Header.Get("Authorization")
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`{"id":"1","key":"TEST-1","self":"http://test"}`))
			})

			_, err := client.CreateIssue(ctx, &jira.CreateIssueOptions{
				ProjectKey: "TEST",
				IssueType:  "Bug",
				Summary:    "Auth test",
			})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(receivedAuth).Should(HavePrefix("Basic "))
			Expect(receivedAuth).ShouldNot(BeEmpty())
		})
	})

	Describe("GetIssue", func() {
		It("should get an issue by key", func() {
			mux.HandleFunc("/rest/api/2/issue/OHSS-456", func(w http.ResponseWriter, r *http.Request) {
				Expect(r.Method).Should(Equal("GET"))
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"id":"10002","key":"OHSS-456","self":"https://jira.example.com/rest/api/2/issue/10002"}`))
			})

			issue, err := client.GetIssue(ctx, "OHSS-456")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(issue.Key).Should(Equal("OHSS-456"))
			Expect(issue.ID).Should(Equal("10002"))
			Expect(issue.Self).Should(Equal("https://jira.example.com/rest/api/2/issue/10002"))
			Expect(issue.WebURL).Should(ContainSubstring("/browse/OHSS-456"))
		})

		It("should return an error when issue is not found", func() {
			mux.HandleFunc("/rest/api/2/issue/NOPE-999", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(`{"errorMessages":["Issue does not exist or you do not have permission to see it."],"errors":{}}`))
			})

			issue, err := client.GetIssue(ctx, "NOPE-999")

			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("Issue does not exist"))
			Expect(issue).Should(BeNil())
		})
	})

	Describe("AddComment", func() {
		It("should add a comment to an issue", func() {
			var receivedBody map[string]string

			mux.HandleFunc("/rest/api/2/issue/OHSS-123/comment", func(w http.ResponseWriter, r *http.Request) {
				Expect(r.Method).Should(Equal("POST"))

				body, err := io.ReadAll(r.Body)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(json.Unmarshal(body, &receivedBody)).Should(Succeed())

				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`{"id":"20001","body":"CAD investigation complete"}`))
			})

			err := client.AddComment(ctx, "OHSS-123", "CAD investigation complete")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(receivedBody["body"]).Should(Equal("CAD investigation complete"))
		})

		It("should return an error on failure", func() {
			mux.HandleFunc("/rest/api/2/issue/OHSS-123/comment", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(`{"errorMessages":["Comment body can not be empty"],"errors":{}}`))
			})

			err := client.AddComment(ctx, "OHSS-123", "")

			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("Comment body can not be empty"))
		})
	})

	Describe("SearchIssues", func() {
		It("should return issues matching JQL", func() {
			mux.HandleFunc("/rest/api/2/search", func(w http.ResponseWriter, r *http.Request) {
				Expect(r.Method).Should(Equal("GET"))
				Expect(r.URL.Query().Get("jql")).Should(Equal("project = OHSS AND status = Open"))
				Expect(r.URL.Query().Get("maxResults")).Should(Equal("10"))

				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{
					"issues": [
						{"id":"10001","key":"OHSS-100","self":"https://jira.example.com/rest/api/2/issue/10001"},
						{"id":"10002","key":"OHSS-101","self":"https://jira.example.com/rest/api/2/issue/10002"}
					],
					"startAt": 0,
					"maxResults": 10,
					"total": 2
				}`))
			})

			issues, err := client.SearchIssues(ctx, "project = OHSS AND status = Open", 10)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(issues).Should(HaveLen(2))
			Expect(issues[0].Key).Should(Equal("OHSS-100"))
			Expect(issues[0].WebURL).Should(ContainSubstring("/browse/OHSS-100"))
			Expect(issues[1].Key).Should(Equal("OHSS-101"))
			Expect(issues[1].ID).Should(Equal("10002"))
		})

		It("should return empty slice for no results", func() {
			mux.HandleFunc("/rest/api/2/search", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"issues":[],"startAt":0,"maxResults":50,"total":0}`))
			})

			issues, err := client.SearchIssues(ctx, "project = OHSS AND key = NOPE", 50)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(issues).Should(HaveLen(0))
			Expect(issues).ShouldNot(BeNil())
		})

		It("should omit query params when not provided", func() {
			mux.HandleFunc("/rest/api/2/search", func(w http.ResponseWriter, r *http.Request) {
				Expect(r.URL.Query().Get("jql")).Should(BeEmpty())
				Expect(r.URL.Query().Get("maxResults")).Should(BeEmpty())

				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"issues":[],"startAt":0,"maxResults":50,"total":0}`))
			})

			_, err := client.SearchIssues(ctx, "", 0)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("Key validation", func() {
		It("should reject keys with path traversal", func() {
			_, err := client.GetIssue(ctx, "../../admin")
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("invalid Jira issue key"))
		})

		It("should reject keys with slashes", func() {
			err := client.AddComment(ctx, "OHSS-123/../../admin", "test")
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("invalid Jira issue key"))
		})

		It("should reject lowercase keys", func() {
			_, err := client.GetIssue(ctx, "ohss-123")
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("invalid Jira issue key"))
		})

		It("should reject empty keys", func() {
			_, err := client.GetIssue(ctx, "")
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("invalid Jira issue key"))
		})
	})

	Describe("Error handling", func() {
		It("should handle non-JSON error responses", func() {
			mux.HandleFunc("/rest/api/2/issue/ERR-1", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("<html>Internal Server Error</html>"))
			})

			_, err := client.GetIssue(ctx, "ERR-1")

			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("500"))
			Expect(err.Error()).Should(ContainSubstring("Internal Server Error"))
		})

		It("should handle empty error response body", func() {
			mux.HandleFunc("/rest/api/2/issue/ERR-2", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusForbidden)
			})

			_, err := client.GetIssue(ctx, "ERR-2")

			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("403"))
		})
	})
})
