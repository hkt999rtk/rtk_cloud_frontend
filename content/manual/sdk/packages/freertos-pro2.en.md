---
title: "FreeRTOS and Pro2 SDK"
description: "Integrate the Pro2 source demo with board, transport, storage, and WebRTC adapters."
---

## Delivery model

The FreeRTOS/Pro2 deliverable is a device-demo source bundle, not a prebuilt production firmware image. It includes public demo headers, adapter boundaries, host smoke tests, package metadata, manifest, and checksum. The product team supplies the approved vendor SDK and ASDK revisions.

## Board adapter

Implement networking, TLS, secure storage, time, randomness, logging, task synchronization, media storage, and device identity through the board adapter contract. Do not add vendor-specific types to the portable public boundary. Enforce bounded allocations and deterministic cleanup.

## Device workflow

Initialize the board and RTK client configuration, establish device authentication, connect the owner transport, publish supported state/events, receive commands, and shut down cleanly. Snapshot and clip operations must respect memory limits and never place a complete large clip in RAM when streaming storage reads are available.

## WebRTC boundary

The demo exposes an answerer integration boundary. Vendor media capture, codecs, peer connection, ICE, TURN, rendering, and audio/video policy are outside the portable SDK. Validate the exact vendor baseline before claiming board support.

## Validation

Run the CMake host smoke tests for the demo, board adapter, and WebRTC adapter. Record physical-board, camera-sensor, network, and vendor-media validation separately. A host smoke pass does not prove that firmware flashes, boots, captures media, or interoperates with a deployed cloud.
