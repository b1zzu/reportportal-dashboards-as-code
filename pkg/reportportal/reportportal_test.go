package reportportal

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	// baseURLPath is a non-empty Client.BaseURL path to use during tests,
	// to ensure relative URLs are used for all endpoints.
	baseURLPath = "/api-v3"
)

// setup sets up a test HTTP server along with a reportportal.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() (client *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	// We want to ensure that tests catch mistakes where the endpoint URL is
	// specified as absolute rather than relative. It only makes a difference
	// when there's a non-empty base URL path. So, use that.
	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(os.Stderr, "FAIL: Client.BaseURL path prefix is not preserved in the request URL:")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\t"+req.URL.String())
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\tDid you accidentally use an absolute endpoint URL rather than relative?")
		fmt.Fprintln(os.Stderr, "\tSee https://github.com/google/go-github/issues/752 for information.")
		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	// client is the ReportPortal client being tested and is
	// configured to use test server.
	client, _ = NewClient(nil, server.URL+baseURLPath+"/")

	return client, mux, server.URL, server.Close
}

func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

type values map[string]string

func testFormValues(t *testing.T, r *http.Request, values values) {
	t.Helper()
	want := url.Values{}
	for k, v := range values {
		want.Set(k, v)
	}

	r.ParseForm()
	if got := r.Form; !cmp.Equal(got, want) {
		t.Errorf("Request parameters: %v, want %v", got, want)
	}
}

func TestNewClient(t *testing.T) {

	tests := []*struct {
		description string
		httpClient  *http.Client
		baseURL     string

		expectError   string
		expectBaseURL string
	}{
		{
			description: "New client without api/ path should add it",
			httpClient:  &http.Client{},
			baseURL:     "example.com",

			expectBaseURL: "example.com/api/",
		},
		{
			description: "New client without trailing slash should add it",
			httpClient:  &http.Client{},
			baseURL:     "https://example.com/api",

			expectBaseURL: "https://example.com/api/",
		},
		{
			description: "New client without baseURL should fail",
			httpClient:  &http.Client{},
			baseURL:     "",

			expectError: "baseURL is empty",
		},
		{
			description: "New client without http.Client should create it",
			httpClient:  nil,
			baseURL:     "example.com",

			expectBaseURL: "example.com/api/",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {

			c, err := NewClient(test.httpClient, test.baseURL)
			if test.expectError != "" || err != nil {
				if err != nil {
					if err.Error() == test.expectError {
						// nice the test succeded
						return
					}
				}
				t.Fatalf("expected error \"%s\" but got \"%s\"", test.expectError, err)
			}

			if c.BaseURL.String() != test.expectBaseURL {
				t.Errorf("expected baseURL \"%s\" but got \"%s\"", test.expectBaseURL, c.BaseURL.String())
			}
		})
	}
}

func TestNewClientServices(t *testing.T) {

	c, err := NewClient(nil, "localhost")
	if err != nil {
		t.Fatal(err)
	}

	if c.client == nil {
		t.Error("client nil")
	}
	if c.Dashboard == nil {
		t.Error("dashboard service is nil")
	}
	if c.Widget == nil {
		t.Error("widget service is nil")
	}
	if c.Filter == nil {
		t.Error("filter service is nil")
	}
	if c.ProjectSettings == nil {
		t.Error("project settings service is nil")
	}
}
