#!/usr/bin/env bash
set -euo pipefail

bin="${THINGS_CLOUD_CLI:-things-cloud-cli}"

echo "Checking help without credentials..."
env -u THINGS_USERNAME -u THINGS_PASSWORD -u THINGS_TOKEN -u THINGS_CONFIG "$bin" --help >/dev/null

if [[ -z "${THINGS_USERNAME:-}" && -z "${THINGS_CONFIG:-}" ]]; then
  echo "No Things credentials configured. Help check passed; skipping Cloud checks."
  exit 0
fi

echo "Checking compact Today JSON..."
"$bin" today --simple >/tmp/things-cloud-cli-today.json

echo "Checking dry-run write..."
"$bin" create "Agent smoke test dry run" --when today --dry-run >/tmp/things-cloud-cli-dry-run.json

echo "Smoke test completed."
echo "Today output: /tmp/things-cloud-cli-today.json"
echo "Dry-run output: /tmp/things-cloud-cli-dry-run.json"
