package discovery

import (
	"context"
	"path/filepath"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

// DocsStage maps documentation file existence and completeness percentages
type DocsStage struct{}

func (s *DocsStage) Name() string { return "documentation_stage" }
func (s *DocsStage) Execute(ctx context.Context, fs filesystem.FileSystem, pm *ProjectModel) error {
	// Standard project documents list
	docsMap := map[string]string{
		"PRD":             "PRD.md",
		"Architecture":    "docs/Architecture.md",
		"Database":        "docs/Database.md",
		"API":             "docs/API.md",
		"BusinessRules":   "docs/BusinessRules.md",
		"Roadmap":         "docs/Roadmap.md",
		"Troubleshooting": "docs/Troubleshooting.md",
	}

	for label, relativePath := range docsMap {
		fullPath := filepath.Join(pm.RootDir, relativePath)
		spec := DocSpec{
			Name:   label,
			Path:   relativePath,
			Exists: false,
		}

		if fs.Exists(fullPath) {
			spec.Exists = true
			// Calculate basic completeness: check word counts or layout sections
			data, err := fs.ReadFile(fullPath)
			if err == nil {
				length := len(data)
				if length > 500 {
					spec.Completeness = 100.0
				} else if length > 100 {
					spec.Completeness = 50.0
				} else {
					spec.Completeness = 10.0
				}
			}
		}

		pm.Docs[label] = spec
	}

	return nil
}
