---
document_id: stacks-php-conventions
title: PHP Engineering Standard
ecosystem: php-laravel
target_versions:
  php: ">=8.3"
dependencies:
  - core-universal-coding-standards
  - core-architecture-and-simplicity
audience: [human, agent]
last_reviewed: 2026-08-01
---

# PHP Engineering Standard

## Inheritance & Alignment
This document inherits from and extends the [Universal Coding Standards](../../core/05-universal-coding-standards.md) and the [Architecture Standards](../../core/02-architecture-and-simplicity.md). It outlines PHP-specific conventions, design decision matrices, security mitigations, and performance targets.

---

## 1. PHP Philosophy & Modern Standards

Modern PHP is a strongly-typed, compile-optimized object-oriented language. We prioritize **Strict Typing**, **Explicitness**, and **Defensive Type Boundaries**.

### Strict Typing (Mandatory)
Every PHP file must declare strict types before executing any code. This prevents runtime type casting anomalies.
```php
<?php

declare(strict_types=1);

namespace App\Services;
```

---

## 2. Project Structure & Composer Integration

### Namespace & Autoload Standards
- All projects must conform to the **PSR-4** autoloading standard defined inside `composer.json`.
- Keep the `composer.json` clean, optimized, and validated:
  ```json
  "autoload": {
      "psr-4": {
          "App\\": "app/"
      }
  }
  ```
- **Autoload Optimization (Production)**: Always run optimize flags during builds:
  ```bash
  composer dump-autoload -o --apcu --no-dev
  ```

---

## 3. Decision Matrices

When building system layers, developers and AI agents must refer to the following decision tables:

### Matrix A: Class State Properties
| Choice | Use Case | Why / Trade-offs |
| :--- | :--- | :--- |
| **Readonly Class / Properties** | Immutable Data Transfer Objects (DTOs), Configuration blocks | **Why**: Prevents state modification post-instantiation.<br>**Trade-off**: Requires full constructor recreation to alter attributes. |
| **Mutable Properties** | Dynamic cache stores, transient workspace models | **Why**: Allows properties update in-memory.<br>**Trade-off**: Increases state-tracking bugs. |

### Matrix B: Custom Type Groupings
| Choice | Use Case | Why |
| :--- | :--- | :--- |
| **Backing Enum (PHP 8.1+)** | Fixed state groups stored in DB (e.g., `Status::Active`) | Type-safety, strict matches, native database values mapping. |
| **Class Constants** | Low-overhead internal system configurations, math parameters | Lightweight, but lacks runtime type safety. |

### Matrix C: Interface vs. Concrete Class
- **Interface**: Use only when multiple class implementations exist, or when wrapping third-party boundaries (like billing channels) for mock test verification.
- **Concrete Class**: Default to concrete classes for internal business services and logic workflows to avoid abstraction bloat (refer to [core/02-architecture-and-simplicity.md](file:///Users/kodexkode/Documents/workspace/promptengine/core/02-architecture-and-simplicity.md)).

### Matrix D: DTO vs. Associative Array
- **DTO**: Use for all class parameters crossing boundary limits (HTTP payload $\rightarrow$ Action class). Provides IDE completion and type assertions.
- **Associative Array**: Use only for temporary internal array maps within private methods.

---

## 4. Modern Language Features & OOP Design

### Constructor Property Promotion
Declare visibility, types, and properties directly in constructor parameters:
```php
class InvoiceService {
    public function __construct(
        private readonly PaymentGatewayInterface $gateway,
        private readonly LoggerInterface $logger
    ) {}
}
```

### Generics Simulation (PHPDoc Annotations)
Since PHP lacks native runtime generics, type assertions for arrays must be documented in PHPDoc annotations for static analysis engines (PHPStan/Psalm):
```php
/**
 * @param array<int, UserDto> $users
 * @return array<string, Invoice>
 */
public function processBatch(array $users): array { ... }
```

### Traits vs. Composition
- **Rule**: Avoid utilizing Traits. Traits hide dependency parameters, complicate testing mocks, and obscure class inheritance boundaries.
- **Action**: Inject classes via Dependency Injection (DI) instead of importing Traits.

---

## 5. Security Engineering (PHP Specific)

### 1. Unsafe PHP Functions Exclusion
The following functions are strictly prohibited due to code injection risks. Static analysis checkers must block them:

| Prohibited Function | Safe Alternative | Reason |
| :--- | :--- | :--- |
| `eval()` | Native algorithms, dynamic class mapping | Arbitrary string execution vulnerability |
| `unserialize()` | `json_decode()` | PHP Object Injection, remote code execution (RCE) |
| `exec()`, `system()`, `shell_exec()` | Symfony `Process` component | Shell injection, unescaped parameter execution |

### 2. Cryptography and Passwords
- Always hash passwords using `password_hash()` with `PASSWORD_ARGON2ID` or `PASSWORD_BCRYPT`. Never implement custom hashing.
- Generate cryptographically secure random values using `random_bytes()` or `random_int()`. Never use `rand()` or `mt_rand()`.

---

## 6. Performance & Databases Optimizations

### OPcache Optimization (Production)
Configure PHP runtime parameters to cache compiled byte-code in memory:
```ini
opcache.enable=1
opcache.validate_timestamps=0 ; Disables file modification checks in production
opcache.memory_consumption=256
opcache.interned_strings_buffer=16
```

### Streaming Large Datasets (Generators)
- **Problem**: Querying millions of rows using standard array returns blows memory limits (Out of Memory error).
- **Solution**: Use PHP **Generators** (`yield`) to stream database cursors:
  ```php
  /**
   * @return Generator<int, Transaction>
   */
  public function getLargeExport(): Generator {
      $cursor = $this->db->query("SELECT * FROM transactions");
      while ($row = $cursor->fetch()) {
          yield new Transaction($row);
      }
  }
  ```

---

## 7. Static Analysis & Code Quality Tools

All PHP projects must configure strict code validation tools. These pipelines must block commits on failure (refer to [environment/03-ci-cd-pipelines.md](file:///Users/kodexkode/Documents/workspace/promptengine/environment/03-ci-cd-pipelines.md)).

### PHPStan Strict Levels
- **New Projects**: Configure PHPStan `phpstan.neon` to **Level 8** or **Level 9**. Enforce strict null checking and type completions.
- **Legacy Projects**:
  1. Boot PHPStan at **Level 5**.
  2. Generate a baseline configuration (`phpstan-baseline.neon`) to whitelist existing errors.
  3. Incremental updates: Elevate the level step-by-step as you refactor.

### Rector (Automated Refactoring)
Use Rector configurations to automatically modernize legacy code syntax (e.g. upgrading array syntax or migrating from PHP 7.x constructor setups to PHP 8.x).

---

## 8. Anti-Patterns to Avoid

- **Fat Helpers (`Helpers.php`)**: Creating a global file full of procedural functions, generating circular dependencies.
- **Magic Constant Abuse**: Utilizing magic numeric values (`if ($status === 3)`) instead of typed Enums.
- **Underengineering Database Transactions**: Modifying multiple SQL tables without wrapping calls inside database transaction blocks, causing partial data commits.

---

## 9. AI Agent PHP Directives

AI agents modifying PHP files in this repository must follow these rules:
1. **Declare Strict Types**: Verify every new or refactored `.php` file begins with `declare(strict_types=1);`.
2. **Explicit Constructor Promotion**: Use property promotion for class dependencies.
3. **No Unsafe Functions**: Do not use `unserialize()` or shell executions.
4. **Mock repositories**: Do not instantiate raw Eloquent/SQL structures inside unit tests; mock repositories using standard interfaces.

---

## 10. PHP Code Review Checklist

- [ ] **Types**: Are parameters, properties, and method returns fully typed?
- [ ] **OPcache**: Are settings optimized for production?
- [ ] **Transactions**: Are modifications wrapped inside database transactions?
- [ ] **Static Checks**: Does the code run green through PHPStan at the designated project level?
- [ ] **Security**: Are all inputs validated and output parameters escaped?

---

## References
- Universal Naming Rules: [core/05-universal-coding-standards.md](file:///Users/kodexkode/Documents/workspace/promptengine/core/05-universal-coding-standards.md)
- Database Optimizations: [performance/02-database-optimization.md](file:///Users/kodexkode/Documents/workspace/promptengine/performance/02-database-optimization.md)
- Environment Setup: [environment/01-local-dev-standards.md](file:///Users/kodexkode/Documents/workspace/promptengine/environment/01-local-dev-standards.md)
