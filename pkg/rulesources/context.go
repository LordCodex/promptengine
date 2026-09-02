package rulesources

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

type CachedRuleFile struct {
	SourceID  string
	RulePath  string
	CachePath string
	Required  bool
}

type Resolution struct {
	ProfileID          string
	SourceIDs          []string
	Files              []CachedRuleFile
	MissingSources     []string
	MissingEntrypoints []string
}

type Resolver struct {
	Registry  *Registry
	Profiles  []Profile
	FS        filesystem.FileSystem
	CacheRoot string
}

func NewResolver(registry *Registry, profiles []Profile, fsys filesystem.FileSystem) *Resolver {
	return &Resolver{Registry: registry, Profiles: profiles, FS: fsys, CacheRoot: DefaultCacheRoot}
}

func (r *Resolver) Resolve(technologies []string, intent string) (*Resolution, error) {
	if r == nil || r.Registry == nil {
		return nil, fmt.Errorf("rule resolver has no registry")
	}
	if r.FS == nil {
		return nil, fmt.Errorf("rule resolver has no filesystem")
	}
	profile, ok := MatchProfile(r.Profiles, technologies)
	if !ok {
		return &Resolution{}, nil
	}
	sourceIDs, err := r.Registry.ResolveProfile(profile)
	if err != nil {
		return nil, err
	}

	resolution := &Resolution{ProfileID: profile.ID, SourceIDs: sourceIDs}
	terms := ruleIntentTerms(intent)
	seenFiles := map[string]bool{}
	for _, sourceID := range sourceIDs {
		source := r.Registry.Sources[sourceID]
		snapshot, err := LoadSnapshot(r.FS, r.cacheRoot(), sourceID, source)
		if err != nil {
			resolution.MissingSources = append(resolution.MissingSources, sourceID)
			continue
		}

		available := map[string]bool{}
		for _, file := range snapshot.Files {
			available[file] = true
		}
		for _, required := range profile.RequiredRuleEntrypoints[sourceID] {
			if !available[required] {
				resolution.MissingEntrypoints = append(resolution.MissingEntrypoints, sourceID+":"+required)
				continue
			}
			cachePath := path.Join(r.cacheRoot(), sourceID, source.Ref, required)
			if !seenFiles[cachePath] {
				seenFiles[cachePath] = true
				resolution.Files = append(resolution.Files, CachedRuleFile{SourceID: sourceID, RulePath: required, CachePath: cachePath, Required: true})
			}
		}

		for _, file := range snapshot.Files {
			if availableEntrypoint(profile.RequiredRuleEntrypoints[sourceID], file) || !rulePathMatchesTerms(file, terms) {
				continue
			}
			cachePath := path.Join(r.cacheRoot(), sourceID, source.Ref, file)
			if seenFiles[cachePath] {
				continue
			}
			seenFiles[cachePath] = true
			resolution.Files = append(resolution.Files, CachedRuleFile{SourceID: sourceID, RulePath: file, CachePath: cachePath})
		}
	}

	sort.Strings(resolution.MissingSources)
	sort.Strings(resolution.MissingEntrypoints)
	return resolution, nil
}

func (r *Resolver) cacheRoot() string {
	if strings.TrimSpace(r.CacheRoot) == "" {
		return DefaultCacheRoot
	}
	return strings.Trim(strings.ReplaceAll(r.CacheRoot, "\\", "/"), "/")
}

func availableEntrypoint(entrypoints []string, file string) bool {
	for _, entrypoint := range entrypoints {
		if entrypoint == file {
			return true
		}
	}
	return false
}

func rulePathMatchesTerms(rulePath string, terms []string) bool {
	if len(terms) == 0 {
		return false
	}
	haystack := strings.ToLower(strings.ReplaceAll(rulePath, "_", "-"))
	for _, term := range terms {
		if strings.Contains(haystack, term) {
			return true
		}
	}
	return false
}

func ruleIntentTerms(intent string) []string {
	aliases := map[string][]string{
		"auth":          {"auth", "authentication", "authorization", "security"},
		"login":         {"auth", "authentication", "security"},
		"permission":    {"authorization", "security"},
		"role":          {"authorization", "security"},
		"database":      {"database", "data"},
		"migration":     {"database", "migration", "data"},
		"query":         {"database", "performance", "query"},
		"cache":         {"cache", "performance"},
		"queue":         {"queue", "reliability"},
		"job":           {"queue", "reliability"},
		"api":           {"api", "security"},
		"endpoint":      {"api", "routing"},
		"route":         {"routing", "api"},
		"security":      {"security", "authentication", "authorization"},
		"test":          {"test", "testing"},
		"performance":   {"performance", "cache"},
		"deploy":        {"deployment", "production", "reliability"},
		"production":    {"production", "deployment", "reliability", "observability"},
		"inertia":       {"inertia", "component"},
		"component":     {"component", "frontend"},
		"vue":           {"vue", "component", "frontend"},
		"react":         {"react", "component", "frontend"},
		"livewire":      {"livewire", "component"},
		"blade":         {"component", "blade"},
		"feature":       {"feature", "architecture"},
		"architecture":  {"architecture", "feature"},
		"observability": {"observability", "reliability"},
	}
	stop := map[string]bool{
		"with": true, "from": true, "into": true, "this": true, "that": true,
		"using": true, "change": true, "update": true, "implement": true, "create": true,
	}
	seen := map[string]bool{}
	var out []string
	for _, token := range strings.FieldsFunc(strings.ToLower(intent), func(r rune) bool {
		return (r < 'a' || r > 'z') && (r < '0' || r > '9')
	}) {
		if len(token) < 3 || stop[token] {
			continue
		}
		candidates := append([]string{token}, aliases[token]...)
		for _, candidate := range candidates {
			if !seen[candidate] {
				seen[candidate] = true
				out = append(out, candidate)
			}
		}
	}
	return out
}
