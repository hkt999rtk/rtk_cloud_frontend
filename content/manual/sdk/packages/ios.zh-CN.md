---
title: "iOS Swift SDK"
description: "集成 Swift package、async operation、Live signaling、clip browser、protected key 与 AVPlayer。"
---

## 安装与 client

以 approved release tag 加入 `RTKCloudClient` Swift package，或使用内部 source archive。支持 iOS 13+ 与 macOS 12+。在 network-facing service layer 创建 client，优先使用 Swift concurrency async overload；捕获 `RTKCloudError` 并保留 stable status。

## Live 与 stored video

RTK SDK 提供 authentication、ICE、offer/answer/close 与 signaling session helper；App 将 session data 接到 iOS platform WebRTC component，并负责 peer connection、track、renderer、audio session 与 UIKit/SwiftUI lifecycle。

`listClips` 返回 typed `ClipPage`；按 `nextSkip` 分页并在 main actor 外加载 thumbnail。以 Keychain/Secure Enclave-backed `PlaybackKeyProvider` 返回 wrapped clip key 与 ephemeral public key，再将 `PlaybackSession.playbackURL` 交给 `AVPlayerItem`。URL 过期或 selection 变更时创建新 session/item。

## 验证

执行 `swift test`、`swift build` 与 simulator build/install/launch/UI test。iOS Playback 状态为 Real SDK + AVPlayer；Live 页面为 RTK signaling integration/fixture UI，尚不代表完整 media rendering。
