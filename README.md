# Realtek Connect+

Realtek Connect+ is a Go-rendered HTTP website for a Realtek-style IoT cloud platform concept. It uses `net/http`, `html/template`, SQLite, and static CSS. There is no npm, React, Tailwind, or frontend build step.

## Project Status

Current status: **v0.1 Marketing Foundation**.

This repository currently contains a working marketing website foundation, a developer docs portal structure, per-page SEO/social metadata, sitemap and robots endpoints, contact lead capture, SQLite storage, admin lead review, CSV export, and health check. It is not yet a complete ESP RainMaker parity website, IoT console, user authentication service, real OTA service, device provisioning backend, or telemetry platform.

The full roadmap and developer issue backlog live in [`docs/SPEC.md`](docs/SPEC.md).

## Run

```bash
go run ./cmd/server
```

Open:

```text
http://localhost:8080
```

## Configuration

Environment variables:

- `PORT`: HTTP port, default `8080`.
- `DATABASE_PATH`: SQLite database path, default `data/connectplus.db`.
- `ADMIN_TOKEN`: enables protected lead viewing and CSV export.

## Routes

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
- `GET /admin/leads`, requires `ADMIN_TOKEN`
- `GET /admin/leads.csv`, requires `ADMIN_TOKEN`
- `GET /static/...`

New public routes should be documented here and in `docs/SPEC.md`.

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

## Build

```bash
go build -o bin/realtek-connect ./cmd/server
```
