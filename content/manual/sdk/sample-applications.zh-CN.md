---
title: "Sample applications"
description: "构建 reference app，并按可见状态正确解读 SDK integration evidence。"
---

| Sample | 可见状态 | 验证内容 | 不代表 |
| --- | --- | --- | --- |
| Android Playback | Real SDK + Media3 | Clip list、playback session、range playback | Live media rendering |
| iOS Playback | Real SDK + AVPlayer | Clip list、playback session、range playback | Live media rendering |
| Android/iOS Live | RTK signaling integration / fixture UI | Session data、offer/answer/close 与 UI lifecycle | 完整 WebRTC rendering |
| WebApp Ops Lab | Fixture-backed | Ops workflow 与 signaling helper demonstration | Production WebRTC peer |
| Linux simulator | Device workflow simulator | Command、state、log、report 与 snapshot behavior | Camera frame 或 WebRTC signaling |
| PRO2 host smoke | Adapter + lifecycle smoke | Adapter boundary 与 signaling lifecycle | Physical media validation |

Android 以 `gradle -p samples/android assembleDebug` 与 `testDebugUnitTest` 验证；iOS 以 `swift test --package-path samples/ios` 与 `swift build --package-path samples/ios` 验证。Runtime credential 与 playback URL 必须排除于 debug export。

Sample 用于复制 architecture 与 error-handling pattern，不可复制 credential 或 deployment constant。Production firmware 仍需提供 board I/O、secure storage、time、network、camera/audio media 与 vendor WebRTC implementation。
