---
title: "Sample applications"
description: "Build the reference applications and use them as executable SDK integration examples."
---

## Android

Run `gradle -p samples/android assembleDebug` and `gradle -p samples/android testDebugUnitTest`. Camera > Playback accepts a Video Cloud base URL, device ID, and app bearer token at runtime. It lists clips, creates a playback session, and renders with Media3/ExoPlayer. Credentials remain in Compose state and are excluded from debug exports.

Visible status: **Playback — Real SDK + Media3. Live — RTK signaling integration / fixture UI; complete media rendering is not implemented.**

## iOS

Run `swift test --package-path samples/ios` and `swift build --package-path samples/ios`. The simulator helper can build, install, and launch an application bundle when an iOS Simulator SDK is installed. Camera > Playback uses the typed SDK browser and AVPlayer. Credentials and playback URLs are excluded from reports.

Visible status: **Playback — Real SDK + AVPlayer. Live — RTK signaling integration / fixture UI; complete media rendering is not implemented.**

## Web and Linux

The WebApp Ops Lab is a fixture-backed operations and signaling-helper demonstration, not a production WebRTC peer. The Linux simulator demonstrates device transport, command handling, state, logs, and snapshot behavior; it does not simulate camera frames or WebRTC signaling.

## FreeRTOS/Pro2

The Pro2 demo is a source integration reference. Its host smoke tests validate adapter boundaries and signaling lifecycle without claiming physical-board or vendor-media validation. Production firmware must supply board I/O, secure storage, time, networking, media, and vendor WebRTC implementations.

## Using samples safely

Copy architecture and error-handling patterns, not credentials or deployment constants. Keep sample runtime profiles separate from production configuration. Re-run the package and sample tests after upgrading the SDK.
