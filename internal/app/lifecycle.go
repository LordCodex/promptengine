package app

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/LordCodex/promptengine/internal/config"
	"github.com/LordCodex/promptengine/internal/domain/agents"
	"github.com/LordCodex/promptengine/internal/domain/discovery"
	"github.com/LordCodex/promptengine/internal/eventbus"
	"github.com/LordCodex/promptengine/internal/filesystem"
	"github.com/LordCodex/promptengine/internal/history"
	"github.com/LordCodex/promptengine/pkg/manifest"
	"github.com/spf13/cobra"
)

type LifecycleContext struct {
	Ctx      context.Context
	Cmd      *cobra.Command
	Out      io.Writer
	FS       filesystem.FileSystem
	Config   *config.AppConfig
	Model    *discovery.ProjectModel
	Manifest *manifest.Manifest
	Logger   *slog.Logger
}

type LifecycleRunner func(lc *LifecycleContext, args []string) error
type LifecycleWrapper func(runner LifecycleRunner) func(cmd *cobra.Command, args []string) error

func (a *App) EnforceLifecycle(runner LifecycleRunner) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		started := time.Now().UTC()
		var lifecycleErr error
		defer func() {
			a.recordHistory(cmd, started, lifecycleErr)
		}()
		logger := a.Container.Logger
		logger.Debug("lifecycle initialize")

		if a.Config == nil {
			lifecycleErr = fmt.Errorf("lifecycle: application configuration not loaded")
			return lifecycleErr
		}

		logger.Debug("lifecycle load configuration")
		logger.Debug("lifecycle initialize services")
		activeManifest, err := a.loadManifestForLifecycle()
		if err != nil {
			a.EventBus.Publish(eventbus.Event{
				Type:    eventbus.ManifestValidationFailed,
				Message: "manifest validation failed",
				Payload: err,
			})
			lifecycleErr = err
			return lifecycleErr
		}
		logger.Debug("lifecycle load plugins")
		if a.Plugins != nil {
			for id, raw := range a.Config.Plugins {
				if cfg, ok := raw.(map[string]interface{}); ok {
					a.Plugins.Configure(id, cfg)
				}
			}
		}
		if a.Agents != nil {
			for id, profile := range a.Config.Agents {
				if profile.InstructionFile == "" {
					continue
				}
				a.Agents.Register(agents.AgentProfile{ID: id, Name: id, InstructionFile: profile.InstructionFile, Format: profile.Format})
			}
		}
		var projectModel *discovery.ProjectModel
		if a.Discovery != nil {
			var err error
			projectModel, err = a.Discovery.Execute(cmd.Context(), a.FS, ".")
			if err != nil {
				lifecycleErr = err
				return lifecycleErr
			}
		}
		if err := a.activateAuthoritativeRules(projectModel); err != nil {
			lifecycleErr = err
			return lifecycleErr
		}
		if a.Manifest != nil {
			activeManifest = a.Manifest.ActiveManifest()
		}

		a.EventBus.Publish(eventbus.Event{
			Type:    eventbus.CommandStarted,
			Message: "command started",
			Payload: map[string]string{"command": cmd.Name()},
		})

		lc := &LifecycleContext{
			Ctx:      cmd.Context(),
			Cmd:      cmd,
			Out:      a.Out,
			FS:       a.FS,
			Config:   a.Config,
			Model:    projectModel,
			Manifest: activeManifest,
			Logger:   logger,
		}

		execErr := runner(lc, args)
		lifecycleErr = execErr
		a.EventBus.Publish(eventbus.Event{
			Type:    eventbus.CommandCompleted,
			Message: "command completed",
			Payload: map[string]interface{}{
				"command": cmd.Name(),
				"error":   execErr,
			},
		})

		logger.Debug("lifecycle render output")
		logger.Debug("lifecycle cleanup")
		return execErr
	}
}

func (a *App) recordHistory(cmd *cobra.Command, started time.Time, err error) {
	if a.History == nil || cmd == nil {
		return
	}
	status := "completed"
	errText := ""
	if err != nil {
		status = "failed"
		errText = err.Error()
	}
	entry := history.Entry{
		Command:   cmd.CommandPath(),
		Status:    status,
		StartedAt: started,
		EndedAt:   time.Now().UTC(),
		Error:     errText,
		Metadata:  map[string]string{"phase": "lifecycle"},
	}
	if recordErr := a.History.Record(entry); recordErr != nil {
		a.Container.Logger.Debug("audit history write failed", "error", recordErr)
	}
}

func (a *App) loadManifestForLifecycle() (*manifest.Manifest, error) {
	if a.Manifest == nil {
		return nil, nil
	}
	loader := manifest.NewLoader(a.FS)
	path, ok := loader.Discover(".")
	if !ok {
		return nil, nil
	}
	loaded, err := loader.Load(path)
	if err != nil {
		return nil, err
	}
	if err := a.Manifest.Register("project", manifest.SourceProject, loaded); err != nil {
		return nil, err
	}
	a.EventBus.Publish(eventbus.Event{
		Type:    eventbus.ManifestLoaded,
		Message: "manifest loaded",
		Payload: map[string]string{"path": path},
	})
	return a.Manifest.ActiveManifest(), nil
}
