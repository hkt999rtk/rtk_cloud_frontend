---
title: "能力工作流"
description: "使用 SDK 集成 provisioning、device operation、telemetry、firmware、command 与 WebRTC signaling。"
---

## Provisioning、device 与 OTA

Account Manager 负责 product claim、account binding 与 readiness；Video Cloud 负责 video-device activation 与 runtime state。使用 SDK helper 并以 bounded deadline 查询 readiness，不可由 product app 调用 service-internal activation route。

呈现控制项前先读取 device 与 product capability。Command 的 HTTP acceptance 与 asynchronous device result 应分开处理并使用 correlation ID。Telemetry/log 必须 typed 且 redacted。Firmware SDK 提供 campaign、rollout、report 与 cancellation vocabulary，但 application 仍负责 rollout policy、approval、signature verification 与 restart recovery。

## Video Cloud 责任

| Layer | 提供内容 | 不负责 |
| --- | --- | --- |
| Cloud | HTTPS signaling、current-owner MQTT/WebSocket delivery、TURN credential、session state | Live media frame 与 renderer |
| RTK SDK | Authentication、ICE、offer/answer/close、error 与 clip workflow | 完整 peer connection/media engine |
| Product app | Platform WebRTC、renderer、audio policy 与 UX lifecycle | Device camera/codec |
| Device SDK/firmware | Offer、answer、camera/audio track 与 resource limit | Product app UI |

RTK SDK 会交换 offer、answer、ICE 与 session lifecycle。Product app 不需要自行创建 signaling backend，但仍需将 RTK SDK session data 接到 platform WebRTC/media component。Stored-video playback 是独立 HTTP workflow。
