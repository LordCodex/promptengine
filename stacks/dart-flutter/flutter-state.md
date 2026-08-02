---
document_id: stacks-flutter-state
title: Flutter State Management Standards
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

# Flutter State Management Standards

## Inheritance
This document inherits from and extends the [Universal Coding Standards](../../core/05-universal-coding-standards.md). Refer to the universal standards for class encapsulation, data mutation limits, and decoupling conventions. This page specifies only Flutter state tools configurations.

## Purpose
This document defines standards for state separation, event routing, and provider configurations using BLoC and Riverpod.

## Scope
Applies to client business state wrappers (`lib/application/`, `lib/bloc/`, `lib/providers/`).

---

## Directives

### 1. Unidirectional Data Flow
- **Standard**: User actions must dispatch events (or call state functions). The state manager processes the request and emits a new immutable state to the UI view. The UI remains read-only.

### 2. BLoC (Business Logic Component) Rules
- **Rule 1**: BLoCs must be event-driven. UI files should never invoke methods on a BLoC class directly. Use events:
  ```dart
  // UI Event dispatch
  context.read<AuthBloc>().add(const AuthLoginRequested(email, password));
  ```
- **Rule 2**: States emitted by BLoCs must be immutable. Define properties as `final` and implement `copyWith` for transitions:
  ```dart
  class AuthState {
    final AuthStatus status;
    final User? user;

    const AuthState({required this.status, this.user});

    AuthState copyWith({AuthStatus? status, User? user}) {
      return AuthState(
        status: status ?? this.status,
        user: user ?? this.user,
      );
    }
  }
  ```

### 3. Riverpod Provider Standards
- Use `ref.watch()` within widget build methods to observe state changes dynamically.
- Use `ref.read()` only within button callbacks or triggers to execute actions. Do not use `read` inside build cycles, as it does not trigger repaints.
- Use `select` to filter properties and prevent widget rebuilds on unrelated state changes:
  ```dart
  final username = ref.watch(userProvider.select((user) => user.username));
  ```

---

## Common Mistakes & Anti-Patterns
- **Direct State Mutation**: Modifying field properties directly in BLoC or provider instances (`bloc.state.status = Status.success`) instead of emitting new instances.
- **UI Logic in State**: Coupling layout formatting (e.g. returning a `Text` widget from a Bloc state property) inside application business classes.
- **Duplicate Provider Registration**: Registering multiple instances of providers or blocs without managing their lifecycle, causing stale data cache bugs.

---

## References
- Widget optimization: [flutter-widgets.md](file:///Users/kodexkode/Documents/workspace/promptengine/stacks/dart-flutter/flutter-widgets.md)
- REST API integration: [bridges/laravel-api-flutter.md](file:///Users/kodexkode/Documents/workspace/promptengine/bridges/laravel-api-flutter.md)
