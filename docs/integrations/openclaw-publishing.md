# OpenClaw and ClawHub Publishing

This repository includes a publishable Things Cloud agent skill at:

```text
skills/things-cloud/SKILL.md
```

The skill is designed for OpenClaw, ClawHub, Codex, Claude Code, and other
agent runtimes that can consume `SKILL.md` style instructions. It wraps the
maintained `things-cloud-cli` and `things-mcp` tools instead of duplicating
runtime code.

It works with OpenClaw, Codex, and Claude Code through MCP when available, with
a CLI fallback for hosts that prefer shell commands.

## Install from This Repository

After the skill is merged to `main`, OpenClaw users can install from the raw
skill URL if their OpenClaw version supports URL installs:

```bash
openclaw skills install https://raw.githubusercontent.com/pdurlej/things-cloud-sdk/main/skills/things-cloud/SKILL.md --as things-cloud
```

For local testing from a checkout:

```bash
openclaw skills install ./skills/things-cloud --as things-cloud
```

## Publish to ClawHub

Prerequisites:

```bash
npm i -g clawhub
clawhub login
```

Recommended dry-run/sync check:

```bash
clawhub sync --all --dry-run
```

Publish the skill:

```bash
clawhub skill publish ./skills/things-cloud \
  --slug things-cloud \
  --name "Things Cloud" \
  --version 0.1.1 \
  --tags latest \
  --changelog "Document OpenClaw, Codex, and Claude Code support via MCP with CLI fallback." \
  --clawscan-note "This skill installs and calls the maintained things-cloud-cli and things-mcp Go binaries. It requires Things Cloud credentials through THINGS_USERNAME plus THINGS_TOKEN or THINGS_PASSWORD. Agent writes should use dry-run before execution."
```

Expected listing copy:

```text
Things Cloud

Manage Things 3 tasks through Things Cloud with a maintained Go CLI and MCP
server. Supports safe dry-run writes, compact JSON task reads, completed/logbook
evidence, recurring tasks, and agent-friendly credentials.
```

Suggested search terms:

```text
things, things3, things-cloud, task-management, productivity, mcp, cli, agents, openclaw
```

## OpenClaw Integration Request

The OpenClaw integrations page links "Request Integration" to the
`openclaw/openclaw` issue tracker. Submit an issue after the skill is merged and
published or at least installable from `main`.

Suggested title:

```text
Add Things Cloud / Things 3 as a Productivity integration
```

Suggested body:

~~~markdown
## Integration

Things Cloud / Things 3 task management

## What it does

Adds an agent-friendly way to read and safely update Things 3 tasks through
Things Cloud.

## Links

- SDK/CLI/MCP: https://github.com/pdurlej/things-cloud-sdk
- Release: https://github.com/pdurlej/things-cloud-sdk/releases/tag/v0.2.3
- Skill: https://github.com/pdurlej/things-cloud-sdk/tree/main/skills/things-cloud

## Capabilities

- list Today, Inbox, Anytime, Someday, Upcoming
- search tasks
- create/edit/complete/trash/move tasks
- dry-run writes before user confirmation
- completed/logbook evidence for feedback loops
- recurring task creation via CLI
- stdio MCP server for agent hosts

## Setup

```bash
go install github.com/pdurlej/things-cloud-sdk/cmd/things-cloud-cli@v0.2.3
go install github.com/pdurlej/things-cloud-sdk/cmd/things-mcp@v0.2.3
export THINGS_USERNAME="you@example.com"
export THINGS_TOKEN="your-things-cloud-password-or-token-alias"
```

## Safety

The local Things SQLite adapter is read-only. Agent writes go through Things
Cloud and should use dry-run before executing user-visible changes.
```
~~~
