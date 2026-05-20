package local

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func TestReaderTasks(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.sqlite")
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if _, err := db.Exec(`
		CREATE TABLE TMTask (
			uuid TEXT PRIMARY KEY,
			title TEXT,
			status INTEGER,
			notes TEXT
		);
		INSERT INTO TMTask (uuid, title, status, notes) VALUES
			('task-2', 'Bravo', 3, 'done'),
			('task-1', 'Alpha', 0, 'needle');
	`); err != nil {
		t.Fatalf("seed sqlite: %v", err)
	}
	db.Close()

	reader, err := Open(path)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer reader.Close()

	tasks, err := reader.Tasks(context.Background(), Query{})
	if err != nil {
		t.Fatalf("Tasks failed: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("len(tasks) = %d, want 1", len(tasks))
	}
	if tasks[0].UUID != "task-1" || tasks[0].Status != "open" {
		t.Fatalf("task = %#v, want task-1 open", tasks[0])
	}

	tasks, err = reader.Tasks(context.Background(), Query{IncludeCompleted: true, Search: "bravo"})
	if err != nil {
		t.Fatalf("Tasks search failed: %v", err)
	}
	if len(tasks) != 1 || tasks[0].UUID != "task-2" || tasks[0].Status != "completed" {
		t.Fatalf("tasks = %#v, want completed task-2", tasks)
	}
}

func TestNormalizeStatus(t *testing.T) {
	tests := map[string]string{
		"0":          "open",
		"incomplete": "open",
		"2":          "canceled",
		"3":          "completed",
	}
	for raw, want := range tests {
		if got := normalizeStatus(raw); got != want {
			t.Fatalf("normalizeStatus(%q) = %q, want %q", raw, got, want)
		}
	}
}
