# PromptEngine Implementation Notes

These items are intentionally deferred until after the initial implementation.

Their absence does not block the architecture or the first production release.

## Deferred Improvements

### 1. Public Go SDK

**Status**: Deferred

PromptEngine will initially be developed as a CLI application.

After the CLI architecture stabilizes, identify which packages should become part of the public SDK.

Do not expose internal packages prematurely.

**Goals**:
- Stable public APIs
- Minimal exported surface
- Backward compatibility
- Clear separation between internal implementation and external SDK

### 2. Scheduling Framework

**Status**: Deferred

The architecture supports future scheduling.

Implement only when real scheduling requirements exist.

**Possible future jobs**:
- Documentation synchronization
- Health recalculation
- Cache cleanup
- Update checks
- Repository indexing

The scheduler should remain optional and should not require background daemons.

### 3. Configuration Versioning

**Status**: Deferred

The initial configuration schema will remain simple.

When configuration evolves, introduce:
- configuration version numbers
- migration pipeline
- automatic upgrade path
- rollback support

### 4. Structured Application Errors

**Status**: Deferred

Quality findings already provide rich metadata.

Evaluate later whether application-level errors should include:
- severity
- category
- recommendation
- documentation reference
- automatic fix identifier

Avoid unnecessary complexity until a real need exists.

### 5. Command Aliases

**Status**: Implementation Improvement

Use Cobra's built-in aliases instead of duplicate commands where appropriate.

**Examples**:
- `detect` → `scan`
- `repair` → `doctor`
- `bootstrap` → `init`

This keeps the command tree smaller and easier to maintain.

### 6. Public Package Review

**Status**: Post-v1 Review

After implementation:
- review package boundaries
- determine which packages belong under `pkg/`
- keep implementation details under `internal/`
- expose only stable APIs

## Guiding Principle

Implementation comes before optimization.

Architecture should evolve only when real implementation experience demonstrates the need for change.

Avoid speculative features and premature abstractions.
