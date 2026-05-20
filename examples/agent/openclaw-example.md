# OpenClaw Agent Example

Use this repository through the Cloud API path when you want cross-platform
agent access to Things. Use the local SQLite reader only for read-only macOS
inspection.

## Recommended Tooling

- MCP server: `things-mcp`
- CLI fallback: `things-cloud-cli`
- Task list format: `--simple`
- Write policy: dry-run first, then confirm, then execute

## Example Agent Policy

```text
When the user asks to inspect Things tasks, call list/search tools first.
When the user asks to create or modify a task, run the operation in dry-run mode
first and summarize the planned user-visible change. Execute the non-dry-run
write only after confirmation.
Never write to local Things SQLite.
```

## Example Flows

List Today:

```bash
things-cloud-cli today --simple
```

Search:

```bash
things-cloud-cli search "invoice" --simple
```

Create with preview:

```bash
things-cloud-cli create "Follow up with Marta" --when today --dry-run
things-cloud-cli create "Follow up with Marta" --when today
```

Complete with preview:

```bash
things-cloud-cli complete <task-uuid> --dry-run
things-cloud-cli complete <task-uuid>
```

Read completion evidence:

```bash
things-cloud-cli completed --since 2026-05-20T00:00:00Z --format full
```

Create recurring task with preview:

```bash
things-cloud-cli create "Check car listings" --repeat every-day --dry-run
things-cloud-cli create "Check car listings" --repeat every-day
```

## Notes for OpenClaw Skills

- Prefer MCP tools when available.
- Keep raw command output in JSON and let the host render summaries.
- Use task UUIDs from `--simple` output for follow-up write commands.
- If the task reference is ambiguous, search first and ask the user to choose.
- Do not infer completion from absence in Today. Use `completed`/`logbook`.
