# Workflow Engine Architecture Specification

The Workflow Engine orchestrates the execution of multi-step PromptEngine tasks (such as project bootstrapping, code changes analysis, and documentation synchronization). It isolates step boundaries, tracks pipeline states, and dispatches events.

---

## 1. Execution Lifecycle stages

Every workflow runs as a sequence of stages that fulfill the `PipelineStage` interface:

```text
[PreconditionStage: Files & Configuration checks]
   └── [DiscoveryStage: Telemetry Project scan]
          └── [ContextStage: Playbooks load & Token minimize]
                 └── [PostconditionStage: Compliance & Heritage verification]
```

### PipelineStage Contract
```go
type PipelineStage interface {
    Name() string
    Run(ctx context.Context, fs filesystem.FileSystem, flow *FlowContext) error
}
```

---

## 2. Preconditions & Postconditions

- **Preconditions**: Asserts that baseline workspace parameters exist before running the execution. For example, verifying a stack composer manifest exists before migrating Laravel settings.
- **Postconditions**: Asserts completion compliance criteria after task completion:
  1. Checks that required target documentation files were successfully updated.
  2. Asserts that generated constitutional `AGENTS.md` rules inherit relative links to the universal coding standard playbooks:
     - `core/05-universal-coding-standards.md`
     Inheriting this reference is a mandatory rule that prevents duplicates and maintains repository guidelines structure.

---

## 3. Event Publisher System

Workflows publish telemetry updates dynamically to registered subscribers:
- **`WorkflowStarted`**: Fires when starting a pipeline run.
- **`ContextGenerated`**: Fires after the Context Engine calculates token payloads.
- **`ValidationPassed`**: Fires when precondition audits validate correctly.
- **`WorkflowCompleted`**: Fires after postcondition verification success.
- **`WorkflowFailed`**: Fires upon runtime execution exceptions.

Plugins, loggers, or progress UI adapters subscribe to these event channels to render status updates asynchronously without coupling.
