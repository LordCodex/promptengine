---
document_id: stacks-flutter-testing
title: Flutter Testing and Mocking Conventions
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

# Flutter Testing and Mocking Conventions

## Inheritance
This document inherits from and extends the [Universal Coding Standards](../../core/05-universal-coding-standards.md). Refer to the universal standards for testing pyramids, behavior-oriented verification, and mocking limits. This page specifies only Flutter-specific testing harnesses.

## Purpose
This document defines standards for testing Flutter interfaces, executing widget checks, performing golden tests, and configuring mock classes.

## Scope
Applies to Flutter test files (`test/`, `integration_test/`).

---

## Directives

### 1. Widget Testing (PumpWidget)
- **Standard**: Test UI component interactivity by simulating widget frames in memory.
- **Rule**: Always wrap widgets with a root test container providing required routing and localization contexts:
  ```dart
  void main() {
    testWidgets('UserCard displays title and triggers event', (tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: UserCard(title: 'John Doe'),
          ),
        ),
      );

      // Verify widget renders title
      expect(find.text('John Doe'), findsOneWidget);

      // Trigger tap and verify
      await tester.tap(find.byType(UserCard));
      await tester.pump(); // Pump frame
    });
  }
  ```

### 2. Golden UI Component Verification
- **Standard**: Golden tests verify visual rendering layouts by comparing active pixels against a saved image baseline.
- **Rule 1**: Only use goldens for cohesive, low-level UI elements (e.g. custom buttons, status icons). Do not run goldens on large full pages containing dynamic content.
- **Rule 2**: Run tests under a fixed screen size constraint:
  ```dart
  await tester.binding.setSurfaceSize(const Size(800, 600));
  ```

### 3. Mocking Native Channels
- When running test scripts, override Method Channels to simulate native device metrics (e.g. hardware sensors, platform versions):
  ```dart
  tester.binding.defaultBinaryMessenger.setMockMethodCallHandler(
    const MethodChannel('playbook.dev/battery'),
    (MethodCall methodCall) async {
      if (methodCall.method == 'getBatteryLevel') {
        return 42; // Mock return
      }
      return null;
    },
  );
  ```

---

## Common Mistakes & Anti-Patterns
- **Unawaited Pump Calls**: Forgetting to execute `await tester.pump()` or `await tester.pumpAndSettle()` after an interaction, causing verification code to run on a legacy layout state.
- **Dynamic Asset Failures**: Running golden tests containing online network image lookups, which fail dynamically under offline build environments.
- **Mocking State Classes**: Mocking state notifier providers or BLoCs directly in widget tests instead of mock-testing the data repositories and letting the state manager behave natively.

---

## References
- Testing philosophy: [core/04-testing-philosophy.md](file:///Users/kodexkode/Documents/workspace/promptengine/core/04-testing-philosophy.md)
- Flutter widgets: [flutter-widgets.md](file:///Users/kodexkode/Documents/workspace/promptengine/stacks/dart-flutter/flutter-widgets.md)
