---
title: "能力工作流"
description: "使用 SDK 整合 provisioning、device operation、telemetry、firmware、command 與 WebRTC signaling。"
---

## Provisioning、device 與 OTA

Account Manager 負責 product claim、account binding 與 readiness；Video Cloud 負責 video-device activation 與 runtime state。使用 SDK helper 並以 bounded deadline 查詢 readiness，不可由 product app 呼叫 service-internal activation route。

呈現控制項前先讀取 device 與 product capability。Command 的 HTTP acceptance 與 asynchronous device result 應分開處理並使用 correlation ID。Telemetry/log 必須 typed 且 redacted。Firmware SDK 提供 campaign、rollout、report 與 cancellation vocabulary，但 application 仍負責 rollout policy、approval、signature verification 與 restart recovery。

## Video Cloud 權責

| Layer | 提供內容 | 不負責 |
| --- | --- | --- |
| Cloud | HTTPS signaling、current-owner MQTT/WebSocket delivery、TURN credential、session state | Live media frame 與 renderer |
| RTK SDK | Authentication、ICE、offer/answer/close、error 與 clip workflow | 完整 peer connection/media engine |
| Product app | Platform WebRTC、renderer、audio policy 與 UX lifecycle | Device camera/codec |
| Device SDK/firmware | Offer、answer、camera/audio track 與 resource limit | Product app UI |

RTK SDK 會交換 offer、answer、ICE 與 session lifecycle。Product app 不需自行建立 signaling backend，但仍需將 RTK SDK session data 接到 platform WebRTC/media component。Stored-video playback 是獨立 HTTP workflow。
