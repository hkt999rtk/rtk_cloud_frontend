#!/usr/bin/env bash
set -euo pipefail

base_url="${REALTEK_CONNECT_VERIFY_BASE_URL:-${PUBLIC_BASE_URL:-http://127.0.0.1:8080}}"
prefix="${REALTEK_CONNECT_DEPLOY_PREFIX:-/opt/realtek-connect}"

curl_common=(curl -fsS --max-time 20)
health="$("${curl_common[@]}" "$base_url/healthz")"
[[ "$health" == "ok" ]]

home="$(mktemp)"
trap 'rm -f "$home" "$headers"' EXIT
"${curl_common[@]}" "$base_url/" > "$home"
grep -q "Realtek Connect" "$home"
grep -q "Contact Us" "$home"
if grep -q "Contact Sales" "$home"; then
  echo "stale Contact Sales copy found" >&2
  exit 1
fi

headers="$(mktemp)"
curl -fsSI --max-time 20 "$base_url/static/assets/realtek-logo.png" > "$headers"
grep -qi "content-type: image/" "$headers"

curl -fsSI --max-time 20 "$base_url/static/assets/realtek-brand-film.mp4" > "$headers"
grep -qi "content-type: video/mp4" "$headers"
awk 'BEGIN { ok = 0 } tolower($1) == "content-length:" && $2 + 0 > 1000000 { ok = 1 } END { exit ok ? 0 : 1 }' "$headers"

if [[ -f "$prefix/current-version" ]]; then
  test -s "$prefix/current-version"
fi

echo "realtek-connect verify ok: $base_url"
