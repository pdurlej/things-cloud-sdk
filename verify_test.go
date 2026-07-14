package thingscloud

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Verify(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		server := fakeServer(fakeResponse{200, "verify-success.json"})
		defer server.Close()

		c := New(fmt.Sprintf("http://%s", server.Listener.Addr().String()), "martin@example.com", "")
		v, err := c.Verify()
		if err != nil {
			t.Fatalf("Expected Verification to succeed, but didn't: %q", err.Error())
		}
		if v.Status != AccountStatusActive {
			t.Errorf("Expected account to be %q, but got %q", AccountStatusActive, v.Status)
		}
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()
		server := fakeServer(fakeResponse{401, "error.json"})
		defer server.Close()

		c := New(fmt.Sprintf("http://%s", server.Listener.Addr().String()), "unknown@example.com", "")
		_, err := c.Verify()
		if err == nil {
			t.Error("Expected Verification to fail, but didn't")
		}
	})

	t.Run("InvalidEmail", func(t *testing.T) {
		c := New("https://example.com", "invalid\nemail", "secret")
		if _, err := c.Verify(); err == nil {
			t.Fatal("Verify succeeded with an invalid email")
		}
	})

	t.Run("MalformedJSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("not-json"))
		}))
		defer server.Close()

		if _, err := New(server.URL, "test@example.com", "secret").Verify(); err == nil {
			t.Fatal("Verify accepted malformed JSON")
		}
	})
}
