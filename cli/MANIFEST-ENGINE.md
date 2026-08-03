# Manifest Engine Architecture Specification

The Manifest Engine acts as the central declarative intelligence of PromptEngine. It unifies core playbooks, plugin additions, organization guides, and project settings into a single versioned schema, eliminating hardcoded workflow rules.

---

## 1. Concurrency-Safe Registries & Cache

The manifest parses multiple sources asynchronously and caches the compiled result:

```text
[Core Manifest] ──┐
[Org Manifest]  ──┼─> [Registry Merging & Overrides] ─> [Merged cache mapping]
[Project Conf]  ──┘
```

Startup optimizations:
1. **Invalidation Flags**: The compiled cache remains in memory. It is marked dirty and rebuilt only when a new manifest file is loaded.
2. **File Caching**: Reads are parsed once, resolving subsequent queries in memory without filesystem calls.

---

## 2. Precedence Overrides & Conflict Resolution

When merging multiple manifests, the engine resolves key conflicts using a strict hierarchy mapping:

```text
[Highest priority: Local Project Manifest] 
   └── [Plugin Extensions Manifest] 
          └── [Organization Manifest] 
                 └── [Lowest priority: Core PromptEngine Manifest]
```
If the project manifest defines a custom `tech-stack` or overrides standard `health_rules` thresholds, the engine overrides the corresponding keys natively.

---

## 3. Query API Capabilities

Other CLI modules consult the Manifest Query API instead of parsing files:
- **`FindWorkflow(id)`**: Fetches pipeline pre/postconditions.
- **`FindStandard(id)`**: Fetches priority score mapping for the Context Engine.
- **`FindTech(id)`**: Resolves tech-specific playbooks lists.
- **`VerifyCompatibility(cliVersion)`**: Runs version checks.
