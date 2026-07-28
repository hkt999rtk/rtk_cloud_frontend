---
title: "生命周期与错误"
description: "一致处理 client ownership、async execution、retry、cancellation 与 shutdown。"
---

## Ownership 与 error

在 application service boundary 创建 client，不要每个 request 创建一次。Device session 属于单一 client 且不可比它存活更久。Disconnect 必须容忍 partial connection；重复 shutdown 不应 double-free 或发出两次 terminal event。

所有 package 都区分 invalid argument/state、timeout、authentication、transport、protocol、unsupported capability、platform failure、cancellation、memory exhaustion 与 internal failure。Application telemetry 应保留 stable SDK category；raw HTTP body 只作诊断 context。

## Retry、cancel 与 callback

只 retry 明确 idempotent 或可安全重复的 operation，对暂时 transport failure、429 或 5xx 使用 bounded exponential backoff with jitter。没有 idempotency 保证时，不可自动重试 upload authorization、clip delete、device activation 或 command。

将 Kotlin coroutine、Swift Task、Go context、JavaScript AbortSignal 与 native callback 的 cancellation 传到底层。Cancellation 不等于 upload success，应先查询 upload state。Callback payload 若无另行说明，只在 callback 期间有效；UI update 必须切回 platform main thread。
