---
title: "Troubleshooting and compatibility"
description: "Diagnose authentication, transport, upload, browsing, playback, and version problems."
---

## Request failures

| Symptom | Likely cause | Action |
| --- | --- | --- |
| TLS or hostname failure | Wrong base URL, CA, system time, or proxy | Verify endpoint, certificate chain, hostname, and clock |
| 401 | Missing or expired credential | Refresh the token or repair device certificate selection |
| 403 | Missing permission or product capability | Check identity, ownership, and Account Manager capability data |
| 404 | Wrong device/clip or deleted resource | Refresh device and clip state |
| 410 from clip upload | Legacy upload route after direct-upload cutover | Use authorize, presigned PUT, complete, and status APIs |
| Timeout | Network path, server load, or an operation deadline that is too short | Capture correlation IDs and retry only when safe |

## Upload failures

Presigned PUT failures commonly result from a changed signed header, wrong content length, bearer token leakage to object storage, expired URL, early EOF, or a ciphertext hash mismatch. Do not complete an upload unless the PUT returned success. Poll the upload record to distinguish verification failure from expiration.

## Playback failures

Create a new session when a URL expires. For encrypted clips, verify that the provider used the active playback public key and returned both wrapped fields. Confirm that the platform player supports the returned content type and HTTP range requests. Never replace playback-session authorization with the legacy query token.

## Version mismatch

Use the documentation manifest to compare SDK commit and version. Regenerate references when exported symbols change. A server contract change must update the normative contract and SDK compatibility tests before user documentation claims support.

## Support evidence

Provide SDK version, server version, platform/toolchain version, operation name, stable status, HTTP status, sanitized correlation IDs, and reproduction steps. Do not attach tokens, private keys, presigned URLs, wrapped keys, or customer media.
