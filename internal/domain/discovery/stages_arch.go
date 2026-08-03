package discovery

import (
	"context"
	"path/filepath"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

// ArchStage scans folder structures to heuristically infer codebase architectures
type ArchStage struct{}

func (s *ArchStage) Name() string { return "architecture_stage" }
func (s *ArchStage) Execute(ctx context.Context, fs filesystem.FileSystem, pm *ProjectModel) error {
	// Check MVC (typical of Laravel, Rails, or Express routing structures)
	hasControllers := fs.Exists(filepath.Join(pm.RootDir, "app/Http/Controllers")) || fs.Exists(filepath.Join(pm.RootDir, "controllers")) || fs.Exists(filepath.Join(pm.RootDir, "src/controllers"))
	hasModels := fs.Exists(filepath.Join(pm.RootDir, "app/Models")) || fs.Exists(filepath.Join(pm.RootDir, "models")) || fs.Exists(filepath.Join(pm.RootDir, "src/models"))
	hasViews := fs.Exists(filepath.Join(pm.RootDir, "resources/views")) || fs.Exists(filepath.Join(pm.RootDir, "views")) || fs.Exists(filepath.Join(pm.RootDir, "templates"))

	if hasControllers && hasModels && hasViews {
		pm.Architectures = append(pm.Architectures, ArchitectureInference{
			Style:      "MVC",
			Confidence: 0.85,
			Reason:     "Identified classic controllers, models, and views folders layout.",
		})
	}

	// Check Clean Architecture (typical Go/TS domains, entities, usecases layouts)
	hasDomain := fs.Exists(filepath.Join(pm.RootDir, "internal/domain")) || fs.Exists(filepath.Join(pm.RootDir, "src/domain"))
	hasUsecases := fs.Exists(filepath.Join(pm.RootDir, "internal/usecase")) || fs.Exists(filepath.Join(pm.RootDir, "src/usecases")) || fs.Exists(filepath.Join(pm.RootDir, "internal/app"))
	hasAdapters := fs.Exists(filepath.Join(pm.RootDir, "internal/adapters")) || fs.Exists(filepath.Join(pm.RootDir, "internal/filesystem"))

	if hasDomain && (hasUsecases || hasAdapters) {
		pm.Architectures = append(pm.Architectures, ArchitectureInference{
			Style:      "Clean Architecture",
			Confidence: 0.90,
			Reason:     "Spotted domain logic modules segregated from filesystem adapters.",
		})
	}

	// Check Domain-Driven Design (DDD)
	hasBoundedContexts := fs.Exists(filepath.Join(pm.RootDir, "src/Domain")) && fs.Exists(filepath.Join(pm.RootDir, "src/Infrastructure")) && fs.Exists(filepath.Join(pm.RootDir, "src/Application"))
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
