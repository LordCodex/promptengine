---
document_id: stacks-laravel-testing
title: Laravel Testing Conventions
ecosystem: php-laravel
target_versions:
  laravel: ">=10.0"
dependencies:
  - core-universal-coding-standards
  - stacks-php-conventions
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Laravel Testing Conventions

## Inheritance
This document inherits from and extends the [Universal Coding Standards](../../core/05-universal-coding-standards.md). Refer to the universal standards for testing scopes, assertions boundaries, and mock configurations. This page specifies only Laravel-specific testing tools (Pest/PHPUnit).

## Purpose
This document establishes rules for writing feature and unit tests in Laravel, setting up database testing states, and mocking third-party service gateways.

## Scope
Applies to test suites (`tests/Feature/`, `tests/Unit/`) using Pest or PHPUnit.

---

## Directives

### 1. Database Transaction Isolation
- **Rule**: Every integration/feature test that accesses the database must use the `Illuminate\Foundation\Testing\RefreshDatabase` or `DatabaseTransactions` trait.
- **Why**: Enforces test isolation. Transactions roll back database changes automatically at the end of each test run, preventing state pollution.

### 2. Standardized Feature Test Structure
- **Focus**: Test API endpoints, controllers, routing configurations, and middleware blocks.
- **Code Syntax (Pest)**:
  ```php
  use App\Models\User;
  use function Pest\Laravel\{actingAs, postJson};

  uses(Tests\TestCase::class, Illuminate\Foundation\Testing\RefreshDatabase::class);

  test('user can register with valid parameters', function () {
      $payload = [
          'email' => 'newuser@example.com',
          'name' => 'John Doe',
          'password' => 'SecurePassword123!',
      ];

      postJson(route('api.register'), $payload)
          ->assertStatus(201)
          ->assertJsonStructure([
              'data' => ['id', 'email', 'name']
          ]);

      $this->assertDatabaseHas('users', ['email' => 'newuser@example.com']);
  });
  ```

### 3. Model Factories for State Definition
- **Rule**: Never use raw database inserts (`DB::insert`) or manual user models generation in test cases. Always use Laravel Model Factories.
- **Code Syntax**:
  ```php
  // Safe model state generation
  $admin = User::factory()->admin()->create();
  ```

### 4. Mocking External Services
- Use standard facade testing fakes (`Http::fake()`, `Mail::fake()`, `Queue::fake()`) instead of instantiating manually mocked interfaces (refer to [core/04-testing-philosophy.md](file:///Users/kodexkode/Documents/workspace/promptengine/core/04-testing-philosophy.md)).

---

## Common Mistakes & Anti-Patterns
- **Testing Live Stripe Gateways**: Failing to fake HTTP integrations, causing test suites to trigger actual web network calls to third-party endpoints.
- **Leftover DB Records**: Forgetting to include the `RefreshDatabase` trait, causing tests to fail randomly due to duplicate keys remaining in the local DB.
- **Unit Testing Eloquent**: Writing unit tests for model relations without booting the application environment, throwing fatal driver errors.

---

## References
- Testing philosophy: [core/04-testing-philosophy.md](file:///Users/kodexkode/Documents/workspace/promptengine/core/04-testing-philosophy.md)
- Controller parameters: [laravel-routing.md](file:///Users/kodexkode/Documents/workspace/promptengine/stacks/php-laravel/laravel-routing.md)
