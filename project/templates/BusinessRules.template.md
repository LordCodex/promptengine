# Business Rules Document

## Playbook Metadata
- **Purpose**: Authoritative template capturing the business logic, financial rules, state transitions, and validation rules that govern system behavior.
- **Scope**: Reusable across any software project domain (e.g. SaaS, mobile, e-commerce, financial tools) regardless of technology stack.
- **When to Read**: Prior to writing feature code, designing database schemas, or executing business validation updates.
- **Related Playbooks**: [Project Overview](../README.md), [Project Documentation Standard](../02-project-documentation-standard.md).
- **Version**: 1.0.0
- **Last Reviewed**: 2026-08-03

---

## Document Metadata
- **Project Name**: [Enter Project Name]
- **Document Version**: 1.0.0
- **Status**: [Draft / In Review / Approved]
- **Owner**: [Enter Product/Business Owner Name]
- **Reviewers**: [Enter Reviewers]
- **Last Updated**: [YYYY-MM-DD]
- **Related Documents**: [PRD.md](PRD.md) | [Architecture.md](Architecture.md) | [Database.md](Database.md) | [API.md](API.md)

---

## 1. Business Overview
- **Business Domain**: [Provide a brief explanation of the business domain, e.g. B2B SaaS logistics, marketplace retail.]
- **Primary Objectives**: [Summarize the main business drivers.]
- **Actors Vocabulary**: [Identify key user groups and system actors.]

---

## 2. Core Principles
- [High-level guidelines that represent invariant business constraints. Use placeholders.]
- **Principle 1**: [e.g., A payment transaction must never be executed more than once.]
- **Principle 2**: [e.g., Refunds are irreversible and must trigger financial reconciliation.]
- **Principle 3**: [e.g., Product stock counts cannot be updated below zero in active catalogs.]

---

## 3. Actors
Document the characteristics, permissions, and boundaries of each system participant:

### Actor A: [Actor Name / Role]
- **Description**: [e.g., Registered Customer]
- **Responsibilities**: [Key actions expected, e.g. places order, initiates payment]
- **Permissions**: [e.g., Can view own invoice history]
- **Restrictions**: [e.g., Cannot modify product prices or coupon values]
- **Relationships**: [e.g., Belongs to a single billing account]

---

## 4. Business Modules
Rules are organized logically by feature domain. Use placeholders:

### Module A: [Module Name / Domain]
- **Overview**: [Summary of the module's business purpose.]

#### Rule ID: BR-MODA-001: [Rule Title]
- **Title**: [e.g., Coupon Expiry Valuation]
- **Description**: [Detailed description of the rule behavior.]
- **Reason**: [Explain the business rationale behind the rule.]
- **Affected Modules**: [e.g., Billing, Checkout, Notifications]
- **Preconditions**: [State of the system required before validation, e.g. user is checking out]
- **Business Logic**:
  - [Write the exact logical flow, e.g., "If current date > coupon.expiration_date, discount = 0."]
- **Exceptions**: [e.g., Extended coupons granted by administrators override standard expiry.]
- **Failure Behavior**: [e.g., Reject checkout, return coupon invalid notification.]
- **Success Behavior**: [e.g., Deduct discount amount from order subtotal.]
- **Dependencies**: [e.g., Rule BR-MODB-005]
- **Priority**: [Critical / Secondary]
- **References**: [e.g., PRD Section 4.2]
- **Notes**: [Additional context]

---

## 5. State Transitions
Define application entity lifecycles and workflow rules:

```text
[State A] ──> [State B] ──> [State C] ──> [State D]
```

### Aggregate Entity: [Entity Name]
- **Transitions Map**:
  - **Draft** $\rightarrow$ **Pending**: Triggered when the user submits data for approval.
  - **Pending** $\rightarrow$ **Approved**: Triggered when an administrator logs verification.
  - **Approved** $\rightarrow$ **Completed**: Triggered when the transaction commits.
  - **Any State** $\rightarrow$ **Archived**: Safe retention workflow rules.

---

## 6. Financial Rules
- **Pricing Models**: [e.g., Flat fee tiering, volume-based pricing limits.]
- **Discounts & Coupons**: [e.g., Rounding policies when applying percentage discounts.]
- **Taxes**: [e.g., Regional tax brackets based on checkout shipping destination.]
- **Refunds & Commissions**: [e.g., Refund window limits, calculation of partner payout cuts.]
- **Wallets & Invoicing**: [e.g., Invoice generation deadlines, maximum wallet balances.]
- **Rounding & Currency**: [e.g., Round up financial transactions to the nearest two decimal places.]

---

## 7. Authorization Rules
- **Ownership Policies**: [e.g., Users can only view or modify entities created by their own organization.]
- **Approval Workflows**: [e.g., Single purchases over $10,000 require dual admin approvals.]
- **Visibility Gates**: [e.g., Staff members can only view anonymized customer profiles.]

---

## 8. Validation Rules
- **Domain Eligibility**: [e.g., Minimum age requirements to create accounts.]
- **Quantity Limits**: [e.g., Maximum items allowed per single transaction.]
- **Business Deadlines**: [e.g., Cancellations are only permitted up to 24 hours before service starts.]

---

## 9. Notifications & Communications
- **Notification Triggers**: [e.g., Dispatch account activation alert when signup completes.]
- **Delivery Channels**: [e.g., SMS, Email, In-app push notifications.]
- **Timing & Retries**: [e.g., Queue retry attempts up to 3 times for SMS dispatch failures.]

---

## 10. Integrations & Third-Party Expectations
- **Partner Expectations**: [e.g., Third-party credit checks must return response in under 5 seconds.]
- **Failure Strategies**: [e.g., If background check provider is offline, fall back to manual queue review.]
- **Ownership Boundaries**: [Identify which module manages integration API keys.]

---

## 11. Reporting & Analytics Rules
- **Business Metrics**: [Define the KPIs and metrics calculations.]
- **Aggregation Heuristics**: [e.g., Group transaction totals daily by UTC timestamp.]
- **Retention Rules**: [e.g., Monthly performance logs are archived after 3 years.]

---

## 12. Regulatory & Compliance Rules
- **Compliance Standards**: [e.g., Must support user data deletion requests under GDPR.]
- **Audit Trails**: [e.g., Audit logs must permanently record user log-ins and parameter edits.]

---

## 13. System Exceptions
- **Legacy Workarounds**: [Document temporary rules supporting backward-compatibility migrations.]
- **Feature Flags**: [Rules that alter behavior dynamically based on environment configs.]

---

## 14. Glossary & Vocabulary

- **Term A**: [Clear business definition, synonyms, and context.]
- **Term B**: [Clear business definition, synonyms, and context.]

---

## 15. Related Documents
- **PRD**: [PRD.md](PRD.md)
- **Architecture**: [Architecture.md](Architecture.md)
- **Database Schema**: [Database.md](Database.md)
- **API Contracts**: [API.md](API.md)

---

## AI Guidance

When reading or updating business rules, follow these instructions:
- **Never Invent Logic**: Do not guess calculation rates, age limits, or status workflows. If logic is missing, ask.
- **Isolate Adjustments**: When editing business rules, update only the target modules. Do not rewrite unrelated rules.
- **Priority over Tech**: Business rules always override implementation convenience. Do not suggest compromises in business logic to suit framework limitations.
- **Distinguish Inferences**: Clearly mark any inferred logic or assumptions using the `[ASSUMPTION]` prefix.

---

## Developer Guidance

- **Keep Identifiers Stable**: Ensure rule IDs (e.g. `BR-MODA-001`) remain stable so they can be referenced inside tests and code comments.
- **Document Changes**: Record intentional adjustments in the document version metadata history.
- **No Implementation Duplication**: Focus strictly on business constraints; do not describe SQL columns, classes, or code variables here.
