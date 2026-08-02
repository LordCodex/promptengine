---
document_id: core-mobile-application-engineering-standard
title: Mobile Application Engineering Standard
ecosystem: cross-cutting
dependencies:
  - core-universal-coding-standards
  - core-architecture-and-simplicity
  - core-database-engineering-standard
  - core-api-engineering-standard
  - core-security-engineering-standard
  - core-testing-engineering-standard
  - stacks-flutter-dart-engineering-standard
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Mobile Application Engineering Standard

## Purpose & Inheritance
This document defines the core standards for mobile client application engineering across Flutter, native Android, native iOS, and hybrid frameworks. It inherits from and extends the [Universal Coding Standards](05-universal-coding-standards.md), the [Architecture Standards](02-architecture-and-simplicity.md), the [Security Engineering Standard](08-security-engineering-standard.md), the [Testing Engineering Standard](11-testing-engineering-standard.md), and the [Flutter and Dart Engineering Standard](../stacks/dart-flutter/flutter-dart-engineering-standard.md). It establishes implementation protocols for offline synchronization, local secure cache management, push notification deep links, and client-side transaction verification safety.

---

## 1. Mobile Engineering Philosophy

Unlike web clients running in stable browser tabs, mobile applications operate in highly unpredictable client environments:
- **Network Instability**: Devices transition between fast Wi-Fi, high-latency 3G, and complete signal loss (tunnels, elevators).
- **Physical Constraints**: Hardware resources are shared and finite. High memory consumption triggers immediate OS termination, and excessive CPU loops degrade battery life.
- **Outdated Clients**: Users may refuse to update their application version for months or years. The backend API must support deprecated API endpoints.
- **Implicit Termination**: The mobile OS can suspend or terminate background applications at any moment to reclaim resources. State restoration must be handled transparently.

---

## 2. Mobile Application Architecture

We enforce a decoupled, layered architecture to prevent visual layouts from becoming tightly coupled to business calculations.

```text
  [Presentation Layer] (Widgets, UI Controllers, Forms)
         │
  [Business Logic Layer] (Domain Actions, State Managers)
         │
  [Data Layer] (Repositories Interfaces)
    ├── [Network Client] (Dio/Http Session Manager) ──> Backend API
    └── [Local Storage] (SQLite/Hive Cache DB) ──> Device Storage
```

### Layer Responsibilities
- **Presentation**: Renders layouts based on reactive state models. Handles user inputs (clicks, gestures) and delegates operations to the business logic layer.
- **Business Logic (Domain)**: Enforces validation rules, state transitions, and business calculations.
- **Data (Repository)**: Decides whether to return data from local database caches or fetch fresh payloads from remote network clients.

---

## 3. Mobile Lifecycle Management

Every mobile application transitions through key OS state cycles.

```text
[Inactive / Background] <──> [Active (Foreground)] ──> [Suspended / Terminated]
```

### State Management Directives
- **Launch**: Preload lightweight configurations. Defer heavy initializations (like telemetry analytics) until after the initial UI frame renders to keep app boot time low.
- **Background**: Halt UI animations, pause video playback, and cancel non-critical network requests to conserve battery.
- **State Restoration**: Save the current user context (e.g. current navigation target and active form inputs) to temporary storage when backgrounded. If the OS terminates the application, restore this state on the next boot.
- **Network Recovery**: Listen to system network transitions. Automatically resume queued mutations when the device reconnects to a network interface.

---

## 4. Offline-First Engineering

An offline-first strategy treats local database caches as the **primary source of truth for UI presentation**, updating network records asynchronously in the background.

### Sync Configurations
- **Offline-First**: Recommended for apps requiring high operational readiness (e.g., field notes, messaging apps). Read and write operations execute locally first.
- **Online-First**: Best for transactional resources requiring real-time verification (e.g., ticket purchases, banking transfers). Deny operations if a network connection is unavailable.
- **Local Storage Scoping**:
  - *Save locally*: Session layouts, feature configuration settings, offline draft models.
  - *Never save locally*: Unencrypted password strings, raw database encryption keys, full payment cards PIN values.

---

## 5. Network Resiliency & API Design

Mobile network clients must expect connections to fail.

### Resiliency Protocol
- **Dio Client Timeout Boundaries**: Set connect timeouts to $10\text{ seconds}$ and read/write timeouts to $15\text{ seconds}$.
- **Exponential Backoff**: When retrying failed network requests, implement a randomized exponential backoff formula to prevent flooding the backend server:
  $$T_{\text{wait}} = 2^{\text{retry\_count}} + \text{random\_jitter}$$
- **Request Cancellation**: Cancel active network requests when a user navigates away from a page to save bandwidth.

### Mobile-Optimized API Design
- **Minimize Payloads**: Strip unused properties from JSON API response payloads to reduce download size.
- **Cursor Pagination**: Enforce cursor-based pagination (using token timestamps) instead of offset limits to prevent duplicate records rendering when lists update.
- **Backward Compatibility**: Do not deprecate API fields without incrementing route versions (e.g. `/api/v2/...`). Expect older clients to query v1 endpoints indefinitely.

---

## 6. Synchronization & Conflicts Resolution

When offline mutations are pushed to the server, merge conflicts are inevitable.

### Resolution Policies
- **Last-Write-Wins (LWW)**: The server overwrites existing database attributes with the latest client payload based on timestamps. Best for non-collaborative data structures.
- **Merge Strategy**: For collaborative data structures, implement field-level merging or store modifications in conflict tables for manual user resolution.
- **Idempotency Keys**: Append a unique UUID payload key (`Idempotency-Key`) to all state-mutating requests (POST/PUT). The backend API must cache this key to prevent duplicate records from being created if the client retries a request.

---

## 7. Push Notifications Engineering

```text
Backend Server ──> Notification Gateway (FCM / APNs) ──> Target Mobile OS Device
```

### Notification Guidelines
- **Permission Requests**: Do not request push notification permissions on the initial application launch screen. Wait until the user performs an action that benefits from notifications (e.g. placing an order).
- **Deep Linking Safety**: Verify that deep links parsed from notification payloads validate user permissions before routing them directly to private screens.
- **Rate Limit Notifications**: Never spam notifications. Group alerts on the backend before dispatching them to the gateway.

---

## 8. Mobile Authentication & Local Security

- **Secure Session Token Refresh**: Store JWT session tokens in secure Keychains (iOS) or KeyStore (Android) containers. Use refresh tokens to fetch new access tokens in the background:

```dart
// Good: Secure token lookup and rotation flow
Future<String?> getValidAccessToken() async {
  final token = await secureStorage.read(key: 'access_token');
  if (isTokenExpired(token)) {
    final refreshToken = await secureStorage.read(key: 'refresh_token');
    return await rotateToken(refreshToken);
  }
  return token;
}
```

- **Obfuscation & Obfuscate Flags**: Run native binary obfuscations during compilation to hide API routes, string keys, and class names from decompilation tools (e.g., ProGuard for Android, symbol stripping for iOS).
- **SSL Pinning**: Pin server certificates inside the mobile HTTP client to block Man-in-the-Middle (MitM) traffic interceptions.

---

## 9. Payments & Financial Transactions

Mobile clients are highly susceptible to instrumentation and hacking. **You must never trust transaction validation checks performed by the mobile application.**

### Financial Rules
- **Server-Side Validation Only**: The mobile client simply initiates the payment flow and receives a transaction receipt signature. The server must verify this receipt signature directly with the payment provider's API (e.g. Apple App Store, Google Play Billing) before unlocking access.
- **Receipt Replay Guards**: Cache validated receipt IDs in the database to prevent attackers from using the same transaction receipt twice (receipt replay attacks).
- **Idempotency checks**: Generate and verify unique transaction hashes before processing payments to prevent duplicate charges.

---

## 10. Mobile Performance

- **Memory Management**: Avoid loading high-resolution images in memory. Set `cacheWidth` and `cacheHeight` on image widgets to downscale assets to match their target rendering sizes.
- **Identify Jank (Frame Drops)**: Avoid executing heavy CPU tasks (like parsing large JSON responses or compressing images) on the UI thread. Use background threads or Isolates.
- **Startup Latency**: Keep initial app loading times under $2\text{ seconds}$. Defer third-party SDK initializations to keep startup fast.

---

## 11. Accessibility (A11y)

- **Accessible Touch Targets**: Interactive elements (buttons, inputs, checkboxes) must have a minimum touch target size of $48 \times 48\text{ dp}$ (density-independent pixels) to allow easy tapping.
- **Font Scaling**: Use scalable text units (like `sp` in Android or dynamic types in iOS) rather than fixed pixel heights to ensure text wraps correctly when users increase their system font size.
- **Contrast Check**: Maintain a minimum contrast ratio of $4.5:1$ for normal text against its background.

---

## 12. Decision Matrices

Use these matrices to identify the correct mobile application engineering decision based on project context.

### Matrix 1: Offline-First vs. Online-First
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Core features require constant availability, offline data input | **Offline-First** | Writes data locally first, providing zero latency. |
| Financial transfers, ticket purchases, real-time inventories | **Online-First** | Requires real-time server verification to prevent fraud. |

### Matrix 2: Local Cache Storage vs. Server-Only Fetching
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| User settings, cached profiles, localized static data | **Local Cache** | Reduces API payload size and network consumption. |
| Account balances, transaction details, secure keys | **Server-Only** | Prevents sensitive data from being leaked if the device is lost. |

### Matrix 3: Push Notification vs. Email Alerts
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Time-sensitive alerts, chat updates, order delivery progress | **Push Notification** | Delivers real-time notifications directly to the device. |
| Receipts, account invoices, monthly activity reports | **Email** | Provides a permanent, searchable record for the user. |

### Matrix 4: Native Features vs. Cross-Platform Frameworks
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Unified codebases, cross-platform layouts, rapid iterations | **Cross-Platform (Flutter)** | Shares up to 90% of code across platforms, reducing costs. |
| Complex native views, platform-specific hardware APIs | **Native (Swift/Kotlin)** | Direct access to native OS components and SDK updates. |

### Matrix 5: Local Cache vs. Fresh Network Data
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| User profile views, catalog lists, slow-changing content | **Local Cache** | Fast page loads; updates cache in the background. |
| Product inventory counts, live pricing, account balances | **Fresh Network** | Displays only validated, real-time data from the server. |

### Matrix 6: Synchronization Policy Selection
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Simple user inputs, personal draft updates, single-user fields | **Last-Write-Wins (LWW)**| Simple implementation; low overhead. |
| Collaborative document edits, inventory updates | **Merge / Conflict Queue**| Prevents changes from being accidentally overwritten. |

---

## 13. AI Mobile Rules

AI agents modifying or writing mobile application code must follow these rules:

1. **Never Hardcode API Keys**: Do not write sensitive configuration variables or API secrets inside Flutter Dart files or Android/iOS build files.
2. **Never Trust Client Inputs**: Enforce server-side receipt validation for all in-app purchases.
3. **No UI Thread Blocks**: Offload CPU-heavy calculations (e.g. file serialization, encryption) to background threads.
4. **Use Const Constructor Declarations**: Enforce immutable class const initializers.
5. **No Blind Permission Requests**: Do not request device permissions (e.g., location, camera) unless the active feature explicitly requires it.

---

## 14. Mobile Review Checklist

Use this checklist during code review to evaluate mobile application changes.

### Architecture & Lifecycles
- [ ] Are business logic calculations separated from UI widget code?
- [ ] Does the application handle background state transitions and suspend cleanups?

### Network & Data Resiliency
- [ ] Do network calls use timeouts, retries, and exponential backoff?
- [ ] Does the sync queue append unique idempotency keys to POST requests?
- [ ] Do list layouts use pagination (no large, un-paginated requests)?

### Security & Privacy
- [ ] Are credentials, tokens, and sensitive data stored exclusively in secure Keychains/KeyStores?
- [ ] Is binary obfuscation enabled for release builds?
- [ ] Is SSL pinning configured in the HTTP client?

### Payments & Transactions
- [ ] Are all purchase validations processed server-side (no client-side validation)?
- [ ] Is there protection against receipt replay attacks?

### Performance & Usability
- [ ] Are CPU-heavy tasks offloaded to background threads (Isolates)?
- [ ] Are images downscaled (`cacheWidth`/`cacheHeight`) before being cached?
- [ ] Do interactive elements meet the minimum touch target size ($48 \times 48\text{ dp}$)?

---

## References
- Secure Database Schemas: [core/06-database-engineering-standard.md](06-database-engineering-standard.md)
- Security Engineering: [core/08-security-engineering-standard.md](08-security-engineering-standard.md)
- Flutter Development: [stacks/dart-flutter/flutter-dart-engineering-standard.md](../stacks/dart-flutter/flutter-dart-engineering-standard.md)
