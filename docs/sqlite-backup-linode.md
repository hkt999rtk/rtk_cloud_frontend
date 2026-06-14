# SQLite Backup Policy

Realtek Connect+ stores runtime website data in SQLite. In Kubernetes v1 these
files live on the frontend `/data` PVC; in legacy native deployments they live
on the deployment host. Release artifacts and container images must never
include runtime database files.

Kubernetes database paths:

- `/data/connectplus.db`
- `/data/analytics.db`

Legacy native database paths:

- `/var/lib/realtek-connect/connectplus.db`
- `/var/lib/realtek-connect/analytics.db`

Default backup directory:

- `/var/lib/realtek-connect/backups`

## Backup Command

Use SQLite's online backup API instead of copying live DB files directly:

```sh
sudo install -d -m 0750 /var/lib/realtek-connect/backups
ts="$(date -u +%Y%m%dT%H%M%SZ)"
sudo sqlite3 /var/lib/realtek-connect/connectplus.db \
  ".backup '/var/lib/realtek-connect/backups/connectplus-$ts.db'"
sudo sqlite3 /var/lib/realtek-connect/analytics.db \
  ".backup '/var/lib/realtek-connect/backups/analytics-$ts.db'"
```

Verify each backup:

```sh
sudo sqlite3 "/var/lib/realtek-connect/backups/connectplus-$ts.db" "PRAGMA integrity_check;"
sudo sqlite3 "/var/lib/realtek-connect/backups/analytics-$ts.db" "PRAGMA integrity_check;"
```

The expected output is `ok`.

## Retention

Recommended first-version retention:

- keep daily backups for 14 days
- keep weekly backups for 8 weeks
- keep monthly backups for 6 months if the environment is used for real leads

Local retention cleanup can start with:

```sh
sudo find /var/lib/realtek-connect/backups -type f -name '*.db' -mtime +60 -delete
```

Adjust retention before public launch based on privacy and legal requirements.

## Optional Object Storage Copy

For production, copy verified backups to a private Linode Object Storage bucket
prefix such as `sqlite-backups/<hostname>/`. Do not use the public release
artifact prefix for database backups.

```sh
go run ./cmd/linode-object-storage put \
  --file "/var/lib/realtek-connect/backups/connectplus-$ts.db" \
  --key "sqlite-backups/$(hostname)/connectplus-$ts.db"
```

## Restore Procedure

Restores should be deliberate because they replace runtime data. For Kubernetes,
scale the frontend Deployment to zero or otherwise ensure no pod is writing the
SQLite files before replacing data on the PVC. The example below is the legacy
native host form:

```sh
sudo systemctl stop realtek-connect
sudo cp /var/lib/realtek-connect/connectplus.db \
  /var/lib/realtek-connect/connectplus.db.pre-restore
sudo cp /var/lib/realtek-connect/backups/connectplus-YYYYMMDDTHHMMSSZ.db \
  /var/lib/realtek-connect/connectplus.db
sudo sqlite3 /var/lib/realtek-connect/connectplus.db "PRAGMA integrity_check;"
sudo systemctl start realtek-connect
curl -fsS https://example.com/healthz
```

Repeat the same pattern for `analytics.db` when analytics data must be
restored.

## Operational Notes

- Backups are runtime data management, not release packaging or image creation.
- Kubernetes v1 must keep one frontend replica while SQLite is writable.
- Do not commit `.db`, `.db-wal`, `.db-shm`, or backup files to git.
- Do not include database files in `realtek-connect-<version>.tar.gz`.
- Redact lead emails and analytics data before attaching backup-derived output
  to GitHub issues.
