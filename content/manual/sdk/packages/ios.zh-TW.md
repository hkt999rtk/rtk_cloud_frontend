---
title: "iOS Swift SDK"
description: "整合 Swift package、async operation、Live signaling、clip browser、protected key 與 AVPlayer。"
---

## 安裝與 client

以 approved release tag 加入 `RTKCloudClient` Swift package，或使用內部 source archive。支援 iOS 13+ 與 macOS 12+。在 network-facing service layer 建立 client，優先使用 Swift concurrency async overload；捕捉 `RTKCloudError` 並保留 stable status。

## Live 與 stored video

RTK SDK 提供 authentication、ICE、offer/answer/close 與 signaling session helper；App 將 session data 接到 iOS platform WebRTC component，並負責 peer connection、track、renderer、audio session 與 UIKit/SwiftUI lifecycle。

`listClips` 回傳 typed `ClipPage`；依 `nextSkip` 分頁並在 main actor 外載入 thumbnail。以 Keychain/Secure Enclave-backed `PlaybackKeyProvider` 回傳 wrapped clip key 與 ephemeral public key，再將 `PlaybackSession.playbackURL` 交給 `AVPlayerItem`。URL 過期或 selection 變更時建立新 session/item。

## 驗證

執行 `swift test`、`swift build` 與 simulator build/install/launch/UI test。iOS Playback 狀態為 Real SDK + AVPlayer；Live 頁為 RTK signaling integration/fixture UI，尚不代表完整 media rendering。
