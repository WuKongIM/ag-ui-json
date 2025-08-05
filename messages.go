package agui

import (
	"fmt"
)

// Message represents any type of message in the AG-UI system.
// This is implemented as an interface to support the union type from the JSON schema.
type Message interface {
	GetID() string
	GetRole() Role
	GetName() string
	Validate() error
	// MessageType returns the concrete type name for type switching
	MessageType() string
}

// BaseMessage contains common properties shared by all message types.
type BaseMessage struct {
	ID   string `json:"id"`             // Unique identifier for the message
	Role Role   `json:"role"`           // Role of the message sender
	Name string `json:"name,omitempty"` // Optional name of the sender
}

// GetID returns the message ID.
func (b *BaseMessage) GetID() string {
	return b.ID
}

// GetRole returns the message role.
func (b *BaseMessage) GetRole() Role {
	return b.Role
}

// GetName returns the message sender name.
func (b *BaseMessage) GetName() string {
	return b.Name
}

// Validate checks if the BaseMessage is valid.
func (b *BaseMessage) Validate() error {
	if b.ID == "" {
		return fmt.Errorf("message ID is required")
	}
	if !b.Role.IsValid() {
		return fmt.Errorf("invalid message role: %s", b.Role)
	}
	return nil
}

// DeveloperMessage represents a message from a developer.
type DeveloperMessage struct {
	BaseMessage
	Content string `json:"content"` // Text content of the message
}

// MessageType returns the concrete type name.
func (d *DeveloperMessage) MessageType() string {
	return "DeveloperMessage"
}

// Validate checks if the DeveloperMessage is valid.
func (d *DeveloperMessage) Validate() error {
	if err := d.BaseMessage.Validate(); err != nil {
		return err
	}
	if d.Role != RoleDeveloper {
		return fmt.Errorf("developer message must have developer role, got: %s", d.Role)
	}
	if d.Content == "" {
		return fmt.Errorf("developer message content is required")
	}
	return nil
}

// SystemMessage represents a system message.
type SystemMessage struct {
	BaseMessage
	Content string `json:"content"` // Text content of the message
}

// MessageType returns the concrete type name.
func (s *SystemMessage) MessageType() string {
	return "SystemMessage"
}

// Validate checks if the SystemMessage is valid.
func (s *SystemMessage) Validate() error {
	if err := s.BaseMessage.Validate(); err != nil {
		return err
	}
	if s.Role != RoleSystem {
		return fmt.Errorf("system message must have system role, got: %s", s.Role)
	}
	if s.Content == "" {
		return fmt.Errorf("system message content is required")
	}
	return nil
}

// AssistantMessage represents a message from an assistant.
type AssistantMessage struct {
	BaseMessage
	Content   string     `json:"content,omitempty"`   // Text content of the message
	ToolCalls []ToolCall `json:"toolCalls,omitempty"` // Tool calls made in this message
}

// MessageType returns the concrete type name.
func (a *AssistantMessage) MessageType() string {
	return "AssistantMessage"
}

// Validate checks if the AssistantMessage is valid.
func (a *AssistantMessage) Validate() error {
	if err := a.BaseMessage.Validate(); err != nil {
		return err
	}
	if a.Role != RoleAssistant {
		return fmt.Errorf("assistant message must have assistant role, got: %s", a.Role)
	}

	// Validate tool calls if present
	for i, toolCall := range a.ToolCalls {
		if err := toolCall.Validate(); err != nil {
			return fmt.Errorf("invalid tool call at index %d: %w", i, err)
		}
	}

	return nil
}

// UserMessage represents a message from a user.
type UserMessage struct {
	BaseMessage
	Content string `json:"content"` // Text content of the message
}

// MessageType returns the concrete type name.
func (u *UserMessage) MessageType() string {
	return "UserMessage"
}

// Validate checks if the UserMessage is valid.
func (u *UserMessage) Validate() error {
	if err := u.BaseMessage.Validate(); err != nil {
		return err
	}
	if u.Role != RoleUser {
		return fmt.Errorf("user message must have user role, got: %s", u.Role)
	}
	if u.Content == "" {
		return fmt.Errorf("user message content is required")
	}
	return nil
}

// ToolMessage represents a message from a tool.
type ToolMessage struct {
	BaseMessage
	Content    string `json:"content"`         // Text content of the message
	ToolCallID string `json:"toolCallId"`      // ID of the tool call this message responds to
	Error      string `json:"error,omitempty"` // Error message if the tool call failed
}

// MessageType returns the concrete type name.
func (t *ToolMessage) MessageType() string {
	return "ToolMessage"
}

// Validate checks if the ToolMessage is valid.
func (t *ToolMessage) Validate() error {
	if err := t.BaseMessage.Validate(); err != nil {
		return err
	}
	if t.Role != RoleTool {
		return fmt.Errorf("tool message must have tool role, got: %s", t.Role)
	}
	if t.Content == "" {
		return fmt.Errorf("tool message content is required")
	}
	if t.ToolCallID == "" {
		return fmt.Errorf("tool message toolCallId is required")
	}
	return nil
}

// MessageWrapper is used for JSON marshaling/unmarshaling of the Message interface.
type MessageWrapper struct {
	Role Role `json:"role"`
	*DeveloperMessage
	*SystemMessage
	*AssistantMessage
	*UserMessage
	*ToolMessage
}
