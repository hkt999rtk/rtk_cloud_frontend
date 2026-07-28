---
title: "Go SDK"
description: "使用 draft pure Go SDK 集成 device client、automation、telemetry、firmware 与 signaling。"
---

## Status 与安装

Go SDK 是 draft/internal package，不在第一批五 package user delivery bundle。以 Go 1.21+ import `github.com/hkt999rtk/rtk_cloud_client/packages/golang/rtkc`，并 pin approved commit/release。

## Client、auth 与能力

以 `rtkc.NewClient`/`ClientConfig` 创建 client。所有 network call 接受 `context.Context`；传递 deadline/cancellation 并确定性地 close session/client。使用 token refresh helper 或 `auth` package 的 PKI CSR/mTLS helper，保护 PEM/PKCS#11 key。以 `rtkc.WithCorrelation` 加入 sanitized correlation ID。

Module 包含 device lifecycle、owner WebSocket transport、MQTT hook、WebRTC signaling、telemetry 与 firmware vocabulary。Application 仍负责 media engine、persistent state、rollout policy 与 product authorization。

## 验证

执行 `CGO_ENABLED=0 go build ./...` 与 `CGO_ENABLED=0 go test ./...`。Release manifest 未改变 status 前，不宣称 public distribution。
