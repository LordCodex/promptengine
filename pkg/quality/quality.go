package quality

import (
	"github.com/LordCodex/promptengine/internal/domain/quality"
	"github.com/LordCodex/promptengine/internal/domain/quality/doctor"
	"github.com/LordCodex/promptengine/internal/domain/quality/validation"
	"github.com/LordCodex/promptengine/internal/filesystem"
)

type Finding = quality.Finding
type Severity = quality.Severity

type Checker struct {
	docEngine *doctor.DoctorEngine
	valReg    *validation.Registry
}

func NewChecker() *Checker {
	return &Checker{
		docEngine: doctor.NewDoctorEngine(),
		valReg:    validation.NewRegistry(),
	}
}

func (c *Checker) RunDoctor(fs filesystem.FileSystem) (*quality.Report, error) {
	return c.docEngine.Diagnose(fs)
}

func (c *Checker) RunValidation(fs filesystem.FileSystem) ([]Finding, error) {
	return c.valReg.Run(fs)
}
