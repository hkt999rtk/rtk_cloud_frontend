---
title: "Deployment Notes"
description: "Track rollout assumptions, content refresh steps, and manual publishing checks."
---

## Before you publish

- Confirm the page has valid YAML front matter.
- Review the rendered output on `/manual/deployment-notes`.
- Make sure the manual index still points to the intended chapter list.

## Runtime refresh

Use the protected admin content reload endpoint after updating files on disk.

1. Update the Markdown content.
1. Save the manual index if the section list changed.
1. Reload the content cache with a valid admin token.

## Rollout checklist

| Check | Result |
| --- | --- |
| English source exists | Ready |
| Locale fallback works | Ready |
| Navigation entry resolves | Ready |

