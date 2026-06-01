# Service Logging Migration

Status: implementation handoff.

Owner: `rtk_cloud_frontend`.

## Goal

Make the Realtek Connect+ frontend runtime logs compatible with the RTK Cloud
central service logger while keeping website analytics, lead storage, and
service logs as separate concerns.

## Required Changes

- Replace stdlib Go logging in `cmd/server` and `cmd/search-index` with
  `rtk_cloud_logger` zap logging.
- Emit server request logs as single-line JSON on stdout/stderr.
- Preserve `request_id` and `trace_id` for documentation/search requests when
  present.
- Log search-index and content-load failures with structured fields.
- Define forwarder labels for any non-Go web/runtime logs collected from nginx
  or deployment-managed services.
- Do not log lead form raw payloads, cookies, OpenAI/API keys, SMTP secrets, or
  SQLite connection details with credentials.

## Entrypoints To Cover

- `cmd/server`
- `cmd/search-index`
- `deploy/install.sh`
- `deploy/*.service` templates or generated units
- nginx access/error logs in Linode deployments

## Acceptance Criteria

- `realtek-connect.service` emits JSON zap logs from Go processes.
- HTTP request logs include status, latency, sanitized path, remote address,
  and request id.
- Search-index runs can be traced by `service`, `component`, and build version.
- `go test ./...` and frontend build/test checks continue to pass.
