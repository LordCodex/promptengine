package sync

import (
	"testing"
)

func TestChangeDetector_Detect_SingleSignal(t *testing.T) {
	detector := NewChangeDetector()
	recs := detector.Detect([]ChangeSignal{SignalNewMigration})
	if len(recs) == 0 {
		t.Error("expected recommendations for new-migration signal")
	}
	// Database should always be in the affected docs for a migration
	found := false
	for _, r := range recs {
		if r.DocID == "database" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'database' to be in recommendations for new-migration")
	}
}

func TestChangeDetector_Detect_NoDuplicates(t *testing.T) {
	detector := NewChangeDetector()
	// Two identical signals should not double the recommendations
	recs := detector.Detect([]ChangeSignal{SignalNewMigration, SignalNewMigration})
	seen := make(map[string]int)
	for _, r := range recs {
		seen[r.DocID]++
	}
	for docID, count := range seen {
		if count > 1 {
			t.Errorf("expected no duplicate recommendation for '%s', got %d", docID, count)
		}
	}
}

func TestChangeDetector_RegisterCustomRule(t *testing.T) {
	detector := NewChangeDetector()
	detector.RegisterRule(SyncRule{
		Signal:       ChangeSignal("new-feature-flag"),
		AffectedDocs: []string{"decisions", "progress"},
		Description:  "New feature flag requires Decisions and Progress updates.",
	})
	recs := detector.Detect([]ChangeSignal{ChangeSignal("new-feature-flag")})
	if len(recs) != 2 {
		t.Errorf("expected 2 recommendations from custom rule, got %d", len(recs))
	}
}

func TestSyncEngine_DryRun(t *testing.T) {
	engine := NewSyncEngine(NewChangeDetector())
	result := engine.Run([]ChangeSignal{SignalArchitectureChange}, true)
	if !result.DryRun {
		t.Error("expected DryRun to be true in dry-run mode")
	}
	if len(result.Applied) != 0 {
		t.Error("expected no applied documents in dry-run mode")
	}
	if len(result.Pending) == 0 {
		t.Error("expected pending recommendations in dry-run mode")
	}
}

func TestSyncEngine_AutoApply_SafeSignals(t *testing.T) {
	engine := NewSyncEngine(NewChangeDetector())
	result := engine.Run([]ChangeSignal{SignalNewDocumentation}, false)
	if len(result.Applied) == 0 {
		t.Error("expected auto-apply for documentation signal")
	}
}
