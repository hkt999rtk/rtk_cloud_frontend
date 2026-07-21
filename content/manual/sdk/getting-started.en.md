---
title: "Getting started"
description: "Install an SDK package, configure a client, and verify a first authenticated operation."
---

## Before you start

Obtain the Video Cloud HTTPS base URL, a device identifier, and an appropriate short-lived bearer token or device certificate. Never embed production tokens, private keys, or customer media in source control. Confirm that the selected product enables the capability you plan to call.

## Integration sequence

1. Add the package using CMake, Gradle, SwiftPM, npm, or Go modules.
2. Create one client with the Video Cloud base URL and platform TLS configuration.
3. Obtain or inject authentication material.
4. call a read-only operation such as server version, camera info, or clip enumeration.
5. Map SDK errors to user-safe UI and redacted diagnostics.
6. Close sessions and clients deterministically.

## Development profiles

Use separate profiles for local simulation, staging, and production. A profile may contain public endpoints and non-secret commit identifiers. Device IDs, tokens, passwords, client certificates, private keys, presigned URLs, and wrapped media keys must be supplied at runtime and excluded from logs and reports.

## First verification

A successful first request should prove TLS hostname verification, authentication, request serialization, response parsing, and cancellation. Record only the SDK version, server version, request correlation IDs, status category, and sanitized timestamps. Do not record authorization headers or full response bodies that can contain download URLs.

## Package-specific setup

Continue with the package chapter for exact build commands and entry points. The generated API reference is authoritative for exported names at the SDK commit recorded in the documentation manifest.
