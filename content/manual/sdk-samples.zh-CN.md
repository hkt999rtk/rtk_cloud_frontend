---
title: "SDK 范例应用"
description: "使用 Realtek Connect+ SDK 范例生态系验证 App 与装置流程。"
---

![Realtek Connect+ 范例生态系](/static/assets/connectplus-sample-ecosystem-corporate-v2.jpg)

Realtek Connect+ 客户评估应先从可执行的参考范例开始，再决定正式品牌 App、量产韧体或私有部署计划。详细 source of truth 仍在 [`rtk_cloud_client`](https://github.com/hkt999rtk/rtk_cloud_client) repository，尤其是 [`docs/SAMPLE_APPLICATIONS.md`](https://github.com/hkt999rtk/rtk_cloud_client/blob/main/docs/SAMPLE_APPLICATIONS.md) 与各 sample README。

## 范例家族

| 家族 | 范例 | 目的 |
| --- | --- | --- |
| 家庭 App 参考客户端 | Android 智能家庭范例、iOS 智能家庭范例、WebApp Ops Lab 范例 | 验证 App 端 SDK 用法，包含 token setup、配网、装置列表/细节、灯具与空调控制、相机监看、debug report 与 evidence collection。 |
| 装置参考客户端 | Linux 模拟器、PRO2 相机装置范例 | 验证装置端 command handling、sample MQTT payload 行为、状态/log/event 回报、snapshot upload 与 WebRTC Video over TURN answerer / ICE/TURN 边界。 |

## 客户评估路径

1. 先跑 Linux simulator，在没有实体硬体时验证灯具、空调与相机 command handling。
2. 使用 WebApp Ops Lab 验证云端侧 onboarding、device registry 检视、MQTT payload inspection、模拟控制、相机 helper flow 与 debug report。
3. 执行 Android 与 iOS 智能家庭范例，验证 native mobile flow、setup profile、装置控制、相机监看边界与 redacted debug evidence。
4. 有 PRO2 硬体时接上 PRO2 camera device demo，验证 device-bound token、owner transport、snapshot upload、相机 status/log/event 回报与 WebRTC Video over TURN answerer 行为。

## 边界

- 这些范例是 SDK usage references，不是 production app-store apps、white-label release packages 或 customer release artifacts。
- WebApp 范例展示 browser-side cloud flow，不实现 BLE 或 SoftAP onboarding。
- 灯具与空调控制的 sample command schema 是 sample-local reference，不是正式云端 wire contract。
- Video streaming 对外文案只使用 WebRTC Video over TURN；TURN 与 coturn 仅描述为 WebRTC ICE infrastructure，不定位成独立 relay 产品。
- 产品团队仍需拥有 production UX、App 签章、App 上架、家庭分享、推播、情境、排程、自动化与市场 release policy。

## Source documents

| 主题 | 来源 |
| --- | --- |
| 范例生态系总览 | [`docs/SAMPLE_APPLICATIONS.md`](https://github.com/hkt999rtk/rtk_cloud_client/blob/main/docs/SAMPLE_APPLICATIONS.md) |
| Home app 行为 | [`docs/SAMPLE_HOME_APP_SPEC.md`](https://github.com/hkt999rtk/rtk_cloud_client/blob/main/docs/SAMPLE_HOME_APP_SPEC.md) |
| Device reference 行为 | [`docs/SAMPLE_DEVICE_APP_SPEC.md`](https://github.com/hkt999rtk/rtk_cloud_client/blob/main/docs/SAMPLE_DEVICE_APP_SPEC.md) |
| Android sample | [`samples/android/README.md`](https://github.com/hkt999rtk/rtk_cloud_client/blob/main/samples/android/README.md) |
| iOS sample | [`samples/ios/README.md`](https://github.com/hkt999rtk/rtk_cloud_client/blob/main/samples/ios/README.md) |
| WebApp sample | [`samples/webapp/README.md`](https://github.com/hkt999rtk/rtk_cloud_client/blob/main/samples/webapp/README.md) |
| Linux simulator | [`samples/linux-simulator/README.md`](https://github.com/hkt999rtk/rtk_cloud_client/blob/main/samples/linux-simulator/README.md) |
| PRO2 device demo | [`packages/freertos/pro2_demo/README.md`](https://github.com/hkt999rtk/rtk_cloud_client/blob/main/packages/freertos/pro2_demo/README.md) |

## 下一步

先用这些范例判断哪一条 App 与装置路径符合你的产品评估，再[联络 Realtek Connect+ 团队](/zh-cn/contact)讨论 SDK 范围、装置类型与部署假设。
