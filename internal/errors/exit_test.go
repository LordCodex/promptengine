package errors

import (
	stderrors "errors"
	"testing"
)

func TestAppError_WrapsCause(t *testing.T) {
	cause := stderrors.New("missing config")
	err := New(CategoryConfiguration, ExitConfiguration, "configuration failed", cause)

	if !stderrors.Is(err, cause) {
		t.Fatal("expected AppError to unwrap cause")
	}
	if ExitCode(err) != ExitConfiguration {
		t.Fatalf("expected configuration exit code, got %d", ExitCode(err))
	}
	if err.Category != CategoryConfiguration {
		t.Fatalf("expected configuration category, got %s", err.Category)
	}
}
