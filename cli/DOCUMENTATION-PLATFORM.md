# Documentation Platform Architecture

PromptEngine's Documentation Platform is the single system responsible for creating, maintaining, validating, synchronising and evolving every document in a project. It integrates with the Discovery Engine, Context Engine, Workflow Engine and Manifest Engine.

---

## 1. Documentation Engine (`internal/domain/docs/`)

The Documentation Engine is the orchestration layer for all document operations.

**`DocumentSpec`** describes a registered document type:

| Field | Purpose |
|---|---|
| `ID` | Unique identifier (e.g., `"database"`) |
| `RequiredSections` | Sections every file of this type must contain |
| `DependsOn` | IDs of other documents this doc depends on |
| `UpdateTriggers` | Change signals that should trigger an update |
| `Priority` | Health score weight |

**`DocRegistry`** is the central catalogue:
- `Register(spec)` — plugins and manifest can add types without modifying core code.
- `AffectedBy(trigger)` — returns all doc IDs that declared a given trigger.

**`Engine`** wires all sub-systems:
- `Status()` — scans the filesystem and classifies every registered spec as `synced`, `missing`, or `stale`.
- `Validate(docID)` — checks required sections exist in the file.
- `RecordGeneration(docID, version, content)` — writes a `VersionedDoc` entry for audit and rollback.

---

## 2. Template Engine (`internal/domain/docs/template/`)

**`Template`** supports:
- **Inheritance** via `ParentID` — child templates extend parent sections.
- **Conditional Sections** via `Condition` — a section is omitted unless the named variable is non-empty.
- **Partials** via `{{partial:name}}` — inline another registered template.
- **Variables** resolved with `{variable_name}` syntax.

**`TemplateRegistry`** accepts templates from any source: `core`, `org`, `plugin`, `technology`, `marketplace`.

**`Renderer`** resolves the full inheritance chain, evaluates conditionals, inlines partials, and injects all variables before returning the final string.

---

## 3. Generator Engine (`internal/domain/docs/generator/`)

**`Generator`** interface — any type implementing `DocType()` and `Generate(input)` can be registered.

**`GeneratorRegistry`** is plugin-extensible: no switch statements, no hardcoded list.

**Standard generators** registered by `RegisterDefaults()`:

| DocType | Output file |
|---|---|
| `architecture` | `docs/Architecture.md` |
| `business-rules` | `docs/BusinessRules.md` |
| `database` | `docs/Database.md` |
| `api` | `docs/API.md` |
| `prd` | `docs/PRD.md` |
| `progress` | `docs/Progress.md` |
| `roadmap` | `docs/Roadmap.md` |
| `deployment` | `docs/Deployment.md` |
| `troubleshooting` | `docs/Troubleshooting.md` |
| `decisions` | `docs/Decisions.md` |
| `security` | `docs/Security.md` |
| `testing` | `docs/Testing.md` |

Plugins contribute additional generators by calling `registry.Register(myGenerator)`.

---

## 4. Prompt Engine (`internal/domain/prompts/`)

**`PromptDef`** is the full definition of a reusable prompt:

| Field | Purpose |
|---|---|
| `Workflow` | Engineering activity (e.g., `new-project`, `bug-fix`) |
| `Source` | Origin (`core`, `org`, `plugin`, `technology`, `provider`) |
| `RequiredContext` | Context keys that must be populated before building |
| `Template` | Prompt body with `{variable}` placeholders |
| `ProviderHints` | Per-provider prefix/suffix formatting |

**`PromptRegistry`** supports queries by workflow and by source, with no hardcoded prompt selection.

**`PromptBuilder`** pipeline:
1. Merges `ContextPackage` (from Context Engine) into template variables.
2. Validates required context keys are populated.
3. Injects variables using `Inject()`.
4. Applies provider-specific hint formatting.
5. Returns a `Prompt` ready for copy-paste or AI Provider Engine submission.

**Supported workflows**: `new-project`, `existing-project`, `migration`, `feature-development`, `bug-fix`, `refactoring`, `architecture-review`, `security-review`, `performance-review`, `deployment`, `documentation-review`, `project-audit`, `ai-session-bootstrap`, `context-refresh`, `prompt-improvement`.

---

## 5. Synchronisation Engine (`internal/domain/docs/sync/`)

**`ChangeSignal`** is a detected project change event.

**`SyncRule`** maps each signal to the document IDs that must be reviewed:

| Signal | Affected Documents |
|---|---|
| `new-migration` | database, architecture, api, deployment, progress, roadmap |
| `architecture-change` | architecture, decisions, business-rules, deployment |
| `new-api` | api, architecture, deployment, progress |
| `new-service` | architecture, deployment, database, api |
| `new-technology` | architecture, deployment, decisions |

**`SyncEngine`** modes:
- **Dry-run**: Recommends updates without applying any changes.
- **Auto-apply**: Applies changes where `AutoApply = true` (safe signals only).
- **Pending approval**: Queues changes requiring human review.

Plugins and workflows register additional `SyncRule` declarations.

---

## 6. Validation Engine (`internal/domain/docs/validation/`)

**`ValidationRule`** interface — plugins register additional rules without modifying core code.

**Built-in rules**:

| Rule | Severity | Description |
|---|---|---|
| `missing-document` | error | File does not exist on disk |
| `missing-sections` | error | Document has no headings |
| `stale-document` | warning | Contains placeholder text (`_Define`, `TBD`, `TODO`) |
| `broken-references` | warning | Markdown links pointing to non-existent files |
| `duplicate-content` | warning | Repeated top-level headings |
| `orphaned-document` | info | Very short document likely not connected to the manifest |

Every finding includes a `Suggestion` field with a concrete, actionable recommendation.

---

## Registries

All four registries follow the same pattern:
1. **Core defaults** registered at startup.
2. **Plugin contributions** registered during plugin `OnEnable()`.
3. **Org overrides** registered by the Manifest Engine on startup.

No hardcoded selection — all lookups go through the registry.

---

## Versioning

- `VersionedDoc` records every generation event (doc ID, version, timestamp, content).
- The Update Engine manages migration between schema versions.
- Templates carry a `Version` field for compatibility tracking.

---

## AI Integration

The Documentation Platform maintains a strict separation of responsibility:

| Layer | Responsibility |
|---|---|
| **Context Engine** | Determines _what_ information is required |
| **Prompt Engine** | Determines _how_ the AI is instructed |
| **AI Provider Engine** | Determines _where_ the prompt is executed |
| **Documentation Engine** | Validates and persists the resulting documentation |

The Documentation Platform never communicates directly with an AI provider.

---

## Plugin Extension Points

Plugins may contribute:
- **Templates** → `TemplateRegistry.Register()`
- **Generators** → `GeneratorRegistry.Register()`
- **Prompts** → `PromptRegistry.Register()`
- **Validation Rules** → `Validator.Register()`
- **Sync Rules** → `ChangeDetector.RegisterRule()`
- **Document Types** → `DocRegistry.Register()`
