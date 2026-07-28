---
title: "FreeRTOS/Pro2 SDK"
description: "集成 Pro2 source demo 的 board、transport、storage 与 WebRTC adapter。"
---

## 交付与 board adapter

FreeRTOS/Pro2 是 device-demo source bundle，不是 production firmware image。交付物包含 public demo header、adapter boundary、host smoke、package metadata、manifest 与 checksum；product team 提供 approved vendor SDK/ASDK。

Board adapter 必须实现 network、TLS、secure storage、time、random、log、task sync、media storage 与 device identity。保持 portable public boundary、bounded allocation 与 deterministic cleanup。Clip operation 应 streaming read，避免把大型 clip 全部放入 RAM。

## WebRTC boundary

Demo 暴露 answerer integration boundary。Device firmware 接收 cloud 通过 current-owner transport 送达的 offer，连接 camera/audio track，创建 answer 并经 HTTPS 回复。Vendor media capture、codec、peer connection、ICE/TURN 与 media policy 都不在 portable SDK 内。

## 验证

CMake host smoke 只验证 demo、board adapter 与 signaling lifecycle；不代表 firmware 可 flash/boot、camera frame、audio、codec 或 deployed-cloud media interoperability。Physical-board validation 必须单独记录。
