package docs

import (
	"fmt"

	"github.com/LordCodex/promptengine/internal/domain/docs/generator"
)

type DocType = generator.DocType
type GeneratorInput = generator.GeneratorInput
type GeneratorOutput = generator.GeneratorOutput

type Engine struct {
	registry *generator.GeneratorRegistry
}

func NewEngine() *Engine {
	reg := generator.NewGeneratorRegistry()
	generator.RegisterDefaults(reg)
	return &Engine{registry: reg}
}

func (e *Engine) Generate(docType DocType, input GeneratorInput) (GeneratorOutput, error) {
	g, ok := e.registry.Get(docType)
	if !ok {
		return GeneratorOutput{}, fmt.Errorf("no generator for type %s", docType)
	}
	return g.Generate(input)
}
