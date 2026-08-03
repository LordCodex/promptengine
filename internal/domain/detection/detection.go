package detection

import (
	"github.com/LordCodex/promptengine/internal/filesystem"
)

// ProjectMetadata details stack parameters on disk
type ProjectMetadata struct {
	Languages     []string
	Frameworks    []string
	Databases     []string
	Queues        []string
	Caches        []string
	TestingFrames []string
	CIs           []string
	HasDocker     bool
	IsMonorepo    bool
	IsMobile      bool
	HasAgentsMD   bool
	HasDocs       bool
	HasConfig     bool
}

// Detector is a modular scanning unit
type Detector interface {
	Name() string
	Detect(fs filesystem.FileSystem) (bool, error)
	Apply(fs filesystem.FileSystem, meta *ProjectMetadata) error
}

// Registry manages project stack detectors list
type Registry struct {
	detectors []Detector
}

func NewRegistry() *Registry {
	return &Registry{
		detectors: make([]Detector, 0),
	}
}

func (r *Registry) Register(d Detector) {
	r.detectors = append(r.detectors, d)
}

func (r *Registry) DetectAll(fs filesystem.FileSystem) (*ProjectMetadata, error) {
	meta := &ProjectMetadata{
		Languages:     make([]string, 0),
		Frameworks:    make([]string, 0),
		Databases:     make([]string, 0),
		Queues:        make([]string, 0),
		Caches:        make([]string, 0),
		TestingFrames: make([]string, 0),
		CIs:           make([]string, 0),
	}

	for _, d := range r.detectors {
		ok, err := d.Detect(fs)
		if err != nil {
			return nil, err
		}
		if ok {
			if err := d.Apply(fs, meta); err != nil {
				return nil, err
			}
		}
	}

	return meta, nil
}
