// Package agui provides Go types and encoding/decoding functionality for the Agent User Interaction Protocol.
//
// The Agent User Interaction Protocol (AG-UI) is a standardized protocol for communication between
// front-end applications and AI agents. This package implements the complete AG-UI specification
// including all event types, message types, and core data structures.
//
// # Overview
//
// AG-UI follows an event-driven architecture where agents communicate with clients through a stream
// of typed events. The protocol supports:
//
//   - Lifecycle events (run started/finished, step tracking)
//   - Text message streaming (start, content chunks, end)
//   - Tool call execution (start, arguments, end, results)
//   - State management (snapshots and deltas)
//   - Custom and raw events for extensibility
//
// # Core Types
//
// The package defines several core types that represent the fundamental data structures:
//
//   - EventType: Enumeration of all possible event types
//   - Role: Enumeration of message sender roles (developer, system, assistant, user, tool)
//   - ToolCallType: Type of tool calls (currently only "function")
//   - State: Flexible type for agent state (can be any data structure)
//   - Context: Contextual information provided to agents
//   - Tool: Tool definitions that agents can use
//   - ToolCall: Tool calls made by agents
//   - RunAgentInput: Input parameters for running an agent
//
// # Events
//
// Events are the primary communication mechanism in AG-UI. All events implement the Event interface
// and include common properties like type, timestamp, and optional raw event data.
//
// ## Lifecycle Events
//
//   - RunStartedEvent: Signals the start of an agent run
//   - RunFinishedEvent: Signals successful completion of an agent run
//   - RunErrorEvent: Signals an error during an agent run
//   - StepStartedEvent: Signals the start of a step within an agent run
//   - StepFinishedEvent: Signals completion of a step within an agent run
//
// ## Text Message Events
//
//   - TextMessageStartEvent: Signals the start of a text message
//   - TextMessageContentEvent: Represents a chunk of content in a streaming text message
//   - TextMessageEndEvent: Signals the end of a text message
//
// ## Tool Call Events
//
//   - ToolCallStartEvent: Signals the start of a tool call
//   - ToolCallArgsEvent: Represents a chunk of argument data for a tool call
//   - ToolCallEndEvent: Signals the end of a tool call
//   - ToolCallResultEvent: Provides the result of a tool call execution
//
// ## State Management Events
//
//   - StateSnapshotEvent: Provides a complete snapshot of an agent's state
//   - StateDeltaEvent: Provides partial updates using JSON Patch operations
//   - MessagesSnapshotEvent: Provides a snapshot of all messages in a conversation
//
// ## Special Events
//
//   - RawEvent: Used to pass through events from external systems
//   - CustomEvent: Used for application-specific custom events
//
// # Messages
//
// Messages represent different types of communication in conversations. All messages implement
// the Message interface and include common properties like ID, role, and optional name.
//
//   - DeveloperMessage: Messages from developers
//   - SystemMessage: System messages
//   - AssistantMessage: Messages from AI assistants
//   - UserMessage: Messages from users
//   - ToolMessage: Messages from tools
//
// # Encoding and Decoding
//
// The package provides comprehensive encoding and decoding functionality:
//
//   - Encoder/Decoder: Stream-based encoding/decoding to/from io.Reader/Writer
//   - EncodeEvent/DecodeEventFromBytes: Direct encoding/decoding of events
//   - EncodeMessage/DecodeMessageFromBytes: Direct encoding/decoding of messages
//   - StreamDecoder: Specialized decoder for handling streaming events and messages
//
// # Validation
//
// All types include validation methods that check for required fields and constraints
// based on the AG-UI schema specifications. Validation is automatically performed
// during encoding operations.
//
// # Factory Functions
//
// The package provides convenient factory functions for creating properly initialized
// events and messages:
//
//   - NewRunStartedEvent, NewTextMessageStartEvent, etc.
//   - NewDeveloperMessage, NewUserMessage, etc.
//   - GenerateMessageID, GenerateRunID, etc.
//
// # Usage Examples
//
// ## Basic Event Encoding/Decoding
//
//	// Create an event
//	event := agui.NewTextMessageStartEvent("msg_123")
//
//	// Encode to JSON
//	data, err := agui.EncodeEvent(event)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Decode from JSON
//	decoded, err := agui.DecodeEventFromBytes(data)
//	if err != nil {
//		log.Fatal(err)
//	}
//
// ## Streaming Events
//
//	// Create a stream decoder
//	decoder := agui.NewStreamDecoder(reader)
//	eventChan, errorChan := decoder.DecodeEvents()
//
//	// Process events as they arrive
//	for {
//		select {
//		case event := <-eventChan:
//			// Handle event
//		case err := <-errorChan:
//			// Handle error
//		}
//	}
//
// ## Creating Messages
//
//	// Create a user message
//	message := agui.NewUserMessage("msg_456", "Hello!", "user_123")
//
//	// Validate the message
//	if err := message.Validate(); err != nil {
//		log.Fatal(err)
//	}
//
//	// Encode to JSON
//	data, err := agui.EncodeMessage(message)
//	if err != nil {
//		log.Fatal(err)
//	}
//
// ## Tool Call Flow
//
//	// Start a tool call
//	startEvent := agui.NewToolCallStartEvent("tool_call_1", "search", "msg_parent")
//
//	// Stream arguments
//	argsEvent := agui.NewToolCallArgsEvent("tool_call_1", `{"query": "test"}`)
//
//	// End the tool call
//	endEvent := agui.NewToolCallEndEvent("tool_call_1")
//
//	// Provide results
//	resultEvent := agui.NewToolCallResultEvent("msg_result", "tool_call_1", "Search results")
//
// ## State Management
//
//	// Create state snapshot
//	state := map[string]interface{}{"key": "value"}
//	snapshotEvent := agui.NewStateSnapshotEvent(state)
//
//	// Create state delta (JSON Patch)
//	delta := []interface{}{
//		map[string]interface{}{
//			"op":    "replace",
//			"path":  "/key",
//			"value": "new_value",
//		},
//	}
//	deltaEvent := agui.NewStateDeltaEvent(delta)
//
// # Error Handling
//
// The package defines several error types for different failure scenarios:
//
//   - ErrInvalidEventType: Invalid or unknown event type
//   - ErrInvalidMessageType: Invalid or unknown message type
//   - ErrInvalidStructure: Invalid message structure
//   - ErrUnmarshalFailed: JSON unmarshaling failed
//   - ErrMarshalFailed: JSON marshaling failed
//   - ErrValidationFailed: Validation failed
//
// # Thread Safety
//
// The encoding and decoding functions are thread-safe. However, individual event and message
// instances are not thread-safe and should not be modified concurrently.
//
// # Performance Considerations
//
// The package is designed for high-performance streaming scenarios:
//
//   - Minimal allocations during encoding/decoding
//   - Efficient type switching for union types
//   - Streaming support for large message flows
//   - Validation caching where possible
//
// # Compatibility
//
// This implementation is fully compatible with the AG-UI protocol specification and
// can interoperate with other AG-UI implementations in different languages.
//
// For more information about the AG-UI protocol, see: https://docs.ag-ui.com/
package agui
