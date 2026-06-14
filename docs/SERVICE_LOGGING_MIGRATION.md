# Service Logging Migration

Status: implemented.

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
- Keep Kubernetes runtime logs on stdout/stderr for cluster log collection.
- Treat host/nginx forwarder labels as legacy native deployment compatibility.
- Do not log lead form raw payloads, cookies, OpenAI/API keys, SMTP secrets, or
  SQLite connection details with credentials.

## Entrypoints To Cover

- `cmd/server`
- `cmd/search-index`
- `deploy/install.sh` legacy native unit generation
- `deploy/*.service` templates or generated units, when native artifacts are used

## Forwarder Labels

The LKE path should collect the container's stdout/stderr records and attach
Kubernetes metadata such as namespace, pod, container, image, and deployment
labels. The app itself must not depend on host-local journald or reverse-proxy
logs.

For legacy native artifacts, the install script writes these low-cardinality
label sets into the generated `realtek-connect.service` unit for the host log
forwarder:

- Go runtime journald records: `service=realtek-connect`,
  `unit=realtek-connect.service`, `component=server`
- legacy nginx access log records: `service=realtek-connect`, `unit=nginx.service`,
  `component=nginx-access`
- legacy nginx error log records: `service=realtek-connect`, `unit=nginx.service`,
  `component=nginx-error`

High-cardinality values such as request ids, trace ids, paths, and remote
addresses remain structured log fields, not default forwarder labels.

## Acceptance Criteria

- Kubernetes frontend pods emit JSON zap logs from Go processes on stdout/stderr.
- Legacy `realtek-connect.service` emits the same JSON zap logs when native
  artifacts are used.
- HTTP request logs include status, latency, sanitized path, remote address,
  and request id.
- Search-index runs can be traced by `service`, `component`, and build version.
- `go test ./...` and frontend build/test checks continue to pass.
