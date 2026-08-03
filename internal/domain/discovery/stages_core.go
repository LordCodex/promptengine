package discovery

import (
	"context"
	"path/filepath"

	"github.com/LordCodex/promptengine/internal/filesystem"
	"github.com/LordCodex/promptengine/internal/version"
)

// BaseStage checks standard VCS and repository properties
type BaseStage struct{}

func (s *BaseStage) Name() string { return "base_stage" }
func (s *BaseStage) Execute(ctx context.Context, fs filesystem.FileSystem, pm *ProjectModel) error {
	if len(pm.Repository.Files) == 0 && len(pm.Repository.Directories) == 0 {
		scanned, err := NewRepositoryScanner().Scan(ctx, fs, pm.RootDir)
		if err != nil {
			return err
		}
		pm.Repository = scanned
	}
	if fs.Exists(filepath.Join(pm.RootDir, ".git")) || hasAnyPath(pm.Repository.Directories, ".git") {
		pm.HasGit = true
	}

	if fs.Exists(filepath.Join(pm.RootDir, "Dockerfile")) || fs.Exists(filepath.Join(pm.RootDir, "docker-compose.yml")) || hasAnyPath(pm.Repository.Files, "Dockerfile", "docker-compose.yml", "docker-compose.yaml") {
		pm.HasDocker = true
	}

	return nil
}

type RepositoryScanStage struct {
	Scanner *RepositoryScanner
}

func (s *RepositoryScanStage) Name() string { return "repository_scan" }
func (s *RepositoryScanStage) Execute(ctx context.Context, fs filesystem.FileSystem, pm *ProjectModel) error {
	scanner := s.Scanner
	if scanner == nil {
		scanner = NewRepositoryScanner()
	}
	info, err := scanner.Scan(ctx, fs, pm.RootDir)
	if err != nil {
		return err
	}
	pm.Repository = info
	pm.Project.RootPath = info.RootPath
	pm.RootDir = info.RootPath
	return nil
}

// PromptEngineStage checks installation statuses
type PromptEngineStage struct{}

func (s *PromptEngineStage) Name() string { return "promptengine_stage" }
func (s *PromptEngineStage) Execute(ctx context.Context, fs filesystem.FileSystem, pm *ProjectModel) error {
	hasManifest := fs.Exists(filepath.Join(pm.RootDir, "playbook-manifest.json"))
	hasConfig := fs.Exists(filepath.Join(pm.RootDir, ".promptengine.yaml")) || fs.Exists(filepath.Join(pm.RootDir, ".promptengine.json"))

	if hasManifest {
		pm.PromptEngine.Installed = true
		pm.PromptEngine.Version = version.Version
	}

	if hasConfig {
		pm.PromptEngine.HasConfig = true
		pm.PromptEngine.ConfigVersion = version.Version
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
