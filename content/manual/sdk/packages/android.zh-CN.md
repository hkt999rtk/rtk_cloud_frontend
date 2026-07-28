---
title: "Android Kotlin SDK"
description: "集成 Android AAR、coroutine、Live signaling、clip browser、secure key 与 Media3。"
---

## 安装与 client

使用发布的 Maven coordinate 或 package script 产生的 local repository。交付物包含 AAR、sources JAR、POM、Gradle metadata、manifest、checksum 与 consumer smoke report。启用 Internet permission 并保持正常 TLS validation。

创建 `RtkCloudClient` 并从 ViewModel/service 使用 async 或 coroutine API。将 `RtkCloudException.status` 映射为 stable application state，diagnostics 前先 redaction response body。

## Live 与 stored video

RTK SDK 提供 authentication、ICE、offer/answer/close 与 signaling session helper；App 将 session data 接到 Android platform WebRTC component，并负责 peer connection、track、renderer、audio policy 与 lifecycle。SDK 不包含完整 media renderer。

以 `ClipQuery` 调用 `listClips`，处理 `ClipPage`、`nextSkip`、filter、empty state、thumbnail 与 delete confirmation。`PlaybackKeyProvider` 应由 Android Keystore-backed code 实现；将 `playbackUrl` 交给 Media3/ExoPlayer，URL 过期后更新并随 UI lifecycle release player。

## 验证

执行 package unit test、可用时的 instrumentation test、AAR consumer smoke、sample unit/UI smoke。Android Playback 状态为 Real SDK + Media3；Live 页面为 RTK signaling integration/fixture UI，尚不代表完整 media rendering。
