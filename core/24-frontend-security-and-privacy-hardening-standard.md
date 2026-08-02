---
document_id: core-frontend-security
title: Frontend Security and Privacy Hardening Standard
ecosystem: cross-cutting
dependencies:
  - core-security-engineering-standard
  - core-frontend-architecture
  - core-universal-coding-standards
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Frontend Security and Privacy Hardening Standard

## Inheritance
This document inherits from and extends the [Security Engineering Standard](08-security-engineering-standard.md) and the [Frontend Architecture Standard](23-frontend-architecture-standard.md). It defines client-side browser and device security enforcement policies, threat mitigations, and privacy-preservation guidelines.

---

## 1. The Core Frontend Security Principle

**The frontend is a completely untrusted environment.**

Everything compiled into, executed within, or stored inside the client browser or device is public and modifiable by the client. Never assume:
- A hidden UI component prevents access (e.g. hiding an admin tab based on a local flag).
- A disabled button prevents execution.
- Client-side validation protects the database.
- A client-side API route guard enforces authorization.
- Users cannot capture, inspect, and replay network requests.

**Rule**: All authorization, verification, sanitization, and data validations must be executed and enforced on the server.

---

## 2. Cross-Site Scripting (XSS) Prevention

Do not insert user-controlled content as executable HTML.

### Prohibited Unescaped Sinks
Unless data is strictly sanitized via an approved HTML sanitizer (e.g. DOMPurify) and explicitly required by the product design, you must never use the following sinks:

- **Vanilla CSS/JS**: `Element.innerHTML`, `Element.outerHTML`, `document.write()`.
- **React**: `dangerouslySetInnerHTML`.
- **Vue**: `v-html` directive.
- **Flutter**: Render HTML extensions or raw webview execution without origin restrictions.

### Safe Alternatives
- Use text rendering nodes: `Element.textContent`, `Element.innerText`.
- Use standard framework escaping (e.g. `{{ expression }}` in Vue, `{expression}` in React).
- Use native browser DOM APIs like `document.createElement()` and `Node.appendChild()`.

---

## 3. Content Security Policy (CSP)

A Content Security Policy must be defined and enforced on all web servers.

### Recommended Baseline Header Configuration
```text
default-src 'self';
script-src 'self';
style-src 'self';
img-src 'self' data:;
font-src 'self';
connect-src 'self';
object-src 'none';
base-uri 'self';
frame-ancestors 'self';
form-action 'self';
```

### CSP Constraints
- **Never allow `unsafe-inline` or `unsafe-eval`** in `script-src` or `style-src` unless an explicit, documented architectural exception is approved.
- Use cryptographic nonces (`nonce-[value]`) or SHA-256 hashes (`sha256-[value]`) for approved inline scripts or styles that cannot be moved to external files.
- Constrain `frame-ancestors` to `'self'` or explicit approved origins to prevent clickjacking.

---

## 4. Security Headers

Ensure the server deployment config enforces the following headers for all frontend responses:

| Header | Value | Purpose |
| :--- | :--- | :--- |
| **Content-Security-Policy** | (See baseline above) | Restricts resource loading vectors. |
| **Strict-Transport-Security** (HSTS) | `max-age=63072000; includeSubDomains; preload` | Forces HTTPS communication. |
| **X-Content-Type-Options** | `nosniff` | Prevents MIME-type sniffing attacks. |
| **X-Frame-Options** | `DENY` or `SAMEORIGIN` | Mitigates clickjacking. |
| **Referrer-Policy** | `strict-origin-when-cross-origin` | Protects referrer leakage to third-parties. |
| **Permissions-Policy** | `geolocation=(), camera=(), microphone=()` | Restricts access to sensitive browser APIs. |

### Header Exposure Rule
Do not expose headers that leak tech stack internals, versions, or build dates (e.g., `X-Powered-By`, `Server` details).

---

## 5. Cross-Site Request Forgery (CSRF)

When using cookie-based authentication:
- **CSRF Tokens**: Force the server to generate and validate unique anti-CSRF tokens for all state-changing requests (POST, PUT, DELETE, PATCH).
- **SameSite Attribute**: Configure all session cookies with the `SameSite=Lax` or `SameSite=Strict` attributes.
- **Secure Attribute**: Mark all session cookies as `Secure` (HTTPS only) and `HttpOnly` (inaccessible to JS).
- **Origin Validation**: Verify the `Origin` and `Referer` headers on all state-changing requests.

---

## 6. Authentication Security

Sensitive authorization tokens must not be exposed to XSS data exfiltration.

### Storage Rules
- **Do not store** access tokens, refresh tokens, passwords, or recovery codes in **`localStorage`** or **`sessionStorage`**. These storages are globally accessible to any script running on the origin (XSS vector).
- **Prefer HttpOnly Cookies**: Store authentication tokens in `Secure`, `HttpOnly`, `SameSite=Lax` cookies, preventing JavaScript access entirely.
- **Memory Storage Fallback**: If cookies are not supported, store tokens in-memory (e.g., local state variables) and handle session recovery via a secure refresh token endpoint.

## 7. Sensitive Data Handling

### 7a. Secret and API Key Protection
Frontend code must never contain or expose:
- Private or secret API keys.
- Database credentials or connection URIs.
- Payment provider secret keys.
- Admin or superuser authentication tokens.
- Service account credential JSONs or files.
- Symmetric or private encryption keys.
- Internal infrastructure or server credentials.

Never assume environment variables are secret. Any variable exposed to frontend build tools (e.g. prefix-exposed) is compiled into the client-side JavaScript assets and becomes publicly visible.

**Unsafe Environment Variables (Never compile these):**
- `VITE_SECRET_KEY=value`
- `NEXT_PUBLIC_PRIVATE_KEY=value`
- `NUXT_PUBLIC_SECRET=value`

Do not place secrets in:
- `.env` files that are loaded or bundled into frontend builds.
- Public configuration or layout files.
- JavaScript/TypeScript constants or environment wrappers.
- Client-side configuration settings or JSON files.

### 7b. Public vs. Private Configuration
Only expose values that are intentionally public.

| Configuration Type | Safe to Expose | Unsafe (Keep Server-Side) |
| :--- | :--- | :--- |
| **System Identity** | Application name, version, public support link | Server internals, source paths, server operating system |
| **APIs and Endpoints** | Public API base URL, public analytics keys | Database connection strings, database URLs, microservice auth tokens |
| **Third-Party Keys** | Public stripe keys, user-facing client IDs | Payment secret keys, private client secret keys, webhooks secrets |
| **Feature Control** | User-facing feature flags, active theme | Admin flags, billing override flags |

### 7c. API Communication & Backend Proxy Rule
If the frontend needs access to a protected external service:
- Do not call the service directly from the browser with secret credentials.
- All sensitive integrations must go through a backend proxy endpoint that appends credentials securely server-side.

```text
[ Browser / Frontend Client ]
             ↓
[ Backend API Proxy (appends secret key) ]
             ↓
[ External Third-Party Service ]
```

**Required Proxy Implementations:**
- **Payments**: Frontend triggers the backend to create a secure checkout session or payment intent. Browser does not call Stripe with the secret API key.
- **AI Integrations**: Frontend requests inferences from the application backend. Browser does not pass the OpenAI or Gemini API secret directly.
- **Email Services**: Frontend triggers transaction actions via backend controllers. Browser does not interface with SendGrid or Mailgun API keys.

### 7d. Response Data Security
Backend APIs must be designed to restrict the data returned to the frontend. Never return:
- Password hashes or reset tokens.
- Internal configurations, connection secrets, or private keys.
- Unnecessary personal information (PII) that the active view does not render.
- Administrative permission structures or utility fields.

**Rule**: Only return the absolute minimum database fields required by the active UI template.

### 7e. Frontend Build Audit
Before deploying any frontend build to production:
- **Exposed Secret Audit**: Search compiled JavaScript bundles (`/dist`, `/build`, `.next`) for accidentally exposed key prefixes or strings.
- **Environment Variable Audit**: Run a sanity check on all environment variables active in the build runner.
- **Source Map Verification**: Confirm production source maps do not expose sensitive files or code blocks to the public web inspector unless explicitly allowed by team architecture.

---

## 8. URL and Protocol Security

User-controlled URLs and redirects are common entry points for attacks.

### Open Redirect Prevention
Do not redirect users to URLs provided in query parameters (e.g., `?next=http://malicious.com`) without strictly validating the destination origin.

### Protocol Whitelisting
When rendering dynamic links (e.g. `<a :href="userLink">`):
- Only accept explicitly whitelisted protocols: `https://`, `http://`, `mailto:`, `tel:`.
- **Always reject** execution schemes: `javascript:`, `data:`, `file:`.

```typescript
// Correct: whitelist protocol validation
function sanitizeUrl(inputUrl: string): string {
  const normalized = inputUrl.trim().toLowerCase();
  if (normalized.startsWith('https://') || normalized.startsWith('http://')) {
    return inputUrl;
  }
  return 'about:blank';
}
```

---

## 9. File Upload Security (Client-Side)

Client-side upload validations are strictly for user experience, not security:
- **File constraints**: Validate file size and select allowed MIME types/extensions on the client to give fast feedback.
- **File previews**: Render image previews using object URLs (`URL.createObjectURL()`) to prevent executing raw script paths.
- **Backend Rules**: Backend must independently re-validate MIME types, file sizes, and rename files upon storage.

---

## 10. Third-Party Scripts & Dependency Auditing

### Script Evaluation
Before loading any external script (e.g., Google Analytics, Intercom, Facebook Pixel):
- Evaluate the security, privacy, and performance impact.
- Avoid loading scripts from third-party CDNs. Prefer hosting scripts locally within the application bundle.
- Apply Subresource Integrity (SRI) hashes to all external scripts to guarantee they have not been tampered with:

```html
<script
  src="https://cdn.example.com/library-1.0.0.js"
  integrity="sha384-oqVuAfXRKap7fdgcCY5uykM6+R9GqQ8K/uxy9rx7HNQlGYl1kPzQho1wx4JwY8wC"
  crossorigin="anonymous"
></script>
```

### Dependency Hardening
- Regularly run dependency audits: `npm audit`, `pnpm audit`, or `yarn audit`.
- Automate security scans in CI pipelines; fail builds on serious or critical vulnerability findings.

---

## 11. DOM Security & Execution Risks

Avoid unsafe JavaScript patterns:
- **Do not use `eval()`** or the `Function` constructor.
- Do not build HTML strings manually by concatenating user strings (`const html = '<div>' + input + '</div>'`).
- Do not use dynamic script insertion (`document.createElement('script')` with user inputs).

---

## 12. Browser Storage Rules

Ask three questions before writing data to `localStorage` or `sessionStorage`:
1. *Is this data sensitive?* (If yes, do not store)
2. *Does it need persistence across browser restarts?* (If no, use memory state)
3. *Can the server manage this state instead?* (If yes, keep it server-side)

**Prohibited Client Storage**: Passwords, access/refresh tokens, credit card details, government IDs, and highly sensitive personally identifiable information (PII).

---

## 13. Privacy Engineering and Data Minimization

- **Data Minimization**: Collect, store, and transmit only the absolute minimum data required to satisfy the feature.
- **No Hidden Analytics**: Do not track user behaviors, inputs, or selections without explicit consent (cookie banner, privacy policy).
- **Masking Inputs**: Use correct input type masks (`type="password"`) and disable input autofill for sensitive fields where appropriate.

---

## 14. Payment Security (PCI Compliance)

Frontend applications must never touch or handle raw credit card numbers:
- Use payment gateways (Stripe, Braintree, PayPal) that render hosted input fields inside isolated iframe elements (e.g., Stripe Elements).
- Never trust product prices or payment amounts displayed in client-side HTML. All payment totals must be recalculated and validated server-side based on database records before creating the payment intent.

---

## 15. Performance Security

Security measures must not cripple frontend performance:
- Avoid loading heavy, bloated client-side encryption libraries inside the rendering path.
- Keep dependency audits fast and lightweight. Do not block initial render or hydration with security checks.

---

## Review Checklist

Prior to frontend deployment, verify against this security checklist:
- [ ] **No unsafe HTML rendering**: Dynamic user content is not rendered via `innerHTML`, `v-html`, or `dangerouslySetInnerHTML` without DOMPurify sanitization.
- [ ] **No exposed secrets**: All API keys, passwords, database credentials, and client secrets are absent from build assets.
- [ ] **CSP considered**: Content Security Policy is configured on the host server; no `unsafe-inline` or `unsafe-eval` without justification.
- [ ] **Security headers configured**: HSTS, X-Frame-Options, X-Content-Type-Options, and Referrer-Policy are present.
- [ ] **Authentication handled safely**: Tokens are not stored in localStorage/sessionStorage. Secure, HttpOnly cookies are used.
- [ ] **CSRF protection considered**: CSRF tokens or SameSite attributes protect state-changing endpoints.
- [ ] **No unsafe browser APIs**: Visually hidden inputs, geolocation, and clipboard access are minimized.
- [ ] **Dependencies reviewed**: Dependency audit check passed with zero serious/critical warnings.
- [ ] **User content escaped**: Dynamic attributes and text bindings are escaped correctly.
- [ ] **Error safety**: Stack traces and raw database errors are suppressed from user-facing screens.
- [ ] **Privacy impact considered**: Data collection follows minimization principles; consent is obtained where required.

---

## References
- Security Standard: [core/08-security-engineering-standard.md](08-security-engineering-standard.md)
- Frontend Architecture: [core/23-frontend-architecture-standard.md](23-frontend-architecture-standard.md)
- DOMPurify: [https://github.com/cure53/DOMPurify](https://github.com/cure53/DOMPurify)
- OWASP Top 10 client-side: [https://owasp.org/www-project-top-ten/](https://owasp.org/www-project-top-ten/)
