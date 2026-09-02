package app

import (
	"fmt"

	"github.com/LordCodex/promptengine/pkg/rulesources"
	"github.com/spf13/cobra"
)

type ruleSourceStatus struct {
	ID         string `json:"id" yaml:"id"`
	Repository string `json:"repository" yaml:"repository"`
	Ref        string `json:"ref" yaml:"ref"`
	Synced     bool   `json:"synced" yaml:"synced"`
}

type ruleResolutionResult struct {
	Profile            string             `json:"profile" yaml:"profile"`
	DetectedTechnology []string           `json:"detected_technologies" yaml:"detected_technologies"`
	Sources            []ruleSourceStatus `json:"sources" yaml:"sources"`
	MissingSources     []string           `json:"missing_sources,omitempty" yaml:"missing_sources,omitempty"`
	MissingEntrypoints []string           `json:"missing_entrypoints,omitempty" yaml:"missing_entrypoints,omitempty"`
	FallbackActive     bool               `json:"bundled_fallback_active" yaml:"bundled_fallback_active"`
}

func NewRulesCommand(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rules",
		Short: "Resolve and synchronize authoritative engineering rules",
	}
	cmd.AddCommand(newRulesResolveCommand(app))
	cmd.AddCommand(newRulesSyncCommand(app))
	cmd.AddCommand(newRulesStatusCommand(app))
	return cmd
}

func newRulesResolveCommand(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "resolve",
		Short: "Show the authoritative rule profile for the detected project stack",
		RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
			result, err := app.ruleResolution(lc)
			if err != nil {
				return err
			}
			return app.Renderer.Render(app.Out, result)
		}),
	}
}

func newRulesStatusCommand(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show whether the detected stack's pinned rule sources are synchronized",
		RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
			result, err := app.ruleResolution(lc)
			if err != nil {
				return err
			}
			return app.Renderer.Render(app.Out, result)
		}),
	}
}

func newRulesSyncCommand(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "Fetch pinned authoritative rules for the detected project stack",
		Long: "Fetch the pinned rule repositories for the detected stack into .promptengine/rules. " +
			"For private GitHub repositories, set GITHUB_TOKEN or GH_TOKEN with read access. Tokens are read from the environment and are never written to the rule cache.",
		RunE: app.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
			if app.RuleSources == nil {
				return fmt.Errorf("authoritative rule source service is unavailable")
			}
			technologies := detectedRuleTechnologies(lc.Model)
			report, err := app.RuleSources.Sync(lc.Ctx, technologies)
			if err != nil {
				return err
			}
			if err := app.activateAuthoritativeRules(lc.Model); err != nil {
				return err
			}
			return app.Renderer.Render(app.Out, report)
		}),
	}
}

func (a *App) ruleResolution(lc *LifecycleContext) (*ruleResolutionResult, error) {
	if a.RuleSources == nil {
		return nil, fmt.Errorf("authoritative rule source service is unavailable")
	}
	technologies := detectedRuleTechnologies(lc.Model)
	resolution, err := a.RuleSources.Resolve(technologies, "")
	if err != nil {
		return nil, err
	}
	result := &ruleResolutionResult{DetectedTechnology: technologies}
	if resolution == nil || resolution.ProfileID == "" {
		result.FallbackActive = true
		return result, nil
	}
	result.Profile = resolution.ProfileID
	result.MissingSources = resolution.MissingSources
	result.MissingEntrypoints = resolution.MissingEntrypoints
	result.FallbackActive = len(resolution.MissingSources) > 0 || len(resolution.MissingEntrypoints) > 0
	missing := map[string]bool{}
	for _, id := range resolution.MissingSources {
		missing[id] = true
	}
	for _, id := range resolution.SourceIDs {
		source := a.RuleSources.Registry.Sources[id]
		result.Sources = append(result.Sources, ruleSourceStatus{
			ID:         id,
			Repository: source.Repository,
			Ref:        source.Ref,
			Synced:     !missing[id] && sourceSnapshotExists(a.FS, id, source),
		})
	}
	return result, nil
}

func sourceSnapshotExists(fsys interface{ ReadFile(string) ([]byte, error) }, sourceID string, source rulesources.Source) bool {
	// Keep the command independent of concrete filesystem implementations while
	// using the same cache layout as the rule source package.
	path := rulesources.DefaultCacheRoot + "/" + sourceID + "/" + source.Ref + "/.source.yaml"
	_, err := fsys.ReadFile(path)
	return err == nil
}
