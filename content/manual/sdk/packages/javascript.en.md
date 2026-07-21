---
title: "JavaScript and TypeScript SDK"
description: "Integrate the typed ESM SDK in Node.js and browser-facing applications."
---

## Install

The internal release is an npm-compatible tarball containing built ESM JavaScript, TypeScript declarations, package metadata, manifest, checksum, and isolated-consumer smoke report. Node.js 20 or later is required by the current package metadata.

## Client usage

Import from `@rtk-cloud/client`, create a client with the HTTPS endpoint and adapters, then call the typed request methods. Keep bearer tokens outside committed configuration and avoid exposing device credentials to browser code unless the architecture explicitly supports user-scoped tokens.

## Asynchronous behavior

Every network operation returns a promise. Await it inside bounded application workflows, convert SDK errors at the service boundary, and prevent unhandled promise rejections. Use cancellation support where exposed and stop subscriptions or sessions when their owner is disposed.

## Browser security

Apply a restrictive Content Security Policy, use HTTPS, and avoid localStorage for long-lived credentials. The SDK does not make an administrative token safe for browser use. Use a trusted backend when an operation requires service credentials or cross-account authority.

## Validation

Run `npm ci`, the TypeScript build, unit tests, package creation, and the isolated consumer smoke. Treat generated declarations as the public TypeScript contract and regenerate the API reference whenever exported declarations change.
