package discovery

import (
	"context"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

// ArchStage scans folder structures to heuristically infer codebase architectures
type ArchStage struct{}

func (s *ArchStage) Name() string { return "architecture_stage" }
func (s *ArchStage) Execute(ctx context.Context, fs filesystem.FileSystem, pm *ProjectModel) error {
	// Check MVC (typical of Laravel, Rails, or Express routing structures)
	hasControllers := fs.Exists("app/Http/Controllers") || fs.Exists("controllers") || fs.Exists("src/controllers")
	hasModels := fs.Exists("app/Models") || fs.Exists("models") || fs.Exists("src/models")
	hasViews := fs.Exists("resources/views") || fs.Exists("views") || fs.Exists("templates")

	if hasControllers && hasModels && hasViews {
		pm.Architectures = append(pm.Architectures, ArchitectureInference{
			Style:      "MVC",
			Confidence: 0.85,
			Reason:     "Identified classic controllers, models, and views folders layout.",
		})
	}

	// Check Clean Architecture (typical Go/TS domains, entities, usecases layouts)
	hasDomain := fs.Exists("internal/domain") || fs.Exists("src/domain")
	hasUsecases := fs.Exists("internal/usecase") || fs.Exists("src/usecases") || fs.Exists("internal/app")
	hasAdapters := fs.Exists("internal/adapters") || fs.Exists("internal/filesystem")

	if hasDomain && (hasUsecases || hasAdapters) {
		pm.Architectures = append(pm.Architectures, ArchitectureInference{
			Style:      "Clean Architecture",
			Confidence: 0.90,
			Reason:     "Spotted domain logic modules segregated from filesystem adapters.",
		})
	}

	// Check Domain-Driven Design (DDD)
	hasBoundedContexts := fs.Exists("src/Domain") && fs.Exists("src/Infrastructure") && fs.Exists("src/Application")
	if hasBoundedContexts {
		pm.Architectures = append(pm.Architectures, ArchitectureInference{
			Style:      "Domain-Driven Design (DDD)",
			Confidence: 0.80,
			Reason:     "Spotted typical Application/Domain/Infrastructure boundaries structure.",
		})
	}

	// Fallback to basic monolithic layout
	if len(pm.Architectures) == 0 {
		pm.Architectures = append(pm.Architectures, ArchitectureInference{
			Style:      "Monolithic / Flat Layout",
			Confidence: 0.50,
			Reason:     "No specialized architectural folder structures discovered.",
		})
	}

	return nil
}
