---
title: "iOS Swift SDK"
description: "Integrate the Swift package, asynchronous operations, clip browser, protected keys, and AVPlayer."
---

## Install

Add the `RTKCloudClient` Swift package at the approved release tag or use the internal source archive. Supported package platforms are iOS 13 or later and macOS 12 or later. Import `RTKCloudClient` in the target that owns the network-facing service layer.

## Client usage

Create `RTKCloudClient` with the HTTPS base URL and required device context. Use throwing synchronous APIs only away from the main thread and prefer async overloads from Swift concurrency code. Catch `RTKCloudError`, preserve its stable status for diagnostics, and show user-safe messages.

## Stored video browser

Create `ClipQuery` with bearer token, optional filters, limit, and skip. `listClips` returns a typed `ClipPage`; display its `ClipSummary` values and follow `nextSkip`. Load thumbnails away from the main actor. Do not trust a cached clip URL as a playback session.

## Encrypted playback

Implement `PlaybackKeyProvider` in code backed by Keychain and, where appropriate, Secure Enclave key operations. The provider receives the selected clip and current server key and returns wrapped clip key plus ephemeral public key. Give the `PlaybackSession.playbackURL` to a new `AVPlayerItem`. Request another session after expiry and replace the item when the selection changes.

## PKI

Use platform certificate and identity storage. Keep private keys non-exportable where product requirements allow. Validate server trust normally; never add a production trust callback that accepts arbitrary certificates.

## Validation

Run `swift test` and `swift build` for the package and sample. Run simulator build, installation, launch, and UI tests on a macOS host with an iOS Simulator SDK. Deployed-cloud validation supplies credentials to the simulator sandbox at runtime and must not compile them into the app.
