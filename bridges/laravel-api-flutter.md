---
document_id: bridge-laravel-api-flutter
title: Laravel API and Flutter Integration Bridge
ecosystem: cross-cutting
target_versions:
  laravel: ">=10.0"
  flutter: "^3.0"
dependencies:
  - stacks-laravel-routing
  - stacks-flutter-architecture
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Laravel API and Flutter Integration Bridge

## Purpose
This document defines standards for authentication handshakes, payload formatting, and network data serialization between a Laravel backend API and a Flutter mobile client.

## Scope
Applies to API Controllers (Laravel) and Data Repositories/API clients (Flutter/Dart).

---

## Directives

### 1. Stateless Authentication (Bearer Tokens)
- **Standard**: Authenticate Flutter mobile requests using bearer tokens (e.g. Laravel Sanctum or JWT). Do not use cookie session states.
- **Rules**:
  - Flutter must store tokens securely on-device using a keychain wrapper (e.g. `flutter_secure_storage`).
  - Attach the token to all outgoing HTTP headers: `Authorization: Bearer <token>`.
  - Handle `401 Unauthorized` responses gracefully by redirecting the user to the Login screen or executing a refresh-token handshake.

### 2. Payload Case Reconciliation
- **Rule**:
  - Laravel API endpoints must return and ingest `snake_case` JSON properties.
  - Flutter Dart models must parse this payload into strongly-typed `camelCase` properties internally, converting it back to `snake_case` when serializing payloads to the server.
- **Code Syntax (Dart Model parsing)**:
  ```dart
  class UserDto {
    final String id;
    final String emailAddress; // camelCase internal

    const UserDto({required this.id, required this.emailAddress});

    factory UserDto.fromJson(Map<String, dynamic> json) {
      return UserDto(
        id: json['id'] as String,
        emailAddress: json['email_address'] as String, // parses snake_case API
      );
    }

    Map<String, dynamic> toJson() {
      return {
        'id': id,
        'email_address': emailAddress, // serializes to snake_case API
      };
    }
  }
  ```

### 3. API Contract Validation
- Use a mock server or contract testing setup to verify that changes to database columns do not break mobile serialization pipelines (refer to [legacy/02-backward-compatibility.md](file:///Users/kodexkode/Documents/workspace/promptengine/legacy/02-backward-compatibility.md)).

---

## Common Mistakes & Anti-Patterns
- **Dynamic Type Parsing**: Parsing dynamic arrays directly in Dart without explicit type checks (`json['items'] as List<dynamic>`), triggering runtime casting crashes.
- **Ignoring API Error Envelopes**: Reading response bytes on the mobile client without checking HTTP status codes, leading to parsing errors on `500` or `422` error pages.
- **Plaintext Storage**: Storing access tokens in standard unencrypted storage plugins (like `shared_preferences`) on local mobile device configurations.

---

## References
- REST API standard envelopes: [core/03-data-and-api-modeling.md](file:///Users/kodexkode/Documents/workspace/promptengine/core/03-data-and-api-modeling.md)
- Authorization controls: [security/02-auth-and-permissions.md](file:///Users/kodexkode/Documents/workspace/promptengine/security/02-auth-and-permissions.md)
