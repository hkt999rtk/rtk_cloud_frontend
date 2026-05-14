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

test -f "$bundle"
test -f "$checksum"
test -f "$manifest"

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

if [[ -z "${AWS_ACCESS_KEY_ID:-}" || -z "${AWS_SECRET_ACCESS_KEY:-}" || -z "${LINODE_OBJ_BUCKET:-}" || -z "${LINODE_OBJ_ENDPOINT:-}" ]]; then
  test -n "${LINODE_TOKEN:-}"
  command -v curl >/dev/null 2>&1
  command -v jq >/dev/null 2>&1

  buckets_json="$(curl -fsS -H "Authorization: Bearer $LINODE_TOKEN" "$linode_api/object-storage/buckets")"
  if [[ -z "${LINODE_OBJ_BUCKET:-}" ]]; then
    bucket_count="$(jq '.data | length' <<<"$buckets_json")"
    if [[ "$bucket_count" != "1" ]]; then
      echo "LINODE_OBJ_BUCKET is required when the account has $bucket_count Object Storage buckets" >&2
      jq -r '.data[].label' <<<"$buckets_json" >&2
      exit 2
    fi
    LINODE_OBJ_BUCKET="$(jq -r '.data[0].label' <<<"$buckets_json")"
  fi

  bucket_json="$(jq -e --arg bucket "$LINODE_OBJ_BUCKET" '.data[] | select(.label == $bucket)' <<<"$buckets_json")"
  bucket_region="$(jq -r '.region // .cluster // empty' <<<"$bucket_json")"
  bucket_hostname="$(jq -r '.hostname // .s3_endpoint // empty' <<<"$bucket_json")"
  test -n "$bucket_region"
  if [[ -z "${LINODE_OBJ_ENDPOINT:-}" ]]; then
    test -n "$bucket_hostname"
    LINODE_OBJ_ENDPOINT="https://${bucket_hostname#${LINODE_OBJ_BUCKET}.}"
  fi

  key_label="realtek-connect-$version-$(date -u +%Y%m%d%H%M%S)"
  key_payload="$(jq -n \
    --arg label "$key_label" \
    --arg bucket "$LINODE_OBJ_BUCKET" \
    --arg region "$bucket_region" \
    '{label: $label, bucket_access: [{bucket_name: $bucket, region: $region, permissions: "read_write"}]}')"
  key_json="$(curl -fsS -X POST \
    -H "Authorization: Bearer $LINODE_TOKEN" \
    -H "Content-Type: application/json" \
    --data "$key_payload" \
    "$linode_api/object-storage/keys")"
  temp_key_id="$(jq -r '.id // empty' <<<"$key_json")"
  AWS_ACCESS_KEY_ID="$(jq -r '.access_key // empty' <<<"$key_json")"
  AWS_SECRET_ACCESS_KEY="$(jq -r '.secret_key // empty' <<<"$key_json")"
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

(
  cd "$dist_dir"
  shasum -a 256 -c "realtek-connect-$version.tar.gz.sha256"
)

prefix="s3://$LINODE_OBJ_BUCKET/releases/$version"

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

s3_put "$bundle" "releases/$version/realtek-connect-$version.tar.gz"
s3_put "$checksum" "releases/$version/realtek-connect-$version.tar.gz.sha256"
s3_put "$manifest" "releases/$version/manifest.json"

echo "uploaded realtek-connect $version to $prefix"
