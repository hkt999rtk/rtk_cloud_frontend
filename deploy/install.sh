#!/usr/bin/env bash
set -euo pipefail

# Legacy/native installer for website-test or recovery hosts. Official LKE
# deployments must not use host-local systemd installation.

release_dir="${1:-}"
prefix="${2:-/opt/realtek-connect}"
etc_dir="${3:-/etc/realtek-connect}"
systemd_dir="${4:-/etc/systemd/system}"
data_dir="${5:-/var/lib/realtek-connect}"

if [[ -z "$release_dir" ]]; then
  echo "usage: $0 <release-dir> [prefix] [etc-dir] [systemd-dir] [data-dir]" >&2
  exit 2
fi

version="$(cat "$release_dir/VERSION")"
install_dir="$prefix/releases/$version"
service_file="$systemd_dir/realtek-connect.service"
env_file="$etc_dir/realtek-connect.env"

install -d -m 0755 "$prefix/releases" "$etc_dir" "$systemd_dir"
install -d -m 0750 "$data_dir" "$data_dir/backups"
rm -rf "$install_dir"
install -d -m 0755 "$install_dir"
cp -R "$release_dir"/. "$install_dir/"
chmod 0755 "$install_dir/bin/realtek-connect" "$install_dir/deploy/"*.sh

if [[ ! -f "$env_file" ]]; then
  cat > "$env_file" <<EOF
PORT=8080
DATABASE_PATH=$data_dir/connectplus.db
ANALYTICS_DATABASE_PATH=$data_dir/analytics.db
SEARCH_DATABASE_PATH=$prefix/current/data/search.db
SEARCH_ENABLED=false
PUBLIC_BASE_URL=
DISABLE_SEARCH_INDEXING=true
ENABLE_ASSET_FINGERPRINTS=true
ENABLE_CDN_CACHE_HEADERS=false
ANALYTICS_ENABLED=true
ADMIN_TOKEN=change-me-before-public-use
EOF
  chmod 0640 "$env_file"
fi

cat > "$service_file" <<EOF
[Unit]
Description=Realtek Connect+ website
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
EnvironmentFile=$env_file
Environment=REALTEK_CONNECT_VERSION=$version
Environment=RTK_LOG_FORWARDER_JOURNAL_LABELS=service=realtek-connect,unit=realtek-connect.service,component=server
Environment=RTK_LOG_FORWARDER_NGINX_ACCESS_LABELS=service=realtek-connect,unit=nginx.service,component=nginx-access
Environment=RTK_LOG_FORWARDER_NGINX_ERROR_LABELS=service=realtek-connect,unit=nginx.service,component=nginx-error
WorkingDirectory=$prefix/current
ExecStart=$prefix/current/bin/realtek-connect
Restart=on-failure
RestartSec=5
ReadWritePaths=$data_dir
NoNewPrivileges=true

[Install]
WantedBy=multi-user.target
EOF

ln -sfn "$install_dir" "$prefix/current"
echo "installed realtek-connect $version into $install_dir"
