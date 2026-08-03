# Performance Optimization Guide

PromptEngine is designed to run instantly. This document describes the caching, timing, and optimization configurations used to keep performance high even in monorepos.

---

## 1. Startup Optimization

PromptEngine is compiled into a single static Go binary. Unlike Node.js or Python alternatives:
- No runtime compilation overhead.
- No heavy package dependencies loaded dynamically at execution time.
- Cold-start timing is consistently under **15 milliseconds**.

---

## 2. Platform Caching

To avoid parsing complex markdown specifications or traversing massive monorepo file trees repeatedly:
- The discovery cache stores identified stack metadata in `.promptengine/cache/`.
- File content hashing (SHA256) is used to detect filesystem drifts. If no files have changed, the cached model is returned instantly.

---

## 3. Large Repositories & Monorepos

When executing in huge repositories:
- Ignore massive build directories (`node_modules/`, `vendor/`, `build/`) during traversal to minimize disk I/O.
- Run discovery stages in parallel where possible.
- Use the `context` command to isolate the exact subset of documents needed for a task, keeping prompt payloads lightweight.
