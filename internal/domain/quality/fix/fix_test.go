package fix

import (
	"testing"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

func TestFixEngine_Preview(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	engine := NewFixEngine()

	preview, err := engine.Preview("create-docs-dir", fs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if preview == "" {
		t.Error("expected non-empty preview description")
	}
}

func TestFixEngine_Apply_DocsDir(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	engine := NewFixEngine()

	result := engine.Apply("create-docs-dir", fs, false)
	if result.Error != nil {
		t.Fatalf("expected apply to succeed, got: %v", result.Error)
	}
	if !result.Applied {
		t.Error("expected Applied=true after fix")
	}
	if !fs.Exists("docs/.keep") {
		t.Error("expected docs/.keep to be created")
	}
}

func TestFixEngine_DryRun_NoChanges(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	engine := NewFixEngine()

	result := engine.Apply("create-docs-dir", fs, true)
	if result.Error != nil {
		t.Fatalf("unexpected error in dry-run: %v", result.Error)
	}
	if result.Applied {
		t.Error("expected Applied=false in dry-run mode")
	}
	if !result.DryRun {
		t.Error("expected DryRun=true")
	}
	if fs.Exists("docs/.keep") {
		t.Error("expected dry-run to NOT create files")
	}
}

func TestFixEngine_Rollback(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	engine := NewFixEngine()

	// Apply then rollback
	engine.Apply("create-docs-dir", fs, false)
	result := engine.Rollback("create-docs-dir", fs)
	if result.Error != nil {
		t.Fatalf("expected rollback to succeed: %v", result.Error)
	}
	if !result.Rolled {
		t.Error("expected Rolled=true after rollback")
	}
}

func TestFixEngine_Disabled_BlocksAutoFix(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	engine := NewFixEngine()
	engine.SetEnabled(false)

	result := engine.Apply("create-docs-dir", fs, false)
	if result.Applied {
		t.Error("expected fix to be blocked when auto-fixes are disabled")
	}
}

func TestFixEngine_UnknownAction(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	engine := NewFixEngine()

	result := engine.Apply("nonexistent-fix", fs, false)
	if result.Error == nil {
		t.Error("expected error for unknown repair action")
	}
}

// customRepairAction simulates a plugin-contributed repair action
type customRepairAction struct{ applied bool }

func (a *customRepairAction) ID() string          { return "custom-fix" }
func (a *customRepairAction) Description() string { return "Custom fix" }
func (a *customRepairAction) Safety() FixSafety   { return SafeAuto }
func (a *customRepairAction) Preview(_ filesystem.FileSystem) (string, error) {
	return "Will apply custom fix.", nil
}
func (a *customRepairAction) Apply(_ filesystem.FileSystem) error {
	a.applied = true
	return nil
}
func (a *customRepairAction) Rollback(_ filesystem.FileSystem) error { return nil }

func TestFixEngine_PluginAction(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	engine := NewFixEngine()
	action := &customRepairAction{}
	engine.Register(action)

	result := engine.Apply("custom-fix", fs, false)
	if result.Error != nil {
		t.Fatalf("unexpected error: %v", result.Error)
	}
	if !action.applied {
		t.Error("expected plugin repair action to be applied")
	}
}

func TestFixEngine_ManifestAction_RequiresReview(t *testing.T) {
	// create-manifest has SafeReview — should still apply when enabled
	fs := filesystem.NewMockFileSystem()
	engine := NewFixEngine()

	result := engine.Apply("create-manifest", fs, false)
	if result.Error != nil {
		t.Fatalf("unexpected error: %v", result.Error)
	}
	if !result.Applied {
		t.Error("expected SafeReview action to apply when engine is enabled")
	}
	if !fs.Exists("playbook-manifest.json") {
		t.Error("expected manifest file to be created")
	}
}
