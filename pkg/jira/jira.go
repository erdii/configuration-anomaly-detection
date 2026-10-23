// Package jira contains Jira API client for creating and managing issues.
//
// Implementation inspired by github.com/andygrunwald/go-jira (MIT License).
// Original Copyright (c) 2015 Andy Grunwald. See LICENSE-THIRD-PARTY for details.
// Only the minimal REST plumbing for issue CRUD and search was extracted;
// all external dependencies were replaced with stdlib equivalents.
package jira

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const clientTimeout = 30 * time.Second

// maxErrorResponseBytes limits how much of a Jira error response body we read
// to prevent memory exhaustion from large HTML error pages.
const maxErrorResponseBytes = 1 << 20 // 1 MB

var jiraKeyPattern = regexp.MustCompile(`^[A-Z][A-Z0-9]+-\d+$`)

//go:generate mockgen --build_flags=--mod=readonly -source $GOFILE -destination ./mock/jiramock.go -package jiramock

// Client is the interface exposing Jira functions
type Client interface {
	CreateIssue(ctx context.Context, opts *CreateIssueOptions) (*Issue, error)
	GetIssue(ctx context.Context, key string) (*Issue, error)
	AddComment(ctx context.Context, issueKey string, body string) error
	SearchIssues(ctx context.Context, jql string, maxResults int) ([]*Issue, error)
}

// CreateIssueOptions contains the fields for creating a Jira issue
type CreateIssueOptions struct {
	ProjectKey   string
	IssueType    string
	Summary      string
	Description  string
	Priority     string
	Labels       []string
	Components   []string
	CustomFields map[string]any
}

var _ Client = (*SdkClient)(nil)

// Issue represents a Jira issue
type Issue struct {
	Key    string
	ID     string
	Self   string
	WebURL string
}

// SdkClient implements Client using the Jira REST API v2 with basic auth
type SdkClient struct {
	baseURL    *url.URL
	httpClient *http.Client
}

// New creates a new Jira SdkClient.
// baseURL is the Jira instance URL (e.g. "https://redhat.atlassian.net").
// email and token are used for HTTP Basic Authentication (API token auth).
func New(baseURL, email, token string) (*SdkClient, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("baseURL is required")
	}
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if token == "" {
		return nil, fmt.Errorf("token is required")
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid baseURL: %w", err)
	}
	if !strings.HasSuffix(u.Path, "/") {
		u.Path += "/"
	}

	return &SdkClient{
		baseURL: u,
		httpClient: &http.Client{
			Timeout: clientTimeout,
			Transport: &basicAuthTransport{
				username:  email,
				apiToken:  token,
				transport: http.DefaultTransport.(*http.Transport).Clone(),
			},
		},
	}, nil
}

func (c *SdkClient) CreateIssue(ctx context.Context, opts *CreateIssueOptions) (*Issue, error) {
	payload := c.buildCreatePayload(opts)

	body, err := c.do(ctx, http.MethodPost, "rest/api/2/issue", payload)
	if err != nil {
		return nil, fmt.Errorf("create issue: %w", err)
	}
	defer body.Close()

	var resp issueResponse
	if err := json.NewDecoder(body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("decode create response: %w", err)
	}

	return &Issue{
		Key:    resp.Key,
		ID:     resp.ID,
		Self:   resp.Self,
		WebURL: c.issueWebURL(resp.Key),
	}, nil
}

func (c *SdkClient) GetIssue(ctx context.Context, key string) (*Issue, error) {
	if !jiraKeyPattern.MatchString(key) {
		return nil, fmt.Errorf("invalid Jira issue key: %q", key)
	}
	body, err := c.do(ctx, http.MethodGet, fmt.Sprintf("rest/api/2/issue/%s", key), nil)
	if err != nil {
		return nil, fmt.Errorf("get issue %s: %w", key, err)
	}
	defer body.Close()

	var resp issueResponse
	if err := json.NewDecoder(body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("decode get response: %w", err)
	}

	return &Issue{
		Key:    resp.Key,
		ID:     resp.ID,
		Self:   resp.Self,
		WebURL: c.issueWebURL(resp.Key),
	}, nil
}

func (c *SdkClient) AddComment(ctx context.Context, issueKey string, commentBody string) error {
	if !jiraKeyPattern.MatchString(issueKey) {
		return fmt.Errorf("invalid Jira issue key: %q", issueKey)
	}
	payload := map[string]string{"body": commentBody}

	body, err := c.do(ctx, http.MethodPost, fmt.Sprintf("rest/api/2/issue/%s/comment", issueKey), payload)
	if err != nil {
		return fmt.Errorf("add comment to %s: %w", issueKey, err)
	}
	body.Close()
	return nil
}

func (c *SdkClient) SearchIssues(ctx context.Context, jql string, maxResults int) ([]*Issue, error) {
	endpoint := "rest/api/2/search"
	params := url.Values{}
	if jql != "" {
		params.Set("jql", jql)
	}
	if maxResults > 0 {
		params.Set("maxResults", strconv.Itoa(maxResults))
	}
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	body, err := c.do(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("search issues: %w", err)
	}
	defer body.Close()

	var resp searchResponse
	if err := json.NewDecoder(body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("decode search response: %w", err)
	}

	issues := make([]*Issue, len(resp.Issues))
	for i, raw := range resp.Issues {
		issues[i] = &Issue{
			Key:    raw.Key,
			ID:     raw.ID,
			Self:   raw.Self,
			WebURL: c.issueWebURL(raw.Key),
		}
	}
	return issues, nil
}

// do executes an HTTP request against the Jira API.
// If payload is non-nil it is JSON-encoded as the request body.
// Returns the response body (caller must close) or an error.
func (c *SdkClient) do(ctx context.Context, method, path string, payload any) (io.ReadCloser, error) {
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	rel.Path = strings.TrimPrefix(rel.Path, "/")
	u := c.baseURL.ResolveReference(rel)

	var reqBody io.Reader
	if payload != nil {
		buf := new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(payload); err != nil {
			return nil, fmt.Errorf("encode request: %w", err)
		}
		reqBody = buf
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
	if err != nil {
		return nil, err
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		defer resp.Body.Close()
		return nil, parseErrorResponse(resp)
	}

	return resp.Body, nil
}

func (c *SdkClient) buildCreatePayload(opts *CreateIssueOptions) map[string]any {
	fields := map[string]any{
		"project":   map[string]string{"key": opts.ProjectKey},
		"issuetype": map[string]string{"name": opts.IssueType},
		"summary":   opts.Summary,
	}

	if opts.Description != "" {
		fields["description"] = opts.Description
	}
	if opts.Priority != "" {
		fields["priority"] = map[string]string{"name": opts.Priority}
	}
	if len(opts.Labels) > 0 {
		fields["labels"] = opts.Labels
	}
	if len(opts.Components) > 0 {
		components := make([]map[string]string, len(opts.Components))
		for i, name := range opts.Components {
			components[i] = map[string]string{"name": name}
		}
		fields["components"] = components
	}
	for k, v := range opts.CustomFields {
		fields[k] = v
	}

	return map[string]any{"fields": fields}
}

func (c *SdkClient) issueWebURL(key string) string {
	return fmt.Sprintf("%s://%s/browse/%s", c.baseURL.Scheme, c.baseURL.Host, key)
}

// issueResponse is the minimal Jira issue JSON response we parse
type issueResponse struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Self string `json:"self"`
}

// searchResponse is the Jira search endpoint response
type searchResponse struct {
	Issues     []issueResponse `json:"issues"`
	StartAt    int             `json:"startAt"`
	MaxResults int             `json:"maxResults"`
	Total      int             `json:"total"`
}

// jiraErrorResponse mirrors Jira's error JSON format
type jiraErrorResponse struct {
	ErrorMessages []string          `json:"errorMessages"`
	Errors        map[string]string `json:"errors"`
}

func parseErrorResponse(resp *http.Response) error {
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxErrorResponseBytes))
	if err != nil {
		return fmt.Errorf("jira API error (status %d): could not read body: %w", resp.StatusCode, err)
	}

	contentType := resp.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		var jerr jiraErrorResponse
		if jsonErr := json.Unmarshal(body, &jerr); jsonErr == nil {
			if len(jerr.ErrorMessages) > 0 {
				return fmt.Errorf("jira API error (status %d): %s", resp.StatusCode, jerr.ErrorMessages[0])
			}
			for field, msg := range jerr.Errors {
				return fmt.Errorf("jira API error (status %d): %s - %s", resp.StatusCode, field, msg)
			}
		}
	}

	return fmt.Errorf("jira API error (status %d): %s", resp.StatusCode, string(body))
}

// basicAuthTransport is an http.RoundTripper that adds HTTP Basic Authentication
// using an email and API token (Jira Cloud pattern).
type basicAuthTransport struct {
	username  string
	apiToken  string
	transport http.RoundTripper
}

func (t *basicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := req.Clone(req.Context())
	req2.SetBasicAuth(t.username, t.apiToken)
	return t.transport.RoundTrip(req2)
}
