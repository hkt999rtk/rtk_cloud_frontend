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

SDK helpers exchange offers, answers, ICE data, and session lifecycle messages. The application owns the WebRTC peer connection, media engine, tracks, rendering, audio policy, TURN behavior, and user-visible call state. Stored-video playback is a separate HTTP workflow.
