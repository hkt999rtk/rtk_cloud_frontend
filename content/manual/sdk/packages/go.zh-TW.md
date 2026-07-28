---
title: "Go SDK"
description: "使用 draft pure Go SDK 整合 device client、automation、telemetry、firmware 與 signaling。"
---

## Status 與安裝

Go SDK 是 draft/internal package，不在第一波五 package user delivery bundle。以 Go 1.21+ import `github.com/hkt999rtk/rtk_cloud_client/packages/golang/rtkc`，並 pin approved commit/release。

## Client、auth 與能力

以 `rtkc.NewClient`/`ClientConfig` 建立 client。所有 network call 接受 `context.Context`；傳遞 deadline/cancellation 並確定性地 close session/client。使用 token refresh helper 或 `auth` package 的 PKI CSR/mTLS helper，保護 PEM/PKCS#11 key。以 `rtkc.WithCorrelation` 加入 sanitized correlation ID。

Module 包含 device lifecycle、owner WebSocket transport、MQTT hook、WebRTC signaling、telemetry 與 firmware vocabulary。Application 仍負責 media engine、persistent state、rollout policy 與 product authorization。

## 驗證

執行 `CGO_ENABLED=0 go build ./...` 與 `CGO_ENABLED=0 go test ./...`。Release manifest 未改變 status 前，不宣稱 public distribution。
