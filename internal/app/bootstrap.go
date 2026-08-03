package app

import (
	"context"
	"io"
	"log/slog"

	"github.com/LordCodex/promptengine/internal/cache"
	"github.com/LordCodex/promptengine/internal/config"
	"github.com/LordCodex/promptengine/internal/container"
	contextpkg "github.com/LordCodex/promptengine/internal/domain/context"
	"github.com/LordCodex/promptengine/internal/domain/discovery"
	docengine "github.com/LordCodex/promptengine/internal/domain/docs"
	"github.com/LordCodex/promptengine/internal/domain/docs/generator"
	docsync "github.com/LordCodex/promptengine/internal/domain/docs/sync"
	healthpkg "github.com/LordCodex/promptengine/internal/domain/health"
	"github.com/LordCodex/promptengine/internal/domain/hooks"
	"github.com/LordCodex/promptengine/internal/domain/installer"
	"github.com/LordCodex/promptengine/internal/domain/plugins"
	"github.com/LordCodex/promptengine/internal/domain/prompts"
	"github.com/LordCodex/promptengine/internal/domain/quality/audit"
	"github.com/LordCodex/promptengine/internal/domain/quality/compliance"
	"github.com/LordCodex/promptengine/internal/domain/quality/doctor"
	"github.com/LordCodex/promptengine/internal/domain/quality/fix"
	"github.com/LordCodex/promptengine/internal/domain/quality/report"
	"github.com/LordCodex/promptengine/internal/domain/quality/validation"
	"github.com/LordCodex/promptengine/internal/domain/review"
	"github.com/LordCodex/promptengine/internal/domain/updater"
	"github.com/LordCodex/promptengine/internal/domain/workflows"
	"github.com/LordCodex/promptengine/internal/eventbus"
	"github.com/LordCodex/promptengine/internal/filesystem"
	"github.com/LordCodex/promptengine/internal/output"
	"github.com/LordCodex/promptengine/internal/telemetry"
	"github.com/spf13/cobra"
)

// App contains wired application state, container services, and all engine references
type App struct {
	Config   *config.AppConfig
	FS       filesystem.FileSystem
	Logger   *slog.Logger
	Renderer output.Renderer

	// Platform Caching, Telemetry, and Event Bus
	Cache     *cache.Cache
	Telemetry *telemetry.Telemetry
	EventBus  *eventbus.EventBus
	Container *container.Container

	// Quality Platform
	Doctor     *doctor.DoctorEngine
	Health     *healthpkg.Registry
	Reviewer   *review.Registry
	Auditor    *audit.AuditEngine
	Compliance *compliance.ComplianceEngine
	Validator  *validation.Registry
	Fixer      *fix.FixEngine
	Reporter   *report.RendererRegistry

	// Documentation Platform
	DocRegistry   *docengine.DocRegistry
	DocEngine     *docengine.Engine
	GenRegistry   *generator.GeneratorRegistry
	PromptReg     *prompts.PromptRegistry
	PromptBuilder *prompts.PromptBuilder
	SyncEngine    *docsync.SyncEngine

	// Extensibility
	PluginRegistry *plugins.Registry
	HookRegistry   *hooks.Registry
	Installer      *installer.LocalInstaller
	Updater        *updater.UpdateEngine

	// Workflows
	WorkflowRegistry *workflows.Registry
	WorkflowEngine   *workflows.Engine

	// Discovery & Context
	Discovery      *discovery.Pipeline
	ContextBuilder *contextpkg.Builder

	RootCmd *cobra.Command
}

// Bootstrap wires all engines and instantiates the CLI command router
func Bootstrap(out io.Writer, verbose bool) (*App, error) {
	// Initialize the Service Container
	c, err := container.NewContainer(verbose)
	if err != nil {
		return nil, err
	}

	// ── Quality Platform ────────────────────────────────────────────────────
	doctorEngine := doctor.NewDoctorEngine()
	healthReg := healthpkg.NewRegistry()
	reviewReg := review.NewRegistry()
	review.RegisterDefaultReviewers(reviewReg)
	auditEngine := audit.NewAuditEngine()
	complianceEngine := compliance.NewComplianceEngine()
	validationReg := validation.NewRegistry()
	fixEngine := fix.NewFixEngine()
	reportReg := report.NewRendererRegistry()

	// ── Documentation Platform ─────────────────────────────────────────────
	docReg := docengine.NewDocRegistry()
	docEng := docengine.NewEngine(c.FS, docReg)
	genReg := generator.NewGeneratorRegistry()
	generator.RegisterDefaults(genReg)

	promptReg := prompts.NewPromptRegistry()
	promptBld := prompts.NewPromptBuilder(promptReg)

	changeDetector := docsync.NewChangeDetector()
	syncEng := docsync.NewSyncEngine(changeDetector)

	// ── Extensibility ──────────────────────────────────────────────────────
	pluginReg := plugins.NewRegistry()
	hookReg := hooks.NewRegistry(c.FS)
	localInstaller := installer.NewLocalInstaller(c.FS, ".promptengine/plugins")
	updateEngine := updater.NewUpdateEngine()

	// ── Workflows ──────────────────────────────────────────────────────────
	workflowsReg := workflows.NewRegistry()
	workflowsPub := workflows.NewPublisher()
	workflowsEngine := workflows.NewEngine(c.FS, workflowsReg, workflowsPub)

	// ── Discovery ──────────────────────────────────────────────────────────
	pipeline := discovery.NewPipeline()
	pipeline.Register(
		&discovery.BaseStage{},
		&discovery.TechStage{},
		&discovery.DocsStage{},
		&discovery.ArchStage{},
		&discovery.PromptEngineStage{},
	)

	// Context Builder is wired to a manifest at command execution time
	ctxBuilder := contextpkg.NewBuilder(c.FS, nil)

	// ── Commands ───────────────────────────────────────────────────────────
	builder := NewCommandBuilder(out)
	rootCmd := builder.BuildRoot()

	app := &App{
		Config:    c.Config,
		FS:        c.FS,
		Logger:    c.Logger,
		Renderer:  c.Renderer,
		Cache:     c.Cache,
		Telemetry: c.Telemetry,
		EventBus:  c.EventBus,
		Container: c,

		Doctor:     doctorEngine,
		Health:     healthReg,
		Reviewer:   reviewReg,
		Auditor:    auditEngine,
		Compliance: complianceEngine,
		Validator:  validationReg,
		Fixer:      fixEngine,
		Reporter:   reportReg,

		DocRegistry:   docReg,
		DocEngine:     docEng,
		GenRegistry:   genReg,
		PromptReg:     promptReg,
		PromptBuilder: promptBld,
		SyncEngine:    syncEng,

		PluginRegistry: pluginReg,
		HookRegistry:   hookReg,
		Installer:      localInstaller,
		Updater:        updateEngine,

		WorkflowRegistry: workflowsReg,
		WorkflowEngine:   workflowsEngine,

		Discovery:      pipeline,
		ContextBuilder: ctxBuilder,

		RootCmd: rootCmd,
	}

	// ── Primary Commands ──────────────────────────────────────────────────
	rootCmd.AddCommand(builder.BuildInit(app))
	rootCmd.AddCommand(builder.BuildMigrate(app))
	rootCmd.AddCommand(builder.BuildDoctor(app))
	rootCmd.AddCommand(builder.BuildAnalyze(app))
	rootCmd.AddCommand(builder.BuildScan(app))
	rootCmd.AddCommand(builder.BuildReview(app))
	rootCmd.AddCommand(builder.BuildSync(app))
	rootCmd.AddCommand(builder.BuildUpdate(app))
	rootCmd.AddCommand(builder.BuildPrompt(app))
	rootCmd.AddCommand(builder.BuildContext(app))
	rootCmd.AddCommand(builder.BuildDocs(app))
	rootCmd.AddCommand(builder.BuildGenerate(app))
	rootCmd.AddCommand(builder.BuildConfig(app))
	rootCmd.AddCommand(builder.BuildInstall(app))
	rootCmd.AddCommand(builder.BuildPlugins(app))
	rootCmd.AddCommand(builder.BuildHooks(app))
	rootCmd.AddCommand(builder.BuildHealth(app))
	rootCmd.AddCommand(builder.BuildAudit(app))
	rootCmd.AddCommand(builder.BuildVersion())
	rootCmd.AddCommand(builder.BuildCompletion())

	// ── Missing Command Platform Additions ───────────────────────────────
	rootCmd.AddCommand(builder.BuildDetect(app))
	rootCmd.AddCommand(builder.BuildValidate(app))
	rootCmd.AddCommand(builder.BuildWorkflow(app))
	rootCmd.AddCommand(builder.BuildUninstall(app))
	rootCmd.AddCommand(builder.BuildPlugin(app))
	rootCmd.AddCommand(builder.BuildHook(app))
	rootCmd.AddCommand(builder.BuildManifest(app))
	rootCmd.AddCommand(builder.BuildStatus(app))
	rootCmd.AddCommand(builder.BuildList(app))
	rootCmd.AddCommand(builder.BuildClean(app))
	rootCmd.AddCommand(builder.BuildRepair(app))
	rootCmd.AddCommand(builder.BuildBootstrap(app))
	rootCmd.AddCommand(builder.BuildAI(app))
	rootCmd.AddCommand(builder.BuildTemplate(app))
	rootCmd.AddCommand(builder.BuildReport(app))

	return app, nil
}

func (a *App) Execute(ctx context.Context, args []string) error {
	a.RootCmd.SetArgs(args)
	return a.RootCmd.ExecuteContext(ctx)
}
