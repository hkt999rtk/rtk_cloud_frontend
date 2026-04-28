# Realtek Connect+ Specification

## Summary

Realtek Connect+ is an English-first B2B website for an end-to-end IoT cloud platform. It is positioned for IoT product teams evaluating device onboarding, OTA, fleet operations, mobile app enablement, insights, private deployment, and ecosystem integrations.

The first version is implemented as an HTTP-only Go application using `net/http`, `html/template`, and SQLite. It does not use npm, React, Tailwind, or any frontend build step.

## Visual Direction

The site uses a modern, minimal, direct enterprise style aligned with Realtek's official web presence: white content areas, clear navigation, product-led messaging, and blue-green brand accents. Compared with a traditional corporate site, Realtek Connect+ uses more whitespace, denser feature grids, stronger calls to action, and a clear platform architecture visual.

Color system:

- Base: white and near-white surfaces.
- Text: deep navy and charcoal.
- Accent: Realtek-style blue-green / teal for buttons, links, icons, and highlights.
- Secondary: pale blue-gray panels for architecture, feature detail, and contact surfaces.

UI system:

- Header: Realtek logo slot or text wordmark, Connect+ brand, Features, Architecture, Contact.
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
- OTA: firmware upload, campaign rollout, job status, cancel/archive, version validation, and Immediate/Scheduled/User-controlled rollout modes.
- Fleet Management: device registry, groups, metadata/tags, batch operations, timezone, device sharing.
- App SDK: iOS/Android SDK, sample app, rebrand/customize path, push notifications, app publishing path.
- Insights: activation statistics, firmware distribution, logs, crash reports, reboot reasons, RSSI/memory metrics.
- Private Cloud: enterprise deployment, data ownership, custom domain, cloud customization, commercial support.
- Integrations: Alexa, Google Assistant, Matter, REST API, MQTT over TLS, webhooks.

## HTTP Interface

Routes:

- `GET /`: homepage.
- `GET /features`: feature overview.
- `GET /features/{slug}`: feature detail pages.
- `GET /contact`: contact / early access registration form.
- `POST /contact`: validate and store a lead in SQLite.
- `GET /static/...`: CSS and asset files.

Feature slugs:

- `provision`
- `ota`
- `fleet-management`
- `app-sdk`
- `insights`
- `private-cloud`
- `integrations`

Environment:

- `PORT`, default `8080`.
- `DATABASE_PATH`, default `data/connectplus.db`.

Commands:

- `go run ./cmd/server`
- `go test ./...`
- `go build -o bin/realtek-connect ./cmd/server`

## SQLite

SQLite stores website leads only. It does not store real IoT telemetry or device state.

Default database path:

```text
data/connectplus.db
```

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
- HTTP route tests for `/`, `/features`, all feature detail pages, and `/contact`.
- Unknown feature slug returns 404.
- Valid contact POST writes SQLite and shows success.
- Invalid contact POST shows validation errors and does not write a lead.
- Manual browser checks verify desktop/mobile layout, Realtek-style white + teal palette, generated image loading, and no npm/React/Tailwind dependency.

## First-Version Limits

- HTTP only; TLS can be handled later by reverse proxy or deployment platform.
- Product pages describe platform capabilities; they do not perform real IoT cloud operations.
- Realtek official logo is expected to be provided by the user. Until then the site uses a text wordmark.
