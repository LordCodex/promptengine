---
document_id: performance-concurrency-async
title: Concurrency and Asynchronous Processing
ecosystem: cross-cutting
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Concurrency and Asynchronous Processing

## Purpose
This document outlines standards for writing non-blocking asynchronous operations, isolating heavy computation threads, and managing event loops.

## Scope
Applies to JavaScript asynchronous functions, Dart Isolates/Streams, and modern PHP multi-process runtime configurations (Swoole, RoadRunner, Laravel Octane).

---

## Directives

### 1. The Event Loop Principle
- **Rule**: Never block the main event loop / UI thread with heavy calculations or synchronous I/O.
- **JavaScript (Vue)**: Use asynchronous promises and `await` operators for api operations. Avoid synchronous loops (`for (let i = 0; i < 1e9; i++)`) on the client thread.
- **Dart (Flutter)**: Keep widget build methods pure. If parsing large JSON files (> 1MB) or running encryption logic, spawn a separate background **Isolate** to prevent frame drops.

### 2. Multi-Process PHP (Swoole / Laravel Octane)
- Modern PHP runtimes boot the application code only once in memory and process concurrent requests across persistent state workers.
- **Memory Leak Alert**: Never store request-specific parameters in static class properties or singletons. Since the process memory persists across requests, state will leak between separate users.
- **Reset State**: Always clear local state in Service Providers or Middleware after each request terminates.

### 3. Stream Backpressure
- When using Streams or Event Observers, implement throttling or debouncing to prevent client-side performance degradation.
- **Flutter UI Input**: Debounce text input streams (e.g. search auto-complete) by at least 300ms before dispatching HTTP calls.

---

## Common Mistakes & Anti-Patterns
- **Octane State Pollution**: Storing an authenticated user instance inside a singleton service class, exposing that user's session data to subsequent requests.
- **Blocking the Main Thread**: Running computational encryption (e.g. bcrypt/scrypt) directly on the Flutter main UI thread, causing the screen to freeze.
- **Unhandled Promise Rejections**: Failing to attach `.catch()` blocks or `try/catch` wrappers around asynchronous calls, leading to silent application state crashes.

---

## References
- Code isolation: [core/02-architecture-and-simplicity.md](file:///Users/kodexkode/Documents/workspace/promptengine/core/02-architecture-and-simplicity.md)
- Flutter performance: [stacks/dart-flutter/flutter-widgets.md](file:///Users/kodexkode/Documents/workspace/promptengine/stacks/dart-flutter/flutter-widgets.md)
