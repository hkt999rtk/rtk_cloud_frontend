# Portal design implementation — 2026-09-05

## Upstream integration follow-up

This Portal refresh has been preserved on latest main `302cb31`, including the
upstream documentation changes from PR #147. The isolated local integration branch
is `codex/ui-quality-integration-20260905`; the original checkout remains intact.
The integrated source passed `GOWORK=off go test ./...` and the Portal visual review
(144 page/viewport cases plus navigation, contents and form-validation interactions).
The current merged preview uses port 18088 with disposable local data and analytics
disabled. No GitHub push or deployment is part of this integration.

The workspace report `docs/ui-quality-integration-20260905.md` contains the current
cross-repository evidence and release boundaries. Earlier results below are retained
as the initial implementation record.

## Delivered

Public Portal styling now follows a shared Realtek blue, white, and neutral system. The implementation retains Go templates and the existing JavaScript; no new frontend dependency, external font, tracking service, or generated image was introduced.

| Page family | Implemented improvement |
| --- | --- |
| Shared shell | Proportional Realtek branding, explicit Manual entry, current-section indicators, quieter footer and controls |
| Homepage | Responsive headline hierarchy, balanced hero, complete-aspect-ratio product imagery, simplified deployment and architecture panels |
| Features | Consistent card spacing, legible descriptions, contained previews, breadcrumbs and calmer detail panels |
| Documentation | Manual-first primary action, three-column reading tracks, customer-facing copy replacing internal publishing details in all locales |
| Manuals / SDK | Reading-width article, responsive diagrams/code/tables, generated localized contents navigation, consistent SDK cards |
| Contact | Clear context/form layout, stronger input boundaries, retained server-side validation and accessible errors |
| Privacy / evaluation terms | Shared typography and surfaces; legal wording and acceptance semantics preserved |

The style source is `static/portal-ui.css`; component rules and publishing criteria are documented in [Portal UI style guideline](portal-ui-style-guideline.md). Public templates load this layer through the existing asset-versioning helper. Private Portal administration is excluded from the stylesheet. Service login URLs and localized routes are unchanged.

The documentation manual-entry CTA uses the new allowed analytics event `docs_cta_manual`; the historical `docs_cta_primary` event remains accepted for backwards compatibility.

## Verification

- `GOWORK=off go test ./...` passed across the repository, including the new Portal and analytics regression tests.
- The extended Chrome smoke review passed **144 page/viewport cases**, covering English, Traditional Chinese, Simplified Chinese, and 1440 / 768 / 390 / 320px widths. Each case checks a rendered heading, loaded target image, document title, and absence of document-level horizontal overflow.
- Browser interaction checks passed: mobile menu opening and focus, Escape dismissal, desktop arrow-key tabs, mobile accordion, localized manual contents and anchor focus, and empty contact-form validation.
- Visual inspection included desktop/mobile homepage, tablet homepage, desktop documentation, desktop manual, desktop contact, mobile feature cards, and narrow SDK documentation screenshots.
- Added server regression coverage for public styling, localized manual links, current navigation, manual enhancement markup, private-admin style isolation, manual CTA analytics, and production asset fingerprinting.
- All browser actions used a local preview with an isolated temporary database. No customer record or deployed website was modified. The empty-form check is expected to fail validation without creating a lead.

Reproduce from the frontend repository:

```sh
GOWORK=off go test ./...
GOWORK=off go run ./cmd/visual-smoke -portal-review -timeout 8m -screenshot-dir /tmp/portal-review
```

Use a loopback `-base-url` for interactive review. Remote interaction checks refuse to submit forms. The browser command's default local server uses the repository content and does not need deployment credentials.

## Review boundaries and release notes

- This is a **local implementation**, not a deployment or a full accessibility certification.
- The responsive matrix covers page families, not every individual feature/manual article. Existing server tests provide broader route and business-logic coverage.
- Search is disabled in the selected preview and is omitted from the visual matrix. Its backend/disabled behavior remains covered by existing tests; the shared search form styling should be visually checked in an enabled-search environment before release.
- The SDK download catalog is disabled in this preview. Existing automated tests retain coverage of consent/version/download handling; an enabled-catalog visual review is still needed before release.
- Some generic manual chapters still contain placeholder imagery and sample text. These need technical-owner review; this redesign deliberately does not invent setup instructions or treat placeholders as validated documentation.
- Privacy contacts and evaluation/legal text need the appropriate owner's approval. Existing legal language, including English evaluation terms on localized routes, was preserved.
- No source commit, pull request, staging rollout, or production publication was performed.
