---
title: "JavaScript/TypeScript SDK"
description: "在 Node.js 与 browser-facing application 集成 typed ESM SDK。"
---

## 安装与使用

内部 release 是 npm-compatible tarball，包含 built ESM JavaScript、TypeScript declaration、metadata、manifest、checksum 与 isolated-consumer smoke report；当前需要 Node.js 20+。从 `@rtk-cloud/client` import，以 HTTPS endpoint 与 adapter 创建 client。Bearer token 不可放入 committed config，除非架构明确支持 user-scoped token，否则不可把 device/service credential 暴露给 browser。

所有 network operation 返回 Promise。请在 bounded workflow 中 `await`、在 service boundary 转换 SDK error、处理 cancellation，并在 owner dispose 时停止 subscription/session。

## Video 与 browser security

WebApp Ops Lab 是 fixture-backed operations 与 signaling helper demonstration，不是 production WebRTC peer。Product app 若要呈现 live media，仍需集成 browser WebRTC peer connection 与 renderer。使用 HTTPS、restrictive CSP，且不可用 localStorage 保存 long-lived credential。

## 验证

执行 `npm ci`、TypeScript build、unit test、package build 与 isolated consumer smoke；export declaration 改变时重新生成 API reference。
