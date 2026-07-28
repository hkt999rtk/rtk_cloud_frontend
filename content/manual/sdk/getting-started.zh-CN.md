---
title: "开始使用"
description: "安装 SDK package、设置 client，并验证第一次 authenticated operation。"
---

## 准备

取得 Video Cloud HTTPS base URL、device identifier，以及适当的 short-lived bearer token 或 device certificate。确认产品已启用需要的 capability。不可将 production token、private key 或 customer media 放入 source control。

## 集成顺序

1. 以 CMake、Gradle、SwiftPM、npm 或 Go modules 加入 package。
2. 以 Video Cloud base URL 与正常 TLS validation 创建一个 client。
3. 在 runtime 注入 authentication material。
4. 调用 server version、camera info 或 clip list 等 read-only operation。
5. 将 SDK error 映射为 user-safe UI 与 redacted diagnostics。
6. 确定性地关闭 session 与 client。

## Environment 与验证

Local simulation、staging 与 production 使用独立 profile。Device ID、token、password、certificate、private key、presigned URL 与 wrapped media key 必须由 runtime 提供并排除于 log/report。第一次成功 request 应证明 TLS hostname validation、authentication、serialization、parsing 与 cancellation；只记录 SDK/server version、correlation ID、status category 与 sanitized timestamp。
