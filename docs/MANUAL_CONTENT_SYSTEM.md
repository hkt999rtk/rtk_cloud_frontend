# File-Based Manual / Documentation Content System

## Overview

This document specifies a lightweight, file-based content system for serving user
manuals and documentation pages on Realtek Connect+. Content lives in the repository
as plain files (YAML + Markdown), is loaded at server startup, and is served under the
`/manual` route. No CMS, no admin UI, and no content database are required.

## Guiding Principles

- Content is Git-native: all updates go through PR review.
- No runtime dependency beyond the Go binary and the `content/` directory.
- Structured metadata is YAML; prose is Markdown (with YAML front matter).
- Images inline in Markdown are served via a dedicated static sub-route.
- Multi-language: one file per locale per page.
- Reload without restart: a protected admin endpoint re-reads content from disk.

## Routes

| Route | Description |
|---|---|
| `/manual` | Index: lists all manual sections/chapters. |
| `/manual/{slug}` | Detail: renders a single page by slug. |
| `/content-assets/{path}` | Serves images stored inside `content/` (restricted to `content/assets/`). |

The `/docs` route continues to exist as the developer-portal stub. `/manual` is the
new file-based surface for user-facing manuals and documentation. In a future release
`/manual` may absorb or replace `/docs` after the content migration is validated.

## File Layout

```
content/                          ← repo root, alongside templates/ and static/
├── manual/
│   ├── index.en.yaml             ← section/chapter index (English)
│   ├── index.zh-TW.yaml          ← section/chapter index (Traditional Chinese)
│   ├── index.zh-CN.yaml          ← section/chapter index (Simplified Chinese)
│   ├── getting-started.en.md     ← prose page (English)
│   ├── getting-started.zh-TW.md  ← prose page (Traditional Chinese)
│   ├── getting-started.zh-CN.md  ← prose page (Simplified Chinese)
│   └── …
└── assets/
    ├── manual/                   ← images referenced from manual Markdown
    │   ├── setup-diagram.png
    │   └── …
    └── …
```

Convention: `{slug}.{locale}.md` or `{slug}.{locale}.yaml`.
Supported locale codes: `en`, `zh-TW`, `zh-CN`.

## Index File Format (YAML)

`content/manual/index.{locale}.yaml` defines the ordered list of sections and their
pages. This is what the `/manual` listing page renders.

```yaml
# content/manual/index.en.yaml
title: "User Manual"
description: "Step-by-step guides and reference documentation for Realtek Connect+."
sections:
  - slug: getting-started
    title: "Getting Started"
    summary: "Set up your first device and verify connectivity."
  - slug: ota-updates
    title: "OTA Updates"
    summary: "Push firmware updates to your device fleet."
  - slug: fleet-management
    title: "Fleet Management"
    summary: "Organize devices into groups and monitor status at scale."
```

Fields:

| Field | Type | Required | Description |
|---|---|---|---|
| `title` | string | yes | Page `<h1>` for the index. |
| `description` | string | yes | Meta description and subtitle. |
| `sections[].slug` | string | yes | Matches `{slug}.{locale}.md` filename and URL path. |
| `sections[].title` | string | yes | Section heading on the index page. |
| `sections[].summary` | string | yes | One-sentence description shown in the listing. |

## Prose Page Format (Markdown + YAML Front Matter)

`content/manual/{slug}.{locale}.md` is a standard Markdown file with a YAML front
matter block delimited by `---`.

```markdown
---
title: "Getting Started"
description: "Set up your first device and verify connectivity in under 10 minutes."
---

## Prerequisites

- A Realtek Connect+ account (see [Sign Up](/contact)).
- A supported Realtek Wi-Fi module.

## Step 1: Flash the SDK

Download the SDK and flash the default firmware:

![SDK download screen](../../content-assets/manual/setup-diagram.png)
```

Front matter fields:

| Field | Type | Required | Description |
|---|---|---|---|
| `title` | string | yes | `<title>` tag and `<h1>` on the detail page. |
| `description` | string | yes | Meta description. |

Markdown body: any CommonMark / GFM syntax supported by
[goldmark](https://github.com/yuin/goldmark). Inline images use paths relative to
`/content-assets/` as shown above.

## Inline Images

Images embedded in Markdown are referenced with a URL starting with
`/content-assets/`:

```markdown
![alt text](/content-assets/manual/my-diagram.png)
```

The server maps `/content-assets/` to the `content/assets/` directory at the repo
root. Only `content/assets/` is exposed this way — the rest of `content/` (YAML,
Markdown source files) is never served directly.

## Markdown Rendering

Library: **goldmark** (`github.com/yuin/goldmark`).
Extensions enabled: GFM tables, strikethrough, autolinks, task lists.
HTML output is sanitized before template insertion (allow-list of safe tags).

Goldmark is a pure-Go library with no C dependencies and passes the CommonMark
spec.

## Multi-Language (i18n)

One file per locale per page. Locale resolution follows the same logic as the rest of
the site:

1. Request arrives at `/manual/{slug}` (no prefix → `en`).
2. Request arrives at `/zh-tw/manual/{slug}` → `zh-TW`.
3. Request arrives at `/zh-cn/manual/{slug}` → `zh-CN`.

If a locale variant does not exist, the server falls back to `en`. If `en` also does
not exist, the server returns 404.

Index file resolution follows the same pattern: `/manual` loads
`content/manual/index.{locale}.yaml`.

## Content Loading and Reload

**Startup:** the server reads all files under `content/` into memory during
`Server.init()`. Parsing errors are logged but do not crash the server; the affected
page returns 404.

**Hot reload:** a `POST /admin/reload-content` endpoint re-reads all content files
from disk without restarting the process. The endpoint is protected by the same
admin-token contract used by the leads admin: `X-Admin-Token` header or
`?token=` query parameter backed by `ADMIN_TOKEN`.

This strategy means:
- Zero latency for page renders (no disk I/O per request).
- Content updates in production require either a process restart or a call to the
  reload endpoint.
- The reload endpoint is not exposed publicly.

## Go Package Layout

```
internal/
└── manual/
    ├── loader.go     ← reads content/ files, returns ManualIndex and []ManualPage
    ├── page.go       ← ManualIndex, ManualSection, ManualPage types
    └── render.go     ← goldmark Markdown → template.HTML conversion
```

The `web.Server` struct gains a `manualLoader *manual.Loader` field loaded at startup.
The reload endpoint calls `loader.Reload()` under a mutex.

## Templates

```
templates/
├── manual_index.html   ← /manual listing page
└── manual_page.html    ← /manual/{slug} detail page
```

Both extend `layout.html` via `{{template "layout" .}}`.

## Placeholder Implementation (v0)

The first implementation is a demo that proves the routing, file loading, Markdown
rendering, and template system work end-to-end. It ships with:

- `content/manual/index.en.yaml` — three placeholder sections.
- `content/manual/getting-started.en.md` — one placeholder prose page with a dummy
  inline image.
- `content/assets/manual/placeholder.png` — a static placeholder image.
- `templates/manual_index.html` and `templates/manual_page.html`.
- `internal/manual/loader.go`, `page.go`, `render.go`.
- `/manual` and `/manual/{slug}` routes wired in `internal/web/server.go`.
- `/content-assets/` static file route restricted to `content/assets/`.
- `POST /admin/reload-content` endpoint.

No existing content is extracted or migrated in this phase.

## Dependency

Add to `go.mod`:

```
github.com/yuin/goldmark v1.7.x
```

No other new runtime dependencies.

## Out of Scope for v0

- Migration of existing `/docs` content into the new system.
- Simplified Chinese or Traditional Chinese placeholder content.
- Search / full-text index.
- Table-of-contents auto-generation.
- Versioned docs.
