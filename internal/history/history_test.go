package history

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

func TestRecorderRecordsRedactedJSONL(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	recorder := NewRecorder(fs)

	err := recorder.Record(Entry{
		Command:   "promptengine prompt",
		Status:    "failed",
		StartedAt: time.Now().UTC(),
		Error:     "provider failed with api_key=sk_live_secret",
	})
	if err != nil {
		t.Fatalf("record failed: %v", err)
	}
	data, err := fs.ReadFile(DefaultPath)
	if err != nil {
		t.Fatalf("history file missing: %v", err)
	}
	if strings.Contains(string(data), "sk_live_secret") {
		t.Fatalf("history leaked secret: %s", data)
	}
	var entry Entry
	if err := json.Unmarshal([]byte(strings.TrimSpace(string(data))), &entry); err != nil {
		t.Fatalf("invalid history json: %v", err)
	}
	if entry.Command != "promptengine prompt" || entry.Status != "failed" || entry.Duration == "" {
		t.Fatalf("unexpected history entry: %#v", entry)
	}
}
