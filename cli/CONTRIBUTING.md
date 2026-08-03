# PromptEngine CLI Developer Contribution Guidelines

Welcome to the PromptEngine CLI developer guide! This document explains how the CLI project is architected and how to contribute new features, command specs, and stack detectors.

---

## 1. Clean Architecture Design

The CLI uses Go clean architecture conventions separating presentation, orchestration logic, and domain engines:

- **Command Router Layer (`cmd/`, `internal/app/`)**: Built on Cobra command trees. Subcommands delegate inputs to coordinators.
- **TUI UI Layer (`internal/ui/`)**: Built on Charm bubbletea. Contains reusable questions models (dropdowns, multiselect lists, spinner loaders).
- **Core Domain Layer (`internal/domain/`)**: Pure business logic (Context Builder token calculations, Health scoring, project detectors interfaces).

### Key Rules
1. **Zero Global States**: Packages must avoid using package-level mutable variables. Always instantiate structures using factory functions.
2. **Interface Segregation**: All filesystem and provider calls must use defined boundary interfaces to allow unit mock testing.
3. **No Circular Imports**: Dependencies must flow inwards. Pure domain structures cannot import UI or application coordinator packages.

---

## 2. Contributing a New Command

To add a command (e.g. `generate`):
1. **Define command specifications**: Create `cli/specs/generate.md` defining inputs, outputs, TTY flows, and error maps.
2. **Register Command**: Add the builder method in `internal/app/commands.go` returning `*cobra.Command`.
3. **Wire command**: Register the builder inside the bootstrap mapping sequence in `internal/app/bootstrap.go`.
4. **Implement core use-case logic**: Put the logic inside `internal/domain/docs/` or separate engine modules, never in the cobra command `Run` function itself.

---

## 3. Contributing a Stack Detector

To write a detector (e.g. detecting Python Flask stack):
1. **Create the file**: Put it inside `internal/domain/detection/`.
2. **Implement `Detector` interface**:
   ```go
   type Detector interface {
       Name() string
       Detect(fs filesystem.FileSystem) (bool, error)
       Apply(fs filesystem.FileSystem, meta *ProjectMetadata) error
   }
   ```
3. **Register Detector**: Call `registry.Register()` during discovery engine bootstrapping.

---

## 4. Running Tests

- Run all unit tests:
  ```bash
  go test ./...
  ```
- Run tests with coverage output:
  ```bash
  go test -cover ./...
  ```
