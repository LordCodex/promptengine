# Troubleshooting & Diagnostics Guide

## Playbook Metadata
- **Purpose**: Authoritative reference template defining the project's diagnostic workflows, environment status verification checklists, common error indices, and incident lessons logs.
- **Scope**: Reusable for cloud, mobile, database, API, and frontend applications.
- **When to Read**: Immediately upon observing system anomalies, during active incidents, or when documenting resolved bug fixes.
- **Related Playbooks**: [Project Overview](../README.md), [Project Documentation Standard](../02-project-documentation-standard.md), [Observability Standard](../../core/29-observability-and-operational-excellence-standard.md).
- **Version**: 1.0.0
- **Last Reviewed**: 2026-08-03

---

## Document Metadata
- **Project Name**: [Enter Project Name]
- **Document Version**: 1.0.0
- **Status**: [Active / Under Revision]
- **Owner**: [Enter Tech Lead / Lead Support Engineer Role]
- **Reviewers**: [Enter Reviewers]
- **Last Updated**: [YYYY-MM-DD]
- **Related Documents**: [PRD.md](PRD.md) | [Architecture.md](Architecture.md) | [Database.md](Database.md) | [API.md](API.md) | [Deployment.md](Deployment.md)

---

## 1. How to Use This Guide
- **Search First**: Search for issue symptoms, keywords, or error codes in the Issue Index table before debugging.
- **Verify Environment**: Complete the Environment Verification checklist to rule out deployment or connection issues.
- **Follow Order**: Execute diagnostic verification steps in chronological order.
- **Fact-Based Fixes**: Confirm the root cause before applying resolutions. Document only verified solutions; avoid speculation.

---

## 2. Environment Verification

Prior to debugging application code, verify baseline system status:
- [ ] **Application Version**: Confirm the running container version matches build release tags.
- [ ] **Active Environment**: Identify target system boundaries (Dev, Staging, Production).
- [ ] **Runtime Configuration**: Check for missing environment variables or configuration maps.
- [ ] **Feature Flags**: Verify flags configuration state.
- [ ] **Database Connectivity**: Confirm ping responses from database instances.
- [ ] **Cache State**: Check Redis/Memcached availability and memory footprint.
- [ ] **Queue Status**: Check for queue backlog size spikes.
- [ ] **External Services API**: Confirm regional latency and auth keys status on partner services (e.g. Stripe).
- [ ] **Network Connectivity**: Validate VPC routing rules and public-facing DNS.
- [ ] **Infrastructure Health**: Verify CPU, memory, and disk usage are within safe limits.

---

## 3. Issue Index

| Issue ID | Symptom / Error | Category | Severity | Status | Last Verified | Related Docs |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| **TSG-001** | [e.g. API returns 504 gateway timeout] | [e.g. API] | [Critical] | [Resolved] | [YYYY-MM-DD] | [API.md](API.md) |

---

## 4. Troubleshooting Entry Template

Copy and fill this template for each verified issue:

---

### Issue ID: `[TSG-XXX: e.g. TSG-001]`

#### Title: `[Short descriptive title detailing symptom]`
- **Category**: [Application / Database / API / Authentication / Authorization / Infrastructure / Performance / Deployment / UI / Mobile / Security / Networking]
- **Severity**: [Critical / High / Medium / Low]
- **Status**: [Investigating / Active Workaround / Resolved]
- **Last Verified**: [YYYY-MM-DD]

#### Symptoms
- **Observable Behavior**: [What the user or system experiences.]
- **Error Messages**: `[Copy raw error log details]`
- **Logs**:
  ```text
  [Copy log trace lines]
  ```
- **Screenshots / Visual Evidence**: [Placeholder for screenshots if applicable]

#### Impact Scope
- **Affected Components**: [List modules, services, database tables, or SQS queues.]
- **Environment**: [e.g. Production and Staging]
- **Affected Users**: [e.g. Guest users checkout path.]

#### Root Cause Analysis
- **Verified Explanation**: [What caused the issue.]
- *(Tag clearly if details represent a confirmed root cause or a working hypothesis)*

#### Diagnostic Steps
1. **Verification**: [e.g. Trigger request manually via Postman/cURL.]
2. **Log Inspection**: [e.g. Search Datadog logs for specific transaction UUID.]
3. **Metrics Analysis**: [e.g. Read database connection counts charts.]
4. **Configuration Check**: [e.g. Verify environment credentials in secrets manager.]

#### Resolution Execution
[Detail step-by-step instructions to fix the issue, e.g. environment variable update, service restart, database index build.]

#### Recovery Verification
- **Verification Commands**: [e.g. run health test command or integration script.]
- **Expected Outcomes**: [e.g. Response returns 200 OK within 150ms.]
- **Monitoring Indicators**: [e.g. Check Sentry error charts for drop in exception count.]

#### Prevention & Automation
- **Prevention Recommendations**: [How to avoid this issue in the future.]
- **Monitoring Improvements**: [e.g. Add alert threshold trigger on database connection spikes.]
- **Testing Improvements**: [e.g. Add regression test covering this checkout edge case.]

#### Related Issues
- **Similar Patterns**: Cross-reference similar issues: `[e.g., TSG-005]`.

#### References & Links
- **PRD**: [PRD.md](PRD.md#section)
- **Architecture**: [Architecture.md](Architecture.md#section)
- **Database**: [Database.md](Database.md#section)

---

## 5. Diagnostic checklist

A reusable workflow for engineers and AI agents investigating system errors:
- [ ] **Confirm Reproduction**: Verify the error can be replicated in development or staging.
- [ ] **Verify Environment**: Walk through the Environment Verification checklist.
- [ ] **Review Recent Changes**: Scan git history, recent PRs, and infrastructure deployments.
- [ ] **Inspect Logs**: Check application, database, and load balancer logs.
- [ ] **Check APM Metrics**: Verify latency charts, memory graphs, and throughput indicators.
- [ ] **Validate Configurations**: Confirm environment variables comply with required schemas.
- [ ] **Test Dependencies**: Check third-party services status and database locks.
- [ ] **Confirm Fix**: Run integration tests to ensure the solution resolves the error.
- [ ] **Update Documentation**: Record the symptom, root cause, and resolution in this guide.

---

## 6. Incident Lessons Learned
- **Summary**: [Summary of what went wrong.]
- **Root Cause**: [What triggered the incident.]
- **Resolution**: [Actions taken to restore system operations.]
- **System Improvements**: [List follow-up actions to prevent recurrence, e.g. refactor database index logic.]

---

## 7. Monitoring & Alerting Opportunities
- **Telemetry Improvements**: [List requests for better dashboards, tracing spans, or logs format changes.]
- **New Alerts**: [e.g. Alert when cron schedules fail to execute 2 times consecutively.]

---

## 8. Known Operational Limitations
- **Accepted Limitations**: [e.g., High latency on large reporting downloads exceeding 50MB.]
- **Legacy Behavior Constraints**: [e.g. User login routes utilize older password hash formats for legacy accounts compatibility.]
- **Technical Debt Side Effects**: [Describe how existing tech debt impacts diagnosis speed.]

---

## 9. Related Documents
- **PRD**: [PRD.md](PRD.md)
- **Architecture Spec**: [Architecture.md](Architecture.md)
- **Database Schema**: [Database.md](Database.md)
- **API Contracts**: [API.md](API.md)
- **Deployment & Operations**: [Deployment.md](Deployment.md)
- **ADR Logs**: [Decisions/README.md](Decisions/README.md)

---

## AI Guidance

When assisting with troubleshooting tasks, comply with these guidelines:
- **Never Invent Root Causes**: Do not make up explanations. Distinguish between confirmed root causes, working hypotheses, and assumptions.
- **Prefer Documented Solutions**: Prioritize documented fixes over speculative suggestions.
- **Check Environment First**: Always check environment status before suggesting deep application code changes.
- **Update Post-Fix**: Suggest updating this troubleshooting guide immediately after verifying a bug fix.
- **Do Not Erase History**: Keep historical issue logs intact to preserve development institutional knowledge.

---

## Developer Guidance

- **Record Verified Fixes**: Document resolutions here immediately after resolving a production bug or Sev incident.
- **Prefer Step-by-Step Lists**: Write steps using numbered lists.
- **Audit Logging**: Confirm that manual database patches or state changes executed during incidents are logged in SRE trackers.
