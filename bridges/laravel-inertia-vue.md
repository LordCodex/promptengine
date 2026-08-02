---
document_id: bridge-laravel-inertia-vue
title: Laravel Inertia Vue Integration Bridge
ecosystem: cross-cutting
target_versions:
  laravel: ">=10.0"
  vue: "^3.3"
  inertia: "^1.0"
dependencies:
  - stacks-laravel-routing
  - stacks-vue-components
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Laravel Inertia Vue Integration Bridge

## Purpose
This document defines standards for sharing data, handling form state validations, and managing client-server routing using Inertia.js as a bridge between Laravel and Vue 3.

## Scope
Applies to Laravel HandleInertiaRequests Middleware and Vue Inertia page components.

---

## Directives

### 1. Unified Share Middleware
- **Rule**: All shared global data (e.g. authenticated user profile, flash notifications, active configuration flags) must be declared in Laravel's `HandleInertiaRequests` middleware.
- **Why**: Keeps sharing logic centralized and ensures properties are hydrated on every Inertia page transition without manual controller injections.
- **Code Syntax**:
  ```php
  // app/Http/Middleware/HandleInertiaRequests.php
  public function share(Request $request): array {
      return array_merge(parent::share($request), [
          'auth' => [
              'user' => $request->user() ? [
                  'id' => $request->user()->id,
                  'name' => $request->user()->name,
                  'role' => $request->user()->role,
              ] : null,
          ],
          'flash' => [
              'success' => fn () => $request->session()->get('success'),
              'error' => fn () => $request->session()->get('error'),
          ],
      ]);
  }
  ```

### 2. Vue Form Submission Lifecycle
- **Rule**: Use Inertia's `useForm` helper inside Vue pages to manage form data, track validation states, and submit requests.
- **Why**: Handles loading states, automatically intercepts and displays Laravel validation validation messages, and prevents page refresh.
- **Code Syntax**:
  ```html
  <script setup lang="ts">
  import { useForm } from '@inertiajs/vue3';

  const form = useForm({
    email: '',
    name: '',
  });

  function submit() {
    form.post(route('users.store'), {
      onSuccess: () => form.reset(),
      onError: (errors) => console.log('Validation failed', errors),
    });
  }
  </script>
  ```

### 3. Preserving Page State (Preventing Reloads)
- **Standard**: Always use Inertia's `<Link>` component for routing inside Vue components. Never use raw anchor `<a>` tags.
- **Why**: Standard `<a>` tags trigger full browser page reloads, destroying the client-side state in memory and reloading the entire bundle.

---

## Common Mistakes & Anti-Patterns
- **Massive Share Payloads**: Sharing huge model relations (like a user's entire purchase history) inside the `HandleInertiaRequests` middleware on every page request. Keep shared props minimal.
- **Duplicate Router Hooks**: Combining Vue Router inside an Inertia application. Inertia handles client routing natively using Laravel route mappings.
- **Direct Controller HTML Returns**: Returning a standard `view('users.index')` instead of `Inertia::render('Users/Index')` for Inertia-targeted routes.

---

## References
- Vue components conventions: [vue-components.md](file:///Users/kodexkode/Documents/workspace/promptengine/stacks/js-ts-vue-nuxt/vue-components.md)
- Laravel route validation: [laravel-routing.md](file:///Users/kodexkode/Documents/workspace/promptengine/stacks/php-laravel/laravel-routing.md)
