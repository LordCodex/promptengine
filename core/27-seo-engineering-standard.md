---
document_id: core-seo-engineering
title: SEO Engineering and Search Optimization Standard
ecosystem: cross-cutting
dependencies:
  - core-frontend-architecture
  - core-frontend-performance
  - core-frontend-security
audience: [human, agent]
last_reviewed: 2026-08-01
---

# SEO Engineering and Search Optimization Standard

## Purpose & Inheritance
This document defines client-side search engine optimization (SEO) standards, indexing policies, and content metadata structures. It inherits from the [Frontend Architecture Standard](23-frontend-architecture-standard.md), the [Frontend Security Standard](24-frontend-security-and-privacy-hardening-standard.md), and the [Frontend Performance Standard](26-frontend-performance-engineering-standard.md), ensuring applications are discoverable, indexable, fast, and structured for search crawler parsing.

---

## 1. The Core SEO Principle

**SEO is not an afterthought.**

Search engines index and evaluate applications using technical metrics: semantic DOM clarity, rendering speeds, accessibility tags, layout stability, and structured schema graphs.
- **Do not rely on hacks** (e.g. dynamic crawler redirects, hidden text).
- **Rule**: A crawler-friendly application is almost always a better application for human users. Optimization decisions must serve both.

---

## 2. Semantic Document Structure

Search crawlers parse semantic relationships to identify primary page topics:
- **Layout Landmarkers**: Use `<main>`, `<article>`, `<section>`, `<nav>`, `<header>`, and `<footer>` to divide layout segments. Avoid using nested generic `<div>` tags for page layouts.
- **Section Headings**:
  - Every indexable page must contain **exactly one `<h1>`** tag representing the main topic.
  - Nest subheadings (`<h2>`, `<h3>`, `<h4>`) sequentially.
  - **Never use headings for font styling alone.** Use CSS utility classes to style typography sizes without altering semantic document hierarchy.

---

## 3. Metadata and Indexing Control

### 3a. Page Title & Meta Description Rules

Every indexable route must supply a unique page title and meta description:

| Attribute | Constraint | Recommended Value Pattern |
| :--- | :--- | :--- |
| **Page Title** | 50–60 characters; unique per page | `[Page Topic] | [Secondary Context] - [Brand Name]` |
| **Meta Description** | 150–160 characters; descriptive summary | A human-readable call to action matching user search intent. |

- Avoid generic page titles (`"Home"`, `"Details"`, `"Page 3"`).
- Avoid keyword stuffing or duplicate metadata across different page lists.

### 3b. Indexing Boundaries & Control
Prevent crawlers from indexing private, duplicate, or administrative pages:
- **Meta Robots Tag**: Apply `<meta name="robots" content="noindex, nofollow">` to private screens (user accounts, billing, settings) and intermediate transition pages (payment success/auth gates).
- **Authentication**: Gate private routes behind server-side authentication protocols. Do not rely solely on `robots.txt` rules to hide private URLs.
- **Canonical URLs**: Every indexable page must declare a self-referential canonical URL (`<link rel="canonical" href="https://example.com/page">`) to prevent crawler fragmentation from trailing slashes or UTM parameter URLs.

---

## 4. URL Routing Architecture

Configure route URLs to be human-readable, stable, and semantic:
- **Clean Slugs**: Prefer `/books/clean-code` over `/page?id=239482` or `/books/id/239482`.
- **Lowercase**: Enforce lowercase routes.
- **URL Stability**: Never change existing routes without setting up HTTP `301 Permanent Redirect` rules.

---

## 5. Crawl Budget Utilities (Sitemaps & Robots)

### Sitemaps
- Generate dynamic XML sitemaps automatically during build pipelines or server routines.
- Include only public, canonical, and indexable URLs.
- Exclude all authenticated dashboards, admin layouts, and duplicate routes.
- Split sitemaps exceeding 50,000 URLs or 50MB files using a Sitemap Index.

### Robots.txt
Maintain a clean root `/robots.txt` configuration:
- Explicitly declare the URL location of your primary sitemap file.
- Disallow crawlers from accessing private directories (e.g. `/admin/`, `/api/`, `/user/`).
- **Never disallow asset paths** (CSS, JS, fonts) required by modern search engines to render the application interface.

---

## 6. Schema.org Structured Data

Provide search engines with structured context using JSON-LD metadata schema.

### Example Schema Implementation (JSON-LD Product)
```html
<script type="application/ld+json">
{
  "@context": "https://schema.org",
  "@type": "Product",
  "name": "Clean Code Handbook",
  "image": "https://example.com/images/clean-code.jpg",
  "description": "A handbook of agile software craftsmanship.",
  "sku": "978-0132350884",
  "offers": {
    "@type": "Offer",
    "url": "https://example.com/books/clean-code",
    "priceCurrency": "USD",
    "price": "39.99",
    "availability": "https://schema.org/InStock"
  }
}
</script>
```

- **Rules**:
  - Implement dynamic schema matching the page topic (Product, Book, Organization, Article, FAQ).
  - Only mark up values that are visible to the user on the screen. Fake or hidden structured data triggers search penalties.

---

## 7. Social Metadata (Rich Cards)

Ensure all public content pages contain Open Graph (OG) and Twitter card parameters to render rich cards when shared:

```html
<!-- Open Graph -->
<meta property="og:title" content="Affordable Digital Books Marketplace | Example">
<meta property="og:description" content="Browse and purchase agile programming books and design patterns.">
<meta property="og:image" content="https://example.com/assets/og-home.jpg">
<meta property="og:url" content="https://example.com/">
<meta property="og:type" content="website">

<!-- Twitter Cards -->
<meta name="twitter:card" content="summary_large_image">
<meta name="twitter:title" content="Affordable Digital Books Marketplace | Example">
<meta name="twitter:description" content="Browse and purchase agile programming books and design patterns.">
<meta name="twitter:image" content="https://example.com/assets/og-home.jpg">
```

---

## 8. Images and Asset Search Optimization

- **Descriptive Names**: Never use camera-generated filenames (e.g. `IMG_0482.jpg`). Use semantic hyphenated names (e.g. `clean-code-book-cover.jpg`).
- **Alt Attribute**: Enforce `alt` descriptions on all assets. Use empty values (`alt=""`) only for decorative assets.
- **Image Performance**: Apply responsive images (`srcset`) to ensure Google Mobile-First Indexing checks pass.

---

## 9. JavaScript SEO and Rendering Architecture

Crawlers render pages dynamically, but execution delays consume "crawl budget."

### Rendering Strategies
Choose the correct rendering model based on the indexability requirement:

| Rendering Model | Best For | SEO Impact |
| :--- | :--- | :--- |
| **Static Site Generation** (SSG) | Marketing pages, documentation, static product catalogs | Excellent (instant HTML parsing, low crawl budget impact) |
| **Server-Side Rendering** (SSR) | Dynamic public content, e-commerce stores, news portals | Excellent (fresh data generated server-side, instant crawler access) |
| **Single Page App** (SPA) | Authenticated dashboards, internal business tools | Poor (crawlers must execute JavaScript; prone to blank indexing) |

**Rule**: Publicly discoverable content must be built using SSR or SSG architectures. Avoid using SPAs for SEO-critical pages.

---

## 10. Framework Metadata Implementation

### 10a. Vue 3 / Nuxt
Use Nuxt's head metadata composables to inject SEO attributes cleanly:

```typescript
// Correct: Nuxt head configuration
useHead({
  title: 'Affordable Digital Books Marketplace',
  meta: [
    { name: 'description', content: 'Browse and purchase agile programming books.' }
  ],
  link: [
    { rel: 'canonical', href: 'https://example.com/books' }
  ]
})
```

### 10b. React / Next.js
Use Next.js Page Router metadata objects or App Router metadata configs:

```typescript
// Correct: Next.js App Router metadata config
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Affordable Digital Books Marketplace',
  description: 'Browse and purchase agile programming books.',
  alternates: {
    canonical: 'https://example.com/books',
  }
}
```

---

## SEO Review Checklist

Before completing public-facing frontend work, verify against this SEO checklist:
- [ ] **Unique page title**: Title is specific, contains relevant terms, and stays under 60 characters.
- [ ] **Meta description exists**: Description is unique, summarizes the page, and stays under 160 characters.
- [ ] **Correct heading hierarchy**: Exactly one `<h1>` is present; subheadings nested sequentially.
- [ ] **Semantic HTML used**: Structural elements (`<main>`, `<article>`, `<nav>`) are used in place of divs.
- [ ] **URL structure reviewed**: URL routes are lowercase, human-readable, and stable.
- [ ] **Canonical handling considered**: Canonical tags are configured to prevent duplicate indexing.
- [ ] **Sitemap considered**: The route is integrated into the sitemap logic if the page is public.
- [ ] **Robots rules reviewed**: Admin and user account screens include a `noindex` robots meta tag.
- [ ] **Structured data accurate**: Valid Schema.org JSON-LD matches visible on-screen elements.
- [ ] **Social metadata added**: Open Graph and Twitter Card tags are populated with title, description, and image.
- [ ] **Images optimized**: File names are descriptive, alt tags are present, and sizes are compressed.
- [ ] **Performance considered**: Core Web Vitals targets (LCP, INP, CLS) are met to preserve crawler score.
- [ ] **Accessibility considered**: Semantic markup, keyboard access, and input labels are fully functional.

---

## References
- Frontend Architecture: [core/23-frontend-architecture-standard.md](core/23-frontend-architecture-standard.md)
- Frontend Performance: [core/26-frontend-performance-engineering-standard.md](core/26-frontend-performance-engineering-standard.md)
- Schema.org markup references: [https://schema.org](https://schema.org)
- Google Search Console guidelines: [https://search.google.com/search-console/about](https://search.google.com/search-console/about)
