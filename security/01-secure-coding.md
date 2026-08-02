---
document_id: security-secure-coding
title: Secure Coding Standards
ecosystem: cross-cutting
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Secure Coding Standards

## Purpose
This document defines input validation, output encoding, and SQL/XSS injection prevention guidelines to maintain code security across all platforms.

## Scope
Applies to backend request processors (Laravel, PHP core) and frontend templates (Vue 3, Flutter).

---

## Directives

### 1. Input Validation (The Trust Boundary)
- **Standard**: Treat all data coming from a web request, API payload, file upload, or user input as hostile and untrusted.
- **Rule**: Validate all inputs at the boundary (e.g. inside Laravel FormRequests or Dart input formatters) for:
  - Type (e.g., string, integer, boolean).
  - Bounds (e.g., character length ranges, minimum/maximum values).
  - Format (e.g., standard email regex, alphanumeric patterns).

### 2. Preventing SQL Injections
- **Rule**: Never concatenate variable parameters directly inside database queries.
- **Good**: Use framework query builders or PDO parameters that implement bound sql inputs:
  ```php
  // Safe Eloquent query binding
  User::where('email', $request->input('email'))->first();
  ```
- **Bad**:
  ```php
  // Vulnerable to injection
  DB::select("SELECT * FROM users WHERE email = '" . $request->input('email') . "'");
  ```

### 3. Cross-Site Scripting (XSS) Prevention
- **Rule**: Escape and encode variables before rendering them in views or templates.
- **Vue**: Use standard Mustache syntax `{{ value }}` which escapes values automatically. Never use `v-html` for user-generated strings unless it is sanitized using a verified library like DOMPurify.
- **Laravel**: Use standard Blade output braces `{{ $value }}` rather than unescaped `{!! $value !!}` tags.

---

## Common Mistakes & Anti-Patterns
- **Frontend-Only Validation**: Relying solely on client-side JS/Vue validation. Attackers can bypass the frontend client and post raw invalid payloads directly to your API routes.
- **Concatenated Queries**: Writing raw SQL fragments with concatenated variables to bypass Eloquent migration limitations.
- **Raw File Execution**: Allowing users to upload files with execute permissions (e.g., `.php`, `.js` files) to web-accessible local directories.

---

## References
- Manifest mapping: [playbook-manifest.json](file:///Users/kodexkode/Documents/workspace/promptengine/playbook-manifest.json)
- API structure validation: [core/03-data-and-api-modeling.md](file:///Users/kodexkode/Documents/workspace/promptengine/core/03-data-and-api-modeling.md)
