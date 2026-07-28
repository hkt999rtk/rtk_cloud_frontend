---
title: "生命週期與錯誤"
description: "一致處理 client ownership、async execution、retry、cancellation 與 shutdown。"
---

## Ownership 與 error

在 application service boundary 建立 client，不要每個 request 建立一次。Device session 屬於單一 client 且不可比它存活更久。Disconnect 必須容忍 partial connection；重複 shutdown 不應 double-free 或發出兩次 terminal event。

所有 package 都區分 invalid argument/state、timeout、authentication、transport、protocol、unsupported capability、platform failure、cancellation、memory exhaustion 與 internal failure。Application telemetry 應保留 stable SDK category；raw HTTP body 只作診斷 context。

## Retry、cancel 與 callback

只 retry 明確 idempotent 或可安全重複的 operation，對暫時 transport failure、429 或 5xx 使用 bounded exponential backoff with jitter。未有 idempotency 保證時，不可自動重試 upload authorization、clip delete、device activation 或 command。

將 Kotlin coroutine、Swift Task、Go context、JavaScript AbortSignal 與 native callback 的 cancellation 傳到底層。Cancellation 不等於 upload success，應先查 upload state。Callback payload 若無另行說明，只在 callback 期間有效；UI update 必須切回 platform main thread。
