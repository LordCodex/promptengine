# Dart Flutter Sandbox Reference

This directory serves as the executable reference implementation of the Dart/Flutter playbook standards. It provides a functional, linted boilerplate mapping exactly to the conventions defined in:
- [dart-conventions.md](../../stacks/dart-flutter/dart-conventions.md)
- [flutter-widgets.md](../../stacks/dart-flutter/flutter-widgets.md)
- [flutter-state.md](../../stacks/dart-flutter/flutter-state.md)

---

## Structure
- `/lib/user_dto.dart`: Data Transfer Object class showing type-safe JSON mapping.
- `/lib/user_badge.dart`: Simple widget illustrating key usage, `const` instantiation, and rebuild isolation.

---

## Validation and Analysis
To run static analysis checks against this sandbox block:
```bash
# Execute static analysis
flutter analyze
```
