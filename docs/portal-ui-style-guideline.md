# Realtek Connect+ Portal — Web UI style guideline

Status: implemented public-Portal baseline, 2026-09-05. This is a product UI guideline, not an official Realtek corporate identity manual.

## Direction and scope

The Portal introduces Realtek Connect+ to enterprise buyers and developers. It should be easy to evaluate, easy to navigate, and comfortable to read. Use the visual discipline of developer documentation products, while keeping Realtek's logo and blue brand identity. Do not copy another company's branding.

Scope: homepage, feature index and detail pages, documentation index and details, manual/SDK pages, contact, and shared legal-page presentation. English, Traditional Chinese, and Simplified Chinese use one design system. Private Portal administration screens are excluded from the new stylesheet. The service login and customer console remain separate applications.

## Design tokens

| Role | Value | Usage |
| --- | --- | --- |
| Brand blue | `#0068b7` | Primary actions, focus rings, selected navigation |
| Deep blue | `#035390` | Links, hero emphasis, managed-cloud band |
| Cyan | `#6dcedd` | Decorative accent only; not small text on white |
| Heading | `#162b40` | Titles and strong labels |
| Body | `#242f3e` | Primary reading content |
| Secondary | `#526172` | Supporting text, metadata |
| Canvas | `#ffffff` | Main surface |
| Neutral surface | `#f6f8fa` | Footer, code blocks, section separation |
| Blue surface | `#edf5fb` | Selected/featured contextual panels |
| Border | `#dce3ea` | Cards and dividers; not form-control boundaries |
| Input border | `#7b8998` | Visible form boundaries |
| Error | `#b42318` | Invalid field boundaries, accompanied by error text |

Tokens are defined on `.portal-ui` in `static/portal-ui.css`. Keep the original logo asset intact. Brand reference: [Realtek corporate website](https://www.realtek.com/). The neutral colors and component rules above are Connect+ product decisions, not claimed corporate standards.

## Typography and layout

- Use the native system sans-serif stack with Traditional Chinese fallbacks. Do not require an external font download.
- Body text: 16px, line-height 1.65; manual prose: line-height 1.8 and roughly 78 characters per line maximum.
- Homepage headline: 36–48px, with a smaller product descriptor. Interior h1: 32–44px; section h2: 28–36px; card titles: 22–24px.
- Use weight 600–650 for hierarchy. Reserve uppercase letter spacing for short English category labels, never long body paragraphs.
- Content maximum width: 1200px. Desktop gutters: 32px; mobile gutters: 20px. Typical spacing: 8, 12, 16, 20, 24, 32, 40, 48, 64px.
- Use 6px control corners and 8px card corners. Prefer a fine border to a shadow. A subtle shadow is acceptable on a product demonstration frame.
- Avoid forced no-wrap headings, fixed-height text panels, perspective transformations, oversized empty cards, and screenshots cropped by arbitrary heights.

## Shared components

### Navigation

Keep brand, documentation, manual, features, architecture, contact, service entry, and locale selection discoverable. Mark the current section using `aria-current="page"`, blue text, and an underline or selected background. Detail pages use breadcrumbs. Never change the configured service login URL to a hardcoded environment.

Mobile navigation is a disclosure, not a modal dialog: the trigger exposes expanded state; opening moves focus into the links; Escape closes it and restores trigger focus. Outside-click and selecting a link close it. Without JavaScript, navigation links remain available.

### Buttons and links

Use one primary task per action group. Filled Realtek blue is the primary action; a bordered white button is secondary. Controls should be at least 44px high. Use descriptive labels rather than repeated “Learn more” where context is unclear. Documentation's main action leads to the manual, not a sales form. SDK/legal consent links must remain clearly distinguishable.

### Cards, diagrams, and images

Use consistent padding, aligned headings, one category label, concise summary, and a clear destination. Feature imagery uses a consistent preview frame without cropping technical content. Keep demo imagery identifiable as demonstration material. Do not imply real customer metrics, production status, certifications, or security guarantees through decorative screenshots.

Architecture uses three stages (device, cloud, application), laid out horizontally on desktop and vertically on mobile. Core capability tabs remain keyboard navigable and become accordions on mobile.

### Documentation and SDK

Use real headings, lists, links, tables, and code elements. Manual pages with two or more h2 headings receive an automatically generated “On this page” navigation. This enhancement uses DOM text, preserves existing heading IDs, and does not rewrite technical content. At narrow widths, the contents list moves above the article. Without JavaScript, the article remains fully available.

Technical diagrams retain intrinsic proportions. Long code and tables scroll inside their container, not across the whole page. Keep package versions, checksums, validation status, terms versions, and explicit consent controls visible. Do not preselect consent or weaken downloads' server-side checks.

### Forms, errors, and legal pages

Labels stay visible above controls. Keep autocomplete, required state, max lengths, submitted values, honeypot handling, and localized server errors. Invalid fields use `aria-invalid` and linked error descriptions; the summary remains an alert with links to affected fields. Do not use color alone to convey an error or success.

Legal text must remain intact. Styling does not constitute legal approval. Do not invent privacy contacts, rewrite license obligations, or remove evaluation warnings as a cosmetic change.

## Localization and accessibility

- All new navigation labels must exist in English and Traditional Chinese; the existing catalog derives Simplified Chinese. File-owned content requires all three locale files.
- Keep localized destinations, canonical links, alternate-language metadata, and route slugs unchanged.
- Verify 1440px, 768px, 390px, and 320px layouts, keyboard controls, long Chinese headings, and focus visibility.
- Target WCAG AA text contrast. Cyan is not a text color on white. Validate actual component combinations, including hover and invalid states, before introducing new colors.
- Keep the skip link, semantic headings, accessible names, reduced-motion support, and progressive enhancement. Automated checks are not a complete accessibility certification.

## Implementation and verification

`templates/layout.html` loads the versioned `portal-ui.css` layer only for public pages. Existing `styles.css` remains the base for compatibility and private administration. Add Portal overrides to this one documented layer rather than spreading new styles across templates.

Run from the frontend repository:

```sh
GOWORK=off go test ./...
GOWORK=off go run ./cmd/visual-smoke -portal-review -timeout 8m -screenshot-dir /tmp/portal-review
```

The visual command uses the existing Chrome/chromedp dependency and starts a local server unless `-base-url` is specified. Use a loopback preview for interaction checks. The contact check submits an empty form and expects validation errors, not a stored lead. Search is excluded from the default visual matrix because the feature is disabled by default; enabled-search behavior has separate existing server tests.

Review screenshots as well as test results. Check header, first viewport, lower sections, footer, cards, forms, and the longest SDK content. Preserve analytics hooks and do not reinterpret historical CTA events when changing a destination; the new manual CTA uses `docs_cta_manual`.

## Content-quality gate before public release

Visual consistency alone does not make documentation production-ready. Existing generic manual chapters contain placeholder imagery and sample instructions. Some legal/contact content still requires owner approval. Replace placeholders with reviewed technical guidance, approve actual privacy contact details and terms, and audit claims against available functionality before publishing. Do not silently hide these gaps with polished styling.
