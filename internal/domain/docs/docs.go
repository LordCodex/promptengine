package docs

import (
	"fmt"
	"time"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

// SpecStatus remains the original status taxonomy
type SpecStatus string

const (
	StatusSynced  SpecStatus = "synced"
	StatusDrift   SpecStatus = "drift"
	StatusMissing SpecStatus = "missing"
	StatusStale   SpecStatus = "stale"
)

// RequiredSection declares a section every document of this type must contain
type RequiredSection struct {
	Heading     string
	Description string
}

// DocumentSpec describes a document type registered in the platform
type DocumentSpec struct {
	ID               string
	Name             string
	Purpose          string
	DefaultPath      string
	Priority         int // higher = more critical to health score
	RequiredSections []RequiredSection
	DependsOn        []string // IDs of other docs this doc depends on
	UpdateTriggers   []string // change signals that should trigger an update
}

// SpecMetadata represents the runtime state of a concrete document file
type SpecMetadata struct {
	Name        string
	Path        string
	Status      SpecStatus
	LastEdited  string
	GeneratedBy string
	Version     string
}

// VersionedDoc records a generation event for audit and rollback purposes
type VersionedDoc struct {
	DocID       string
	Version     string
	GeneratedAt time.Time
	Content     []byte
}

// DocRegistry is the central catalogue of all known document types
type DocRegistry struct {
	specs map[string]DocumentSpec
}

func NewDocRegistry() *DocRegistry {
	return &DocRegistry{specs: make(map[string]DocumentSpec)}
}

func (r *DocRegistry) Register(spec DocumentSpec) {
	r.specs[spec.ID] = spec
}

func (r *DocRegistry) Get(id string) (DocumentSpec, bool) {
	s, ok := r.specs[id]
	return s, ok
}

func (r *DocRegistry) All() []DocumentSpec {
	var list []DocumentSpec
	for _, s := range r.specs {
		list = append(list, s)
	}
	return list
}

// AffectedBy returns all document IDs that list the given trigger in their UpdateTriggers
func (r *DocRegistry) AffectedBy(trigger string) []string {
	var affected []string
	for id, spec := range r.specs {
		for _, t := range spec.UpdateTriggers {
			if t == trigger {
				affected = append(affected, id)
				break
			}
		}
	}
	return affected
}

// --- Legacy interfaces preserved for backward compatibility ---

// Syncer manages synchronizing specifications with code changes
type Syncer interface {
	Sync(fs filesystem.FileSystem, codeDiff string) (modifiedPaths []string, err error)
}

// Manager orchestrates documentation state checks
type Manager interface {
	Generate(fs filesystem.FileSystem, templateName string) (string, error)
	Validate(fs filesystem.FileSystem, docPath string) (bool, []string, error)
	Status(fs filesystem.FileSystem) ([]SpecMetadata, error)
}

// Engine wires together the doc registry with filesystem operations
type Engine struct {
	registry *DocRegistry
	fs       filesystem.FileSystem
	history  []*VersionedDoc
}

func NewEngine(fs filesystem.FileSystem, reg *DocRegistry) *Engine {
	return &Engine{fs: fs, registry: reg, history: make([]*VersionedDoc, 0)}
}

func (e *Engine) RecordGeneration(docID, version string, content []byte) {
	e.history = append(e.history, &VersionedDoc{
		DocID:       docID,
		Version:     version,
		GeneratedAt: time.Now(),
		Content:     content,
	})
}

func (e *Engine) Status() ([]SpecMetadata, error) {
	var result []SpecMetadata
	for _, spec := range e.registry.All() {
		status := StatusMissing
		if e.fs.Exists(spec.DefaultPath) {
			status = StatusSynced
		}
		result = append(result, SpecMetadata{
			Name:   spec.Name,
			Path:   spec.DefaultPath,
			Status: status,
		})
	}
	return result, nil
}

func (e *Engine) Validate(docID string) ([]string, error) {
	spec, ok := e.registry.Get(docID)
	if !ok {
		return nil, fmt.Errorf("document type '%s' not registered", docID)
	}
	if !e.fs.Exists(spec.DefaultPath) {
		return []string{fmt.Sprintf("document '%s' is missing at %s", spec.Name, spec.DefaultPath)}, nil
	}
	data, err := e.fs.ReadFile(spec.DefaultPath)
	if err != nil {
		return nil, err
	}
	content := string(data)

	var issues []string
	for _, section := range spec.RequiredSections {
		if !containsHeading(content, section.Heading) {
			issues = append(issues, fmt.Sprintf("missing required section '# %s' in %s", section.Heading, spec.Name))
		}
	}
	return issues, nil
}

// StandardDocManager is kept for backward compatibility
type StandardDocManager struct{}

func (m *StandardDocManager) Generate(fs filesystem.FileSystem, name string) (string, error) {
	return "", nil
}
func (m *StandardDocManager) Validate(fs filesystem.FileSystem, path string) (bool, []string, error) {
	return true, nil, nil
}
func (m *StandardDocManager) Status(fs filesystem.FileSystem) ([]SpecMetadata, error) {
	return make([]SpecMetadata, 0), nil
}

func containsHeading(content, heading string) bool {
	needle := "# " + heading
	for i := 0; i <= len(content)-len(needle); i++ {
		if content[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}
