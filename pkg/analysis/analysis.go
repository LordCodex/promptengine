package analysis

import (
	"github.com/LordCodex/promptengine/internal/domain/health"
	"github.com/LordCodex/promptengine/internal/domain/quality"
	"github.com/LordCodex/promptengine/internal/domain/quality/audit"
	"github.com/LordCodex/promptengine/internal/domain/quality/compliance"
	"github.com/LordCodex/promptengine/internal/domain/quality/doctor"
	"github.com/LordCodex/promptengine/internal/domain/quality/validation"
	"github.com/LordCodex/promptengine/internal/filesystem"
)

type AnalysisResult struct {
	Findings []quality.Finding
	Health   *health.Result
}

type Analyzer struct {
	docEngine  *doctor.DoctorEngine
	auditEng   *audit.AuditEngine
	compEngine *compliance.ComplianceEngine
	valReg     *validation.Registry
	healthReg  *health.Registry
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{
		docEngine:  doctor.NewDoctorEngine(),
		auditEng:   audit.NewAuditEngine(),
		compEngine: compliance.NewComplianceEngine(),
		valReg:     validation.NewRegistry(),
		healthReg:  health.NewRegistry(),
	}
}

func (a *Analyzer) Analyze(fs filesystem.FileSystem) (*AnalysisResult, error) {
	var findings []quality.Finding
	if r, err := a.docEngine.Diagnose(fs); err == nil {
		findings = append(findings, r.Findings...)
	}
	if r, err := a.auditEng.Run(fs); err == nil {
		findings = append(findings, r.Findings...)
	}
	if r, err := a.compEngine.Run(fs); err == nil {
		for _, pr := range r.ProfileResults {
			findings = append(findings, pr.Findings...)
		}
	}
	if f, err := a.valReg.Run(fs); err == nil {
		findings = append(findings, f...)
	}
	h, err := a.healthReg.Evaluate(fs)
	if err != nil {
		return nil, err
	}
	return &AnalysisResult{
		Findings: findings,
		Health:   h,
	}, nil
}
