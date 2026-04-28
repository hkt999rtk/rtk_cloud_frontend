# Realtek Connect+

Realtek Connect+ is a Go-rendered HTTP website for a Realtek-style IoT cloud platform concept. It uses `net/http`, `html/template`, SQLite, and static CSS. There is no npm, React, Tailwind, or frontend build step.

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
- `GET /features`
- `GET /features/{slug}`
- `GET /contact`
- `POST /contact`
- `GET /healthz`
- `GET /admin/leads`, requires `ADMIN_TOKEN`
- `GET /admin/leads.csv`, requires `ADMIN_TOKEN`
- `GET /static/...`

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
