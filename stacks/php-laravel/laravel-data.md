---
document_id: stacks-laravel-data
title: Laravel Eloquent and Database Migrations
ecosystem: php-laravel
target_versions:
  laravel: ">=10.0"
dependencies:
  - core-universal-coding-standards
  - stacks-php-conventions
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Laravel Eloquent and Database Migrations

## Inheritance
This document inherits from and extends the [Universal Coding Standards](../../core/05-universal-coding-standards.md). Refer to the universal standards for core naming, database modeling, and integrity principles. This page specifies only Laravel-specific database patterns.

## Purpose
This document defines standards for database tables configuration, Eloquent model relationships, model scopes, and data seeders.

## Scope
Applies to database migrations (`database/migrations/`), Eloquent models (`app/Models/`), and database seeders (`database/seeders/`).

---

## Directives

### 1. Migration Best Practices
- **Define Explicit Foreign Key Actions**: Always link relations explicitly with correct cascading logic:
  ```php
  Schema::create('orders', function (Blueprint $table) {
      $table->id();
      $table->foreignId('user_id')->constrained()->onDelete('restrict');
      $table->unsignedInteger('amount_cents');
      $table->timestamps();
  });
  ```
- **Never Modify Migrations**: Once a migration is pushed or deployed to staging/production, never modify it. Create a new migration file to alter schema properties (refer to [legacy/02-backward-compatibility.md](file:///Users/kodexkode/Documents/workspace/promptengine/legacy/02-backward-compatibility.md)).

### 2. Eloquent Model Standards
- **Enforce Strict Attributes**: Add type casts to model properties explicitly to ensure values returned from database queries match system types:
  ```php
  protected $casts = [
      'amount_cents' => 'integer',
      'is_active' => 'boolean',
      'activated_at' => 'datetime',
  ];
  ```
- **Use Local Scopes**: Encapsulate database query filtering logic within reusable model scopes rather than writing custom `where` chains in service layers:
  ```php
  // Model Local Scope
  public function scopeActive(Builder $query): Builder {
      return $query->where('is_active', true);
  }
  ```

### 3. Model Relationships
- **Specify Keys Explicitly**: When defining complex relationships (e.g. `hasManyThrough`, `belongsToMany`), write out the key references explicitly to prevent auto-naming guess bugs during schema refactoring.
- **Set Up Eager Loading**: Prevent N+1 queries by defining `$with` configurations on models or utilizing `with()` on query calls (refer to [performance/02-database-optimization.md](file:///Users/kodexkode/Documents/workspace/promptengine/performance/02-database-optimization.md)).

---

## Common Mistakes & Anti-Patterns
- **Polymorphic Eloquent Models**: Creating polymorphic relations without database triggers or reference integrity checks.
- **Unencrypted Attributes**: Storing sensitive parameters (e.g., access tokens, user details) in plain database text instead of using Eloquent's `$casts = ['property' => 'encrypted']`.
- **Logic in Eloquent Models**: Writing complex business calculations or external API calls inside model event listeners (`creating`, `updating`). Keep models focused purely on data structures.

---

## References
- Database indexes: [performance/02-database-optimization.md](file:///Users/kodexkode/Documents/workspace/promptengine/performance/02-database-optimization.md)
- Safe database upgrades: [legacy/02-backward-compatibility.md](file:///Users/kodexkode/Documents/workspace/promptengine/legacy/02-backward-compatibility.md)
