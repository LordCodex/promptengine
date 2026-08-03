# Extensibility Platform Architecture

PromptEngine is designed so that all new capabilities are added through the Extensibility Platform rather than by modifying core code. The platform consists of five cooperating subsystems.

---

## 1. Plugin Engine (`internal/domain/plugins/`)

Plugins contribute to PromptEngine through declared **ContributionPoints**:

| Point | Examples |
|---|---|
| `Commands` | `pe-laravel`, `pe-vue-analyze` |
| `Workflows` | `laravel-feature`, `vue-component` |
| `Standards` | `laravel-conventions.md` |
| `Detectors` | Laravel framework detector |
| `Validators` | Blade template validator |
| `AIProviders` | Custom LLM adapters |
| `HookTypes` | Custom automation hooks |

**Lifecycle**: `OnInstall → OnEnable → [running] → OnDisable → OnUninstall → OnUpdate`

**Dependency Resolution**: The `Registry.ResolveLoadOrder()` method performs a topological sort, detecting cyclic dependencies before any plugin activates.

**Version Compatibility**: `CheckCompatibility(meta, coreVer)` guards against installing plugins that require a newer core.

---

## 2. Hook Engine (`internal/domain/hooks/`)

Hooks run automated checks at Git or CI trigger points.

**Policy Modes**:
- `PolicyWarn` (default) — emits a warning but allows the operation to continue.
- `PolicyEnforce` — blocks the Git operation or CI step on failure.

**Git Hook Script**: `GitHook.Install()` writes a managed shell script to `.git/hooks/`. Scripts include PromptEngine attribution comments to avoid manual overwrites.

**CI Template Generator**: `CITemplate.Generate()` produces ready-to-use YAML for GitHub Actions, GitLab CI, and Azure Pipelines.

**Plugin-Driven**: Plugins register their own hook types by contributing `HookTypes` IDs to the Plugin Registry.

---

## 3. Installer Engine (`internal/domain/installer/`)

The `Installer` interface supports multiple installation backends:

| Backend | Status |
|---|---|
| `LocalInstaller` | Implemented — installs from local filesystem paths |
| `RemoteInstaller` | Future — downloads from marketplace or remote URLs |

**Transaction Safety**: Every installation writes a `.installed` marker file. `IsInstalled()` checks this marker without needing a database. `Uninstall()` removes the marker for clean state.

**Supported Kinds**: `plugin`, `technology-pack`, `org-pack`, `template`, `workflow-pack`, `prompt-library`, `ai-provider`, `standard`.

---

## 4. Update Engine (`internal/domain/updater/`)

The Update Engine uses a **plan → apply → rollback** model:

1. **Plan (dry-run)**: `Plan(req)` returns a `UpdateReport` with compatibility issues and migration strategy warnings. No changes are made.
2. **Apply**: `Apply(req, currentData)` snapshots the current state then performs the update.
3. **Rollback**: `Rollback(id)` restores the pre-update snapshot.

**Migration Strategies**: Components register `MigrationStrategy` declarations describing the steps required for a specific version transition. If no strategy exists for a version pair, the engine emits a warning advisory.

---

## 5. Marketplace Foundation (`internal/domain/marketplace/`)

The `MarketplaceClient` interface defines the contract for any future marketplace backend:

```go
type MarketplaceClient interface {
    Search(ctx context.Context, filter SearchFilter) ([]PackageListing, error)
    GetPackage(ctx context.Context, id string) (*PackageListing, error)
    Download(ctx context.Context, req InstallRequest) ([]byte, error)
    Publish(ctx context.Context, manifest []byte) error
}
```

No concrete implementation is provided yet. The interface is defined now so that callers in the CLI can depend on it without breaking changes when a real marketplace is connected.
