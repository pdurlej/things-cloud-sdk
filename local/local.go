// Package local provides a read-only adapter for local Things SQLite databases.
package local

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "modernc.org/sqlite"
)

type Reader struct {
	db *sql.DB
}

type Query struct {
	Search           string
	Limit            int
	IncludeCompleted bool
}

type Task struct {
	UUID   string `json:"uuid"`
	Title  string `json:"title"`
	Status string `json:"status"`
	Notes  string `json:"notes,omitempty"`
}

func DefaultDatabasePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Library", "Group Containers", "JLMPQHK86H.com.culturedcode.ThingsMac", "Things Database.thingsdatabase", "main.sqlite"), nil
}

func OpenDefault() (*Reader, error) {
	path, err := DefaultDatabasePath()
	if err != nil {
		return nil, err
	}
	return Open(path)
}

func Open(path string) (*Reader, error) {
	u := url.URL{Scheme: "file", Path: path, RawQuery: "mode=ro&immutable=1"}
	db, err := sql.Open("sqlite", u.String())
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return &Reader{db: db}, nil
}

func (r *Reader) Close() error {
	return r.db.Close()
}

func (r *Reader) Tasks(ctx context.Context, q Query) ([]Task, error) {
	table, cols, err := r.detectTaskTable(ctx)
	if err != nil {
		return nil, err
	}

	uuidCol := firstColumn(cols, "uuid", "uid", "id")
	titleCol := firstColumn(cols, "title", "name")
	statusCol := firstColumn(cols, "status")
	notesCol := firstColumn(cols, "notes", "note", "notesplain", "notetext")
	if uuidCol == "" || titleCol == "" {
		return nil, fmt.Errorf("task table %s is missing uuid/title columns", table)
	}

	query := fmt.Sprintf("SELECT %s, %s, %s, %s FROM %s",
		quoteIdent(uuidCol),
		quoteIdent(titleCol),
		selectColumn(statusCol, "0"),
		selectColumn(notesCol, "''"),
		quoteIdent(table),
	)
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	search := strings.ToLower(strings.TrimSpace(q.Search))
	var tasks []Task
	for rows.Next() {
		var uuid, title, statusRaw, notes string
		if err := rows.Scan(&uuid, &title, &statusRaw, &notes); err != nil {
			return nil, err
		}
		status := normalizeStatus(statusRaw)
		if !q.IncludeCompleted && status != "open" {
			continue
		}
		if search != "" && !strings.Contains(strings.ToLower(title), search) && !strings.Contains(strings.ToLower(notes), search) {
			continue
		}
		tasks = append(tasks, Task{UUID: uuid, Title: title, Status: status, Notes: notes})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	sort.Slice(tasks, func(i, j int) bool {
		if tasks[i].Title != tasks[j].Title {
			return tasks[i].Title < tasks[j].Title
		}
		return tasks[i].UUID < tasks[j].UUID
	})
	if q.Limit > 0 && q.Limit < len(tasks) {
		tasks = tasks[:q.Limit]
	}
	return tasks, nil
}

func (r *Reader) detectTaskTable(ctx context.Context) (string, map[string]string, error) {
	tables, err := r.tables(ctx)
	if err != nil {
		return "", nil, err
	}
	for _, preferred := range []string{"TMTask", "Task", "tasks"} {
		for _, table := range tables {
			if strings.EqualFold(table, preferred) {
				cols, err := r.columns(ctx, table)
				return table, cols, err
			}
		}
	}
	for _, table := range tables {
		cols, err := r.columns(ctx, table)
		if err != nil {
			return "", nil, err
		}
		if firstColumn(cols, "uuid", "uid", "id") != "" && firstColumn(cols, "title", "name") != "" {
			return table, cols, nil
		}
	}
	return "", nil, fmt.Errorf("no task-like table found")
}

func (r *Reader) tables(ctx context.Context) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT name FROM sqlite_master WHERE type='table'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		tables = append(tables, name)
	}
	return tables, rows.Err()
}

func (r *Reader) columns(ctx context.Context, table string) (map[string]string, error) {
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf("PRAGMA table_info(%s)", quoteIdent(table)))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cols := map[string]string{}
	for rows.Next() {
		var cid int
		var name, typ string
		var notNull int
		var defaultValue any
		var pk int
		if err := rows.Scan(&cid, &name, &typ, &notNull, &defaultValue, &pk); err != nil {
			return nil, err
		}
		cols[strings.ToLower(name)] = name
	}
	return cols, rows.Err()
}

func firstColumn(cols map[string]string, names ...string) string {
	for _, name := range names {
		if col, ok := cols[strings.ToLower(name)]; ok {
			return col
		}
	}
	return ""
}

func selectColumn(col, fallback string) string {
	if col == "" {
		return fallback
	}
	return quoteIdent(col)
}

func quoteIdent(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}

func normalizeStatus(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", "0", "open", "incomplete", "pending":
		return "open"
	case "2", "canceled", "cancelled":
		return "canceled"
	case "3", "completed", "done":
		return "completed"
	default:
		return raw
	}
}
