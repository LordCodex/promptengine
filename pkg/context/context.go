package context

import (
	contextpkg "github.com/LordCodex/promptengine/internal/domain/context"
	"github.com/LordCodex/promptengine/internal/filesystem"
	"github.com/LordCodex/promptengine/pkg/manifest"
)

type Builder struct {
	inner *contextpkg.Builder
}

func NewBuilder(fs filesystem.FileSystem, manifest *manifest.PlaybookManifest) *Builder {
	return &Builder{
		inner: contextpkg.NewBuilder(fs, manifest),
	}
}

func (b *Builder) BuildContext(task contextpkg.TaskType, extraFiles []string) ([]string, error) {
	return b.inner.BuildContext(task, extraFiles)
}
