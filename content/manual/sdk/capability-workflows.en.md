---
title: "Capability workflows"
description: "Use the SDK for provisioning, device operations, telemetry, firmware, commands, and WebRTC signaling."
---

## Provisioning and activation

Account Manager owns product claim, account binding, and readiness. Video Cloud owns video-device activation and runtime state. Use the SDK provisioning helpers where available and poll readiness with a bounded deadline. Do not call service-internal activation routes from a product application.

## Device state and commands

Read device information before presenting controls. Send commands through the typed or documented SDK helper, attach correlation IDs, and handle asynchronous device results separately from HTTP acceptance. Product capability data controls whether the UI exposes commands.

## Telemetry and logs

Use typed telemetry events for supported payload families. Runtime device logs and central cloud-service logs are separate systems. Redact secrets before submission and set retention appropriate to the product.

## Firmware campaigns

Firmware helpers expose campaign, rollout, report, and cancellation vocabulary. The SDK does not choose rollout policy, approve firmware, or replace artifact-signature verification. Applications must persist enough state to recover safely after restart.

## WebRTC signaling

Realtek Connect+ provides the signaling service: HTTPS APIs create and manage sessions, while the current device owner's MQTT or WebSocket transport receives the same `webrtc_offer` payload. SDK helpers exchange offers, answers, ICE data, and session lifecycle messages, so product teams do not need to build a separate signaling backend.

The application owns the platform WebRTC peer connection, media engine, tracks, rendering, audio policy, and user-visible call state. Device SDK or firmware owns offer handling, answer generation, camera/audio tracks, codecs, and resource limits. Stored-video playback is a separate HTTP workflow.

The current release does not bundle server-side transcoding, S3 multipart upload, simulcast negotiation, renegotiation, or a complete in-SDK WebRTC media renderer.
