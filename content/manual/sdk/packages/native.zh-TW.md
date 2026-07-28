---
title: "Native C/C++ SDK"
description: "整合 stable C ABI、C++ wrapper、transport、signaling 與 direct media upload。"
---

## 安裝與 ABI

Release archive 包含 `librtk_cloud_client.a`、public header 與匯出 `rtk::cloud_client` 的 CMake metadata。C 使用 `rtkc.h`，C++ 使用 thin wrapper `rtkc.hpp`。所有 public request struct 必須以對應 `rtkc_*_init` 初始化；handle 視為 opaque，SDK allocation 用配對 API release。

## Lifecycle 與能力

以 validated endpoint、auth、callback、transport 與 timeout 建立 client。Session 必須依序 connect、process、disconnect、destroy，最後 destroy client。Public API 包含 token、device lifecycle、event/log、command、telemetry、firmware、snapshot、media 與 WebRTC signaling helper。

## Direct clip upload

先 authorize，再以 `rtkc_presigned_stream_put_request_t` 設定 HTTPS URL、signed header、exact body size、read/progress/cancel callback 與 bounded buffer；PUT 成功後 complete 並查詢 readiness。不得把 Video Cloud bearer token 送到 object storage，也不支援 S3 multipart upload。

## 驗證

以 CMake/CTest、package consumer smoke、platform transport test 與 approved deployed-server harness 驗證。WebRTC media engine、track 與 renderer 仍由 product/device integration 負責。
