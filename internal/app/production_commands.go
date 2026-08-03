package app

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/LordCodex/promptengine/internal/config"
	"github.com/LordCodex/promptengine/internal/domain/agents"
	"github.com/LordCodex/promptengine/internal/domain/ai"
	contextengine "github.com/LordCodex/promptengine/internal/domain/context"
	"github.com/LordCodex/promptengine/internal/domain/docs"
	"github.com/LordCodex/promptengine/internal/domain/installer"
	"github.com/LordCodex/promptengine/internal/domain/intelligence"
	"github.com/LordCodex/promptengine/internal/domain/personal"
	"github.com/LordCodex/promptengine/internal/domain/plugins"
	"github.com/LordCodex/promptengine/internal/domain/quality/fix"
	"github.com/LordCodex/promptengine/internal/domain/updater"
	"github.com/LordCodex/promptengine/internal/domain/workflows"
	apperrors "github.com/LordCodex/promptengine/internal/errors"
	"github.com/LordCodex/promptengine/internal/eventbus"
	"github.com/LordCodex/promptengine/internal/filesystem"
	"github.com/LordCodex/promptengine/internal/output"
	"github.com/LordCodex/promptengine/pkg/manifest"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewProductionCommands(app *App) []*cobra.Command {
	return []*cobra.Command{
		newInitCommand(app),
		newMigrateCommand(app),
		newScanCommand(app),
		newContextCommand(app),
		newWorkflowCommand(app),
		newDocsCommand(app),
		newQualityCommand(app, "doctor"),
		newQualityCommand(app, "review"),
		newQualityCommand(app, "audit"),
		newQualityCommand(app, "health"),
		newPromptCommand(app),
		newAICommand(app),
		newPluginCommand(app),
		newPluginsAliasCommand(app),
		newHooksCommand(app),
		newInstallCommand(app),
		newUpdateCommand(app),
		newAnalyzeCommand(app),
		newSyncCommand(app),
		newAgentsCommand(app),
		newProfileCommand(app),
		newTaskCommand(app),
		newVerifyCommand(app),
		newMemoryCommand(app),
		newInsightsCommand(app),
		newDecisionsCommand(app),
		newImpactCommand(app),
		newConfigCommand(app),
		newCompletionCommand(app),
	}
}

type statusResult struct {
	Status         string   `json:"status" yaml:"status"`
	Message        string   `json:"message,omitempty" yaml:"message,omitempty"`
	Files          []string `json:"files,omitempty" yaml:"files,omitempty"`
	Recommendation string   `json:"recommendation,omitempty" yaml:"recommendation,omitempty"`
}

func renderCLI(app *App, lc *LifecycleContext, data interface{}, human func() string) error {
	if configured, ok := app.Renderer.(*output.ConfiguredRenderer); ok && configured.Format != output.FormatText && configured.Format != output.FormatHuman {
		return app.Renderer.Render(app.Out, data)
	}
	if app.Config.CLI.JSON {
		return app.Renderer.Render(app.Out, data)
	}
	if human != nil {
		fmt.Fprintln(app.Out, human())
		return nil
	}
	return app.Renderer.Render(app.Out, data)
}

func newInitCommand(app *App) *cobra.Command {
	var overwrite bool
	var agentProfiles []string
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize PromptEngine in this project",
		RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
			files := []string{}
			if overwrite || !lc.FS.Exists(".promptengine.yaml") {
				data := []byte("version: 1.0.0\nmode: production\nproject:\n  name: PromptEngineProject\n  version: 1.0.0\n  promptengine_path: ./promptengine\n")
				if err := lc.FS.WriteFile(".promptengine.yaml", data, 0644); err != nil {
					return appErr("initialize project configuration", err, "Check write permissions in the project root.")
				}
				files = append(files, ".promptengine.yaml")
			}
			if overwrite || !lc.FS.Exists("playbook-manifest.json") {
				m := manifest.Manifest{}
				m.Metadata.Name = lc.Config.Project.Name
				m.Metadata.Version = lc.Config.Project.Version
				m.Metadata.SchemaVersion = manifest.SupportedSchemaVersion
				m.Metadata.GeneratedAt = time.Now().UTC()
				data, err := json.MarshalIndent(m, "", "  ")
				if err != nil {
					return appErr("initialize project manifest", err, "Retry initialization.")
				}
				if err := lc.FS.WriteFile("playbook-manifest.json", data, 0644); err != nil {
					return appErr("initialize project manifest", err, "Check write permissions in the project root.")
				}
				files = append(files, "playbook-manifest.json")
			}
			agentTargets := expandAgentProfiles(agentProfiles)
			if len(agentTargets) > 0 {
				if err := lc.FS.MkdirAll(".promptengine", 0755); err != nil {
					return appErr("initialize PromptEngine directory", err, "Check write permissions in the project root.")
				}
				files = append(files, ".promptengine/")
				for _, profile := range agentTargets {
					profile = strings.TrimSpace(profile)
					if profile == "" {
						continue
					}
					prefs := personalPrefs(app)
					generated, err := app.Agents.Generate(agents.InstructionRequest{Profile: profile, Project: lc.Model, Manifest: lc.Manifest, Preferences: prefs})
					if err != nil {
						return appErr("generate agent instructions", err, "Use a supported agent profile or run promptengine agents list.")
					}
					files = append(files, generated.Path)
				}
			}
			if app.Workflow != nil {
				flow := workflows.NewFlowContext("new-project")
				flow.Project = lc.Model
				_, _ = app.Workflow.Execute(lc.Ctx, "new-project", flow)
			}
			result := statusResult{Status: "initialized", Files: files}
			if len(files) == 0 {
				result.Status = "ready"
				result.Message = "PromptEngine is already initialized."
			}
			return renderCLI(app, lc, result, func() string {
				if len(files) == 0 {
					return result.Message
				}
				return "Initialized PromptEngine:\n" + strings.Join(prefixLines(files, "  "), "\n")
			})
		}),
	}
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "replace existing PromptEngine files")
	cmd.Flags().StringSliceVar(&agentProfiles, "agents", nil, "agent instruction profiles to create: codex, claude, codex-md, cursor, windsurf, all")
	return cmd
}

func newMigrateCommand(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Upgrade existing PromptEngine project configuration",
		RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
			if err := lc.Config.Migrate(); err != nil {
				return appErr("migrate configuration", err, "Review the project configuration file.")
			}
			report, err := app.Quality.Validate(lc.Ctx)
			if err != nil {
				return appErr("validate migrated project", err, "Run promptengine doctor for details.")
			}
			return renderCLI(app, lc, map[string]interface{}{"status": "migrated", "version": config.CurrentVersion, "validation": report}, func() string {
				return fmt.Sprintf("Migration complete\nVersion: %s\nValidation findings: %d", config.CurrentVersion, len(report.Findings))
			})
		}),
	}
}

func newScanCommand(app *App) *cobra.Command {
	return &cobra.Command{
		Use:     "scan",
		Aliases: []string{"detect"},
		Short:   "Analyze project structure and technology stack",
		RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
			return renderCLI(app, lc, lc.Model, func() string {
				return fmt.Sprintf("Project: %s\nRoot: %s\nLanguages: %s\nFrameworks: %s\nDatabases: %s\nArchitecture: backend=%t frontend=%t mobile=%t",
					lc.Model.Project.Name, lc.Model.RootDir, join(lc.Model.Languages), join(lc.Model.Frameworks), join(lc.Model.Databases), lc.Model.Architecture.Backend, lc.Model.Architecture.Frontend, lc.Model.Architecture.Mobile)
			})
		}),
	}
}

func newContextCommand(app *App) *cobra.Command {
	var task, workflow, intent, budget string
	var maxBytes int
	cmd := &cobra.Command{
		Use:   "context",
		Short: "Build an optimized context package",
		RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
			pkg, err := app.Context.Build(lc.Ctx, contextengine.ContextRequest{
				TaskType:     contextengine.TaskType(task),
				WorkflowType: workflow,
				Project:      lc.Model,
				UserIntent:   intent,
				MaxBytes:     maxBytes,
				Budget:       contextengine.BudgetType(budget),
			})
			if err != nil {
				return appErr("build context package", err, "Reduce context limits or verify discovery output.")
			}
			return renderCLI(app, lc, pkg, func() string {
				return "Selected Context:\n" + strings.Join(prefixLines(append(pkg.SelectedFiles, pkg.SelectedDocs...), "  "), "\n") + "\n\nReason:\n" + strings.Join(pkg.Reasoning, "\n")
			})
		}),
	}
	cmd.Flags().StringVar(&task, "task", "feature", "task type")
	cmd.Flags().StringVar(&workflow, "workflow", "", "workflow type")
	cmd.Flags().StringVar(&intent, "intent", "", "user intent")
	cmd.Flags().StringVar(&budget, "budget", string(contextengine.BudgetMedium), "context budget")
	addContextLimitFlags(cmd, &maxBytes)
	cmd.AddCommand(newContextExportCommand(app))
	return cmd
}

func newContextExportCommand(app *App) *cobra.Command {
	var task, intent, budget, format, agent string
	var maxBytes int
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export optimized context for an AI coding agent",
		RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
			pkg, err := app.Context.Build(lc.Ctx, contextengine.ContextRequest{
				TaskType:   contextengine.TaskType(task),
				Project:    lc.Model,
				UserIntent: intent,
				Budget:     contextengine.BudgetType(budget),
				MaxBytes:   maxBytes,
			})
			if err != nil {
				return appErr("build context export", err, "Reduce context limits or verify discovery output.")
			}
			export, err := app.Agents.ExportContext(agents.ContextExportRequest{Task: task, Agent: agent, Format: format, Package: pkg})
			if err != nil {
				return appErr("export context", err, "Use markdown, json, or yaml.")
			}
			return renderCLI(app, lc, export, func() string {
				return fmt.Sprintf("Context exported\nFile: %s\nAgent: %s\nEstimated context size: %d bytes", export.File, export.Agent, export.EstimatedContextBytes)
			})
		}),
	}
	cmd.Flags().StringVar(&task, "task", "feature", "task type")
	cmd.Flags().StringVar(&intent, "intent", "", "user intent")
	cmd.Flags().StringVar(&budget, "budget", string(contextengine.BudgetMedium), "context budget")
	cmd.Flags().StringVar(&format, "format", "markdown", "context export format: markdown, json, yaml")
	cmd.Flags().StringVar(&agent, "agent", "generic", "agent format: claude, codex, cursor, windsurf, generic")
	addContextLimitFlags(cmd, &maxBytes)
	return cmd
}

func newWorkflowCommand(app *App) *cobra.Command {
	var id string
	cmd := &cobra.Command{
		Use:   "workflow",
		Short: "Execute a registered workflow",
		RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
			if id == "" && len(args) > 0 {
				id = args[0]
			}
			if id == "" {
				id = "feature-implementation"
			}
			if lc.Manifest != nil {
				app.Workflow.LoadFromManifest(lc.Manifest)
			}
			flow := workflows.NewFlowContext(id)
			flow.Project = lc.Model
			exec, err := app.Workflow.Execute(lc.Ctx, id, flow)
			if err != nil {
				return appErr("execute workflow", err, "Verify the workflow is registered and all step handlers are available.")
			}
			return renderCLI(app, lc, exec, func() string {
				return fmt.Sprintf("Workflow %s: %s\nSteps completed: %d", exec.WorkflowID, exec.Status, len(exec.Results))
			})
		}),
	}
	cmd.Flags().StringVar(&id, "id", "", "workflow identifier")
	return cmd
}

func newDocsCommand(app *App) *cobra.Command {
	var docID string
	var overwrite bool
	cmd := &cobra.Command{Use: "docs", Short: "Generate, validate, and sync documentation", RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
		doc, err := app.Docs.Generate(lc.Ctx, docs.GenerationRequest{DocumentID: normalizeID(docID), Project: lc.Model, Manifest: lc.Manifest, Overwrite: overwrite})
		if err != nil {
			return appErr("generate documentation", err, "Check template definitions and document identifiers.")
		}
		return renderCLI(app, lc, doc, func() string { return fmt.Sprintf("Generated %s\n%s", doc.Name, doc.Path) })
	})}
	generate := &cobra.Command{Use: "generate", Short: "Generate project documentation", RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
		doc, err := app.Docs.Generate(lc.Ctx, docs.GenerationRequest{DocumentID: normalizeID(docID), Project: lc.Model, Manifest: lc.Manifest, Overwrite: overwrite})
		if err != nil {
			return appErr("generate documentation", err, "Check template definitions and document identifiers.")
		}
		return renderCLI(app, lc, doc, func() string { return fmt.Sprintf("Generated %s\n%s", doc.Name, doc.Path) })
	})}
	generate.Flags().StringVar(&docID, "doc", "architecture", "document id")
	generate.Flags().BoolVar(&overwrite, "overwrite", false, "replace existing documents")
	cmd.Flags().StringVarP(&docID, "doc-name", "d", "architecture", "document id")
	cmd.Flags().BoolVarP(&overwrite, "force", "f", false, "replace existing documents")
	validate := &cobra.Command{Use: "validate", Short: "Validate project documentation", RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
		report, err := app.Docs.Validate(lc.Model)
		if err != nil {
			return appErr("validate documentation", err, "Fix broken document references or missing sections.")
		}
		return renderCLI(app, lc, report, func() string {
			return fmt.Sprintf("Documentation Status\nDocuments: %d\nFindings: %d", len(report.Documents), len(report.Findings))
		})
	})}
	var changed []string
	var dryRun bool
	syncCmd := &cobra.Command{Use: "sync", Short: "Detect documentation drift", RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
		result := app.Docs.Sync(changed, dryRun)
		return renderCLI(app, lc, result, func() string {
			return fmt.Sprintf("Documentation sync\nPending: %d\nRecommendations: %d", len(result.Pending), len(result.Recommendations))
		})
	})}
	syncCmd.Flags().StringSliceVar(&changed, "changed", nil, "changed files")
	syncCmd.Flags().BoolVar(&dryRun, "dry-run", true, "report changes without writing")
	cmd.AddCommand(generate, validate, syncCmd)
	return cmd
}

func newQualityCommand(app *App, name string) *cobra.Command {
	var fixIssues bool
	var path string
	cmd := &cobra.Command{
		Use:   name,
		Short: qualityShort(name),
		RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
			if name == "doctor" && fixIssues {
				engine := fix.NewFixEngine()
				results := []fix.FixResult{
					engine.Apply("create-promptengine-dir", lc.FS, false),
					engine.Apply("create-docs-dir", lc.FS, false),
				}
				return renderCLI(app, lc, results, func() string { return fmt.Sprintf("Applied doctor fixes: %d", len(results)) })
			}
			_ = path
			var report interface{ Human() string }
			var err error
			switch name {
			case "doctor":
				report, err = app.Quality.Validate(lc.Ctx)
			case "review":
				report, err = app.Quality.Review(lc.Ctx)
			case "audit":
				report, err = app.Quality.Audit(lc.Ctx)
			case "health":
				report, err = app.Quality.Health(lc.Ctx)
			}
			if err != nil {
				return appErr("run "+name, err, "Review the findings and rerun the command.")
			}
			return renderCLI(app, lc, report, func() string { return report.Human() })
		}),
	}
	if name == "doctor" {
		cmd.Flags().BoolVar(&fixIssues, "fix", false, "apply safe automated fixes")
	}
	if name == "review" {
		cmd.Flags().StringVar(&path, "path", ".", "path to review")
	}
	return cmd
}

func newPromptCommand(app *App) *cobra.Command {
	var task, workflow, request, provider, model string
	var format, client, outPath, budget string
	var copyPrompt bool
	var maxBytes int
	cmd := &cobra.Command{
		Use:   "prompt",
		Short: "Generate an AI-ready prompt for external AI clients",
		RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
			if request == "" {
				request = strings.Join(args, " ")
			}
			if workflow != "" && task == "feature" {
				task = workflow
			}
			pkg, err := app.Context.Build(lc.Ctx, contextengine.ContextRequest{
				TaskType:     contextengine.TaskType(task),
				WorkflowType: workflow,
				Project:      lc.Model,
				UserIntent:   request,
				Budget:       contextengine.BudgetType(budget),
				MaxBytes:     maxBytes,
			})
			if err != nil {
				return appErr("build prompt context", err, "Verify project discovery output.")
			}
			compiled, err := app.AI.Compile(ai.CompileInput{Provider: provider, Model: model, UserRequest: request, ContextPackage: pkg, Manifest: lc.Manifest})
			if err != nil {
				return appErr("compile prompt", err, "Provide a non-empty prompt request.")
			}
			if app.Intelligence != nil {
				git := personal.DetectGitContext(lc.Ctx, ".")
				insights := app.Intelligence.Analyze(lc.Model)
				insights.References = app.Intelligence.FindSimilar(request, lc.Model)
				impact := app.Intelligence.AnalyzeImpact(git, lc.Model)
				compiled.Context += intelligence.FormatPromptEnhancement(insights, impact)
			}
			prefs := personalPrefs(app)
			promptPkg := buildPromptPackage(client, task, compiled, pkg, prefs)
			rendered, ext, err := renderPromptPackage(promptPkg, format)
			if err != nil {
				return appErr("render prompt package", err, "Use markdown, text, or json.")
			}
			if outPath == "" {
				outPath = defaultPromptExportPath(task, ext)
			}
			if err := lc.FS.WriteFile(outPath, rendered, 0644); err != nil {
				return appErr("export prompt package", err, "Check write permissions or choose another --out path.")
			}
			app.EventBus.Publish(eventbus.Event{Type: eventbus.PromptPackageCreated, Message: "prompt package created", Payload: map[string]interface{}{"file": outPath, "client": promptPkg.Client, "format": format}})
			if copyPrompt {
				if err := copyToClipboard(string(rendered)); err != nil {
					return appErr("copy prompt package", err, "Install a clipboard utility or omit --copy.")
				}
			}
			result := map[string]interface{}{
				"status":                 "exported",
				"file":                   outPath,
				"copied":                 copyPrompt,
				"format":                 format,
				"client":                 promptPkg.Client,
				"estimated_context_size": promptPkg.EstimatedContextSize,
			}
			return renderCLI(app, lc, result, func() string {
				lines := []string{
					"Prompt package exported",
					"File: " + outPath,
					fmt.Sprintf("Client: %s", promptPkg.Client),
					fmt.Sprintf("Estimated context size: %d bytes", promptPkg.EstimatedContextSize),
				}
				if copyPrompt {
					lines = append(lines, "Copied to clipboard")
				}
				return strings.Join(lines, "\n")
			})
		}),
	}
	cmd.Flags().StringVar(&task, "task", "feature", "task type")
	cmd.Flags().StringVar(&workflow, "workflow", "", "workflow type or prompt workflow")
	cmd.Flags().StringVar(&request, "request", "", "user request")
	cmd.Flags().StringVar(&provider, "provider", "", "AI provider")
	cmd.Flags().StringVar(&model, "model", "", "AI model")
	cmd.Flags().StringVar(&format, "format", "markdown", "prompt export format: markdown, text, json")
	cmd.Flags().StringVar(&client, "client", "generic", "external AI client template: claude, codex, chatgpt, generic")
	cmd.Flags().StringVar(&outPath, "out", "", "prompt export output path")
	cmd.Flags().BoolVar(&copyPrompt, "copy", false, "copy generated prompt package to clipboard")
	cmd.Flags().StringVar(&budget, "budget", string(contextengine.BudgetMedium), "context budget")
	addContextLimitFlags(cmd, &maxBytes)
	return cmd
}

func newAICommand(app *App) *cobra.Command {
	var provider, model, prompt string
	cmd := &cobra.Command{Use: "ai", Short: "Execute AI provider requests", RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
		if prompt == "" {
			prompt = strings.Join(args, " ")
		}
		resp, err := app.AI.Generate(lc.Ctx, ai.Request{Provider: provider, Model: model, Prompt: prompt})
		if err != nil {
			return appErr("execute AI request", err, "Configure provider credentials or choose an available local provider.")
		}
		return renderCLI(app, lc, resp, func() string { return resp.Content })
	})}
	cmd.Flags().StringVar(&provider, "provider", "ollama", "AI provider")
	cmd.Flags().StringVar(&model, "model", "", "AI model")
	cmd.Flags().StringVar(&prompt, "prompt", "", "prompt text")
	return cmd
}

func newPluginCommand(app *App) *cobra.Command {
	cmd := &cobra.Command{Use: "plugin", Short: "Manage PromptEngine plugins"}
	cmd.AddCommand(&cobra.Command{Use: "list", Short: "List plugins", RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
		if err := loadInstalledPlugins(lc, app); err != nil {
			return appErr("load installed plugins", err, "Verify installed plugin manifests under .promptengine/plugins.")
		}
		list := app.Plugins.List()
		return renderCLI(app, lc, list, func() string {
			if len(list) == 0 {
				return "No plugins registered."
			}
			var lines []string
			for _, p := range list {
				lines = append(lines, fmt.Sprintf("%s %s %s", p.ID, p.Version, p.Status))
			}
			return strings.Join(lines, "\n")
		})
	})})
	cmd.AddCommand(pluginLifecycleCommand(app, "install"), pluginLifecycleCommand(app, "remove"), pluginLifecycleCommand(app, "enable"), pluginLifecycleCommand(app, "disable"))
	cmd.AddCommand(&cobra.Command{Use: "health", Short: "Check plugin health", RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
		if err := loadInstalledPlugins(lc, app); err != nil {
			return appErr("load installed plugins", err, "Verify installed plugin manifests under .promptengine/plugins.")
		}
		findings := app.Plugins.Health(lc.Ctx, lc.FS)
		return renderCLI(app, lc, findings, func() string { return fmt.Sprintf("Plugin health findings: %d", len(findings)) })
	})})
	return cmd
}

func newPluginsAliasCommand(app *App) *cobra.Command {
	cmd := newPluginCommand(app)
	cmd.Use = "plugins"
	cmd.Short = "Manage PromptEngine plugins"
	return cmd
}

func newAnalyzeCommand(app *App) *cobra.Command {
	return &cobra.Command{Use: "analyze", Short: "Analyze current project changes", RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
		git := personal.DetectGitContext(lc.Ctx, ".")
		report := app.Intelligence.AnalyzeImpact(git, lc.Model)
		return renderCLI(app, lc, report, func() string {
			return "Change Analysis\nChanged files:\n" + strings.Join(prefixLines(report.ChangedFiles, "  "), "\n") + "\nAffected areas:\n" + strings.Join(prefixLines(report.AffectedAreas, "  "), "\n")
		})
	})}
}

func newSyncCommand(app *App) *cobra.Command {
	var changed []string
	var dryRun bool
	cmd := &cobra.Command{Use: "sync", Short: "Detect documentation synchronization needs", RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
		result := app.Docs.Sync(changed, dryRun)
		return renderCLI(app, lc, result, func() string {
			return fmt.Sprintf("Documentation sync\nPending: %d\nRecommendations: %d", len(result.Pending), len(result.Recommendations))
		})
	})}
	cmd.Flags().StringSliceVar(&changed, "changed", nil, "changed files")
	cmd.Flags().BoolVar(&dryRun, "dry-run", true, "report changes without writing")
	return cmd
}

func newHooksCommand(app *App) *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{Use: "hooks", Short: "Install PromptEngine git hook foundations", RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
		path := filepath.Join(".git", "hooks", "pre-commit")
		content := []byte("#!/bin/sh\npromptengine verify\n")
		if dryRun {
			return renderCLI(app, lc, statusResult{Status: "preview", Files: []string{path}}, func() string { return "Would install " + path })
		}
		if err := lc.FS.WriteFile(path, content, 0755); err != nil {
			return appErr("install git hook", err, "Run inside a git repository with writable .git/hooks.")
		}
		return renderCLI(app, lc, statusResult{Status: "installed", Files: []string{path}}, func() string { return "Installed " + path })
	})}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview hook changes")
	return cmd
}

func newInstallCommand(app *App) *cobra.Command {
	var manifestPath string
	cmd := &cobra.Command{Use: "install [id]", Short: "Install a local PromptEngine package", RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
		if manifestPath == "" {
			return appErr("install package", fmt.Errorf("package manifest path is required"), "Use plugin install --manifest or pass --manifest.")
		}
		meta, err := plugins.NewLoader(lc.FS).LoadManifest(manifestPath)
		if err != nil {
			return appErr("load plugin manifest", err, "Verify the plugin manifest path and schema.")
		}
		if err := app.Plugins.Register(plugins.NewManifestPlugin(meta)); err != nil {
			return appErr("register plugin", err, "Verify plugin compatibility, permissions, and dependencies.")
		}
		if err := app.Plugins.Install(lc.Ctx, meta.ID, lc.FS); err != nil {
			return appErr("install plugin", err, "Verify the plugin can install safely.")
		}
		record, err := app.Installer.Install(installer.PackageManifest{ID: meta.ID, Name: meta.Name, Kind: installer.KindPlugin, Version: meta.Version, Files: meta.Files})
		if err != nil {
			return appErr("install package", err, "Verify the package is not already installed.")
		}
		if err := persistPluginMetadata(lc.FS, meta); err != nil {
			return appErr("persist plugin metadata", err, "Check write permissions for .promptengine/plugins.")
		}
		return renderCLI(app, lc, record, func() string { return "Installed plugin " + meta.ID })
	})}
	cmd.Flags().StringVar(&manifestPath, "manifest", "", "local package or plugin manifest")
	return cmd
}

func newUpdateCommand(app *App) *cobra.Command {
	var target, id, from, to string
	var dryRun bool
	cmd := &cobra.Command{Use: "update", Short: "Plan or apply safe PromptEngine updates", RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
		if id == "" && len(args) > 0 {
			id = args[0]
		}
		if target == "" {
			target = string(updater.TargetCLI)
		}
		if id == "" {
			id = target
		}
		engine := updater.NewUpdateEngine()
		req := updater.UpdateRequest{Target: updater.UpdateTarget(target), ID: id, FromVer: from, ToVer: to, DryRun: dryRun}
		report := engine.Plan(req)
		if !dryRun {
			var err error
			report, err = engine.Apply(req, nil)
			if err != nil {
				return err
			}
		}
		return renderCLI(app, lc, report, func() string {
			return fmt.Sprintf("Update %s\nCan update: %t\nIssues: %d", id, report.CanUpdate, len(report.Issues))
		})
	})}
	cmd.Flags().StringVar(&target, "target", "cli", "update target: cli, plugin, manifest, template, schema")
	cmd.Flags().StringVar(&id, "id", "", "component id")
	cmd.Flags().StringVar(&from, "from", "", "current version")
	cmd.Flags().StringVar(&to, "to", "", "target version")
	cmd.Flags().BoolVar(&dryRun, "dry-run", true, "plan update without applying changes")
	return cmd
}

func newAgentsCommand(app *App) *cobra.Command {
	cmd := &cobra.Command{Use: "agents", Short: "Generate and synchronize AI coding agent instructions"}
	cmd.AddCommand(&cobra.Command{Use: "list", Short: "List supported agent profiles", RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
		profiles := app.Agents.Profiles()
		return renderCLI(app, lc, profiles, func() string {
			var lines []string
			for _, profile := range profiles {
				lines = append(lines, fmt.Sprintf("%s -> %s", profile.ID, profile.InstructionFile))
			}
			return strings.Join(lines, "\n")
		})
	})})
	var profile string
	syncCmd := &cobra.Command{Use: "sync", Short: "Generate or update agent instruction files", RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
		report, err := app.Agents.Sync(agents.InstructionRequest{Profile: profile, Project: lc.Model, Manifest: lc.Manifest, Preferences: personalPrefs(app)})
		if err != nil {
			return appErr("sync agent instructions", err, "Use a supported agent profile or all.")
		}
		return renderCLI(app, lc, report, func() string {
			return fmt.Sprintf("Agent instructions synchronized\nGenerated: %d\nUpdated: %d\nCurrent: %d", len(report.Generated), len(report.Updated), len(report.Current))
		})
	})}
	syncCmd.Flags().StringVar(&profile, "agent", "all", "agent profile to sync or all")
	cmd.AddCommand(syncCmd)
	return cmd
}

func newProfileCommand(app *App) *cobra.Command {
	cmd := &cobra.Command{Use: "profile", Short: "Manage the local personal developer profile"}
	cmd.AddCommand(&cobra.Command{Use: "show", Short: "Show personal profile", RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
		profile, err := app.Personal.LoadProfile("")
		if err != nil {
			return appErr("load personal profile", err, "Check .promptengine/profile.yaml.")
		}
		return renderCLI(app, lc, profile, func() string {
			return "Developer Preferences:\n" + strings.Join(prefixLines(profile.PreferenceLines(), "  "), "\n")
		})
	})})
	cmd.AddCommand(&cobra.Command{Use: "init", Short: "Create a default personal profile", RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
		profile := personal.DefaultProfile()
		if err := app.Personal.SaveProfile("", profile); err != nil {
			return appErr("save personal profile", err, "Check write permissions for .promptengine/profile.yaml.")
		}
		return renderCLI(app, lc, statusResult{Status: "created", Files: []string{personal.DefaultProfilePath}}, func() string { return "Created " + personal.DefaultProfilePath })
	})})
	return cmd
}

func newTaskCommand(app *App) *cobra.Command {
	var template, client, format, outPath, budget string
	var maxBytes int
	cmd := &cobra.Command{Use: "task [description]", Short: "Prepare a personal AI task package", RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
		description := strings.Join(args, " ")
		if description == "" {
			return appErr("prepare task", fmt.Errorf("task description is required"), "Pass a task description in quotes.")
		}
		profile, err := app.Personal.LoadProfile("")
		if err != nil {
			return appErr("load personal profile", err, "Run promptengine profile init.")
		}
		memory, err := app.Personal.LoadMemory("")
		if err != nil {
			return appErr("load personal memory", err, "Check .promptengine/memory.yaml.")
		}
		tmpl := personal.ResolveTemplate(template)
		git := personal.DetectGitContext(lc.Ctx, ".")
		pkg, err := app.Context.Build(lc.Ctx, contextengine.ContextRequest{TaskType: contextengine.TaskType(tmpl.TaskType), Project: lc.Model, UserIntent: description, AffectedFiles: git.ChangedFiles, Budget: contextengine.BudgetType(budget), MaxBytes: maxBytes})
		if err != nil {
			return appErr("build task context", err, "Reduce context limits or verify project discovery.")
		}
		taskPkg := personal.BuildTaskPackage(personal.TaskRequest{Description: description, Template: template, Profile: profile, Memory: memory, Git: git, Project: lc.Model, Context: pkg})
		if app.Intelligence != nil {
			insights := app.Intelligence.Analyze(lc.Model)
			insights.References = app.Intelligence.FindSimilar(description, lc.Model)
			impact := app.Intelligence.AnalyzeImpact(git, lc.Model)
			taskPkg.Prompt += intelligence.FormatPromptEnhancement(insights, impact)
		}
		compiled, err := app.AI.Compile(ai.CompileInput{Provider: "", UserRequest: description, ContextPackage: pkg, Manifest: lc.Manifest, SystemInstructions: "Use the personal developer profile and PromptEngine standards."})
		if err == nil {
			promptPkg := buildPromptPackage(client, tmpl.TaskType, compiled, pkg, profile.PreferenceLines())
			rendered, ext, err := renderPromptPackage(promptPkg, format)
			if err != nil {
				return appErr("render task prompt", err, "Use markdown, text, or json.")
			}
			if outPath == "" {
				outPath = defaultPromptExportPath(tmpl.TaskType, ext)
			}
			if err := lc.FS.WriteFile(outPath, rendered, 0644); err != nil {
				return appErr("write task prompt", err, "Check write permissions or pass --out.")
			}
		}
		return renderCLI(app, lc, taskPkg, func() string {
			return fmt.Sprintf("Task: %s\nWorkflow: %s\nSelected files:\n%s\nSelected docs:\n%s\nPrompt: %s", taskPkg.Summary, taskPkg.RecommendedWorkflow, strings.Join(prefixLines(taskPkg.SelectedFiles, "  "), "\n"), strings.Join(prefixLines(taskPkg.SelectedDocuments, "  "), "\n"), outPath)
		})
	})}
	cmd.Flags().StringVar(&template, "template", "", "task template: feature, bug_fix, refactor, review, architecture, docs, security")
	cmd.Flags().StringVar(&client, "client", "codex", "AI client prompt format")
	cmd.Flags().StringVar(&format, "format", "markdown", "prompt export format")
	cmd.Flags().StringVar(&outPath, "out", "", "prompt output path")
	cmd.Flags().StringVar(&budget, "budget", string(contextengine.BudgetMedium), "context budget")
	addContextLimitFlags(cmd, &maxBytes)
	return cmd
}

func addContextLimitFlags(cmd *cobra.Command, target *int) {
	cmd.Flags().IntVar(target, "max-bytes", 0, "maximum context size in bytes")
	cmd.Flags().IntVar(target, "limit", 0, "alias for --max-bytes")
}

func newInsightsCommand(app *App) *cobra.Command {
	return &cobra.Command{Use: "insights", Short: "Show local project patterns and recommendations", RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
		insights := app.Intelligence.Analyze(lc.Model)
		return renderCLI(app, lc, insights, func() string {
			return fmt.Sprintf("Insights\nPatterns: %d\nDecisions: %d\nRecommendations: %d", len(insights.Patterns), len(insights.Decisions), len(insights.Recommendations))
		})
	})}
}

func newDecisionsCommand(app *App) *cobra.Command {
	cmd := &cobra.Command{Use: "decisions", Short: "Manage local architecture decision memory"}
	cmd.AddCommand(&cobra.Command{Use: "list", Short: "List stored decisions", RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
		decisions, err := app.Intelligence.ListDecisions()
		if err != nil {
			return appErr("list decisions", err, "Check .promptengine/decisions.yaml.")
		}
		return renderCLI(app, lc, decisions, func() string {
			if len(decisions) == 0 {
				return "No decisions stored."
			}
			var lines []string
			for _, d := range decisions {
				lines = append(lines, d.Title+": "+d.Reason)
			}
			return strings.Join(lines, "\n")
		})
	})})
	var title, reason string
	var affected []string
	store := &cobra.Command{Use: "store", Short: "Store an architecture decision", RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
		if title == "" && len(args) > 0 {
			title = strings.Join(args, " ")
		}
		if err := app.Intelligence.StoreDecision(intelligence.Decision{Title: title, Reason: reason, Affected: affected}); err != nil {
			return appErr("store decision", err, "Provide --title and --reason.")
		}
		return renderCLI(app, lc, statusResult{Status: "stored", Message: title}, func() string { return "Stored decision: " + title })
	})}
	store.Flags().StringVar(&title, "title", "", "decision title")
	store.Flags().StringVar(&reason, "reason", "", "decision reason")
	store.Flags().StringSliceVar(&affected, "affected", nil, "affected areas")
	cmd.AddCommand(store)
	return cmd
}

func newImpactCommand(app *App) *cobra.Command {
	return &cobra.Command{Use: "impact", Short: "Analyze potential impact of current changes", RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
		git := personal.DetectGitContext(lc.Ctx, ".")
		report := app.Intelligence.AnalyzeImpact(git, lc.Model)
		return renderCLI(app, lc, report, func() string {
			return "Impact Analysis\nChanged files:\n" + strings.Join(prefixLines(report.ChangedFiles, "  "), "\n") + "\nAffected areas:\n" + strings.Join(prefixLines(report.AffectedAreas, "  "), "\n")
		})
	})}
}

func newVerifyCommand(app *App) *cobra.Command {
	return &cobra.Command{Use: "verify", Short: "Run personal pre-commit verification", RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
		audit, err := app.Quality.Audit(lc.Ctx)
		if err != nil {
			return appErr("run verification", err, "Review quality platform errors.")
		}
		docsReport, err := app.Docs.Validate(lc.Model)
		if err != nil {
			return appErr("verify documentation", err, "Fix documentation validation errors.")
		}
		result := map[string]interface{}{"quality": audit, "documentation_findings": len(docsReport.Findings), "passed": audit.Score.Passed && len(docsReport.Findings) == 0}
		return renderCLI(app, lc, result, func() string {
			return fmt.Sprintf("Verification\nHealth Score: %d/100\nQuality findings: %d\nDocumentation findings: %d", audit.Score.Overall, len(audit.Findings), len(docsReport.Findings))
		})
	})}
}

func newMemoryCommand(app *App) *cobra.Command {
	cmd := &cobra.Command{Use: "memory", Short: "Manage local non-sensitive personal memory"}
	cmd.AddCommand(&cobra.Command{Use: "show", Short: "Show memory", RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
		mem, err := app.Personal.LoadMemory("")
		if err != nil {
			return appErr("load memory", err, "Check .promptengine/memory.yaml.")
		}
		return renderCLI(app, lc, mem, nil)
	})})
	var key, value string
	add := &cobra.Command{Use: "add", Short: "Add a local memory note", RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
		if key == "" || value == "" {
			return appErr("add memory", fmt.Errorf("key and value are required"), "Pass --key and --value.")
		}
		mem, err := app.Personal.AddMemory(key, value)
		if err != nil {
			return appErr("add memory", err, "Do not store secrets, API keys, credentials, or tokens.")
		}
		return renderCLI(app, lc, mem, func() string { return "Stored memory: " + key })
	})}
	add.Flags().StringVar(&key, "key", "", "memory key")
	add.Flags().StringVar(&value, "value", "", "memory value")
	cmd.AddCommand(add)
	return cmd
}

func pluginLifecycleCommand(app *App, action string) *cobra.Command {
	var manifestPath string
	cmd := &cobra.Command{Use: action + " [id]", Short: action + " plugin", RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
		id := ""
		if len(args) > 0 {
			id = args[0]
		}
		switch action {
		case "install":
			if manifestPath == "" {
				return appErr("install plugin", fmt.Errorf("plugin manifest path is required"), "Pass --manifest with a local plugin manifest.")
			}
			meta, err := plugins.NewLoader(lc.FS).LoadManifest(manifestPath)
			if err != nil {
				return appErr("load plugin manifest", err, "Verify the plugin manifest path and schema.")
			}
			plugin := plugins.NewManifestPlugin(meta)
			if err := app.Plugins.Register(plugin); err != nil {
				return appErr("register plugin", err, "Verify plugin compatibility, permissions, and dependencies.")
			}
			if err := app.Plugins.Install(lc.Ctx, meta.ID, lc.FS); err != nil {
				return appErr("install plugin", err, "Verify the plugin can install safely.")
			}
			record, err := app.Installer.Install(installer.PackageManifest{ID: meta.ID, Name: meta.Name, Kind: installer.KindPlugin, Version: meta.Version, Files: meta.Files})
			if err != nil {
				return appErr("install plugin", err, "Verify the plugin is not already installed.")
			}
			if err := persistPluginMetadata(lc.FS, meta); err != nil {
				return appErr("persist plugin metadata", err, "Check write permissions for .promptengine/plugins.")
			}
			return renderCLI(app, lc, record, func() string { return "Installed plugin " + meta.ID })
		case "remove":
			if err := loadInstalledPlugins(lc, app); err != nil {
				return appErr("load installed plugins", err, "Verify installed plugin manifests under .promptengine/plugins.")
			}
			if err := app.Plugins.Remove(lc.Ctx, id, lc.FS); err != nil {
				return appErr("remove plugin", err, "Verify the plugin is installed.")
			}
			if err := app.Installer.Uninstall(id); err != nil {
				return appErr("remove plugin", err, "Verify the plugin is installed.")
			}
		case "enable":
			if err := loadInstalledPlugins(lc, app); err != nil {
				return appErr("load installed plugins", err, "Verify installed plugin manifests under .promptengine/plugins.")
			}
			if err := app.Plugins.Enable(id); err != nil {
				return appErr("enable plugin", err, "Install the plugin before enabling it.")
			}
			if err := app.Installer.Enable(id); err != nil {
				return appErr("enable plugin", err, "Install the plugin before enabling it.")
			}
		case "disable":
			if err := loadInstalledPlugins(lc, app); err != nil {
				return appErr("load installed plugins", err, "Verify installed plugin manifests under .promptengine/plugins.")
			}
			if err := app.Plugins.Disable(id); err != nil {
				return appErr("disable plugin", err, "Install the plugin before disabling it.")
			}
			if err := app.Installer.Disable(id); err != nil {
				return appErr("disable plugin", err, "Install the plugin before disabling it.")
			}
		}
		return renderCLI(app, lc, statusResult{Status: action, Message: id}, func() string { return strings.Title(action) + " plugin " + id })
	})}
	cmd.Flags().StringVar(&manifestPath, "manifest", "", "plugin manifest path")
	return cmd
}

func persistPluginMetadata(fs filesystem.FileSystem, meta plugins.PluginMetadata) error {
	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	return fs.WriteFile(filepath.Join(".promptengine", "plugins", meta.ID, "plugin.json"), data, 0644)
}

func loadInstalledPlugins(lc *LifecycleContext, app *App) error {
	metas, err := plugins.NewLoader(lc.FS).LoadManifestsFrom(filepath.Join(".promptengine", "plugins"))
	if err != nil {
		return err
	}
	for _, meta := range metas {
		markerPath := filepath.Join(".promptengine", "plugins", meta.ID, ".installed")
		marker, markerErr := lc.FS.ReadFile(markerPath)
		if markerErr != nil || strings.TrimSpace(string(marker)) == "uninstalled" {
			continue
		}
		if _, exists := app.Plugins.Get(meta.ID); exists {
			continue
		}
		if err := app.Plugins.Register(plugins.NewManifestPlugin(meta)); err != nil {
			return err
		}
	}
	return nil
}

func newConfigCommand(app *App) *cobra.Command {
	cmd := &cobra.Command{Use: "config", Short: "View and update configuration"}
	cmd.AddCommand(&cobra.Command{Use: "view", Short: "View active configuration", RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
		return renderCLI(app, lc, lc.Config, nil)
	})})
	var key, value string
	set := &cobra.Command{Use: "set", Short: "Set a project configuration value", RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
		if key == "" || value == "" {
			return appErr("set configuration", fmt.Errorf("key and value are required"), "Pass --key and --value.")
		}
		if err := setConfigValue(lc.Config, key, value); err != nil {
			return appErr("set configuration", err, "Use a supported project, docs, or mode key.")
		}
		data, err := yaml.Marshal(lc.Config)
		if err != nil {
			return appErr("serialize configuration", err, "Review active configuration values.")
		}
		path := lc.Config.CLI.Config
		if path == "" {
			path = ".promptengine.yaml"
		}
		if err := lc.FS.WriteFile(path, data, 0644); err != nil {
			return appErr("write configuration", err, "Check write permissions for the project configuration file.")
		}
		return renderCLI(app, lc, statusResult{Status: "updated", Message: key}, func() string { return "Updated " + key })
	})}
	set.Flags().StringVar(&key, "key", "", "configuration key")
	set.Flags().StringVar(&value, "value", "", "configuration value")
	cmd.AddCommand(set)
	return cmd
}

func newCompletionCommand(app *App) *cobra.Command {
	cmd := &cobra.Command{Use: "completion [bash|zsh|fish|powershell]", Short: "Generate shell completion", Args: cobra.ExactArgs(1), RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
		switch args[0] {
		case "bash":
			return lc.Cmd.Root().GenBashCompletion(app.Out)
		case "zsh":
			return lc.Cmd.Root().GenZshCompletion(app.Out)
		case "fish":
			return lc.Cmd.Root().GenFishCompletion(app.Out, true)
		case "powershell":
			return lc.Cmd.Root().GenPowerShellCompletion(app.Out)
		default:
			return fmt.Errorf("unsupported shell %q", args[0])
		}
	})}
	return cmd
}

func qualityShort(name string) string {
	switch name {
	case "doctor":
		return "Run project diagnostics"
	case "review":
		return "Run engineering review"
	case "audit":
		return "Run a full quality audit"
	case "health":
		return "Calculate project health score"
	default:
		return name
	}
}

func appErr(action string, err error, rec string) error {
	return apperrors.New(apperrors.CategoryCommand, apperrors.ExitGeneralError, action, err).WithRecommendation(rec)
}

func join(items []string) string {
	if len(items) == 0 {
		return "none"
	}
	sort.Strings(items)
	return strings.Join(items, ", ")
}

func prefixLines(items []string, prefix string) []string {
	if len(items) == 0 {
		return []string{prefix + "none"}
	}
	out := make([]string, len(items))
	for i, item := range items {
		out[i] = prefix + item
	}
	return out
}

func expandAgentProfiles(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	var out []string
	for _, value := range values {
		for _, item := range strings.Split(value, ",") {
			item = strings.TrimSpace(item)
			if item == "" {
				continue
			}
			if item == "all" {
				out = append(out, "codex", "claude", "codex-md", "cursor", "windsurf")
				continue
			}
			out = append(out, item)
		}
	}
	return dedupeStrings(out)
}

func personalPrefs(app *App) []string {
	if app == nil || app.Personal == nil {
		return nil
	}
	profile, err := app.Personal.LoadProfile("")
	if err != nil {
		return nil
	}
	return profile.PreferenceLines()
}

func setConfigValue(cfg *config.AppConfig, key, value string) error {
	switch key {
	case "mode":
		cfg.Mode = value
	case "project.name":
		cfg.Project.Name = value
	case "project.version":
		cfg.Project.Version = value
	case "project.description":
		cfg.Project.Description = value
	case "project.promptengine_path":
		cfg.Project.PromptEnginePath = value
	case "docs.root_dir":
		cfg.Docs.RootDir = value
	case "docs.agents_constitution":
		cfg.Docs.AgentsConstitution = value
	case "docs.api_spec":
		cfg.Docs.APISpec = value
	case "docs.database_spec":
		cfg.Docs.DatabaseSpec = value
	case "docs.business_rules":
		cfg.Docs.BusinessRules = value
	default:
		return fmt.Errorf("unsupported configuration key %q", key)
	}
	return nil
}

func normalizeID(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.ReplaceAll(value, "_", "-")
	value = strings.ReplaceAll(value, " ", "-")
	return value
}
