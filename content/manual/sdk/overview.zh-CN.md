---
title: "SDK 总览"
description: "选择 RTK Cloud SDK package，并理解 Cloud、SDK、Product App 与 Device 的责任边界。"
---

RTK Cloud SDK 为 device、mobile、web 与 automation application 提供一致的 client boundary。请使用 SDK，不要自行组合 Video Cloud wire request；SDK 会统一 authentication、error、timeout、payload 与 compatibility behavior。

![RTK Cloud SDK package map](/content-assets/manual/sdk/sdk-package-map.png)

## Package 选择

| Package | 主要用户 | 重要边界 |
| --- | --- | --- |
| Native C/C++ | Embedded Linux 与 native product | Stable C ABI 与 thin C++ wrapper |
| Android Kotlin | Android product app | App 负责 platform WebRTC、Media3、UI 与 secure key storage |
| iOS Swift | iPhone/iPad product app | App 负责 platform WebRTC、AVPlayer、UI 与 key policy |
| JavaScript/TypeScript | Browser 与 Node.js tool | SDK 不负责 application state 或 UI |
| Go | Device client、CI 与 automation | Pure Go；目前为 draft/internal |
| FreeRTOS/Pro2 | AmebaPro2 camera firmware | Board、media engine、storage 与 vendor SDK 由产品负责 |

## 能力与责任

Cloud 提供 HTTPS signaling、当前 owner 的 MQTT/WebSocket delivery、TURN credential、session state 与 stored-video API。RTK SDK 提供 authentication、ICE、offer/answer/close helper、stable error 与 clip workflow。Product app 集成 platform WebRTC、renderer、audio policy 与 UX lifecycle。Device SDK/firmware 负责 offer、answer、camera/audio track 与 resource limit。

当前版本不包含 server-side transcoding、S3 multipart upload、simulcast negotiation、renegotiation，也不提供完整的 SDK 内置 WebRTC media renderer。Live WebRTC 不会自动形成 Stored Clip。

## 版本

每份发布手册的 manifest 都记录 SDK commit 与 release version。请确认手册与实际 package 相符；deprecated API 仅用于兼容性维护。
