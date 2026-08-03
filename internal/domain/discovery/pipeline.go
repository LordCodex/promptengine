package discovery

import (
	"context"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

// Stage is a modular discovery run phase
type Stage interface {
	Name() string
	Execute(ctx context.Context, fs filesystem.FileSystem, pm *ProjectModel) error
}

// Pipeline coordinates stage executions
type Pipeline struct {
	stages []Stage
}

func NewPipeline() *Pipeline {
	return &Pipeline{
		stages: make([]Stage, 0),
	}
}

// Register appends a new execution stage to the pipeline
func (p *Pipeline) Register(stages ...Stage) {
	p.stages = append(p.stages, stages...)
}

func (p *Pipeline) Execute(ctx context.Context, fs filesystem.FileSystem, rootPath string) (*ProjectModel, error) {
	pm := NewProjectModel(rootPath)

	for _, stage := range p.stages {
		// Run stage audit
		if err := stage.Execute(ctx, fs, pm); err != nil {
			return nil, err
		}
	}

	// Dynamic project categorization consolidation
	p.consolidateClassifications(pm)

	return pm, nil
}

func (p *Pipeline) consolidateClassifications(pm *ProjectModel) {
	// Greenfield vs Existing
	if len(pm.Languages) == 0 && len(pm.Frameworks) == 0 {
		pm.Classifications = append(pm.Classifications, ClassGreenfield)
	} else {
		pm.Classifications = append(pm.Classifications, ClassExisting)
	}

	// PromptEngine status
	if pm.PromptEngine.Installed {
		pm.Classifications = append(pm.Classifications, ClassPromptEngine)
	} else if len(pm.Languages) > 0 {
		pm.Classifications = append(pm.Classifications, ClassNonPromptEngine)
	}

	// Library vs CLI vs Backend vs SSR
	for _, f := range pm.Frameworks {
		switch f {
		case "Laravel", "NestJS", "Express":
			pm.Classifications = append(pm.Classifications, ClassBackendAPI)
		case "Nuxt", "Next.js":
			pm.Classifications = append(pm.Classifications, ClassSSRApplication)
		case "Vue", "React":
			pm.Classifications = append(pm.Classifications, ClassFrontendSPA)
		case "Flutter":
			pm.Classifications = append(pm.Classifications, ClassMobileApplication)
		}
	}
}
