---
document_id: core-testing-philosophy
title: Testing Philosophy and Patterns
ecosystem: cross-cutting
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Testing Philosophy and Patterns

## Purpose
This document establishes the testing principles required for all projects. It guides developers and AI agents to write maintainable, fast, and behavior-oriented test suites, preventing brittle tests that block code refactoring.

## Scope
Applies to backend unit/feature tests (Laravel Pest/PHPUnit) and frontend component/unit tests (Vue Vitest, Flutter Widget/Unit testing).

---

## The Testing Pyramid

We organize tests into three primary layers:

```text
      /\
     /  \   End-to-End / UI (5%) - Critical user flows, slow, high cost
    /----\
   /      \  Integration / Feature (30%) - API endpoints, DB queries, UI bindings
  /--------\
 /          \ Unit (65%) - Pure logic, calculations, zero I/O, fast
/____________\
```

### 1. Unit Tests
- **Focus**: Pure functions, mathematical algorithms, state transitions.
- **Constraints**: No database access, no network requests, no filesystem operations.
- **Execution Speed**: Under 10ms per test.

### 2. Integration / Feature Tests
- **Focus**: Testing components together (e.g., API controller routing to database, widget rendering triggering state manager changes).
- **Strategy**: Hit a test database (in transactions to auto-reset state) or local mock server. Mock external network gateways (Stripe, SendGrid) only.

### 3. End-to-End (E2E) / UI Boundary Tests
- **Focus**: Verify the exact integration of the systems (frontend, backend, database).
- **Strategy**: Limit these to core paths (e.g., "User signup to payment completion").

---

## Test Behavior, Not Implementation

> [!IMPORTANT]
> Tests must assert **what** the code does (public contract/behavior), not **how** it does it (private methods/internal state).

- **Bad**: Testing that a private method `$this->calculateDiscount()` is called with specific parameters.
- **Good**: Testing that an order with a premium user status returns a final total with 15% discount.
- *Reasoning*: If you refactor the internal implementation of `calculateDiscount` but the behavior remains unchanged, a behavior-oriented test passes. An implementation-oriented test breaks, causing cognitive overhead and slowing down development.

---

## Mocking, Stubbing, and Faking

Avoid over-mocking. If a unit/class calls a simple utility class, do not mock it; use the real implementation.

| Test Double | Purpose | When to Use |
| :--- | :--- | :--- |
| **Fake** | Working implementation with simplified logic | In-memory database, fake mail sender, fake local storage |
| **Stub** | Hardcoded return value to satisfy dependencies | Simulating a payment success response from Stripe API |
| **Mock** | Asserts calls and interactions occur | Asserting that an event-listener dispatches a background job |

---

## Example: Brittle vs. Robust Testing

### Scenario
An invoice service registers a customer purchase and emails a copy of the receipt.

#### The Brittle Approach (Avoid)
Mocking every class using dependency assertion:
```php
public function test_purchase_sends_email() {
    $mailMock = Mockery::mock(Mailer::class);
    $mailMock->shouldReceive('send')
             ->once()
             ->with('invoice_template', ['amount' => 100])
             ->andReturn(true);

    $service = new InvoiceService($mailMock);
    $service->processPurchase(100);
}
```
*Why this is bad*: If the template variable changes or the mailer class method signature gets refactored, the test fails, even if the user still receives the correct invoice.

#### The Robust Approach (Prefer)
Use system assertion and fakes:
```php
public function test_purchase_sends_email() {
    // Laravel Fake mailer setup
    Mail::fake();

    $service = resolve(InvoiceService::class);
    $service->processPurchase(100);

    // Assert target state occurred rather than class parameters
    Mail::assertSent(InvoiceMail::class, function ($mail) {
        return $mail->amount === 100;
    });
}
```
*Why this is good*: Verifies the business intent (that an InvoiceMail gets dispatched) without locking the service to specific mail engine signatures.

---

## References
- Planning Test Cases: [01-thinking-and-planning.md](file:///Users/kodexkode/Documents/workspace/promptengine/core/01-thinking-and-planning.md)
- Laravel concrete Pest setup: [stacks/php-laravel/laravel-testing.md](file:///Users/kodexkode/Documents/workspace/promptengine/stacks/php-laravel/laravel-testing.md)
- Flutter test implementations: [stacks/dart-flutter/flutter-testing.md](file:///Users/kodexkode/Documents/workspace/promptengine/stacks/dart-flutter/flutter-testing.md)
