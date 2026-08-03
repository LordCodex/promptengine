package discovery

import (
	"context"
	"path/filepath"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

// BaseStage checks standard VCS and repository properties
type BaseStage struct{}

func (s *BaseStage) Name() string { return "base_stage" }
func (s *BaseStage) Execute(ctx context.Context, fs filesystem.FileSystem, pm *ProjectModel) error {
	// Check Git
	if fs.Exists(filepath.Join(pm.RootDir, ".git")) {
		pm.HasGit = true
	}

	// Check Docker
	if fs.Exists(filepath.Join(pm.RootDir, "Dockerfile")) || fs.Exists(filepath.Join(pm.RootDir, "docker-compose.yml")) {
		pm.HasDocker = true
	}

	return nil
}

// PromptEngineStage checks installation statuses
type PromptEngineStage struct{}

func (s *PromptEngineStage) Name() string { return "promptengine_stage" }
func (s *PromptEngineStage) Execute(ctx context.Context, fs filesystem.FileSystem, pm *ProjectModel) error {
	hasManifest := fs.Exists(filepath.Join(pm.RootDir, "playbook-manifest.json"))
	hasConfig := fs.Exists(filepath.Join(pm.RootDir, ".promptengine.json"))

	if hasManifest {
		pm.PromptEngine.Installed = true
		pm.PromptEngine.Version = "0.1.0"
	}

	if hasConfig {
		pm.PromptEngine.HasConfig = true
		pm.PromptEngine.ConfigVersion = "1"
	}

	// Check project constitution AGENTS.md
	agentsMD := filepath.Join(pm.RootDir, "AGENTS.md")
	if fs.Exists(agentsMD) {
		pm.PromptEngine.AgentsMDPresent = true
	} else {
		// Crawl common paths
		commonPaths := []string{".github/AGENTS.md", ".agents/AGENTS.md", "config/AGENTS.md"}
		for _, cp := range commonPaths {
			if fs.Exists(filepath.Join(pm.RootDir, cp)) {
				pm.PromptEngine.AgentsMDPresent = true
				break
			}
		}
	}

	return nil
}
