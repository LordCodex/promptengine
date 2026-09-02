package app

import (
	"context"
	"io"

	"github.com/LordCodex/promptengine/internal/config"
	"github.com/LordCodex/promptengine/internal/container"
	"github.com/LordCodex/promptengine/internal/domain/agents"
	"github.com/LordCodex/promptengine/internal/domain/ai"
	contextengine "github.com/LordCodex/promptengine/internal/domain/context"
	"github.com/LordCodex/promptengine/internal/domain/discovery"
	"github.com/LordCodex/promptengine/internal/domain/docs"
	"github.com/LordCodex/promptengine/internal/domain/hooks"
	"github.com/LordCodex/promptengine/internal/domain/installer"
	"github.com/LordCodex/promptengine/internal/domain/intelligence"
	"github.com/LordCodex/promptengine/internal/domain/personal"
	"github.com/LordCodex/promptengine/internal/domain/plugins"
	"github.com/LordCodex/promptengine/internal/domain/quality"
	"github.com/LordCodex/promptengine/internal/domain/workflows"
	"github.com/LordCodex/promptengine/internal/eventbus"
	"github.com/LordCodex/promptengine/internal/filesystem"
	"github.com/LordCodex/promptengine/internal/history"
	"github.com/LordCodex/promptengine/internal/output"
	"github.com/LordCodex/promptengine/pkg/manifest"
	"github.com/LordCodex/promptengine/pkg/rulesources"
	"github.com/spf13/cobra"
)

type BootstrapOptions struct {
	Out    io.Writer
	Err    io.Writer
	Flags  config.Flags
	Args   []string
	Config string
}

type App struct {
	Config   *config.AppConfig
	FS       filesystem.FileSystem
	Renderer output.Renderer
	Out      io.Writer
	Err      io.Writer
	History  *history.Recorder

	EventBus     *eventbus.EventBus
	Container    *container.Container
	Manifest     *manifest.Engine
	ManifestQ    *manifest.QueryEngine
	RuleSources  *rulesources.Service
	Discovery    *discovery.Pipeline
	Context      *contextengine.Engine
	Workflow     *workflows.Engine
	Docs         *docs.Platform
	Quality      *quality.Platform
	AI           *ai.Platform
	Plugins      *plugins.Registry
	Hooks        *hooks.Registry
	Installer    *installer.LocalInstaller
	Agents       *agents.Platform
	Personal     *personal.Platform
	Intelligence *intelligence.Platform
	RootCmd      *cobra.Command
}

func Bootstrap(opts BootstrapOptions) (*App, error) {
	loader := config.NewConfigLoader()
	flags := opts.Flags
	if opts.Config != "" {
		flags.Config = opts.Config
	}

	cfg, err := loader.LoadWithFlags(".promptengine.yaml", "~/.promptengine/config.yaml", flags)
	if err != nil {
		return nil, err
	}

	c, err := container.NewContainer(container.Options{
		Config: cfg,
		Out:    opts.Out,
		Err:    opts.Err,
	})
	if err != nil {
		return nil, err
	}

	app := &App{
		Config:       c.Config,
		FS:           c.FS,
		Renderer:     c.Renderer,
		Out:          opts.Out,
		Err:          opts.Err,
		History:      c.History,
		EventBus:     c.EventBus,
		Container:    c,
		Manifest:     c.Manifest,
		ManifestQ:    manifest.NewQueryEngine(c.Manifest),
		RuleSources:  c.RuleSources,
		Discovery:    c.Discovery,
		Context:      c.Context,
		Workflow:     c.Workflow,
		Docs:         c.Docs,
		Quality:      c.Quality,
		AI:           c.AI,
		Plugins:      c.Plugins,
		Hooks:        c.Hooks,
		Installer:    c.Installer,
		Agents:       c.Agents,
		Personal:     c.Personal,
		Intelligence: c.Intelligence,
	}
	app.RootCmd = NewRootCommand(app)

	app.EventBus.Publish(eventbus.Event{Type: eventbus.ApplicationStarted, Message: "application started"})
	app.EventBus.Publish(eventbus.Event{Type: eventbus.ApplicationReady, Message: "application ready"})

	return app, nil
}

func (a *App) Execute(ctx context.Context, args []string) error {
	a.RootCmd.SetArgs(args)
	err := a.RootCmd.ExecuteContext(ctx)
	a.EventBus.Publish(eventbus.Event{Type: eventbus.ApplicationShutdown, Message: "application shutdown"})
	return err
}
