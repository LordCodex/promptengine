---
document_id: performance-db-optimization
title: Database Optimization and Indexing
ecosystem: cross-cutting
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Database Optimization and Indexing

## Purpose
This document defines index strategies and query design standards to prevent slow database queries and maintain short locks under high concurrency.

## Scope
Applies to database schema migrations and SQL query constructions.

---

## Directives

### 1. Indexing Strategies
- **Foreign Key Indexing**: Every column containing a foreign key reference must be indexed. Databases do not automatically index foreign keys, leading to full-table scans on join operations.
- **Compound Indexes**: When queries filter by multiple fields (e.g., `WHERE status = 'active' AND created_at > ?`), create a compound index. Order compound index columns from lowest cardinality (most general, e.g., `status`) to highest cardinality (most specific, e.g., `created_at`).
- **Index Overhead**: Indexes speed up reads but slow down writes. Avoid adding indexes to columns that are rarely queried.

### 2. Query Analysis (EXPLAIN)
- **Standard**: Run the `EXPLAIN` or `EXPLAIN ANALYZE` command against any SQL query running in a high-traffic endpoint.
- **Goal**: Verify that the query executes an `Index Scan` rather than a `Seq Scan` (Sequence Scan / Full Table Scan). Ensure the database uses the intended index keys.

### 3. Avoiding the N+1 Query Trap
- **Rule**: Never run database queries inside loops.
- **Bad (N+1 Queries)**:
  ```php
  $books = Book::all(); // 1 query
  foreach ($books as $book) {
      echo $book->author->name; // N queries
  }
  ```
- **Good (Eager Loading)**:
  ```php
  $books = Book::with('author')->get(); // 2 queries total
  foreach ($books as $book) {
      echo $book->author->name;
  }
  ```

---

## Common Mistakes & Anti-Patterns
- **The "Select *" Abuse**: Querying all columns (`SELECT *`) from large tables when only one column is needed. This inflates memory usage and database serialization overhead.
- **Ignoring Database Wildcards**: Querying strings using leading wildcards (`LIKE '%search_term'`). This prevents the database database engine from traversing indexes, forcing full-table scans.
- **Massive Offset Queries**: Using high page offsets (`OFFSET 100000 LIMIT 10`) for pagination. This forces the database to read and discard 100,000 records before returning the final 10. Use Cursor-Based Pagination instead.

---

## References
- Relational schema structure: [core/03-data-and-api-modeling.md](file:///Users/kodexkode/Documents/workspace/promptengine/core/03-data-and-api-modeling.md)
- Laravel Eloquent configuration: [stacks/php-laravel/laravel-data.md](file:///Users/kodexkode/Documents/workspace/promptengine/stacks/php-laravel/laravel-data.md)
