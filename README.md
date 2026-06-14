# Realtek Connect+

Realtek Connect+ is a Go-rendered HTTP website for a Realtek-style IoT cloud platform concept. It uses `net/http`, `html/template`, SQLite, and static CSS. There is no npm, React, Tailwind, or frontend build step.

## Project Status

Current status: **v0.1 Marketing Foundation**.

This repository currently contains a working marketing website foundation, a developer docs portal structure, feature pages for provisioning, OTA, fleet management, smart home experience, user management, app SDK, insights, private cloud, and integrations, a file-backed manual surface with SDK sample application guidance, multilingual public routes for English, Traditional Chinese, and Simplified Chinese, per-page SEO/social metadata, sitemap and robots endpoints, a privacy notice, contact lead capture, SQLite storage, admin lead review with filtering and pagination, filtered CSV export, health check, and a container deployment recipe. It is not yet a complete IoT console, user authentication service, real OTA service, device provisioning backend, telemetry platform, or production mobile app package.

The App SDK, SDK docs, homepage, and manual now summarize the `rtk_cloud_client` sample ecosystem: Android and iOS Home Automation samples, WebApp Ops Lab, Linux device simulator, and PRO2 camera device demo. The website treats those samples as customer-facing SDK usage references, describes video streaming as WebRTC Video over TURN, and points to `rtk_cloud_client/docs/SAMPLE_APPLICATIONS.md` plus the sample README files as the deeper source of truth.

The homepage includes a locally hosted Realtek corporate brand film at `static/assets/realtek-brand-film.mp4`, with a generated poster image and `preload="metadata"`. The video supports brand trust after the platform architecture section and is not used as autoplay hero media.

Privacy readiness is intentionally lightweight for this prototype: `/privacy` describes contact form data, first-party SQLite analytics when `ANALYTICS_ENABLED=true`, OpenAI-backed documentation query behavior when search is enabled, the analytics event types collected, referrer-origin-only handling, ephemeral session ids, 90-day raw analytics event retention, the 24-month lead retention intent, data request handling, admin protection, no third-party analytics or advertising pixels or fingerprinting, and local video behavior. Replace the placeholder `privacy@example.com` contact before public launch and complete legal review.

The full roadmap and developer issue backlog live in [`docs/SPEC.md`](docs/SPEC.md).
The service logging migration to `rtk_cloud_logger` zap and central journald
forwarding is documented in
[`docs/SERVICE_LOGGING_MIGRATION.md`](docs/SERVICE_LOGGING_MIGRATION.md).
Website-local SQLite persistence and the low-priority Redis/cache boundary are
documented in
[`docs/PERSISTENCE_CACHE_BOUNDARIES.md`](docs/PERSISTENCE_CACHE_BOUNDARIES.md).
The implemented website HTTP API is documented in
[`docs/API_REFERENCE.md`](docs/API_REFERENCE.md), with a machine-readable
OpenAPI contract in [`docs/openapi.yaml`](docs/openapi.yaml).

Content classification terms for future project discussions:

- **Fixed Content**: content written directly in templates, hard-coded HTML, or Go code. Use this term when referring to page structure, shared UI copy, or fixed text that is not driven by content files.
- **Managed Content**: content maintained in content files, Markdown, YAML, content catalog data, or structured Go data and then rendered by templates. Use this term when referring to feature/docs/manual/localized content that can be centrally managed or validated.

Navigation conventions:

- The page footer contains a localized human-readable sitemap for public pages: platform entry points, features, developer docs, manual chapters, contact, and privacy.
- `/sitemap.xml` remains the crawler-facing XML sitemap. It is not linked as a normal website page and should not include admin, health, static, or internal API routes.
- Fixed Content owns footer layout and detail-page related navigation. Managed Content owns Markdown links, which should use root-relative public paths so they can be localized when rendered.

Tracked validation reports:

- `docs/TEST_REPORT.md`: deterministic CI / PR validation report.
- `docs/READINESS_TEST_REPORT.md`: deterministic legacy website-test readiness report. Official LKE runtime rollout evidence is produced by the workspace deployment flow.
- CI/CD generate sanitized candidates under `.artifacts/report-candidates/docs/` and upload them as the `report-candidates` artifact.
- Use the `Import Report Candidate` workflow to import only `docs/TEST_REPORT.md` or `docs/READINESS_TEST_REPORT.md` from a selected workflow run into a target PR branch or explicit branch after heading, redaction, and path validation.

## Run

```bash
go run ./cmd/server
```

Open:

```text
http://localhost:8080
```

Localized public entry points:

```text
http://localhost:8080/
http://localhost:8080/zh-tw/
http://localhost:8080/zh-cn/
```

## Configuration

Environment variables:

- `PORT`: HTTP port, default `8080`.
- `DATABASE_PATH`: SQLite database path, default `data/connectplus.db`.
- `ANALYTICS_ENABLED`: analytics toggle, default `true`.
- `ANALYTICS_DATABASE_PATH`: SQLite analytics database path, default `data/analytics.db`.
- `ANALYTICS_RETENTION_DAYS`: raw analytics retention window, default `90`.
- `ADMIN_TOKEN`: enables protected lead viewing and CSV export.
- `OPENAI_API_KEY`: required by `cmd/search-index` and by runtime search answering when `SEARCH_ENABLED=true`.
- `SEARCH_DATABASE_PATH`: SQLite documentation search database path, default `data/search.db`.
- `SEARCH_ENABLED`: public documentation query toggle, default `false`. Enable only after `search.db` has been built and `OPENAI_API_KEY` is configured.
- `SEARCH_EMBEDDING_MODEL`: OpenAI embedding model, default `text-embedding-3-small`.
- `SEARCH_ANSWER_MODEL`: OpenAI answer model, default `gpt-4.1-mini`.
- `DISABLE_SEARCH_INDEXING`: set to `true` on private/test deployments to emit `X-Robots-Tag: noindex, nofollow, noarchive`, add page-level `robots` meta tags, disallow all crawling in `/robots.txt`, and hide `/sitemap.xml`.
- `PUBLIC_BASE_URL`: optional public origin such as `https://webtest.mgmeet.io`. When empty, canonical URLs, social image URLs, `hreflang`, robots sitemap references, and sitemap locations are built from the incoming request host and forwarded headers.
- `ENABLE_ASSET_FINGERPRINTS`: optional. Set to `true` to append content hashes to template-rendered static URLs, for example `/static/styles.css?v=<hash>`.
- `ENABLE_CDN_CACHE_HEADERS`: optional. Set to `true` to emit CDN-friendly cache headers for static assets, public HTML, admin/contact POST responses, health, robots, and sitemap.

Runtime behavior:

- HTTP server uses read, write, and idle timeouts for a safer default operational baseline.
- `SIGINT` and `SIGTERM` trigger graceful shutdown with a bounded drain window before exit.

Search indexing:

- Test and preview deployments should run with `DISABLE_SEARCH_INDEXING=true`.
- Public launch deployments can omit it when the site is approved for indexing.

Documentation query:

- Build the local website-only search index with `OPENAI_API_KEY=... go run ./cmd/search-index`.
- `deploy/package.sh` builds a precomputed `data/search.db` into the release bundle when `OPENAI_API_KEY` is available. If the key is absent, packaging succeeds without the index so normal PR validation is not blocked.
- Start the server with `SEARCH_ENABLED=true` only after `SEARCH_DATABASE_PATH` points at a built index and `OPENAI_API_KEY` is configured.
- Search answers are generated only when retrieved website content clears the relevance threshold. No matching documentation returns a controlled no-hit answer without calling the answer model.
- Search query text is sent to OpenAI for embeddings. When sources are found, the query and retrieved snippets are sent to the OpenAI Responses API to generate a source-grounded answer. Raw query text is not stored in the analytics event payload.

Persistence and cache boundary:

- Contact leads, first-party analytics, and the local documentation search index remain concrete SQLite repositories unless a measured website bottleneck or concrete operational need justifies different cache work.
- Do not add Redis or Redis-compatible caching to this repository only because the broader platform may use it.
- Website SQLite stores must not hold authoritative IoT telemetry, product device state, customer account state, fleet data, OTA execution state, or production mobile app user state.
- CDN/static cache headers remain separate from application Redis/cache decisions.

CDN readiness:

- CDN-specific behavior is implemented but disabled by default.
- Keep `PUBLIC_BASE_URL`, `ENABLE_ASSET_FINGERPRINTS`, and `ENABLE_CDN_CACHE_HEADERS` unset until a CDN is ready to sit in front of the site.
- When enabling a CDN later, set `PUBLIC_BASE_URL` to the public HTTPS origin, enable asset fingerprints, enable CDN cache headers, and configure the CDN to cache `/static/*` aggressively while leaving HTML, `/contact`, `/admin/*`, and `/healthz` uncacheable.
- This project does not assume a specific CDN provider. Configure TLS, cache rules, purge behavior, compression, and origin forwarding in the chosen CDN or reverse proxy.

## Routes

- `GET /`
- `GET /docs`
- `GET /docs/{slug}`
- `GET /manual`
- `GET /manual/{slug}`
- `GET /features`
- `GET /features/{slug}`
- `GET /contact`
- `GET /privacy`
- `GET /search`
- `GET /robots.txt`
- `GET /sitemap.xml`
- `POST /contact`
- `POST /api/search`
- `GET /healthz`
- `GET /admin/leads`, requires `ADMIN_TOKEN`
- `GET /admin/leads.csv`, requires `ADMIN_TOKEN`
- `GET /static/...`

New public routes should be documented here, in `docs/SPEC.md`, and in
`docs/openapi.yaml` when they expose an API or form contract.

## Multilingual Public Site

The public website supports:

- English: unprefixed canonical URLs, for example `/features/ota`.
- Traditional Chinese: `/zh-tw/...`, for example `/zh-tw/features/ota`.
- Simplified Chinese: `/zh-cn/...`, for example `/zh-cn/features/ota`.

Feature and docs slugs stay in English across all locales. This keeps links stable and avoids translating route identifiers. Public pages emit canonical URLs plus `hreflang` alternates for `en`, `zh-Hant`, `zh-Hans`, and `x-default`. `/sitemap.xml` includes all localized public URLs when search indexing is enabled.

The admin, health, and static routes are not localized.

## Content Authoring (Docs Placeholder, Issue #75)

This project is adding a lightweight docs content source in two phases:

- Phase 1 (this issue): `/docs` uses content files only, as a layout/rendering preview.
- Phase 2 (later issues): expand the same model to other docs pages and features.

Phase 1 conventions:

- Content format: YAML frontmatter + Markdown body.
- File location: `content/docs/<locale>/docs.yaml` (locale examples: `en`, `zh-TW`, `zh-CN`).
- The Markdown body supports embedded media links, for example `![alt text](/static/assets/...)`.
- Reload behavior: content is cached at startup, with a manual admin reload endpoint.

Admin helper for development/test refresh:

- `POST /admin/reload-content` (requires `ADMIN_TOKEN`) clears cached content and reloads from files.
- Existing routes and locales remain unchanged for all other pages in this issue.

Admin requests can authenticate with either:

```bash
curl -H "X-Admin-Token: $ADMIN_TOKEN" http://localhost:8080/admin/leads
```

or:

```bash
curl "http://localhost:8080/admin/leads.csv?token=$ADMIN_TOKEN"
```

## Test

```bash
go test ./...
```

## Visual Smoke Check

```bash
go run ./cmd/visual-smoke
```

Notes:

- The command starts the website in-process by default and checks English, Traditional Chinese, and Simplified Chinese public pages at desktop and mobile widths.
- The check covers localized homepages plus representative localized feature detail pages, verifies key images load, and fails on horizontal overflow.
- A local Chrome install is required. Override detection with `CHROME_PATH=/path/to/chrome` or `go run ./cmd/visual-smoke -chrome-path /path/to/chrome`.
- Use `go run ./cmd/visual-smoke -base-url http://localhost:8080` to target an already running server instead of the in-process test server.

## Build

```bash
go build -o bin/realtek-connect ./cmd/server
```

## Kubernetes Deployment

Build the container image:

```bash
docker build -t realtek-connect .
```

Run it with a persistent SQLite volume:

```bash
docker run --rm -p 8080:8080 \
  -e ADMIN_TOKEN=change-me \
  -v "$(pwd)/data:/data" \
  realtek-connect
```

Kubernetes deployment notes:

- The image keeps application state only in SQLite under `/data/connectplus.db`; mount `/data` to persist leads across restarts.
- When analytics is enabled, the app also keeps first-party event data in `/data/analytics.db`.
- The official LKE deployment entry is `rtk_cloud_workspace`; it builds this repository's `Dockerfile`, exports `LKE_FRONTEND_IMAGE`, and deploys that image into the `<stack>-frontend` namespace.
- Keep Kubernetes `replicas: 1` while leads and analytics use SQLite. Do not horizontally scale the frontend until persistence moves to an external database or another multi-writer-safe design.
- Mount `/data` from a PVC and set `DATABASE_PATH=/data/connectplus.db`, `ANALYTICS_DATABASE_PATH=/data/analytics.db`, and `SEARCH_DATABASE_PATH=/data/search.db` or an immutable bundled search index path.
- The container serves HTTP on port `8080` and exposes `/healthz`. TLS termination, public routing, DNS, and certificates are handled by Ingress/Gateway, Linode NodeBalancer, cert-manager, and the selected CDN or reverse proxy layer.
- The app logs to stdout/stderr; Kubernetes log collection and central forwarding are deployment concerns.
- Native builds and release bundles remain supported only for local diagnostics, legacy website-test validation, or non-K8s recovery use. They are not the official staging or production rollout path.

Kubernetes runbooks:

- `docs/deployment-k8s.md`
- `docs/deployment-promotion-rollback.md`
- `docs/sqlite-backup-linode.md`

Legacy/native artifact tools:

- `deploy/package.sh <version>` creates `dist/realtek-connect-<version>.tar.gz`, a checksum, and a manifest. When `OPENAI_API_KEY` is present, the bundle includes `data/search.db` as a precomputed documentation search index.
- `deploy/check-release.sh`, `deploy/install.sh`, and `deploy/verify.sh` validate or install native bundles for legacy website-test or recovery environments.
- `deploy/upload-linode-artifact.sh <version>` remains available for archiving native bundles to Linode Object Storage, but Object Storage bundles are not the LKE runtime rollout artifact.
