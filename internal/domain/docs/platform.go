package docs

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/LordCodex/promptengine/internal/domain/discovery"
	docsync "github.com/LordCodex/promptengine/internal/domain/docs/sync"
	doctemplate "github.com/LordCodex/promptengine/internal/domain/docs/template"
	docvalidation "github.com/LordCodex/promptengine/internal/domain/docs/validation"
	"github.com/LordCodex/promptengine/internal/domain/workflows"
	"github.com/LordCodex/promptengine/internal/eventbus"
	"github.com/LordCodex/promptengine/internal/filesystem"
	"github.com/LordCodex/promptengine/pkg/manifest"
)

type DocumentType string

const (
	DocumentREADME          DocumentType = "README"
	DocumentPRD             DocumentType = "PRD"
	DocumentArchitecture    DocumentType = "Architecture"
	DocumentDatabase        DocumentType = "Database"
	DocumentAPI             DocumentType = "API"
	DocumentDeployment      DocumentType = "Deployment"
	DocumentDecisions       DocumentType = "Decisions"
	DocumentBusinessRules   DocumentType = "Business Rules"
	DocumentTroubleshooting DocumentType = "Troubleshooting"
	DocumentProgress        DocumentType = "Progress"
)

type DocumentStatus string

const (
	DocumentStatusMissing        DocumentStatus = "missing"
	DocumentStatusValid          DocumentStatus = "valid"
	DocumentStatusInvalid        DocumentStatus = "invalid"
	DocumentStatusRequiresUpdate DocumentStatus = "requires_update"
	DocumentStatusGenerated      DocumentStatus = "generated"
)

type Document struct {
	ID                 string         `json:"id"`
	Name               string         `json:"name"`
	Type               DocumentType   `json:"type"`
	Path               string         `json:"path"`
	Version            string         `json:"version"`
	Status             DocumentStatus `json:"status"`
	GeneratedAt        time.Time      `json:"generated_at,omitempty"`
	LastSynchronizedAt time.Time      `json:"last_synchronized_at,omitempty"`
}

type GenerationRequest struct {
	DocumentID string
	Project    *discovery.ProjectModel
	Manifest   *manifest.Manifest
	Variables  map[string]string
	Overwrite  bool
}

type DocumentationReport struct {
	Documents       []Document                        `json:"documents"`
	Findings        []docvalidation.ValidationFinding `json:"findings,omitempty"`
	Recommendations []docsync.SyncRecommendation      `json:"recommendations,omitempty"`
}

type Platform struct {
	fs        filesystem.FileSystem
	registry  *DocRegistry
	templates *doctemplate.TemplateRegistry
	validator *docvalidation.Validator
	sync      *docsync.SyncEngine
	events    *eventbus.EventBus
	manifest  *manifest.Engine
}

func NewPlatform(fs filesystem.FileSystem, events *eventbus.EventBus, manifestEngine *manifest.Engine) *Platform {
	reg := NewDocRegistry()
	RegisterDefaultDocuments(reg)
	return &Platform{
		fs:        fs,
		registry:  reg,
		templates: doctemplate.NewTemplateRegistry(),
		validator: docvalidation.NewValidator(),
		sync:      docsync.NewSyncEngine(docsync.NewChangeDetector()),
		events:    events,
		manifest:  manifestEngine,
	}
}

func RegisterDefaultDocuments(reg *DocRegistry) {
	defaults := []DocumentSpec{
		{ID: "readme", Name: "README", DefaultPath: "README.md", Priority: 100, RequiredSections: []RequiredSection{{Heading: "Overview"}}, UpdateTriggers: []string{"new-technology"}},
		{ID: "prd", Name: "PRD", DefaultPath: "docs/PRD.md", Priority: 90, RequiredSections: []RequiredSection{{Heading: "Product Requirements Document"}}},
		{ID: "architecture", Name: "Architecture", DefaultPath: "docs/Architecture.md", Priority: 100, RequiredSections: []RequiredSection{{Heading: "Architecture"}}, UpdateTriggers: []string{"architecture-change", "new-service", "new-technology"}},
		{ID: "database", Name: "Database", DefaultPath: "docs/Database.md", Priority: 90, RequiredSections: []RequiredSection{{Heading: "Database"}}, UpdateTriggers: []string{"new-migration"}},
		{ID: "api", Name: "API", DefaultPath: "docs/API.md", Priority: 90, RequiredSections: []RequiredSection{{Heading: "API"}}, UpdateTriggers: []string{"new-api"}},
		{ID: "deployment", Name: "Deployment", DefaultPath: "docs/Deployment.md", Priority: 75, RequiredSections: []RequiredSection{{Heading: "Deployment"}}},
		{ID: "decisions", Name: "Decisions", DefaultPath: "docs/Decisions.md", Priority: 70, RequiredSections: []RequiredSection{{Heading: "Decisions"}}},
		{ID: "business-rules", Name: "Business Rules", DefaultPath: "docs/BusinessRules.md", Priority: 85, RequiredSections: []RequiredSection{{Heading: "Business Rules"}}},
		{ID: "troubleshooting", Name: "Troubleshooting", DefaultPath: "docs/Troubleshooting.md", Priority: 65, RequiredSections: []RequiredSection{{Heading: "Troubleshooting"}}},
		{ID: "progress", Name: "Progress", DefaultPath: "docs/Progress.md", Priority: 60, RequiredSections: []RequiredSection{{Heading: "Progress"}}},
	}
	for _, spec := range defaults {
		reg.Register(spec)
	}
}

func (p *Platform) Registry() *DocRegistry { return p.registry }

func (p *Platform) LoadTemplates() error {
	loader := doctemplate.NewLoader(p.fs, "project/templates")
	templates, err := loader.LoadAll()
	if err != nil {
		return err
	}
	for _, tmpl := range templates {
		if err := tmpl.Validate(); err != nil {
			return err
		}
		if _, exists := p.templates.Get(tmpl.ID); exists {
			continue
		}
		if err := p.templates.Register(tmpl); err != nil {
			return err
		}
	}
	return nil
}

func (p *Platform) ApplyManifestRules(m *manifest.Manifest) {
	if m == nil {
		return
	}
	for _, tmpl := range m.Templates {
		id := strings.ToLower(strings.ReplaceAll(tmpl.Name, " ", "-"))
		path := tmpl.Location
		if path == "" {
			continue
		}
		p.registry.Register(DocumentSpec{ID: id, Name: tmpl.Name, DefaultPath: outputPathForTemplate(tmpl), Priority: 50})
	}
}

func (p *Platform) Discover(project *discovery.ProjectModel) []Document {
	var docs []Document
	seen := map[string]bool{}
	for _, spec := range p.registry.All() {
		docs = append(docs, p.documentFromSpec(spec))
		seen[spec.DefaultPath] = true
	}
	if project != nil {
		for _, path := range project.Repository.DocumentationFiles {
			if seen[path] {
				continue
			}
			docs = append(docs, Document{ID: docIDFromPath(path), Name: docNameFromPath(path), Type: docTypeFromPath(path), Path: path, Status: statusForPath(p.fs, path)})
		}
	}
	sort.Slice(docs, func(i, j int) bool { return docs[i].Path < docs[j].Path })
	return docs
}

func (p *Platform) Generate(ctx context.Context, req GenerationRequest) (Document, error) {
	if err := ctx.Err(); err != nil {
		return Document{}, err
	}
	if err := p.LoadTemplates(); err != nil {
		return Document{}, err
	}
	if req.Manifest != nil {
		p.ApplyManifestRules(req.Manifest)
	}
	spec, ok := p.registry.Get(req.DocumentID)
	if !ok {
		return Document{}, fmt.Errorf("document %q is not registered", req.DocumentID)
	}
	if p.fs.Exists(spec.DefaultPath) && !req.Overwrite {
		return p.documentFromSpec(spec), nil
	}
	templateID := req.DocumentID
	vars := p.variables(req)
	renderer := doctemplate.NewRenderer(p.templates)
	content, err := renderer.Render(templateID, vars)
	if err != nil {
		content = fallbackContent(spec, vars)
	}
	if err := p.fs.WriteFile(spec.DefaultPath, []byte(content), 0644); err != nil {
		return Document{}, err
	}
	doc := p.documentFromSpec(spec)
	doc.Status = DocumentStatusGenerated
	doc.GeneratedAt = time.Now().UTC()
	doc.LastSynchronizedAt = doc.GeneratedAt
	p.publish(eventbus.DocumentationGenerated, "documentation generated", doc)
	return doc, nil
}

func (p *Platform) Validate(project *discovery.ProjectModel) (DocumentationReport, error) {
	report := DocumentationReport{Documents: p.Discover(project)}
	for _, doc := range report.Documents {
		findings, err := p.validator.Validate(p.fs, doc.Path)
		if err != nil {
			return report, err
		}
		report.Findings = append(report.Findings, findings...)
	}
	for _, finding := range report.Findings {
		if finding.Severity == docvalidation.SeverityError {
			p.publish(eventbus.DocumentationValidationFailed, "documentation validation failed", finding)
			break
		}
	}
	return report, nil
}

func (p *Platform) Sync(changedFiles []string, dryRun bool) docsync.SyncResult {
	signals := SignalsFromChangedFiles(changedFiles)
	result := p.sync.Run(signals, dryRun)
	if len(result.Pending) > 0 || len(result.Recommendations) > 0 {
		p.publish(eventbus.DocumentationSyncRequired, "documentation sync required", result)
	}
	return result
}

func (p *Platform) Update(ctx context.Context, req GenerationRequest) (Document, error) {
	req.Overwrite = true
	doc, err := p.Generate(ctx, req)
	if err != nil {
		return doc, err
	}
	p.publish(eventbus.DocumentationUpdated, "documentation updated", doc)
	return doc, nil
}

func (p *Platform) WorkflowHandlers() map[string]workflows.StepHandler {
	return map[string]workflows.StepHandler{
		"generate_documentation": workflows.StepHandlerFunc(func(ctx context.Context, step workflows.WorkflowStep, flow *workflows.FlowContext) (interface{}, error) {
			docID := step.InputMapping["document_id"]
			if docID == "" {
				docID = "architecture"
			}
			return p.Generate(ctx, GenerationRequest{DocumentID: docID, Project: flow.Project, Overwrite: true})
		}),
		"validate_documentation": workflows.StepHandlerFunc(func(ctx context.Context, step workflows.WorkflowStep, flow *workflows.FlowContext) (interface{}, error) {
			return p.Validate(flow.Project)
		}),
		"synchronize_documentation": workflows.StepHandlerFunc(func(ctx context.Context, step workflows.WorkflowStep, flow *workflows.FlowContext) (interface{}, error) {
			var files []string
			if raw, ok := flow.Inputs["changed_files"].([]string); ok {
				files = raw
			}
			return p.Sync(files, true), nil
		}),
		"update_documentation": workflows.StepHandlerFunc(func(ctx context.Context, step workflows.WorkflowStep, flow *workflows.FlowContext) (interface{}, error) {
			docID := step.InputMapping["document_id"]
			if docID == "" {
				docID = "architecture"
			}
			return p.Update(ctx, GenerationRequest{DocumentID: docID, Project: flow.Project})
		}),
	}
}

func SignalsFromChangedFiles(files []string) []docsync.ChangeSignal {
	var signals []docsync.ChangeSignal
	for _, file := range files {
		path := filepath.ToSlash(file)
		switch {
		case strings.Contains(path, "database/migrations"):
			signals = append(signals, docsync.SignalNewMigration)
		case strings.Contains(path, "routes/"), strings.Contains(strings.ToLower(path), "controller"):
			signals = append(signals, docsync.SignalNewAPI)
		case strings.Contains(strings.ToLower(path), "service"):
			signals = append(signals, docsync.SignalNewService)
		case strings.Contains(path, "Dockerfile"), strings.Contains(path, ".github/workflows"):
			signals = append(signals, docsync.SignalArchitectureChange)
		case strings.HasSuffix(path, ".md"):
			signals = append(signals, docsync.SignalNewDocumentation)
		}
	}
	return signals
}

func (p *Platform) documentFromSpec(spec DocumentSpec) Document {
	return Document{ID: spec.ID, Name: spec.Name, Type: docTypeFromName(spec.Name), Path: spec.DefaultPath, Version: "1.0.0", Status: statusForPath(p.fs, spec.DefaultPath)}
}

func statusForPath(fs filesystem.FileSystem, path string) DocumentStatus {
	if fs.Exists(path) {
		return DocumentStatusValid
	}
	return DocumentStatusMissing
}

func (p *Platform) variables(req GenerationRequest) map[string]string {
	vars := map[string]string{"project_name": "Project", "stack": ""}
	if req.Project != nil {
		if req.Project.Project.Name != "" {
			vars["project_name"] = req.Project.Project.Name
		}
		vars["stack"] = strings.Join(append(req.Project.Languages, req.Project.Frameworks...), ", ")
	}
	for k, v := range req.Variables {
		vars[k] = v
	}
	return vars
}

func (p *Platform) publish(t eventbus.EventType, msg string, payload interface{}) {
	if p.events != nil {
		p.events.Publish(eventbus.Event{Type: t, Message: msg, Payload: payload})
	}
}

func fallbackContent(spec DocumentSpec, vars map[string]string) string {
	return fmt.Sprintf("# %s\n\n## Overview\n\nGenerated documentation for %s.\n", spec.Name, vars["project_name"])
}

func outputPathForTemplate(t manifest.TemplateDefinition) string {
	base := strings.TrimSuffix(filepath.Base(t.Location), ".template.md")
	if base == "" || base == filepath.Base(t.Location) {
		base = t.Name
	}
	if strings.EqualFold(base, "README") {
		return "README.md"
	}
	return filepath.Join("docs", base+".md")
}

func docIDFromPath(path string) string {
	base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	return strings.ToLower(strings.ReplaceAll(base, " ", "-"))
}

func docNameFromPath(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}

func docTypeFromPath(path string) DocumentType {
	return docTypeFromName(docNameFromPath(path))
}

func docTypeFromName(name string) DocumentType {
	switch strings.ToLower(strings.ReplaceAll(name, " ", "")) {
	case "readme":
		return DocumentREADME
	case "prd":
		return DocumentPRD
	case "architecture":
		return DocumentArchitecture
	case "database":
		return DocumentDatabase
	case "api":
		return DocumentAPI
	case "deployment":
		return DocumentDeployment
	case "decisions":
		return DocumentDecisions
	case "businessrules":
		return DocumentBusinessRules
	case "troubleshooting":
		return DocumentTroubleshooting
	case "progress":
		return DocumentProgress
	default:
		return DocumentType(name)
	}
}
