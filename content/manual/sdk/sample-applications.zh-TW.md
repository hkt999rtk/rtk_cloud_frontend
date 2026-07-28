---
title: "Sample applications"
description: "建置 reference app，並依可見狀態正確解讀 SDK integration evidence。"
---

| Sample | 可見狀態 | 驗證內容 | 不代表 |
| --- | --- | --- | --- |
| Android Playback | Real SDK + Media3 | Clip list、playback session、range playback | Live media rendering |
| iOS Playback | Real SDK + AVPlayer | Clip list、playback session、range playback | Live media rendering |
| Android/iOS Live | RTK signaling integration / fixture UI | Session data、offer/answer/close 與 UI lifecycle | 完整 WebRTC rendering |
| WebApp Ops Lab | Fixture-backed | Ops workflow 與 signaling helper demonstration | Production WebRTC peer |
| Linux simulator | Device workflow simulator | Command、state、log、report 與 snapshot behavior | Camera frame 或 WebRTC signaling |
| PRO2 host smoke | Adapter + lifecycle smoke | Adapter boundary 與 signaling lifecycle | Physical media validation |

Android 以 `gradle -p samples/android assembleDebug` 與 `testDebugUnitTest` 驗證；iOS 以 `swift test --package-path samples/ios` 與 `swift build --package-path samples/ios` 驗證。Runtime credential 與 playback URL 必須排除於 debug export。

Sample 用於複製 architecture 與 error-handling pattern，不可複製 credential 或 deployment constant。Production firmware 仍需提供 board I/O、secure storage、time、network、camera/audio media 與 vendor WebRTC implementation。
