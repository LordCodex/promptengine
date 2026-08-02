---
document_id: stacks-nuxt-ts-engineering-standard
title: Nuxt 3 Engineering Standard
ecosystem: js-ts-vue-nuxt
target_versions:
  nuxt: "^3.10"
  vue: "^3.4"
  typescript: "^5.0"
dependencies:
  - core-universal-coding-standards
  - core-architecture-and-simplicity
  - core-security-engineering-standard
  - core-performance-engineering-standard
  - stacks-js-ts-conventions
  - stacks-vue-ts-engineering-standard
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Nuxt 3 Engineering Standard

## Purpose & Inheritance
This document defines the core standards for Nuxt 3 full-stack development. It inherits from and extends the [Universal Coding Standards](../../core/05-universal-coding-standards.md), the [Architecture Standards](../../core/02-architecture-and-simplicity.md), the [Security Engineering Standard](../../core/08-security-engineering-standard.md), the [Performance Engineering Standard](../../core/10-performance-engineering-standard.md), and the [Vue 3 and TypeScript Engineering Standard](vue-ts-engineering-standard.md). It establishes patterns for server-side rendering (SSR), hydration safety, Nitro server routers, data fetching, and SEO optimization.

---

## 1. Nuxt Philosophy

Nuxt is not just Vue with an automatic router; it is a **full-stack client-server framework**.
- **SSR Execution Boundary**: Developers must understand that Nuxt code executes twice: first on the Node/Nitro server (SSR phase) and then in the client browser (hydration phase). Code must be resilient to executing in both environments.
- **Auto-Imports Optimization**: Leverage Nuxt's auto-import feature for composables and components, but maintain clear boundaries. Do not duplicate global names, and ensure auto-imported modules have type declarations.
- **Decoupled Backend Engine**: While Nitro can serve full API routes, avoid placing your entire database-backed core business logic inside Nuxt server directories. Use Nuxt server routes as a lightweight BFF (Backend-for-Frontend) layer to call external backend services (e.g., Laravel, Go).

---

## 2. Directory Structure Conventions

Nuxt enforces a file-system directory hierarchy. We structure our projects as follows:

```text
src/
├── assets/       # Uncompiled source assets (Sass, CSS, images)
├── components/   # UI components (automatically imported)
├── composables/  # Stateful Composition API functions (automatically imported)
├── layouts/      # Visual page layouts (e.g., default.vue, auth.vue)
├── middleware/   # Client and server-side navigation route guards
├── pages/        # File-based routing definitions (e.g., index.vue, users/[id].vue)
├── plugins/      # Custom Vue extensions and third-party initializers
├── public/       # Static assets served at the domain root (e.g., favicon.ico)
├── server/       # Nitro server-side layer
│   ├── api/      # Server endpoints accessible via HTTP clients (e.g., /api/pay)
│   ├── middleware/ # Nitro server-side request hooks
│   └── utils/    # Server-only utilities and typings
├── stores/       # Global Pinia state stores
└── utils/        # Stateless helper functions
```

---

## 3. Rendering Strategies Matrix

Nuxt supports multiple rendering strategies. Choose the strategy that matches your feature requirements:

| Rendering Mode | Description | Use Cases | Benefits | Trade-offs |
| :--- | :--- | :--- | :--- | :--- |
| **SSR** (Server-Side) | Node server renders HTML on every request. | E-commerce product lists, public portals, SaaS login walls. | Optimal SEO index speed, fast initial page load. | High server CPU load, complex node environment hosting. |
| **CSR** (Client-Side) | Standard SPA behavior; the browser renders HTML. | Private dashboards, internal tools, profile settings. | Lowest server compute cost; fast subsequent routing. | Poor SEO indexing; slower initial loading speeds. |
| **SSG** (Static Gen) | HTML pages pre-built at compile time. | Marketing landing pages, documentation sites, blogs. | Blazing-fast loading; zero runtime server compute costs. | Rebuild required for data changes. |

### Hybrid Rendering Configuration
Configure different rendering strategies per route using `routeRules` inside `nuxt.config.ts`:
```typescript
// nuxt.config.ts
export default defineNuxtConfig({
  routeRules: {
    '/': { prerender: true },              // SSG at build time
    '/blog/**': { isr: true },             // Incremental Static Gen (cache updates)
    '/admin/**': { ssr: false },           // CSR (Client-Side SPA Only)
    '/api/**': { cors: true }              // CORS API handlers
  }
});
```

---

## 4. Data Fetching Patterns

Do not use raw `fetch()` or Axios inside Nuxt page templates. Enforce Nuxt data-fetching composables to prevent duplicate calls during hydration.

- **`useFetch`**: The default composable for fetching data during page initialization. It prevents duplicate data fetches by sharing state between the server render phase and client hydration.
- **`useAsyncData`**: Use when the fetch operation requires wrapper logic or combining multiple async tasks (e.g., querying local content and an external API simultaneously).
- **`$fetch`**: Use for client-side user event triggers (e.g., form submissions, button clicks) where SSR execution is not needed.

```vue
<!-- Good: Data Fetching with useFetch and Type Safety -->
<script setup lang="ts">
interface Invoice {
  id: string;
  amount_cents: number;
}

// Fetch data safely on SSR and pass to client without duplicate calls
const { data: invoices, pending, error } = await useFetch<Invoice[]>('/api/invoices', {
  lazy: true, // Do not block route navigation during loading
  transform: (data) => data.map(inv => ({ ...inv, formatted: `$${inv.amount_cents / 100}` }))
});
</script>

<template>
  <div>
    <div v-if="pending">Loading invoices...</div>
    <div v-else-if="error">Failed to load invoices.</div>
    <ul v-else>
      <li v-for="inv in invoices" :key="inv.id">
        Invoice #{{ inv.id }}: {{ inv.formatted }}
      </li>
    </ul>
  </div>
</template>
```

---

## 5. Nitro Server Routes

Nitro server routes run inside a sandboxed serverless Node environment.

### Server Route Rules
- **Validate Inputs**: Enforce schema validations on requests using lightweight libraries (e.g., `zod` or `h3` body readers).
- **Enforce Business Boundaries**: Do not execute direct database queries from client-facing page files. Route queries through server-only paths (`/server/api/...`) to protect database credentials.

```typescript
// server/api/invoices.get.ts
import { defineEventHandler, createError } from 'h3';

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig(event);
  
  try {
    // Call the external API secure backend with private credentials
    return await $fetch('/invoices', {
      baseURL: config.apiBaseUrl,
      headers: {
        Authorization: `Bearer ${config.apiToken}`
      }
    });
  } catch (error) {
    throw createError({
      statusCode: 500,
      statusMessage: 'Failed to retrieve invoice records'
    });
  }
});
```

---

## 6. Nuxt with External Backends & Authentication

When consuming external APIs (like Laravel or Go):

### API Configuration
Configure a reverse proxy pattern using Nitro's route rules or fetch wrappers to prevent CORS (Cross-Origin Resource Sharing) blockages during local development.

### Authentication Boundary
- **HttpOnly Cookies**: Do not store JWTs or session tokens in client-accessible storage (like `localStorage` or `sessionStorage`). This exposes tokens to XSS attacks.
- **SSR-Friendly Auth**: Pass authentication cookies through useFetch request headers during SSR execution to ensure the Node server can verify permissions:
  ```typescript
  // Pass client cookies to external API requests during SSR
  const headers = useRequestHeaders(['cookie']);
  const { data } = await useFetch('/api/user', { headers });
  ```

---

## 7. Security Hardening

- **Runtime Configuration Security**: Keep sensitive variables (like API keys, private keys, database passwords) private. Define them in `runtimeConfig.private` in `nuxt.config.ts`. Only variables in `runtimeConfig.public` are exposed to the browser client.
  ```typescript
  // nuxt.config.ts
  export default defineNuxtConfig({
    runtimeConfig: {
      apiToken: '', // Private (only accessible server-side)
      public: {
        apiBaseUrl: '' // Public (accessible client-side)
      }
    }
  });
  ```
- **XSS and CSP Hardening**: Configure Content Security Policies (CSP) to restrict script sources. Ensure all dynamic data renders safely without executing injected scripts.
- **CSRF Protection**: Enable CSRF token checks on all state-mutating requests (POST, PUT, DELETE).

---

## 8. SEO Engineering

SEO is a primary reason to choose Nuxt over vanilla Vue. Implement these meta tag configurations:

- **Strict Head Management**: Use `useSeoMeta` for fast, declarative SEO meta tagging.
- **JSON-LD Structured Data**: Inject structured schema definitions using `useHead` scripts to improve search results indexing.

```vue
<!-- Good: SEO Meta tagging & Schema Integration -->
<script setup lang="ts">
useSeoMeta({
  title: 'Invoice Ledger - Professional Billing Dashboard',
  ogTitle: 'Invoice Ledger - Professional Billing Dashboard',
  description: 'Manage payments and view financial ledgers in real-time.',
  ogDescription: 'Manage payments and view financial ledgers in real-time.',
  ogImage: 'https://domain.com/social-share.png',
  twitterCard: 'summary_large_image'
});

useHead({
  script: [
    {
      type: 'application/ld+json',
      innerHTML: JSON.stringify({
        '@context': 'https://schema.org',
        '@type': 'SoftwareApplication',
        'name': 'Invoice Ledger',
        'applicationCategory': 'BusinessApplication'
      })
    }
  ]
});
</script>
```

---

## 9. Performance & Hydration Safety

Hydration is the process where Vue matches client-side DOM nodes with server-rendered HTML. Hydration mismatches degrade rendering speeds and create layout glitches.

### Hydration Guidelines
- **Avoid Environment-Specific Output**: Do not render date formats or client-specific objects (like `window.innerWidth`) in SSR templates. These values differ on the server and client, causing hydration mismatches. Wrap client-only code in `<ClientOnly>` components:
  ```vue
  <!-- Good: Wrap client-only code in ClientOnly to prevent hydration mismatches -->
  <ClientOnly>
    <div>Window width: {{ windowWidth }}</div>
  </ClientOnly>
  ```
- **Lazy Load Components**: Prefix dynamic components with `Lazy` (e.g., `<LazyInvoiceModal>`) to defer loading their scripts until they are rendered.
- **Image Optimization**: Use the `@nuxt/image` module to serve responsive, optimized image formats (WebP/AVIF) matching target viewport sizes.

---

## 10. Middleware & Plugins

- **Route Middleware**: Best for navigation-level access control (e.g., checking if a user is authenticated before routing to `/dashboard`). Do not execute heavy business calculations inside middleware scripts.
- **Plugins**: Use plugins to initialize global Vue directives or register third-party libraries (e.g., analytical tools) at startup. Keep plugin files lightweight to avoid blocking initial app renders.

---

## 11. TypeScript in Nuxt

Nuxt automatically generates type schemas for your project files, including runtime config variables and API endpoints.

- **Auto-Imports Typings**: Run the Nuxt build compiler (`npx nuxi prepare`) locally to regenerate TypeScript helper declarations (`.nuxt/tsconfig.json` and `.nuxt/nuxt.d.ts`).
- **Response Schemas**: Define strict interfaces for API responses and use them with `useFetch`:
  ```typescript
  const { data } = await useFetch<InvoiceResponse>('/api/invoice');
  ```

---

## 12. Deployment Targets

Nuxt applications can be compiled for multiple execution environments.

### Deployment Configurations
- **Node Server VPS**: Run standard build actions (`npm run build`) and deploy the resulting `.output/server/index.mjs` node script on a VPS server behind an Nginx proxy.
- **Serverless Edge (Vercel, Netlify, Cloudflare Workers)**: Nitro compiles build outputs to match target edge platform runtime signatures automatically.
- **Runtime Variables**: Pass production configurations to Nitro deployments using prefixed environment variables (e.g., `NUXT_API_TOKEN` overrides `runtimeConfig.apiToken`).

---

## 13. Decision Matrices

Use these matrices to identify the correct Nuxt engineering decision based on project context.

### Matrix 1: Nuxt vs. Vue (Vanilla SPA)
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Public index portals, e-commerce, content-focused sites | **Nuxt 3** | Built-in Server-Side Rendering (SSR) and SEO optimizations. |
| Private internal dashboard clients, tools with complex state | **Vanilla Vue 3** | Avoids server-side hydration issues; lower server costs. |

### Matrix 2: SSR vs. CSR (Client-Side SPA Only)
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Public landing pages, blog grids, directories | **SSR** | Allows fast indexing and provides optimal search metadata rendering. |
| Authenticated pages, settings forms, user profile editors | **CSR (Client-Side)** | Zero server rendering cost; pages load instantly after login. |

### Matrix 3: useFetch vs. useAsyncData vs. $fetch
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Standard page updates, simple queries | **useFetch** | Default, clean wrapper; handles SSR state sharing. |
| Bundled requests, caching wrappers | **useAsyncData**| Provides flexibility when wrapping multiple operations. |
| Button triggers, user input mutations | **$fetch** | Prevents SSR overhead; runs strictly client-side. |

### Matrix 4: Nitro Server Route vs. External API Backend (Laravel / Go)
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Simple data proxies, credential masking, token refreshes | **Nitro Server Route** | Low latency; hides sensitive credentials from client requests. |
| Heavy database logic, queues processing, transactional tables | **External API Backend**| Decouples compute logic from rendering layers; scales independently. |

### Matrix 5: Composable vs. Plugin
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| UI state handlers, forms data validations | **Composable** | Automatically imported; lifecycle-aware; scope-isolated. |
| Global integrations, global utility definitions | **Plugin** | Boots once at startup; extends Vue constructor objects. |

### Matrix 6: Pinia vs. Server State Caches (e.g., Nuxt useFetch Caching)
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Interactive client states, user shopping carts | **Pinia** | Persistent client storage with unified mutations. |
| Static list payloads, API get queries | **Server State** | Defer states caching to useFetch key configurations. |

---

## 14. AI Nuxt Rules

AI agents modifying or writing Nuxt code in this repository must follow these rules:

1. **Verify Hydration Safety**: Do not write client-side variables (such as `localStorage` or `window`) inside global components without wrapping them in `<ClientOnly>` tags.
2. **Never Expose Sensitive Keys**: Ensure all API tokens and secret config keys are defined in private runtime environments.
3. **Use useFetch for Initial Page Loads**: Do not suggest using raw Axios or browser `fetch()` for initial page load requests.
4. **No Direct DB Queries in Client Pages**: Restrict database calls to the `/server` directory; client templates must consume data via API endpoints.
5. **No Blind Module Additions**: Do not add third-party modules to `nuxt.config.ts` without validating existing utilities first.

---

## 15. Nuxt Review Checklist

Use this checklist during code review to evaluate Nuxt applications.

### Architecture & Directory Alignment
- [ ] Are files located in their correct directories (e.g., pages in `/pages`, endpoints in `/server`)?
- [ ] Is business logic separated from client UI templates?

### Rendering & Hydration
- [ ] Are client-only API parameters (e.g., innerWidth) wrapped in `<ClientOnly>` tags?
- [ ] Have hydration mismatches been verified and resolved?

### Data Fetching & Caching
- [ ] Do page initialization calls use `useFetch` or `useAsyncData` (no raw Axios)?
- [ ] Are user interaction triggers using `$fetch` (no SSR overhead)?

### SEO & Metadata
- [ ] Are SEO tags configured using `useSeoMeta`?
- [ ] Is JSON-LD structured schema code injected on public pages?

### Security Hardening
- [ ] Are API secrets kept out of public runtime configurations?
- [ ] Are authentication tokens stored in secure, HttpOnly cookies?

### Performance & Build
- [ ] Are dynamic UI dialog components lazy-loaded?
- [ ] Do images use responsive WebP/AVIF formats?

---

## References
- Universal Naming Rules: [core/05-universal-coding-standards.md](../../core/05-universal-coding-standards.md)
- Security Engineering: [core/08-security-engineering-standard.md](../../core/08-security-engineering-standard.md)
- Vue Component Conventions: [vue-ts-engineering-standard.md](vue-ts-engineering-standard.md)
