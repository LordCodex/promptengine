package validation

import (
	"fmt"
	"strings"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

// FindingSeverity classifies validation findings
type FindingSeverity string

const (
	SeverityError   FindingSeverity = "error"
	SeverityWarning FindingSeverity = "warning"
	SeverityInfo    FindingSeverity = "info"
)

// ValidationFinding is one actionable finding from a validation rule
type ValidationFinding struct {
	Rule        string
	DocPath     string
	Severity    FindingSeverity
	Message     string
	Suggestion  string
}

// ValidationRule is the interface all validation rules must implement
type ValidationRule interface {
	Name() string
	Run(fs filesystem.FileSystem, docPath string, content string) []ValidationFinding
}

// Validator applies all registered rules to a document
type Validator struct {
	rules []ValidationRule
}

func NewValidator() *Validator {
	v := &Validator{}
	v.RegisterDefaults()
	return v
}

func (v *Validator) Register(rule ValidationRule) {
	v.rules = append(v.rules, rule)
}

func (v *Validator) RegisterDefaults() {
	v.rules = []ValidationRule{
		&missingSectionsRule{},
		&staleDocumentRule{},
		&brokenReferencesRule{},
		&duplicateContentRule{},
		&orphanedDocumentRule{},
	}
}

// Validate runs all rules and returns aggregated findings
func (v *Validator) Validate(fs filesystem.FileSystem, docPath string) ([]ValidationFinding, error) {
	var findings []ValidationFinding

	if !fs.Exists(docPath) {
		return []ValidationFinding{{
			Rule:       "missing-document",
			DocPath:    docPath,
			Severity:   SeverityError,
			Message:    fmt.Sprintf("document '%s' does not exist", docPath),
			Suggestion: fmt.Sprintf("run 'promptengine generate' to create the missing document at %s", docPath),
		}}, nil
	}

	data, err := fs.ReadFile(docPath)
	if err != nil {
		return nil, err
	}
	content := string(data)

	for _, rule := range v.rules {
		f := rule.Run(fs, docPath, content)
		findings = append(findings, f...)
	}
	return findings, nil
}

// --- Concrete validation rules ---

// missingSectionsRule detects documents lacking any h1 or h2 headings
type missingSectionsRule struct{}

func (r *missingSectionsRule) Name() string { return "missing-sections" }
func (r *missingSectionsRule) Run(_ filesystem.FileSystem, docPath string, content string) []ValidationFinding {
	if !strings.Contains(content, "# ") {
		return []ValidationFinding{{
			Rule:       r.Name(),
			DocPath:    docPath,
			Severity:   SeverityError,
			Message:    "document contains no headings",
			Suggestion: "ensure the document has at least one top-level heading (# Heading)",
		}}
	}
	return nil
}

// staleDocumentRule flags documents with placeholder text indicating they were never filled in
type staleDocumentRule struct{}

func (r *staleDocumentRule) Name() string { return "stale-document" }
func (r *staleDocumentRule) Run(_ filesystem.FileSystem, docPath string, content string) []ValidationFinding {
	stalePhrases := []string{
		"_Define ", "_Document ", "_Track ", "TBD", "TODO", "FIXME",
	}
	for _, phrase := range stalePhrases {
		if strings.Contains(content, phrase) {
			return []ValidationFinding{{
				Rule:       r.Name(),
				DocPath:    docPath,
				Severity:   SeverityWarning,
				Message:    fmt.Sprintf("document contains stale placeholder text ('%s')", phrase),
				Suggestion: "complete the document section or remove placeholder text",
			}}
		}
	}
	return nil
}

// brokenReferencesRule detects markdown links pointing to non-existent files
type brokenReferencesRule struct{}

func (r *brokenReferencesRule) Name() string { return "broken-references" }
func (r *brokenReferencesRule) Run(fs filesystem.FileSystem, docPath string, content string) []ValidationFinding {
	var findings []ValidationFinding
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		// Simple markdown link detection: [text](path)
		for _, part := range strings.Split(line, "](") {
			if !strings.Contains(part, ")") {
				continue
			}
			href := strings.SplitN(part, ")", 2)[0]
			if strings.HasPrefix(href, "http") || strings.HasPrefix(href, "#") {
				continue
			}
			if !fs.Exists(href) {
				findings = append(findings, ValidationFinding{
					Rule:       r.Name(),
					DocPath:    docPath,
					Severity:   SeverityWarning,
					Message:    fmt.Sprintf("broken reference to '%s'", href),
					Suggestion: fmt.Sprintf("update or remove the link to '%s'", href),
				})
			}
		}
	}
	return findings
}

// duplicateContentRule flags duplicate top-level headings within a document
type duplicateContentRule struct{}

func (r *duplicateContentRule) Name() string { return "duplicate-content" }
func (r *duplicateContentRule) Run(_ filesystem.FileSystem, docPath string, content string) []ValidationFinding {
	headings := make(map[string]int)
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, "# ") {
			headings[line]++
		}
	}
	var findings []ValidationFinding
	for heading, count := range headings {
		if count > 1 {
			findings = append(findings, ValidationFinding{
				Rule:       r.Name(),
				DocPath:    docPath,
				Severity:   SeverityWarning,
				Message:    fmt.Sprintf("duplicate heading '%s' appears %d times", heading, count),
				Suggestion: "merge duplicate sections or rename headings",
			})
		}
	}
	return findings
}

// orphanedDocumentRule detects documents not referenced by any manifest entry
// (In the current implementation this emits an advisory — future versions will cross-reference the manifest)
type orphanedDocumentRule struct{}

func (r *orphanedDocumentRule) Name() string { return "orphaned-document" }
func (r *orphanedDocumentRule) Run(fs filesystem.FileSystem, docPath string, content string) []ValidationFinding {
	// Advisory only — checking against the manifest requires a registry reference
	// which is wired in by the calling engine. This rule emits info-level.
	if len(content) < 50 {
		return []ValidationFinding{{
			Rule:       r.Name(),
			DocPath:    docPath,
			Severity:   SeverityInfo,
			Message:    "document is very short and may be orphaned or empty",
			Suggestion: "verify this document is referenced from the manifest and has substantive content",
		}}
	}
	return nil
}
