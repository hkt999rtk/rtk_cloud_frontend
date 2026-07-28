---
title: "JavaScript/TypeScript SDK"
description: "在 Node.js 與 browser-facing application 整合 typed ESM SDK。"
---

## 安裝與使用

內部 release 是 npm-compatible tarball，包含 built ESM JavaScript、TypeScript declaration、metadata、manifest、checksum 與 isolated-consumer smoke report；目前需要 Node.js 20+。從 `@rtk-cloud/client` import，以 HTTPS endpoint 與 adapter 建立 client。Bearer token 不可放入 committed config，除非架構明確支援 user-scoped token，否則不可把 device/service credential 暴露給 browser。

所有 network operation 回傳 Promise。請在 bounded workflow 中 `await`、於 service boundary 轉換 SDK error、處理 cancellation，並在 owner dispose 時停止 subscription/session。

## Video 與 browser security

WebApp Ops Lab 是 fixture-backed operations 與 signaling helper demonstration，不是 production WebRTC peer。Product app 若要呈現 live media，仍需整合 browser WebRTC peer connection 與 renderer。使用 HTTPS、restrictive CSP，且不可用 localStorage 保存 long-lived credential。

## 驗證

執行 `npm ci`、TypeScript build、unit test、package build 與 isolated consumer smoke；export declaration 改變時重新產生 API reference。
