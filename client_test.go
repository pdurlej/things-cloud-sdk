package thingscloud

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func stringVal(str string) *string {
	return &str
}

type fakeResponse struct {
	statusCode int
	file       string
}

func fakeServer(t fakeResponse) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open(fmt.Sprintf("tapes/%s", t.file))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Printf("Unable to open fixture %q path %q", t.file, r.URL.Path)
			return
		}
		defer f.Close()

		content, err := io.ReadAll(f)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Printf("Unable to load fixture %q path %q", t.file, r.URL.Path)
			return
		}
		w.WriteHeader(t.statusCode)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(content))
	}))
}

func TestClient_UserAgent(t *testing.T) {
	var capturedHeaders http.Header
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeaders = r.Header.Clone()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	c := New(ts.URL, "test@example.com", "password")

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.do(req)
	if err != nil {
		t.Fatal(err)
	}

	// Verify User-Agent is the updated value
	got := capturedHeaders.Get("User-Agent")
	want := "ThingsMac/32209501"
	if got != want {
		t.Errorf("User-Agent = %q, want %q", got, want)
	}

	// Verify things-client-info header is set and non-empty
	clientInfo := capturedHeaders.Get("Things-Client-Info")
	if clientInfo == "" {
		t.Error("things-client-info header is missing or empty")
	}
}

func TestClient_DebugLoggingRedactsSensitiveData(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"private":"response-body"}`))
	}))
	defer ts.Close()

	var logs bytes.Buffer
	oldOutput := log.Writer()
	log.SetOutput(&logs)
	t.Cleanup(func() { log.SetOutput(oldOutput) })

	c := New(ts.URL, "test@example.com", "secret-password")
	c.Debug = true
	req, err := http.NewRequest("POST", "/test", strings.NewReader(`{"password":"new-secret"}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Password secret-password")

	resp, err := c.do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	got := logs.String()
	for _, sensitive := range []string{"secret-password", "new-secret", "response-body", "Authorization"} {
		if strings.Contains(got, sensitive) {
			t.Errorf("debug log contains sensitive value %q", sensitive)
		}
	}
	if !strings.Contains(got, "REQUEST: POST") || !strings.Contains(got, "RESPONSE: 200 OK") {
		t.Errorf("debug log is missing request or response metadata: %q", got)
	}
}
