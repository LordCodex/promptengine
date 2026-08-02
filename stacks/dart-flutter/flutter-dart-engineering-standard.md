---
document_id: stacks-flutter-dart-engineering-standard
title: Flutter and Dart Engineering Standard
ecosystem: dart-flutter
target_versions:
  flutter: "^3.19"
  dart: "^3.3"
dependencies:
  - core-universal-coding-standards
  - core-architecture-and-simplicity
  - core-security-engineering-standard
  - core-testing-engineering-standard
  - stacks-dart-conventions
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Flutter and Dart Engineering Standard

## Purpose & Inheritance
This document defines the core standards for mobile application development using Flutter and Dart. It inherits from and extends the [Universal Coding Standards](../../core/05-universal-coding-standards.md), the [Architecture Standards](../../core/02-architecture-and-simplicity.md), the [Security Engineering Standard](../../core/08-security-engineering-standard.md), and the [Testing Engineering Standard](../../core/11-testing-engineering-standard.md). It establishes implementation protocols for widget layouts, state controllers, local keychains storage, background isolate tasks, and offline repositories synchronization.

---

## 1. Flutter Engineering Philosophy

A Flutter application is a **highly stateful, multi-threaded client system, not just a set of UI screens**.
- **Separate Presentation from Logic**: UI widgets must be completely free of business logic. Widgets should only render layouts based on state models and delegate user interactions to controller layers.
- **Model Immutable States**: Treat application state as immutable data frames. Modifying state must trigger a predictable redraw flow rather than mutations of inline variables.
- **Isolate Network & Storage Overhead**: UI rendering cycles execute at $60\text{Hz}$ or $120\text{Hz}$. Any slow operations (e.g. database parsing, HTTP requests, file compression) must execute on background threads (Isolates) to prevent frame drops.

---

## 2. Dart Language Standards & Code Quality

Dart is a strongly typed language. We enforce its features to catch errors at compile time.

### Dart Type System
- **Strict Null Safety**: Never use null-assertion operators (`!`) unless you can guarantee the variable is not null (e.g. immediately following a null check). Prefer conditional checks (`?.`) or fallback default bindings (`??`).
- **Records & Pattern Matching**: Leverage Dart records to return multiple values from functions cleanly. Use pattern matching to destructure data:
  ```dart
  // Good: Return record type and match pattern in receiver
  (int code, String message) fetchStatus() {
    return (200, 'Success');
  }

  void main() {
    final (code, message) = fetchStatus();
    print('Code: $code, Msg: $message');
  }
  ```
- **Prefer Extensions Over Utility Classes**: If a helper method is strictly bound to a specific class (like formatting a `DateTime` object), implement it as an `extension` rather than creating a global utility class:
  ```dart
  extension DateTimeFormatting on DateTime {
    String get toIsoDateString => '${year.toString()}-${month.toString().padLeft(2, '0')}-${day.toString().padLeft(2, '0')}';
  }
  ```
- **Asynchronous Error Handling**: Always wrap async operations in `try-catch` blocks and return typed Result models rather than throwing raw, untyped errors across architectural layers.

---

## 3. Project Directory Structures

Select the project layout structure that matches your application's complexity:

### Option A: Feature-Based Architecture (Recommended Default)
Group files by feature domains. Highly modular; easy to scale across multiple teams.

```text
lib/
├── core/             # Shared API clients, global routing configs, themes
├── shared/           # Generic buttons, input widgets, dialog components
└── features/
    ├── invoicing/
    │   ├── data/          # API models, database adapters, providers
    │   ├── domain/        # Business models, validation services
    │   ├── presentation/  # Widgets, screen layouts, controller states
    │   └── invoicing_feature.dart
```

### Option B: Clean Architecture
Enforces strict boundaries between Presentation, Domain, and Data. Recommended only for enterprise projects with large datasets and complex business rules. Avoid for simple CRUD applications as it introduces excessive boilerplate.

---

## 4. State Management Standards

Select a state management framework that matches the scale and requirements of your project:

### State Management Frameworks Matrix
- **Riverpod**: The recommended default for medium-to-large projects. It provides compile-safe dependency injection, global state isolation, and handles asynchronous network caching out of the box.
- **Bloc / Cubit**: Recommended for enterprise teams that require strict event-driven state transitions, formal states logging, and strict state machine transitions.
- **ValueNotifier / StatefulWidget**: Best for localized page-specific states (e.g., controlling a cursor position, text input updates, tab changes). Do not use for global application state.

### State Scoping Rules
- **Local State**: Keep UI interaction state (e.g. whether a modal is open) local to the widget using `StatefulWidget` or `ValueNotifier`.
- **Global Shared State**: User accounts, billing scopes, and local databases caches must reside inside global Riverpod providers or Bloc controllers.

---

## 5. Widget Engineering & Composition

### Widget Construction Rules
- **Const Constructors Everywhere**: Always prefix constructor declarations with the `const` keyword if the widget has immutable properties. This allows Flutter to cache the widget instance and bypass redundant redraw cycles.
- **Decompose Large Widgets**: Keep build methods short (under 100 lines). If a widget tree has deep nested structures, extract sub-trees to separate stateless widgets.
- **Responsive Layouts**: Do not hardcode layout pixel sizes. Use `MediaQuery` or layout builders (`LayoutBuilder`) to ensure visual screens render correctly across multiple viewport sizes (mobile, tablet, desktop).
- **Accessibility (A11y)**: Wrap images and icons in `Semantics` widgets to provide descriptive labels for screen readers.

### Widget Anti-Patterns
- **Massive Build Blocks**: Putting hundreds of layout widgets and inline functions inside a single build method.
- **Inline Business Logic**: Fetching API data or calculating payment values inside a button's `onPressed` callback. Call a controller instead.
- **Rebuilding Unnecessary Widgets**: Calling `setState` at the root widget level to modify a nested child element.

---

## 6. Data Layer & API Integrations

The data layer abstracts network APIs and local storage engines behind clean Repository interfaces.

### API Client (Dio) Configuration
- **Use Dio Client**: Use the `Dio` library instead of Dart's raw HTTP client. Dio supports network interceptors, custom timeout limits, and global error handling out of the box.
- **JWT Token Refresh Interceptor**: Automate JWT token refreshes using Dio request interceptors to retry failed requests seamlessly:

```dart
// Good: JWT Auto-Refresh Interceptor
class AuthInterceptor extends Interceptor {
  final Dio dio;
  
  AuthInterceptor(this.dio);

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) async {
    if (err.response?.statusCode == 401) {
      final success = await refreshSessionToken();
      if (success) {
        // Retry the original failed request
        final response = await dio.fetch(err.requestOptions);
        return handler.resolve(response);
      }
    }
    return handler.next(err);
  }
}
```

- **JSON Serialization**: Do not parse JSON maps manually. Enforce automatic serialization code generation using `json_serializable` or `freezed` schemas to catch structural mismatches at compile time.

---

## 7. Local Storage & Mobile Security

### Sensitive vs. Non-Sensitive Data
- **Non-Sensitive Data**: Use `shared_preferences` or `Hive` to cache theme states, user preferences, and app configurations.
- **Sensitive Data**: Store API tokens, passwords, and private identifiers exclusively inside secure OS containers (iOS Keychain, Android KeyStore) using `flutter_secure_storage`. Never write plain secrets to database caches.

### Mobile Security Hardening
- **SSL Pinning**: Configure SSL Pinning inside your HTTP client to prevent Man-in-the-Middle (MitM) packet interceptions on public networks.
- **Obfuscation**: Always build release binaries with obfuscation flags enabled (`flutter build apk --obfuscate --split-debug-info=...`) to make reverse-engineering difficult.
- **Screenshot Protection**: Disable screenshot captures on pages containing sensitive financial or personal data (using secure flags in native Android/iOS window layouts).

---

## 8. Performance Optimization

- **ListView.builder for Lists**: Never render list lists using a basic `ListView` wrapper. Always use `ListView.builder` to dynamically construct list nodes only when they enter the visible screen viewport.
- **Reduce Widget Rebuilds**: Keep state scopes small. In Riverpod, use the `select()` operator to listen only to the specific state properties your widget requires:
  ```dart
  // Redraws only when the customer's payment status changes
  final isPaid = ref.watch(invoiceProvider.select((inv) => inv.isPaid));
  ```
- **Offload Heavy Work to Isolates**: Do not run CPU-intensive operations (such as parsing large JSON files, compressing images, or running encryption algorithms) on the main thread. Run them on a separate thread using `compute()` or raw Isolates:
  ```dart
  // Offloads heavy calculations to a background thread
  final processedData = await compute(parseLargeJson, rawString);
  ```

---

## 9. Testing Strategy

### Automated Testing Protocols
- **Unit Tests**: Test repositories, state controllers, and pure business algorithms.
- **Widget Tests**: Verify single components layout structures, tap behaviors, and text render outputs:
  ```dart
  testWidgets('InvoiceRow displays amount and triggers tap', (WidgetTester tester) async {
    bool isTapped = false;
    await tester.pumpWidget(MaterialApp(
      home: Scaffold(
        body: InvoiceRow(
          amountCents: 5000,
          onTap: () => isTapped = true,
        ),
      ),
    ));

    expect(find.text('\$50.00'), findsOneWidget);
    await tester.tap(find.byType(InvoiceRow));
    expect(isTapped, isTrue);
  });
  ```
- **Golden UI Tests**: Assert pixel-perfect visual alignments of critical UI primitives across screen adjustments.
- **Integration Tests**: Boot the application on real devices or emulator pipelines to verify end-to-end flows.

---

## 10. Release Engineering & Build Configurations

- **App Flavors**: Enforce build flavors (`dev`, `staging`, `prod`) to keep API urls, app icons, and signing credentials separate:
  ```bash
  flutter build apk --flavor prod -t lib/main_prod.dart
  ```
- **App Signing**: Store signing keystores and private certs outside the repository. Inject signing keys using CI environment variables during build steps.
- **Crashlytics Monitoring**: Integrate crash reporters (like Firebase Crashlytics) at startup to capture runtime exceptions in production.

---

## 11. Legacy Refactoring

When working on legacy Flutter codebases containing outdated packages, deprecated Dart options, or massive unstructured files:
1. **Enable Null Safety**: Ensure the application compile configuration enforces strict null safety.
2. **Isolate Widgets First**: Extract large nested widgets from massive build methods before refactoring state layers.
3. **Write Widget Behavior Tests**: Write widget tests as code baselines before swapping out deprecated state libraries.
4. **Upgrade Dependencies Incrementally**: Update dependency versions one by one rather than modifying `pubspec.yaml` in bulk.

---

## 12. Decision Matrices

Use these matrices to identify the correct Flutter engineering decision based on project context.

### Matrix 1: Riverpod vs. Bloc (State Management)
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Modular setups, API data integrations, rapid features prototyping | **Riverpod** | Auto-dispose configurations; simple API state caching. |
| Strict transactional states, complex logging pipelines, large enterprise | **Bloc** | Enforces event-driven state transitions; highly auditable. |

### Matrix 2: Clean Architecture vs. Simple Feature-Based Structure
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Small-to-medium apps, simple CRUD screens | **Simple Feature** | Low boilerplate; fast feature iteration. |
| Enterprise apps, multiple integration APIs, complex business logic | **Clean Architecture** | Enforces strict boundaries; decoupled domain entities. |

### Matrix 3: Local State vs. Global Shared State
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Page text fields inputs, tab selections, visual animations | **Local State** | Fast; localized to the widget lifecycle. |
| Session tokens, shopping carts data, global notifications | **Global State** | Syncs state across separate screen scopes. |

### Matrix 4: Repository vs. Direct API Client Queries
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Production code, unit-testable data structures | **Repository** | Decouples data retrieval from business logic. |
| Rapid UI prototyping, one-off throwaway test screens | **Direct Queries** | Low initial boilerplate setup. |

### Matrix 5: Future vs. Stream
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| One-off API GET requests, file load operations | **Future** | Simple trigger-and-return async pattern. |
| Live chat messages, web sockets feeds, database query listener | **Stream** | Listens to a continuous flow of asynchronous data updates. |

### Matrix 6: Local Database Caches vs. Network-Only Fetching
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Offline-first features, poor internet areas applications | **Local Database** | Data is read instantly from local storage; syncs in background. |
| Core real-time transactional data (e.g. transfers lists) | **Network-Only** | Guarantees that users always view the most up-to-date data. |

---

## 13. AI Flutter Rules

AI agents modifying or writing Flutter code in this repository must follow these rules:

1. **Keep Build Methods Short**: Do not generate build methods exceeding 100 lines. Extract nested sub-trees to separate stateless widgets automatically.
2. **Never mutation props directly**: Do not modify widget properties inside child components. Delegate updates to controller interfaces.
3. **No raw secure variables**: Ensure sensitive keys use secure storage plugins (never save API keys in SharedPreferences).
4. **Use Const Constructors**: Auto-apply the `const` keyword to immutable widget constructors.
5. **Never run CPU-intensive tasks on UI thread**: Wrap file parsers or JSON translations inside `compute()` isolates.

---

## 14. Mobile Code Review Checklist

Use this checklist during code review to evaluate Flutter and Dart changes.

### Dart & Code Quality
- [ ] Are all variables typed securely (no raw `dynamic` or `any` values)?
- [ ] Does async code wrap API calls in try-catch logic (returning typed Result envelopes)?
- [ ] Are extension methods used for type-specific formatting (no global utilities classes)?

### Architecture & Directory Alignment
- [ ] Is layout code completely separated from business logic?
- [ ] Are files placed in their correct feature directories?

### Widgets & Layout
- [ ] Are build methods short and readable (under 100 lines)?
- [ ] Are widget constructors configured as `const` where possible?
- [ ] Do lists use `ListView.builder` for virtualized rendering?

### State Management
- [ ] Is Riverpod or Bloc used correctly (no unnecessary global states)?
- [ ] Do widgets use `select()` to watch only the specific state fields they need?

### Security Hardening
- [ ] Are API tokens and sensitive credentials stored in secure storage?
- [ ] Is SSL pinning configured in the HTTP client (Dio)?
- [ ] Is obfuscation enabled for release builds?

### Testing
- [ ] Do widget tests verify visual outputs and user inputs?
- [ ] Are CPU-heavy tasks offloaded to background threads (Isolates)?

---

## References
- Universal Naming Rules: [core/05-universal-coding-standards.md](../../core/05-universal-coding-standards.md)
- Security Engineering: [core/08-security-engineering-standard.md](../../core/08-security-engineering-standard.md)
- Automated CI Pipelines: [core/13-cicd-and-deployment-standard.md](../../core/13-cicd-and-deployment-standard.md)
