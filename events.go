package agui

import (
	"fmt"
	"time"
)

// Event represents any event in the AG-UI system.
// This is implemented as an interface to support the union type from the JSON schema.
type Event interface {
	GetType() EventType
	GetTimestamp() *int64
	GetRawEvent() interface{}
	Validate() error
	// EventType returns the concrete type name for type switching
	EventTypeName() string
}

// BaseEvent contains common properties shared by all event types.
type BaseEvent struct {
	Type      EventType   `json:"type"`                // The type of event
	Timestamp *int64      `json:"timestamp,omitempty"` // Timestamp when the event was created
	RawEvent  interface{} `json:"rawEvent,omitempty"`  // Original event data if this event was transformed
}

// GetType returns the event type.
func (b *BaseEvent) GetType() EventType {
	return b.Type
}

// GetTimestamp returns the event timestamp.
func (b *BaseEvent) GetTimestamp() *int64 {
	return b.Timestamp
}

// GetRawEvent returns the raw event data.
func (b *BaseEvent) GetRawEvent() interface{} {
	return b.RawEvent
}

// Validate checks if the BaseEvent is valid.
func (b *BaseEvent) Validate() error {
	if !b.Type.IsValid() {
		return fmt.Errorf("invalid event type: %s", b.Type)
	}
	return nil
}

// SetTimestamp sets the timestamp to the current time.
func (b *BaseEvent) SetTimestamp() {
	now := time.Now().UnixMilli()
	b.Timestamp = &now
}

// Lifecycle Events

// RunStartedEvent signals the start of an agent run.
type RunStartedEvent struct {
	BaseEvent
	ThreadID string `json:"threadId"` // ID of the conversation thread
	RunID    string `json:"runId"`    // ID of the agent run
}

// EventTypeName returns the concrete type name.
func (r *RunStartedEvent) EventTypeName() string {
	return "RunStartedEvent"
}

// Validate checks if the RunStartedEvent is valid.
func (r *RunStartedEvent) Validate() error {
	if err := r.BaseEvent.Validate(); err != nil {
		return err
	}
	if r.Type != EventTypeRunStarted {
		return fmt.Errorf("run started event must have RUN_STARTED type, got: %s", r.Type)
	}
	if r.ThreadID == "" {
		return fmt.Errorf("thread ID is required")
	}
	if r.RunID == "" {
		return fmt.Errorf("run ID is required")
	}
	return nil
}

// RunFinishedEvent signals the successful completion of an agent run.
type RunFinishedEvent struct {
	BaseEvent
	ThreadID string      `json:"threadId"`         // ID of the conversation thread
	RunID    string      `json:"runId"`            // ID of the agent run
	Result   interface{} `json:"result,omitempty"` // Result data from the agent run
}

// EventTypeName returns the concrete type name.
func (r *RunFinishedEvent) EventTypeName() string {
	return "RunFinishedEvent"
}

// Validate checks if the RunFinishedEvent is valid.
func (r *RunFinishedEvent) Validate() error {
	if err := r.BaseEvent.Validate(); err != nil {
		return err
	}
	if r.Type != EventTypeRunFinished {
		return fmt.Errorf("run finished event must have RUN_FINISHED type, got: %s", r.Type)
	}
	if r.ThreadID == "" {
		return fmt.Errorf("thread ID is required")
	}
	if r.RunID == "" {
		return fmt.Errorf("run ID is required")
	}
	return nil
}

// RunErrorEvent signals an error during an agent run.
type RunErrorEvent struct {
	BaseEvent
	Message string `json:"message"`        // Error message
	Code    string `json:"code,omitempty"` // Error code
}

// EventTypeName returns the concrete type name.
func (r *RunErrorEvent) EventTypeName() string {
	return "RunErrorEvent"
}

// Validate checks if the RunErrorEvent is valid.
func (r *RunErrorEvent) Validate() error {
	if err := r.BaseEvent.Validate(); err != nil {
		return err
	}
	if r.Type != EventTypeRunError {
		return fmt.Errorf("run error event must have RUN_ERROR type, got: %s", r.Type)
	}
	if r.Message == "" {
		return fmt.Errorf("error message is required")
	}
	return nil
}

// StepStartedEvent signals the start of a step within an agent run.
type StepStartedEvent struct {
	BaseEvent
	StepName string `json:"stepName"` // Name of the step
}

// EventTypeName returns the concrete type name.
func (s *StepStartedEvent) EventTypeName() string {
	return "StepStartedEvent"
}

// Validate checks if the StepStartedEvent is valid.
func (s *StepStartedEvent) Validate() error {
	if err := s.BaseEvent.Validate(); err != nil {
		return err
	}
	if s.Type != EventTypeStepStarted {
		return fmt.Errorf("step started event must have STEP_STARTED type, got: %s", s.Type)
	}
	if s.StepName == "" {
		return fmt.Errorf("step name is required")
	}
	return nil
}

// StepFinishedEvent signals the completion of a step within an agent run.
type StepFinishedEvent struct {
	BaseEvent
	StepName string `json:"stepName"` // Name of the step
}

// EventTypeName returns the concrete type name.
func (s *StepFinishedEvent) EventTypeName() string {
	return "StepFinishedEvent"
}

// Validate checks if the StepFinishedEvent is valid.
func (s *StepFinishedEvent) Validate() error {
	if err := s.BaseEvent.Validate(); err != nil {
		return err
	}
	if s.Type != EventTypeStepFinished {
		return fmt.Errorf("step finished event must have STEP_FINISHED type, got: %s", s.Type)
	}
	if s.StepName == "" {
		return fmt.Errorf("step name is required")
	}
	return nil
}

// Text Message Events

// TextMessageStartEvent signals the start of a text message.
type TextMessageStartEvent struct {
	BaseEvent
	MessageID string `json:"messageId"` // Unique identifier for the message
	Role      Role   `json:"role"`      // Role is always "assistant" for text messages
}

// EventTypeName returns the concrete type name.
func (t *TextMessageStartEvent) EventTypeName() string {
	return "TextMessageStartEvent"
}

// Validate checks if the TextMessageStartEvent is valid.
func (t *TextMessageStartEvent) Validate() error {
	if err := t.BaseEvent.Validate(); err != nil {
		return err
	}
	if t.Type != EventTypeTextMessageStart {
		return fmt.Errorf("text message start event must have TEXT_MESSAGE_START type, got: %s", t.Type)
	}
	if t.MessageID == "" {
		return fmt.Errorf("message ID is required")
	}
	if t.Role != RoleAssistant {
		return fmt.Errorf("text message role must be assistant, got: %s", t.Role)
	}
	return nil
}

// TextMessageContentEvent represents a chunk of content in a streaming text message.
type TextMessageContentEvent struct {
	BaseEvent
	MessageID string `json:"messageId"` // Matches the ID from TextMessageStartEvent
	Delta     string `json:"delta"`     // Text content chunk (non-empty)
}

// EventTypeName returns the concrete type name.
func (t *TextMessageContentEvent) EventTypeName() string {
	return "TextMessageContentEvent"
}

// Validate checks if the TextMessageContentEvent is valid.
func (t *TextMessageContentEvent) Validate() error {
	if err := t.BaseEvent.Validate(); err != nil {
		return err
	}
	if t.Type != EventTypeTextMessageContent {
		return fmt.Errorf("text message content event must have TEXT_MESSAGE_CONTENT type, got: %s", t.Type)
	}
	if t.MessageID == "" {
		return fmt.Errorf("message ID is required")
	}
	if t.Delta == "" {
		return fmt.Errorf("delta must not be empty")
	}
	return nil
}

// TextMessageEndEvent signals the end of a text message.
type TextMessageEndEvent struct {
	BaseEvent
	MessageID string `json:"messageId"` // Matches the ID from TextMessageStartEvent
}

// EventTypeName returns the concrete type name.
func (t *TextMessageEndEvent) EventTypeName() string {
	return "TextMessageEndEvent"
}

// Validate checks if the TextMessageEndEvent is valid.
func (t *TextMessageEndEvent) Validate() error {
	if err := t.BaseEvent.Validate(); err != nil {
		return err
	}
	if t.Type != EventTypeTextMessageEnd {
		return fmt.Errorf("text message end event must have TEXT_MESSAGE_END type, got: %s", t.Type)
	}
	if t.MessageID == "" {
		return fmt.Errorf("message ID is required")
	}
	return nil
}

// Tool Call Events

// ToolCallStartEvent signals the start of a tool call.
type ToolCallStartEvent struct {
	BaseEvent
	ToolCallID      string `json:"toolCallId"`                // Unique identifier for the tool call
	ToolCallName    string `json:"toolCallName"`              // Name of the tool being called
	ParentMessageID string `json:"parentMessageId,omitempty"` // ID of the parent message
}

// EventTypeName returns the concrete type name.
func (t *ToolCallStartEvent) EventTypeName() string {
	return "ToolCallStartEvent"
}

// Validate checks if the ToolCallStartEvent is valid.
func (t *ToolCallStartEvent) Validate() error {
	if err := t.BaseEvent.Validate(); err != nil {
		return err
	}
	if t.Type != EventTypeToolCallStart {
		return fmt.Errorf("tool call start event must have TOOL_CALL_START type, got: %s", t.Type)
	}
	if t.ToolCallID == "" {
		return fmt.Errorf("tool call ID is required")
	}
	if t.ToolCallName == "" {
		return fmt.Errorf("tool call name is required")
	}
	return nil
}

// ToolCallArgsEvent represents a chunk of argument data for a tool call.
type ToolCallArgsEvent struct {
	BaseEvent
	ToolCallID string `json:"toolCallId"` // Matches the ID from ToolCallStartEvent
	Delta      string `json:"delta"`      // Argument data chunk
}

// EventTypeName returns the concrete type name.
func (t *ToolCallArgsEvent) EventTypeName() string {
	return "ToolCallArgsEvent"
}

// Validate checks if the ToolCallArgsEvent is valid.
func (t *ToolCallArgsEvent) Validate() error {
	if err := t.BaseEvent.Validate(); err != nil {
		return err
	}
	if t.Type != EventTypeToolCallArgs {
		return fmt.Errorf("tool call args event must have TOOL_CALL_ARGS type, got: %s", t.Type)
	}
	if t.ToolCallID == "" {
		return fmt.Errorf("tool call ID is required")
	}
	// Delta can be empty for tool call args
	return nil
}

// ToolCallEndEvent signals the end of a tool call.
type ToolCallEndEvent struct {
	BaseEvent
	ToolCallID string `json:"toolCallId"` // Matches the ID from ToolCallStartEvent
}

// EventTypeName returns the concrete type name.
func (t *ToolCallEndEvent) EventTypeName() string {
	return "ToolCallEndEvent"
}

// Validate checks if the ToolCallEndEvent is valid.
func (t *ToolCallEndEvent) Validate() error {
	if err := t.BaseEvent.Validate(); err != nil {
		return err
	}
	if t.Type != EventTypeToolCallEnd {
		return fmt.Errorf("tool call end event must have TOOL_CALL_END type, got: %s", t.Type)
	}
	if t.ToolCallID == "" {
		return fmt.Errorf("tool call ID is required")
	}
	return nil
}

// ToolCallResultEvent provides the result of a tool call execution.
type ToolCallResultEvent struct {
	BaseEvent
	MessageID  string `json:"messageId"`      // ID of the conversation message this result belongs to
	ToolCallID string `json:"toolCallId"`     // Matches the ID from the corresponding ToolCallStartEvent
	Content    string `json:"content"`        // The actual result/output content from the tool execution
	Role       Role   `json:"role,omitempty"` // Role identifier, typically "tool" for tool results
}

// EventTypeName returns the concrete type name.
func (t *ToolCallResultEvent) EventTypeName() string {
	return "ToolCallResultEvent"
}

// Validate checks if the ToolCallResultEvent is valid.
func (t *ToolCallResultEvent) Validate() error {
	if err := t.BaseEvent.Validate(); err != nil {
		return err
	}
	if t.Type != EventTypeToolCallResult {
		return fmt.Errorf("tool call result event must have TOOL_CALL_RESULT type, got: %s", t.Type)
	}
	if t.MessageID == "" {
		return fmt.Errorf("message ID is required")
	}
	if t.ToolCallID == "" {
		return fmt.Errorf("tool call ID is required")
	}
	if t.Content == "" {
		return fmt.Errorf("content is required")
	}
	if t.Role != "" && !t.Role.IsValid() {
		return fmt.Errorf("invalid role: %s", t.Role)
	}
	return nil
}

// State Management Events

// StateSnapshotEvent provides a complete snapshot of an agent's state.
type StateSnapshotEvent struct {
	BaseEvent
	Snapshot State `json:"snapshot"` // Complete state snapshot
}

// EventTypeName returns the concrete type name.
func (s *StateSnapshotEvent) EventTypeName() string {
	return "StateSnapshotEvent"
}

// Validate checks if the StateSnapshotEvent is valid.
func (s *StateSnapshotEvent) Validate() error {
	if err := s.BaseEvent.Validate(); err != nil {
		return err
	}
	if s.Type != EventTypeStateSnapshot {
		return fmt.Errorf("state snapshot event must have STATE_SNAPSHOT type, got: %s", s.Type)
	}
	if s.Snapshot == nil {
		return fmt.Errorf("snapshot is required")
	}
	return nil
}

// StateDeltaEvent provides a partial update to an agent's state using JSON Patch.
type StateDeltaEvent struct {
	BaseEvent
	Delta []interface{} `json:"delta"` // Array of JSON Patch operations (RFC 6902)
}

// EventTypeName returns the concrete type name.
func (s *StateDeltaEvent) EventTypeName() string {
	return "StateDeltaEvent"
}

// Validate checks if the StateDeltaEvent is valid.
func (s *StateDeltaEvent) Validate() error {
	if err := s.BaseEvent.Validate(); err != nil {
		return err
	}
	if s.Type != EventTypeStateDelta {
		return fmt.Errorf("state delta event must have STATE_DELTA type, got: %s", s.Type)
	}
	if s.Delta == nil {
		return fmt.Errorf("delta is required")
	}
	return nil
}

// MessagesSnapshotEvent provides a snapshot of all messages in a conversation.
type MessagesSnapshotEvent struct {
	BaseEvent
	Messages []Message `json:"messages"` // Array of message objects
}

// EventTypeName returns the concrete type name.
func (m *MessagesSnapshotEvent) EventTypeName() string {
	return "MessagesSnapshotEvent"
}

// Validate checks if the MessagesSnapshotEvent is valid.
func (m *MessagesSnapshotEvent) Validate() error {
	if err := m.BaseEvent.Validate(); err != nil {
		return err
	}
	if m.Type != EventTypeMessagesSnapshot {
		return fmt.Errorf("messages snapshot event must have MESSAGES_SNAPSHOT type, got: %s", m.Type)
	}
	if m.Messages == nil {
		return fmt.Errorf("messages are required")
	}

	// Validate each message
	for i, msg := range m.Messages {
		if err := msg.Validate(); err != nil {
			return fmt.Errorf("invalid message at index %d: %w", i, err)
		}
	}

	return nil
}

// Special Events

// RawEvent is used to pass through events from external systems.
type RawEvent struct {
	BaseEvent
	Event  interface{} `json:"event"`            // Original event data
	Source string      `json:"source,omitempty"` // Source of the event
}

// EventTypeName returns the concrete type name.
func (r *RawEvent) EventTypeName() string {
	return "RawEvent"
}

// Validate checks if the RawEvent is valid.
func (r *RawEvent) Validate() error {
	if err := r.BaseEvent.Validate(); err != nil {
		return err
	}
	if r.Type != EventTypeRaw {
		return fmt.Errorf("raw event must have RAW type, got: %s", r.Type)
	}
	if r.Event == nil {
		return fmt.Errorf("event is required")
	}
	return nil
}

// CustomEvent is used for application-specific custom events.
type CustomEvent struct {
	BaseEvent
	Name  string      `json:"name"`  // Name of the custom event
	Value interface{} `json:"value"` // Value associated with the event
}

// EventTypeName returns the concrete type name.
func (c *CustomEvent) EventTypeName() string {
	return "CustomEvent"
}

// Validate checks if the CustomEvent is valid.
func (c *CustomEvent) Validate() error {
	if err := c.BaseEvent.Validate(); err != nil {
		return err
	}
	if c.Type != EventTypeCustom {
		return fmt.Errorf("custom event must have CUSTOM type, got: %s", c.Type)
	}
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}
	if c.Value == nil {
		return fmt.Errorf("value is required")
	}
	return nil
}
