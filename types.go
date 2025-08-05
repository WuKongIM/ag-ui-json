// Package agui provides Go types and encoding/decoding functionality for the Agent User Interaction Protocol.
// This package implements the AG-UI protocol data structures based on the JSON schemas defined in the
// official AG-UI specification.
package agui

import (
	"encoding/json"
	"fmt"
)

// EventType represents all possible event types in the AG-UI protocol.
type EventType string

// Event type constants as defined in the AG-UI protocol specification.
const (
	EventTypeTextMessageStart   EventType = "TEXT_MESSAGE_START"
	EventTypeTextMessageContent EventType = "TEXT_MESSAGE_CONTENT"
	EventTypeTextMessageEnd     EventType = "TEXT_MESSAGE_END"
	EventTypeToolCallStart      EventType = "TOOL_CALL_START"
	EventTypeToolCallArgs       EventType = "TOOL_CALL_ARGS"
	EventTypeToolCallEnd        EventType = "TOOL_CALL_END"
	EventTypeToolCallResult     EventType = "TOOL_CALL_RESULT"
	EventTypeStateSnapshot      EventType = "STATE_SNAPSHOT"
	EventTypeStateDelta         EventType = "STATE_DELTA"
	EventTypeMessagesSnapshot   EventType = "MESSAGES_SNAPSHOT"
	EventTypeRaw                EventType = "RAW"
	EventTypeCustom             EventType = "CUSTOM"
	EventTypeRunStarted         EventType = "RUN_STARTED"
	EventTypeRunFinished        EventType = "RUN_FINISHED"
	EventTypeRunError           EventType = "RUN_ERROR"
	EventTypeStepStarted        EventType = "STEP_STARTED"
	EventTypeStepFinished       EventType = "STEP_FINISHED"
)

// IsValid checks if the EventType is a valid AG-UI event type.
func (e EventType) IsValid() bool {
	switch e {
	case EventTypeTextMessageStart, EventTypeTextMessageContent, EventTypeTextMessageEnd,
		EventTypeToolCallStart, EventTypeToolCallArgs, EventTypeToolCallEnd, EventTypeToolCallResult,
		EventTypeStateSnapshot, EventTypeStateDelta, EventTypeMessagesSnapshot,
		EventTypeRaw, EventTypeCustom,
		EventTypeRunStarted, EventTypeRunFinished, EventTypeRunError,
		EventTypeStepStarted, EventTypeStepFinished:
		return true
	default:
		return false
	}
}

// Role represents the possible roles for message senders in conversations.
type Role string

// Role constants as defined in the AG-UI protocol specification.
const (
	RoleDeveloper Role = "developer"
	RoleSystem    Role = "system"
	RoleAssistant Role = "assistant"
	RoleUser      Role = "user"
	RoleTool      Role = "tool"
)

// IsValid checks if the Role is a valid AG-UI role.
func (r Role) IsValid() bool {
	switch r {
	case RoleDeveloper, RoleSystem, RoleAssistant, RoleUser, RoleTool:
		return true
	default:
		return false
	}
}

// ToolCallType represents the type of tool call.
type ToolCallType string

// Tool call type constants as defined in the AG-UI protocol specification.
const (
	ToolCallTypeFunction ToolCallType = "function"
)

// IsValid checks if the ToolCallType is valid.
func (t ToolCallType) IsValid() bool {
	return t == ToolCallTypeFunction
}

// State represents the state of an agent during execution.
// It can hold any data structure needed by the agent implementation.
type State interface{}

// Context represents a piece of contextual information provided to an agent.
type Context struct {
	Description string `json:"description"` // Description of what this context represents
	Value       string `json:"value"`       // The actual context value
}

// Validate checks if the Context is valid according to AG-UI schema constraints.
func (c *Context) Validate() error {
	if c.Description == "" {
		return fmt.Errorf("context description is required")
	}
	if c.Value == "" {
		return fmt.Errorf("context value is required")
	}
	return nil
}

// Tool defines a tool that can be called by an agent.
type Tool struct {
	Name        string      `json:"name"`        // Name of the tool
	Description string      `json:"description"` // Description of what the tool does
	Parameters  interface{} `json:"parameters"`  // JSON Schema defining the parameters for the tool
}

// Validate checks if the Tool is valid according to AG-UI schema constraints.
func (t *Tool) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("tool name is required")
	}
	if t.Description == "" {
		return fmt.Errorf("tool description is required")
	}
	if t.Parameters == nil {
		return fmt.Errorf("tool parameters are required")
	}
	return nil
}

// FunctionCall represents function name and arguments in a tool call.
type FunctionCall struct {
	Name      string `json:"name"`      // Name of the function to call
	Arguments string `json:"arguments"` // JSON-encoded string of arguments to the function
}

// Validate checks if the FunctionCall is valid according to AG-UI schema constraints.
func (f *FunctionCall) Validate() error {
	if f.Name == "" {
		return fmt.Errorf("function name is required")
	}
	if f.Arguments == "" {
		return fmt.Errorf("function arguments are required")
	}
	// Validate that arguments is valid JSON
	var args interface{}
	if err := json.Unmarshal([]byte(f.Arguments), &args); err != nil {
		return fmt.Errorf("function arguments must be valid JSON: %w", err)
	}
	return nil
}

// ToolCall represents a tool call made by an agent.
type ToolCall struct {
	ID       string       `json:"id"`       // Unique identifier for the tool call
	Type     ToolCallType `json:"type"`     // Type of the tool call
	Function FunctionCall `json:"function"` // Details about the function being called
}

// Validate checks if the ToolCall is valid according to AG-UI schema constraints.
func (t *ToolCall) Validate() error {
	if t.ID == "" {
		return fmt.Errorf("tool call ID is required")
	}
	if !t.Type.IsValid() {
		return fmt.Errorf("invalid tool call type: %s", t.Type)
	}
	return t.Function.Validate()
}

// RunAgentInput represents input parameters for running an agent.
type RunAgentInput struct {
	ThreadID       string      `json:"threadId"`       // ID of the conversation thread
	RunID          string      `json:"runId"`          // ID of the current run
	State          State       `json:"state"`          // Current state of the agent
	Messages       []Message   `json:"messages"`       // List of messages in the conversation
	Tools          []Tool      `json:"tools"`          // List of tools available to the agent
	Context        []Context   `json:"context"`        // List of context objects provided to the agent
	ForwardedProps interface{} `json:"forwardedProps"` // Additional properties forwarded to the agent
}

// Validate checks if the RunAgentInput is valid according to AG-UI schema constraints.
func (r *RunAgentInput) Validate() error {
	if r.ThreadID == "" {
		return fmt.Errorf("thread ID is required")
	}
	if r.RunID == "" {
		return fmt.Errorf("run ID is required")
	}

	// Validate messages
	for i, msg := range r.Messages {
		if err := msg.Validate(); err != nil {
			return fmt.Errorf("invalid message at index %d: %w", i, err)
		}
	}

	// Validate tools
	for i, tool := range r.Tools {
		if err := tool.Validate(); err != nil {
			return fmt.Errorf("invalid tool at index %d: %w", i, err)
		}
	}

	// Validate context
	for i, ctx := range r.Context {
		if err := ctx.Validate(); err != nil {
			return fmt.Errorf("invalid context at index %d: %w", i, err)
		}
	}

	return nil
}
