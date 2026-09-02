package app

import (
	"fmt"
	"strings"

	"github.com/LordCodex/promptengine/internal/domain/discovery"
	"github.com/LordCodex/promptengine/pkg/manifest"
)

func detectedRuleTechnologies(model *discovery.ProjectModel) []string {
	if model == nil {
		return nil
	}
	seen := map[string]bool{}
	var out []string
	for _, value := range append(append([]string{}, model.Frameworks...), model.Languages...) {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		key := strings.ToLower(value)
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, value)
	}
	return out
}

// activateAuthoritativeRules overlays the bundled PromptEngine technology
// mappings only when every pinned source required by the detected stack has
// been synced successfully. An incomplete cache deliberately leaves the bundled
// playbooks active as the preservation-first fallback.
func (a *App) activateAuthoritativeRules(model *discovery.ProjectModel) error {
	if a == nil || a.RuleSources == nil || a.Manifest == nil || model == nil {
		return nil
	}
	overlay, resolution, err := a.RuleSources.BuildManifest(detectedRuleTechnologies(model))
	if err != nil {
		return fmt.Errorf("resolve authoritative rules: %w", err)
	}
	if overlay == nil {
		if resolution != nil && resolution.ProfileID != "" && a.Container != nil && a.Container.Logger != nil {
			a.Container.Logger.Debug("authoritative rules not fully synced; using bundled fallback",
				"profile", resolution.ProfileID,
				"missing_sources", resolution.MissingSources,
				"missing_entrypoints", resolution.MissingEntrypoints,
			)
		}
		return nil
	}
	if err := a.Manifest.Register("authoritative-rules", manifest.SourceOrganization, overlay); err != nil {
		return fmt.Errorf("register authoritative rules: %w", err)
	}
	return nil
}
