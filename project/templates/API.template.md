# API Documentation

## Playbook Metadata
- **Purpose**: Authoritative reference template defining the project's external and internal API contracts, authentication lifecycles, error payloads, and protocol architectures.
- **Scope**: Reusable for REST, GraphQL, gRPC, WebSockets, Server-Sent Events (SSE), and webhook pipelines.
- **When to Read**: Prior to implementing new endpoints, modifying request parameters, or refactoring authorization filters.
- **Related Playbooks**: [Project Overview](../README.md), [Project Documentation Standard](../02-project-documentation-standard.md), [API Engineering Standard](../../core/07-api-engineering-standard.md).
- **Version**: 1.0.0
- **Last Reviewed**: 2026-08-03

---

## Document Metadata
- **Project Name**: [Enter Project Name]
- **Document Version**: 1.0.0
- **Status**: [Draft / In Review / Approved]
- **Owner**: [Enter API Lead / Tech Owner Role]
- **API Version**: [e.g. v1 / v2]
- **Last Updated**: [YYYY-MM-DD]
- **Reviewers**: [Enter Reviewers]
- **Related Documents**: [PRD.md](PRD.md) | [Architecture.md](Architecture.md) | [Database.md](Database.md) | [BusinessRules.md](BusinessRules.md)

---

## 1. API Overview
- **Overview**: [Provide a high-level description of the API scope and system footprint.]
- **Primary Consumers**: [Identify target integration consumers, e.g. web frontend, iOS client, third-party partners.]
- **Supported Protocols**: [e.g., REST (JSON), GraphQL, WebSockets]
- **Primary Use Cases**: [e.g. customer billing ingestion, real-time status telemetry.]

---

## 2. API Architecture
Define the structural protocols and pipeline targets:
- **REST Endpoints**: [Identify REST base routing policies, e.g., JSON payloads over HTTPS.]
- **GraphQL Schema**: [Outline endpoints query layout, complexity limits.]
- **WebSockets / Event Stream**: [Highlight real-time channels footprint.]
- **RPC / gRPC Services**: [Identify internal service mesh contracts.]
- **Outgoing Webhooks**: [Outline callback event triggers.]

---

## 3. Versioning Strategy
- **Version Format**: [e.g. URI-based versioning (api/v1/...), header-based negotiation.]
- **Breaking Changes Policy**: [Define what constitutes a breaking change, e.g. deleting columns, modifying response keys.]
- **Deprecation & Sunset Schedule**: [e.g. Deprecated APIs remain supported for 6 months; Sunset alerts dispatch header notices.]
- **Backward Compatibility guarantees**: [e.g. Optional inputs must not be changed to required status.]

---

## 4. Authentication
- **Supported Authentication Methods**:
  - [e.g. Bearer stateless JWT tokens, OAuth2 flows, SHA-256 HMAC client signatures, Signed URLs.]
- **Token Lifecycle**: [e.g. Access token TTL: 15 minutes. Refresh token TTL: 14 days.]
- **Token Refresh & Revocation Strategy**: [Explain token invalidation policies.]

---

## 5. Authorization & Multi-Tenancy
- **Ownership Constraints**: [e.g. Consumers can only read/write resources linked to their own account UUID.]
- **Role-Based Access Control (RBAC)**: [e.g. Scope lists required per endpoint class, e.g. read:orders, write:orders.]
- **Multi-Tenant Routing Policies**: [e.g. Isolation rules preventing cross-tenant leakage.]

---

## 6. Base URLs

- **Development**: `http://localhost:8000/api/v1`
- **Testing**: `https://test.api.domain.com/v1`
- **Staging**: `https://staging.api.domain.com/v1`
- **Production**: `https://api.domain.com/v1`
- **Regional / Private Endpoints**: [Define internal routing setups]

---

## 7. Common Headers

| Header Name | Required | Type | Purpose |
| :--- | :--- | :--- | :--- |
| **Authorization** | Yes (auth routes) | String | Bearer token authorization |
| **X-Correlation-ID** | Yes | UUID | End-to-end trace correlation tracking |
| **X-Idempotency-Key** | Yes (mutations) | String | Prevents double-processing write payloads |
| **Content-Type** | Yes | String | Must be `application/json` |

---

## 8. Request & Response Standards
- **Input Parameters**: [Naming style rules, e.g. query parameters in snake_case, paths in snake_case.]
- **Pagination Heuristics**: [e.g., Default cursor pagination for bulk streams. Standard query limit bounds: max 100 records.]
- **Sorting & Filtering**: [e.g. `?sort=-created_at&filter[status]=active` standard layouts.]
- **File Upload Policies**: [e.g. Multipart form uploads must restrict sizes and pass magic byte validations.]

---

## 9. Standard Error Format (RFC 7807)

Every validation failure or application exception must return a JSON response complying with **RFC 7807 Problem Details**:

```json
{
  "type": "https://api.domain.com/errors/validation-failed",
  "title": "Validation Failed",
  "status": 422,
  "detail": "The request body failed parameter constraints.",
  "instance": "/api/v1/orders",
  "error_code": "VAL_REQ_005",
  "validation_errors": {
    "email": [
      "The email field must be a valid email address."
    ]
  },
  "correlation_id": "8f9a2e6b-4c5d-4f3a-2b1c-0e9f8a7b6c5d"
}
```

---

## 10. Endpoint Reference Template

Create a section for each endpoint in the API:

### Endpoint: `[Endpoint ID: e.g. API-ORD-001]`
- **Name**: [e.g. Create Order]
- **Method & Path**: `[POST /api/v1/orders]`
- **Purpose**: [Detailed business description of what the endpoint does.]
- **Authentication**: [Yes/No] (type: `[Bearer JWT / API Key]`)
- **Authorization Scopes**: `[write:orders]`
- **Request Parameters (Query/Path)**:
  - `[parameter]` (type: `[type]`, e.g. String, Required): [Description]
- **Request Body Payload**:
  ```json
  {
    "user_id": "UUID",
    "items": []
  }
  ```
- **Validation Rules**:
  - `[field_name]`: [e.g. Required, UUID, exists in users table.]
- **Success Response (201 Created)**:
  ```json
  {
    "order_id": "UUID",
    "status": "pending"
  }
  ```
- **Error Responses**:
  - **401 Unauthorized**: Token is expired or missing.
  - **422 Unprocessable**: Validation failure (RFC 7807 payload).
- **Business Rules Linked**: Refer to rule [BR-BILL-001](BusinessRules.md#BR-BILL-001) for payment validation rules.
- **Side Effects**: Dispatches order activation emails, hooks Stripe checkout API.
- **Idempotency & Rate Limits**:
  - Idempotency: Required via `X-Idempotency-Key` header.
  - Rate Limits: 60 requests per minute per IP address.
- **Related Endpoints**: `[GET /api/v1/orders/{id}]`

---

## 11. GraphQL Specifications
- **Schema Boundary**: [Mermaid or text schema layout maps.]
- **Query Complexity Rules**: [e.g., Maximum query depth allowed is 5 to prevent DOS loops.]
- **Subscriptions**: [Websocket telemetry parameters.]

---

## 12. WebSockets & Real-Time Specifications
- **Heartbeat & Pings**: [e.g., Pings dispatched every 30 seconds; connection drops if 2 miss.]
- **Events Map**:
  - **Client $\rightarrow$ Server Events**: `[join_room]`, `[leave_room]`.
  - **Server $\rightarrow$ Client Events**: `[order_status_updated]`.

---

## 13. Outgoing Webhooks & Callbacks
- **Webhook Events List**: `[order.created]`, `[payment.failed]`.
- **Signature Verification**: [e.g. HMAC SHA256 signatures dispatched in the `X-Webhook-Signature` header.]
- **Retry Schedule**: [e.g. Exponential backoff retry scheme: retry up to 5 times over 24 hours.]

---

## 14. Performance Optimization
- **Caching**: [e.g. Catalog query views utilize ETags; cache headers: `Cache-Control: public, max-age=3600`.]
- **Compression**: [e.g. Gzip compression enabled for payload responses exceeding 10KB.]

---

## 15. Security Hardening
- **CORS Policies**: [e.g. REST calls restricted to domains: `*.domain.com`.]
- **CSRF Mitigations**: [e.g. Required state verification tokens for Web session endpoints.]
- **PII Hardening**: [e.g. Under no circumstances must passwords or auth tokens be recorded in request trace logs.]

---

## 16. Observability & SRE Metrics
- **Performance Targets**: [e.g. 95% of reads must return within 150ms; 95% of writes must commit within 300ms.]
- **Audit Traces**: [e.g. All API mutations must generate database audit logs tracking consumer UUIDs.]

---

## 17. Related Documents
- **PRD**: [PRD.md](PRD.md)
- **Architecture**: [Architecture.md](Architecture.md)
- **Database Schema**: [Database.md](Database.md)
- **ADR Logs**: [Decisions/README.md](Decisions/README.md)

---

## AI Guidance

When reading or updating API documentation, follow these rules:
- **Never Guess Endpoints**: Do not create paths, query parameters, or payload fields. Verify them in routes files or schemas.
- **Strict Backward Compatibility**: Never suggest breaking shifts in parameters or keys without explicit developer approval.
- **Align with Security playbooks**: Ensure rate limits, headers, and CORS structures comply with [API Security Standard](../../core/08-security-engineering-standard.md).
- **Maintain Schema Accuracy**: Verify that REST request payloads align with database schema definitions in [Database.md](Database.md).

---

## Developer Guidance

- **Keep Contracts Synced**: Always update the API markdown document in the same commit that adds or edits code routes.
- **Document Validation Rules**: Detail validation parameters (required states, bounds) clearly in this file so frontend clients avoid validation failures.
- **Prefer Standard Error Specs**: Enforce consistency across all controllers; reject PRs returning custom exceptions outside the RFC 7807 payload.
