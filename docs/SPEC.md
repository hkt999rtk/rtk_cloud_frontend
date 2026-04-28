# Realtek Connect+ Specification

## Summary

Realtek Connect+ is an English-first B2B website for an end-to-end IoT cloud platform. It is positioned for IoT product teams evaluating device onboarding, OTA, fleet operations, mobile app enablement, insights, private deployment, and ecosystem integrations.

The first version is implemented as an HTTP-only Go application using `net/http`, `html/template`, and SQLite. It does not use npm, React, Tailwind, or any frontend build step.

The current project status is **v0.1 Marketing Foundation**. It is a working website foundation and lead-capture backend, not a complete ESP RainMaker parity website, IoT console, authentication system, or real device cloud implementation.

## Current Implementation

Implemented today:

- Go HTTP server using `net/http`.
- Runtime entrypoint under `cmd/server` with request logging, graceful shutdown, and baseline read/write/idle timeouts.
- Server-rendered pages using `html/template`.
- Static CSS with a Realtek-style white, deep navy, and blue-green/teal visual system.
- Generated hero/platform image stored in `static/assets/connectplus-hero.png`.
- Per-page title, description, canonical, Open Graph, and Twitter card metadata.
- Developer docs landing and detail pages covering Product Overview, Development, APIs, SDKs, Firmware, CLI, Deployment, and Release Notes.
- Feature overview and detail pages for Provision, OTA, Fleet Management, User Management, App SDK, Insights, Private Cloud, and Integrations, including production-grade OTA rollout detail and a structured private deployment comparison story.
- `robots.txt` and `sitemap.xml` routes for crawl and link discovery.
- Contact / early access registration form.
- SQLite lead capture through `DATABASE_PATH`, defaulting to `data/connectplus.db`.
- Protected admin lead review with filtering, pagination, and filtered CSV export when `ADMIN_TOKEN` is set.
- Plain-text `/healthz` endpoint.
- Multi-stage Docker packaging that ships the server binary plus runtime templates/static assets and defaults SQLite storage to `/data/connectplus.db` inside the container.

Current routes:

- `GET /`
- `GET /docs`
- `GET /docs/{slug}`
- `GET /features`
- `GET /features/{slug}`
- `GET /contact`
- `GET /robots.txt`
- `GET /sitemap.xml`
- `POST /contact`
- `GET /healthz`
- `GET /admin/leads`
- `GET /admin/leads.csv`
- `GET /static/...`

This implementation is enough to demonstrate the Realtek Connect+ direction and collect leads. It is not yet content-complete against ESP RainMaker's public website and documentation surface.

## Visual Direction

The site uses a modern, minimal, direct enterprise style aligned with Realtek's official web presence: white content areas, clear navigation, product-led messaging, and blue-green brand accents. Compared with a traditional corporate site, Realtek Connect+ uses more whitespace, denser feature grids, stronger calls to action, and a clear platform architecture visual.

Color system:

- Base: white and near-white surfaces.
- Text: deep navy and charcoal.
- Accent: Realtek-style blue-green / teal for buttons, links, icons, and highlights.
- Secondary: pale blue-gray panels for architecture, feature detail, and contact surfaces.

UI system:

- Header: Realtek logo slot or text wordmark, Connect+ brand, Docs, Features, Architecture, Contact.
- Buttons: solid teal primary and restrained outline secondary.
- Cards: 6-8px radius, thin border, minimal shadow.
- Hero: direct product positioning with a platform architecture poster image.

Assets:

- Generated hero/platform imagery lives in `static/assets/`.
- The generated image style must be clean B2B technology: white background, teal/blue accents, device-cloud-app-dashboard flow, no stock-photo people, no third-party marks.
- Video is optional. If a ChatGPT video generation tool is available later, the site can add a short product loop with a poster image fallback. If no video tool is available, CSS motion or static generated imagery is sufficient.

## Feature Scope

Realtek Connect+ presents the following ESP RainMaker-aligned capabilities:

- Platform Overview: device firmware, cloud backend, mobile app, dashboard, and third-party integrations.
- Provision: Wi-Fi/BLE onboarding, device binding, activation, user-device association.
- OTA: firmware upload, campaign rollout, job status, cancel/archive, version validation, and force/normal/scheduled/user-controlled/time-window rollout modes.
- Fleet Management: device registry, groups, metadata/tags, batch operations, timezone, device sharing.
- User Management: sign up, sign in, OTP verification, third-party login, password recovery/change, and account deletion framed as platform capabilities rather than current website authentication flows.
- App SDK: iOS/Android SDK, sample app, rebrand/customize path, push notifications, app publishing path.
- Insights: activation statistics, firmware distribution, logs, crash reports, reboot reasons, RSSI/memory metrics.
- Private Cloud: public evaluation versus private commercial deployment, data ownership, custom domain, regional placement, upgrade path, deployment FAQ, and commercial support positioning.
- Integrations: Alexa, Google Assistant, Matter, REST API, MQTT over TLS, webhooks.

## RainMaker Parity Gap

Status values:

- `Implemented`: working website/backend behavior exists.
- `Content Partial`: page or section exists, but parity depth is incomplete.
- `Planned`: required for RainMaker parity, not yet represented deeply enough.
- `Out of Scope for website v1`: acknowledged capability, but not intended as a working implementation in the public website.

| Area | Current Status | Gap To Close |
| --- | --- | --- |
| Platform Overview | Content Partial | Add stronger end-to-end architecture story, developer entry points, security/scalability/cost positioning, and deployment model comparison. |
| Provision | Content Partial | Add QR code, SoftAP/BLE flow, device claiming, activation, user-node association, local setup failure states, and product onboarding diagrams. |
| OTA | Implemented | `/features/ota` now covers firmware upload, extracted release metadata, model/version targeting, force/normal/scheduled/user-controlled/time-window rollouts, dynamic OTA eligibility, job detail, cancellation, archive flow, and a structured rollout strategy table. |
| Fleet Management / Admin Operations | Content Partial | Add node registration, certificate flow, device registry, groups, metadata, batch operations, node summary widgets, activation statistics, and admin console concept. |
| User Management | Content Partial | Feature content now covers sign up, sign in, OTP verification, third-party login, forgot/change password, account deletion, and account lifecycle boundaries; follow-on work can add deeper support models, session behavior, and visual diagrams. |
| End-user Smart Home Features | Planned | Add remote control, local control, scheduling, scenes, grouping, node sharing, push notifications, alerts, and mobile user workflows. |
| Mobile App SDK | Content Partial | Add iOS SDK, Android SDK, sample app, customization/rebrand, push notification, app store publishing, widgets/settings, and app developer roadmap. |
| Insights | Content Partial | Add logs, crash reports, reboot reasons, custom metrics, RSSI/memory metrics, firmware distribution, support workflows, and dashboard visuals. |
| Private Cloud / Deployment | Implemented | `/features/private-cloud` now compares public evaluation, managed private deployment, and customer-operated private regions with explicit coverage for data ownership, custom domains, regional placement, upgrade path, deployment FAQ, and production support boundaries. |
| Matter / Ecosystem Integrations | Content Partial | Add Matter ecosystem positioning, Matter Fabric concept, voice assistants, MQTT over TLS, REST APIs, webhooks, cloud-to-cloud integration, and protocol diagrams. |
| Developer Docs / APIs / SDKs / CLI | Content Partial | Docs portal structure now exists across Product Overview, Development, APIs, SDKs, Firmware, CLI, Deployment, and Release Notes; deeper implementation detail and reference content still needs follow-on work. |
| SEO / Launch Readiness | Content Partial | Metadata, sitemap, and robots now exist; remaining work includes accessibility pass, visual smoke checks, deployment packaging, and CI. |
| Real IoT Cloud Operations | Out of Scope for website v1 | The public website will describe platform capabilities; it will not implement real device provisioning, OTA delivery, user auth, or telemetry ingestion in v1. |

## Website Completion Roadmap

- **v0.1 Marketing Foundation**: current state. Working Go site, product positioning, core feature pages, generated hero asset, contact lead capture, admin lead export.
- **v0.2 RainMaker Parity Content**: expand remaining content depth so Realtek Connect+ can credibly map to ESP RainMaker's feature set across user, admin, mobile SDK, deployment, Matter, and developer documentation areas.
- **v0.3 Launch Readiness**: improve SEO, accessibility, visual assets, deployment packaging, CI, form hardening, admin usability, and server operations.
- **v1.0 Public Website Candidate**: polished public-facing site with complete parity content, Realtek brand assets, reliable contact/admin workflows, and deployment documentation.

## HTTP Interface

Routes:

- `GET /`: homepage.
- `GET /docs`: developer/documentation portal landing page.
- `GET /docs/{slug}`: developer/documentation section detail pages.
- `GET /features`: feature overview.
- `GET /features/{slug}`: feature detail pages.
- `GET /contact`: contact / early access registration form.
- `GET /robots.txt`: crawl directives for search bots.
- `GET /sitemap.xml`: sitemap covering public marketing and docs pages.
- `POST /contact`: validate and store a lead in SQLite.
- `GET /healthz`: plain-text health check.
- `GET /admin/leads`: protected lead review page, enabled only when `ADMIN_TOKEN` is set.
- `GET /admin/leads.csv`: protected CSV export, enabled only when `ADMIN_TOKEN` is set.
- `GET /static/...`: CSS and asset files.

Feature slugs:

- `provision`
- `ota`
- `fleet-management`
- `user-management`
- `app-sdk`
- `insights`
- `private-cloud`
- `integrations`

Documentation slugs:

- `product-overview`
- `development`
- `apis`
- `sdks`
- `firmware`
- `cli`
- `deployment`
- `release-notes`

Environment:

- `PORT`, default `8080`.
- `DATABASE_PATH`, default `data/connectplus.db`.
- `ADMIN_TOKEN`, optional. When set, enables protected lead review and CSV export.

Operational behavior:

- The runtime uses configured `http.Server` read, write, and idle timeouts.
- Requests are logged with method, path, response status, and elapsed time.
- `SIGINT` and `SIGTERM` trigger graceful shutdown with a bounded timeout.
- Production TLS is expected to terminate at a reverse proxy, ingress controller, or hosting platform in front of the app.

Commands:

- `go run ./cmd/server`
- `go test ./...`
- `go build -o bin/realtek-connect ./cmd/server`
- `docker build -t realtek-connect .`
- `docker run --rm -p 8080:8080 -v "$(pwd)/data:/data" realtek-connect`

## SQLite

SQLite stores website leads only. It does not store real IoT telemetry or device state.

Default database path:

```text
data/connectplus.db
```

Container default database path:

```text
/data/connectplus.db
```

Container deployment notes:

- `Dockerfile` copies the compiled server together with `templates/` and `static/`, which are runtime dependencies for page rendering and asset delivery.
- `/data` is declared as the persistent volume for SQLite-backed lead storage.
- HTTPS is intentionally out of process and should be handled by the reverse proxy or deployment platform instead of the Go app directly.

Schema:

```sql
CREATE TABLE IF NOT EXISTS leads (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  company TEXT,
  email TEXT NOT NULL,
  interest TEXT NOT NULL,
  message TEXT,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

Contact form fields:

- `name`, required.
- `company`, optional.
- `email`, required.
- `interest`, required.
- `message`, optional.

## Test Plan

- `go test ./...`
- HTTP route tests for `/`, `/docs`, `/features`, feature/detail pages, `/contact`, `/robots.txt`, and `/sitemap.xml`.
- Unknown feature slug returns 404.
- Valid contact POST writes SQLite and shows success.
- Invalid contact POST shows validation errors and does not write a lead.
- Admin lead routes require `ADMIN_TOKEN`; unauthorized requests return 401, disabled admin routes return 404.
- Manual browser checks verify desktop/mobile layout, Realtek-style white + teal palette, generated image loading, and no npm/React/Tailwind dependency.

## Definition of Done

For any future issue that changes website behavior or public content:

- `go test ./...` passes.
- `go build ./cmd/server` passes.
- Go source is formatted with `gofmt`.
- Desktop and mobile visual smoke checks confirm no blank page, missing hero asset, or horizontal overflow.
- No npm, React, Tailwind, or frontend build step is introduced.
- Any new public route is documented in both `docs/SPEC.md` and `README.md`.
- Runtime SQLite files are not committed.
- Generated website assets are stored under `static/assets/` if referenced by templates or CSS.

## Developer Issue Backlog

The following GitHub issues are the source-of-truth backlog for moving from v0.1 to v1.0.

### 1. Expand RainMaker parity feature matrix

Labels: `documentation`, `enhancement`

Goal: Complete the Realtek Connect+ vs ESP RainMaker parity matrix.

Acceptance criteria:

- Matrix covers user management, smart home features, admin operations, mobile SDK, deployment, Matter, Insights, OTA, docs, APIs, SDKs, and CLI.
- Each capability has a status and notes for website v1.
- Matrix avoids claiming real cloud operations are implemented.

### 2. Add developer and documentation portal structure

Labels: `documentation`, `enhancement`

Goal: Add a developer/docs section similar to RainMaker documentation entry points.

Acceptance criteria:

- Developer/docs structure covers Product Overview, Development, APIs, SDKs, Firmware, CLI, Deployment, and Release Notes.
- Navigation exposes a Developer or Docs entry.
- First implementation can be static Go template content.

### 3. Add user management feature content

Labels: `enhancement`

Goal: Add platform content for account lifecycle capabilities.

Acceptance criteria:

- Content covers sign up, sign in, OTP, third-party login, forgot password, change password, and delete account.
- Content clearly states these are platform capabilities, not current website authentication flows.
- Feature data and templates remain server-rendered.

### 4. Add smart home end-user feature pages

Labels: `enhancement`

Goal: Expand consumer-facing IoT workflow content.

Acceptance criteria:

- Content covers remote control, local control, scheduling, scenes, grouping, node sharing, push notifications, and alerts.
- Homepage or feature overview links to these capabilities.
- Copy uses Realtek Connect+ naming and does not copy ESP RainMaker wording.

### 5. Add admin and device operations content

Labels: `enhancement`

Goal: Present admin/device operations beyond the website lead admin page.

Acceptance criteria:

- Content covers node registration, certificates, device registry, OTA jobs, firmware images, batch operations, and statistics widgets.
- Adds or expands an admin/platform operations page or section.
- Clearly separates website lead admin from a future IoT platform admin console.

### 6. Expand OTA content to production-grade detail

Labels: `enhancement`

Goal: Make OTA parity stronger and more credible.

Acceptance criteria:

- Covers firmware upload, metadata extraction, version/model targeting, force/normal/user-controlled/time-window OTA, dynamic OTA, job details, cancel, and archive.
- Includes a visual, table, or structured explanation of rollout strategies.
- `/features/ota` remains the canonical OTA page.

### 7. Add mobile SDK and app publishing content

Labels: `documentation`, `enhancement`

Goal: Match RainMaker's mobile app development story.

Acceptance criteria:

- Covers iOS SDK, Android SDK, sample app, customization/rebrand, push notifications, and app store publishing.
- Adds clear CTAs for app developers and product teams.
- Does not introduce npm or a client-side app framework.

### 8. Add private deployment and commercial cloud content

Labels: `enhancement`

Goal: Strengthen enterprise private cloud positioning.

Acceptance criteria:

- Covers public evaluation vs private commercial deployment, data ownership, custom domain, regional deployment, upgrade path, and deployment FAQ.
- Expands `/features/private-cloud` or adds a dedicated page.
- Production TLS remains documented as a reverse proxy or deployment-platform concern.

### 9. Add Matter and ecosystem integrations content

Labels: `enhancement`

Goal: Expand integrations beyond a shallow feature list.

Acceptance criteria:

- Covers Matter ecosystem positioning, Matter Fabric concept, voice assistants, MQTT over TLS, REST APIs, and webhooks.
- Adds at least one diagram or structured table for ecosystem integration paths.
- Avoids unsupported implementation promises.

### 10. Generate additional Realtek-style visual assets

Labels: `enhancement`

Goal: Use ChatGPT image generation for production-looking website visuals.

Acceptance criteria:

- Generate and store provisioning, OTA, insights, and private cloud visuals under `static/assets/`.
- Assets match white/near-white, blue-green/teal, enterprise technology styling.
- Templates reference only workspace-local assets.

### 11. Add SEO and social sharing metadata

Labels: `enhancement`

Goal: Improve public website readiness.

Acceptance criteria:

- Add per-page title and description metadata.
- Add Open Graph and Twitter card metadata.
- Add `/robots.txt` and `/sitemap.xml`.
- Tests cover new routes.

### 12. Accessibility and responsive UX audit

Labels: `enhancement`

Goal: Improve keyboard, semantic, contrast, form, and mobile behavior.

Acceptance criteria:

- Add skip link and visible focus states.
- Improve form error semantics.
- Verify primary text and controls meet WCAG AA contrast.
- Desktop/mobile checks show no horizontal overflow.

### 13. Harden contact form and SQLite lead capture

Labels: `enhancement`

Goal: Make lead capture less fragile and less spam-prone.

Acceptance criteria:

- Add max length validation.
- Add honeypot spam field.
- Add basic in-memory rate limiting.
- Tests cover valid, invalid, oversized, and spam submissions.

### 14. Improve admin leads management

Labels: `enhancement`

Goal: Make lead admin usable beyond a basic table.

Acceptance criteria:

- Add pagination.
- Add filtering/search by email, company, and interest.
- CSV export respects active filters.
- `ADMIN_TOKEN` protection remains required.

### 15. Add deployment packaging

Labels: `documentation`, `enhancement`

Goal: Make the Go site easy to run outside local development.

Acceptance criteria:

- Add Dockerfile or equivalent deployment recipe.
- Document reverse proxy TLS assumption.
- Document persistent SQLite volume.
- `go build` remains supported.

### 16. Add GitHub Actions CI

Labels: `enhancement`

Goal: Run checks on push and pull requests.

Acceptance criteria:

- Workflow runs `go test ./...`.
- Workflow runs `go build ./cmd/server`.
- Workflow fails on unformatted Go code.

### 17. Add server timeouts, request logging, and graceful shutdown

Labels: `enhancement`

Goal: Improve backend operational baseline.

Acceptance criteria:

- Use configured `http.Server` with read, write, and idle timeouts.
- Add simple request logging middleware.
- Handle interrupt for graceful shutdown.
- Existing routes and tests continue passing.

### 18. Add visual regression smoke checks

Labels: `enhancement`

Goal: Catch blank pages, missing hero assets, and mobile overflow.

Acceptance criteria:

- Add a documented smoke check command or script.
- Check homepage desktop and mobile widths.
- Verify hero image loads and no horizontal overflow.
- Avoid adding npm project dependencies.

## First-Version Limits

- HTTP only; TLS can be handled later by reverse proxy or deployment platform.
- Product pages describe platform capabilities; they do not perform real IoT cloud operations.
- Realtek official logo is expected to be provided by the user. Until then the site uses a text wordmark.
