package container

import (
	"context"
	"io"
	"log/slog"

	"github.com/LordCodex/promptengine/internal/cache"
	"github.com/LordCodex/promptengine/internal/config"
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
)

type Options struct {
	Config *config.AppConfig
	Out    io.Writer
	Err    io.Writer
}

type Container struct {
	Config       *config.AppConfig
	Logger       *slog.Logger
	FS           filesystem.FileSystem
	Cache        *cache.Cache
	History      *history.Recorder
	Renderer     output.Renderer
	EventBus     *eventbus.EventBus
	Manifest     *manifest.Engine
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
}

func NewContainer(opts Options) (*Container, error) {
	cfg := opts.Config
	if cfg == nil {
		cfg = config.DefaultConfig()
	}

	logWriter := opts.Err
	if logWriter == nil {
		logWriter = io.Discard
	}

	level := slog.LevelInfo
	if cfg.CLI.Debug {
		level = slog.LevelDebug
	}

	mode := cfg.Mode
	if cfg.CLI.Debug {
		mode = "debug"
	} else if cfg.CLI.Verbose {
		mode = "development"
	}

	handlerOpts := &slog.HandlerOptions{Level: level}
	var handler slog.Handler
	if mode == "production" {
		handler = slog.NewJSONHandler(logWriter, handlerOpts)
	} else {
		handler = slog.NewTextHandler(logWriter, handlerOpts)
	}

	format := output.FormatText
	if cfg.CLI.JSON {
		format = output.FormatJSON
	}

	fs := &filesystem.OSFileSystem{BaseDir: "."}

	events := eventbus.NewEventBus()
	manifestEngine := manifest.NewEngineWithFS(fs)
	pluginRegistry := plugins.NewRegistryWithEvents(events)
	hookRegistry := hooks.NewRegistry(fs)
	hookRegistry.Attach(events)
	pluginInstaller := installer.NewLocalInstaller(fs, ".promptengine/plugins")
	agentPlatform := agents.NewPlatform(fs, events)
	personalPlatform := personal.NewPlatform(fs, events)
	intelligencePlatform := intelligence.NewPlatform(fs, events)

	workflowRegistry := workflows.NewRegistry()
	docsPlatform := docs.NewPlatform(fs, events, manifestEngine)
	qualityPlatform := quality.NewPlatform(fs, events)
	aiPlatform := ai.NewPlatform(nil, events)
	workflowEngine := workflows.NewEngine(fs, workflowRegistry, events)
	workflowEngine.RegisterHandler("discovery", workflows.StepHandlerFunc(func(ctx context.Context, step workflows.WorkflowStep, flow *workflows.FlowContext) (interface{}, error) {
		model, err := discovery.NewDefaultPipeline(events, manifestEngine).Execute(ctx, fs, ".")
		if err != nil {
			return nil, err
		}
		flow.Project = model
		return model, nil
	}))
	workflowEngine.RegisterHandler("context", workflows.StepHandlerFunc(func(ctx context.Context, step workflows.WorkflowStep, flow *workflows.FlowContext) (interface{}, error) {
		if flow.Project == nil {
			model, err := discovery.NewDefaultPipeline(events, manifestEngine).Execute(ctx, fs, ".")
			if err != nil {
				return nil, err
			}
			flow.Project = model
		}
		pkg, err := contextengine.NewEngine(fs, contextengine.WithManifestQuery(manifest.NewQueryEngine(manifestEngine)), contextengine.WithEventBus(events)).Build(ctx, contextengine.ContextRequest{
			TaskType: contextengine.TaskType(flow.TaskName),
			Project:  flow.Project,
			Budget:   contextengine.BudgetSmall,
		})
		if err != nil {
			return nil, err
		}
		flow.SelectedContext = pkg
		return pkg, nil
	}))
	for action, handler := range docsPlatform.WorkflowHandlers() {
		workflowEngine.RegisterHandler(action, handler)
	}

	return &Container{
		Config:       cfg,
		Logger:       slog.New(handler),
		FS:           fs,
		Cache:        cache.NewCache("."),
		History:      history.NewRecorder(fs),
		Renderer:     output.NewConfiguredRenderer(format, false, cfg.CLI.Verbose),
		EventBus:     events,
		Manifest:     manifestEngine,
		Discovery:    discovery.NewDefaultPipeline(events, manifestEngine),
		Context:      contextengine.NewEngine(fs, contextengine.WithManifestQuery(manifest.NewQueryEngine(manifestEngine)), contextengine.WithEventBus(events)),
		Workflow:     workflowEngine,
		Docs:         docsPlatform,
		Quality:      qualityPlatform,
		AI:           aiPlatform,
		Plugins:      pluginRegistry,
		Hooks:        hookRegistry,
		Installer:    pluginInstaller,
		Agents:       agentPlatform,
		Personal:     personalPlatform,
		Intelligence: intelligencePlatform,
	}, nil
}
