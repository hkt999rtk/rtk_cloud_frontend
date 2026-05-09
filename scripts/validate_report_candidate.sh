#!/usr/bin/env bash
set -euo pipefail

candidate="${1:-}"
expected_path="${2:-}"

if [[ -z "$candidate" || -z "$expected_path" ]]; then
  echo "usage: $0 <candidate-file> <docs/TEST_REPORT.md|docs/READINESS_TEST_REPORT.md>" >&2
  exit 2
fi

case "$expected_path" in
  docs/TEST_REPORT.md|docs/READINESS_TEST_REPORT.md)
    ;;
  *)
    echo "report path is not allowlisted: $expected_path" >&2
    exit 1
    ;;
esac

expected_name="$(basename "$expected_path")"
candidate_name="$(basename "$candidate")"
if [[ "$candidate_name" != "$expected_name" ]]; then
  echo "candidate filename $candidate_name does not match expected $expected_name" >&2
  exit 1
fi

if [[ ! -f "$candidate" ]]; then
  echo "candidate file not found: $candidate" >&2
  exit 1
fi

required_headings=(
  "## Summary"
  "## Source Anchors"
  "## Environment"
  "## Commands"
  "## Result Summary"
  "## Detailed Results"
  "## Correctness / Behavior Gates"
  "## Coverage / Metrics"
  "## Skips And Blocks"
  "## Failures"
  "## Artifacts And Logs"
  "## Sign-off"
)

for heading in "${required_headings[@]}"; do
  if ! grep -qxF "$heading" "$candidate"; then
    echo "missing required heading: $heading" >&2
    exit 1
  fi
done

if grep -n -E 'Authorization: Bearer [A-Za-z0-9._~+/=-]+|password=[^<[:space:]]+|token=[^<[:space:]]+|BEGIN [A-Z ]*PRIVATE KEY|postgres://[^<[:space:]]+@|mysql://[^<[:space:]]+@' "$candidate"; then
  echo "candidate appears to contain unredacted credential material" >&2
  exit 1
fi

if grep -n -E '/Users/[^[:space:]]+|/home/[^[:space:]]+|/root/[^[:space:]]+' "$candidate"; then
  echo "candidate appears to contain local absolute paths" >&2
  exit 1
fi

if ! grep -q 'Overall result: PASS\|Overall result: FAIL\|Overall result: BLOCKED' "$candidate"; then
  echo "candidate has invalid overall result" >&2
  exit 1
fi
