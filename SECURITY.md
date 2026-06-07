# Security Policy

This is an unofficial, reverse-engineered Things Cloud SDK. It may handle Things
Cloud credentials, local cache files, and task data.

## Supported Versions

Security fixes target the latest tagged release and `main`.

## Reporting a Vulnerability

Do not paste credentials, tokens, HAR captures, or private task data into a
public issue.

If GitHub private vulnerability reporting is available for this repository, use
it. Otherwise, open a minimal public issue without secrets and ask for a private
contact path.

If a credential may have been exposed, rotate or revoke it before debugging.

## Handling Sensitive Repros

Preferred repro formats:

- sanitized CLI commands
- sanitized dry-run JSON
- minimal redacted snippets from Things Cloud payloads
- clear expected vs actual behavior

Avoid attaching full local Things databases or network captures.
