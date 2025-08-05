package agui

import (
	"fmt"
	"time"
)

// Event Factory Functions
// These functions provide convenient ways to create properly initialized events.

// NewRunStartedEvent creates a new RunStartedEvent with the current timestamp.
func NewRunStartedEvent(threadID, runID string) *RunStartedEvent {
	event := &RunStartedEvent{
		BaseEvent: BaseEvent{
			Type: EventTypeRunStarted,
		},
		ThreadID: threadID,
		RunID:    runID,
	}
	event.SetTimestamp()
	return event
}

// NewRunFinishedEvent creates a new RunFinishedEvent with the current timestamp.
func NewRunFinishedEvent(threadID, runID string, result interface{}) *RunFinishedEvent {
	event := &RunFinishedEvent{
		BaseEvent: BaseEvent{
			Type: EventTypeRunFinished,
		},
		ThreadID: threadID,
		RunID:    runID,
		Result:   result,
	}
	event.SetTimestamp()
	return event
}

// NewRunErrorEvent creates a new RunErrorEvent with the current timestamp.
func NewRunErrorEvent(message, code string) *RunErrorEvent {
	event := &RunErrorEvent{
		BaseEvent: BaseEvent{
			Type: EventTypeRunError,
		},
		Message: message,
		Code:    code,
	}
	event.SetTimestamp()
	return event
}

// NewStepStartedEvent creates a new StepStartedEvent with the current timestamp.
func NewStepStartedEvent(stepName string) *StepStartedEvent {
	event := &StepStartedEvent{
		BaseEvent: BaseEvent{
			Type: EventTypeStepStarted,
		},
		StepName: stepName,
	}
	event.SetTimestamp()
	return event
}

// NewStepFinishedEvent creates a new StepFinishedEvent with the current timestamp.
func NewStepFinishedEvent(stepName string) *StepFinishedEvent {
	event := &StepFinishedEvent{
		BaseEvent: BaseEvent{
			Type: EventTypeStepFinished,
		},
		StepName: stepName,
	}
	event.SetTimestamp()
	return event
}

// NewTextMessageStartEvent creates a new TextMessageStartEvent with the current timestamp.
func NewTextMessageStartEvent(messageID string) *TextMessageStartEvent {
	event := &TextMessageStartEvent{
		BaseEvent: BaseEvent{
			Type: EventTypeTextMessageStart,
		},
		MessageID: messageID,
		Role:      RoleAssistant,
	}
	event.SetTimestamp()
	return event
}

// NewTextMessageContentEvent creates a new TextMessageContentEvent with the current timestamp.
func NewTextMessageContentEvent(messageID, delta string) *TextMessageContentEvent {
	event := &TextMessageContentEvent{
		BaseEvent: BaseEvent{
			Type: EventTypeTextMessageContent,
		},
		MessageID: messageID,
		Delta:     delta,
	}
	event.SetTimestamp()
	return event
}

// NewTextMessageEndEvent creates a new TextMessageEndEvent with the current timestamp.
func NewTextMessageEndEvent(messageID string) *TextMessageEndEvent {
	event := &TextMessageEndEvent{
		BaseEvent: BaseEvent{
			Type: EventTypeTextMessageEnd,
		},
		MessageID: messageID,
	}
	event.SetTimestamp()
	return event
}

// NewToolCallStartEvent creates a new ToolCallStartEvent with the current timestamp.
func NewToolCallStartEvent(toolCallID, toolCallName, parentMessageID string) *ToolCallStartEvent {
	event := &ToolCallStartEvent{
		BaseEvent: BaseEvent{
			Type: EventTypeToolCallStart,
		},
		ToolCallID:      toolCallID,
		ToolCallName:    toolCallName,
		ParentMessageID: parentMessageID,
	}
	event.SetTimestamp()
	return event
}

// NewToolCallArgsEvent creates a new ToolCallArgsEvent with the current timestamp.
func NewToolCallArgsEvent(toolCallID, delta string) *ToolCallArgsEvent {
	event := &ToolCallArgsEvent{
		BaseEvent: BaseEvent{
			Type: EventTypeToolCallArgs,
		},
		ToolCallID: toolCallID,
		Delta:      delta,
	}
	event.SetTimestamp()
	return event
}

// NewToolCallEndEvent creates a new ToolCallEndEvent with the current timestamp.
func NewToolCallEndEvent(toolCallID string) *ToolCallEndEvent {
	event := &ToolCallEndEvent{
		BaseEvent: BaseEvent{
			Type: EventTypeToolCallEnd,
		},
		ToolCallID: toolCallID,
	}
	event.SetTimestamp()
	return event
}

// NewToolCallResultEvent creates a new ToolCallResultEvent with the current timestamp.
func NewToolCallResultEvent(messageID, toolCallID, content string) *ToolCallResultEvent {
	event := &ToolCallResultEvent{
		BaseEvent: BaseEvent{
			Type: EventTypeToolCallResult,
		},
		MessageID:  messageID,
		ToolCallID: toolCallID,
		Content:    content,
		Role:       RoleTool,
	}
	event.SetTimestamp()
	return event
}

// NewStateSnapshotEvent creates a new StateSnapshotEvent with the current timestamp.
func NewStateSnapshotEvent(snapshot State) *StateSnapshotEvent {
	event := &StateSnapshotEvent{
		BaseEvent: BaseEvent{
			Type: EventTypeStateSnapshot,
		},
		Snapshot: snapshot,
	}
	event.SetTimestamp()
	return event
}

// NewStateDeltaEvent creates a new StateDeltaEvent with the current timestamp.
func NewStateDeltaEvent(delta []interface{}) *StateDeltaEvent {
	event := &StateDeltaEvent{
		BaseEvent: BaseEvent{
			Type: EventTypeStateDelta,
		},
		Delta: delta,
	}
	event.SetTimestamp()
	return event
}

// NewMessagesSnapshotEvent creates a new MessagesSnapshotEvent with the current timestamp.
func NewMessagesSnapshotEvent(messages []Message) *MessagesSnapshotEvent {
	event := &MessagesSnapshotEvent{
		BaseEvent: BaseEvent{
			Type: EventTypeMessagesSnapshot,
		},
		Messages: messages,
	}
	event.SetTimestamp()
	return event
}

// NewRawEvent creates a new RawEvent with the current timestamp.
func NewRawEvent(event interface{}, source string) *RawEvent {
	rawEvent := &RawEvent{
		BaseEvent: BaseEvent{
			Type: EventTypeRaw,
		},
		Event:  event,
		Source: source,
	}
	rawEvent.SetTimestamp()
	return rawEvent
}

// NewCustomEvent creates a new CustomEvent with the current timestamp.
func NewCustomEvent(name string, value interface{}) *CustomEvent {
	event := &CustomEvent{
		BaseEvent: BaseEvent{
			Type: EventTypeCustom,
		},
		Name:  name,
		Value: value,
	}
	event.SetTimestamp()
	return event
}

// Message Factory Functions
// These functions provide convenient ways to create properly initialized messages.

// NewDeveloperMessage creates a new DeveloperMessage.
func NewDeveloperMessage(id, content, name string) *DeveloperMessage {
	return &DeveloperMessage{
		BaseMessage: BaseMessage{
			ID:   id,
			Role: RoleDeveloper,
			Name: name,
		},
		Content: content,
	}
}

// NewSystemMessage creates a new SystemMessage.
func NewSystemMessage(id, content, name string) *SystemMessage {
	return &SystemMessage{
		BaseMessage: BaseMessage{
			ID:   id,
			Role: RoleSystem,
			Name: name,
		},
		Content: content,
	}
}

// NewAssistantMessage creates a new AssistantMessage.
func NewAssistantMessage(id, content, name string, toolCalls []ToolCall) *AssistantMessage {
	return &AssistantMessage{
		BaseMessage: BaseMessage{
			ID:   id,
			Role: RoleAssistant,
			Name: name,
		},
		Content:   content,
		ToolCalls: toolCalls,
	}
}

// NewUserMessage creates a new UserMessage.
func NewUserMessage(id, content, name string) *UserMessage {
	return &UserMessage{
		BaseMessage: BaseMessage{
			ID:   id,
			Role: RoleUser,
			Name: name,
		},
		Content: content,
	}
}

// NewToolMessage creates a new ToolMessage.
func NewToolMessage(id, content, toolCallID, errorMsg, name string) *ToolMessage {
	return &ToolMessage{
		BaseMessage: BaseMessage{
			ID:   id,
			Role: RoleTool,
			Name: name,
		},
		Content:    content,
		ToolCallID: toolCallID,
		Error:      errorMsg,
	}
}

// Utility Functions

// GenerateMessageID generates a unique message ID based on the current timestamp.
func GenerateMessageID() string {
	return fmt.Sprintf("msg_%d", time.Now().UnixNano())
}

// GenerateRunID generates a unique run ID based on the current timestamp.
func GenerateRunID() string {
	return fmt.Sprintf("run_%d", time.Now().UnixNano())
}

// GenerateThreadID generates a unique thread ID based on the current timestamp.
func GenerateThreadID() string {
	return fmt.Sprintf("thread_%d", time.Now().UnixNano())
}

// GenerateToolCallID generates a unique tool call ID based on the current timestamp.
func GenerateToolCallID() string {
	return fmt.Sprintf("tool_call_%d", time.Now().UnixNano())
}
