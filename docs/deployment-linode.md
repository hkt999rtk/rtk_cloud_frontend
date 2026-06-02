# Linode Artifact Deployment

This runbook describes the standalone Linode deployment path for Realtek
Connect+. It complements the existing `website-test.local` CD host; it does not
replace that test deployment.

## Deployment Model

Realtek Connect+ uses an artifact-first deployment model:

- CI builds `realtek-connect-ci-<shortsha>.tar.gz` for every PR and push.
- CI uploads every deployable bundle as a GitHub Actions artifact named
  `realtek-connect-ci-<shortsha>`.
- On `main`, tags, and manual runs, CI also uploads the bundle, checksum, and
  manifest to Linode Object Storage under `releases/realtek_connect-ci-<shortsha>/`.
- The release workflow uploads explicit release bundles to GitHub Releases and
  Linode Object Storage under `releases/realtek_connect-<version>/`.
- The Linode deploy workflow installs a selected version onto the VM.
- Rollback deploys a previous published artifact version.
- Runtime Go service logs should emit `rtk_cloud_logger` zap JSON to
  stdout/stderr for journald collection and central forwarding; see
  `docs/SERVICE_LOGGING_MIGRATION.md`.

Do not deploy production by copying a developer checkout to the VM.

## GitHub Configuration

Required secrets:

- `LINODE_TOKEN`, used by CI to create a temporary bucket-scoped Object Storage
  key for artifact upload.
- Optional `LINODE_OBJ_ACCESS_KEY_ID` and `LINODE_OBJ_SECRET_ACCESS_KEY` if a
  long-lived bucket-scoped key is preferred over temporary keys.
- `REALTEK_CONNECT_DEPLOY_HOST`
- `REALTEK_CONNECT_DEPLOY_USER`
- `REALTEK_CONNECT_DEPLOY_SSH_KEY`
- optional `REALTEK_CONNECT_DEPLOY_PORT`

Required variables:

- Optional `LINODE_OBJ_BUCKET`; required only when the Linode account has more
  than one Object Storage bucket.
- Optional `LINODE_OBJ_ENDPOINT`; inferred from the selected bucket hostname
  when omitted.
- `REALTEK_CONNECT_PUBLIC_BASE_URL`

Optional variables:

- `REALTEK_CONNECT_DEPLOY_PREFIX`, default `/opt/realtek-connect`
- `REALTEK_CONNECT_DEPLOY_ETC_DIR`, default `/etc/realtek-connect`
- `REALTEK_CONNECT_DEPLOY_SYSTEMD_DIR`, default `/etc/systemd/system`
- `REALTEK_CONNECT_DEPLOY_DATA_DIR`, default `/var/lib/realtek-connect`
- `REALTEK_CONNECT_DEPLOY_REMOTE_DIR`, default `/tmp/realtek-connect-deploy`

## CI Artifact Contract

Every CI run builds and verifies a deployment-ready bundle:

```text
realtek-connect-ci-<shortsha>
realtek-connect-ci-<shortsha>.tar.gz
realtek-connect-ci-<shortsha>.tar.gz.sha256
realtek-connect-ci-<shortsha>.object-manifest.json
```

For `main`, tags, and manual CI runs, the same files are published to Linode
Object Storage:

```text
releases/realtek_connect-ci-<shortsha>/ci-<shortsha>.tar.gz
releases/realtek_connect-ci-<shortsha>/ci-<shortsha>.tar.gz.sha256
releases/realtek_connect-ci-<shortsha>/manifest.json
```

Pull request CI never uploads to Object Storage and does not require Linode
secrets. The bundle includes the binary, `content/`, `templates/`, `static/`,
`deploy/`, `VERSION`, and release manifests. Runtime SQLite database files must
not be included.

The upload script supports two authentication modes. If
`LINODE_OBJ_ACCESS_KEY_ID` and `LINODE_OBJ_SECRET_ACCESS_KEY` are present, it
uses them directly with the Linode S3-compatible API. Otherwise it uses
`LINODE_TOKEN` to create a temporary limited-access Object Storage key for the
selected bucket, uploads through the S3-compatible API, and deletes the temporary
key before exiting. It signs S3-compatible HTTP requests itself and does not
require the AWS CLI on the runner.

## Host Bootstrap

Create a Linode VM with Ubuntu LTS or Debian stable, then install the runtime
packages:

```sh
sudo apt update
sudo apt install -y ca-certificates curl nginx sqlite3 openssl certbot python3-certbot-nginx
```

Create the runtime directories:

```sh
sudo mkdir -p /opt/realtek-connect/releases
sudo mkdir -p /etc/realtek-connect
sudo mkdir -p /var/lib/realtek-connect/backups
sudo chown -R root:root /opt/realtek-connect /etc/realtek-connect
sudo chmod 0755 /opt/realtek-connect /opt/realtek-connect/releases
sudo chmod 0750 /var/lib/realtek-connect /var/lib/realtek-connect/backups
```

Create a deployment user with non-interactive sudo rights for install,
systemd, and verification commands. A broad first version is acceptable for a
private deployment host:

```sh
sudo useradd -m -s /bin/bash realtek-deploy || true
echo 'realtek-deploy ALL=(root) NOPASSWD:ALL' | sudo tee /etc/sudoers.d/realtek-connect-deploy
sudo chmod 0440 /etc/sudoers.d/realtek-connect-deploy
sudo visudo -cf /etc/sudoers.d/realtek-connect-deploy
```

Add the public key that matches `REALTEK_CONNECT_DEPLOY_SSH_KEY` to
`/home/realtek-deploy/.ssh/authorized_keys`.

## Runtime Environment

The first deployment creates `/etc/realtek-connect/realtek-connect.env` if it
does not already exist. Review it before public launch:

```dotenv
PORT=8080
DATABASE_PATH=/var/lib/realtek-connect/connectplus.db
ANALYTICS_DATABASE_PATH=/var/lib/realtek-connect/analytics.db
SEARCH_DATABASE_PATH=/opt/realtek-connect/current/data/search.db
SEARCH_ENABLED=false
PUBLIC_BASE_URL=https://example.com
DISABLE_SEARCH_INDEXING=true
ENABLE_ASSET_FINGERPRINTS=true
ENABLE_CDN_CACHE_HEADERS=false
ANALYTICS_ENABLED=true
ADMIN_TOKEN=change-me-before-public-use
```

Set `DISABLE_SEARCH_INDEXING=false` only when the site is ready to be indexed.
Keep SQLite files in `/var/lib/realtek-connect`; they are runtime state and are
never part of release artifacts. Documentation search can use the precomputed
`/opt/realtek-connect/current/data/search.db` bundled with a release artifact.
Enable it only after `OPENAI_API_KEY` is present in the env file and the bundle
manifest confirms `search_index.included=true`.

## GoDaddy DNS

In GoDaddy DNS management:

- Create an `A` record for the selected hostname pointing to the Linode public
  IPv4 address.
- Optionally create `www` as a CNAME to the apex or canonical hostname.
- Keep TTL low during initial cutover, for example 600 seconds.

Set `REALTEK_CONNECT_PUBLIC_BASE_URL` and the runtime `PUBLIC_BASE_URL` to the
final HTTPS URL.

## nginx And TLS

The Go application stays HTTP-only on port `8080`. nginx terminates public
HTTP/HTTPS and proxies to the app over loopback; if this host is scraped by the
central Video Cloud Prometheus, allow `8080/tcp` only from the private VPC CIDR.

```nginx
server {
    listen 80;
    server_name example.com www.example.com;

    location = /metrics/prometheus {
        return 404;
    }

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

`GET /metrics/prometheus` is intended for private Prometheus scraping only. Do
not proxy it through the public nginx server block; scrape the app over the
private VPC IP and port `8080`.

Enable the site, then issue TLS with Let’s Encrypt:

```sh
sudo nginx -t
sudo systemctl reload nginx
sudo certbot --nginx -d example.com -d www.example.com
sudo systemctl status certbot.timer --no-pager
```

## First Release And Deploy

Create a release artifact:

```sh
gh workflow run "Release Bundle" \
  --repo hkt999rtk/rtk_cloud_frontend \
  --ref main \
  -f version=v0.1.0-linode-test
```

Deploy that version:

```sh
gh workflow run "Deploy Linode" \
  --repo hkt999rtk/rtk_cloud_frontend \
  --ref main \
  -f environment=staging \
  -f version=v0.1.0-linode-test \
  -f run_verify=true
```

After deployment:

```sh
sudo systemctl status realtek-connect --no-pager
cat /opt/realtek-connect/current-version
curl -fsS https://example.com/healthz
```

## Related Runbooks

- [deployment-promotion-rollback.md](deployment-promotion-rollback.md)
- [sqlite-backup-linode.md](sqlite-backup-linode.md)
