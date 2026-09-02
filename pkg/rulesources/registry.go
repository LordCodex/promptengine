package rulesources

import (
	"fmt"
	"io/fs"
	"strings"

	"gopkg.in/yaml.v3"
)

type Registry struct {
	Version            int                   `yaml:"version"`
	UpdatedAt          string                `yaml:"updated_at"`
	Sources            map[string]Source     `yaml:"sources"`
	PromptEngine       PromptEngineOwnership `yaml:"promptengine"`
	PreservationPolicy PreservationPolicy    `yaml:"preservation_policy"`
}

type Source struct {
	Repository string   `yaml:"repository"`
	Ref        string   `yaml:"ref"`
	Role       string   `yaml:"role"`
	Owns       []string `yaml:"owns"`
	Inherits   []string `yaml:"inherits"`
}

type PromptEngineOwnership struct {
	Role string   `yaml:"role"`
	Owns []string `yaml:"owns"`
}

type PreservationPolicy struct {
	Classifications []string `yaml:"classifications"`
	Invariant       string   `yaml:"invariant"`
}

type Profile struct {
	ID                      string                  `yaml:"id"`
	Version                 int                     `yaml:"version"`
	Match                   ProfileMatch            `yaml:"match"`
	Inheritance             []string                `yaml:"inheritance"`
	Presentation            ProfilePresentation     `yaml:"presentation"`
	RequiredRuleEntrypoints map[string][]string     `yaml:"required_rule_entrypoints"`
	ResolutionPolicy        ProfileResolutionPolicy `yaml:"resolution_policy"`
	Precedence              []string                `yaml:"precedence"`
	Notes                   []string                `yaml:"notes"`
}

type ProfileMatch struct {
	RequiredTechnologies []string `yaml:"required_technologies"`
}

type ProfilePresentation struct {
	BackendAuthority string `yaml:"backend_authority"`
	Integration      string `yaml:"integration"`
	Frontend         string `yaml:"frontend"`
}

type ProfileResolutionPolicy struct {
	ProjectRulesLast              bool `yaml:"project_rules_last"`
	SelectTaskRelevantRules       bool `yaml:"select_task_relevant_rules"`
	ConcatenateEntireRepositories bool `yaml:"concatenate_entire_repositories"`
	PreservePromptEngineBridges   bool `yaml:"preserve_promptengine_bridges"`
	PreserveProjectSpecificRules  bool `yaml:"preserve_project_specific_rules"`
}

func LoadRegistry(data []byte) (*Registry, error) {
	var registry Registry
	if err := yaml.Unmarshal(data, &registry); err != nil {
		return nil, fmt.Errorf("decode rule source registry: %w", err)
	}
	if err := registry.Validate(); err != nil {
		return nil, err
	}
	return &registry, nil
}

func LoadRegistryFS(fsys fs.FS, path string) (*Registry, error) {
	data, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, fmt.Errorf("read rule source registry %q: %w", path, err)
	}
	return LoadRegistry(data)
}

func LoadProfile(data []byte) (*Profile, error) {
	var profile Profile
	if err := yaml.Unmarshal(data, &profile); err != nil {
		return nil, fmt.Errorf("decode rule profile: %w", err)
	}
	if err := profile.Validate(); err != nil {
		return nil, err
	}
	return &profile, nil
}

func LoadProfileFS(fsys fs.FS, path string) (*Profile, error) {
	data, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, fmt.Errorf("read rule profile %q: %w", path, err)
	}
	return LoadProfile(data)
}

func (r *Registry) Validate() error {
	if r == nil {
		return fmt.Errorf("rule source registry is nil")
	}
	if r.Version <= 0 {
		return fmt.Errorf("rule source registry version must be positive")
	}
	if len(r.Sources) == 0 {
		return fmt.Errorf("rule source registry has no sources")
	}
	for id, source := range r.Sources {
		if strings.TrimSpace(id) == "" {
			return fmt.Errorf("rule source id must not be empty")
		}
		if strings.TrimSpace(source.Repository) == "" {
			return fmt.Errorf("rule source %q has no repository", id)
		}
		if strings.TrimSpace(source.Ref) == "" {
			return fmt.Errorf("rule source %q has no pinned ref", id)
		}
		for _, parent := range source.Inherits {
			if _, ok := r.Sources[parent]; !ok {
				return fmt.Errorf("rule source %q inherits unknown source %q", id, parent)
			}
		}
	}
	for id := range r.Sources {
		if _, err := r.Resolve(id); err != nil {
			return err
		}
	}
	return nil
}

func (p *Profile) Validate() error {
	if p == nil {
		return fmt.Errorf("rule profile is nil")
	}
	if strings.TrimSpace(p.ID) == "" {
		return fmt.Errorf("rule profile id must not be empty")
	}
	if p.Version <= 0 {
		return fmt.Errorf("rule profile %q version must be positive", p.ID)
	}
	if len(p.Inheritance) == 0 {
		return fmt.Errorf("rule profile %q has no inheritance sources", p.ID)
	}
	return nil
}

// Resolve returns source IDs in dependency-first order and includes the requested source.
// For example, resolving Laravel yields universal, PHP, Laravel with the current registry.
func (r *Registry) Resolve(sourceID string) ([]string, error) {
	if r == nil {
		return nil, fmt.Errorf("rule source registry is nil")
	}
	if _, ok := r.Sources[sourceID]; !ok {
		return nil, fmt.Errorf("unknown rule source %q", sourceID)
	}

	state := map[string]uint8{}
	seen := map[string]bool{}
	var out []string
	var visit func(string) error
	visit = func(id string) error {
		switch state[id] {
		case 1:
			return fmt.Errorf("rule source inheritance cycle detected at %q", id)
		case 2:
			return nil
		}
		state[id] = 1
		for _, parent := range r.Sources[id].Inherits {
			if _, ok := r.Sources[parent]; !ok {
				return fmt.Errorf("rule source %q inherits unknown source %q", id, parent)
			}
			if err := visit(parent); err != nil {
				return err
			}
		}
		state[id] = 2
		if !seen[id] {
			seen[id] = true
			out = append(out, id)
		}
		return nil
	}
	if err := visit(sourceID); err != nil {
		return nil, err
	}
	return out, nil
}

// ResolveProfile expands every profile source through the registry inheritance graph,
// preserving profile order while de-duplicating inherited sources.
func (r *Registry) ResolveProfile(profile *Profile) ([]string, error) {
	if err := profile.Validate(); err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	var out []string
	for _, sourceID := range profile.Inheritance {
		resolved, err := r.Resolve(sourceID)
		if err != nil {
			return nil, fmt.Errorf("resolve profile %q: %w", profile.ID, err)
		}
		for _, id := range resolved {
			if seen[id] {
				continue
			}
			seen[id] = true
			out = append(out, id)
		}
	}
	return out, nil
}

func (p *Profile) Matches(technologies []string) bool {
	available := map[string]bool{}
	for _, technology := range technologies {
		available[strings.ToLower(strings.TrimSpace(technology))] = true
	}
	for _, required := range p.Match.RequiredTechnologies {
		if !available[strings.ToLower(strings.TrimSpace(required))] {
			return false
		}
	}
	return true
}
