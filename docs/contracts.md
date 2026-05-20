# Agent Contracts

This document defines stable JSON shapes that agents can rely on when using the
CLI or MCP server. It is intentionally narrower than the full Things data model.

## CLI Task List, Simple Format

Commands:

```bash
things-cloud-cli today --simple
things-cloud-cli inbox --simple
things-cloud-cli anytime --simple
things-cloud-cli someday --simple
things-cloud-cli upcoming --simple
things-cloud-cli search "query" --simple
```

Shape:

```json
[
  {
    "uuid": "string",
    "title": "string",
    "status": "open|completed|canceled|trashed"
  }
]
```

Rules:

- `uuid` is the task identifier to pass to write commands.
- `title` is the user-visible Things title.
- `status` is normalized for agents.

## CLI Task List, Full Format

Commands:

```bash
things-cloud-cli today --format full
things-cloud-cli show <task-uuid> --format full
```

Shape:

```json
[
  {
    "uuid": "string",
    "title": "string",
    "note": "string",
    "status": 0,
    "inTrash": false,
    "isProject": false,
    "schedule": 1,
    "scheduledDate": "YYYY-MM-DD",
    "deadlineDate": "YYYY-MM-DD",
    "areaIds": ["string"],
    "parentIds": ["string"]
  }
]
```

Optional fields may be omitted when empty.

## CLI Write Success

Create task:

```json
{
  "status": "created",
  "uuid": "string",
  "title": "string"
}
```

Update task:

```json
{
  "status": "updated",
  "uuid": "string"
}
```

Complete task:

```json
{
  "status": "completed",
  "uuid": "string"
}
```

Trash task:

```json
{
  "status": "trashed",
  "uuid": "string"
}
```

Move to Today:

```json
{
  "status": "moved-to-today",
  "uuid": "string"
}
```

## CLI Dry Run

Shape:

```json
{
  "status": "dry-run",
  "operation": "string",
  "items": [
    {
      "t": 0,
      "e": "Task6",
      "p": {}
    }
  ]
}
```

Rules:

- Dry-run output is a preview. It must not write to Things Cloud.
- `items` contains Things Cloud wire payloads for review/debugging.
- Agents should summarize the user-visible effect, not expose raw payloads by
  default.

## MCP Tool Result

The MCP server returns tool content as JSON text. Agents should parse the first
text content item as JSON when they need structured data.

Task list shape:

```json
[
  {
    "uuid": "string",
    "title": "string",
    "status": "open|completed|canceled|trashed",
    "view": "inbox|today|anytime|someday|upcoming",
    "scheduledDate": "YYYY-MM-DD",
    "deadlineDate": "YYYY-MM-DD"
  }
]
```

Write dry-run shape:

```json
{
  "status": "dry-run",
  "item": {
    "t": 0,
    "e": "Task6",
    "p": {}
  }
}
```

Write success shape:

```json
{
  "status": "created|completed|updated|trashed|moved-to-today",
  "uuid": "string"
}
```

## Stability Notes

- The simple task list shape is the preferred stable contract for agents.
- Full format may grow additional fields over time.
- Raw wire payload fields in dry-run output are for diagnostics and may be more
  Things-specific than agent-friendly.
- Prefer `docs/agent-cookbook.md` for workflow guidance.
