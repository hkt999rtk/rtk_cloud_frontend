# SQLite Backup Policy

Status: active service-local procedure.

Owner: rtk_cloud_frontend.

Last reviewed: 2026-08-31.

For coordinated workspace/LKE recovery, follow
[Core Backup and Restore](../../../docs/backup-restore.md). That v1 procedure
uses a manual maintenance window and a matched encrypted core archive with
both frontend databases, Cloud Admin, PostgreSQL, OpenBao and durable Redis
state. Restore after deployment remains fenced until verification and explicit
resume. Object payloads and independent escrow are outside its core scope.

The online SQLite/native-host commands below are **service-local alternatives**,
not the workspace backup entry point and not evidence of cross-service
consistency. Do not copy live database/WAL files independently. Workspace v1
stops all writers, archives the dedicated PVC and checks restored copies with
`PRAGMA integrity_check`.

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

Service-local retention suggestions (workspace v1 has no automatic pruning):

- keep daily backups for 14 days
- keep weekly backups for 8 weeks
- keep monthly backups for 6 months if the environment is used for real leads

Only for the legacy native-host directory, after reviewing the exact inventory
and independent recovery copy, local cleanup can start with:

```sh
sudo find /var/lib/realtek-connect/backups -type f -name '*.db' -mtime +60 -delete
```

Adjust retention before public launch based on privacy and legal requirements.

## Optional Object Storage Copy

For a service-local deployment, encrypt verified backups before transferring
them to an independently protected private bucket. The native helper example
below shows transfer syntax only, **not encryption**; do not use it to upload
plaintext production lead data. For workspace production, use the matched
encrypted archive and private destination in the workspace recovery procedure.

The legacy transfer example uses a private Linode Object Storage bucket
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
