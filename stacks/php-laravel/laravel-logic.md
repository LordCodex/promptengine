---
document_id: stacks-laravel-logic
title: Laravel Business Logic and Service Layers
ecosystem: php-laravel
target_versions:
  laravel: ">=10.0"
dependencies:
  - core-universal-coding-standards
  - stacks-php-conventions
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Laravel Business Logic and Service Layers

## Inheritance
This document inherits from and extends the [Universal Coding Standards](../../core/05-universal-coding-standards.md). Refer to the universal standards for class encapsulation, single responsibility functions, and background queue designs. This page specifies only Laravel-specific logic patterns.

## Purpose
This document defines standards for separating business logic from framework boundaries, utilizing single-purpose Action classes, and managing background queue workers.

## Scope
Applies to service classes (`app/Services/`), Action modules (`app/Actions/`), Event listeners, and queued jobs.

---

## Directives

### 1. Action Classes Pattern
- **Standard**: Encapsulate transactional business logic in single-purpose Action classes. An Action class must expose exactly one public method (usually `execute()` or `__invoke()`).
- **Benefits**: Simplifies unit testing, enforces Single Responsibility Principle, and makes business actions highly reusable.
- **Code Syntax**:
  ```php
  namespace App\Actions;

  use App\Models\User;
  use App\Events\UserRegistered;

  class RegisterUserAction {
      public function execute(array $data): User {
          // Wrap state alterations in transactions
          return DB::transaction(function () use ($data) {
              $user = User::create([
                  'email' => $data['email'],
                  'name' => $data['name'],
                  'password' => bcrypt($data['password']),
              ]);

              event(new UserRegistered($user));

              return $user;
          });
      }
  }
  ```

### 2. Event-Driven Decoupling
- **Rule**: Do not perform secondary side-effects (e.g. sending welcome emails, dispatching tracking parameters to third-party endpoints) inside Action classes or Controllers.
- **Action**: Dispatch events (`UserRegistered`) and configure listeners to run these operations asynchronously in the background.

### 3. Background Job Queues
- **Payload Purity**: Pass only primary keys (identifiers) to Job constructor actions. Let the worker pull fresh data during execution:
  ```php
  class SendInvoiceEmail implements ShouldQueue {
      use Dispatchable, InteractsWithQueue, Queueable, SerializesModels;

      public function __construct(private int $invoiceId) {}

      public function handle(InvoiceRepository $repository): void {
          $invoice = $repository->findOrFail($this->invoiceId);
          // execute email dispatch
      }
  }
  ```

---

## Common Mistakes & Anti-Patterns
- **Fat Models**: Writing database queries, API logic, and notification triggers directly inside Eloquent models.
- **Synchronous API Delays**: Invoking slow external APIs (e.g., dispatching mailing templates via SendGrid) inside the HTTP thread instead of queuing a job.
- **State Leakage in Octane**: Preserving user-specific data in class variables within Singleton services, which leaks sessions to subsequent requests.

---

## References
- Simplicity over abstraction: [core/02-architecture-and-simplicity.md](file:///Users/kodexkode/Documents/workspace/promptengine/core/02-architecture-and-simplicity.md)
- Queue worker performance: [performance/03-caching-and-queues.md](file:///Users/kodexkode/Documents/workspace/promptengine/performance/03-caching-and-queues.md)
