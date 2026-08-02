---
document_id: stacks-nextjs-api-routes
title: Next.js API Routes, Route Handlers, and Server Actions
ecosystem: react-next
dependencies:
  - core-universal-coding-standards
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Next.js API Routes, Route Handlers, and Server Actions

## Inheritance & Constraints
This document inherits from the [Universal Coding Standards](../../core/05-universal-coding-standards.md). It outlines standards for Route Handlers and Server Actions, avoiding general tutorials.

---

## 1. Route Handlers (`route.ts`)

Use `route.ts` files to create API endpoints:
- **HTTP Methods**: Export distinct verb handlers: `export async function GET(request: Request) {}`.
- **Payload Verification**: Always validate incoming query parameters and JSON request payloads against type validation schemas (e.g. Zod) before logic execution.
- **Unified JSON Response**: Ensure success and error channels return standard JSON objects with appropriate HTTP status codes:
  - `NextResponse.json({ error: 'Unauthenticated' }, { status: 401 })`

---

## 2. Server Actions

Use Server Actions (`'use server'`) for secure form actions, database queries, and client-server mutations:
- **Explicit Declaration**: Declare `'use server'` at the absolute top of the server action file or async function.
- **Input Validation**: **Never trust parameter structures sent from Client Components.** Re-validate all arguments inside the Server Action block.
- **State Feedback**: Use hooks like `useActionState` (or the experimental `useFormState`) on the client side to track pending states, results, and validation messages.

---

## Review Checklist

Verify backend endpoints against this checklist:
- [ ] Inherits correctly from Universal Coding Standards.
- [ ] Route Handler outputs return standard `NextResponse.json()` packages.
- [ ] Server Actions re-validate input parameters server-side.
- [ ] Error statuses reflect validation rules.
