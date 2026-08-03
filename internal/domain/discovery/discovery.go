package discovery

import (
	"os"
	"path/filepath"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

// Finder locates the PromptEngine core playbooks path
type Finder struct {
	fs filesystem.FileSystem
}

func NewFinder(fs filesystem.FileSystem) *Finder {
	return &Finder{fs: fs}
}

func (f *Finder) Find(localConfigPath, localParamPath string) (string, error) {
	// 1. Env override
	if env := os.Getenv("PROMPTENGINE_PATH"); env != "" {
		if f.fs.Exists(env) {
			return env, nil
		}
	}

	// 2. Explicit argument parameter path
	if localParamPath != "" {
		if f.fs.Exists(localParamPath) {
			return localParamPath, nil
		}
	}

	// 3. Local configuration path
	if localConfigPath != "" {
		if f.fs.Exists(localConfigPath) {
			return localConfigPath, nil
		}
	}

	// 4. Relative directory crawler scan
	cwd, err := os.Getwd()
	if err == nil {
		curr := cwd
		for {
			target := filepath.Join(curr, "promptengine")
			if f.fs.Exists(target) {
				return target, nil
			}
			parent := filepath.Dir(curr)
			if parent == curr {
				break
			}
			curr = parent
		}
	}

	// 5. Global user dir fallback
	home, err := os.UserHomeDir()
	if err == nil {
		fallback := filepath.Join(home, ".promptengine", "core")
		if f.fs.Exists(fallback) {
			return fallback, nil
		}
	}

	return "", filepath.ErrBadPattern
}
