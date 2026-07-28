---
title: "Native C/C++ SDK"
description: "集成 stable C ABI、C++ wrapper、transport、signaling 与 direct media upload。"
---

## 安装与 ABI

Release archive 包含 `librtk_cloud_client.a`、public header 与导出 `rtk::cloud_client` 的 CMake metadata。C 使用 `rtkc.h`，C++ 使用 thin wrapper `rtkc.hpp`。所有 public request struct 必须以对应 `rtkc_*_init` 初始化；handle 视为 opaque，SDK allocation 用配对 API release。

## Lifecycle 与能力

以 validated endpoint、auth、callback、transport 与 timeout 创建 client。Session 必须依次 connect、process、disconnect、destroy，最后 destroy client。Public API 包含 token、device lifecycle、event/log、command、telemetry、firmware、snapshot、media 与 WebRTC signaling helper。

## Direct clip upload

先 authorize，再以 `rtkc_presigned_stream_put_request_t` 设置 HTTPS URL、signed header、exact body size、read/progress/cancel callback 与 bounded buffer；PUT 成功后 complete 并查询 readiness。不得把 Video Cloud bearer token 发送到 object storage，也不支持 S3 multipart upload。

## 验证

以 CMake/CTest、package consumer smoke、platform transport test 与 approved deployed-server harness 验证。WebRTC media engine、track 与 renderer 仍由 product/device integration 负责。
