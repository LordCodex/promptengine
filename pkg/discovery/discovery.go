package discovery

import (
	"context"

	"github.com/LordCodex/promptengine/internal/domain/discovery"
	"github.com/LordCodex/promptengine/internal/filesystem"
)

// ProjectModel describes the detected project stack metadata
type ProjectModel = discovery.ProjectModel

// Pipeline is the public entry point for stack discovery
type Pipeline struct {
	inner *discovery.Pipeline
}

// NewPipeline instantiates and registers default detection stages
func NewPipeline() *Pipeline {
	p := discovery.NewPipeline()
	p.Register(
		&discovery.BaseStage{},
		&discovery.TechStage{},
		&discovery.DocsStage{},
		&discovery.ArchStage{},
		&discovery.PromptEngineStage{},
	)
	return &Pipeline{inner: p}
}

// Execute runs the discovery pipeline against a target folder path
func (p *Pipeline) Execute(ctx context.Context, fs filesystem.FileSystem, path string) (*ProjectModel, error) {
	return p.inner.Execute(ctx, fs, path)
}
