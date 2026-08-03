package ai

type Capability string

const (
	CapInteractiveChat   Capability = "interactive_chat"
	CapStreaming         Capability = "streaming"
	CapToolCalling       Capability = "tool_calling"
	CapStructuredOutput  Capability = "structured_output"
	CapJSONMode          Capability = "json_mode"
	CapImageInput        Capability = "image_input"
	CapImageGeneration   Capability = "image_generation"
	CapVision            Capability = "vision"
	CapCodeExecution     Capability = "code_execution"
	CapLongContext       Capability = "long_context"
	CapReasoningControls Capability = "reasoning_controls"
	CapMCPSupport        Capability = "mcp_support"
)

// ProviderMetadata specifies operational and billing properties
type ProviderMetadata struct {
	ID          string
	Name        string
	ContextSize int  // token limit
	CostRank    int  // 1 (low) to 5 (high)
	SpeedRank   int  // 1 (slow) to 5 (fast)
	IsOffline   bool // supports air-gapped run
}
