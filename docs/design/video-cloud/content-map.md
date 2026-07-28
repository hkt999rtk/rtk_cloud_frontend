# Video Cloud Product Page Content Map

Status: implemented responsive design source for `/features/video-cloud`.

## Design source

The production Go template, localized feature catalogs, and shared responsive
CSS are the design source. Screenshots in this directory are review artifacts,
not a second implementation.

| Surface | Source |
| --- | --- |
| Page content and data structures | `internal/features/features.go` |
| Traditional/Simplified Chinese | `internal/content/zh.go` |
| Responsive page composition | `templates/feature.html` |
| Flow, table, and related-link styling | `static/styles.css` |
| Text-free hero illustration | `static/assets/connectplus-video-cloud-corporate-v1.png` |

## Page sequence

1. Two-path hero: Live WebRTC and encrypted Stored Video.
2. Three summary cards: available cloud behavior, integration capabilities,
   and product outcome.
3. Eight-step Live signaling flow.
4. Six-step encrypted clip lifecycle.
5. Cloud / RTK SDK / Product App / Device responsibility matrix.
6. Sample application truth matrix.
7. Snapshot versus Clip comparison.
8. Retention, URL, upload-lifecycle, and object-size technical defaults.
9. Current-release exclusions and product-app/device ownership.
10. Related App SDK and SDK Manual entry points.

## Required product truths

- Realtek Connect+ supplies the HTTPS and current-owner MQTT/WebSocket
  signaling service; users do not build a separate signaling backend.
- Android and iOS SDKs expose signaling and stored-video workflow helpers.
- Product apps connect SDK session data to platform WebRTC/media components.
- Device SDK/firmware owns answer generation and physical camera/audio tracks.
- Live media is not received or stored by the cloud and does not automatically
  become a Stored Clip.
- The current release excludes server-side transcoding, S3 multipart upload,
  simulcast negotiation, renegotiation, and a complete in-SDK media renderer.
- 1/7/30-day retention, MP4 256 MiB, JPEG 5 MiB, 10-minute signed URLs, and a
  30-minute upload lifecycle are deployment-configurable technical defaults,
  not price, quota, backup, region, or SLA commitments.

## Responsive acceptance

- Desktop uses four-column flow steps and full-width horizontally scrollable
  matrices.
- Tablet uses two-column flow steps.
- Mobile uses one-column flows and related cards; tables remain within the
  existing overflow container.
- English, Traditional Chinese, and Simplified Chinese use the same section
  order and assets.

## Generated review artifacts

- `desktop.png`: English desktop viewport.
- `mobile.png`: English mobile viewport.
- `zh-tw-desktop.png`: Traditional Chinese desktop viewport and glyph check.
