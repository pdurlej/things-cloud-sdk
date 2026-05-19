package thingscloud

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHistory_Items(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		server := fakeServer(fakeResponse{200, "history-items-success.json"})
		defer server.Close()

		c := New(fmt.Sprintf("http://%s", server.Listener.Addr().String()), "martin@example.com", "")
		h := &History{
			Client: c,
			ID:     "33333abb-bfe4-4b03-a5c9-106d42220c72",
		}
		items, _, err := h.Items(ItemsOptions{})
		if err != nil {
			t.Fatalf("Expected items request to succeed, but didn't: %q", err.Error())
		}

		if len(items) < 1 {
			t.Fatalf("Expected items, but got none: %#v", items)
		}
	})

	t.Run("SetsServerIndexFromOuterItems", func(t *testing.T) {
		t.Parallel()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"items":[{"task-a":{"e":"Task6","t":0,"p":{"tt":"Task A"}},"task-b":{"e":"Task6","t":0,"p":{"tt":"Task B"}}},{"task-c":{"e":"Task6","t":0,"p":{"tt":"Task C"}}}],"current-item-index":12,"schema":301}`)
		}))
		defer server.Close()

		c := New(server.URL, "martin@example.com", "")
		h := &History{
			Client: c,
			ID:     "33333abb-bfe4-4b03-a5c9-106d42220c72",
		}

		items, _, err := h.Items(ItemsOptions{StartIndex: 10})
		if err != nil {
			t.Fatalf("Items failed: %v", err)
		}
		if len(items) != 3 {
			t.Fatalf("expected 3 flattened items, got %d", len(items))
		}

		indexByUUID := make(map[string]int)
		for _, item := range items {
			if !item.HasServerIndex {
				t.Fatalf("item %s is missing server index metadata", item.UUID)
			}
			indexByUUID[item.UUID] = item.ServerIndex
		}

		for _, uuid := range []string{"task-a", "task-b"} {
			if indexByUUID[uuid] != 10 {
				t.Errorf("%s ServerIndex = %d, want 10", uuid, indexByUUID[uuid])
			}
		}
		if indexByUUID["task-c"] != 11 {
			t.Errorf("task-c ServerIndex = %d, want 11", indexByUUID["task-c"])
		}
	})
}
