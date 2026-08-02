---
document_id: core-data-and-api-modeling
title: Database Schema Design and API Modeling
ecosystem: cross-cutting
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Database Schema Design and API Modeling

## Purpose
This document outlines standard conventions for database normalization, foreign keys, transaction handling, and API response structures across all technology stacks.

## Scope
Applies to database migrations, Eloquent schema setups, REST APIs, JSON serialization, and Dart/Vue data models.

---

## Database Design Standards

### 1. Relational Integrity
- **Always use foreign keys**: Never rely on application-level logic to preserve data integrity. Use database foreign key constraints (`ON DELETE RESTRICT` or `ON DELETE CASCADE` where appropriate).
- **Index foreign keys**: Every column used as a foreign key must be indexed to prevent slow queries on joins (refer to database index performance guidelines).
- **Enforce strict constraints**: Use `UNIQUE` constraints and `NOT NULL` columns aggressively. Do not allow nullable fields unless the null state represents a specific, valid business case.

### 2. Normalization vs. Denormalization Trade-offs
- **Default to 3NF**: Structure tables in Third Normal Form (3NF) to prevent write anomalies and duplicate data.
- **Denormalize only for read optimization**: Denormalize columns (such as cache sums, status counts) only when:
  1. Profile data confirms joins are causing slow response times under load.
  2. The denormalized column value is updated safely via atomic database triggers or database transactions.

---

## API Contract Standards

### 1. JSON Payload Conventions
- **Naming Case consistency**:
  - **API fields**: Use `snake_case` for all public JSON properties in REST endpoints.
  - **Reconciliation**: Convert JavaScript/Dart client models to `camelCase` internally, but serialize them to `snake_case` before transmitting payloads.
- **Null values**: Omit null fields from JSON payloads rather than sending `"property": null` unless explicit null resets are required.

### 2. Standard Envelope Structure
Every API endpoint must return a predictable structure.

#### Success Envelope
```json
{
  "data": {
    "id": 124,
    "email": "user@example.com",
    "status": "active"
  }
}
```

#### Error Envelope
```json
{
  "error": {
    "code": "validation_failed",
    "message": "The provided password is too weak.",
    "details": [
      {
        "field": "password",
        "message": "Must be at least 12 characters long."
      }
    ]
  }
}
```

---

## API Design Best Practices
- **Use standard HTTP status codes**:
  - `200 OK` for successful fetches or updates.
  - `201 Created` for resource generation.
  - `400 Bad Request` for invalid operations.
  - `401 Unauthorized` for missing/invalid auth tokens.
  - `403 Forbidden` for authenticated users trying to access unauthorized resources.
  - `422 Unprocessable Entity` for field validations failures.
- **Avoid exposing primary increments directly**: In public URLs or APIs, prefer using UUIDs, ULIDs, or random hashed tokens instead of exposure of raw incremental IDs (`/api/v1/users/1` $\rightarrow$ `/api/v1/users/01h8x12bf...`).

---

## Common Mistakes & Anti-Patterns
- **Polymorphic Database Relations**: Creating generic tables that hook to multiple parent models using a text string column (`commentable_type`, `commentable_id`) without physical foreign keys.
- **Schema-less Blobs**: Storing complex relational data in unstructured JSON columns when it should be normalized in separate tables.
- **Fat Responses**: Returning database models directly in API responses without a transformer layer (like Laravel API Resources), exposing internal fields like password hashes or secret tokens.

---

## References
- Caching API endpoints: [performance/03-caching-and-queues.md](file:///Users/kodexkode/Documents/workspace/promptengine/performance/README.md)
- Laravel data layer implementation: [stacks/php-laravel/laravel-data.md](file:///Users/kodexkode/Documents/workspace/promptengine/stacks/php-laravel/laravel-data.md)
