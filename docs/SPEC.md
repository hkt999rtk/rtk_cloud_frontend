# Realtek Connect+ Specification

## Summary

Realtek Connect+ is an English-first B2B website for an end-to-end IoT cloud platform. It is positioned for IoT product teams evaluating device onboarding, OTA, fleet operations, mobile app enablement, insights, private deployment, and ecosystem integrations.

The first version is implemented as an HTTP-only Go application using `net/http`, `html/template`, and SQLite. It does not use npm, React, Tailwind, or any frontend build step.

The current project status is **v0.1 Marketing Foundation**. It is a working website foundation and lead-capture backend, not a complete IoT console, authentication system, or real device cloud implementation.

## Current Implementation

Implemented today:

- Go HTTP server using `net/http`.
- Runtime entrypoint under `cmd/server` with request logging, graceful shutdown, and baseline read/write/idle timeouts.
- Server-rendered pages using `html/template`.
- Static CSS with a Realtek-style white, deep navy, and blue-green/teal visual system.
- Artifact-first Linode deployment path with release bundles, Linode Object Storage publication, manual VM deploy workflow, rollback runbook, and SQLite backup policy.
- Corporate hero/platform image stored in `static/assets/connectplus-hero-corporate-v2.jpg`.
- Corporate feature and platform visuals stored under `static/assets/`, including the SDK sample ecosystem diagram at `static/assets/connectplus-sample-ecosystem-corporate-v2.jpg`.
- Per-page title, description, canonical, Open Graph, and Twitter card metadata.
- Developer docs landing and detail pages covering Product Overview, Development, APIs, SDKs, Firmware, CLI, Deployment, and Release Notes.
- Feature overview and detail pages for Provision, OTA, Fleet Management, Smart Home Experience, User Management, App SDK, Insights, Private Cloud, and Integrations, including production-grade OTA rollout detail, a structured end-user smart-home workflow story, App SDK reference sample application coverage, a structured private deployment comparison story, and ecosystem integration coverage across Matter Fabric, voice assistants, REST APIs, MQTT over TLS, and webhooks.
- Public authorization wording is intentionally bounded: human product roles and permission assignments belong to the account-side product authorization contract, service bearer scopes remain integration credentials, and this marketing site does not claim to implement or expose the product ACL system.
- File-backed manual routes for customer-facing guides, including `/manual/sdk-samples` for Android, iOS, WebApp, Linux simulator, and PRO2 device sample evaluation paths.
- Locale-aware public site support for English, Traditional Chinese, and Simplified Chinese. English keeps the existing unprefixed URL structure; Traditional Chinese uses `/zh-tw/...`; Simplified Chinese uses `/zh-cn/...`.
- Language switcher in the shared header that points to the same public page in each supported locale.
- Localized public page metadata with canonical URLs, `hreflang` alternates, and localized sitemap entries.
- `robots.txt` and `sitemap.xml` routes for crawl and link discovery.
- Shared footer sitemap for human navigation across platform entry points, features, developer docs, manual chapters, contact, and privacy.
- Localized privacy notice routes describing contact form data, first-party SQLite analytics, analytics event types, referrer-origin-only handling, ephemeral session ids, retention intent, data request handling, admin protection, no third-party analytics or advertising pixels or fingerprinting, and local video behavior.
- Contact / early access registration form.
- SQLite lead capture through `DATABASE_PATH`, defaulting to `data/connectplus.db`.
- Protected admin lead review with filtering, pagination, and filtered CSV export when `ADMIN_TOKEN` is set.
- Plain-text `/healthz` endpoint.
- Multi-stage Docker packaging that ships the server binary plus runtime templates/static assets and defaults SQLite storage to `/data/connectplus.db` inside the container.

Current routes:

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
- `GET /admin/leads`
- `GET /admin/leads.csv`
- `GET /static/...`
- `GET /zh-tw`
- `GET /zh-tw/docs`
- `GET /zh-tw/docs/{slug}`
- `GET /zh-tw/manual`
- `GET /zh-tw/manual/{slug}`
- `GET /zh-tw/features`
- `GET /zh-tw/features/{slug}`
- `GET /zh-tw/contact`
- `GET /zh-tw/privacy`
- `GET /zh-tw/search`
- `POST /zh-tw/contact`
- `GET /zh-cn`
- `GET /zh-cn/docs`
- `GET /zh-cn/docs/{slug}`
- `GET /zh-cn/manual`
- `GET /zh-cn/manual/{slug}`
- `GET /zh-cn/features`
- `GET /zh-cn/features/{slug}`
- `GET /zh-cn/contact`
- `GET /zh-cn/privacy`
- `GET /zh-cn/search`
- `POST /zh-cn/contact`

This implementation is enough to demonstrate the Realtek Connect+ direction and collect leads. It is not yet content-complete as a full public IoT cloud platform website and documentation surface.

## CI/CD Test Report Workflow

Realtek Connect+ follows the shared test report policy from
`hkt999rtk/rtk_cloud_contracts_doc` PR #24.

- Canonical tracked reports:
  - `docs/TEST_REPORT.md` for CI and PR validation.
  - `docs/READINESS_TEST_REPORT.md` for CD and deployed readiness evidence.
- CI generates `.artifacts/report-candidates/docs/TEST_REPORT.md` covering:
  - brand film asset check
  - Go formatting
  - `go test ./...`
  - `go build ./cmd/server`
  - `go run ./cmd/visual-smoke`
- CD generates `.artifacts/report-candidates/docs/READINESS_TEST_REPORT.md` covering:
  - deployment bundle build
  - website-test deployment result
  - public `/healthz`
  - public homepage verification
  - deployed video asset verification
  - stale-copy checks
- CI/CD upload report candidates as the `report-candidates` artifact and keep raw logs in workflow logs or separate artifacts.
- Drift checks compare committed reports with generated candidates for the selected validation profile.
- Manual report import is handled by the `Import Report Candidate` workflow. It can update only `docs/TEST_REPORT.md` or `docs/READINESS_TEST_REPORT.md` in a target PR branch or explicit branch after allowlist, heading, redaction, path, and whitespace validation.
- Report candidates are deterministic. They intentionally avoid volatile timestamps, local absolute paths, and raw credentials so committed reports can be reviewed and drift-checked reliably.

## Visual Direction

The site uses a modern, minimal, direct enterprise style aligned with Realtek's official web presence: white content areas, clear navigation, product-led messaging, and blue-green brand accents. Compared with a traditional corporate site, Realtek Connect+ uses more whitespace, denser feature grids, stronger calls to action, and a clear platform architecture visual.

Footer navigation uses a human-readable sitemap pattern. It is Fixed Content in the shared layout, populated from localized feature/docs/manual catalogs, and intentionally excludes admin, health, static, internal API, robots, and XML sitemap routes. The XML `/sitemap.xml` remains crawler-facing infrastructure rather than a normal visitor page.

Color system:

- Base: white and near-white surfaces.
- Text: Realtek-style charcoal and deep blue-gray.
- Accent: official-site Realtek blue `#0068b7`, deep navigation blue `#035390`, and cyan highlight `#6dcedd` for buttons, links, icons, and highlights.
- Secondary: pale blue-gray panels for architecture, feature detail, and contact surfaces.

UI system:

- Header: Realtek official logo asset from `static/assets/realtek-logo.png`, Connect+ wordmark, Docs, Features, Architecture, Contact.
- Buttons: solid teal primary and restrained outline secondary.
- Cards: 6-8px radius, thin border, minimal shadow.
- Hero: direct product positioning with a platform architecture poster image.
- Iconography: repo-native inline SVG icons only; no npm, external icon package, React, or client-side build step. Icons use a 24px viewbox, rounded caps/joins, approximately 2px stroke, and navy/teal color treatment inside pale teal 8px containers where emphasis is needed.
- Icon coverage: hero chips, primary CTAs, feature cards, architecture nodes, architecture module chips, protocol rails, use cases, docs cards, contact/admin controls, and text links must use semantic icons rather than letter placeholders.
- Feature hierarchy: homepage cards should separate primary product surfaces such as Provision, OTA, Fleet Management, and Private Cloud from secondary surfaces with stronger card scale, larger icons, and subtle platform-diagram styling.
- Architecture visual language: the homepage architecture area should show device, cloud, and app/dashboard nodes with smaller semantic module chips for provisioning, identity, telemetry, registry, OTA, APIs, insights, and support workflows. Protocol labels such as MQTT over TLS, signed firmware image, REST APIs, and webhooks should appear as scannable rails where relevant.
- Page rhythm: use stronger full-width section contrast inspired by modern IoT platform pages, including pale blue enterprise bands, design-principle panels, product surface imagery, component/service cards with visual previews, and public-vs-private deployment comparison cards.
- Typography: keep the Inter/system stack for performance and operational familiarity, but use a clearer type scale. H1 is strong but restrained so product visuals stay visible; section H2 uses tighter hierarchy; card titles and body copy must be distinct enough for scanning. Eyebrow text should be used sparingly with modest letter spacing so pages do not read like a specification dump.

Assets:

- Generated hero/platform and feature imagery lives in `static/assets/`.
- The generated image style must be clean B2B technology: white background, teal/blue accents, credible product UI mockups and technical diagrams, no stock-photo people, no third-party marks, no neon sci-fi/glass-cloud treatment, and no readable fake text.
- Bitmap imagery supports platform context only; semantic navigation and feature recognition should come from the inline SVG icon system.
- Current generated homepage assets:
  - `static/assets/connectplus-hero-corporate-v2.jpg`: text-free semiconductor-to-gateway-to-dashboard hero visual.
  - `static/assets/connectplus-platform-surfaces-corporate-v2.jpg`: enterprise product surface visual showing registry, dashboards, charts, and mobile companion context.
  - `static/assets/connectplus-sample-ecosystem-corporate-v2.jpg`: generic app/client/simulator/reference-device sample ecosystem diagram without third-party platform marks.
- Homepage brand film:
  - The homepage may include the official Realtek corporate brand film as a trust-building section after Architecture and before Deployment.
  - The current implementation uses a local MP4 asset at `static/assets/realtek-brand-film.mp4`, tracked with Git LFS.
  - The section uses a poster image at `static/assets/realtek-brand-film-poster.jpg`, native controls, and `preload="metadata"`.
  - The film must not replace the product hero, must not autoplay, and must keep Realtek Connect+ product CTAs as the primary conversion path.
  - The video must remain responsive at 16:9 and localized with accessible title text.
  - The first public website version does not need a third-party media iframe for the brand film.
- Video is optional. If a ChatGPT video generation tool is available later, the site can add a short product loop with a poster image fallback. If no video tool is available, CSS motion or static generated imagery is sufficient.

## Multilingual Architecture

Supported locales:

- `en`: default locale, unprefixed URLs such as `/`, `/features`, `/docs/apis`, and `/contact`.
- `zh-TW`: Traditional Chinese public URLs under `/zh-tw`, such as `/zh-tw/features/ota`.
- `zh-CN`: Simplified Chinese public URLs under `/zh-cn`, such as `/zh-cn/features/ota`.

Routing rules:

- Locale is resolved from the path prefix only.
- English routes remain unprefixed for backwards compatibility.
- No cookie, session, or automatic `Accept-Language` redirect is used in v1.
- Unknown locale prefixes such as `/fr/...` return 404.
- `/admin/*`, `/healthz`, `/robots.txt`, `/sitemap.xml`, and `/static/...` are global routes and are not localized.

Content rules:

- Public copy is provided through the internal content catalog.
- Feature and documentation slugs remain English and stable across all locales.
- Localized feature and docs catalogs must preserve the same slug set and ordering as English.
- Static image assets are shared across locales; `alt` text is localized through the catalog.
- Contact form service options display localized titles but submit canonical feature slugs to SQLite, avoiding mixed-language lead interest values.
- Admin lead review remains English-only in v1.

## Lightweight YAML/Markdown Content Source (Issue #75)

### Scope

- For this issue, introduce a lightweight content source only for the `/docs` page as a placeholder.
- Do **not** migrate existing content from Go literals or current `docs`/`features` in this step.
- Keep current `/docs/{slug}` and feature rendering paths untouched for now.
- `/docs` is the first demo target for the new content system (layout and rendering only).

### File Format

- Use locale-separated files with YAML frontmatter + Markdown body.
- File path proposal:
  - `content/docs/en/docs.yaml`
  - `content/docs/zh-TW/docs.yaml`
  - `content/docs/zh-CN/docs.yaml`
- Keep frontmatter schema aligned and minimal:
  - `title`, `subtitle`
  - `seo.meta_title`, `seo.meta_description`, `seo.social_image`
  - `hero_image` (optional), with optional `hero_image_alt`
  - `sections[]` for structured block rendering
- Markdown body is required for narrative content and can embed images:
  - `![alt text](/static/assets/...)`

### i18n and fallback

- English is source-of-truth for this phase (`content/docs/en/docs.yaml`).
- Missing locale file falls back to English.
- Locale file parity (fields, slug order, section keys) must be maintained when content exists.
- i18n routing remains path-prefix based.

### Reload strategy

- Content is loaded into memory at process start.
- Add manual cache refresh endpoint:
  - `POST /admin/reload-content`
- Requires valid `ADMIN_TOKEN`.
- On success, clear and rebuild cached docs content from `content/docs`.
- This supports iterative changes without restart during development/review while preserving production behavior.

### Placement

- Place content under repository root `content/docs/...` for non-technical editing.
- Keep templates in `templates/` and media assets in `static/`.

### Acceptance (Issue #75)

- `GET /docs` reads from `content/docs/<locale>/docs.yaml`.
- Missing locale file falls back to English correctly.
- Markdown body is rendered and image embeds are passed through to HTML safely.
- `POST /admin/reload-content` is token-gated and refreshes cache.
- Add tests for parser, render path, fallback, and reload auth behavior.

## Privacy / GDPR-Lite Handling

The website applies privacy information globally instead of using EU-only IP detection.

- Public routes include `/privacy`, `/zh-tw/privacy`, and `/zh-cn/privacy`.
- The footer links to the localized privacy notice on every page.
- The contact form includes a localized privacy note linking to the privacy notice.
- The privacy notice explains contact form fields, inquiry handling purpose, first-party SQLite analytics when `ANALYTICS_ENABLED=true`, OpenAI-backed documentation query behavior when `SEARCH_ENABLED=true`, collected analytics event types, referrer-origin-only handling, ephemeral session id behavior, 90-day raw analytics event retention through `ANALYTICS_RETENTION_DAYS`, 24-month lead retention intent, data access/correction/deletion request handling, admin token protection, no third-party analytics or advertising pixels or fingerprinting, and local video behavior.
- The first implementation uses `privacy@example.com` as a placeholder privacy contact. This must be replaced with an official contact address before public launch.
- The privacy notice is GDPR-lite readiness for the website prototype, not a complete legal compliance package.
- The homepage brand film is served as a local MP4 asset and does not create a YouTube iframe.

## OpenAI Documentation Query

The website includes a first-version public documentation query tool scoped only to content in this repository.

- Public routes include `/search`, `/zh-tw/search`, and `/zh-cn/search`.
- The query API is `POST /api/search` with JSON body `{ "query": "string", "locale": "en|zh-TW|zh-CN" }`.
- The implemented website HTTP API contract is maintained in `docs/openapi.yaml`
  and summarized for humans in `docs/API_REFERENCE.md`.
- The search index is CLI-only and is not rebuilt during server startup. Run `OPENAI_API_KEY=... go run ./cmd/search-index` before enabling runtime search.
- Release packaging can build a precomputed `data/search.db` into the deployable artifact when `OPENAI_API_KEY` is available. This keeps fixed website-document embeddings out of git while allowing deployment bundles to be search-ready.
- The indexer scans website content from `content/`, `docs/`, `README.md`, feature/docs/manual catalog data, and localized public content, then stores documents, chunks, embeddings, and source metadata in SQLite.
- The default search database path is `data/search.db`, overrideable with `SEARCH_DATABASE_PATH`.
- Embeddings use OpenAI `text-embedding-3-small` by default through `SEARCH_EMBEDDING_MODEL`.
- Runtime answers use OpenAI Responses API with `gpt-4.1-mini` by default through `SEARCH_ANSWER_MODEL`.
- `SEARCH_ENABLED` defaults to `false`. When disabled, `/search` still renders but explains that the index is unavailable, and `/api/search` returns a controlled disabled JSON response.
- Runtime retrieval creates a query embedding, searches local SQLite chunks, and filters by score. If no retrieved chunk clears the threshold, the API returns `answer_found=false`, an empty `sources` array, and a localized no-hit answer without calling the answer model.
- When sources are found, only the user query and retrieved snippets are sent to the answer model. The prompt instructs the model to answer only from sources and to say no matching documentation was found if the answer is not present.
- The endpoint limits query length, rate-limits requests by client address, and limits returned sources.
- Raw search queries must not be stored in analytics event payloads. Only aggregate search count, success, or no-hit metrics may be tracked in later analytics work.

SEO rules:

- Public pages emit localized `<html lang>`, title, description, canonical URL, Open Graph, Twitter card metadata, and `hreflang` alternates.
- `hreflang` values are `en`, `zh-Hant`, `zh-Hans`, and `x-default`.
- `/sitemap.xml` includes all public locale variants when indexing is enabled.
- `DISABLE_SEARCH_INDEXING=true` still disables sitemap exposure and emits noindex signals for all localized public pages.

Translation maintenance:

- English remains the source content baseline.
- Traditional Chinese is maintained as authored localized content.
- Simplified Chinese is currently generated from the Traditional Chinese catalog using a local character conversion map and should be reviewed before public launch.
- Any new public page, feature, docs section, form label, validation message, or metadata string must be added to all supported locales in the same change.

## Feature Scope

Realtek Connect+ presents the following Realtek IoT cloud platform capabilities:

- Platform Overview: device firmware, cloud backend, mobile app, dashboard, and third-party integrations.
- Provision: Wi-Fi/BLE onboarding, device binding, activation, user-device association.
- OTA: firmware lifecycle foundation for upload/catalog, target enablement, rollout status, report, cancel, and download; implemented campaign policy surfaces for scheduled, time-window, user-consent, and archive are available now, while advanced rollout operations remain roadmap scope.
- Fleet Management: device registry, groups, metadata/tags, batch operations, timezone, device sharing.
- Smart Home Experience: remote control, local control fallback, schedules, scenes, grouping, node sharing, push notifications, alerts, and household sharing flows framed as product capabilities rather than current website app behavior.
- User Management: sign up, sign in, OTP verification, third-party login, password recovery/change, and account deletion framed as platform capabilities rather than current website authentication flows.
- App SDK: iOS/Android SDK layers, Android/iOS/WebApp reference app samples, Linux simulator, PRO2 camera device demo, push notifications, app publishing path, and launch-readiness ownership boundaries.
- Insights: activation statistics, firmware distribution, logs, crash reports, reboot reasons, RSSI/memory metrics.
- Private Cloud: VM/container cloud-agnostic deployment (GCP/Azure/AWS/on-premises) versus serverless-locked alternatives, public evaluation tier (5-device default, up to 200 on request, non-commercial), private commercial tier (one-time license + annual maintenance, no minimum scale, contact sales for pricing), SDK licensing posture (currently proprietary, open-source release planned at GA), and tier-aware support boundary (community-tier for evaluation, contract-defined for commercial). Tier numbers, pricing structure, and SDK licensing source-of-truth live in `rtk_cloud_workspace/docs/business-model.md`; the website mirrors that document and must not introduce conflicting figures.
- Integrations: Matter Fabric positioning, Alexa/Google Assistant paths, REST APIs, MQTT over TLS, webhooks, and cloud-to-cloud integration boundaries.

## Client Sample Ecosystem

The website summarizes the `rtk_cloud_client` sample ecosystem as a customer enablement path, not as code hosted or executed by this Go website.

- Public surfaces: homepage sample proof point, `/features/app-sdk`, `/docs/sdks`, and `/manual/sdk-samples`.
- App-side samples: Android Home Automation sample, iOS Home Automation sample, and WebApp Ops Lab sample.
- Device-side samples: Linux simulator and PRO2 camera device demo.
- Validation scope: provisioning adapter states, token/session setup, device list/detail, light and AC command flow, camera monitor, snapshot upload, WebRTC Video over TURN signaling and ICE/TURN boundary, MQTT payload inspection, status/log/event reporting, and redacted debug reports.
- Boundaries: samples are SDK usage references, not production app-store apps, white-label release packages, customer release artifacts, or formal cloud wire contracts. The WebApp sample does not implement BLE or SoftAP onboarding. Public streaming copy must describe WebRTC Video over TURN only; TURN and coturn wording is allowed only as WebRTC ICE infrastructure, not as a standalone relay product.
- Source of truth: deeper sample details remain in `rtk_cloud_client/docs/SAMPLE_APPLICATIONS.md`, `rtk_cloud_client/docs/SAMPLE_HOME_APP_SPEC.md`, `rtk_cloud_client/docs/SAMPLE_DEVICE_APP_SPEC.md`, and the sample README files under `rtk_cloud_client/samples/...` plus `rtk_cloud_client/packages/freertos/pro2_demo/README.md`.

## Platform Completion Gap

Status values:

- `Implemented`: working website/backend behavior exists.
- `Content Partial`: page or section exists, but parity depth is incomplete.
- `Planned`: required for a complete Realtek Connect+ public website, not yet represented deeply enough.
- `Out of Scope for website v1`: acknowledged capability, but not intended as a working implementation in the public website.

The matrix below tracks website v1 representation, not live cloud-service implementation. `Implemented` means the public site content or first-party backend behavior exists; it does not mean the marketing site ships a working device cloud, mobile app, or operator console for that domain.

| Capability | Website v1 Status | Current Surface | Website v1 Notes |
| --- | --- | --- | --- |
| Platform Overview | Content Partial | Homepage, `/features`, docs landing page | The site now explains the device-cloud-app-dashboard story, but follow-on work can deepen security, scalability, cost, and deployment comparison narratives. |
| Provision | Content Partial | `/features/provision`, [`PRODUCT_ONBOARDING.md`](https://github.com/hkt999rtk/rtk_cloud_contracts_doc/blob/main/PRODUCT_ONBOARDING.md) | Provisioning copy now distinguishes the contract-backed cloud registry, cross-service activation, service-scoped provisioning credential, and transport-readiness foundation from integration-ready claim material interfaces and roadmap local onboarding work. Local Wi-Fi/BLE setup, QR/SoftAP UX, ownership transfer, factory reset policy, and aggregate product readiness are not described as generally available until the owner repositories land those implementations. |
| OTA | Content Partial | `/features/ota`, [`FIRMWARE_CAMPAIGN.md`](https://github.com/hkt999rtk/rtk_cloud_contracts_doc/blob/main/FIRMWARE_CAMPAIGN.md) | OTA copy now distinguishes the available firmware lifecycle foundation for upload, catalog, target enablement, rollout status, report, cancel, and download from the implemented campaign policy surfaces for scheduled, time-window, user-consent, and archive. Approval workflow, dashboards, analytics, and staged percentage rollout remain roadmap scope rather than generally available phase-one implementation. |
| Fleet Management | Implemented | `/features/fleet-management` | The fleet page covers node registration, bootstrap certificates, registry, groups, metadata, OTA jobs, firmware images, batch operations, and operator statistics widgets as website content without claiming a live console implementation. |
| Admin Operations | Content Partial | `/features/fleet-management`, `/admin/leads`, `/admin/leads.csv` | Website v1 ships lead-review tooling plus admin-operations product copy, but it does not ship the full fleet console described by the marketing content. |
| User Management | Content Partial | `/features/user-management` | The feature page covers sign up, sign in, OTP verification, third-party login, forgot/change password, delete account, and account lifecycle boundaries; follow-on work can add deeper session behavior, support workflows, and visuals. |
| End-user Smart Home Features | Content Partial | `/features/smart-home` | The smart-home page now covers remote control, local fallback, scheduling, scenes, grouping, node sharing, push notifications, alerts, and household workflow boundaries, with room for richer mobile personas and control-state diagrams. |
| Mobile App SDK | Implemented | `/features/app-sdk`, `/manual/sdk-samples` | The app SDK page and manual now cover iOS and Android SDK layers, Android/iOS/WebApp reference samples, Linux simulator, PRO2 camera device demo, push notifications, and App Store/Google Play publishing guidance without introducing a client-side framework into the website itself. The website summarizes `rtk_cloud_client`; deeper sample specs remain in `rtk_cloud_client/docs/SAMPLE_APPLICATIONS.md` and the sample README files. |
| Insights | Content Partial | `/features/insights` | Insights copy covers activation statistics, firmware distribution, logs, crash reports, reboot reasons, RSSI, and memory signals, but the website still needs stronger dashboard visuals and deeper support/metrics storytelling. |
| Private Cloud / Deployment | Implemented | `/features/private-cloud`, `/docs/deployment` | Website v1 now compares public evaluation, managed private deployment, and customer-operated private regions with coverage for VM/container deployment substrate, GCP/Azure/AWS/on-premises targets, evaluation device limits (5 default / 200 max), commercial pricing structure (license + maintenance, no minimum scale, contact sales for figures), SDK licensing posture, and tier-aware support boundaries. The page mirrors `rtk_cloud_workspace/docs/business-model.md`; self-service signup is owned by `rtk_cloud_admin` and `rtk_account_manager` and is referred to from the marketing site rather than implemented inside it. |
| Matter / Ecosystem Integrations | Website Content Implemented | `/features/integrations` | The integrations page covers Matter Fabric positioning, voice assistants, MQTT over TLS, REST APIs, webhooks, and a structured integration-path comparison. This is public website content only: Matter, Alexa, and Google Assistant remain roadmap service capabilities per `SMART_HOME_ECOSYSTEM.md`; push alerts are integration-ready; scenes, schedules, and household sharing are not generally available platform features. |
| Developer Docs Portal | Content Partial | `/docs`, `/docs/product-overview`, `/docs/development`, `/docs/deployment`, `/docs/release-notes` | The portal structure exists and is navigable, but follow-on work can deepen setup guides, architecture diagrams, and operational runbooks. |
| APIs | Content Partial | `/docs/apis`, `/features/integrations` | API positioning exists and now separates product authorization roles from service bearer scopes. Website v1 still lacks reference-grade endpoint coverage, finalized auth flows, webhook payload examples, and error-model detail. |
| SDK Reference | Content Partial | `/docs/sdks`, `/features/app-sdk`, `/manual/sdk-samples` | The docs, feature, and manual surfaces now include a sample matrix covering Android Home Automation, iOS Home Automation, WebApp Ops Lab, Linux simulator, and PRO2 camera device demo. Follow-on depth can add versioned install guides and language-specific code walkthroughs. |
| CLI | Content Partial | `/docs/cli` | The CLI section exists as part of the docs portal, but website v1 still needs command catalogs, auth/session examples, and operator workflow walkthroughs. |
| SEO / Launch Readiness | Content Partial | Shared layout metadata, `/robots.txt`, `/sitemap.xml`, `go run ./cmd/visual-smoke` | Metadata, sitemap, robots, CI, deployment packaging, and desktop/mobile visual smoke checks for English, Traditional Chinese, and Simplified Chinese public pages now exist; remaining work is broader launch polish such as expanded product visuals plus final parity and documentation close-out. |
| Real IoT Cloud Operations | Out of Scope for website v1 | Public marketing and docs copy only | The public website will describe platform capabilities, but it will not implement real device provisioning, OTA delivery, user auth, telemetry ingestion, or a production device-operations control plane in v1. |

## Website Completion Roadmap

- **v0.1 Marketing Foundation**: current state. Working Go site, product positioning, core feature pages, generated hero asset, contact lead capture, admin lead export.
- **v0.2 Platform Content Depth**: expand remaining content depth so Realtek Connect+ can credibly present its user, admin, mobile SDK, deployment, Matter, and developer documentation areas.
- **v0.3 Launch Readiness**: improve SEO, accessibility, visual assets, deployment packaging, CI, form hardening, admin usability, and server operations.
- **v1.0 Public Website Candidate**: polished public-facing site with complete parity content, Realtek brand assets, reliable contact/admin workflows, and deployment documentation.

## HTTP Interface

Routes:

- `GET /`: homepage.
- `GET /docs`: developer/documentation portal landing page.
- `GET /docs/{slug}`: developer/documentation section detail pages.
- `GET /manual`: file-backed customer manual landing page.
- `GET /manual/{slug}`: file-backed customer manual detail pages, including `/manual/sdk-samples`.
- `GET /features`: feature overview.
- `GET /features/{slug}`: feature detail pages.
- `GET /contact`: contact / early access registration form.
- `GET /privacy`: privacy notice.
- `GET /search`: OpenAI-backed documentation query UI.
- `GET /robots.txt`: crawl directives for search bots.
- `GET /sitemap.xml`: sitemap covering public marketing and docs pages.
- `POST /contact`: validate and store a lead in SQLite.
- `POST /api/search`: JSON documentation query endpoint.
- `GET /healthz`: plain-text health check.
- `GET /admin/leads`: protected lead review page, enabled only when `ADMIN_TOKEN` is set.
- `GET /admin/leads.csv`: protected CSV export, enabled only when `ADMIN_TOKEN` is set.
- `GET /static/...`: CSS and asset files.

The OpenAPI contract in `docs/openapi.yaml` covers API/form/admin endpoints
implemented by this website runtime. It intentionally does not describe the
future IoT product control-plane APIs.

Localized public route variants:

- Traditional Chinese mirrors public routes under `/zh-tw`.
- Simplified Chinese mirrors public routes under `/zh-cn`.
- Examples: `/zh-tw/features/ota`, `/zh-cn/docs/apis`, `/zh-tw/manual/sdk-samples`, `/zh-tw/contact`, `/zh-cn/search`.
- Privacy examples: `/zh-tw/privacy`, `/zh-cn/privacy`.
- Localized `POST /contact` variants write to the same SQLite lead table.
- Feature and documentation slugs remain English across all locales.

Feature slugs:

- `provision`: aligns public provisioning availability wording with the product onboarding interface contract, separating the cloud-side activation foundation from integration-ready claim material and roadmap local onboarding/readiness work.
- `ota`: aligns public OTA availability wording with the firmware campaign interface contract, promoting the implemented campaign policy surfaces while keeping roadmap operations explicit.
- `fleet-management`
- `smart-home`
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
- `ANALYTICS_ENABLED`, optional and enabled by default. When false-like, analytics database setup is skipped until the later event-ingestion work lands.
- `ANALYTICS_DATABASE_PATH`, optional and default `data/analytics.db`.
- `ANALYTICS_RETENTION_DAYS`, optional and default `90`.
- `ADMIN_TOKEN`, optional. When set, enables protected lead review and CSV export.
- `DISABLE_SEARCH_INDEXING`, optional. When truthy, marks the site as non-indexable with HTTP `X-Robots-Tag`, page-level robots meta tags, `/robots.txt` `Disallow: /`, and a disabled `/sitemap.xml`.
- `OPENAI_API_KEY`, required by `cmd/search-index` and by runtime search answering when `SEARCH_ENABLED=true`.
- `SEARCH_DATABASE_PATH`, optional and default `data/search.db`.
- `SEARCH_ENABLED`, optional and disabled by default.
- `SEARCH_EMBEDDING_MODEL`, optional and default `text-embedding-3-small`.
- `SEARCH_ANSWER_MODEL`, optional and default `gpt-4.1-mini`.
- `PUBLIC_BASE_URL`, optional. When empty, canonical URLs, social image URLs, `hreflang` alternates, robots sitemap references, and sitemap locations are generated from the incoming request host and forwarded headers. When set, generated public absolute URLs use this fixed base URL.
- `ENABLE_ASSET_FINGERPRINTS`, optional and disabled by default. When truthy, template-rendered `/static/...` URLs receive a `?v=<content-hash>` query string based on file contents.
- `ENABLE_CDN_CACHE_HEADERS`, optional and disabled by default. When truthy, the app emits provider-neutral CDN cache headers.

Operational behavior:

- The runtime uses configured `http.Server` read, write, and idle timeouts.
- Requests are logged with method, path, response status, and elapsed time.
- `SIGINT` and `SIGTERM` trigger graceful shutdown with a bounded timeout.
- Production TLS is expected to terminate at a reverse proxy, ingress controller, or hosting platform in front of the app.
- Preview and test deployments should set `DISABLE_SEARCH_INDEXING=true` until the site is approved for public search indexing.

CDN-ready behavior:

- The runtime is provider-neutral and does not assume Cloudflare, CloudFront, Fastly, or any other CDN.
- With default configuration, CDN-specific behavior is off and current deployment behavior is preserved.
- `PUBLIC_BASE_URL` should be set only when the public CDN/reverse-proxy origin is known, such as `https://webtest.mgmeet.io`.
- `ENABLE_ASSET_FINGERPRINTS=true` appends content hashes to existing static asset URLs without changing the underlying `/static/...` route.
- `ENABLE_CDN_CACHE_HEADERS=true` applies:
  - `/static/*`: `Cache-Control: public, max-age=31536000, immutable` for existing static files.
  - Public GET HTML pages: `Cache-Control: no-store`.
  - `POST /contact`, localized contact POST variants, `/admin/*`, and `/healthz`: `Cache-Control: no-store`.
  - `/robots.txt` and `/sitemap.xml`: `Cache-Control: public, max-age=300`.
- The CD workflow intentionally does not enable the CDN env vars yet. CDN provider selection, DNS cutover, cache rules, purge strategy, compression, and TLS settings should be completed as a deployment decision.

Persistence and cache boundaries:

- Contact leads, first-party analytics events, and the precomputed documentation search index are website-local SQLite repositories.
- Redis-compatible caching is not a first-priority dependency for this repository. Future application cache work must be justified by a measured website bottleneck or concrete operational requirement.
- Website SQLite stores must not become authoritative storage for real IoT telemetry, product device state, customer account state, fleet inventory or control data, OTA campaign execution state, or production mobile app user state.
- CDN/static cache headers, asset fingerprints, reverse proxy rules, and CDN provider configuration are deployment concerns separate from application Redis/cache work.

Commands:

- `go run ./cmd/server`
- `go test ./...`
- `go build -o bin/realtek-connect ./cmd/server`
- `docker build -t realtek-connect .`
- `docker run --rm -p 8080:8080 -v "$(pwd)/data:/data" realtek-connect`

## SQLite

SQLite stores website leads only. It does not store real IoT telemetry or device state.

When analytics is enabled, first-party event telemetry uses a separate SQLite database so raw analytics data stays isolated from lead data.

Redis-compatible caching is low priority for this repository. Leads, analytics,
and documentation search should remain concrete SQLite repositories unless a
specific runtime bottleneck appears. Any future cache work should document the
need first and must not turn the public website into a product device-state or
telemetry source of truth. The cross-repository roadmap is maintained in
`../../docs/persistence-cache-refactor-roadmap.md`.

Default database path:

```text
data/connectplus.db
```

Container default database path:

```text
/data/connectplus.db
```

Default analytics database path:

```text
data/analytics.db
```

Container default analytics database path:

```text
/data/analytics.db
```

Container deployment notes:

- `Dockerfile` copies the compiled server together with `templates/` and `static/`, which are runtime dependencies for page rendering and asset delivery.
- The native CD bundle includes a writable `data/` directory so default `data/connectplus.db` and `data/analytics.db` paths can initialize on the website test host. When `OPENAI_API_KEY` is available during packaging, the bundle also includes precomputed `data/search.db` for documentation search. Production native hosts should set `DATABASE_PATH` and `ANALYTICS_DATABASE_PATH` to persistent service-owned storage while allowing `SEARCH_DATABASE_PATH` to point at the immutable bundled index if desired.
- `/data` is declared as the persistent volume for SQLite-backed lead and analytics storage.
- HTTPS is intentionally out of process and should be handled by the reverse proxy or deployment platform instead of the Go app directly.

Linode artifact deployment notes:

- `deploy/package.sh <version>` builds `dist/realtek-connect-<version>.tar.gz`, `realtek-connect-<version>.tar.gz.sha256`, and `realtek-connect-<version>.object-manifest.json`.
- Release artifacts contain the server binary, `content/`, `templates/`, `static/`, `deploy/`, `VERSION`, and release manifest metadata. SQLite runtime DB files are excluded; the only SQLite DB allowed in a release bundle is optional precomputed `data/search.db`.
- CI validates source, tests, server build, and visual smoke evidence, but does not upload deployable release bundles to Linode Object Storage.
- `.github/workflows/release.yml` publishes versioned artifacts to GitHub Releases and Linode Object Storage under `releases/realtek_connect-<version>/`.
- `.github/workflows/deploy-linode.yml` installs a selected version to `/opt/realtek-connect/releases/<version>`, updates `/opt/realtek-connect/current`, restarts `realtek-connect.service`, and runs public readiness checks.
- Linode host bootstrap, GoDaddy DNS, nginx reverse proxy, Let’s Encrypt TLS, rollback, and SQLite backup are documented in `docs/deployment-linode.md`, `docs/deployment-promotion-rollback.md`, and `docs/sqlite-backup-linode.md`.

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

Analytics schema:

```sql
CREATE TABLE IF NOT EXISTS analytics_events (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  ts INTEGER NOT NULL,
  event TEXT NOT NULL,
  page TEXT NOT NULL,
  cta TEXT,
  percent INTEGER,
  duration INTEGER,
  variant TEXT,
  referrer_origin TEXT,
  session_id TEXT,
  created_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_analytics_events_ts
  ON analytics_events(ts);

CREATE INDEX IF NOT EXISTS idx_analytics_events_event_page
  ON analytics_events(event, page);
```

Contact form fields:

- `name`, required.
- `company`, optional.
- `email`, required.
- `interest`, required.
- `message`, optional.

## Test Plan

- `go test ./...`
- HTTP route tests for `/`, `/docs`, `/features`, feature/detail pages, `/contact`, `/privacy`, `/search`, `/robots.txt`, and `/sitemap.xml`.
- Localized HTTP route tests for `/zh-tw`, `/zh-tw/features/{slug}`, `/zh-tw/docs/{slug}`, `/zh-tw/contact`, `/zh-tw/privacy`, `/zh-tw/search`, `/zh-cn`, `/zh-cn/features/{slug}`, `/zh-cn/docs/{slug}`, `/zh-cn/contact`, `/zh-cn/privacy`, and `/zh-cn/search`.
- Localized metadata tests for `<html lang>`, canonical URL, `hreflang` alternates, and language switcher current-state links.
- Unknown feature slug returns 404.
- Unknown locale prefix and unknown localized feature/docs slugs return 404.
- Valid contact POST writes SQLite and shows success.
- Invalid contact POST shows validation errors and does not write a lead.
- Localized contact POSTs show localized success/error messaging and store canonical service slugs in SQLite.
- CDN readiness tests verify default behavior remains unfingerprinted and cache-neutral when CDN env vars are unset.
- CDN readiness tests verify `PUBLIC_BASE_URL` affects canonical URLs, social image URLs, `hreflang`, robots sitemap references, and sitemap locations.
- CDN readiness tests verify `ENABLE_ASSET_FINGERPRINTS=true` adds content hashes to rendered static asset URLs.
- CDN readiness tests verify `ENABLE_CDN_CACHE_HEADERS=true` applies the expected static, public HTML, contact, admin, health, robots, and sitemap cache headers.
- Privacy tests verify localized privacy notice routes, footer links, contact form notice links, sitemap privacy URLs, and privacy metadata.
- Search tests verify content extraction, SQLite search schema, no-hit behavior without answer generation, source-only prompt construction, query validation, rate limiting, disabled behavior, localized `/search` routes, cited source responses, and no-hit JSON responses.
- Homepage brand film tests verify the local MP4 source, poster image, native video metadata preload, no YouTube iframe, and localized section copy between Architecture and Deployment.
- Admin lead routes require `ADMIN_TOKEN`; unauthorized requests return 401, disabled admin routes return 404.
- `go run ./cmd/visual-smoke`
- The visual smoke command checks English, Traditional Chinese, and Simplified Chinese public pages at desktop/mobile widths, verifies representative hero/feature images load, and fails on horizontal overflow without adding npm dependencies.

## Definition of Done

For any future issue that changes website behavior or public content:

- `go test ./...` passes.
- `go build ./cmd/server` passes.
- Go source is formatted with `gofmt`.
- Desktop and mobile visual smoke checks confirm no blank page, missing representative image asset, localized route regression, or horizontal overflow.
- No npm, React, Tailwind, or frontend build step is introduced.
- Any new public route is documented in both `docs/SPEC.md` and `README.md`.
- Any new public page or public text string is added to every supported locale.
- Any new feature or docs item keeps stable English slugs across locales and is covered by catalog parity tests.
- Runtime SQLite files are not committed.
- Generated website assets are stored under `static/assets/` if referenced by templates or CSS.

## Developer Issue Backlog

The following GitHub issues are the source-of-truth backlog for moving from v0.1 to v1.0.

### 1. Expand Realtek Connect+ feature matrix

Labels: `documentation`, `enhancement`

Goal: Complete the Realtek Connect+ public capability matrix.

Acceptance criteria:

- Matrix covers user management, smart home features, admin operations, mobile SDK, deployment, Matter, Insights, OTA, docs, APIs, SDKs, and CLI.
- Each capability has a status and notes for website v1.
- Matrix avoids claiming real cloud operations are implemented.

### 2. Add developer and documentation portal structure

Labels: `documentation`, `enhancement`

Goal: Add a developer/docs section with clear product, engineering, API, SDK, and deployment entry points.

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
- Copy uses Realtek Connect+ naming and reinforces Realtek product ownership.

### 5. Add admin and device operations content

Labels: `enhancement`

Goal: Present admin/device operations beyond the website lead admin page.

Acceptance criteria:

- Content covers node registration, certificates, device registry, OTA jobs, firmware images, batch operations, and statistics widgets.
- Adds or expands an admin/platform operations page or section.
- Clearly separates website lead admin from a future IoT platform admin console.

### 6. Expand OTA content to production-grade detail

Labels: `enhancement`

Goal: Make the OTA story stronger and more credible.

Acceptance criteria:

- Covers firmware upload, metadata extraction, version/model targeting, rollout status, report, cancel, and download as the available firmware lifecycle foundation.
- Includes a visual, table, or structured explanation that labels scheduled, time-window, user-consent, archive, approval workflow, analytics, dashboards, and staged percentage rollout according to implementation status.
- `/features/ota` remains the canonical OTA page.

### 7. Add mobile SDK and app publishing content

Labels: `documentation`, `enhancement`

Goal: Present a complete Realtek Connect+ mobile app development story.

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
- Document the production deployment profile, including backup and restore notes.
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

Implementation notes:

- `go run ./cmd/visual-smoke` starts the existing Go-rendered server in-process by default and drives local Chrome headlessly for desktop/mobile checks across English, Traditional Chinese, and Simplified Chinese public pages.
- The command can also target an already running server with `-base-url`.

## First-Version Limits

- HTTP only; TLS can be handled later by reverse proxy or deployment platform.
- Product pages describe platform capabilities; they do not perform real IoT cloud operations.
- Realtek official logo is expected to be provided by the user. Until then the site uses a text wordmark.
