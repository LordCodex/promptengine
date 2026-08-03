package generator

import (
	"fmt"
	"strings"
	"sync"
)

// DocType identifies a standard PromptEngine document
type DocType string

const (
	DocPRD             DocType = "prd"
	DocArchitecture    DocType = "architecture"
	DocBusinessRules   DocType = "business-rules"
	DocDatabase        DocType = "database"
	DocAPI             DocType = "api"
	DocProgress        DocType = "progress"
	DocRoadmap         DocType = "roadmap"
	DocDeployment      DocType = "deployment"
	DocTroubleshooting DocType = "troubleshooting"
	DocDecisions       DocType = "decisions"
	DocEnvironment     DocType = "environment"
	DocSecurity        DocType = "security"
	DocTesting         DocType = "testing"
	DocObservability   DocType = "observability"
	DocCICD            DocType = "ci-cd"
)

// GeneratorInput carries project metadata for generators
type GeneratorInput struct {
	ProjectName string
	Stack       []string
	Variables   map[string]string
}

// GeneratorOutput is the result of a document generation
type GeneratorOutput struct {
	DocType  DocType
	Filename string
	Content  string
}

// Generator defines the interface any document generator must implement
type Generator interface {
	DocType() DocType
	Generate(input GeneratorInput) (GeneratorOutput, error)
}

// GeneratorRegistry is the plugin-extensible catalogue of generators
type GeneratorRegistry struct {
	mu         sync.RWMutex
	generators map[DocType]Generator
}

func NewGeneratorRegistry() *GeneratorRegistry {
	return &GeneratorRegistry{generators: make(map[DocType]Generator)}
}

func (r *GeneratorRegistry) Register(g Generator) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.generators[g.DocType()]; exists {
		return fmt.Errorf("generator for doc type '%s' already registered", g.DocType())
	}
	r.generators[g.DocType()] = g
	return nil
}

func (r *GeneratorRegistry) Get(dt DocType) (Generator, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	g, ok := r.generators[dt]
	return g, ok
}

// All returns every registered generator in deterministic order
func (r *GeneratorRegistry) All() []Generator {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Generator, 0, len(r.generators))
	for _, g := range r.generators {
		out = append(out, g)
	}
	return out
}

// RegisterDefaults loads all standard PromptEngine generators
func RegisterDefaults(reg *GeneratorRegistry) {
	defaults := []Generator{
		&architectureGenerator{},
		&businessRulesGenerator{},
		&databaseGenerator{},
		&apiGenerator{},
		&prdGenerator{},
		&progressGenerator{},
		&roadmapGenerator{},
		&deploymentGenerator{},
		&troubleshootingGenerator{},
		&decisionsGenerator{},
		&securityGenerator{},
		&testingGenerator{},
	}
	for _, g := range defaults {
		_ = reg.Register(g) // ignore duplicate errors when called multiple times
	}
}

// --- Standard generators (minimal scaffold output, AI fills the details) ---

type architectureGenerator struct{}

func (g *architectureGenerator) DocType() DocType { return DocArchitecture }
func (g *architectureGenerator) Generate(in GeneratorInput) (GeneratorOutput, error) {
	content := fmt.Sprintf("# Architecture\n\n## Overview\n\n%s is built using %s.\n\n## Components\n\n_Define components here._\n\n## Decisions\n\n_See Decisions.md._\n",
		in.ProjectName, strings.Join(in.Stack, ", "))
	return GeneratorOutput{DocType: DocArchitecture, Filename: "docs/Architecture.md", Content: content}, nil
}

type businessRulesGenerator struct{}

func (g *businessRulesGenerator) DocType() DocType { return DocBusinessRules }
func (g *businessRulesGenerator) Generate(in GeneratorInput) (GeneratorOutput, error) {
	content := fmt.Sprintf("# Business Rules\n\n## Core Rules\n\n_Document business invariants for %s here._\n\n## Exceptions\n\n_Document approved exceptions with rationale._\n", in.ProjectName)
	return GeneratorOutput{DocType: DocBusinessRules, Filename: "docs/BusinessRules.md", Content: content}, nil
}

type databaseGenerator struct{}

func (g *databaseGenerator) DocType() DocType { return DocDatabase }
func (g *databaseGenerator) Generate(in GeneratorInput) (GeneratorOutput, error) {
	content := "# Database\n\n## Schema Overview\n\n_Document tables, relationships and constraints._\n\n## Migrations\n\n_Track all schema migrations here._\n\n## Indexes\n\n_Document indexes and rationale._\n"
	return GeneratorOutput{DocType: DocDatabase, Filename: "docs/Database.md", Content: content}, nil
}

type apiGenerator struct{}

func (g *apiGenerator) DocType() DocType { return DocAPI }
func (g *apiGenerator) Generate(in GeneratorInput) (GeneratorOutput, error) {
	content := "# API\n\n## Endpoints\n\n_Document all API endpoints, request/response schemas._\n\n## Authentication\n\n_Document auth mechanisms._\n\n## Versioning\n\n_Document API versioning strategy._\n"
	return GeneratorOutput{DocType: DocAPI, Filename: "docs/API.md", Content: content}, nil
}

type prdGenerator struct{}

func (g *prdGenerator) DocType() DocType { return DocPRD }
func (g *prdGenerator) Generate(in GeneratorInput) (GeneratorOutput, error) {
	content := fmt.Sprintf("# Product Requirements Document\n\n## Project\n\n%s\n\n## Problem Statement\n\n_Define the problem._\n\n## Goals\n\n_Define measurable goals._\n\n## Non-Goals\n\n_Explicitly out of scope._\n", in.ProjectName)
	return GeneratorOutput{DocType: DocPRD, Filename: "docs/PRD.md", Content: content}, nil
}

type progressGenerator struct{}

func (g *progressGenerator) DocType() DocType { return DocProgress }
func (g *progressGenerator) Generate(in GeneratorInput) (GeneratorOutput, error) {
	content := "# Progress\n\n## Completed\n\n_Track completed milestones._\n\n## In Progress\n\n_Current work._\n\n## Blocked\n\n_Blockers and resolutions._\n"
	return GeneratorOutput{DocType: DocProgress, Filename: "docs/Progress.md", Content: content}, nil
}

type roadmapGenerator struct{}

func (g *roadmapGenerator) DocType() DocType { return DocRoadmap }
func (g *roadmapGenerator) Generate(in GeneratorInput) (GeneratorOutput, error) {
	content := "# Roadmap\n\n## Near Term\n\n_Next 4 weeks._\n\n## Mid Term\n\n_Next quarter._\n\n## Long Term\n\n_Strategic horizon._\n"
	return GeneratorOutput{DocType: DocRoadmap, Filename: "docs/Roadmap.md", Content: content}, nil
}

type deploymentGenerator struct{}

func (g *deploymentGenerator) DocType() DocType { return DocDeployment }
func (g *deploymentGenerator) Generate(in GeneratorInput) (GeneratorOutput, error) {
	content := "# Deployment\n\n## Environments\n\n_Define environments (dev, staging, prod)._\n\n## Release Process\n\n_Step-by-step deployment procedure._\n\n## Rollback\n\n_Rollback strategy._\n"
	return GeneratorOutput{DocType: DocDeployment, Filename: "docs/Deployment.md", Content: content}, nil
}

type troubleshootingGenerator struct{}

func (g *troubleshootingGenerator) DocType() DocType { return DocTroubleshooting }
func (g *troubleshootingGenerator) Generate(in GeneratorInput) (GeneratorOutput, error) {
	content := "# Troubleshooting\n\n## Common Issues\n\n_Document known issues and resolutions._\n\n## Debug Procedures\n\n_Step-by-step diagnostic procedures._\n"
	return GeneratorOutput{DocType: DocTroubleshooting, Filename: "docs/Troubleshooting.md", Content: content}, nil
}

type decisionsGenerator struct{}

func (g *decisionsGenerator) DocType() DocType { return DocDecisions }
func (g *decisionsGenerator) Generate(in GeneratorInput) (GeneratorOutput, error) {
	content := "# Decisions\n\n## Architecture Decisions\n\n_Record all significant architecture decisions with rationale._\n\n## Rejected Alternatives\n\n_Document alternatives considered and why they were rejected._\n"
	return GeneratorOutput{DocType: DocDecisions, Filename: "docs/Decisions.md", Content: content}, nil
}

type securityGenerator struct{}

func (g *securityGenerator) DocType() DocType { return DocSecurity }
func (g *securityGenerator) Generate(in GeneratorInput) (GeneratorOutput, error) {
	content := "# Security\n\n## Threat Model\n\n_Define attack surface and mitigations._\n\n## Authentication & Authorization\n\n_Document auth mechanisms._\n\n## Data Protection\n\n_Encryption, PII handling._\n"
	return GeneratorOutput{DocType: DocSecurity, Filename: "docs/Security.md", Content: content}, nil
}

type testingGenerator struct{}

func (g *testingGenerator) DocType() DocType { return DocTesting }
func (g *testingGenerator) Generate(in GeneratorInput) (GeneratorOutput, error) {
	content := "# Testing\n\n## Strategy\n\n_Unit, integration, e2e approach._\n\n## Coverage Requirements\n\n_Minimum coverage targets._\n\n## Test Data\n\n_Fixtures and factories._\n"
	return GeneratorOutput{DocType: DocTesting, Filename: "docs/Testing.md", Content: content}, nil
}
