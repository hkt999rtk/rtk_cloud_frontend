---
title: "Video Cloud 工作流"
description: "整合 Live WebRTC signaling、encrypted clip、snapshot、browser 與 playback。"
---

![Stored video upload, browser, and playback flow](/content-assets/manual/sdk/video-flow.png)

## Live WebRTC signaling

1. App 透過 RTK SDK 呼叫 HTTPS ICE API。
2. App 的 platform WebRTC component 建立 SDP offer。
3. RTK SDK 透過 HTTPS 建立 signaling session。
4. Cloud 經目前 device owner 的 MQTT 或 WebSocket transport 送出 `webrtc_offer`；不 fan-out，也不自動 fallback 到非 owner transport。
5. Device SDK/firmware 接收 offer、連接 camera/audio track 並建立 answer。
6. Device 經 HTTPS answer API 回覆。
7. App 經 RTK SDK 取得 answer並完成 media negotiation。
8. App 或 device 關閉 session；expired、timeout、offline、busy 與 unsupported capability 都必須成為可見 terminal/error state。

Cloud 協調 signaling，但不接收或儲存 Live media frame。Live WebRTC 不會自動形成 Stored Clip。

## Encrypted clip upload

完成錄製 MP4 並加密 stored bytes，計算 ciphertext size 與 base64 SHA-256，再建立 encryption descriptor。透過 SDK authorize upload，將所有 signed header 與 encrypted bytes 送到 HTTPS presigned PUT URL，complete 後輪詢到 `ready`、`failed` 或 `expired`。目前 technical defaults 為 MP4 256 MiB、JPEG 5 MiB、signed URL 10 分鐘、upload lifecycle 30 分鐘；retention 可依 deployment 設定 1/7/30 天，均不是價格或 SLA。

目前不支援 S3 multipart upload，streaming PUT 的 read source 必須產生精確 `body_size`。

## Browser、thumbnail、playback 與 delete

Android/iOS `listClips` 回傳 typed `ClipPage`，以 `skip`、`limit`、filter 與 `nextSkip` 分頁；empty page 是正常狀態。Thumbnail 必須透過 SDK 下載。Playback session 回傳短效、支援 HTTP range 的 URL，App 交給 Media3/ExoPlayer 或 AVPlayer，過期後 refresh，不可持久保存。Delete 前確認使用者意圖，404 表示 stale/deleted clip。

Snapshot 是獨立 JPEG workflow。Live、recording、stored clip 與 snapshot 必須在 UI 與 operational state 中分開。

## Current release boundary

目前版本不包含 server-side transcoding、S3 multipart upload、simulcast negotiation、renegotiation，也不提供完整的 SDK 內建 WebRTC media renderer。
