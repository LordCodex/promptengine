package app

import (
	"fmt"

	"github.com/LordCodex/promptengine/internal/output"
	"github.com/LordCodex/promptengine/internal/version"
	"github.com/spf13/cobra"
)

type versionResult struct {
	Version string `json:"version" yaml:"version"`
}

func NewRootCommand(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "promptengine",
		Short:         "PromptEngine CLI foundation",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
			showVersion, err := lc.Cmd.Flags().GetBool("version")
			if err != nil {
				return err
			}
			if showVersion {
				if !app.Config.CLI.JSON {
					fmt.Fprintln(app.Out, version.String())
					return nil
				}
				return app.Renderer.Render(app.Out, versionResult{Version: version.String()})
			}
			if app.Config.CLI.JSON {
				return app.Renderer.Render(app.Out, map[string]string{"status": "ready"})
			}
			return nil
		}),
	}

	cmd.SetOut(app.Out)
	cmd.SetErr(app.Err)
	cmd.PersistentFlags().BoolVarP(&app.Config.CLI.Verbose, "verbose", "v", app.Config.CLI.Verbose, "enable verbose logging")
	cmd.PersistentFlags().BoolVar(&app.Config.CLI.Debug, "debug", app.Config.CLI.Debug, "enable debug logging")
	cmd.PersistentFlags().BoolVar(&app.Config.CLI.JSON, "json", app.Config.CLI.JSON, "render command output as JSON")
	cmd.PersistentFlags().StringVar(&app.Config.CLI.Config, "config", app.Config.CLI.Config, "path to project configuration file")
	cmd.PersistentFlags().String("output", "", "render command output as text, json, or yaml")
	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		format, err := cmd.Root().PersistentFlags().GetString("output")
		if err != nil {
			return err
		}
		switch format {
		case "json":
			app.Config.CLI.JSON = true
			app.Renderer = output.NewConfiguredRenderer(output.FormatJSON, false, app.Config.CLI.Verbose)
		case "yaml":
			app.Renderer = output.NewConfiguredRenderer(output.FormatYAML, false, app.Config.CLI.Verbose)
		}
		return nil
	}

	cmd.Flags().Bool("version", false, "print version and exit")

	cmd.AddCommand(&cobra.Command{
		Use:           "version",
		Short:         "Print PromptEngine version information",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
			if !app.Config.CLI.JSON {
				fmt.Fprintln(app.Out, version.String())
				return nil
				}
				return app.Renderer.Render(app.Out, versionResult{Version: version.String()})
		}),
	})
	cmd.AddCommand(NewRulesCommand(app))
	cmd.AddCommand(NewProductionCommands(app)...)

	return cmd
}
