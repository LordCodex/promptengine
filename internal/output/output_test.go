package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestConfiguredRenderer_JSON(t *testing.T) {
	var buf bytes.Buffer
	renderer := NewConfiguredRenderer(FormatJSON, false, false)

	if err := renderer.Render(&buf, map[string]string{"status": "ready"}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(buf.String(), `"status": "ready"`) {
		t.Fatalf("expected JSON output, got %q", buf.String())
	}
}

func TestConfiguredRenderer_YAML(t *testing.T) {
	var buf bytes.Buffer
	renderer := NewConfiguredRenderer(FormatYAML, false, false)

	if err := renderer.Render(&buf, map[string]string{"status": "ready"}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(buf.String(), "status: ready") {
		t.Fatalf("expected YAML output, got %q", buf.String())
	}
}

func TestConfiguredRenderer_Text(t *testing.T) {
	var buf bytes.Buffer
	renderer := NewConfiguredRenderer(FormatText, false, false)

	if err := renderer.Render(&buf, "ready"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if buf.String() != "ready\n" {
		t.Fatalf("expected text output, got %q", buf.String())
	}
}
