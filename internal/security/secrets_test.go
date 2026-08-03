package security

import (
	"strings"
	"testing"
)

func TestRedactSecrets(t *testing.T) {
	input := "API_KEY=sk_test_1234567890\npassword: hunter2\nAuthorization: Bearer abcdefghijklmnopqrstuvwxyz"
	out, changed := RedactSecrets(input)
	if !changed {
		t.Fatal("expected secrets to be redacted")
	}
	for _, leak := range []string{"sk_test_1234567890", "hunter2", "abcdefghijklmnopqrstuvwxyz"} {
		if strings.Contains(out, leak) {
			t.Fatalf("redacted output leaked %q: %s", leak, out)
		}
	}
	if !strings.Contains(out, "[REDACTED]") {
		t.Fatalf("expected redaction marker, got %s", out)
	}
}

func TestSensitivePathDetection(t *testing.T) {
	for _, path := range []string{".env", "config/secrets/app.yaml", ".ssh/id_ed25519", "project/.npmrc"} {
		if !IsSensitivePath(path) {
			t.Fatalf("expected %s to be sensitive", path)
		}
	}
	if IsSensitivePath("internal/domain/context/engine.go") {
		t.Fatal("ordinary source file should not be sensitive")
	}
}
