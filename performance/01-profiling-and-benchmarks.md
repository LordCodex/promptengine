---
document_id: performance-profiling
title: Application Profiling and Benchmarking
ecosystem: cross-cutting
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Application Profiling and Benchmarking

## Purpose
This document establishes standards for profiling applications, finding memory leaks, and establishing performance benchmarks before writing optimization patches.

## Scope
Applies to backend trace profiles (PHP Xdebug, Blackfire) and client-side profiling tools (Vue DevTools, Flutter DevTools Profiler).

---

## Directives

### 1. Establish Performance Baselines
- **Rule**: Do not attempt to optimize code unless you have collected metrics proving a bottleneck exists.
- Measure and record baseline values:
  - **Latency**: Mean, 95th percentile, and 99th percentile response times.
  - **Memory Allocation**: Peak RAM usage during the request execution.
  - **Database Queries**: Total count of queries and execution time.

### 2. Standard Profiling Tools
Depending on the technology stack, utilize the following tools:

| Ecosystem | Profiler Tool | Focus Metric |
| :--- | :--- | :--- |
| **PHP / Laravel** | Xdebug / Blackfire.io | CPU call count, memory allocation traces |
| **JS / TS / Vue** | Chrome DevTools Performance Panel | Heap snapshots, CPU flamegraphs, render lag |
| **Dart / Flutter** | Flutter DevTools Profiler | Widget build time, UI thread vs Raster thread frames |

### 3. Memory Leak Prevention
- **Vue**: Always destroy listeners, timers (`setInterval`), and event emitters inside the `onUnmounted` lifecycle hooks.
- **Flutter**: Dispose of `ScrollController` and `TextEditingController` instances within stateful widget `dispose` actions (refer to [stacks/dart-flutter/flutter-widgets.md](file:///Users/kodexkode/Documents/workspace/promptengine/stacks/dart-flutter/flutter-widgets.md)).
- **PHP**: Avoid generating huge static arrays that persist across long-running queue processes or octane servers.

---

## Common Mistakes & Anti-Patterns
- **Premature Caching**: Wrapping a slow function in a cache wrapper without diagnosing why the function is slow. This hides database inefficiencies and leads to stale data bugs.
- **Optimizing Non-Critical Paths**: Spending days optimizing a console script that runs once a week for 5 seconds. Focus optimizations on high-traffic web endpoints.
- **Ignoring GC Overhead**: Allocating thousands of short-lived objects inside hot loops (e.g. inside Flutter's `build` method), causing frequent Garbage Collector pauses that degrade frame rates.

---

## References
- Core testing standards: [core/04-testing-philosophy.md](file:///Users/kodexkode/Documents/workspace/promptengine/core/04-testing-philosophy.md)
- Flutter widgets lifecycle: [stacks/dart-flutter/flutter-widgets.md](file:///Users/kodexkode/Documents/workspace/promptengine/stacks/dart-flutter/flutter-widgets.md)
