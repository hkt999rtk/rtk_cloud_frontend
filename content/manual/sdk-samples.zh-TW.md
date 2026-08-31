---
title: "SDK 範例應用"
description: "使用 Realtek Connect+ SDK 範例生態系驗證 App 與裝置流程。"
---

![Realtek Connect+ 範例生態系](/static/assets/connectplus-sample-ecosystem-corporate-v2.jpg)

Realtek Connect+ 客戶評估應先從可執行的參考範例開始，再決定正式品牌 App、量產韌體或私有部署計畫。詳細 source of truth 仍在 [`rtk_cloud_client`](https://github.com/hkt999rtk/rtk_cloud_client) repository，尤其是 [`docs/sample_applications.md`](https://github.com/hkt999rtk/rtk_cloud_client/blob/main/docs/sample_applications.md) 與各 sample README。

## 範例家族

| 家族 | 範例 | 目的 |
| --- | --- | --- |
| 家庭 App 參考客戶端 | Android 智慧家庭範例、iOS 智慧家庭範例、WebApp Ops Lab 範例 | 驗證 App 端 SDK 用法，包含 token setup、配網、裝置列表/細節、燈具與空調控制、相機監看、debug report 與 evidence collection。 |
| 裝置參考客戶端 | Linux 模擬器、PRO2 相機裝置範例 | Linux 驗證 command/state/report/snapshot metadata，不模擬 camera frame 或 WebRTC signaling；PRO2 host smoke 驗證 adapter 與 signaling lifecycle，實體硬體才驗證 media。 |

## 客戶評估路徑

1. 先跑 Linux simulator，在沒有實體硬體時驗證燈具、空調、meter、state、report 與 snapshot metadata；它不模擬 camera frame 或 WebRTC signaling。
2. 使用 WebApp Ops Lab 驗證雲端側 onboarding、device registry 檢視、MQTT payload inspection、模擬控制、相機 helper flow 與 debug report。
3. 執行 Android 與 iOS 智慧家庭範例，驗證 native mobile flow、setup profile、裝置控制、相機監看邊界與 redacted debug evidence。
4. PRO2 host smoke 只驗證 adapter 與 signaling lifecycle；有實體 PRO2 camera device 時，再驗證 snapshot、camera/audio、codec 與 WebRTC answerer media 行為。

## 邊界

- 這些範例是 SDK usage references，不是 production app-store apps、white-label release packages 或 customer release artifacts。
- WebApp 範例展示 browser-side cloud flow，不實作 BLE 或 SoftAP onboarding。
- 燈具與空調控制的 sample command schema 是 sample-local reference，不是正式雲端 wire contract。
- Video streaming 對外文案只使用 WebRTC Video over TURN；TURN 與 coturn 僅描述為 WebRTC ICE infrastructure，不定位成獨立 relay 產品。
- 產品團隊仍需擁有 production UX、App 簽章、App 上架、家庭分享、推播、情境、排程、自動化與市場 release policy。

## Source documents

| 主題 | 來源 |
| --- | --- |
| 範例生態系總覽 | [`docs/sample_applications.md`](https://github.com/hkt999rtk/rtk_cloud_client/blob/main/docs/sample_applications.md) |
| Home app 行為 | [`docs/sample_home_app_spec.md`](https://github.com/hkt999rtk/rtk_cloud_client/blob/main/docs/sample_home_app_spec.md) |
| Device reference 行為 | [`docs/sample_device_app_spec.md`](https://github.com/hkt999rtk/rtk_cloud_client/blob/main/docs/sample_device_app_spec.md) |
| Android sample | [`samples/android/README.md`](https://github.com/hkt999rtk/rtk_cloud_client/blob/main/samples/android/README.md) |
| iOS sample | [`samples/ios/README.md`](https://github.com/hkt999rtk/rtk_cloud_client/blob/main/samples/ios/README.md) |
| WebApp sample | [`samples/webapp/README.md`](https://github.com/hkt999rtk/rtk_cloud_client/blob/main/samples/webapp/README.md) |
| Linux simulator | [`samples/linux-simulator/README.md`](https://github.com/hkt999rtk/rtk_cloud_client/blob/main/samples/linux-simulator/README.md) |
| PRO2 device demo | [`packages/freertos/pro2_demo/README.md`](https://github.com/hkt999rtk/rtk_cloud_client/blob/main/packages/freertos/pro2_demo/README.md) |

## 下一步

先用這些範例判斷哪一條 App 與裝置路徑符合你的產品評估，再[聯絡 Realtek Connect+ 團隊](/contact)討論 SDK 範圍、裝置類型與部署假設。
