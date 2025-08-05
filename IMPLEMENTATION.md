# AG-UI Go Implementation

This directory contains a complete Go implementation of the Agent User Interaction Protocol (AG-UI) based on the JSON schemas defined in the official AG-UI specification.

## Overview

The implementation provides:

- **Complete type definitions** for all AG-UI protocol components
- **JSON encoding/decoding** with proper validation
- **Streaming support** for real-time event processing
- **Type-safe interfaces** for union types (Events and Messages)
- **Comprehensive validation** based on schema constraints
- **Factory functions** for convenient object creation
- **Error handling** with specific error types
- **Full test coverage** with examples

## File Structure

```
pkg/ag-ui/
├── README.md                 # Schema documentation
├── IMPLEMENTATION.md         # This file
├── doc.go                   # Package documentation
├── types.go                 # Core types and enums
├── messages.go              # Message type definitions
├── events.go                # Event type definitions
├── codec.go                 # Encoding/decoding functionality
├── helpers.go               # Factory functions and utilities
├── codec_test.go            # Comprehensive tests
├── example_test.go          # Usage examples
├── *.schema.json           # JSON schema definitions
```

## Core Components

### 1. Enums and Constants

- **EventType**: All possible event types (17 types)
- **Role**: Message sender roles (developer, system, assistant, user, tool)
- **ToolCallType**: Tool call types (currently only "function")

### 2. Core Data Structures

- **Context**: Contextual information for agents
- **Tool**: Tool definitions with JSON schema parameters
- **ToolCall**: Tool calls made by agents
- **FunctionCall**: Function call details
- **RunAgentInput**: Complete input for running an agent
- **State**: Flexible agent state (interface{})

### 3. Message Types

All messages implement the `Message` interface:

- **DeveloperMessage**: Messages from developers
- **SystemMessage**: System messages
- **AssistantMessage**: Messages from AI assistants (can include tool calls)
- **UserMessage**: Messages from users
- **ToolMessage**: Messages from tools (responses to tool calls)

### 4. Event Types

All events implement the `Event` interface and include:

#### Lifecycle Events
- **RunStartedEvent**: Agent run started
- **RunFinishedEvent**: Agent run completed successfully
- **RunErrorEvent**: Agent run failed
- **StepStartedEvent**: Step within run started
- **StepFinishedEvent**: Step within run completed

#### Text Message Events
- **TextMessageStartEvent**: Text message started
- **TextMessageContentEvent**: Text content chunk (streaming)
- **TextMessageEndEvent**: Text message completed

#### Tool Call Events
- **ToolCallStartEvent**: Tool call started
- **ToolCallArgsEvent**: Tool call arguments chunk (streaming)
- **ToolCallEndEvent**: Tool call completed
- **ToolCallResultEvent**: Tool call result

#### State Management Events
- **StateSnapshotEvent**: Complete state snapshot
- **StateDeltaEvent**: State changes using JSON Patch
- **MessagesSnapshotEvent**: Complete message history

#### Special Events
- **RawEvent**: Pass-through for external events
- **CustomEvent**: Application-specific events

## Encoding/Decoding

### Basic Usage

```go
// Encode an event
event := agui.NewTextMessageStartEvent("msg_123")
data, err := agui.EncodeEvent(event)

// Decode an event
decoded, err := agui.DecodeEventFromBytes(data)

// Encode a message
message := agui.NewUserMessage("msg_456", "Hello!", "user")
data, err := agui.EncodeMessage(message)

// Decode a message
decoded, err := agui.DecodeMessageFromBytes(data)
```

### Streaming

```go
// Stream decoder for continuous event processing
decoder := agui.NewStreamDecoder(reader)
eventChan, errorChan := decoder.DecodeEvents()

for {
    select {
    case event := <-eventChan:
        // Process event
    case err := <-errorChan:
        // Handle error
    }
}
```

### Validation

All types include validation methods:

```go
event := &agui.RunStartedEvent{
    BaseEvent: agui.BaseEvent{Type: agui.EventTypeRunStarted},
    ThreadID:  "thread_123",
    RunID:     "run_456",
}

if err := event.Validate(); err != nil {
    // Handle validation error
}
```

## Factory Functions

Convenient factory functions with automatic timestamp setting:

```go
// Events
event := agui.NewRunStartedEvent("thread_1", "run_1")
textEvent := agui.NewTextMessageContentEvent("msg_1", "Hello")
toolEvent := agui.NewToolCallStartEvent("tool_1", "search", "msg_parent")

// Messages
userMsg := agui.NewUserMessage("msg_1", "Hello", "user_123")
assistantMsg := agui.NewAssistantMessage("msg_2", "Hi there!", "assistant", nil)

// IDs
threadID := agui.GenerateThreadID()
runID := agui.GenerateRunID()
messageID := agui.GenerateMessageID()
```

## Error Handling

Specific error types for different scenarios:

```go
var (
    ErrInvalidEventType   = fmt.Errorf("agui: invalid event type")
    ErrInvalidMessageType = fmt.Errorf("agui: invalid message type")
    ErrInvalidStructure   = fmt.Errorf("agui: invalid message structure")
    ErrUnmarshalFailed    = fmt.Errorf("agui: failed to unmarshal")
    ErrMarshalFailed      = fmt.Errorf("agui: failed to marshal")
    ErrValidationFailed   = fmt.Errorf("agui: validation failed")
)
```

## Union Type Handling

The implementation handles union types (Event and Message interfaces) using:

1. **Type probing**: Examine discriminator fields (type/role)
2. **Type switching**: Decode into specific concrete types
3. **Interface methods**: Common methods for all implementations
4. **Validation**: Type-specific validation rules

## Performance Considerations

- **Minimal allocations** during encoding/decoding
- **Efficient type switching** for union types
- **Streaming support** for large message flows
- **Raw message caching** to avoid re-parsing
- **Validation caching** where possible

## Thread Safety

- Encoding/decoding functions are thread-safe
- Individual instances are not thread-safe
- Use separate encoder/decoder instances per goroutine

## Compatibility

- Fully compatible with AG-UI protocol specification
- Interoperable with other AG-UI implementations
- JSON Schema Draft 7 compliant
- Go 1.19+ compatible

## Testing

Comprehensive test suite includes:

- **Unit tests** for all types and functions
- **Integration tests** for encoding/decoding roundtrips
- **Streaming tests** for event flow scenarios
- **Validation tests** for schema constraints
- **Example tests** demonstrating usage patterns

Run tests with:
```bash
go test -v
```

## Usage Examples

See `example_test.go` for comprehensive usage examples including:

- Basic event and message encoding/decoding
- Streaming event processing
- Tool call flows
- State management
- RunAgentInput creation and validation

## Future Enhancements

Potential improvements:

- **Code generation** from JSON schemas
- **Performance optimizations** for high-throughput scenarios
- **Additional validation** rules
- **Metrics and monitoring** integration
- **Protocol versioning** support
