#!/usr/bin/env bash
set -euo pipefail

# Legacy/native artifact helper. Official LKE runtime rollouts use the
# workspace-built container image exported as LKE_FRONTEND_IMAGE.

version="${1:-${VERSION:-}}"
if [[ -z "$version" ]]; then
  echo "usage: $0 <version>" >&2
  exit 2
fi
case "$version" in
  *[!A-Za-z0-9._-]* | "" | .* | *..*)
    echo "invalid version: use only letters, digits, dot, underscore, and dash" >&2
    exit 2
    ;;
esac

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
dist_dir="$repo_root/dist"
sdk_repo="${SDK_REPO:-$repo_root/../rtk_cloud_client}"
sdk_docs_version="${SDK_DOCS_VERSION:-$version}"
sdk_docs_current="$dist_dir/sdk-docs/current"
sdk_docs_ready="false"
release_name="realtek-connect-$version"
release_dir="$dist_dir/$release_name"
bundle="$dist_dir/$release_name.tar.gz"
checksum="$bundle.sha256"
object_manifest="$dist_dir/$release_name.object-manifest.json"

if [[ -d "$sdk_docs_current/html" && -d "$sdk_docs_current/pdf" && -f "$sdk_docs_current/manifest.json" ]]; then
  if python3 - "$sdk_docs_current/manifest.json" "$sdk_docs_version" <<'PY'
import json
import sys
from pathlib import Path

manifest = json.loads(Path(sys.argv[1]).read_text(encoding="utf-8"))
raise SystemExit(0 if manifest.get("sdk_version") == sys.argv[2] else 1)
PY
  then
    sdk_docs_ready="true"
  fi
fi

if [[ "$sdk_docs_ready" != "true" ]]; then
  [[ -d "$sdk_repo/.git" ]] || {
    echo "SDK documentation is not built and SDK_REPO is unavailable: $sdk_repo" >&2
    exit 1
  }
  python3 "$repo_root/tools/build_sdk_docs.py" build \
    --source "$repo_root/content/manual/sdk" \
    --sdk-repo "$sdk_repo" \
    --version "$sdk_docs_version" \
    --output "$dist_dir/sdk-docs"
fi

rm -rf "$release_dir" "$bundle" "$checksum" "$object_manifest"
mkdir -p "$release_dir/bin" "$release_dir/data"

GOOS="${GOOS:-linux}" GOARCH="${GOARCH:-amd64}" CGO_ENABLED="${CGO_ENABLED:-0}" \
  go build -o "$release_dir/bin/realtek-connect" ./cmd/server

cp -R "$repo_root/content" "$release_dir/"
cp -R "$repo_root/templates" "$release_dir/"
cp -R "$repo_root/static" "$release_dir/"
cp -R "$repo_root/deploy" "$release_dir/"
mkdir -p "$release_dir/dist"
cp -R "$dist_dir/sdk-docs" "$release_dir/dist/"
chmod 0755 "$release_dir/bin/realtek-connect" "$release_dir/deploy/"*.sh
chmod 0775 "$release_dir/data"
printf '%s\n' "$version" > "$release_dir/VERSION"

search_index_included="false"
search_index_sha=""
if [[ "${BUILD_SEARCH_INDEX:-true}" != "false" && -n "${OPENAI_API_KEY:-}" ]]; then
  echo "building precomputed documentation search index"
  SEARCH_DATABASE_PATH="$release_dir/data/search.db" \
    go run ./cmd/search-index \
      -repo-root "$repo_root" \
      -content-root "$repo_root/content" \
      -database "$release_dir/data/search.db"
  search_index_included="true"
  search_index_sha="$(shasum -a 256 "$release_dir/data/search.db" | awk '{print $1}')"
else
  echo "skipping precomputed documentation search index: OPENAI_API_KEY is not set or BUILD_SEARCH_INDEX=false"
fi

source_commit="$(git -C "$repo_root" rev-parse HEAD 2>/dev/null || echo unknown)"
created_at="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"
binary_sha="$(shasum -a 256 "$release_dir/bin/realtek-connect" | awk '{print $1}')"

python3 - "$release_dir/release-manifest.json" "$version" "$source_commit" "$created_at" "$binary_sha" "$search_index_included" "$search_index_sha" <<'PY'
import json
import sys
from pathlib import Path

out, version, source_commit, created_at, binary_sha, search_index_included, search_index_sha = sys.argv[1:]
includes = ["bin/realtek-connect", "content/", "templates/", "static/", "deploy/", "dist/sdk-docs/", "VERSION"]
search_index = {
    "included": search_index_included == "true",
    "path": "data/search.db" if search_index_included == "true" else "",
    "sha256": search_index_sha,
}
if search_index["included"]:
    includes.append("data/search.db")
manifest = {
    "artifact": f"realtek-connect-{version}",
    "binary_sha256": binary_sha,
    "created_at": created_at,
    "includes": includes,
    "search_index": search_index,
    "source_commit": source_commit,
    "version": version,
}
Path(out).write_text(json.dumps(manifest, indent=2, sort_keys=True) + "\n")
PY

tar -C "$dist_dir" -czf "$bundle" "$release_name"
bundle_sha="$(shasum -a 256 "$bundle" | awk '{print $1}')"
printf '%s  %s\n' "$bundle_sha" "$(basename "$bundle")" > "$checksum"

python3 - "$object_manifest" "$version" "$source_commit" "$created_at" "$bundle_sha" <<'PY'
import json
import os
import sys
from pathlib import Path

out, version, source_commit, created_at, bundle_sha = sys.argv[1:]
artifact_name = "realtek_connect"
bundle = f"{version}.tar.gz"
manifest = {
    "artifact_name": artifact_name,
    "artifact_path": f"releases/{artifact_name}-{version}/{bundle}",
    "artifact_type": "release-bundle",
    "bundle": bundle,
    "created_at": created_at,
    "repo": os.environ.get("GITHUB_REPOSITORY", "rtk_cloud_frontend"),
    "run_attempt": os.environ.get("GITHUB_RUN_ATTEMPT", ""),
    "run_id": os.environ.get("GITHUB_RUN_ID", ""),
    "sha256": bundle_sha,
    "source_commit": source_commit,
    "version": version,
    "workflow": os.environ.get("GITHUB_WORKFLOW", ""),
}
Path(out).write_text(json.dumps(manifest, indent=2, sort_keys=True) + "\n")
PY

"$repo_root/deploy/check-release.sh" "$release_dir"

printf '%s\n' "$bundle"
