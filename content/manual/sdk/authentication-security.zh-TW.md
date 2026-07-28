---
title: "認證與安全"
description: "保護 bearer token、device certificate、encryption key、media URL 與 user data。"
---

## Token 與 PKI

只取得 operation 所需的最小 token，並在到期前更新。將 credential 保存在 process memory 或 platform-protected storage，從 log 移除 `Authorization`、cookie、query token 與 response body。401 通常需要 refresh；403 通常表示 permission 或 product capability 不足。

啟用 PKI 時，device 可用 client certificate 認證。正常驗證 server certificate 與 hostname，以 keystore 或 hardware-backed storage 保護 private key，並在到期前 rotation。Certificate identity 必須與 provisioned device 相符。

## Stored-video encryption

Device 在上傳前加密完整 clip，只將 encrypted object、integrity metadata 與 wrapped key material 交給 Cloud。`PlaybackKeyProvider` 的 recipient private key 應留在 Android Keystore、iOS Keychain/Secure Enclave 或等效 boundary。

## Presigned URL 與 logging

Presigned upload/playback URL 是短效 bearer-like secret。只送 authorization response 指定的 signed header，不可將 Video Cloud token 附加到 object-storage request。不可永久保存 URL、放入 analytics/debug report，或記錄 wrapped key、private key、raw media。可安全記錄 operation、stable status、HTTP status、correlation ID、byte count 與 elapsed time。
