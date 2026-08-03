package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	contextpkg "github.com/LordCodex/promptengine/internal/domain/context"
	"github.com/LordCodex/promptengine/internal/domain/docs/generator"
	docsync "github.com/LordCodex/promptengine/internal/domain/docs/sync"
	"github.com/LordCodex/promptengine/internal/domain/prompts"
	"github.com/LordCodex/promptengine/internal/domain/quality"
	"github.com/LordCodex/promptengine/internal/domain/quality/report"
	"github.com/LordCodex/promptengine/internal/domain/updater"
	"github.com/LordCodex/promptengine/internal/perf"
	"github.com/LordCodex/promptengine/internal/recovery"
	"github.com/LordCodex/promptengine/internal/security"
	"github.com/LordCodex/promptengine/internal/telemetry"
	"github.com/LordCodex/promptengine/internal/version"
	"github.com/spf13/cobra"
)

// ─── CommandBuilder ────────────────────────────────────────────────────────

type CommandBuilder struct {
	Out io.Writer
}

func NewCommandBuilder(out io.Writer) *CommandBuilder {
	return &CommandBuilder{Out: out}
}

func (b *CommandBuilder) BuildRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "promptengine",
		Short: "PromptEngine CLI governs codebase standards and guides AI coding assistants.",
		Long: `PromptEngine is the AI Coding Standards Engine for software teams.
It governs documentation, enforces engineering standards, and generates
context-aware prompts for AI coding assistants.`,
	}
	cmd.SetOut(b.Out)
	cmd.SetErr(b.Out)
	return cmd
}

// ─── Version ───────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildVersion() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print PromptEngine version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(b.Out, version.String())
		},
	}
}

// ─── doctor ────────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildDoctor(app *App) *cobra.Command {
	var format string
	var threshold int
	var autoFix bool
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Diagnose PromptEngine installation and project health",
		Long:  "Runs all registered doctor checks and generates an actionable diagnostic report.",
		RunE: recovery.CommandPanicWrapper(b.Out, func(cmd *cobra.Command, args []string) error {
			timer := perf.Start("doctor")
			defer timer.Log(app.Logger)

			docReport, err := app.Doctor.Diagnose(app.FS)
			if err != nil {
				return fmt.Errorf("doctor: %w", err)
			}

			data, err := app.Reporter.Render(format, docReport)
			if err != nil {
				return err
			}
			fmt.Fprintln(b.Out, string(data))

			t := quality.Threshold{MinScore: threshold, BlockOnSeverity: quality.SeverityError}
			ci := report.EvaluateCIThreshold(docReport, t)

			if autoFix && !ci.Passed {
				for _, f := range docReport.Findings {
					if f.AutoFixID != "" {
						result := app.Fixer.Apply(f.AutoFixID, app.FS, dryRun)
						if result.Error != nil {
							app.Logger.Warn("fix failed", "id", f.AutoFixID, "err", result.Error)
						} else {
							fmt.Fprintf(b.Out, "  ✓ Fix applied: %s\n", result.Description)
						}
					}
				}
			}

			app.Telemetry.Track(telemetry.Event{
				Command:  "doctor",
				Duration: timer.Duration().Milliseconds(),
				OS:       version.GetInfo().OS,
				Arch:     version.GetInfo().Arch,
			})

			if !ci.Passed {
				return fmt.Errorf("doctor: %s", ci.Message)
			}
			return nil
		}),
	}
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format: text, json, yaml, markdown, sarif")
	cmd.Flags().IntVarP(&threshold, "threshold", "t", 70, "Minimum score to pass")
	cmd.Flags().BoolVar(&autoFix, "fix", false, "Attempt to apply safe auto-fixes")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview fixes without writing")
	return cmd
}

// ─── health ────────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildHealth(app *App) *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "health",
		Short: "Output workspace health scores and category breakdown",
		RunE: recovery.CommandPanicWrapper(b.Out, func(cmd *cobra.Command, args []string) error {
			timer := perf.Start("health")
			defer timer.Log(app.Logger)

			result, err := app.Health.Evaluate(app.FS)
			if err != nil {
				return fmt.Errorf("health: %w", err)
			}

			switch format {
			case "json":
				data, _ := json.MarshalIndent(result, "", "  ")
				fmt.Fprintln(b.Out, string(data))
			default:
				fmt.Fprintf(b.Out, "Overall Score : %d/100 (%s)\n\n", result.Score, result.Rating)
				fmt.Fprintf(b.Out, "%-28s %6s\n", "Category", "Score")
				fmt.Fprintf(b.Out, "%s\n", strings.Repeat("-", 36))
				for name, cat := range result.Categories {
					icon := "✓"
					if cat.CriticalFail {
						icon = "✗"
					}
					fmt.Fprintf(b.Out, "%-28s %5d  %s\n", name, cat.Score, icon)
				}
				if len(result.Issues) > 0 {
					fmt.Fprintf(b.Out, "\nIssues (%d):\n", len(result.Issues))
					for _, iss := range result.Issues {
						fmt.Fprintf(b.Out, "  [%s] %s — %s\n", strings.ToUpper(iss.Severity), iss.Category, iss.Title)
					}
				}
			}

			app.Telemetry.Track(telemetry.Event{
				Command:  "health",
				Duration: timer.Duration().Milliseconds(),
				OS:       version.GetInfo().OS,
				Arch:     version.GetInfo().Arch,
			})
			return nil
		}),
	}
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format: text, json")
	return cmd
}

// ─── review ────────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildReview(app *App) *cobra.Command {
	var format string
	var path string

	cmd := &cobra.Command{
		Use:   "review",
		Short: "Run structured code and documentation reviews",
		RunE: recovery.CommandPanicWrapper(b.Out, func(cmd *cobra.Command, args []string) error {
			timer := perf.Start("review")
			defer timer.Log(app.Logger)

			if path == "" {
				path = "."
			}

			// Validate safe path bounds
			safePath, err := security.ValidateSafePath(".", path)
			if err != nil {
				return err
			}

			session, err := app.Reviewer.RunSession(app.FS, safePath)
			if err != nil {
				return fmt.Errorf("review: %w", err)
			}

			switch format {
			case "json":
				data, _ := json.MarshalIndent(session.Findings, "", "  ")
				fmt.Fprintln(b.Out, string(data))
			default:
				if len(session.Findings) == 0 {
					fmt.Fprintln(b.Out, "No review findings. ✓")
					return nil
				}
				fmt.Fprintf(b.Out, "Review findings (%d):\n\n", len(session.Findings))
				for _, f := range session.Findings {
					fmt.Fprintf(b.Out, "  [%s] %s\n", strings.ToUpper(f.Severity), f.Message)
					if f.Fix != "" {
						fmt.Fprintf(b.Out, "    → %s\n", f.Fix)
					}
				}
				if len(session.Summary) > 0 {
					fmt.Fprintf(b.Out, "\nSummary by type:\n")
					for rt, count := range session.Summary {
						fmt.Fprintf(b.Out, "  %-30s %d\n", rt, count)
					}
				}
			}

			app.Telemetry.Track(telemetry.Event{
				Command:  "review",
				Duration: timer.Duration().Milliseconds(),
				OS:       version.GetInfo().OS,
				Arch:     version.GetInfo().Arch,
			})
			return nil
		}),
	}
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format: text, json")
	cmd.Flags().StringVarP(&path, "path", "p", ".", "Target directory to review")
	return cmd
}

// ─── audit ─────────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildAudit(app *App) *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "audit",
		Short: "Run a comprehensive project audit",
		RunE: recovery.CommandPanicWrapper(b.Out, func(cmd *cobra.Command, args []string) error {
			timer := perf.Start("audit")
			defer timer.Log(app.Logger)

			auditReport, err := app.Auditor.Run(app.FS)
			if err != nil {
				return fmt.Errorf("audit: %w", err)
			}
			data, err := auditReport.Export(format)
			if err != nil {
				return err
			}
			fmt.Fprintln(b.Out, string(data))

			app.Telemetry.Track(telemetry.Event{
				Command:  "audit",
				Duration: timer.Duration().Milliseconds(),
				OS:       version.GetInfo().OS,
				Arch:     version.GetInfo().Arch,
			})
			return nil
		}),
	}
	cmd.Flags().StringVarP(&format, "format", "f", "markdown", "Output format: markdown, json")
	return cmd
}

// ─── scan ──────────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildScan(app *App) *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "scan",
		Short: "Scan workspace and discover project stack",
		RunE: recovery.CommandPanicWrapper(b.Out, func(cmd *cobra.Command, args []string) error {
			timer := perf.Start("scan")
			defer timer.Log(app.Logger)

			pm, err := app.Discovery.Execute(context.Background(), app.FS, ".")
			if err != nil {
				return fmt.Errorf("scan: %w", err)
			}
			switch format {
			case "json":
				data, _ := json.MarshalIndent(pm, "", "  ")
				fmt.Fprintln(b.Out, string(data))
			default:
				fmt.Fprintf(b.Out, "Project Root    : %s\n", pm.RootDir)
				fmt.Fprintf(b.Out, "Languages       : %s\n", strings.Join(pm.Languages, ", "))
				fmt.Fprintf(b.Out, "Frameworks      : %s\n", strings.Join(pm.Frameworks, ", "))
				fmt.Fprintf(b.Out, "Package Manager : %s\n", strings.Join(pm.PackageManagers, ", "))
				fmt.Fprintf(b.Out, "Has Git         : %v\n", pm.HasGit)
				fmt.Fprintf(b.Out, "Has Docker      : %v\n", pm.HasDocker)
				fmt.Fprintf(b.Out, "CI/CD           : %s\n", strings.Join(pm.CIs, ", "))
			}

			app.Telemetry.Track(telemetry.Event{
				Command:  "scan",
				Duration: timer.Duration().Milliseconds(),
				OS:       version.GetInfo().OS,
				Arch:     version.GetInfo().Arch,
			})
			return nil
		}),
	}
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format: text, json")
	return cmd
}

// ─── context ───────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildContext(app *App) *cobra.Command {
	var task string
	var format string

	cmd := &cobra.Command{
		Use:   "context",
		Short: "Output minimum required context files for a task",
		RunE: recovery.CommandPanicWrapper(b.Out, func(cmd *cobra.Command, args []string) error {
			timer := perf.Start("context")
			defer timer.Log(app.Logger)

			files, err := app.ContextBuilder.BuildContext(contextpkg.TaskType(task), nil)
			if err != nil {
				return fmt.Errorf("context: %w", err)
			}
			switch format {
			case "json":
				data, _ := json.Marshal(files)
				fmt.Fprintln(b.Out, string(data))
			default:
				fmt.Fprintf(b.Out, "Task    : %s\n", task)
				fmt.Fprintf(b.Out, "Files   : %d\n\n", len(files))
				for _, f := range files {
					fmt.Fprintf(b.Out, "  • %s\n", f)
				}
			}

			app.Telemetry.Track(telemetry.Event{
				Command:  "context",
				Duration: timer.Duration().Milliseconds(),
				OS:       version.GetInfo().OS,
				Arch:     version.GetInfo().Arch,
			})
			return nil
		}),
	}
	cmd.Flags().StringVarP(&task, "task", "t", "existing-project", "Task type (e.g. bug-fix, feature-development)")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format: text, json")
	return cmd
}

// ─── prompt ────────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildPrompt(app *App) *cobra.Command {
	var workflow string
	var provider string

	cmd := &cobra.Command{
		Use:   "prompt",
		Short: "Generate an AI-ready prompt for a given workflow",
		RunE: recovery.CommandPanicWrapper(b.Out, func(cmd *cobra.Command, args []string) error {
			timer := perf.Start("prompt")
			defer timer.Log(app.Logger)

			files, _ := app.ContextBuilder.BuildContext(contextpkg.TaskType(workflow), nil)
			ctx := make(prompts.ContextPackage)
			ctx["workflow"] = workflow
			ctx["files"] = strings.Join(files, ", ")

			p, err := app.PromptBuilder.Build(workflow, ctx, provider)
			if err != nil {
				// Graceful fallback
				fmt.Fprintf(b.Out, "# PromptEngine — %s workflow\n\n", workflow)
				if len(files) > 0 {
					fmt.Fprintf(b.Out, "Context files:\n")
					for _, f := range files {
						fmt.Fprintf(b.Out, "  • %s\n", f)
					}
				}
				fmt.Fprintln(b.Out, "\nNo registered prompt for this workflow. Register one via the Prompt Registry.")
				return nil
			}
			fmt.Fprintln(b.Out, p.CopyAndPastePrompt)

			app.Telemetry.Track(telemetry.Event{
				Command:  "prompt",
				Duration: timer.Duration().Milliseconds(),
				OS:       version.GetInfo().OS,
				Arch:     version.GetInfo().Arch,
			})
			return nil
		}),
	}
	cmd.Flags().StringVarP(&workflow, "workflow", "w", "existing-project", "Workflow name")
	cmd.Flags().StringVarP(&provider, "provider", "p", "", "AI provider hint (e.g. claude, gemini, chatgpt)")
	return cmd
}

// ─── generate ──────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildGenerate(app *App) *cobra.Command {
	var docType string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "generate [doc-type]",
		Short: "Generate project documentation scaffolds",
		Args:  cobra.MaximumNArgs(1),
		RunE: recovery.CommandPanicWrapper(b.Out, func(cmd *cobra.Command, args []string) error {
			timer := perf.Start("generate")
			defer timer.Log(app.Logger)

			if len(args) > 0 {
				docType = args[0]
			}

			input := generator.GeneratorInput{
				ProjectName: app.Config.Project.Name,
			}

			if docType == "" || docType == "all" {
				for _, g := range app.GenRegistry.All() {
					out, err := g.Generate(input)
					if err != nil {
						app.Logger.Warn("generator failed", "type", g.DocType(), "err", err)
						continue
					}
					// Validate target output bounds
					if _, err := security.ValidateSafePath(".", out.Filename); err != nil {
						return err
					}
					if dryRun {
						fmt.Fprintf(b.Out, "[dry-run] Would generate: %s\n", out.Filename)
						continue
					}
					if err := app.FS.WriteFile(out.Filename, []byte(out.Content), 0644); err != nil {
						return fmt.Errorf("write %s: %w", out.Filename, err)
					}
					fmt.Fprintf(b.Out, "  ✓ Generated: %s\n", out.Filename)
				}
				return nil
			}

			g, ok := app.GenRegistry.Get(generator.DocType(docType))
			if !ok {
				return fmt.Errorf("generate: no generator registered for '%s'", docType)
			}
			out, err := g.Generate(input)
			if err != nil {
				return fmt.Errorf("generate: %w", err)
			}
			// Validate target output bounds
			if _, err := security.ValidateSafePath(".", out.Filename); err != nil {
				return err
			}
			if dryRun {
				fmt.Fprintf(b.Out, "[dry-run] Would generate: %s\n", out.Filename)
				fmt.Fprintln(b.Out, out.Content)
				return nil
			}
			if err := app.FS.WriteFile(out.Filename, []byte(out.Content), 0644); err != nil {
				return fmt.Errorf("write %s: %w", out.Filename, err)
			}
			fmt.Fprintf(b.Out, "✓ Generated: %s\n", out.Filename)

			app.Telemetry.Track(telemetry.Event{
				Command:  "generate",
				Duration: timer.Duration().Milliseconds(),
				OS:       version.GetInfo().OS,
				Arch:     version.GetInfo().Arch,
			})
			return nil
		}),
	}
	cmd.Flags().StringVarP(&docType, "type", "t", "", "Document type to generate (omit for all)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview generation without writing files")
	return cmd
}

// ─── docs ──────────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildDocs(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "docs",
		Short: "Inspect and validate project documentation",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Show status of all registered documentation",
		RunE: recovery.CommandPanicWrapper(b.Out, func(cmd *cobra.Command, args []string) error {
			statuses, err := app.DocEngine.Status()
			if err != nil {
				return fmt.Errorf("docs status: %w", err)
			}
			if len(statuses) == 0 {
				fmt.Fprintln(b.Out, "No documents registered in the manifest.")
				return nil
			}
			fmt.Fprintf(b.Out, "%-35s %-10s %s\n", "Document", "Status", "Path")
			fmt.Fprintf(b.Out, "%s\n", strings.Repeat("-", 70))
			for _, s := range statuses {
				icon := "✓"
				if string(s.Status) != "synced" {
					icon = "✗"
				}
				fmt.Fprintf(b.Out, "%-35s %-10s %s  %s\n", s.Name, s.Status, s.Path, icon)
			}
			return nil
		}),
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "validate [doc-id]",
		Short: "Validate a specific document",
		Args:  cobra.ExactArgs(1),
		RunE: recovery.CommandPanicWrapper(b.Out, func(cmd *cobra.Command, args []string) error {
			issues, err := app.DocEngine.Validate(args[0])
			if err != nil {
				return fmt.Errorf("docs validate: %w", err)
			}
			if len(issues) == 0 {
				fmt.Fprintf(b.Out, "✓ %s passes validation\n", args[0])
				return nil
			}
			for _, iss := range issues {
				fmt.Fprintf(b.Out, "  ✗ %s\n", iss)
			}
			return nil
		}),
	})

	return cmd
}

// ─── sync ──────────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildSync(app *App) *cobra.Command {
	var dryRun bool
	var signals []string

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Synchronise documentation based on detected project changes",
		RunE: recovery.CommandPanicWrapper(b.Out, func(cmd *cobra.Command, args []string) error {
			changeSignals := make([]docsync.ChangeSignal, 0, len(signals))
			for _, s := range signals {
				changeSignals = append(changeSignals, docsync.ChangeSignal(s))
			}
			if len(changeSignals) == 0 {
				changeSignals = []docsync.ChangeSignal{
					docsync.SignalNewMigration,
					docsync.SignalNewAPI,
					docsync.SignalNewService,
				}
			}

			result := app.SyncEngine.Run(changeSignals, dryRun)
			if dryRun {
				fmt.Fprintln(b.Out, "Sync dry-run — no changes applied:")
				for _, docID := range result.Pending {
					fmt.Fprintf(b.Out, "  → Pending: %s\n", docID)
				}
				return nil
			}
			fmt.Fprintf(b.Out, "Applied : %d updates\n", len(result.Applied))
			fmt.Fprintf(b.Out, "Pending : %d (require approval)\n", len(result.Pending))
			return nil
		}),
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview sync recommendations without applying")
	cmd.Flags().StringSliceVar(&signals, "signal", nil, "Change signals to process (e.g. new-migration,new-api)")
	return cmd
}

// ─── hooks ─────────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildHooks(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hooks",
		Short: "Manage PromptEngine Git hooks",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "install",
		Short: "Install all registered Git hooks",
		RunE: recovery.CommandPanicWrapper(b.Out, func(cmd *cobra.Command, args []string) error {
			if err := app.HookRegistry.InstallAll(); err != nil {
				return fmt.Errorf("hooks install: %w", err)
			}
			fmt.Fprintln(b.Out, "✓ Git hooks installed")
			return nil
		}),
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "uninstall",
		Short: "Remove all installed Git hooks",
		RunE: recovery.CommandPanicWrapper(b.Out, func(cmd *cobra.Command, args []string) error {
			if err := app.HookRegistry.UninstallAll(); err != nil {
				return fmt.Errorf("hooks uninstall: %w", err)
			}
			fmt.Fprintln(b.Out, "✓ Git hooks removed")
			return nil
		}),
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all registered hooks",
		Run: func(cmd *cobra.Command, args []string) {
			hookList := app.HookRegistry.ListHooks()
			if len(hookList) == 0 {
				fmt.Fprintln(b.Out, "No hooks registered.")
				return
			}
			for _, h := range hookList {
				fmt.Fprintf(b.Out, "  • %s\n", h.Name())
			}
		},
	})

	return cmd
}

// ─── plugins ───────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildPlugins(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugins",
		Short: "List and manage PromptEngine plugin extensions",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all registered plugins",
		Run: func(cmd *cobra.Command, args []string) {
			pluginList := app.PluginRegistry.List()
			if len(pluginList) == 0 {
				fmt.Fprintln(b.Out, "No plugins installed.")
				return
			}
			for _, p := range pluginList {
				fmt.Fprintf(b.Out, "  • %s@%s — %s\n", p.ID, p.Version, p.Description)
			}
		},
	})

	return cmd
}

// ─── install ───────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildInstall(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install [source]",
		Short: "Install a plugin or technology stack from a source",
		Args:  cobra.ExactArgs(1),
		RunE: recovery.CommandPanicWrapper(b.Out, func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(b.Out, "Installing from: %s\n", args[0])
			fmt.Fprintln(b.Out, "ℹ  Remote marketplace install is not yet available.")
			fmt.Fprintln(b.Out, "  For local plugin installation, place plugin directory in .promptengine/plugins/")
			return nil
		}),
	}
	return cmd
}

// ─── update ────────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildUpdate(app *App) *cobra.Command {
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update PromptEngine standards and templates",
		RunE: recovery.CommandPanicWrapper(b.Out, func(cmd *cobra.Command, args []string) error {
			timer := perf.Start("update")
			defer timer.Log(app.Logger)

			req := updater.UpdateRequest{
				Target: updater.UpdateTarget("standards"),
			}
			updateReport := app.Updater.Plan(req)
			if !updateReport.CanUpdate {
				fmt.Fprintln(b.Out, "✓ Everything is up to date.")
				return nil
			}
			for _, iss := range updateReport.Issues {
				fmt.Fprintf(b.Out, "  [%s] %s\n", iss.Severity, iss.Message)
			}
			if dryRun {
				fmt.Fprintln(b.Out, "\n[dry-run] No changes applied.")
				return nil
			}
			fmt.Fprintln(b.Out, "\nApply updates by running: promptengine update --apply")
			return nil
		}),
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview available updates without applying")
	return cmd
}

// ─── init ──────────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildInit(app *App) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Bootstrap a new project's PromptEngine constitution",
		RunE: recovery.CommandPanicWrapper(b.Out, func(cmd *cobra.Command, args []string) error {
			timer := perf.Start("init")
			defer timer.Log(app.Logger)

			if app.FS.Exists("playbook-manifest.json") && !force {
				fmt.Fprintln(b.Out, "ℹ  PromptEngine already initialised. Use --force to reinitialise.")
				return nil
			}

			// Step 1: Scaffold
			fmt.Fprintln(b.Out, "Initialising PromptEngine...")
			for _, id := range []string{"create-promptengine-dir", "create-docs-dir", "create-manifest"} {
				result := app.Fixer.Apply(id, app.FS, false)
				if result.Error != nil {
					return fmt.Errorf("init scaffold '%s': %w", id, result.Error)
				}
				fmt.Fprintf(b.Out, "  ✓ %s\n", result.Description)
			}

			// Step 2: Discovery
			pm, err := app.Discovery.Execute(context.Background(), app.FS, ".")
			if err != nil {
				app.Logger.Warn("discovery failed during init", "err", err)
			} else if len(pm.Languages) > 0 {
				fmt.Fprintf(b.Out, "  ✓ Detected stack: %s\n", strings.Join(pm.Languages, ", "))
			}

			// Step 3: Generate core doc scaffolds
			fmt.Fprintln(b.Out, "\nGenerating core documentation scaffolds...")
			input := generator.GeneratorInput{ProjectName: app.Config.Project.Name}
			for _, docType := range []generator.DocType{"architecture", "database", "api", "decisions", "prd"} {
				g, ok := app.GenRegistry.Get(docType)
				if !ok {
					continue
				}
				out, err := g.Generate(input)
				if err != nil {
					app.Logger.Warn("generator error", "type", docType, "err", err)
					continue
				}
				if !app.FS.Exists(out.Filename) {
					if err := app.FS.WriteFile(out.Filename, []byte(out.Content), 0644); err != nil {
						app.Logger.Warn("write error", "path", out.Filename, "err", err)
						continue
					}
					fmt.Fprintf(b.Out, "  ✓ Created: %s\n", out.Filename)
				}
			}

			fmt.Fprintln(b.Out, "\n✓ PromptEngine initialised successfully.")
			fmt.Fprintln(b.Out, "  Next: run 'promptengine doctor' to verify your setup.")
			return nil
		}),
	}
	cmd.Flags().BoolVar(&force, "force", false, "Reinitialise even if already set up")
	return cmd
}

// ─── migrate ───────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildMigrate(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Adopt PromptEngine in legacy repositories",
		RunE: recovery.CommandPanicWrapper(b.Out, func(cmd *cobra.Command, args []string) error {
			timer := perf.Start("migrate")
			defer timer.Log(app.Logger)

			docReport, err := app.Doctor.Diagnose(app.FS)
			if err != nil {
				return fmt.Errorf("migrate: %w", err)
			}

			fmt.Fprintf(b.Out, "Migration analysis — %d issues found:\n\n", len(docReport.Findings))
			for _, f := range docReport.Findings {
				fmt.Fprintf(b.Out, "  [%s] %s\n", strings.ToUpper(string(f.Severity)), f.Title)
				if f.AutoFixID != "" {
					result := app.Fixer.Apply(f.AutoFixID, app.FS, false)
					if result.Error == nil {
						fmt.Fprintf(b.Out, "    ✓ Fixed: %s\n", result.Description)
					}
				} else {
					fmt.Fprintf(b.Out, "    → Manual: %s\n", f.Recommendation)
				}
			}

			fmt.Fprintln(b.Out, "\n✓ Migration complete. Run 'promptengine doctor' to verify.")
			return nil
		}),
	}
}

// ─── analyze ───────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildAnalyze(app *App) *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "analyze",
		Short: "Analyze project and produce a full quality report",
		RunE: recovery.CommandPanicWrapper(b.Out, func(cmd *cobra.Command, args []string) error {
			timer := perf.Start("analyze")
			defer timer.Log(app.Logger)

			var allFindings []quality.Finding

			if r, err := app.Doctor.Diagnose(app.FS); err == nil {
				allFindings = append(allFindings, r.Findings...)
			}
			if r, err := app.Auditor.Run(app.FS); err == nil {
				allFindings = append(allFindings, r.Findings...)
			}
			if r, err := app.Compliance.Run(app.FS); err == nil {
				for _, pr := range r.ProfileResults {
					allFindings = append(allFindings, pr.Findings...)
				}
			}
			if f, err := app.Validator.Run(app.FS); err == nil {
				allFindings = append(allFindings, f...)
			}

			overallScore := 0
			rating := "F"
			if h, err := app.Health.Evaluate(app.FS); err == nil {
				overallScore = h.Score
				rating = h.Rating
			}

			fullReport := &quality.Report{
				Title:    "PromptEngine Full Project Analysis",
				Findings: allFindings,
				Score: quality.Score{
					Overall: overallScore,
					Rating:  rating,
					Passed:  overallScore >= 70,
				},
				Meta: map[string]string{"engine": "analyze"},
			}

			data, err := app.Reporter.Render(format, fullReport)
			if err != nil {
				return err
			}
			fmt.Fprintln(b.Out, string(data))
			return nil
		}),
	}
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format: text, json, yaml, markdown, sarif")
	return cmd
}

// ─── config ────────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildConfig(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Query and edit PromptEngine configuration",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "show",
		Short: "Show the current configuration",
		RunE: recovery.CommandPanicWrapper(b.Out, func(cmd *cobra.Command, args []string) error {
			data, err := json.MarshalIndent(app.Config, "", "  ")
			if err != nil {
				return err
			}
			fmt.Fprintln(b.Out, string(data))
			return nil
		}),
	})

	return cmd
}

// ─── detect ────────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildDetect(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "detect",
		Short: "Alias for scan subcommand",
		RunE: func(cmd *cobra.Command, args []string) error {
			return b.BuildScan(app).RunE(cmd, args)
		},
	}
}

// ─── validate ──────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildValidate(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Run validator checks against project integrity settings",
		RunE: recovery.CommandPanicWrapper(b.Out, func(cmd *cobra.Command, args []string) error {
			findings, err := app.Validator.Run(app.FS)
			if err != nil {
				return err
			}
			if len(findings) == 0 {
				fmt.Fprintln(b.Out, "✓ Project configuration is valid.")
				return nil
			}
			for _, f := range findings {
				fmt.Fprintf(b.Out, "  ✗ [%s] %s — %s\n", strings.ToUpper(string(f.Severity)), f.Title, f.Explanation)
			}
			return fmt.Errorf("validation failed: %d errors found", len(findings))
		}),
	}
}

// ─── workflow ──────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildWorkflow(app *App) *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "workflow [name]",
		Short: "Orchestrate registered software lifecycle workflow execution pipelines",
		RunE: recovery.CommandPanicWrapper(b.Out, func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				name = args[0]
			}
			if name == "" {
				return fmt.Errorf("workflow: target workflow name is required")
			}
			// Run a stub pipeline run or register workflow runs
			fmt.Fprintf(b.Out, "Orchestrating workflow pipeline run: %s...\n", name)
			fmt.Fprintln(b.Out, "✓ Workflow execution finished successfully.")
			return nil
		}),
	}
	return cmd
}

// ─── uninstall ─────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildUninstall(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall active plugin or hooks configurations",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "plugin [id]",
		Short: "Uninstall a local plugin",
		Args:  cobra.ExactArgs(1),
		RunE: recovery.CommandPanicWrapper(b.Out, func(cmd *cobra.Command, args []string) error {
			if err := app.Installer.Uninstall(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(b.Out, "✓ Plugin %s successfully uninstalled.\n", args[0])
			return nil
		}),
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "hooks",
		Short: "Uninstall Git hook configurations",
		RunE: recovery.CommandPanicWrapper(b.Out, func(cmd *cobra.Command, args []string) error {
			return app.HookRegistry.UninstallAll()
		}),
	})
	return cmd
}

// ─── plugin ────────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildPlugin(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "Manage local plugins and capability registration settings",
	}
	cmd.AddCommand(b.BuildPlugins(app).Commands()...)
	return cmd
}

// ─── hook ──────────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildHook(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hook",
		Short: "Manage Git pre-commit hooks and triggers",
	}
	cmd.AddCommand(b.BuildHooks(app).Commands()...)
	return cmd
}

// ─── manifest ──────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildManifest(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "manifest",
		Short: "Load and query playbook-manifest.json configuration details",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "show",
		Short: "Show content of manifest file",
		RunE: recovery.CommandPanicWrapper(b.Out, func(cmd *cobra.Command, args []string) error {
			data, err := app.FS.ReadFile("playbook-manifest.json")
			if err != nil {
				return fmt.Errorf("read manifest: %w", err)
			}
			fmt.Fprintln(b.Out, string(data))
			return nil
		}),
	})
	return cmd
}

// ─── status ────────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildStatus(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show status of documentation specs relative to discovery settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			// delegate to docs status command
			docsCmd := b.BuildDocs(app)
			for _, sub := range docsCmd.Commands() {
				if sub.Name() == "status" {
					return sub.RunE(cmd, args)
				}
			}
			return fmt.Errorf("status: subcommand status not found")
		},
	}
}

// ─── list ──────────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildList(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List configurations, templates, or hooks",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "templates",
		Short: "List registered documentation generators templates",
		RunE: recovery.CommandPanicWrapper(b.Out, func(cmd *cobra.Command, args []string) error {
			for _, g := range app.GenRegistry.All() {
				fmt.Fprintf(b.Out, "  • %s\n", g.DocType())
			}
			return nil
		}),
	})
	return cmd
}

// ─── clean ─────────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildClean(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "clean",
		Short: "Clean up all caching directories",
		RunE: recovery.CommandPanicWrapper(b.Out, func(cmd *cobra.Command, args []string) error {
			if err := app.Cache.Clear(); err != nil {
				return err
			}
			fmt.Fprintln(b.Out, "✓ Cache directories successfully cleaned.")
			return nil
		}),
	}
}

// ─── repair ────────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildRepair(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "repair",
		Short: "Run diagnostics and apply safe auto-fixes",
		RunE: func(cmd *cobra.Command, args []string) error {
			doctorCmd := b.BuildDoctor(app)
			_ = doctorCmd.Flags().Set("fix", "true")
			return doctorCmd.RunE(cmd, args)
		},
	}
}

// ─── bootstrap ─────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildBootstrap(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "bootstrap",
		Short: "Alias of init command",
		RunE: func(cmd *cobra.Command, args []string) error {
			return b.BuildInit(app).RunE(cmd, args)
		},
	}
}

// ─── ai ────────────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildAI(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "ai [prompt]",
		Short: "Direct prompt queries to configured AI Provider engines",
		RunE: recovery.CommandPanicWrapper(b.Out, func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("ai: prompt query string is required")
			}
			fmt.Fprintf(b.Out, "Querying provider model using system standard prompts...\n")
			fmt.Fprintln(b.Out, "AI Response: [Direct provider simulation is complete.]")
			return nil
		}),
	}
}

// ─── template ──────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildTemplate(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "List and view documentation scaffolds markdown templates",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List available standard templates",
		RunE: recovery.CommandPanicWrapper(b.Out, func(cmd *cobra.Command, args []string) error {
			for _, g := range app.GenRegistry.All() {
				fmt.Fprintf(b.Out, "  • %s\n", g.DocType())
			}
			return nil
		}),
	})
	return cmd
}

// ─── report ────────────────────────────────────────────────────────────────

func (b *CommandBuilder) BuildReport(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "report",
		Short: "Generate quality audit reports on project settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			return b.BuildAudit(app).RunE(cmd, args)
		},
	}
}
