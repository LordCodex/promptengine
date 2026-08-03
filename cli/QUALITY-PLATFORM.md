# Quality Platform Architecture

PromptEngine's Quality Platform is the single system responsible for all evaluation, validation, auditing, and reporting of software project quality. Every PromptEngine quality-related command consumes this platform rather than implementing its own checks.

---

## Shared Foundation (`internal/domain/quality/`)

**`Finding`** is the canonical observation type emitted by every engine:

| Field | Purpose |
|---|---|
| `Engine` | Source engine (doctor, audit, compliance…) |
| `Rule` | Rule/check ID |
| `Severity` | `critical`, `error`, `warning`, `info`, `suggestion` |
| `Title` | Short human-readable description |
| `Explanation` | Why this is a problem |
| `Impact` | What breaks if not fixed |
| `Recommendation` | What to do |
| `AutoFixID` | ID of repair action (empty = no auto-fix) |

**`CategoryScore`** enables weighted per-dimension scoring. The scoring framework computes a weighted average across all categories. A `CriticalFail` flag on any category forces the overall score to 0.

**`QualityRegistry`** is the top-level plugin registration hub — plugins mount their engine registrar once.

---

## 1. Doctor Engine (`internal/domain/quality/doctor/`)

Diagnoses PromptEngine installation and project configuration health.

**Built-in checks**:
| Check ID | Category | Triggers on |
|---|---|---|
| `manifest-integrity` | integrity | Missing `playbook-manifest.json` |
| `configuration-exists` | configuration | Missing `.promptengine/` |
| `docs-directory` | documentation | Missing `docs/` |
| `plugin-integrity` | integrity | Corrupt plugin directory |
| `broken-references` | integrity | Cross-document broken links |

**Repair Framework**: `RegisterRepair(action)` maps a repair action to a finding rule ID. `Fix(id, fs, dryRun)` applies or previews the fix. Scores are computed across 5 categories with equal weights.

---

## 2. Health Engine (`internal/domain/health/`)

Produces weighted, per-category health scores for the entire project.

**Default categories and weights**:

| Category | Weight |
|---|---|
| documentation | 20% |
| architecture | 15% |
| security | 15% |
| maintainability | 10% |
| testing | 10% |
| project-knowledge | 5% |
| promptengine-adoption | 5% |
| dependency-health | 5% |
| deployment | 5% |
| observability | 5% |
| performance | 5% |

**Manifest-configurable**: `Registry.SetCategories()` accepts a custom `[]HealthCategory` with different weights — CI/org manifests override defaults without code changes.

**`CriticalFail`**: A `block`-severity issue sets the category score to 0 and forces the overall score down.

---

## 3. Review Engine (`internal/domain/review/`)

Structured code and standards reviews grouped into typed sessions.

**ReviewType taxonomy**: `architecture`, `code-organization`, `documentation`, `security`, `performance`, `testing`, `maintainability`, `scalability`, `deployment`, `observability`, `promptengine-compliance`, `org-standards`, `tech-best-practice`.

**`Reviewer`** interface groups multiple `Rule` implementations under one review type.

**`ReviewSession`** accumulates all findings and a per-type summary map. `RunSession()` runs all reviewers.

**Built-in reviewers**: documentation, security, testing, PromptEngine compliance.

---

## 4. Validation Engine (`internal/domain/quality/validation/`)

Fine-grained, independently registerable validators for every project concern.

**Built-in validators**:
| Validator ID | Checks |
|---|---|
| `project-config` | `playbook-manifest.json` exists |
| `promptengine-config` | `.promptengine/` exists |
| `manifest-schema` | Manifest is non-empty JSON |
| `documentation-completeness` | Core docs present (`Architecture.md`, `Database.md`, `API.md`) |
| `hooks-configured` | `.git/` exists for hook installation |
| `technology-compatibility` | Technology compatibility (delegated to Discovery Engine) |

Plugins call `Registry.Register(validator)` to add project-specific validators.

---

## 5. Audit Engine (`internal/domain/quality/audit/`)

Comprehensive project auditing producing exportable reports.

**Built-in rules**:
| Rule ID | Area |
|---|---|
| `project-structure` | `docs/`, `.promptengine/` exist |
| `promptengine-adoption` | Manifest present |
| `missing-documentation` | Critical docs (`Architecture.md`, `Database.md`) |
| `missing-decisions` | `Decisions.md` present |
| `architecture-drift` | Architecture doc vs Discovery Engine diff |
| `config-drift` | Manifest vs actual config diff |

**Export formats**: `json` (structured `quality.Report`) and `markdown` (summary table + findings).

---

## 6. Compliance Engine (`internal/domain/quality/compliance/`)

Runs named compliance profiles; each profile tracks pass/fail independently.

**Built-in profiles**:
| Profile ID | Type | Rules |
|---|---|---|
| `promptengine-core` | promptengine | manifest required, docs/, Decisions.md |
| `security-baseline` | security | Security.md required |

Organisations register custom profiles via `RegisterProfile(ComplianceProfile{...})`.

---

## Reporting (`internal/domain/quality/report/`)

**Supported formats**:
| Format | Use case |
|---|---|
| `text` | Terminal output |
| `json` | API / pipeline consumption |
| `yaml` | GitOps config |
| `markdown` | PR comments / docs |
| `sarif` | GitHub Advanced Security / CI integration |

`EvaluateCIThreshold(report, threshold)` returns a `CIResult` with `Passed bool` for CI integration. Plugins register additional renderers via `RendererRegistry.Register()`.

---

## Auto-Fix Framework (`internal/domain/quality/fix/`)

**Safety levels**:
| Level | Behaviour |
|---|---|
| `auto` | Applied without confirmation |
| `review` | Applied only after human review in non-interactive CI |
| `manual` | Guidance only — no automation |

**`FixEngine` modes**: `Preview()` → dry-run description; `Apply(dryRun=true)` → simulates; `Apply(dryRun=false)` → executes; `Rollback()` → reverses.

**Org disable**: `SetEnabled(false)` blocks all auto-fixes across the engine.

**Built-in actions**: `create-docs-dir`, `create-manifest`, `create-promptengine-dir`.

---

## Plugin Extension Points

Every engine exposes a registration method. Plugins call it during `OnEnable()`:

| Engine | Registration method |
|---|---|
| Doctor | `DoctorEngine.Register(DoctorCheck)` |
| Health | `Registry.Register(Checker)` |
| Review | `Registry.RegisterReviewer(Reviewer)` |
| Validation | `Registry.Register(Validator)` |
| Audit | `AuditEngine.Register(AuditRule)` |
| Compliance | `ComplianceEngine.RegisterProfile(ComplianceProfile)` |
| Reporting | `RendererRegistry.Register(ReportRenderer)` |
| Auto-Fix | `FixEngine.Register(RepairAction)` |

---

## CI/CD Integration

```yaml
# GitHub Actions example
- name: PromptEngine Quality Check
  run: promptengine doctor --format json --threshold 70
```

The CLI exits non-zero when `CIResult.Passed == false`, integrating cleanly with GitHub Actions, GitLab CI, Azure Pipelines, Jenkins, and Buildkite.
