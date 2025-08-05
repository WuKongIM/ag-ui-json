package agui

import (
	"bytes"
	"testing"
)

func TestEventEncoding(t *testing.T) {
	tests := []struct {
		name  string
		event Event
	}{
		{
			name:  "RunStartedEvent",
			event: NewRunStartedEvent("thread_123", "run_456"),
		},
		{
			name:  "RunFinishedEvent",
			event: NewRunFinishedEvent("thread_123", "run_456", map[string]string{"status": "success"}),
		},
		{
			name:  "RunErrorEvent",
			event: NewRunErrorEvent("Something went wrong", "ERR_001"),
		},
		{
			name:  "TextMessageStartEvent",
			event: NewTextMessageStartEvent("msg_789"),
		},
		{
			name:  "TextMessageContentEvent",
			event: NewTextMessageContentEvent("msg_789", "Hello, world!"),
		},
		{
			name:  "TextMessageEndEvent",
			event: NewTextMessageEndEvent("msg_789"),
		},
		{
			name:  "ToolCallStartEvent",
			event: NewToolCallStartEvent("tool_call_123", "search", "msg_parent"),
		},
		{
			name:  "ToolCallArgsEvent",
			event: NewToolCallArgsEvent("tool_call_123", `{"query": "test"}`),
		},
		{
			name:  "ToolCallEndEvent",
			event: NewToolCallEndEvent("tool_call_123"),
		},
		{
			name:  "ToolCallResultEvent",
			event: NewToolCallResultEvent("msg_result", "tool_call_123", "Search results found"),
		},
		{
			name:  "StateSnapshotEvent",
			event: NewStateSnapshotEvent(map[string]interface{}{"key": "value"}),
		},
		{
			name:  "StateDeltaEvent",
			event: NewStateDeltaEvent([]interface{}{map[string]interface{}{"op": "replace", "path": "/key", "value": "new_value"}}),
		},
		{
			name:  "CustomEvent",
			event: NewCustomEvent("test_event", map[string]string{"data": "test"}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test encoding
			data, err := EncodeEvent(tt.event)
			if err != nil {
				t.Fatalf("Failed to encode event: %v", err)
			}

			// Test decoding
			decoded, err := DecodeEventFromBytes(data)
			if err != nil {
				t.Fatalf("Failed to decode event: %v", err)
			}

			// Verify type matches
			if decoded.GetType() != tt.event.GetType() {
				t.Errorf("Event type mismatch: expected %s, got %s", tt.event.GetType(), decoded.GetType())
			}

			// Verify event type name matches
			if decoded.EventTypeName() != tt.event.EventTypeName() {
				t.Errorf("Event type name mismatch: expected %s, got %s", tt.event.EventTypeName(), decoded.EventTypeName())
			}
		})
	}
}

func TestMessageEncoding(t *testing.T) {
	tests := []struct {
		name    string
		message Message
	}{
		{
			name:    "DeveloperMessage",
			message: NewDeveloperMessage("msg_1", "Debug info", "dev_user"),
		},
		{
			name:    "SystemMessage",
			message: NewSystemMessage("msg_2", "System initialization", ""),
		},
		{
			name:    "AssistantMessage",
			message: NewAssistantMessage("msg_3", "Hello! How can I help?", "assistant", nil),
		},
		{
			name:    "UserMessage",
			message: NewUserMessage("msg_4", "What's the weather like?", "user_123"),
		},
		{
			name:    "ToolMessage",
			message: NewToolMessage("msg_5", "Weather is sunny", "tool_call_456", "", "weather_tool"),
		},
		{
			name: "AssistantMessageWithToolCalls",
			message: NewAssistantMessage("msg_6", "Let me search for that", "assistant", []ToolCall{
				{
					ID:   "tool_call_789",
					Type: ToolCallTypeFunction,
					Function: FunctionCall{
						Name:      "search",
						Arguments: `{"query": "test"}`,
					},
				},
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test encoding
			data, err := EncodeMessage(tt.message)
			if err != nil {
				t.Fatalf("Failed to encode message: %v", err)
			}

			// Test decoding
			decoded, err := DecodeMessageFromBytes(data)
			if err != nil {
				t.Fatalf("Failed to decode message: %v", err)
			}

			// Verify role matches
			if decoded.GetRole() != tt.message.GetRole() {
				t.Errorf("Message role mismatch: expected %s, got %s", tt.message.GetRole(), decoded.GetRole())
			}

			// Verify ID matches
			if decoded.GetID() != tt.message.GetID() {
				t.Errorf("Message ID mismatch: expected %s, got %s", tt.message.GetID(), decoded.GetID())
			}

			// Verify message type name matches
			if decoded.MessageType() != tt.message.MessageType() {
				t.Errorf("Message type mismatch: expected %s, got %s", tt.message.MessageType(), decoded.MessageType())
			}
		})
	}
}

func TestStreamDecoding(t *testing.T) {
	// Create a stream of events
	events := []Event{
		&RunStartedEvent{
			BaseEvent: BaseEvent{Type: EventTypeRunStarted},
			ThreadID:  "thread_1",
			RunID:     "run_1",
		},
		&TextMessageStartEvent{
			BaseEvent: BaseEvent{Type: EventTypeTextMessageStart},
			MessageID: "msg_1",
			Role:      RoleAssistant,
		},
		&TextMessageContentEvent{
			BaseEvent: BaseEvent{Type: EventTypeTextMessageContent},
			MessageID: "msg_1",
			Delta:     "Hello",
		},
		&TextMessageContentEvent{
			BaseEvent: BaseEvent{Type: EventTypeTextMessageContent},
			MessageID: "msg_1",
			Delta:     " world!",
		},
		&TextMessageEndEvent{
			BaseEvent: BaseEvent{Type: EventTypeTextMessageEnd},
			MessageID: "msg_1",
		},
		&RunFinishedEvent{
			BaseEvent: BaseEvent{Type: EventTypeRunFinished},
			ThreadID:  "thread_1",
			RunID:     "run_1",
		},
	}

	// Encode events to JSON stream
	var buf bytes.Buffer
	for _, event := range events {
		data, err := EncodeEvent(event)
		if err != nil {
			t.Fatalf("Failed to encode event: %v", err)
		}
		buf.Write(data)
		buf.WriteString("\n")
	}

	// Decode events from stream
	decoder := NewStreamDecoder(&buf)
	eventChan, errorChan := decoder.DecodeEvents()

	var decodedEvents []Event
	for {
		select {
		case event, ok := <-eventChan:
			if !ok {
				goto done
			}
			decodedEvents = append(decodedEvents, event)
		case err := <-errorChan:
			if err != nil {
				t.Fatalf("Stream decoding error: %v", err)
			}
		}
	}

done:
	if len(decodedEvents) != len(events) {
		t.Errorf("Expected %d events, got %d", len(events), len(decodedEvents))
	}

	for i, event := range decodedEvents {
		if event.GetType() != events[i].GetType() {
			t.Errorf("Event %d type mismatch: expected %s, got %s", i, events[i].GetType(), event.GetType())
		}
	}
}

func TestValidation(t *testing.T) {
	tests := []struct {
		name        string
		event       Event
		shouldError bool
	}{
		{
			name: "ValidRunStartedEvent",
			event: &RunStartedEvent{
				BaseEvent: BaseEvent{Type: EventTypeRunStarted},
				ThreadID:  "thread_123",
				RunID:     "run_456",
			},
			shouldError: false,
		},
		{
			name: "InvalidRunStartedEvent_MissingThreadID",
			event: &RunStartedEvent{
				BaseEvent: BaseEvent{Type: EventTypeRunStarted},
				ThreadID:  "",
				RunID:     "run_456",
			},
			shouldError: true,
		},
		{
			name: "InvalidRunStartedEvent_WrongType",
			event: &RunStartedEvent{
				BaseEvent: BaseEvent{Type: EventTypeRunFinished},
				ThreadID:  "thread_123",
				RunID:     "run_456",
			},
			shouldError: true,
		},
		{
			name: "ValidTextMessageContentEvent",
			event: &TextMessageContentEvent{
				BaseEvent: BaseEvent{Type: EventTypeTextMessageContent},
				MessageID: "msg_123",
				Delta:     "Hello",
			},
			shouldError: false,
		},
		{
			name: "InvalidTextMessageContentEvent_EmptyDelta",
			event: &TextMessageContentEvent{
				BaseEvent: BaseEvent{Type: EventTypeTextMessageContent},
				MessageID: "msg_123",
				Delta:     "",
			},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.event.Validate()
			if tt.shouldError && err == nil {
				t.Error("Expected validation error, but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Unexpected validation error: %v", err)
			}
		})
	}
}

func TestRunAgentInputValidation(t *testing.T) {
	validInput := &RunAgentInput{
		ThreadID: "thread_123",
		RunID:    "run_456",
		State:    map[string]interface{}{"key": "value"},
		Messages: []Message{
			NewUserMessage("msg_1", "Hello", "user"),
		},
		Tools: []Tool{
			{
				Name:        "search",
				Description: "Search tool",
				Parameters:  map[string]interface{}{"type": "object"},
			},
		},
		Context: []Context{
			{
				Description: "Test context",
				Value:       "test value",
			},
		},
		ForwardedProps: map[string]interface{}{"prop": "value"},
	}

	if err := validInput.Validate(); err != nil {
		t.Errorf("Valid input should not produce error: %v", err)
	}

	// Test invalid input
	invalidInput := &RunAgentInput{
		ThreadID: "", // Missing thread ID
		RunID:    "run_456",
	}

	if err := invalidInput.Validate(); err == nil {
		t.Error("Invalid input should produce error")
	}
}
