package updater

import "fmt"

// UpdateTarget identifies what is being updated
type UpdateTarget string

const (
	TargetCLI      UpdateTarget = "cli"
	TargetPlugin   UpdateTarget = "plugin"
	TargetManifest UpdateTarget = "manifest"
	TargetTemplate UpdateTarget = "template"
	TargetSchema   UpdateTarget = "schema"
)

// UpdateRequest declares what to update and to which version
type UpdateRequest struct {
	Target  UpdateTarget
	ID      string // component ID (plugin ID, etc.)
	FromVer string
	ToVer   string
	DryRun  bool
}

// CompatibilityIssue describes a blocking or advisory conflict found before updating
type CompatibilityIssue struct {
	Severity string // "error" or "warning"
	Message  string
}

// UpdateReport is the result of a dry-run or actual update
type UpdateReport struct {
	Request           UpdateRequest
	Issues            []CompatibilityIssue
	CanUpdate         bool
	Applied           bool
	RollbackAvailable bool
}

// MigrationStrategy declares how to handle a specific version transition
type MigrationStrategy struct {
	FromVer string
	ToVer   string
	Steps   []string // human-readable migration steps
}

// Snapshot captures pre-update state for rollback
type Snapshot struct {
	ID      string
	Target  UpdateTarget
	Version string
	Data    []byte // serialised state
}

// UpdateEngine orchestrates safe component updates
type UpdateEngine struct {
	snapshots  map[string]*Snapshot
	strategies []MigrationStrategy
}

func NewUpdateEngine() *UpdateEngine {
	return &UpdateEngine{
		snapshots:  make(map[string]*Snapshot),
		strategies: make([]MigrationStrategy, 0),
	}
}

// RegisterMigration allows components to declare how a version transition works
func (e *UpdateEngine) RegisterMigration(m MigrationStrategy) {
	e.strategies = append(e.strategies, m)
}

// Plan runs a dry-run compatibility check without applying changes
func (e *UpdateEngine) Plan(req UpdateRequest) UpdateReport {
	report := UpdateReport{Request: req, CanUpdate: true}

	// Check if downgrade is attempted (simple semver lexicographic guard)
	if req.FromVer > req.ToVer {
		report.Issues = append(report.Issues, CompatibilityIssue{
			Severity: "error",
			Message:  fmt.Sprintf("downgrade from %s to %s is not supported", req.FromVer, req.ToVer),
		})
		report.CanUpdate = false
	}

	// Check for available migration strategies
	found := false
	for _, s := range e.strategies {
		if s.FromVer == req.FromVer && s.ToVer == req.ToVer {
			found = true
			break
		}
	}
	if !found && req.FromVer != "" && req.ToVer != "" {
		report.Issues = append(report.Issues, CompatibilityIssue{
			Severity: "warning",
			Message:  fmt.Sprintf("no migration strategy declared for %s → %s; manual review recommended", req.FromVer, req.ToVer),
		})
	}

	report.RollbackAvailable = e.snapshots[req.ID] != nil
	return report
}

// Apply executes the update after snapshotting current state
func (e *UpdateEngine) Apply(req UpdateRequest, currentData []byte) (UpdateReport, error) {
	report := e.Plan(req)
	if !report.CanUpdate {
		return report, fmt.Errorf("update blocked: %d compatibility issue(s) found", len(report.Issues))
	}
	if req.DryRun {
		return report, nil
	}

	// Save rollback snapshot
	e.snapshots[req.ID] = &Snapshot{
		ID:      req.ID,
		Target:  req.Target,
		Version: req.FromVer,
		Data:    currentData,
	}

	report.Applied = true
	report.RollbackAvailable = true
	return report, nil
}

// Rollback restores pre-update snapshot
func (e *UpdateEngine) Rollback(id string) (*Snapshot, error) {
	snap, ok := e.snapshots[id]
	if !ok {
		return nil, fmt.Errorf("no rollback snapshot found for '%s'", id)
	}
	delete(e.snapshots, id)
	return snap, nil
}
