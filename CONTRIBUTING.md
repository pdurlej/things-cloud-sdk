# Contributing

Thanks for improving this maintained Things Cloud SDK fork.

This project is focused on safe automation around Things Cloud:

- stable CLI and MCP JSON contracts
- dry-run behavior for agent-generated writes
- typed sync changes and persistent sync state
- read-only local Things SQLite inspection
- careful Things Cloud wire-format compatibility

## Before Opening a PR

Run:

```bash
go test ./...
```

For CLI write changes, also verify at least one `--dry-run` command and include
the sanitized output shape in the PR description.

Before changing raw write payloads, read `docs/client-side-bugs.md`. Small
wire-format changes can break Things.app sync behavior.

## Pull Request Notes

Good PR descriptions include:

- the user-facing behavior changed
- CLI or MCP command examples, when relevant
- whether JSON output contracts changed
- the tests or smoke checks that ran

Avoid unrelated refactors in behavior PRs. This repo is small enough that scoped
changes are easier to review and safer to ship.

## Secrets

Do not commit credentials, Things Cloud tokens, HAR captures, or local Things
database files. If a repro needs sensitive data, reduce it to sanitized JSON or
describe the relevant fields.
