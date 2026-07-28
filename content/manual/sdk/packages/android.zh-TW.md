---
title: "Android Kotlin SDK"
description: "整合 Android AAR、coroutine、Live signaling、clip browser、secure key 與 Media3。"
---

## 安裝與 client

使用發布的 Maven coordinate 或 package script 產生的 local repository。交付物包含 AAR、sources JAR、POM、Gradle metadata、manifest、checksum 與 consumer smoke report。啟用 Internet permission 並維持正常 TLS validation。

建立 `RtkCloudClient` 並從 ViewModel/service 使用 async 或 coroutine API。將 `RtkCloudException.status` 映射為 stable application state，diagnostics 前先 redaction response body。

## Live 與 stored video

RTK SDK 提供 authentication、ICE、offer/answer/close 與 signaling session helper；App 將 session data 接到 Android platform WebRTC component，並負責 peer connection、track、renderer、audio policy 與 lifecycle。SDK 不包含完整 media renderer。

以 `ClipQuery` 呼叫 `listClips`，處理 `ClipPage`、`nextSkip`、filter、empty state、thumbnail 與 delete confirmation。`PlaybackKeyProvider` 應由 Android Keystore-backed code 實作；將 `playbackUrl` 交給 Media3/ExoPlayer，URL 過期後更新並隨 UI lifecycle release player。

## 驗證

執行 package unit test、可用時的 instrumentation test、AAR consumer smoke、sample unit/UI smoke。Android Playback 狀態為 Real SDK + Media3；Live 頁為 RTK signaling integration/fixture UI，尚不代表完整 media rendering。
