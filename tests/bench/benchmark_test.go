package bench

import (
	"context"
	"io"
	"testing"

	"github.com/LordCodex/promptengine/internal/app"
	"github.com/LordCodex/promptengine/internal/filesystem"
)

func BenchmarkStartup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := app.Bootstrap(io.Discard, false)
		if err != nil {
			b.Fatalf("failed to bootstrap: %v", err)
		}
	}
}

func BenchmarkDiscoveryExecute(b *testing.B) {
	a, err := app.Bootstrap(io.Discard, false)
	if err != nil {
		b.Fatalf("failed to bootstrap: %v", err)
	}

	fs := filesystem.NewMockFileSystem()
	_ = fs.WriteFile("go.mod", []byte("module myapp\ngo 1.21\n"), 0644)
	_ = fs.WriteFile("playbook-manifest.json", []byte("{}"), 0644)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := a.Discovery.Execute(context.Background(), fs, ".")
		if err != nil {
			b.Fatalf("discovery execution failed: %v", err)
		}
	}
}
