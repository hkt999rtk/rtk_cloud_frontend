---
title: "Android Kotlin SDK"
description: "Integrate the Android AAR, coroutine APIs, clip browser, secure playback keys, and Media3."
---

## Install

Consume the published Maven coordinates or local release repository produced by the package script. The handoff includes the AAR, sources JAR, POM, Gradle module metadata, manifest, checksum, and packaged-consumer smoke report. Configure Internet permission and the deployment's TLS/network security policy without disabling certificate verification.

## Client usage

Create `RtkCloudClient` with the HTTPS base URL, device identifier where required, and the SDK HTTP/WebSocket adapters. Call blocking APIs only from a worker thread. Prefer the `Async` or coroutine variants from view models and services. Map `RtkCloudException.status` to stable application states and redact `responseBody` before diagnostics.

## Stored video browser

Construct `ClipQuery` with bearer token, optional device and event filters, optional time range, positive limit, and non-negative skip. `listClips` and `listClipsAsync` return `ClipPage`, containing typed `ClipSummary` values, total count, and optional `nextSkip`. Use `downloadThumbnail` for preview bytes and `deleteClip` only after explicit confirmation.

## Encrypted playback

Implement `PlaybackKeyProvider` behind Android Keystore-backed code. The provider receives `ClipSummary` and the server `PlaybackKey`, and returns `PlaybackKeyMaterial`. Pass the selected clip to `createPlaybackSession`; the SDK obtains the server key and submits only wrapped material. Give `playbackUrl` to Media3/ExoPlayer and release the player with the UI lifecycle.

## Authentication and PKI

Keep bearer tokens in memory or encrypted app storage. Android PKI helpers integrate platform key references without exporting the private key. Do not put tokens, playback URLs, wrapped keys, or private key identifiers into saved UI state, analytics, or debug exports.

## Validation

Run package unit tests, Android instrumentation tests where a device/emulator is available, the AAR consumer smoke, sample unit tests, and the sample UI smoke. Live-cloud tests require explicit runtime credentials and must never be enabled in untrusted pull requests.
