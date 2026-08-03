package context

import (
	"context"

	"github.com/LordCodex/promptengine/internal/filesystem"
	"github.com/LordCodex/promptengine/pkg/manifest"
)

type TaskType string

const (
	TaskNewProject       TaskType = "new-project"
	TaskExistingProject  TaskType = "existing-project"
	TaskProjectMigration TaskType = "migrate-existing-project"
	TaskAddFeature       TaskType = "add-feature"
	TaskBugFix           TaskType = "bug-fix"
	TaskRefactor         TaskType = "refactor"
	TaskReview           TaskType = "review"
	TaskRelease          TaskType = "release"
	TaskDatabaseChanges  TaskType = "database-changes"
)

// Builder computes the minimum set of playbooks and documents files to load
type Builder struct {
	fs             filesystem.FileSystem
	manifest       *manifest.PlaybookManifest
	manifestEngine *manifest.Engine
}

func NewBuilder(fs filesystem.FileSystem, m *manifest.PlaybookManifest) *Builder {
	return &Builder{
		fs:       fs,
		manifest: m,
	}
}

func NewBuilderWithManifestEngine(fs filesystem.FileSystem, engine *manifest.Engine) *Builder {
	return &Builder{fs: fs, manifestEngine: engine}
}

func (b *Builder) BuildContext(task TaskType, affectedFiles []string) ([]string, error) {
	var resolved []string

	// 1. Core rule maps: always load core agents guidelines
	resolved = append(resolved, "AGENTS.md")

	// 2. Load manifest workflow playbooks mapping task key
	if b.manifest != nil {
		if mapping, ok := b.manifest.TaskMappings[string(task)]; ok {
			for _, pid := range mapping.RequiredPlaybookIDs {
				resolved = append(resolved, pid)
			}
		}
	}

	// 3. Scan affected files list to link specific database or api specs
	for _, f := range affectedFiles {
		if f == "Database" || f == "Database.md" {
			resolved = append(resolved, "docs/Database.md")
		}
		if f == "API" || f == "API.md" {
			resolved = append(resolved, "docs/API.md")
		}
	}

	return resolved, nil
}

func (b *Builder) BuildPackage(req ContextRequest) (*ContextPackage, error) {
	var query *manifest.QueryEngine
	if b.manifestEngine != nil {
		query = manifest.NewQueryEngine(b.manifestEngine)
	}
	return NewEngine(b.fs, WithManifestQuery(query)).Build(context.Background(), req)
}
