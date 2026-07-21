---
title: "Go SDK"
description: "Use the draft pure Go SDK for device clients, automation, telemetry, firmware, and signaling."
---

## Status

The Go SDK is a draft/internal package and is not in the first-wave five-package user delivery bundle. It is documented so Go device clients, administration tools, CI, and server-side automation can use its current supported surface without CGo.

## Install

Import `github.com/hkt999rtk/rtk_cloud_client/packages/golang/rtkc` and the required subpackages. The module targets Go 1.21 or later. Pin an approved commit or release because the public scope is still draft.

## Client and context

Create a client with `rtkc.NewClient` and `ClientConfig`. All network calls accept `context.Context`; propagate request deadlines and cancellation rather than creating background contexts inside integration code. Close sessions and the client deterministically.

## Authentication

Use token request and refresh helpers for bearer flows. The `auth` package provides PKI CSR and mTLS configuration helpers. Keep PEM keys protected or bridge an approved platform/PKCS#11 key provider. Attach correlation IDs with `rtkc.WithCorrelation` for sanitized cross-service diagnostics.

## Capabilities

The module includes device lifecycle, owner WebSocket transport, MQTT hooks, WebRTC signaling, telemetry helpers, and firmware campaign vocabulary. The application owns media engines, rollout policy, persistent state, and product authorization decisions.

## Validation

Run `CGO_ENABLED=0 go build ./...` and `CGO_ENABLED=0 go test ./...`. Do not claim public distribution or compatibility until the release manifest changes the package status.
