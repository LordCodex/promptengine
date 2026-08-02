---
document_id: stacks-nextjs-performance
title: Next.js Performance Engineering
ecosystem: react-next
dependencies:
  - core-universal-coding-standards
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Next.js Performance Engineering

## Inheritance & Constraints
This document inherits from the [Universal Coding Standards](../../core/05-universal-coding-standards.md). It outlines Next.js-specific optimizations, avoiding generic performance advice.

---

## 1. Caching & Data Fetching

Next.js builds cache layers by default. Ensure they are optimized correctly:
- **Server Cache Options**: Configure `revalidate` intervals or tags systematically on `fetch()` calls.
- **Request Memoization**: Next.js automatically deduplicates `fetch` requests with identical signatures inside the same render tree. Avoid manual cache wrappers for nested component fetches.
- **Opt-Out Control**: Mark queries as dynamic when cookies, request headers, or search parameters change parameters:
  - `const cookieStore = await cookies()`

---

## 2. Image Optimization (`next/image`)

- **Explicit Sizes**: Always use the `<Image />` component with set width/height attributes or the `fill` property to prevent Layout Shifts (CLS).
- **Srcset Resolution**: Provide responsive `sizes` mappings (e.g. `sizes="(max-width: 768px) 100vw, 33vw"`) to let the server compress image files matching target viewports.
- **Priority Loading**: Apply `priority` tags to hero images visible on initial viewport load.

---

## 3. SEO & Metadata API

- **Static Metadata**: Export `metadata` objects in layouts/pages.
- **Dynamic Metadata**: Export `generateMetadata` functions to resolve dynamic page parameters (e.g., product titles).
- **Multi-language canonicals**: Map corresponding `alternates` keys to configure `hreflang` tags.

---

## Review Checklist

Verify page performance against this checklist:
- [ ] Inherits correctly from Universal Coding Standards.
- [ ] Images utilize Next.js `<Image>` tags with set `sizes` parameters.
- [ ] Critical viewports contain `priority` image triggers.
- [ ] Metadata configurations are active.
