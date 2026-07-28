---
title: "SDK 總覽"
description: "選擇 RTK Cloud SDK package，並理解 Cloud、SDK、Product App 與 Device 的責任邊界。"
---

RTK Cloud SDK 為 device、mobile、web 與 automation application 提供一致的 client boundary。請使用 SDK，不要自行組合 Video Cloud wire request；SDK 會統一 authentication、error、timeout、payload 與 compatibility behavior。

![RTK Cloud SDK package map](/content-assets/manual/sdk/sdk-package-map.png)

## Package 選擇

| Package | 主要使用者 | 重要邊界 |
| --- | --- | --- |
| Native C/C++ | Embedded Linux 與 native product | Stable C ABI 與 thin C++ wrapper |
| Android Kotlin | Android product app | App 負責 platform WebRTC、Media3、UI 與 secure key storage |
| iOS Swift | iPhone/iPad product app | App 負責 platform WebRTC、AVPlayer、UI 與 key policy |
| JavaScript/TypeScript | Browser 與 Node.js tool | SDK 不負責 application state 或 UI |
| Go | Device client、CI 與 automation | Pure Go；目前為 draft/internal |
| FreeRTOS/Pro2 | AmebaPro2 camera firmware | Board、media engine、storage 與 vendor SDK 由產品負責 |

## 能力與權責

Cloud 提供 HTTPS signaling、目前 owner 的 MQTT/WebSocket delivery、TURN credential、session state 與 stored-video API。RTK SDK 提供 authentication、ICE、offer/answer/close helper、stable error 與 clip workflow。Product app 整合 platform WebRTC、renderer、audio policy 與 UX lifecycle。Device SDK/firmware 負責 offer、answer、camera/audio track 與 resource limit。

目前版本不包含 server-side transcoding、S3 multipart upload、simulcast negotiation、renegotiation，也不提供完整的 SDK 內建 WebRTC media renderer。Live WebRTC 不會自動形成 Stored Clip。

## 版本

每份發布手冊的 manifest 都記錄 SDK commit 與 release version。請確認手冊與實際 package 相符；deprecated API 僅供相容性維持。
