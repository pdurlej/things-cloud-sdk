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

	t.Run("TracksLoadedServerIndexFromStartIndex", func(t *testing.T) {
		t.Parallel()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got := r.URL.Query().Get("start-index"); got != "100" {
				t.Errorf("start-index = %q, want %q", got, "100")
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"items":[{"task-101":{"e":"Task6","t":0,"p":{"tt":"Task 101","tp":0}}},{"task-102":{"e":"Task6","t":0,"p":{"tt":"Task 102","tp":0}}}],"current-item-index":105,"schema":301}`)
		}))
		defer server.Close()

		c := New(server.URL, "martin@example.com", "")
		h := &History{
			Client: c,
			ID:     "33333abb-bfe4-4b03-a5c9-106d42220c72",
		}
		_, more, err := h.Items(ItemsOptions{StartIndex: 100})
		if err != nil {
			t.Fatalf("Items failed: %v", err)
		}
		if h.LoadedServerIndex != 102 {
			t.Errorf("LoadedServerIndex = %d, want %d", h.LoadedServerIndex, 102)
		}
		if h.LatestServerIndex != 105 {
			t.Errorf("LatestServerIndex = %d, want %d", h.LatestServerIndex, 105)
		}
		if !more {
			t.Error("Expected more items")
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

	t.Run("InvalidHistoryID", func(t *testing.T) {
		c := New("https://example.com", "test@example.com", "secret")
		h := c.HistoryWithID("invalid\nvalue")
		if _, _, err := h.Items(ItemsOptions{}); err == nil {
			t.Fatal("Items succeeded with an invalid history ID")
		}
	})

	t.Run("RejectsMalformedPages", func(t *testing.T) {
		tests := []struct {
			name string
			body string
		}{
			{"missing items", `{"current-item-index":0,"schema":301}`},
			{"stalled page", `{"items":[],"current-item-index":2,"schema":301}`},
			{"empty batch", `{"items":[{}],"current-item-index":1,"schema":301}`},
			{"missing kind", `{"items":[{"task-a":{"t":0,"p":{}}}],"current-item-index":1,"schema":301}`},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					fmt.Fprint(w, tt.body)
				}))
				defer server.Close()

				h := New(server.URL, "test@example.com", "secret").HistoryWithID("history-id")
				if _, _, err := h.Items(ItemsOptions{}); err == nil {
					t.Fatal("Items accepted a malformed page")
				}
			})
		}
	})

	t.Run("AllowsEmptyCaughtUpPage", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"items":[],"current-item-index":4,"schema":301}`)
		}))
		defer server.Close()

		h := New(server.URL, "test@example.com", "secret").HistoryWithID("history-id")
		items, more, err := h.Items(ItemsOptions{StartIndex: 4})
		if err != nil {
			t.Fatalf("Items rejected a valid caught-up page: %v", err)
		}
		if len(items) != 0 || more {
			t.Fatalf("Items returned items=%d more=%v, want empty and caught up", len(items), more)
		}
	})

	t.Run("AcceptsDeletionAndTombstoneKinds", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"items":[{"deleted-task":{"e":"Task6","t":2,"p":{}}},{"tombstone":{"e":"Tombstone2","t":0,"p":{"do":"deleted-task"}}}],"current-item-index":2,"schema":301}`)
		}))
		defer server.Close()

		h := New(server.URL, "test@example.com", "secret").HistoryWithID("history-id")
		items, _, err := h.Items(ItemsOptions{})
		if err != nil {
			t.Fatalf("Items rejected valid deletion kinds: %v", err)
		}
		if len(items) != 2 || items[0].Kind != ItemKindTask || items[1].Kind != ItemKindTombstone {
			t.Fatalf("Items returned unexpected deletion kinds: %#v", items)
		}
	})
}
