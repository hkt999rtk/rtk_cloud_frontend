---
title: "Lifecycle and errors"
description: "Apply consistent client ownership, asynchronous execution, retries, cancellation, and shutdown."
---

## Client and session ownership

Create clients at an application-service boundary rather than once per request. A device session belongs to one client and must not outlive it. Disconnect should tolerate a partially connected session, and repeated shutdown must not double-free resources or emit duplicate terminal events.

## Error categories

All packages distinguish invalid arguments, invalid state, timeout, authentication, transport, protocol, unsupported capability, platform failure, cancellation, memory exhaustion, and internal failure. Preserve the stable SDK category in application telemetry. Treat raw HTTP status and response text as diagnostic context only.

## Retry policy

Retry only operations documented as idempotent or safe to repeat. Use bounded exponential backoff with jitter for temporary transport failures and 429 or 5xx responses. Do not automatically repeat direct media upload authorization, clip deletion, device activation, or command operations without an idempotency guarantee. Respect server retry guidance where provided.

## Cancellation

Propagate application cancellation through Kotlin coroutines, Swift tasks, Go contexts, JavaScript abort signals where supported, and native cancellation callbacks. Native streaming upload checks cancellation between buffer reads and reports bytes already sent. Cancellation is not a successful upload; query the upload state before deciding whether to resume or start again.

## Callbacks and threads

Assume callback payloads remain valid only for the callback duration unless the API states otherwise. Do not destroy a client from its currently dispatching callback. Move UI updates to the platform main thread and keep callback work bounded.
