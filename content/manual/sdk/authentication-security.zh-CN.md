---
title: "认证与安全"
description: "保护 bearer token、device certificate、encryption key、media URL 与 user data。"
---

## Token 与 PKI

只取得 operation 所需的最小 token，并在到期前更新。将 credential 保存在 process memory 或 platform-protected storage，从 log 移除 `Authorization`、cookie、query token 与 response body。401 通常需要 refresh；403 通常表示 permission 或 product capability 不足。

启用 PKI 时，device 可用 client certificate 认证。正常验证 server certificate 与 hostname，以 keystore 或 hardware-backed storage 保护 private key，并在到期前 rotation。Certificate identity 必须与 provisioned device 相符。

## Stored-video encryption

Device 在上传前加密完整 clip，只将 encrypted object、integrity metadata 与 wrapped key material 交给 Cloud。`PlaybackKeyProvider` 的 recipient private key 应留在 Android Keystore、iOS Keychain/Secure Enclave 或等效 boundary。

## Presigned URL 与 logging

Presigned upload/playback URL 是短效 bearer-like secret。只发送 authorization response 指定的 signed header，不可将 Video Cloud token 附加到 object-storage request。不可永久保存 URL、放入 analytics/debug report，或记录 wrapped key、private key、raw media。可安全记录 operation、stable status、HTTP status、correlation ID、byte count 与 elapsed time。
