---
title: "调试与兼容性"
description: "诊断 authentication、transport、signaling、upload、playback 与 version 问题。"
---

| Symptom | 常见原因 | 处理方式 |
| --- | --- | --- |
| TLS/hostname failure | Base URL、CA、clock 或 proxy 错误 | 验证 endpoint、chain、hostname 与 system time |
| 401 | Credential 丢失或过期 | Refresh token 或修正 certificate selection |
| 403 | Permission/product capability 不足 | 检查 identity、ownership 与 capability |
| 404 | 错误或已删除的 device/clip | Refresh state |
| Signaling timeout/offline | Device owner transport 不在线 | 显示 terminal state；不可 fan-out 或切到非 owner transport |
| 410 upload | 使用 legacy route | 改用 authorize、presigned PUT、complete 与 status |

Presigned PUT failure 常见于 signed header 被改动、content length 错误、token 泄漏到 object storage、URL 过期、early EOF 或 ciphertext hash mismatch。PUT 未成功不可 complete。Playback URL 过期时创建新 session；encrypted clip 必须同时提供 wrapped clip key 与 ephemeral public key。

Support evidence 可包含 SDK/server/toolchain version、operation、stable status、HTTP status、sanitized correlation ID 与 reproduction step。不可附上 token、private key、presigned URL、wrapped key 或 customer media。用 manifest 比对 SDK commit 与文档版本。
