---
document_id: stacks-laravel-routing
title: Laravel Routing and Request Validation
ecosystem: php-laravel
target_versions:
  laravel: ">=10.0"
dependencies:
  - core-universal-coding-standards
  - stacks-php-conventions
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Laravel Routing and Request Validation

## Inheritance
This document inherits from and extends the [Universal Coding Standards](../../core/05-universal-coding-standards.md). Refer to the universal standards for core naming, function layout, and validation principles. This page specifies only Laravel-specific routing patterns.

## Purpose
This document defines standards for routing architectures, API middlewares, FormRequest validation parameters, and controller designs.

## Scope
Applies to routing files (`routes/web.php`, `routes/api.php`), HTTP Controllers, and FormRequest classes.

---

## Directives

### 1. Route Definitions
- **Group Routes Logically**: Group related routes using route prefixes, namespaces, and middleware bindings rather than declaring them flat.
- **Use Named Routes**: Always attach name tags to routes to prevent broken link patterns when endpoints get refactored:
  ```php
  Route::get('/users/{user}', [UserController::class, 'show'])->name('users.show');
  ```

### 2. Single Responsibility Controllers
- **Rule**: Keep Controllers thin. Controllers should only handle HTTP concerns: authentication checks, request validation triggers, calling business services (Actions), and returning views or JSON payloads.
- **Action Call Pattern**:
  ```php
  public function store(StoreUserRequest $request, RegisterUserAction $registerUser): JsonResponse {
      $user = $registerUser->execute($request->validated());

      return response()->json(new UserResource($user), 201);
  }
  ```

### 3. Dedicated FormRequest Validation
- **Rule**: Never run validation checks directly inside Controller actions. Always instantiate a dedicated FormRequest class to capture validation rules.
- **Why**: Keeps controller actions clean, simplifies unit testing, and ensures validation completes before the controller action boots.
- **Code Syntax**:
  ```php
  namespace App\Http\Requests;

  use Illuminate\Foundation\Http\FormRequest;

  class StoreUserRequest extends FormRequest {
      public function authorize(): bool {
          return $this->user()->can('create', User::class);
      }

      public function rules(): array {
          return [
              'email' => ['required', 'email', 'unique:users,email'],
              'name' => ['required', 'string', 'max:255'],
          ];
      }
  }
  ```

---

## Common Mistakes & Anti-Patterns
- **Fat Controllers**: Writing complex business calculations, SQL queries, or email dispatches directly inside controller actions.
- **Inline Validation**: Running `$request->validate([...])` inside the controller body.
- **Global Auth Checking inside Controllers**: Putting manual user check flags (`if (!auth()->user())`) inside controller actions instead of routing the endpoint through the standard `auth` middleware.

---

## References
- Safe API modeling: [core/03-data-and-api-modeling.md](file:///Users/kodexkode/Documents/workspace/promptengine/core/03-data-and-api-modeling.md)
- Rate limiting rules: [security/04-api-and-infra-security.md](file:///Users/kodexkode/Documents/workspace/promptengine/security/04-api-and-infra-security.md)
