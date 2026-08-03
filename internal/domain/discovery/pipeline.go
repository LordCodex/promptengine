package discovery

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/LordCodex/promptengine/internal/eventbus"
	"github.com/LordCodex/promptengine/internal/filesystem"
	"github.com/LordCodex/promptengine/pkg/manifest"
)

type Stage interface {
	Name() string
	Execute(ctx context.Context, fs filesystem.FileSystem, pm *ProjectModel) error
}

type Pipeline struct {
	stages   []Stage
	events   *eventbus.EventBus
	manifest *manifest.Engine
	rules    *RuleRegistry
}

type Option func(*Pipeline)

func WithEventBus(events *eventbus.EventBus) Option {
	return func(p *Pipeline) { p.events = events }
}

func WithManifestEngine(engine *manifest.Engine) Option {
	return func(p *Pipeline) { p.manifest = engine }
}

func WithRuleRegistry(rules *RuleRegistry) Option {
	return func(p *Pipeline) { p.rules = rules }
}

func NewPipeline(opts ...Option) *Pipeline {
	p := &Pipeline{
		stages: make([]Stage, 0),
		rules:  NewRuleRegistry(),
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func NewDefaultPipeline(events *eventbus.EventBus, engine *manifest.Engine) *Pipeline {
	p := NewPipeline(WithEventBus(events), WithManifestEngine(engine))
	p.Register(&RepositoryScanStage{}, &PromptEngineStage{}, &TechStage{Registry: p.rules, Manifest: engine, Events: events}, &ArchStage{}, &DocsStage{})
	return p
}

func (p *Pipeline) Register(stages ...Stage) {
	p.stages = append(p.stages, stages...)
}

func (p *Pipeline) RegisterRules(rules ...DetectionRule) {
	if p.rules == nil {
		p.rules = NewRuleRegistry()
	}
	p.rules.Register(rules...)
}

func (p *Pipeline) Execute(ctx context.Context, fs filesystem.FileSystem, rootPath string) (*ProjectModel, error) {
	if rootPath == "" {
		rootPath = "."
	}
	pm := NewProjectModel(filepath.Clean(rootPath))
	p.publish(eventbus.ProjectDiscoveryStarted, "project discovery started", rootPath)

	if len(p.stages) == 0 {
		p.Register(&RepositoryScanStage{}, &PromptEngineStage{}, &TechStage{Registry: p.rules, Manifest: p.manifest, Events: p.events}, &ArchStage{}, &DocsStage{})
	}

	for _, stage := range p.stages {
		if err := stage.Execute(ctx, fs, pm); err != nil {
			p.publish(eventbus.ProjectDiscoveryFailed, "project discovery failed", map[string]any{"stage": stage.Name(), "error": err.Error()})
			return nil, fmt.Errorf("discovery stage %s: %w", stage.Name(), err)
		}
		pm.SyncLegacyFields()
	}

	p.consolidateClassifications(pm)
	pm.Project.DetectedType = detectedType(pm)
	p.publish(eventbus.ProjectDetected, "project detected", pm.Project.DetectedType)
	p.publish(eventbus.ProjectDiscoveryCompleted, "project discovery completed", pm)
	return pm, nil
}

func (p *Pipeline) publish(t eventbus.EventType, msg string, payload any) {
	if p.events == nil {
		return
	}
	p.events.Publish(eventbus.Event{Type: t, Message: msg, Payload: payload})
}

func (p *Pipeline) consolidateClassifications(pm *ProjectModel) {
	if len(pm.Languages) == 0 && len(pm.Frameworks) == 0 {
		pm.Classifications = addUniqueClass(pm.Classifications, ClassGreenfield)
	} else {
		pm.Classifications = addUniqueClass(pm.Classifications, ClassExisting)
	}
	if pm.Repository.IsMonorepo {
		pm.Classifications = addUniqueClass(pm.Classifications, ClassMonorepo)
	}
	if pm.PromptEngine.Installed {
		pm.Classifications = addUniqueClass(pm.Classifications, ClassPromptEngine)
	} else if len(pm.Languages) > 0 {
		pm.Classifications = addUniqueClass(pm.Classifications, ClassNonPromptEngine)
	}
	for _, f := range pm.Frameworks {
		switch f {
		case "Laravel", "Django", "Rails", "Express", "Spring":
			pm.Classifications = addUniqueClass(pm.Classifications, ClassBackendAPI)
			pm.Architecture.Backend = true
		case "Nuxt", "Next.js":
			pm.Classifications = addUniqueClass(pm.Classifications, ClassSSRApplication)
			pm.Architecture.Frontend = true
		case "Vue", "React", "Angular":
			pm.Classifications = addUniqueClass(pm.Classifications, ClassFrontendSPA)
			pm.Architecture.Frontend = true
		case "Flutter", "React Native":
			pm.Classifications = addUniqueClass(pm.Classifications, ClassMobileApplication)
			pm.Architecture.Mobile = true
		}
	}
	if len(pm.Technology.Infrastructure) > 0 {
		pm.Architecture.Infrastructure = true
	}
}

func detectedType(pm *ProjectModel) string {
	for _, c := range []ProjectClassification{ClassMonorepo, ClassHybrid, ClassMobileApplication, ClassSSRApplication, ClassFrontendSPA, ClassBackendAPI, ClassCLIApplication, ClassGreenfield} {
		if pm.HasClassification(c) {
			return string(c)
		}
	}
	if len(pm.Frameworks) > 1 {
		return string(ClassHybrid)
	}
	return string(ClassExisting)
}

func addUniqueClass(items []ProjectClassification, item ProjectClassification) []ProjectClassification {
	for _, existing := range items {
		if existing == item {
			return items
		}
	}
	return append(items, item)
}

func hasAnyPath(paths []string, names ...string) bool {
	for _, path := range paths {
		base := filepath.Base(path)
		for _, name := range names {
			if strings.EqualFold(base, name) || strings.EqualFold(path, name) {
				return true
			}
		}
	}
	return false
}
