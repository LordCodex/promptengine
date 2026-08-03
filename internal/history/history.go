package history

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/LordCodex/promptengine/internal/filesystem"
	"github.com/LordCodex/promptengine/internal/security"
)

const DefaultPath = ".promptengine/history.jsonl"

type Entry struct {
	Command   string            `json:"command"`
	Status    string            `json:"status"`
	StartedAt time.Time         `json:"started_at"`
	EndedAt   time.Time         `json:"ended_at"`
	Duration  string            `json:"duration"`
	Error     string            `json:"error,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

type Recorder struct {
	fs   filesystem.FileSystem
	path string
}

func NewRecorder(fs filesystem.FileSystem) *Recorder {
	return &Recorder{fs: fs, path: DefaultPath}
}

func NewRecorderAt(fs filesystem.FileSystem, path string) *Recorder {
	if path == "" {
		path = DefaultPath
	}
	return &Recorder{fs: fs, path: path}
}

func (r *Recorder) Record(entry Entry) error {
	if r == nil || r.fs == nil {
		return nil
	}
	if entry.EndedAt.IsZero() {
		entry.EndedAt = time.Now().UTC()
	}
	if entry.Duration == "" && !entry.StartedAt.IsZero() {
		entry.Duration = entry.EndedAt.Sub(entry.StartedAt).String()
	}
	entry.Error, _ = security.RedactSecrets(entry.Error)
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	line := append(data, '\n')
	if appender, ok := r.fs.(interface {
		AppendFile(string, []byte, os.FileMode) error
	}); ok {
		if err := appender.AppendFile(r.path, line, 0644); err != nil {
			return fmt.Errorf("write audit history: %w", err)
		}
		return nil
	}

	var existing []byte
	if r.fs.Exists(r.path) {
		existing, err = r.fs.ReadFile(r.path)
		if err != nil {
			return err
		}
	}
	content := strings.TrimRight(string(existing), "\n")
	if content != "" {
		content += "\n"
	}
	content += string(line)
	if err := r.fs.WriteFile(r.path, []byte(content), 0644); err != nil {
		return fmt.Errorf("write audit history: %w", err)
	}
	return nil
}
