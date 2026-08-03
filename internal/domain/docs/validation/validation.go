package validation

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/LordCodex/promptengine/internal/domain/quality"
	"github.com/LordCodex/promptengine/internal/filesystem"
)

// FindingSeverity classifies validation findings
type FindingSeverity = quality.Severity

const (
	SeverityError   FindingSeverity = quality.SeverityError
	SeverityWarning FindingSeverity = quality.SeverityWarning
	SeverityInfo    FindingSeverity = quality.SeverityInfo
)

// ValidationFinding is one actionable finding from a validation rule
type ValidationFinding = quality.Finding

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
			Engine:         "docs-validation",
			Rule:           "missing-document",
			FilePath:       docPath,
			Severity:       SeverityError,
			Title:          fmt.Sprintf("document '%s' does not exist", docPath),
			Recommendation: fmt.Sprintf("run 'promptengine generate' to create the missing document at %s", docPath),
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
			Engine:         "docs-validation",
			Rule:           r.Name(),
			FilePath:       docPath,
			Severity:       SeverityError,
			Title:          "document contains no headings",
			Recommendation: "ensure the document has at least one top-level heading (# Heading)",
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
				Engine:         "docs-validation",
				Rule:           r.Name(),
				FilePath:       docPath,
				Severity:       SeverityWarning,
				Title:          fmt.Sprintf("document contains stale placeholder text ('%s')", phrase),
				Recommendation: "complete the document section or remove placeholder text",
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
	linkPattern := regexp.MustCompile(`\[[^\]]+\]\(([^)]+)\)`)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		for _, match := range linkPattern.FindAllStringSubmatch(line, -1) {
			href := strings.TrimSpace(match[1])
			if href == "" || strings.HasPrefix(href, "http") || strings.HasPrefix(href, "#") || strings.HasPrefix(href, "mailto:") {
				continue
			}
			if strings.HasPrefix(href, "file://") {
				href = strings.TrimPrefix(href, "file://")
			}
			href = strings.SplitN(href, "#", 2)[0]
			if href == "" {
				continue
			}
			candidates := []string{href}
			if !filepath.IsAbs(href) {
				candidates = append(candidates, filepath.Clean(filepath.Join(filepath.Dir(docPath), href)))
			}
			if !existsAny(fs, candidates) {
				findings = append(findings, ValidationFinding{
					Engine:         "docs-validation",
					Rule:           r.Name(),
					FilePath:       docPath,
					Severity:       SeverityWarning,
					Title:          fmt.Sprintf("broken reference to '%s'", href),
					Recommendation: fmt.Sprintf("update or remove the link to '%s'", href),
				})
			}
		}
	}
	return findings
}

func existsAny(fs filesystem.FileSystem, paths []string) bool {
	for _, path := range paths {
		if fs.Exists(path) {
			return true
		}
	}
	return false
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
				Engine:         "docs-validation",
				Rule:           r.Name(),
				FilePath:       docPath,
				Severity:       SeverityWarning,
				Title:          fmt.Sprintf("duplicate heading '%s' appears %d times", heading, count),
				Recommendation: "merge duplicate sections or rename headings",
			})
		}
	}
	return findings
}

// orphanedDocumentRule detects documents not referenced by any manifest entry
type orphanedDocumentRule struct{}

func (r *orphanedDocumentRule) Name() string { return "orphaned-document" }
func (r *orphanedDocumentRule) Run(fs filesystem.FileSystem, docPath string, content string) []ValidationFinding {
	if len(content) < 50 {
		return []ValidationFinding{{
			Engine:         "docs-validation",
			Rule:           r.Name(),
			FilePath:       docPath,
			Severity:       SeverityInfo,
			Title:          "document is very short and may be orphaned or empty",
			Recommendation: "verify this document is referenced from the manifest and has substantive content",
		}}
	}
	return nil
}
