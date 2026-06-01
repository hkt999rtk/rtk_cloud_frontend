# Website API Reference

This document summarizes the HTTP API implemented by `rtk_cloud_frontend`.
The machine-readable OpenAPI contract lives in
[`docs/openapi.yaml`](openapi.yaml).

The website API is separate from future Realtek Connect+ product control-plane
APIs. It only covers the marketing/documentation site runtime: analytics,
documentation search, contact lead capture, health, and protected admin
operations.

## Authentication

Public endpoints do not require authentication:

- `POST /api/event`
- `POST /api/search`
- `GET /contact`, `POST /contact`
- `GET /zh-tw/contact`, `POST /zh-tw/contact`
- `GET /zh-cn/contact`, `POST /zh-cn/contact`
- `GET /healthz`

Admin endpoints require `ADMIN_TOKEN` when configured. Prefer the
`X-Admin-Token` header:

```bash
curl -H "X-Admin-Token: $ADMIN_TOKEN" http://localhost:8080/admin/leads
```

The `token` query parameter remains available for manual browser access:

```bash
curl "http://localhost:8080/admin/leads.csv?token=$ADMIN_TOKEN"
```

When `ADMIN_TOKEN` is unset, admin endpoints return `404`.

## Analytics Event API

`POST /api/event`

Stores a privacy-minimal analytics event when analytics storage is enabled.
The endpoint accepts only `application/json` and a maximum 4 KB request body.
Unknown JSON fields are ignored.

```json
{
  "event": "click_cta",
  "page": "home",
  "cta": "contact_us",
  "session_id": "2f3f4b7c-ephemeral"
}
```

Allowed event types are `page_view`, `click_cta`, `scroll`, and `engaged`.
The backend validates event-specific fields:

- `click_cta` requires an allowlisted `cta`.
- `scroll` requires `percent` set to `25`, `50`, `75`, or `100`.
- `engaged` requires `duration` set to `10`, `30`, or `60`.

Successful storage returns:

```json
{
  "status": "ok"
}
```

If analytics storage is disabled, the endpoint returns `204 No Content` and
does not store anything.

## Documentation Search API

`POST /api/search`

Queries the local documentation search index and returns a source-grounded
answer when relevant website content is found.

```json
{
  "query": "How do I enable documentation search?",
  "locale": "en"
}
```

`query` is required and limited to 500 characters. `locale` may be `en`,
`zh-TW`, or `zh-CN`; omitted or unknown values default to `en`.

Successful responses use this shape:

```json
{
  "answer_found": true,
  "answer": "Build the index with cmd/search-index, then start the server with SEARCH_ENABLED=true.",
  "sources": [
    {
      "title": "README",
      "url": "",
      "snippet": "Build the local website-only search index with OPENAI_API_KEY...",
      "locale": "en",
      "source_type": "file",
      "score": 0.82
    }
  ]
}
```

When no retrieved source clears the relevance threshold, the response is still
`200 OK` with `answer_found=false` and an empty `sources` array.

Error responses include the same base fields plus `error`:

```json
{
  "answer_found": false,
  "answer": "Search is not enabled.",
  "sources": [],
  "error": "search_disabled"
}
```

Documented error values are `search_disabled`, `rate_limited`,
`invalid_json`, `invalid_query`, and `search_unavailable`.

## Contact Lead Form

`POST /contact`

Localized variants:

- `POST /zh-tw/contact`
- `POST /zh-cn/contact`

The contact form uses `application/x-www-form-urlencoded` fields:

- `name`, required, max 120 characters.
- `company`, optional, max 160 characters.
- `email`, required valid email.
- `interest`, required. Allowed values: `evaluation-access`,
  `commercial-deployment`, `partnership`, `technical-question`, `other`.
- `message`, optional, max 2000 characters.
- `website`, spam-trap field; leave empty.

Successful and validation responses render HTML. Rate limiting returns
`429 Too Many Requests`.

## Health

`GET /healthz`

Returns `200 OK` and a plain-text body:

```text
ok
```

The response is marked `Cache-Control: no-store`.

## Admin Operations

`GET /admin/leads`

Renders the protected lead review page. Query filters:

- `email`: partial email match.
- `company`: partial company match.
- `interest`: exact contact interest slug.
- `page`: one-based page number.
- `token`: query fallback for `ADMIN_TOKEN`.

`GET /admin/leads.csv`

Exports filtered leads as CSV with these columns:

```text
id,name,company,email,interest,message,created_at
```

`POST /admin/reload-content`

Reloads docs and manual content from disk. On success, the response is:

```text
content reloaded
```

## Privacy And Logging Notes

The website API must not log or store raw lead payloads, cookies, OpenAI/API
keys, SMTP secrets, credentialed connection details, or full referrer URLs.
Search query text may be sent to OpenAI for embeddings and source-grounded
answers when `SEARCH_ENABLED=true`, but raw query text is not stored in the
analytics event payload.
