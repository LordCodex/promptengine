package app

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestContextLimitAliasFlags(t *testing.T) {
	tests := map[string]*cobra.Command{
		"context":        newContextCommand(&App{}),
		"context export": newContextExportCommand(&App{}),
		"prompt":         newPromptCommand(&App{}),
		"task":           newTaskCommand(&App{}),
	}
	for name, cmd := range tests {
		if cmd.Flags().Lookup("max-bytes") == nil {
			t.Fatalf("%s should expose --max-bytes", name)
		}
		if cmd.Flags().Lookup("limit") == nil {
			t.Fatalf("%s should expose --limit alias", name)
		}
	}
}
