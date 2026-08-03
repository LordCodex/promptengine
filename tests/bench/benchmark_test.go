package bench

import (
	"context"
	"io"
	"testing"

	"github.com/LordCodex/promptengine/internal/app"
	contextengine "github.com/LordCodex/promptengine/internal/domain/context"
	"github.com/LordCodex/promptengine/internal/domain/discovery"
	"github.com/LordCodex/promptengine/internal/domain/docs"
	"github.com/LordCodex/promptengine/internal/domain/quality"
	"github.com/LordCodex/promptengine/internal/filesystem"
)

func BenchmarkStartup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := app.Bootstrap(app.BootstrapOptions{Out: io.Discard, Err: io.Discard})
		if err != nil {
			b.Fatalf("failed to bootstrap: %v", err)
		}
	}
}

func BenchmarkRootExecute(b *testing.B) {
	a, err := app.Bootstrap(app.BootstrapOptions{Out: io.Discard, Err: io.Discard})
	if err != nil {
		b.Fatalf("failed to bootstrap: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := a.Execute(context.Background(), []string{}); err != nil {
			b.Fatalf("root execution failed: %v", err)
		}
	}
}

func BenchmarkDiscoveryPipeline(b *testing.B) {
	fs := benchmarkFS()
	pipeline := discovery.NewDefaultPipeline(nil, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := pipeline.Execute(context.Background(), fs, "."); err != nil {
			b.Fatalf("discovery failed: %v", err)
		}
	}
}

func BenchmarkContextBuild(b *testing.B) {
	fs := benchmarkFS()
	pm, err := discovery.NewDefaultPipeline(nil, nil).Execute(context.Background(), fs, ".")
	if err != nil {
		b.Fatalf("discovery failed: %v", err)
	}
	engine := contextengine.NewEngine(fs)
	req := contextengine.ContextRequest{
		TaskType:   contextengine.TaskAddFeature,
		Project:    pm,
		UserIntent: "Add payment support",
		Budget:     contextengine.BudgetSmall,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := engine.Build(context.Background(), req); err != nil {
			b.Fatalf("context build failed: %v", err)
		}
	}
}

func BenchmarkDocsValidate(b *testing.B) {
	fs := benchmarkFS()
	platform := docs.NewPlatform(fs, nil, nil)
	pm, err := discovery.NewDefaultPipeline(nil, nil).Execute(context.Background(), fs, ".")
	if err != nil {
		b.Fatalf("discovery failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := platform.Validate(pm); err != nil {
			b.Fatalf("docs validate failed: %v", err)
		}
	}
}

func BenchmarkQualityAudit(b *testing.B) {
	fs := benchmarkFS()
	platform := quality.NewPlatform(fs, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := platform.Audit(context.Background()); err != nil {
			b.Fatalf("quality audit failed: %v", err)
		}
	}
}

func benchmarkFS() filesystem.FileSystem {
	fs := filesystem.NewMockFileSystem()
	_ = fs.WriteFile("go.mod", []byte("module benchmark"), 0644)
	_ = fs.WriteFile("cmd/promptengine/main.go", []byte("package main\nfunc main() {}\n"), 0644)
	_ = fs.WriteFile("internal/domain/payment/service.go", []byte("package payment\ntype Service struct{}\n"), 0644)
	_ = fs.WriteFile("internal/app/bootstrap.go", []byte("package app\n"), 0644)
	_ = fs.WriteFile("internal/filesystem/filesystem.go", []byte("package filesystem\n"), 0644)
	_ = fs.WriteFile("docs/Architecture.md", []byte("# Architecture\n\nBenchmark architecture."), 0644)
	_ = fs.WriteFile("docs/API.md", []byte("# API\n\nBenchmark API."), 0644)
	_ = fs.WriteFile("docs/Database.md", []byte("# Database\n\nBenchmark database."), 0644)
	_ = fs.WriteFile("docs/Decisions.md", []byte("# Decisions\n\nBenchmark decisions."), 0644)
	_ = fs.WriteFile("docs/Testing.md", []byte("# Testing\n\nBenchmark testing."), 0644)
	_ = fs.WriteFile("SECURITY.md", []byte("# Security\n\nBenchmark security."), 0644)
	_ = fs.WriteFile("playbook-manifest.json", []byte(`{"metadata":{"name":"bench","version":"1","schema_version":"2.0.0"}}`), 0644)
	return fs
}
