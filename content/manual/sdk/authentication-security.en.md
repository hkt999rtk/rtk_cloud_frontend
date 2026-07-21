---
title: "Authentication and security"
description: "Protect bearer tokens, device certificates, encryption keys, media URLs, and user data."
---

## Bearer tokens

Request the narrowest token needed for the operation and refresh it before expiry. Keep tokens in process memory or platform-protected credential storage. Redact `Authorization`, query tokens, cookies, and token response bodies from logs. A 401 normally requires refreshing credentials; a 403 normally means the identity lacks permission or the product capability is unavailable.

## PKI and mutual TLS

Device integrations can authenticate with a client certificate where the deployment enables PKI. Validate the server certificate and hostname, protect the private key with the platform keystore or hardware-backed storage, and rotate certificates before expiry. A certificate common name or subject identity must match the provisioned device identity.

## Stored-video encryption

The device encrypts complete clip bytes before upload and sends the cloud only the encrypted object, required integrity metadata, and wrapped key material. Mobile applications implementing `PlaybackKeyProvider` keep the recipient private key in Android Keystore, iOS Keychain/Secure Enclave-backed code, or an equivalent protected boundary. A provider returns only the wrapped clip key and ephemeral public key required by the playback-session endpoint.

## Presigned URLs

Presigned upload and playback URLs are bearer-like secrets with a short lifetime. Send exactly the signed headers returned by the authorization response. Do not attach the Video Cloud bearer token to object-storage requests. Do not persist URLs, include them in analytics, or export them in debug reports.

## Logging rules

Safe diagnostics include operation name, stable SDK status, HTTP status, correlation IDs, byte counts, and elapsed time. Unsafe diagnostics include tokens, private keys, raw customer media, presigned URLs, unredacted request bodies, and wrapped key material.
