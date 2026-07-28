---
title: "Video Cloud 工作流"
description: "集成 Live WebRTC signaling、encrypted clip、snapshot、browser 与 playback。"
---

![Stored video upload, browser, and playback flow](/content-assets/manual/sdk/video-flow.png)

## Live WebRTC signaling

1. App 通过 RTK SDK 调用 HTTPS ICE API。
2. App 的 platform WebRTC component 创建 SDP offer。
3. RTK SDK 通过 HTTPS 创建 signaling session。
4. Cloud 经当前 device owner 的 MQTT 或 WebSocket transport 发送 `webrtc_offer`；不 fan-out，也不自动 fallback 到非 owner transport。
5. Device SDK/firmware 接收 offer、连接 camera/audio track 并创建 answer。
6. Device 经 HTTPS answer API 回复。
7. App 经 RTK SDK 取得 answer 并完成 media negotiation。
8. App 或 device 关闭 session；expired、timeout、offline、busy 与 unsupported capability 都必须成为可见 terminal/error state。

Cloud 协调 signaling，但不接收或存储 Live media frame。Live WebRTC 不会自动形成 Stored Clip。

## Encrypted clip upload

完成录制 MP4 并加密 stored bytes，计算 ciphertext size 与 base64 SHA-256，再创建 encryption descriptor。通过 SDK authorize upload，将所有 signed header 与 encrypted bytes 发送到 HTTPS presigned PUT URL，complete 后轮询到 `ready`、`failed` 或 `expired`。当前 technical defaults 为 MP4 256 MiB、JPEG 5 MiB、signed URL 10 分钟、upload lifecycle 30 分钟；retention 可按 deployment 设置 1/7/30 天，均不是价格或 SLA。

当前不支持 S3 multipart upload，streaming PUT 的 read source 必须产生精确 `body_size`。

## Browser、thumbnail、playback 与 delete

Android/iOS `listClips` 返回 typed `ClipPage`，以 `skip`、`limit`、filter 与 `nextSkip` 分页；empty page 是正常状态。Thumbnail 必须通过 SDK 下载。Playback session 返回短效、支持 HTTP range 的 URL，App 交给 Media3/ExoPlayer 或 AVPlayer，过期后 refresh，不可永久保存。Delete 前确认用户意图，404 表示 stale/deleted clip。

Snapshot 是独立 JPEG workflow。Live、recording、stored clip 与 snapshot 必须在 UI 与 operational state 中分开。

## Current release boundary

当前版本不包含 server-side transcoding、S3 multipart upload、simulcast negotiation、renegotiation，也不提供完整的 SDK 内置 WebRTC media renderer。
