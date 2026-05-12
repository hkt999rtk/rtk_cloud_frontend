# Linode Promotion And Rollback Runbook

Use this runbook to promote a published Realtek Connect+ artifact to a Linode
environment and roll back to a previous artifact when needed.

## Promotion Prerequisites

Do not deploy until all of these are true:

- CI on the candidate commit is green.
- `Release Bundle` produced `realtek-connect-<version>.tar.gz`.
- The bundle, checksum, and manifest exist in Linode Object Storage under
  `releases/<version>/`.
- `deploy/check-release.sh` passed for the same bundle.
- The Linode VM has nginx, TLS, systemd, SQLite storage, and the deploy user
  configured according to [deployment-linode.md](deployment-linode.md).
- `REALTEK_CONNECT_PUBLIC_BASE_URL` points at the public HTTPS URL.

## Deploy

Run the manual workflow:

```sh
gh workflow run "Deploy Linode" \
  --repo hkt999rtk/rtk_cloud_frontend \
  --ref main \
  -f environment=production \
  -f version=<target-version> \
  -f run_verify=true
```

During the run, confirm:

- checksum validation passed
- release bundle verification passed
- SSH upload succeeded
- `realtek-connect.service` restarted
- `deploy/verify.sh` passed
- deployment summary artifact was uploaded

On the host, confirm:

```sh
sudo systemctl is-active realtek-connect
cat /opt/realtek-connect/current-version
cat /opt/realtek-connect/previous-version 2>/dev/null || true
curl -fsS https://example.com/healthz
```

## Readiness Checks

The verifier checks:

- `GET /healthz` returns `ok`
- homepage contains `Realtek Connect` and `Contact Us`
- homepage does not contain stale `Contact Sales` copy
- Realtek logo static asset is reachable
- local MP4 brand film is served as `video/mp4` and is larger than 1 MB
- `current-version` exists when the deploy workflow has written it

For public launch, also manually check:

- `/robots.txt` matches the intended indexing setting
- `/sitemap.xml` returns the public hostname when indexing is enabled
- `/privacy` is reachable
- contact form writes to SQLite
- admin access requires `ADMIN_TOKEN`

## Rollback

Rollback uses the same deployment workflow with the previous known-good version.
Do not roll back by copying files from a developer machine or checking out old
git state on the VM.

Select the rollback version in this order:

1. `/opt/realtek-connect/previous-version` if the failed deployment completed
   far enough to write version files.
2. The version recorded in the deployment sign-off notes.
3. The most recent known-good GitHub Release / Linode Object Storage artifact.

Run:

```sh
gh workflow run "Deploy Linode" \
  --repo hkt999rtk/rtk_cloud_frontend \
  --ref main \
  -f environment=production \
  -f version=<rollback-version> \
  -f run_verify=true
```

After rollback:

- confirm `current-version` equals the rollback version
- confirm `/healthz` and homepage checks pass
- check `journalctl -u realtek-connect -n 200 --no-pager`
- record the failed version and rollback workflow URL

SQLite databases are not rolled back by release artifact rollback. If data
restore is needed, use [sqlite-backup-linode.md](sqlite-backup-linode.md).

## Diagnostics

Before or after rollback, collect:

```sh
sudo systemctl status realtek-connect --no-pager
sudo journalctl -u realtek-connect -n 200 --no-pager
sudo /opt/realtek-connect/current/deploy/verify.sh
ls -la /opt/realtek-connect/releases
cat /opt/realtek-connect/current-version
cat /opt/realtek-connect/previous-version 2>/dev/null || true
```

Do not paste raw environment files, `ADMIN_TOKEN`, SSH keys, object storage
secrets, or private DNS information into GitHub issues or PRs.
