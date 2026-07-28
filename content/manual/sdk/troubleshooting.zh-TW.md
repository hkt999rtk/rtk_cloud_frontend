---
title: "除錯與相容性"
description: "診斷 authentication、transport、signaling、upload、playback 與 version 問題。"
---

| Symptom | 常見原因 | 處理方式 |
| --- | --- | --- |
| TLS/hostname failure | Base URL、CA、clock 或 proxy 錯誤 | 驗證 endpoint、chain、hostname 與 system time |
| 401 | Credential 遺失或過期 | Refresh token 或修正 certificate selection |
| 403 | Permission/product capability 不足 | 檢查 identity、ownership 與 capability |
| 404 | 錯誤或已刪除的 device/clip | Refresh state |
| Signaling timeout/offline | Device owner transport 不在線 | 顯示 terminal state；不可 fan-out 或切到非 owner transport |
| 410 upload | 使用 legacy route | 改用 authorize、presigned PUT、complete 與 status |

Presigned PUT failure 常見於 signed header 被改動、content length 錯誤、token 洩漏到 object storage、URL 過期、early EOF 或 ciphertext hash mismatch。PUT 未成功不可 complete。Playback URL 過期時建立新 session；encrypted clip 必須同時提供 wrapped clip key 與 ephemeral public key。

Support evidence 可包含 SDK/server/toolchain version、operation、stable status、HTTP status、sanitized correlation ID 與 reproduction step。不可附上 token、private key、presigned URL、wrapped key 或 customer media。用 manifest 比對 SDK commit 與文件版本。
