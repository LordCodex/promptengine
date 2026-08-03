package ai

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type Message struct {
	Role    Role
	Content string
}

// Conversation tracks state history for multi-turn workflows
type Conversation struct {
	Messages []Message
}

func NewConversation(systemPrompt string) *Conversation {
	c := &Conversation{
		Messages: make([]Message, 0),
	}
	if systemPrompt != "" {
		c.Append(RoleSystem, systemPrompt)
	}
	return c
}

func (c *Conversation) Append(role Role, text string) {
	c.Messages = append(c.Messages, Message{Role: role, Content: text})
}

func (c *Conversation) Reset(systemPrompt string) {
	c.Messages = make([]Message, 0)
	if systemPrompt != "" {
		c.Append(RoleSystem, systemPrompt)
	}
}

// EstimateTokens calculates a basic estimation for budgeting inputs (1 token ~ 4 characters)
func (c *Conversation) EstimateTokens() int {
	charCount := 0
	for _, m := range c.Messages {
		charCount += len(m.Content)
	}
	return charCount / 4
}
