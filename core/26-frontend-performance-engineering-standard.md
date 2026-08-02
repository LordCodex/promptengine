---
document_id: core-frontend-performance
title: Frontend Performance Engineering Standard
ecosystem: cross-cutting
dependencies:
  - core-performance-engineering-standard
  - core-frontend-architecture
  - core-frontend-security
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Frontend Performance Engineering Standard

## Inheritance
This document inherits from and extends the [Performance Engineering Standard](10-performance-engineering-standard.md) and the [Frontend Architecture Standard](23-frontend-architecture-standard.md). It defines client-side web and mobile performance rules, Core Web Vitals (CWV) metrics, and asset delivery optimization.

---

## 1. The Core Performance Principle

**Do not optimize blindly.**

Before committing to any performance refactoring, answer three questions:
1. *Is this a verified client-side bottleneck?* (Measure first using profiling tools, not assumptions).
2. *Does this improvement directly benefit the user experience?*
3. *Does the complexity of the optimization justify the maintenance overhead?*

**Rule**: Prefer structural, architectural improvements (reducing network request counts, loading smaller packages) over micro-optimization tricks (rewriting simple loop types).

---

## 2. Core Web Vitals (CWV) Standards

Every frontend web application must actively monitor and optimize for Google's Core Web Vitals:

### 2a. Largest Contentful Paint (LCP)
LCP measures perceived page load speed. It marks the point in the page load timeline when the primary content has likely loaded.
- **Target**: **$LCP \le 2.5\text{s}$** under simulated mobile throttle networks.
- **Optimization Strategy**:
  - Optimize the server response time (Time to First Byte / TTFB).
  - Preload critical above-the-fold resources (e.g. hero images) using `<link rel="preload">` or fetch priority headers.
  - Defer all non-critical scripts and heavy bundles using `defer` or `async` tags.
- **Avoid**: Blocking initial paint with heavy third-party tracking scripts or uncompressed styling sheets.

### 2b. Interaction to Next Paint (INP)
INP measures page responsiveness to user inputs (clicks, key presses, taps). It assesses the latency of all interactions.
- **Target**: **$INP \le 200\text{ms}$**.
- **Optimization Strategy**:
  - Break up long-running JavaScript execution blocks (tasks longer than 50ms) using `requestIdleCallback()` or `setTimeout()`.
  - Debounce or throttle high-frequency input event listeners (e.g. window scrolling, search inputs).
  - Offload heavy client-side data parsing to Web Workers.
- **Avoid**: Running synchronous, expensive calculations on the main UI thread during user interaction events.

### 2c. Cumulative Layout Shift (CLS)
CLS measures visual stability by tracking unexpected layout shifts during the page lifecycle.
- **Target**: **$CLS \le 0.1$**.
- **Optimization Strategy**:
  - Always declare explicit `width` and `height` attributes on images, video tags, and custom icon vectors.
  - Reserve structural layout space for dynamically loaded elements (e.g. ads, late-loaded banner components, async alerts) using placeholder skeletons or minimum height CSS rules.
  - Never insert dynamic content directly above existing content unless triggered by direct user interaction.

---

## 3. JavaScript and Bundle Optimization

JavaScript is the most expensive resource because the browser must download, parse, compile, and execute it.

- **No Bloated Dependencies**: Do not import massive libraries for basic utility functions. (e.g., avoid importing all of Lodash for single operations; use native JS utilities).
- **Import Wisely**: Configure your build tools to enforce tree-shaking.
- **Dynamic Imports**: Split bundles into route-based modules. Dynamically import heavy UI dialogs or rarely accessed screens (e.g., imports using `React.lazy()` or Nuxt `<LazyClientOnly>`).
- **Render Thread Offloading**: Avoid calculating complex business math inside render loop callbacks.

---

## 4. Component Render Performance

Components must manage their reactivity footprint to prevent CPU drain.
- **Prevent Unnecessary Re-renders**: Use component memoization (`React.memo`, custom component updates) only when rendering measurements confirm a performance gain.
- **Reactivity Restraint**: Keep reactive variables shallow. Do not wrap large lookup tables or massive static database matrices in active reactive objects (`ref`, `reactive`, or state stores) as it forces deep tracking overhead.
- **Focused Components**: Build smaller, single-responsibility components. Re-rendering a tiny sub-component is significantly cheaper than re-rendering a massive parent container.

---

## 5. Image & Vector Asset Optimization

Images represent the largest download size on typical websites.

- **Explicit Dimensions**: Declare image dimensions on all tags to prevent layout recalculation (CLS).
- **Modern Formats**: Serve images in modern formats (AVIF or WebP) with standard formats (JPEG/PNG) as fallbacks.
- **Compression**: Enforce asset pipeline compression for all production-bound images.
- **Responsive Images**: Use `srcset` and `sizes` attributes to serve matching image scales to mobile viewports.
- **SVG for UI**: Use vector SVG for all icons, brand logos, and geometric backgrounds.

---

## 6. Lazy Loading Policy

- **Lazy Load Below the Fold**: Automatically apply `loading="lazy"` to images and iframes situated below the initial viewport boundary.
- **Lazy Load Components**: Defer mounting heavy components (e.g. charts, interactive tables, video players) until they scroll into view using an intersection observer.
- **Do Not Lazy Load Above the Fold**: Never apply lazy loading to critical hero images, H1 titles, or initial layout headers. Doing so degrades the LCP score.

---

## 7. Data and Network Optimization

Minimize client-server request chatter.
- **No Duplicate API Queries**: Cache resource requests so that navigating between pages does not trigger identical requests.
- **Pagination**: Implement server-side pagination (keyset or offset) for all lists. Never load entire databases or thousands of rows into the client memory.
- **Request Batching**: Group multiple small data requests into a single query payload where supported by the backend API.
- **Incremental Loading**: Load minimal details first, then retrieve secondary metadata asynchronously.

---

## 8. CSS and Styling Performance

- **Minimize Unused CSS**: Clean up deprecated classes. Use purge tools in the build step to strip unused framework styles.
- **Design Tokens**: Route styling variables through CSS Custom Properties, enabling theme adjustments without duplicating stylesheet rules.
- **Layout-Friendly Animations**: Restrict animation transitions to `transform` and `opacity` properties. Animating attributes like `width`, `height`, `margin`, or `left` forces the browser to recalculate full layout positioning on every frame, causing visual lag.

---

## 9. Web Font Performance

- **Weight Minimization**: Limit font files to only the weights and character sets actively rendered in the UI (typically normal 400 and bold 700).
- **Font Display Swap**: Apply `font-display: swap` to all custom web fonts to ensure text is instantly readable in a system fallback font during load.
- **Preload**: Preload primary fonts to minimize layout shifts on initial load.

---

## 10. Third-Party Script Optimization

- **Privacy & Security Audit**: Audit every third-party tracking, chat, or analytics script for performance impact.
- **Self-Hosting**: Host script assets on the application's origin domain to avoid expensive secondary DNS lookups.
- **Defer Load**: Always load non-essential integrations (chat widgets, customer feedback forms) only after the application's primary content has rendered.

---

## 11. Mobile Device Constraints

Do not design solely for high-end developer workstations. Optimize for mobile:
- **CPU Constraints**: Restrict complex loops and DOM depth. Mobile browsers have limited processing power.
- **Memory Limits**: Clean up event listeners, clear unused caches, and dereference unused variables on component unmount.
- **Battery Preservation**: Avoid continuous background polling loops. Use WebSockets or push notifications instead of short polling.

---

## 12. Caching and Offline Resilience

- **Static Cache**: Cache images, CSS, JS, and font files via HTTP Cache-Control headers.
- **Application Cache**: Store non-sensitive API responses in a local store or IndexDB cache to make secondary page transitions instant.
- **Resilience**: Handle network dropouts gracefully. Ensure the UI notifies the user when they are offline without crashing the current state.

---

## 13. Accessibility is Non-Negotiable

**Do not sacrifice accessibility in the name of performance.**
- Never strip semantic HTML tags (`<button>`, `<fieldset>`) to reduce DOM depth.
- Never remove label tags or hidden screen reader texts to save bytes.
- Accessibility and semantics are critical baseline parameters; performance optimizations must work around them.

---

## 14. Framework Performance Guidance

### 14a. Vue / Nuxt
- **Vue**: Use `<KeepAlive>` to cache stateful component views during navigation. Use `<Suspense>` for async component loading.
- **Nuxt**: Maximize server-side rendering (SSR) benefits. Use Nuxt's hydration rules (`<LazyClientOnly>`) to defer hydration of non-interactive layouts.

### 14b. React / Next.js
- **React**: Use `useMemo` and `useCallback` only when props/renders are measured as expensive.
- **Next.js**: Use Server Components for all data-heavy presentation layouts to minimize client-side bundle sizes.

### 14c. Flutter
- **Stateless Widgets**: Declare widgets as `const` where possible to skip redundant widget builds.
- **List Optimization**: Use `ListView.builder` instead of loading all widgets into memory at once.

---

## Performance Review Checklist

Prior to frontend deployment, verify against this performance checklist:
- [ ] **Bundle impact considered**: Are third-party dependencies minimized and tree-shaken?
- [ ] **Images optimized**: Do images use modern formats (AVIF/WebP), have explicit dimensions, and include `srcset`?
- [ ] **Large lists handled correctly**: Are datasets paginated or virtualized (no massive DOM lists)?
- [ ] **No unnecessary API calls**: Are API requests cached to prevent duplicate fetches?
- [ ] **Loading states exist**: Are skeleton screens configured for async content to maintain layout stability?
- [ ] **Mobile performance considered**: Has the screen been profiled under mobile network throttling?
- [ ] **Rendering cost considered**: Are animations restricted to `transform` and `opacity` properties?
- [ ] **Dependencies justified**: Did we avoid importing a full framework library for a minor utility?
- [ ] **Third-party scripts evaluated**: Are trackers deferred until after initial page paint?
- [ ] **Accessibility preserved**: Are semantic elements and screen reader tags fully intact?

---

## References
- Performance Standard: [core/10-performance-engineering-standard.md](10-performance-engineering-standard.md)
- Frontend Architecture: [core/23-frontend-architecture-standard.md](23-frontend-architecture-standard.md)
- Frontend Testing: [core/25-frontend-testing-standard.md](25-frontend-testing-standard.md)
- Chrome Lighthouse: [https://developer.chrome.com/docs/lighthouse/](https://developer.chrome.com/docs/lighthouse/)
- Web Vitals: [https://web.dev/vitals/](https://web.dev/vitals/)
