package main

import (
	"os"
	"path/filepath"
	"testing"

	thingscloud "github.com/arthursoares/things-cloud-sdk"
	memory "github.com/arthursoares/things-cloud-sdk/state/memory"
)

func TestCommandNeedsHistoryHead(t *testing.T) {
	tests := []struct {
		cmd  string
		want bool
	}{
		{"list", false},
		{"show", false},
		{"areas", false},
		{"projects", false},
		{"tags", false},
		{"create", true},
		{"edit", true},
		{"complete", true},
		{"trash", true},
		{"batch", true},
	}

	for _, tt := range tests {
		t.Run(tt.cmd, func(t *testing.T) {
			if got := commandNeedsHistoryHead(tt.cmd); got != tt.want {
				t.Fatalf("commandNeedsHistoryHead(%q) = %v, want %v", tt.cmd, got, tt.want)
			}
		})
	}
}

func TestCLIStateCachePathUsesEnvOverride(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	t.Setenv("THINGS_CLI_CACHE", path)

	if got := cliStateCachePath(); got != path {
		t.Fatalf("cliStateCachePath() = %q, want %q", got, path)
	}
}

func TestCLIStateCacheMissingFile(t *testing.T) {
	cache, err := loadCLIStateCache(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatalf("loadCLIStateCache failed: %v", err)
	}
	if cache != nil {
		t.Fatalf("cache = %#v, want nil", cache)
	}
}

func TestCLIStateCacheRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "state.json")
	state := memory.NewState()
	state.Tasks["task-1"] = &thingscloud.Task{
		UUID:  "task-1",
		Title: "Cached Task",
	}
	state.Areas["area-1"] = &thingscloud.Area{
		UUID:  "area-1",
		Title: "Cached Area",
	}

	cache := &cliStateCache{
		HistoryID:   "history-1",
		ServerIndex: 42,
		State:       state,
	}
	if err := saveCLIStateCache(path, cache); err != nil {
		t.Fatalf("saveCLIStateCache failed: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat cache failed: %v", err)
	}
	if info.IsDir() {
		t.Fatal("cache path is a directory")
	}

	loaded, err := loadCLIStateCache(path)
	if err != nil {
		t.Fatalf("loadCLIStateCache failed: %v", err)
	}
	if loaded.HistoryID != "history-1" {
		t.Fatalf("HistoryID = %q, want history-1", loaded.HistoryID)
	}
	if loaded.ServerIndex != 42 {
		t.Fatalf("ServerIndex = %d, want 42", loaded.ServerIndex)
	}
	if loaded.State.Tasks["task-1"].Title != "Cached Task" {
		t.Fatalf("task title = %q, want Cached Task", loaded.State.Tasks["task-1"].Title)
	}
	if loaded.State.Areas["area-1"].Title != "Cached Area" {
		t.Fatalf("area title = %q, want Cached Area", loaded.State.Areas["area-1"].Title)
	}
}

func TestCLIStateCacheNormalizesEmptyState(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	cache := &cliStateCache{
		HistoryID:   "history-1",
		ServerIndex: 7,
		State:       &memory.State{},
	}
	if err := saveCLIStateCache(path, cache); err != nil {
		t.Fatalf("saveCLIStateCache failed: %v", err)
	}

	loaded, err := loadCLIStateCache(path)
	if err != nil {
		t.Fatalf("loadCLIStateCache failed: %v", err)
	}
	if loaded.State.Tasks == nil {
		t.Fatal("Tasks map was not initialized")
	}
	if loaded.State.Areas == nil {
		t.Fatal("Areas map was not initialized")
	}
	if loaded.State.Tags == nil {
		t.Fatal("Tags map was not initialized")
	}
	if loaded.State.CheckListItems == nil {
		t.Fatal("CheckListItems map was not initialized")
	}
}
