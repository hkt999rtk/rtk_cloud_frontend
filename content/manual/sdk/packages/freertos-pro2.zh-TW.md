---
title: "FreeRTOS/Pro2 SDK"
description: "整合 Pro2 source demo 的 board、transport、storage 與 WebRTC adapter。"
---

## 交付與 board adapter

FreeRTOS/Pro2 是 device-demo source bundle，不是 production firmware image。交付物包含 public demo header、adapter boundary、host smoke、package metadata、manifest 與 checksum；product team 提供 approved vendor SDK/ASDK。

Board adapter 必須實作 network、TLS、secure storage、time、random、log、task sync、media storage 與 device identity。保持 portable public boundary、bounded allocation 與 deterministic cleanup。Clip operation 應 streaming read，避免把大型 clip 全部放入 RAM。

## WebRTC boundary

Demo 暴露 answerer integration boundary。Device firmware 接收 cloud 透過 current-owner transport 送達的 offer，連接 camera/audio track，建立 answer並經 HTTPS 回覆。Vendor media capture、codec、peer connection、ICE/TURN 與 media policy 都不在 portable SDK 內。

## 驗證

CMake host smoke 只驗證 demo、board adapter 與 signaling lifecycle；不代表 firmware 可 flash/boot、camera frame、audio、codec 或 deployed-cloud media interoperability。Physical-board validation 必須分開記錄。
