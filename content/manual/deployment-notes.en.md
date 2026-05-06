---
title: "Deployment Notes"
description: "Track the production deployment profile, backup and restore steps, and manual publishing checks."
---

## Before you publish

- Confirm the page has valid YAML front matter.
- Review the rendered output on `/manual/deployment-notes`.
- Make sure the manual index still points to the intended chapter list.
- Confirm the production profile still matches the current app packaging and data layout.

## Production profile

The site is expected to run as a Go binary inside a container or VM.

- Keep templates and static assets inside the shipped image or build artifact.
- Mount a persistent SQLite volume for lead records and analytics events.
- Terminate TLS at the reverse proxy, ingress controller, or hosting platform.
- Point health checks at `/healthz` so platform probes stay lightweight.
- Preserve cache headers at the edge instead of adding TLS or cache concerns to the app.

## Runtime refresh

Use the protected admin content reload endpoint after updating files on disk.

1. Update the Markdown content.
1. Save the manual index if the section list changed.
1. Reload the content cache with a valid admin token.

## Backup and restore

- Snapshot the SQLite volume before a production rollout or content migration.
- Restore by attaching a clean volume, copying the database snapshot into place, and restarting the app.
- Verify the restored deployment with `/healthz`, `/admin/leads`, and `/admin/analytics`.
- Re-run the public smoke path after restore so the homepage, docs, and manual routes all render.

## Rollback

- Keep the previous image or deployment manifest available until the new release is verified.
- Roll back by restoring the prior image and the last known-good SQLite snapshot together.
- Confirm the new deployment state still matches the content source files before re-enabling writes.

## Rollout checklist

| Check | Result |
| --- | --- |
| English source exists | Ready |
| Locale fallback works | Ready |
| Navigation entry resolves | Ready |
