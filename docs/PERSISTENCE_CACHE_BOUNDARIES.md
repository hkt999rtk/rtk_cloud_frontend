# Persistence And Cache Boundaries

This repository is a Go-rendered public website and documentation surface. Its
runtime persistence is intentionally small and website-local. Redis-compatible
caching is not a first-priority dependency for this repository unless measured
runtime behavior shows a concrete bottleneck or an operational requirement that
SQLite and existing HTTP/CDN cache controls cannot satisfy.

## Website-Owned SQLite Stores

The website owns these local SQLite stores:

- `connectplus.db`: contact lead capture for public forms and protected admin
  lead review.
- `analytics.db`: first-party website analytics events and aggregate admin
  reporting.
- `search.db`: precomputed documentation search documents, chunks, embeddings,
  and source metadata.

Keep these as concrete SQLite repositories by default. They are small,
deployment-local data stores that match the current website workload and backup
model. Future cache work must start from a measured bottleneck, such as observed
query latency, lock contention, startup cost, operational restore pain, or a
specific multi-instance deployment requirement.

For the Kubernetes v1 deployment profile, `connectplus.db` and `analytics.db`
must live on the `/data` PVC and the frontend Deployment must stay at
`replicas: 1`. Do not horizontally scale the frontend while these writable
stores remain SQLite-backed. A rolling update must avoid two active writers
using the same SQLite files at the same time.

## Out-Of-Scope State

The website SQLite stores must not become authoritative storage for IoT platform
state. In particular, do not store these classes of data in this repository's
website databases:

- real device telemetry
- product device state
- customer account state
- fleet inventory or authoritative fleet control data
- OTA campaign execution state
- production mobile app user state

Those domains belong to the platform services and contracts outside this
marketing/documentation website.

## Cache Boundary

Redis or Redis-compatible cache infrastructure should not be added only because
the broader platform may use it. For this repository, Redis is low priority
until there is a documented operational need.

When proposing Redis or another application cache, document:

- the measured website bottleneck or concrete operational requirement
- why the current SQLite repository, precomputed search index, or static asset
  strategy is insufficient
- the intended invalidation and backup behavior
- how the change preserves public routes, payloads, admin lead behavior,
  analytics behavior, and documentation search behavior

CDN and static asset cache headers are a separate deployment concern. Keep
`ENABLE_CDN_CACHE_HEADERS`, asset fingerprints, reverse proxy rules, and CDN
provider configuration separate from application Redis/cache work.
