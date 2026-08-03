# Discovery & Analysis Engine Specification

The Discovery & Analysis Engine forms the bedrock of the PromptEngine CLI. Rather than implementing individual codebase analysis scripts per command, all commands parse a unified `ProjectModel` compiled by this engine.

---

## 1. Pipeline Execution Stages

The discovery process runs as a sequential, multi-stage pipeline:

```text
[BaseStage: VCS/Git/Docker]
   └── [PromptEngineStage: AGENTS.md / Manifest / Configuration]
          └── [TechStage: Langs / Frameworks / DBs / CIs]
                 └── [ArchStage: MVC / Clean / DDD heuristical scans]
                        └── [DocsStage: Specifications completeness checks]
```

### Stage Boundary
Every stage must satisfy the `Stage` interface:
```go
type Stage interface {
    Name() string
    Execute(ctx context.Context, fs filesystem.FileSystem, pm *ProjectModel) error
}
```

---

## 2. Project Classification Matrix

The engine consolidates discoveries to label the target repository with multiple classifications:
- **Greenfield**: Empty workspaces without languages or frameworks.
- **Existing**: Projects containing identified codebase files.
- **PromptEngine Project**: Projects with initialized `playbook-manifest.json` files.
- **Monorepo**: Workspaces containing nested packages directories.
- **Service Typology**: Categorized into `backend_api`, `ssr_application`, `frontend_spa`, or `mobile_application` based on package inclusions.

---

## 3. Heuristical Architecture Detection

Because architecture can be subjective, the engine runs heuristic scan rules (checking directories paths mapping controllers, domains, bounded contexts) and scores them with **Confidence Levels** (0.0 to 1.0):
- **MVC** (85% Confidence): Matches nested models/controllers directories.
- **Clean Architecture** (90% Confidence): Matches segregated domain and adapter directory boundaries.
- **Domain-Driven Design (DDD)** (80% Confidence): Matches root `Domain`/`Infrastructure`/`Application` namespaces.

---

## 4. Performance & Caching Strategy

To keep execution speeds instantaneous in large monorepos, the pipeline supports:
1. **Shallow Scanning**: Stages prioritize parsing package descriptors (`package.json`, `composer.json`, `go.mod`) rather than recursively walking all directories.
2. **Path Short-circuiting**: Avoids walking ignored paths like `node_modules/`, `vendor/`, or `.git/`.
3. **Incremental Cache (Planned)**: Future builds will save computed discovery digests to `.promptengine/cache.json`, checking file hash changes before running full audits.
