// Package marketplace defines the architectural interface contracts for a future
// PromptEngine Marketplace. No concrete implementation is provided at this stage.
// All types here are intentionally interface-only to allow future implementation
// without changing the contracts that other packages depend on.
package marketplace

import "context"

// PackageKind mirrors the installer's kind taxonomy
type PackageKind string

const (
	KindPlugin         PackageKind = "plugin"
	KindTechnologyPack PackageKind = "technology-pack"
	KindOrgPack        PackageKind = "org-pack"
	KindTemplate       PackageKind = "template"
	KindWorkflowPack   PackageKind = "workflow-pack"
	KindPromptLibrary  PackageKind = "prompt-library"
	KindAIProvider     PackageKind = "ai-provider"
)

// PackageListing is a marketplace search result entry
type PackageListing struct {
	ID          string
	Name        string
	Kind        PackageKind
	Version     string
	Author      string
	Description string
	Downloads   int
	Verified    bool // future: signature-verified
}

// InstallRequest is submitted to the marketplace client to download and install
type InstallRequest struct {
	PackageID string
	Version   string // empty = latest
	Offline   bool
}

// SearchFilter narrows marketplace search results
type SearchFilter struct {
	Kind       PackageKind
	Query      string
	MaxResults int
}

// MarketplaceClient is the interface any future marketplace backend must satisfy.
// Concrete implementations (e.g., PromptEngine Cloud, self-hosted registries) will
// implement this without changing any caller code.
type MarketplaceClient interface {
	Search(ctx context.Context, filter SearchFilter) ([]PackageListing, error)
	GetPackage(ctx context.Context, id string) (*PackageListing, error)
	Download(ctx context.Context, req InstallRequest) ([]byte, error)
	Publish(ctx context.Context, manifest []byte) error
}
