---
title: "開始使用"
description: "安裝 SDK package、設定 client，並驗證第一次 authenticated operation。"
---

## 準備

取得 Video Cloud HTTPS base URL、device identifier，以及適當的 short-lived bearer token 或 device certificate。確認產品已啟用需要的 capability。不可將 production token、private key 或 customer media 放入 source control。

## 整合順序

1. 以 CMake、Gradle、SwiftPM、npm 或 Go modules 加入 package。
2. 以 Video Cloud base URL 與正常 TLS validation 建立一個 client。
3. 在 runtime 注入 authentication material。
4. 呼叫 server version、camera info 或 clip list 等 read-only operation。
5. 將 SDK error 映射為 user-safe UI 與 redacted diagnostics。
6. 確定性地關閉 session 與 client。

## Environment 與驗證

Local simulation、staging 與 production 使用獨立 profile。Device ID、token、password、certificate、private key、presigned URL 與 wrapped media key 必須由 runtime 提供並排除於 log/report。第一次成功 request 應證明 TLS hostname validation、authentication、serialization、parsing 與 cancellation；只記錄 SDK/server version、correlation ID、status category 與 sanitized timestamp。
