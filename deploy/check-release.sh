#!/usr/bin/env bash
set -euo pipefail

# Validates legacy/native release bundles for diagnostics or recovery use.
# It is not the official Kubernetes rollout gate.

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
[[ -f "$release_dir/dist/sdk-docs/current/html/index.html" ]] || fail "missing SDK manual HTML"
[[ -f "$release_dir/dist/sdk-docs/current/pdf/rtk-cloud-sdk-user-manual.pdf" ]] || fail "missing SDK manual PDF"
[[ -f "$release_dir/dist/sdk-docs/current/manifest.json" ]] || fail "missing SDK documentation manifest"

while IFS= read -r sqlite_file; do
  rel="${sqlite_file#"$release_dir"/}"
  case "$rel" in
    data/search.db)
      ;;
    *)
      fail "release bundle must not include SQLite runtime DB files: $rel"
      ;;
  esac
done < <(find "$release_dir" -type f \( -name '*.db' -o -name '*.db-shm' -o -name '*.db-wal' \))

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
search_index = manifest.get("search_index", {"included": False})
if search_index.get("included"):
    if search_index.get("path") != "data/search.db":
        raise SystemExit("manifest search_index path must be data/search.db")
    if "data/search.db" not in manifest["includes"]:
        raise SystemExit("manifest includes must list data/search.db when search index is included")
    search_db = Path(sys.argv[1]).parent / "data" / "search.db"
    if not search_db.is_file():
        raise SystemExit("manifest declares search index but data/search.db is missing")
PY

echo "release bundle ok: $release_dir"
