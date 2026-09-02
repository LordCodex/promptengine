package rulesources

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/LordCodex/promptengine/pkg/manifest"
)

// BuildManifest returns a manifest overlay backed by the pinned, locally synced
// authoritative rule snapshots for the most-specific matching stack profile.
// It refuses to build a partial overlay: if any required source or entrypoint is
// missing, PromptEngine keeps using its bundled playbooks as a safe fallback.
func (r *Resolver) BuildManifest(technologies []string) (*manifest.Manifest, *Resolution, error) {
	resolution, err := r.Resolve(technologies, "")
	if err != nil {
		return nil, nil, err
	}
	if resolution.ProfileID == "" {
		return nil, resolution, nil
	}
	if len(resolution.MissingSources) > 0 || len(resolution.MissingEntrypoints) > 0 {
		return nil, resolution, nil
	}
	profile, ok := profileByID(r.Profiles, resolution.ProfileID)
	if !ok {
		return nil, resolution, fmt.Errorf("matched rule profile %q disappeared", resolution.ProfileID)
	}

	out := &manifest.Manifest{
		Metadata: manifest.ProjectMetadata{
			Name:          "authoritative-rules:" + profile.ID,
			Version:       "1",
			SchemaVersion: manifest.SupportedSchemaVersion,
		},
	}
	playbookIDsBySource := map[string][]string{}
	for _, sourceID := range resolution.SourceIDs {
		source := r.Registry.Sources[sourceID]
		snapshot, err := LoadSnapshot(r.FS, r.cacheRoot(), sourceID, source)
		if err != nil {
			return nil, resolution, fmt.Errorf("load authoritative snapshot %q: %w", sourceID, err)
		}
		required := make(map[string]bool, len(profile.RequiredRuleEntrypoints[sourceID]))
		for _, entrypoint := range profile.RequiredRuleEntrypoints[sourceID] {
			required[entrypoint] = true
		}
		for _, rulePath := range snapshot.Files {
			id := authoritativePlaybookID(sourceID, rulePath, required[rulePath])
			out.Playbooks = append(out.Playbooks, manifest.PlaybookDefinition{
				ID:          id,
				Name:        sourceID + ":" + rulePath,
				Category:    authoritativeCategory(sourceID, rulePath),
				Location:    path.Join(r.cacheRoot(), sourceID, source.Ref, rulePath),
				Description: "Pinned authoritative rule from " + source.Repository + "@" + source.Ref,
				Priority:    100,
			})
			playbookIDsBySource[sourceID] = append(playbookIDsBySource[sourceID], id)
		}
	}

	for _, technology := range profile.Match.RequiredTechnologies {
		var related []string
		for _, sourceID := range resolution.SourceIDs {
			if sourceID == "universal" || sourceOwnsTechnology(r.Registry.Sources[sourceID], technology) {
				related = append(related, playbookIDsBySource[sourceID]...)
			}
		}
		if len(related) == 0 {
			continue
		}
		sort.Strings(related)
		out.Technologies = append(out.Technologies, manifest.TechnologyDefinition{
			ID:               strings.ToLower(technology),
			Framework:        technology,
			RelatedStandards: related,
		})
	}
	return out, resolution, nil
}

func profileByID(profiles []Profile, id string) (Profile, bool) {
	for _, profile := range profiles {
		if profile.ID == id {
			return profile, true
		}
	}
	return Profile{}, false
}

func sourceOwnsTechnology(source Source, technology string) bool {
	needle := strings.ToLower(strings.TrimSpace(technology))
	for _, owned := range source.Owns {
		value := strings.ToLower(strings.TrimSpace(owned))
		if value == needle || strings.Contains(value, needle) || strings.Contains(needle, value) {
			return true
		}
	}
	return false
}

func authoritativePlaybookID(sourceID, rulePath string, required bool) string {
	var b strings.Builder
	b.WriteString("authoritative-")
	b.WriteString(sourceID)
	b.WriteByte('-')
	for _, r := range strings.ToLower(rulePath) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		default:
			if b.Len() > 0 && !strings.HasSuffix(b.String(), "-") {
				b.WriteByte('-')
			}
		}
	}
	id := strings.Trim(b.String(), "-")
	if required {
		id += "-engineering-standard"
	}
	return id
}

func authoritativeCategory(sourceID, rulePath string) manifest.PlaybookCategory {
	lower := strings.ToLower(rulePath)
	switch {
	case strings.Contains(lower, "security"), strings.Contains(lower, "auth"):
		return manifest.CategorySecurity
	case strings.Contains(lower, "performance"), strings.Contains(lower, "cache"):
		return manifest.CategoryPerformance
	case strings.Contains(lower, "test"):
		return manifest.CategoryChecklist
	case sourceID == "universal":
		return manifest.CategoryCore
	default:
		return manifest.CategoryStacks
	}
}
