package agui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// ExampleBasicEventEncoding demonstrates basic event encoding and decoding.
func ExampleBasicEventEncoding() {
	// Create a text message start event
	event := NewTextMessageStartEvent("msg_123")

	// Encode to JSON
	data, err := EncodeEvent(event)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Encoded event: %s\n", string(data))

	// Decode from JSON
	decoded, err := DecodeEventFromBytes(data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Event type: %s\n", decoded.GetType())
	fmt.Printf("Event type name: %s\n", decoded.EventTypeName())

	// Output will be similar to:
	// Encoded event: {"type":"TEXT_MESSAGE_START","timestamp":1234567890123,"messageId":"msg_123","role":"assistant"}
	// Event type: TEXT_MESSAGE_START
	// Event type name: TextMessageStartEvent
}

// ExampleMessageEncoding demonstrates message encoding and decoding.
func ExampleMessageEncoding() {
	// Create a user message
	message := NewUserMessage("msg_456", "Hello, how can you help me?", "user_123")

	// Encode to JSON
	data, err := EncodeMessage(message)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Encoded message: %s\n", string(data))

	// Decode from JSON
	decoded, err := DecodeMessageFromBytes(data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Message role: %s\n", decoded.GetRole())
	fmt.Printf("Message ID: %s\n", decoded.GetID())
	fmt.Printf("Message type: %s\n", decoded.MessageType())

	// Output:
	// Encoded message: {"id":"msg_456","role":"user","name":"user_123","content":"Hello, how can you help me?"}
	// Message role: user
	// Message ID: msg_456
	// Message type: UserMessage
}

// ExampleStreamingEvents demonstrates how to handle streaming events.
func ExampleStreamingEvents() {
	// Simulate a conversation flow with streaming events
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
			Delta:     " there!",
		},
		&TextMessageContentEvent{
			BaseEvent: BaseEvent{Type: EventTypeTextMessageContent},
			MessageID: "msg_1",
			Delta:     " How",
		},
		&TextMessageContentEvent{
			BaseEvent: BaseEvent{Type: EventTypeTextMessageContent},
			MessageID: "msg_1",
			Delta:     " can",
		},
		&TextMessageContentEvent{
			BaseEvent: BaseEvent{Type: EventTypeTextMessageContent},
			MessageID: "msg_1",
			Delta:     " I",
		},
		&TextMessageContentEvent{
			BaseEvent: BaseEvent{Type: EventTypeTextMessageContent},
			MessageID: "msg_1",
			Delta:     " help?",
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

	// Create a JSON stream
	var buf bytes.Buffer
	for _, event := range events {
		data, err := EncodeEvent(event)
		if err != nil {
			log.Fatal(err)
		}
		buf.Write(data)
		buf.WriteString("\n")
	}

	// Decode the stream
	decoder := NewStreamDecoder(&buf)
	eventChan, errorChan := decoder.DecodeEvents()

	// Process events as they arrive
	var messageContent strings.Builder
	for {
		select {
		case event, ok := <-eventChan:
			if !ok {
				fmt.Printf("Final message: %s\n", messageContent.String())
				return
			}

			switch e := event.(type) {
			case *RunStartedEvent:
				fmt.Printf("Run started: %s\n", e.RunID)
			case *TextMessageStartEvent:
				fmt.Printf("Message started: %s\n", e.MessageID)
			case *TextMessageContentEvent:
				messageContent.WriteString(e.Delta)
			case *TextMessageEndEvent:
				fmt.Printf("Message ended: %s\n", e.MessageID)
			case *RunFinishedEvent:
				fmt.Printf("Run finished: %s\n", e.RunID)
			}

		case err := <-errorChan:
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// Output:
	// Run started: run_1
	// Message started: msg_1
	// Message ended: msg_1
	// Run finished: run_1
	// Final message: Hello there! How can I help?
}

// ExampleToolCallFlow demonstrates a complete tool call flow.
func ExampleToolCallFlow() {
	// Create tool call events
	events := []Event{
		&ToolCallStartEvent{
			BaseEvent:       BaseEvent{Type: EventTypeToolCallStart},
			ToolCallID:      "tool_call_1",
			ToolCallName:    "search",
			ParentMessageID: "msg_parent",
		},
		&ToolCallArgsEvent{
			BaseEvent:  BaseEvent{Type: EventTypeToolCallArgs},
			ToolCallID: "tool_call_1",
			Delta:      `{"query":`,
		},
		&ToolCallArgsEvent{
			BaseEvent:  BaseEvent{Type: EventTypeToolCallArgs},
			ToolCallID: "tool_call_1",
			Delta:      `"weather"`,
		},
		&ToolCallArgsEvent{
			BaseEvent:  BaseEvent{Type: EventTypeToolCallArgs},
			ToolCallID: "tool_call_1",
			Delta:      `,"location":`,
		},
		&ToolCallArgsEvent{
			BaseEvent:  BaseEvent{Type: EventTypeToolCallArgs},
			ToolCallID: "tool_call_1",
			Delta:      `"New York"}`,
		},
		&ToolCallEndEvent{
			BaseEvent:  BaseEvent{Type: EventTypeToolCallEnd},
			ToolCallID: "tool_call_1",
		},
		&ToolCallResultEvent{
			BaseEvent:  BaseEvent{Type: EventTypeToolCallResult},
			MessageID:  "msg_result",
			ToolCallID: "tool_call_1",
			Content:    "The weather in New York is sunny, 72°F",
			Role:       RoleTool,
		},
	}

	// Process the tool call flow
	var argsBuilder strings.Builder
	for _, event := range events {
		switch e := event.(type) {
		case *ToolCallStartEvent:
			fmt.Printf("Tool call started: %s (%s)\n", e.ToolCallName, e.ToolCallID)
		case *ToolCallArgsEvent:
			argsBuilder.WriteString(e.Delta)
		case *ToolCallEndEvent:
			fmt.Printf("Tool call arguments: %s\n", argsBuilder.String())
			fmt.Printf("Tool call ended: %s\n", e.ToolCallID)
		case *ToolCallResultEvent:
			fmt.Printf("Tool call result: %s\n", e.Content)
		}
	}

	// Output:
	// Tool call started: search (tool_call_1)
	// Tool call arguments: {"query":"weather","location":"New York"}
	// Tool call ended: tool_call_1
	// Tool call result: The weather in New York is sunny, 72°F
}

// ExampleRunAgentInput demonstrates creating and validating RunAgentInput.
func ExampleRunAgentInput() {
	// Create messages
	messages := []Message{
		NewSystemMessage("msg_1", "You are a helpful assistant.", ""),
		NewUserMessage("msg_2", "What's the capital of France?", "user_123"),
	}

	// Create tools
	tools := []Tool{
		{
			Name:        "search",
			Description: "Search for information",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "The search query",
					},
				},
				"required": []string{"query"},
			},
		},
	}

	// Create context
	context := []Context{
		{
			Description: "User preferences",
			Value:       "The user prefers concise answers",
		},
	}

	// Create RunAgentInput
	input := &RunAgentInput{
		ThreadID:       GenerateThreadID(),
		RunID:          GenerateRunID(),
		State:          map[string]interface{}{"conversation_count": 1},
		Messages:       messages,
		Tools:          tools,
		Context:        context,
		ForwardedProps: map[string]interface{}{"client_version": "1.0.0"},
	}

	// Validate the input
	if err := input.Validate(); err != nil {
		log.Fatal(err)
	}

	// Encode to JSON
	data, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("RunAgentInput JSON:\n%s\n", string(data))

	// Output will be similar to:
	// RunAgentInput JSON:
	// {
	//   "threadId": "thread_1234567890123",
	//   "runId": "run_1234567890123",
	//   "state": {
	//     "conversation_count": 1
	//   },
	//   "messages": [
	//     {
	//       "id": "msg_1",
	//       "role": "system",
	//       "content": "You are a helpful assistant."
	//     },
	//     {
	//       "id": "msg_2",
	//       "role": "user",
	//       "name": "user_123",
	//       "content": "What's the capital of France?"
	//     }
	//   ],
	//   "tools": [
	//     {
	//       "name": "search",
	//       "description": "Search for information",
	//       "parameters": {
	//         "type": "object",
	//         "properties": {
	//           "query": {
	//             "type": "string",
	//             "description": "The search query"
	//           }
	//         },
	//         "required": ["query"]
	//       }
	//     }
	//   ],
	//   "context": [
	//     {
	//       "description": "User preferences",
	//       "value": "The user prefers concise answers"
	//     }
	//   ],
	//   "forwardedProps": {
	//     "client_version": "1.0.0"
	//   }
	// }
}

// ExampleStateManagement demonstrates state snapshot and delta events.
func ExampleStateManagement() {
	// Initial state
	initialState := map[string]interface{}{
		"user_id":            "user_123",
		"conversation_count": 1,
		"preferences": map[string]interface{}{
			"language": "en",
			"theme":    "dark",
		},
	}

	// Create state snapshot event
	snapshotEvent := NewStateSnapshotEvent(initialState)
	fmt.Printf("State snapshot created for timestamp: %d\n", *snapshotEvent.GetTimestamp())

	// Create state delta (JSON Patch operations)
	delta := []interface{}{
		map[string]interface{}{
			"op":    "replace",
			"path":  "/conversation_count",
			"value": 2,
		},
		map[string]interface{}{
			"op":    "replace",
			"path":  "/preferences/theme",
			"value": "light",
		},
	}

	deltaEvent := NewStateDeltaEvent(delta)
	fmt.Printf("State delta created with %d operations\n", len(deltaEvent.Delta))

	// Encode and decode to verify
	snapshotData, _ := EncodeEvent(snapshotEvent)
	deltaData, _ := EncodeEvent(deltaEvent)

	decodedSnapshot, _ := DecodeEventFromBytes(snapshotData)
	decodedDelta, _ := DecodeEventFromBytes(deltaData)

	fmt.Printf("Decoded snapshot type: %s\n", decodedSnapshot.EventTypeName())
	fmt.Printf("Decoded delta type: %s\n", decodedDelta.EventTypeName())

	// Output will be similar to:
	// State snapshot created for timestamp: 1234567890123
	// State delta created with 2 operations
	// Decoded snapshot type: StateSnapshotEvent
	// Decoded delta type: StateDeltaEvent
}
