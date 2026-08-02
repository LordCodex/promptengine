---
document_id: stacks-dart-conventions
title: Modern Dart Core Conventions
ecosystem: dart-flutter
target_versions:
  dart: "^3.0"
dependencies:
  - core-universal-coding-standards
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Modern Dart Core Conventions

## Inheritance
This document inherits from and extends the [Universal Coding Standards](../../core/05-universal-coding-standards.md). Refer to the universal standards for general naming rules, variables scopes, and error structures. This page specifies only Dart-specific typing and asynchronous operators.

## Purpose
This document defines standards for coding in Dart, enforcing sound null safety, compile-time constants, asynchronous Streams/Futures structures, and runtime type checks.

## Scope
Applies to all Dart source code files (`.dart`).

---

## Directives

### 1. Sound Null Safety
- **Rule**: Avoid using the bang operator (`!`) for null assertions unless absolutely unavoidable. Use conditional member access (`?.`), null-coalescing operators (`??`), or early guard clauses instead.
- **Why**: The bang operator throws runtime `NullThrownError` exceptions if the variable is null, bypassing Dart's compile-time safety.
- **Code Syntax**:
  ```dart
  // Safe Null checking
  final String name = payload.username ?? 'guest';
  ```

### 2. Compile-Time Constants (const)
- **Rule**: Use the `const` constructor prefix aggressively for widgets and collections that do not change at runtime.
- **Why**: Allows Dart to cache instance trees in memory, decreasing GC overhead and speeding up UI rendering cycles.
- **Code Syntax**:
  ```dart
  const List<String> roleTypes = ['admin', 'marketer', 'author'];
  ```

### 3. Asynchronous Futures and Streams
- Use `async`/`await` for standard Future operations rather than chaining `.then()` callbacks.
- When generating async sequences, use async generator functions (`async*`) and the `yield` operator:
  ```dart
  Stream<int> countStream(int max) async* {
    for (int i = 1; i <= max; i++) {
      yield i;
    }
  }
  ```

---

## Common Mistakes & Anti-Patterns
- **The Bang Operator Abuse**: Chaining `user!.profile!.email` which crashes at runtime if any property is null.
- **Floating Async Functions**: Dispatched futures that execute without proper error catching blocks (`unawaited(myFutureAction())` with no local try/catch).
- **String Concatenation in Loops**: Constructing strings by appending values in loops (`result += value`) instead of instantiating a clean `StringBuffer`.

---

## References
- System performance guidelines: [performance/04-concurrency-and-async.md](file:///Users/kodexkode/Documents/workspace/promptengine/performance/04-concurrency-and-async.md)
- Flutter widgets usage: [flutter-widgets.md](file:///Users/kodexkode/Documents/workspace/promptengine/stacks/dart-flutter/flutter-widgets.md)
