#!/usr/bin/env bash
set -euo pipefail

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
bundle="$dist_dir/realtek-connect-$version.tar.gz"
checksum="$bundle.sha256"
manifest="$dist_dir/realtek-connect-$version.object-manifest.json"
artifact_name="realtek_connect"
object_bundle="$version.tar.gz"
object_checksum="$object_bundle.sha256"

test -f "$bundle"
test -f "$checksum"
test -f "$manifest"

(
  cd "$dist_dir"
  shasum -a 256 -c "realtek-connect-$version.tar.gz.sha256"
)

expected_source_commit="${GITHUB_SHA:-$(git -C "$repo_root" rev-parse HEAD 2>/dev/null || echo unknown)}"
python3 - "$manifest" "$version" "$expected_source_commit" <<'PY'
import json
import sys
from pathlib import Path

manifest_path, expected_version, expected_source_commit = sys.argv[1:]
manifest = json.loads(Path(manifest_path).read_text(encoding="utf-8"))
required = {
    "artifact_name",
    "artifact_path",
    "artifact_type",
    "bundle",
    "created_at",
    "repo",
    "run_attempt",
    "run_id",
    "sha256",
    "source_commit",
    "version",
    "workflow",
}
missing = sorted(required - set(manifest))
if missing:
    raise SystemExit(f"object manifest missing keys: {', '.join(missing)}")
if manifest["artifact_type"] != "release-bundle":
    raise SystemExit("object manifest artifact_type must be release-bundle")
if manifest["version"] != expected_version:
    raise SystemExit("object manifest version does not match requested upload version")
if manifest["source_commit"] != expected_source_commit:
    raise SystemExit(
        "object manifest source_commit does not match this workflow commit: "
        f"{manifest['source_commit']} != {expected_source_commit}"
    )
expected_path = f"releases/{manifest['artifact_name']}-{expected_version}/{manifest['bundle']}"
if manifest["artifact_path"] != expected_path:
    raise SystemExit("object manifest artifact_path does not match release prefix")
PY

if [[ "${LINODE_UPLOAD_DRY_RUN:-false}" == "true" || "${LINODE_UPLOAD_DRY_RUN:-}" == "1" ]]; then
  echo "dry-run: validated $bundle, $checksum, and $manifest"
  exit 0
fi

linode_api="${LINODE_API_URL:-https://api.linode.com/v4}"
temp_key_id=""
bucket_region="${LINODE_OBJ_REGION:-${AWS_DEFAULT_REGION:-us-sea}}"

delete_temp_key() {
  if [[ -n "$temp_key_id" && -n "${LINODE_TOKEN:-}" ]]; then
    curl -fsS -X DELETE \
      -H "Authorization: Bearer $LINODE_TOKEN" \
      "$linode_api/object-storage/keys/$temp_key_id" >/dev/null || true
  fi
}
trap delete_temp_key EXIT

linode_api_json() {
  local method="$1"
  local path="$2"
  local data="${3:-}"
  local response_file status
  response_file="$(mktemp)"
  if [[ -n "$data" ]]; then
    status="$(curl -sS -o "$response_file" -w '%{http_code}' -X "$method" \
      -H "Authorization: Bearer $LINODE_TOKEN" \
      -H "Content-Type: application/json" \
      --data "$data" \
      "$linode_api$path")"
  else
    status="$(curl -sS -o "$response_file" -w '%{http_code}' -X "$method" \
      -H "Authorization: Bearer $LINODE_TOKEN" \
      "$linode_api$path")"
  fi
  if [[ "$status" -lt 200 || "$status" -ge 300 ]]; then
    printf 'Linode API %s %s failed with HTTP %s\n' "$method" "$path" "$status" >&2
    cat "$response_file" >&2
    printf '\n' >&2
    rm -f "$response_file"
    return 1
  fi
  cat "$response_file"
  rm -f "$response_file"
}

if [[ -z "${AWS_ACCESS_KEY_ID:-}" || -z "${AWS_SECRET_ACCESS_KEY:-}" || -z "${LINODE_OBJ_BUCKET:-}" || -z "${LINODE_OBJ_ENDPOINT:-}" ]]; then
	test -n "${LINODE_TOKEN:-}"
	command -v curl >/dev/null 2>&1
	command -v python3 >/dev/null 2>&1

	buckets_json="$(linode_api_json GET /object-storage/buckets)"
	bucket_info="$(BUCKETS_JSON="$buckets_json" DESIRED_BUCKET="${LINODE_OBJ_BUCKET:-}" python3 <<'PY'
import json
import os
import sys

payload = json.loads(os.environ["BUCKETS_JSON"])
buckets = payload.get("data", [])
desired = os.environ.get("DESIRED_BUCKET", "")
if not desired:
    if len(buckets) != 1:
        print(f"LINODE_OBJ_BUCKET is required when the account has {len(buckets)} Object Storage buckets", file=sys.stderr)
        for bucket in buckets:
            label = bucket.get("label", "")
            if label:
                print(label, file=sys.stderr)
        sys.exit(2)
    selected = buckets[0]
else:
    selected = next((bucket for bucket in buckets if bucket.get("label") == desired), None)
    if selected is None:
        print(f"LINODE_OBJ_BUCKET {desired!r} was not found in Linode Object Storage buckets", file=sys.stderr)
        for bucket in buckets:
            label = bucket.get("label", "")
            if label:
                print(label, file=sys.stderr)
        sys.exit(2)

print(selected.get("label", ""))
region = selected.get("region") or selected.get("cluster") or ""
if region.count("-") >= 2 and region.rsplit("-", 1)[-1].isdigit():
    region = region.rsplit("-", 1)[0]
print(region)
print(selected.get("hostname") or selected.get("s3_endpoint") or "")
PY
	)"
	LINODE_OBJ_BUCKET="$(printf '%s\n' "$bucket_info" | sed -n '1p')"
	bucket_region="$(printf '%s\n' "$bucket_info" | sed -n '2p')"
	bucket_hostname="$(printf '%s\n' "$bucket_info" | sed -n '3p')"
	test -n "$bucket_region"
	if [[ -z "${LINODE_OBJ_ENDPOINT:-}" ]]; then
		test -n "$bucket_hostname"
		LINODE_OBJ_ENDPOINT="https://${bucket_hostname#${LINODE_OBJ_BUCKET}.}"
	fi

	key_label="realtek-connect-$(date -u +%Y%m%d%H%M%S)"
	key_payload="$(python3 - "$key_label" "$LINODE_OBJ_BUCKET" "$bucket_region" <<'PY'
import json
import sys

label, bucket, region = sys.argv[1:]
print(json.dumps({
    "label": label,
    "bucket_access": [{
        "bucket_name": bucket,
        "region": region,
        "permissions": "read_write",
    }],
}, separators=(",", ":")))
PY
	)"
	key_json="$(linode_api_json POST /object-storage/keys "$key_payload")"
	key_info="$(KEY_JSON="$key_json" python3 <<'PY'
import json
import os

payload = json.loads(os.environ["KEY_JSON"])
print(payload.get("id", ""))
print(payload.get("access_key", ""))
print(payload.get("secret_key", ""))
PY
	)"
	temp_key_id="$(printf '%s\n' "$key_info" | sed -n '1p')"
	AWS_ACCESS_KEY_ID="$(printf '%s\n' "$key_info" | sed -n '2p')"
	AWS_SECRET_ACCESS_KEY="$(printf '%s\n' "$key_info" | sed -n '3p')"
	test -n "$AWS_ACCESS_KEY_ID"
	test -n "$AWS_SECRET_ACCESS_KEY"
	export AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY LINODE_OBJ_BUCKET LINODE_OBJ_ENDPOINT
fi

test -n "${AWS_ACCESS_KEY_ID:-}"
test -n "${AWS_SECRET_ACCESS_KEY:-}"
test -n "${LINODE_OBJ_BUCKET:-}"
test -n "${LINODE_OBJ_ENDPOINT:-}"
test -n "$bucket_region"
command -v curl >/dev/null 2>&1

object_checksum_file="$(mktemp)"
trap 'rm -f "$object_checksum_file"; delete_temp_key' EXIT
bundle_sha="$(shasum -a 256 "$bundle" | awk '{print $1}')"
printf '%s  %s\n' "$bundle_sha" "$object_bundle" >"$object_checksum_file"

prefix="s3://$LINODE_OBJ_BUCKET/releases/$artifact_name-$version"

s3_put() {
  local file="$1"
  local key="$2"
  local endpoint_no_scheme="${LINODE_OBJ_ENDPOINT#https://}"
  endpoint_no_scheme="${endpoint_no_scheme#http://}"
  endpoint_no_scheme="${endpoint_no_scheme%%/}"
  endpoint_no_scheme="${endpoint_no_scheme#${LINODE_OBJ_BUCKET}.}"
  local scheme="https"
  if [[ "$LINODE_OBJ_ENDPOINT" == http://* ]]; then
    scheme="http"
  fi
  local host="$LINODE_OBJ_BUCKET.$endpoint_no_scheme"
  local url="$scheme://$host/$key"
  local encoded_key
  encoded_key="$(python3 - "$key" <<'PY'
import sys
from urllib.parse import quote
print("/".join(quote(part, safe="") for part in sys.argv[1].split("/")))
PY
)"
  local signed_headers
  signed_headers="$(python3 - "$file" "$host" "/$encoded_key" "$bucket_region" <<'PY'
import datetime
import hashlib
import hmac
import os
import sys

path, host, canonical_uri, region = sys.argv[1:]
access_key = os.environ["AWS_ACCESS_KEY_ID"]
secret_key = os.environ["AWS_SECRET_ACCESS_KEY"]
service = "s3"
payload_hash = hashlib.sha256(open(path, "rb").read()).hexdigest()
now = datetime.datetime.now(datetime.timezone.utc)
amz_date = now.strftime("%Y%m%dT%H%M%SZ")
date_stamp = now.strftime("%Y%m%d")
canonical_headers = (
    f"host:{host}\n"
    f"x-amz-content-sha256:{payload_hash}\n"
    f"x-amz-date:{amz_date}\n"
)
signed_headers = "host;x-amz-content-sha256;x-amz-date"
canonical_request = "\n".join([
    "PUT",
    canonical_uri,
    "",
    canonical_headers,
    signed_headers,
    payload_hash,
])
credential_scope = f"{date_stamp}/{region}/{service}/aws4_request"
string_to_sign = "\n".join([
    "AWS4-HMAC-SHA256",
    amz_date,
    credential_scope,
    hashlib.sha256(canonical_request.encode()).hexdigest(),
])
def sign(key, msg):
    return hmac.new(key, msg.encode(), hashlib.sha256).digest()
signing_key = sign(sign(sign(sign(("AWS4" + secret_key).encode(), date_stamp), region), service), "aws4_request")
signature = hmac.new(signing_key, string_to_sign.encode(), hashlib.sha256).hexdigest()
authorization = (
    f"AWS4-HMAC-SHA256 Credential={access_key}/{credential_scope}, "
    f"SignedHeaders={signed_headers}, Signature={signature}"
)
print(f"x-amz-date:{amz_date}")
print(f"x-amz-content-sha256:{payload_hash}")
print(f"authorization:{authorization}")
PY
)"
  local amz_date payload_hash authorization
  amz_date="$(awk -F: '$1 == "x-amz-date" {print substr($0, index($0, ":") + 1)}' <<<"$signed_headers")"
  payload_hash="$(awk -F: '$1 == "x-amz-content-sha256" {print substr($0, index($0, ":") + 1)}' <<<"$signed_headers")"
  authorization="$(awk -F: '$1 == "authorization" {print substr($0, index($0, ":") + 1)}' <<<"$signed_headers")"
  curl -fsS -X PUT \
    -H "Host: $host" \
    -H "x-amz-date: $amz_date" \
    -H "x-amz-content-sha256: $payload_hash" \
    -H "Authorization: $authorization" \
    --upload-file "$file" \
    "$url" >/dev/null
}

s3_put "$bundle" "releases/$artifact_name-$version/$object_bundle"
s3_put "$object_checksum_file" "releases/$artifact_name-$version/$object_checksum"
s3_put "$manifest" "releases/$artifact_name-$version/manifest.json"

echo "uploaded realtek-connect $version to $prefix"
