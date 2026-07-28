---
title: "Video Cloud workflows"
description: "Integrate Live WebRTC signaling and upload, browse, delete, and play encrypted stored video."
---

![Stored video upload, browser, and playback flow](/content-assets/manual/sdk/video-flow.png)

## Live WebRTC signaling

1. The app calls the HTTPS ICE API through the RTK SDK.
2. The app's platform WebRTC component creates the SDP offer.
3. The RTK SDK creates the signaling session over HTTPS.
4. The cloud sends `webrtc_offer` through the current device owner's MQTT or WebSocket transport.
5. Device SDK or firmware receives the offer, attaches camera and audio tracks, and creates an answer.
6. The device submits the answer through the HTTPS answer API.
7. The app retrieves the answer through the RTK SDK and completes media negotiation.
8. The app or device closes the session.

Validate close, expiry, timeout, offline, busy, and unsupported-capability states. Signaling delivery does not fan out and does not automatically fall back to a non-owner transport. The cloud coordinates signaling but does not receive or store Live media frames. A Live WebRTC session never becomes a Stored Clip automatically.

## Direct device upload

1. Finish recording a complete MP4 and encrypt the stored bytes.
2. Calculate the ciphertext size and base64 SHA-256 required by the service contract.
3. Create the versioned encryption descriptor and wrapped `clipkey` metadata.
4. Authorize with `POST /v1/devices/{device_id}/clip-uploads` through the native SDK.
5. Copy every returned signed header and PUT only encrypted bytes to the HTTPS presigned URL.
6. Complete the upload through the SDK.
7. Poll upload state until `ready`, `failed`, or `expired`.

The native streaming PUT accepts a rewindable read callback, an exact body size, an optional progress callback, an optional cancellation check, and a configurable bounded buffer. The default implementation uses a small buffer suitable for memory-constrained devices. A zero-byte read before the declared body size is a protocol failure.

Current technical defaults are 256 MiB for MP4 clips, 5 MiB for JPEG snapshots, 10 minutes for signed URLs, and 30 minutes for an upload lifecycle. Retention profiles are deployment-configurable at 1, 7, or 30 days. These values are not pricing, backup, region, quota, or SLA commitments.

## Browser pagination

Android `listClips` and iOS `listClips` return typed `ClipPage` values. Start with `skip = 0` and a positive `limit`, display `ClipSummary` values, and request `nextSkip` until it is absent. Filters can include device, event type, and time range. Handle an empty page as a normal state.

## Thumbnails and deletion

Use `downloadThumbnail` rather than composing download paths in application code. Decode image bytes outside the main thread and apply an application cache policy that does not expose credentials. Confirm user intent before `deleteClip`; deletion removes cloud metadata and media according to service policy and should not be retried blindly.

## Playback sessions

1. Select a typed clip summary.
2. For encrypted media, ask `PlaybackKeyProvider` to prepare wrapped material using the active server playback key.
3. Call `createPlaybackSession` with the device, clip, bearer token, and wrapped fields.
4. Give the returned short-lived range-capable URL to Media3/ExoPlayer or AVPlayer.
5. Refresh the session when it expires instead of persisting the URL.

The player belongs to the application. Release Android player instances with the screen lifecycle and replace the iOS player item when selecting another clip. Treat 401/403 as authentication or authorization failures, 404 as a stale/deleted clip, and an expired playback URL as a request for a new session.

## Legacy behavior

`uploadClip` and `/upload_clip` remain for source compatibility with pre-cutover deployments. Direct-upload-enabled servers return 410 for clip media. New device integrations must use authorize, presigned PUT, complete, and status operations. Snapshot upload remains a separate supported compatibility path.

## Current release boundary

The current release does not bundle server-side transcoding, S3 multipart upload, simulcast negotiation, renegotiation, or a complete in-SDK WebRTC media renderer.
