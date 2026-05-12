#!/usr/bin/env bash
set -euo pipefail

target="${1:-}"
if [[ -z "$target" ]]; then
  echo "usage: $0 <release-dir|release-tar.gz>" >&2
  exit 2
fi

cleanup=""
if [[ -f "$target" ]]; then
  cleanup="$(mktemp -d)"
  tar -C "$cleanup" -xzf "$target"
  release_dir="$(find "$cleanup" -mindepth 1 -maxdepth 1 -type d | head -n 1)"
else
  release_dir="$target"
fi
trap 'if [[ -n "${cleanup:-}" ]]; then rm -rf "$cleanup"; fi' EXIT

fail() {
  echo "check-release failed: $*" >&2
  exit 1
}

[[ -d "$release_dir" ]] || fail "release directory not found: $release_dir"
[[ -x "$release_dir/bin/realtek-connect" ]] || fail "missing executable bin/realtek-connect"
[[ -f "$release_dir/VERSION" ]] || fail "missing VERSION"
[[ -f "$release_dir/release-manifest.json" ]] || fail "missing release-manifest.json"
[[ -d "$release_dir/content" ]] || fail "missing content/"
[[ -d "$release_dir/templates" ]] || fail "missing templates/"
[[ -d "$release_dir/static" ]] || fail "missing static/"
[[ -d "$release_dir/deploy" ]] || fail "missing deploy/"
[[ -x "$release_dir/deploy/install.sh" ]] || fail "missing executable deploy/install.sh"
[[ -x "$release_dir/deploy/verify.sh" ]] || fail "missing executable deploy/verify.sh"
[[ -d "$release_dir/data" ]] || fail "missing data/"

if find "$release_dir" -type f \( -name '*.db' -o -name '*.db-shm' -o -name '*.db-wal' \) | grep -q .; then
  fail "release bundle must not include SQLite runtime DB files"
fi

python3 - "$release_dir/release-manifest.json" "$release_dir/VERSION" <<'PY'
import json
import sys
from pathlib import Path

manifest = json.loads(Path(sys.argv[1]).read_text())
version = Path(sys.argv[2]).read_text().strip()
required = {"artifact", "binary_sha256", "created_at", "includes", "source_commit", "version"}
missing = sorted(required - set(manifest))
if missing:
    raise SystemExit(f"manifest missing keys: {', '.join(missing)}")
if manifest["version"] != version:
    raise SystemExit("manifest version does not match VERSION")
if not isinstance(manifest["includes"], list) or not manifest["includes"]:
    raise SystemExit("manifest includes must be a non-empty list")
PY

echo "release bundle ok: $release_dir"
