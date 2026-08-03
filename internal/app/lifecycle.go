package app

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/LordCodex/promptengine/internal/config"
	"github.com/LordCodex/promptengine/internal/domain/discovery"
	"github.com/LordCodex/promptengine/internal/filesystem"
	"github.com/LordCodex/promptengine/pkg/manifest"
	"github.com/spf13/cobra"
)

// LifecycleContext carries all loaded state through the execution lifecycle.
type LifecycleContext struct {
	Ctx      context.Context
	Out      io.Writer
	FS       filesystem.FileSystem
	Config   *config.AppConfig
	Model    *discovery.ProjectModel
	Manifest *manifest.PlaybookManifest
	Logger   *slog.Logger
}

// LifecycleRunner wraps command execution to enforce loading order and validation rules.
type LifecycleRunner func(lc *LifecycleContext, args []string) error

// EnforceLifecycle executes the standardized command lifecycle phases.
func EnforceLifecycle(app *App, out io.Writer, runner LifecycleRunner) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// Phase 1: Load Config (already bootstrapped, but verified here)
		if app.Config == nil {
			return fmt.Errorf("lifecycle: application configuration not loaded")
		}

		// Phase 2: Discover Project
		app.Logger.Debug("Lifecycle: discovering project stack...")
		pm, err := app.Discovery.Execute(cmd.Context(), app.FS, ".")
		if err != nil {
			app.Logger.Warn("Lifecycle: project discovery warning", "err", err)
		}

		// Phase 3: Load Manifest (if playbook-manifest.json exists)
		var pmf *manifest.PlaybookManifest
		if app.FS.Exists("playbook-manifest.json") {
			app.Logger.Debug("Lifecycle: loading playbook manifest...")
			pmf, _ = manifest.Load("playbook-manifest.json")
		}

		// Phase 4: Construct Lifecycle Context
		lc := &LifecycleContext{
			Ctx:      cmd.Context(),
			Out:      out,
			FS:       app.FS,
			Config:   app.Config,
			Model:    pm,
			Manifest: pmf,
			Logger:   app.Logger,
		}

		// Phase 5: Execute command runner
		app.Logger.Debug("Lifecycle: executing command logic...")
		execErr := runner(lc, args)

		// Phase 6: Cleanup
		app.Logger.Debug("Lifecycle: cleanup completed")
		return execErr
	}
}
