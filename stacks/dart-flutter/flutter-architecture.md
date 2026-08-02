---
document_id: stacks-flutter-architecture
title: Flutter Architectural Layout and Integrations
ecosystem: dart-flutter
target_versions:
  flutter: "^3.0"
dependencies:
  - core-universal-coding-standards
  - stacks-dart-conventions
  - stacks-flutter-widgets
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Flutter Architectural Layout and Integrations

## Inheritance
This document inherits from and extends the [Universal Coding Standards](../../core/05-universal-coding-standards.md). Refer to the universal standards for module boundaries, public interfaces, and class decoupling. This page specifies feature directory structures and native integrations.

## Purpose
This document defines standards for folder layouts, GoRouter navigation configurations, localization setups, and native platform channel bridges.

## Scope
Applies to system layers, package dependencies, and routing setups.

---

## Directives

### 1. Feature-First Folder Structure
- **Standard**: Group code by business feature rather than technical layer. This simplifies navigating related modules.
- **Visual Mapping**:
  ```text
  lib/
  ├── main.dart
  └── features/
      ├── auth/
      │   ├── data/          # Api clients, local database entities
      │   ├── domain/        # Entities, business actions
      │   └── presentation/  # UI Widgets, screens, blocs/providers
      └── payment/
          ├── data/
          ├── domain/
          └── presentation/
  ```

### 2. Navigation with GoRouter
- **Rule 1**: Declare app routes declaratively using `GoRouter`.
- **Rule 2**: Enforce type-safe routing parameters and redirect routes dynamically based on authentication states:
  ```dart
  final appRouter = GoRouter(
    initialLocation: '/login',
    routes: [
      GoRoute(
        path: '/login',
        builder: (context, state) => const LoginScreen(),
      ),
      GoRoute(
        path: '/dashboard',
        builder: (context, state) => const DashboardScreen(),
        redirect: (context, state) {
          final isLoggedIn = checkAuthStatus();
          return isLoggedIn ? null : '/login';
        },
      ),
    ],
  );
  ```

### 3. Native Platform Method Channels
- **Rule**: Encapsulate all native code operations (`MethodChannel`) inside a dedicated repository class. Never invoke channel calls directly from widgets.
- **Why**: Keeps native dependencies mockable during testing.
- **Code Syntax**:
  ```dart
  class BatteryRepository {
    static const _channel = MethodChannel('playbook.dev/battery');

    Future<int> getBatteryLevel() async {
      try {
        final level = await _channel.invokeMethod<int>('getBatteryLevel');
        return level ?? 0;
      } on PlatformException catch (e) {
        logger.error('Failed to get battery', e);
        return -1;
      }
    }
  }
  ```

---

## Common Mistakes & Anti-Patterns
- **Context-Free Navigation**: Storing build context globally to trigger routes from remote service classes, which crashes the widget tree lifecycle.
- **Messy Flat Layout**: Placing all app screens, repositories, and models under a single flat directory structure.
- **Unsafely Typed Channel Returns**: Casting native platform return values directly without verifying type safety and null checks.

---

## References
- Client state transitions: [flutter-state.md](file:///Users/kodexkode/Documents/workspace/promptengine/stacks/dart-flutter/flutter-state.md)
- Testing structures: [flutter-testing.md](file:///Users/kodexkode/Documents/workspace/promptengine/stacks/dart-flutter/flutter-testing.md)
