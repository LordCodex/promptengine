package rulesources

import (
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

const DefaultProfilesDir = "sources/profiles"

// LoadProfilesFS loads every YAML rule profile from dir in deterministic order.
func LoadProfilesFS(fsys fs.FS, dir string) ([]Profile, error) {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil, fmt.Errorf("read rule profiles %q: %w", dir, err)
	}

	var profiles []Profile
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := strings.ToLower(entry.Name())
		if !strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml") {
			continue
		}
		profile, err := LoadProfileFS(fsys, dir+"/"+entry.Name())
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, *profile)
	}

	sort.SliceStable(profiles, func(i, j int) bool {
		if len(profiles[i].Match.RequiredTechnologies) == len(profiles[j].Match.RequiredTechnologies) {
			return profiles[i].ID < profiles[j].ID
		}
		return len(profiles[i].Match.RequiredTechnologies) > len(profiles[j].Match.RequiredTechnologies)
	})
	return profiles, nil
}

// MatchProfile returns the most-specific matching profile. Profiles with more
// required technologies win, making compound stacks such as Laravel + Inertia
// + Vue resolve before generic Vue or Laravel profiles.
func MatchProfile(profiles []Profile, technologies []string) (*Profile, bool) {
	ordered := append([]Profile(nil), profiles...)
	sort.SliceStable(ordered, func(i, j int) bool {
		if len(ordered[i].Match.RequiredTechnologies) == len(ordered[j].Match.RequiredTechnologies) {
			return ordered[i].ID < ordered[j].ID
		}
		return len(ordered[i].Match.RequiredTechnologies) > len(ordered[j].Match.RequiredTechnologies)
	})
	for i := range ordered {
		if ordered[i].Matches(technologies) {
			matched := ordered[i]
			return &matched, true
		}
	}
	return nil, false
}
