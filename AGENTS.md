# Agent Instructions

This repository is an unofficial Things Cloud SDK, CLI, and MCP server. Optimize
changes for safe automation, stable JSON, and clear separation between Cloud
writes and local read-only inspection.

## Default Choices

- Use `things-cloud-cli` in new docs, examples, and agent integrations.
- Keep `things-cli` only as a backward-compatible alias.
- Prefer `things-mcp` for MCP-based agents.
- Prefer `sync.Open(...).Sync()` or `QuickSync()` for long-running services that
  need persisted state and semantic change events.
- Use `local.OpenDefault()` only for read-only macOS SQLite inspection.

## Safety Rules

- Never add code that writes to the local Things SQLite database.
- Use `--dry-run` examples for agent-generated writes.
- Do not commit credentials, account emails, real task UUIDs from a user's
  database, local cache files, or sync DB files.
- Do not invent Things wire payloads casually. Read `docs/client-side-bugs.md`
  before changing write behavior.
- Preserve Things wire-format invariants:
  - Write UUIDs must be Things-compatible Base58.
  - `md` must be `null` on creates.
  - `st` is schedule: `0` Inbox, `1` Anytime/Today, `2` Someday/Upcoming.
  - `ss` is status: `0` pending, `2` canceled, `3` completed.
  - Headings (`tp=2`) must use `st=1`.
  - Tasks under projects, headings, or areas should default to `st=1`.

## Documentation Rules

- Public repository docs should be written in English.
- Keep examples copy-pasteable and safe by default.
- Show compact JSON (`--simple`) for task-listing agent examples.
- Show `--dry-run` before non-read agent write examples.
- Prefer placeholders such as `<task-uuid>` and `you@example.com`.
- Keep docs ASCII unless a file already uses non-ASCII intentionally.

## Development Checks

Run these before committing code changes:

```bash
go test ./...
git diff --check
```

For documentation-only changes, still run:

```bash
git diff --check
files="$(git diff --name-only --diff-filter=ACM)"
test -z "$files" || LC_ALL=C rg -n '[^ -~\t]' $files || true
```

If the non-ASCII scan prints intentional content, document why in the PR.

## Useful Entry Points

- `README.md` - main user and agent onboarding.
- `llms.txt` - short LLM crawler/agent index.
- `docs/agent-cookbook.md` - task-oriented agent recipes.
- `docs/contracts.md` - JSON contracts for CLI and MCP outputs.
- `cmd/things-cloud-cli/` - preferred CLI entrypoint.
- `cmd/things-mcp/` - MCP stdio server.
- `internal/thingscli/` - shared CLI implementation and tests.
- `sync/` - persistent SQLite sync engine.
- `local/` - read-only local Things SQLite reader.
- `docs/client-side-bugs.md` - wire-format and crash analysis.
