package sync

import "fmt"

// ChangeSignal is a detected project change that may require documentation updates
type ChangeSignal string

const (
	SignalNewMigration       ChangeSignal = "new-migration"
	SignalNewAPI             ChangeSignal = "new-api"
	SignalNewService         ChangeSignal = "new-service"
	SignalNewQueue           ChangeSignal = "new-queue"
	SignalNewScheduledJob    ChangeSignal = "new-scheduled-job"
	SignalNewModule          ChangeSignal = "new-module"
	SignalNewPackage         ChangeSignal = "new-package"
	SignalNewTechnology      ChangeSignal = "new-technology"
	SignalArchitectureChange ChangeSignal = "architecture-change"
	SignalNewDocumentation   ChangeSignal = "new-documentation"
)

// SyncRule maps a ChangeSignal to the document IDs that should be updated
type SyncRule struct {
	Signal       ChangeSignal
	AffectedDocs []string // doc IDs (matches DocumentSpec.ID)
	Description  string
}

// SyncRecommendation is one actionable update suggestion
type SyncRecommendation struct {
	DocID     string
	Reason    string
	Signal    ChangeSignal
	AutoApply bool // true = safe to auto-update without approval
}

// SyncResult is the outcome of a sync operation
type SyncResult struct {
	DryRun          bool
	Recommendations []SyncRecommendation
	Applied         []string // doc IDs that were auto-applied
	Pending         []string // doc IDs awaiting approval
}

// ChangeDetector inspects signals and maps them to documentation
type ChangeDetector struct {
	rules []SyncRule
}

func NewChangeDetector() *ChangeDetector {
	return &ChangeDetector{rules: defaultRules()}
}

// RegisterRule allows plugins and workflows to add custom sync rules
func (d *ChangeDetector) RegisterRule(rule SyncRule) {
	d.rules = append(d.rules, rule)
}

// Detect returns all recommendations for a given set of change signals
func (d *ChangeDetector) Detect(signals []ChangeSignal) []SyncRecommendation {
	seen := make(map[string]bool)
	var recs []SyncRecommendation
	for _, signal := range signals {
		for _, rule := range d.rules {
			if rule.Signal == signal {
				for _, docID := range rule.AffectedDocs {
					key := fmt.Sprintf("%s:%s", signal, docID)
					if !seen[key] {
						seen[key] = true
						recs = append(recs, SyncRecommendation{
							DocID:     docID,
							Reason:    rule.Description,
							Signal:    signal,
							AutoApply: signal == SignalNewDocumentation, // only safe for doc additions
						})
					}
				}
			}
		}
	}
	return recs
}

// SyncEngine orchestrates documentation synchronization
type SyncEngine struct {
	detector *ChangeDetector
}

func NewSyncEngine(detector *ChangeDetector) *SyncEngine {
	return &SyncEngine{detector: detector}
}

// Run performs a full sync cycle over the given signals
func (e *SyncEngine) Run(signals []ChangeSignal, dryRun bool) SyncResult {
	recs := e.detector.Detect(signals)
	result := SyncResult{DryRun: dryRun, Recommendations: recs}
	for _, rec := range recs {
		if dryRun {
			result.Pending = append(result.Pending, rec.DocID)
		} else if rec.AutoApply {
			result.Applied = append(result.Applied, rec.DocID)
		} else {
			result.Pending = append(result.Pending, rec.DocID)
		}
	}
	return result
}

// defaultRules are the built-in sync rules encoding documentation dependencies
func defaultRules() []SyncRule {
	return []SyncRule{
		{
			Signal:       SignalNewMigration,
			AffectedDocs: []string{"database", "architecture", "api", "deployment", "progress", "roadmap"},
			Description:  "A new database migration requires updating the Database, Architecture, API, Deployment, Progress and Roadmap documents.",
		},
		{
			Signal:       SignalArchitectureChange,
			AffectedDocs: []string{"architecture", "decisions", "business-rules", "deployment"},
			Description:  "An architecture change requires updating Architecture, Decisions, Business Rules and Deployment documents.",
		},
		{
			Signal:       SignalNewAPI,
			AffectedDocs: []string{"api", "architecture", "deployment", "progress"},
			Description:  "A new API requires updating the API, Architecture, Deployment and Progress documents.",
		},
		{
			Signal:       SignalNewService,
			AffectedDocs: []string{"architecture", "deployment", "database", "api"},
			Description:  "A new service requires updating Architecture, Deployment, Database and API documents.",
		},
		{
			Signal:       SignalNewTechnology,
			AffectedDocs: []string{"architecture", "deployment", "decisions"},
			Description:  "A new technology addition requires updating Architecture, Deployment and Decisions documents.",
		},
		{
			Signal:       SignalNewDocumentation,
			AffectedDocs: []string{"progress"},
			Description:  "New documentation should be reflected in Progress.",
		},
	}
}
