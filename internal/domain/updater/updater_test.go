package updater

import "testing"

func TestUpdateEngine_DryRun(t *testing.T) {
	engine := NewUpdateEngine()
	req := UpdateRequest{
		Target:  TargetPlugin,
		ID:      "laravel-pack",
		FromVer: "1.0.0",
		ToVer:   "2.0.0",
		DryRun:  true,
	}
	report, err := engine.Apply(req, nil)
	if err != nil {
		t.Fatalf("expected dry-run to succeed, got: %v", err)
	}
	if report.Applied {
		t.Error("expected Applied=false on dry-run")
	}
	if report.CanUpdate != true {
		t.Error("expected CanUpdate=true for valid version bump")
	}
}

func TestUpdateEngine_Downgrade_Blocked(t *testing.T) {
	engine := NewUpdateEngine()
	req := UpdateRequest{
		Target:  TargetPlugin,
		ID:      "laravel-pack",
		FromVer: "2.0.0",
		ToVer:   "1.0.0",
	}
	report := engine.Plan(req)
	if report.CanUpdate {
		t.Error("expected downgrade to be blocked")
	}
	if len(report.Issues) == 0 {
		t.Error("expected at least one compatibility issue for downgrade")
	}
}

func TestUpdateEngine_Apply_And_Rollback(t *testing.T) {
	engine := NewUpdateEngine()
	req := UpdateRequest{
		Target:  TargetPlugin,
		ID:      "my-plugin",
		FromVer: "1.0.0",
		ToVer:   "1.1.0",
	}
	originalData := []byte(`{"version":"1.0.0"}`)

	report, err := engine.Apply(req, originalData)
	if err != nil {
		t.Fatalf("expected apply to succeed, got: %v", err)
	}
	if !report.Applied {
		t.Error("expected Applied=true after successful update")
	}
	if !report.RollbackAvailable {
		t.Error("expected rollback to be available after apply")
	}

	snap, err := engine.Rollback("my-plugin")
	if err != nil {
		t.Fatalf("expected rollback to succeed, got: %v", err)
	}
	if string(snap.Data) != string(originalData) {
		t.Errorf("expected rollback data to match original, got: %s", snap.Data)
	}
}

func TestUpdateEngine_MigrationWarning(t *testing.T) {
	engine := NewUpdateEngine()
	// No migration strategy registered — should warn
	req := UpdateRequest{
		Target:  TargetSchema,
		ID:      "core-schema",
		FromVer: "1.0.0",
		ToVer:   "2.0.0",
		DryRun:  true,
	}
	report := engine.Plan(req)
	found := false
	for _, issue := range report.Issues {
		if issue.Severity == "warning" {
			found = true
		}
	}
	if !found {
		t.Error("expected warning when no migration strategy is registered")
	}
}
