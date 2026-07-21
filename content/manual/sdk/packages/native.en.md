---
title: "Native C and C++ SDK"
description: "Integrate the native SDK, stable C ABI, C++ wrapper, transports, and direct media upload."
---

## Install and link

The release archive contains `librtk_cloud_client.a`, public headers, and CMake package metadata exporting `rtk::cloud_client`. Configure the required platform and TLS callbacks, then link the imported target. Include `rtk_cloud_client/rtkc.h` for C or `rtk_cloud_client/rtkc.hpp` for the header-only C++ wrapper.

## ABI conventions

Initialize every public request structure with its matching `rtkc_*_init` function so `struct_size`, defaults, and reserved fields are correct. Treat client and session handles as opaque. Pair SDK allocations with the documented release function and do not retain callback views after callback return.

## Client lifecycle

Create `rtkc_client_t` with validated endpoint, authentication, callback, transport, timeout, and platform configuration. Create or attach sessions beneath the client, connect, drive event processing using the documented blocking entry points, disconnect, destroy the session, and finally destroy the client.

## HTTP and device operations

The public client includes token, device lifecycle, configuration, event, log, command, telemetry, firmware, snapshot, media, and WebRTC signaling helpers. JSON-returning functions populate `rtkc_json_document_t`; inspect status before reading the document and release it using its matching API.

## Streaming direct clip upload

Use `rtkc_client_authorize_direct_clip_upload`, then populate `rtkc_presigned_stream_put_request_t` with the HTTPS URL, signed headers, exact body size, read callback, progress callback, cancellation callback, user data, and buffer size. Call `rtkc_client_put_presigned_object_stream`, complete with `rtkc_client_complete_direct_clip_upload`, and query readiness with `rtkc_client_get_direct_clip_upload`.

Only HTTPS presigned URLs are accepted. The object-storage request must contain no Video Cloud bearer token. The read source must produce exactly `body_size` bytes and remain valid for the synchronous call.

## Validation

Build with CMake and run CTest. The package smoke project validates installed headers, library, and CMake metadata. Run platform transport tests and the deployed-server harness separately when credentials and an approved environment are available.
