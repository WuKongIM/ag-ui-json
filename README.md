# AG-UI Go Package

A complete Go implementation of the Agent User Interaction Protocol (AG-UI) for building AI-powered applications with streaming conversations, tool calls, and state management.

## Quick Start

### Installation

```bash
go get github.com/WuKongIM/WuKongIM/pkg/ag-ui
```

### Basic Usage

```go
package main

import (
    "fmt"
    "log"

    agui "github.com/WuKongIM/WuKongIM/pkg/ag-ui"
)

func main() {
    // Create a user message
    message := agui.NewUserMessage("msg_1", "Hello, how can you help me?", "user_123")

    // Encode to JSON
    data, err := agui.EncodeMessage(message)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Encoded message: %s\n", string(data))

    // Decode back
    decoded, err := agui.DecodeMessageFromBytes(data)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Message from %s: %s\n", decoded.GetRole(), decoded.GetID())
}
```

## Core Concepts

### Events

Events represent real-time interactions in AG-UI. All events implement the `Event` interface:

```go
// Create lifecycle events
runStarted := agui.NewRunStartedEvent("thread_1", "run_1")
runFinished := agui.NewRunFinishedEvent("thread_1", "run_1", nil)

// Create text message events for streaming
textStart := agui.NewTextMessageStartEvent("msg_1")
textContent := agui.NewTextMessageContentEvent("msg_1", "Hello")
textEnd := agui.NewTextMessageEndEvent("msg_1")

// Encode events
data, err := agui.EncodeEvent(textContent)
```

### Messages

Messages represent conversation history. All messages implement the `Message` interface:

```go
// Different message types
systemMsg := agui.NewSystemMessage("msg_1", "You are a helpful assistant.", "")
userMsg := agui.NewUserMessage("msg_2", "What's the weather like?", "user_123")
assistantMsg := agui.NewAssistantMessage("msg_3", "Let me check that for you.", "assistant", nil)

// Assistant message with tool calls
toolCall := agui.ToolCall{
    ID:   "tool_1",
    Type: agui.ToolCallTypeFunction,
    Function: agui.FunctionCall{
        Name:      "get_weather",
        Arguments: `{"location": "New York"}`,
    },
}
assistantWithTools := agui.NewAssistantMessage("msg_4", "", "assistant", []agui.ToolCall{toolCall})
```

### Tool Calls

Handle tool execution with streaming events:

```go
// Tool call flow
toolStart := agui.NewToolCallStartEvent("tool_1", "get_weather", "msg_4")
toolArgs := agui.NewToolCallArgsEvent("tool_1", `{"location": "New York"}`)
toolEnd := agui.NewToolCallEndEvent("tool_1")
toolResult := agui.NewToolCallResultEvent("msg_5", "tool_1", "Sunny, 72°F")
```

## Streaming Events

### Basic Streaming

```go
package main

import (
    "bytes"
    "fmt"
    "log"

    agui "github.com/WuKongIM/WuKongIM/pkg/ag-ui"
)

func streamingExample() {
    // Create event stream
    events := []agui.Event{
        agui.NewRunStartedEvent("thread_1", "run_1"),
        agui.NewTextMessageStartEvent("msg_1"),
        agui.NewTextMessageContentEvent("msg_1", "Hello"),
        agui.NewTextMessageContentEvent("msg_1", " world!"),
        agui.NewTextMessageEndEvent("msg_1"),
        agui.NewRunFinishedEvent("thread_1", "run_1", nil),
    }

    // Encode to stream
    var buf bytes.Buffer
    for _, event := range events {
        data, err := agui.EncodeEvent(event)
        if err != nil {
            log.Fatal(err)
        }
        buf.Write(data)
        buf.WriteString("\n")
    }

    // Decode stream
    decoder := agui.NewStreamDecoder(&buf)
    eventChan, errorChan := decoder.DecodeEvents()

    for {
        select {
        case event, ok := <-eventChan:
            if !ok {
                return
            }
            fmt.Printf("Received event: %s\n", event.GetType())

        case err := <-errorChan:
            if err != nil {
                log.Printf("Stream error: %v", err)
                return
            }
        }
    }
}
```

### Real-time Message Building

```go
func buildStreamingMessage() {
    var messageContent strings.Builder

    decoder := agui.NewStreamDecoder(reader)
    eventChan, errorChan := decoder.DecodeEvents()

    for {
        select {
        case event := <-eventChan:
            switch e := event.(type) {
            case *agui.TextMessageStartEvent:
                fmt.Printf("Message started: %s\n", e.MessageID)
                messageContent.Reset()

            case *agui.TextMessageContentEvent:
                messageContent.WriteString(e.Delta)
                fmt.Printf("Current content: %s\n", messageContent.String())

            case *agui.TextMessageEndEvent:
                fmt.Printf("Final message: %s\n", messageContent.String())
                return
            }

        case err := <-errorChan:
            if err != nil {
                log.Fatal(err)
            }
        }
    }
}
```

## Agent Input & Validation

### Creating RunAgentInput

```go
func createAgentInput() *agui.RunAgentInput {
    // Define tools
    searchTool := agui.Tool{
        Name:        "search",
        Description: "Search for information",
        Parameters: map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "query": map[string]interface{}{
                    "type":        "string",
                    "description": "Search query",
                },
            },
            "required": []string{"query"},
        },
    }

    // Create messages
    messages := []agui.Message{
        agui.NewSystemMessage("msg_1", "You are a helpful assistant.", ""),
        agui.NewUserMessage("msg_2", "Search for Go tutorials", "user_123"),
    }

    // Create context
    context := []agui.Context{
        {
            Description: "User preferences",
            Value:       "Prefers concise explanations",
        },
    }

    // Create complete input
    input := &agui.RunAgentInput{
        ThreadID:       agui.GenerateThreadID(),
        RunID:          agui.GenerateRunID(),
        State:          map[string]interface{}{"step": 1},
        Messages:       messages,
        Tools:          []agui.Tool{searchTool},
        Context:        context,
        ForwardedProps: map[string]interface{}{"version": "1.0"},
    }

    // Validate input
    if err := input.Validate(); err != nil {
        log.Fatalf("Invalid input: %v", err)
    }

    return input
}
```

### Validation Examples

```go
func validationExamples() {
    // Valid event
    event := &agui.RunStartedEvent{
        BaseEvent: agui.BaseEvent{Type: agui.EventTypeRunStarted},
        ThreadID:  "thread_123",
        RunID:     "run_456",
    }

    if err := event.Validate(); err != nil {
        fmt.Printf("Validation failed: %v\n", err)
    } else {
        fmt.Println("Event is valid")
    }

    // Invalid event (missing required field)
    invalidEvent := &agui.TextMessageContentEvent{
        BaseEvent: agui.BaseEvent{Type: agui.EventTypeTextMessageContent},
        MessageID: "msg_123",
        Delta:     "", // Empty delta is invalid
    }

    if err := invalidEvent.Validate(); err != nil {
        fmt.Printf("Expected validation error: %v\n", err)
    }
}
```

## State Management

### State Snapshots and Deltas

```go
func stateManagement() {
    // Initial state
    initialState := map[string]interface{}{
        "user_id":     "user_123",
        "step_count":  1,
        "preferences": map[string]string{"theme": "dark"},
    }

    // Create state snapshot
    snapshot := agui.NewStateSnapshotEvent(initialState)

    // Create state delta (JSON Patch operations)
    delta := []interface{}{
        map[string]interface{}{
            "op":    "replace",
            "path":  "/step_count",
            "value": 2,
        },
        map[string]interface{}{
            "op":    "replace",
            "path":  "/preferences/theme",
            "value": "light",
        },
    }

    deltaEvent := agui.NewStateDeltaEvent(delta)

    // Encode and send
    snapshotData, _ := agui.EncodeEvent(snapshot)
    deltaData, _ := agui.EncodeEvent(deltaEvent)

    fmt.Printf("State snapshot: %s\n", string(snapshotData))
    fmt.Printf("State delta: %s\n", string(deltaData))
}
```

## API Reference

### Core Types

- **EventType**: Enum for all event types (17 types)
- **Role**: Enum for message roles (developer, system, assistant, user, tool)
- **ToolCallType**: Enum for tool call types (function)

### Main Interfaces

```go
type Event interface {
    GetType() EventType
    GetTimestamp() *int64
    GetRawEvent() interface{}
    Validate() error
    EventTypeName() string
}

type Message interface {
    GetID() string
    GetRole() Role
    GetName() string
    Validate() error
    MessageType() string
}
```

### Encoding Functions

```go
// Event encoding
func EncodeEvent(event Event) ([]byte, error)
func DecodeEventFromBytes(data []byte) (Event, error)

// Message encoding
func EncodeMessage(message Message) ([]byte, error)
func DecodeMessageFromBytes(data []byte) (Message, error)

// Stream processing
func NewStreamDecoder(r io.Reader) *StreamDecoder
func (s *StreamDecoder) DecodeEvents() (<-chan Event, <-chan error)
func (s *StreamDecoder) DecodeMessages() (<-chan Message, <-chan error)
```

### Factory Functions

```go
// Event factories
func NewRunStartedEvent(threadID, runID string) *RunStartedEvent
func NewTextMessageContentEvent(messageID, delta string) *TextMessageContentEvent
func NewToolCallStartEvent(toolCallID, toolCallName, parentMessageID string) *ToolCallStartEvent

// Message factories
func NewUserMessage(id, content, name string) *UserMessage
func NewAssistantMessage(id, content, name string, toolCalls []ToolCall) *AssistantMessage

// ID generators
func GenerateMessageID() string
func GenerateRunID() string
func GenerateThreadID() string
func GenerateToolCallID() string
```

## Common Patterns

### Complete Conversation Flow

```go
func conversationFlow() {
    // 1. Start agent run
    runStarted := agui.NewRunStartedEvent("thread_1", "run_1")

    // 2. Stream assistant response
    msgStart := agui.NewTextMessageStartEvent("msg_1")
    msgContent1 := agui.NewTextMessageContentEvent("msg_1", "I'll help you with that. ")
    msgContent2 := agui.NewTextMessageContentEvent("msg_1", "Let me search for information.")
    msgEnd := agui.NewTextMessageEndEvent("msg_1")

    // 3. Make tool call
    toolStart := agui.NewToolCallStartEvent("tool_1", "search", "msg_1")
    toolArgs := agui.NewToolCallArgsEvent("tool_1", `{"query": "Go tutorials"}`)
    toolEnd := agui.NewToolCallEndEvent("tool_1")

    // 4. Tool result
    toolResult := agui.NewToolCallResultEvent("msg_2", "tool_1", "Found 10 Go tutorials")

    // 5. Final response
    finalStart := agui.NewTextMessageStartEvent("msg_3")
    finalContent := agui.NewTextMessageContentEvent("msg_3", "Here are some great Go tutorials...")
    finalEnd := agui.NewTextMessageEndEvent("msg_3")

    // 6. Finish run
    runFinished := agui.NewRunFinishedEvent("thread_1", "run_1", map[string]int{"tutorials_found": 10})

    events := []agui.Event{
        runStarted, msgStart, msgContent1, msgContent2, msgEnd,
        toolStart, toolArgs, toolEnd, toolResult,
        finalStart, finalContent, finalEnd, runFinished,
    }

    // Process events
    for _, event := range events {
        data, _ := agui.EncodeEvent(event)
        fmt.Printf("Event: %s\n", string(data))
    }
}
```

### Error Handling

```go
func errorHandling() {
    // Encoding with validation
    event := agui.NewTextMessageContentEvent("msg_1", "Hello")
    data, err := agui.EncodeEvent(event)
    if err != nil {
        switch {
        case errors.Is(err, agui.ErrValidationFailed):
            log.Printf("Validation error: %v", err)
        case errors.Is(err, agui.ErrMarshalFailed):
            log.Printf("Encoding error: %v", err)
        default:
            log.Printf("Unknown error: %v", err)
        }
        return
    }

    // Decoding with error handling
    decoded, err := agui.DecodeEventFromBytes(data)
    if err != nil {
        switch {
        case errors.Is(err, agui.ErrUnmarshalFailed):
            log.Printf("Decoding error: %v", err)
        case errors.Is(err, agui.ErrInvalidEventType):
            log.Printf("Invalid event type: %v", err)
        default:
            log.Printf("Unknown error: %v", err)
        }
        return
    }

    fmt.Printf("Successfully processed event: %s\n", decoded.GetType())
}
```

## Best Practices

### Performance Tips

- **Reuse decoders**: Create one `StreamDecoder` per connection
- **Batch operations**: Process multiple events together when possible
- **Validate early**: Use validation methods before encoding
- **Handle errors**: Always check encoding/decoding errors

### Memory Management

```go
// Good: Reuse decoder
decoder := agui.NewStreamDecoder(conn)
eventChan, errorChan := decoder.DecodeEvents()

// Good: Process in batches
var events []agui.Event
for i := 0; i < batchSize; i++ {
    select {
    case event := <-eventChan:
        events = append(events, event)
    case <-time.After(timeout):
        break
    }
}
processBatch(events)
```

### Thread Safety

```go
// Safe: Each goroutine has its own decoder
func handleConnection(conn net.Conn) {
    decoder := agui.NewStreamDecoder(conn)
    eventChan, errorChan := decoder.DecodeEvents()

    for {
        select {
        case event := <-eventChan:
            // Process event safely
        case err := <-errorChan:
            // Handle error
            return
        }
    }
}
```

## Troubleshooting

### Common Issues

**1. Validation Errors**
```go
// Problem: Empty required fields
event := &agui.TextMessageContentEvent{
    BaseEvent: agui.BaseEvent{Type: agui.EventTypeTextMessageContent},
    MessageID: "msg_1",
    Delta:     "", // This will fail validation
}

// Solution: Ensure all required fields are set
event.Delta = "Hello world"
```

**2. Type Mismatches**
```go
// Problem: Wrong event type
event := &agui.RunStartedEvent{
    BaseEvent: agui.BaseEvent{Type: agui.EventTypeRunFinished}, // Wrong type
    ThreadID:  "thread_1",
    RunID:     "run_1",
}

// Solution: Use correct type
event.BaseEvent.Type = agui.EventTypeRunStarted
```

**3. JSON Parsing Errors**
```go
// Problem: Invalid JSON in function arguments
toolCall := agui.ToolCall{
    Function: agui.FunctionCall{
        Name:      "search",
        Arguments: `{invalid json}`, // This will fail validation
    },
}

// Solution: Use valid JSON
toolCall.Function.Arguments = `{"query": "search term"}`
```

### Debug Tips

- Use `Validate()` methods to check data before encoding
- Check error types with `errors.Is()` for specific handling
- Enable verbose logging for stream processing
- Test with small data sets first

## Testing

Run the test suite:

```bash
cd pkg/ag-ui
go test -v
```

Run specific tests:

```bash
go test -v -run TestEventEncoding
go test -v -run TestStreamDecoding
```

## Contributing

1. Follow Go conventions and best practices
2. Add tests for new features
3. Update documentation for API changes
4. Validate against AG-UI protocol specification

## License

This implementation follows the AG-UI protocol specification and is compatible with other AG-UI implementations.

For more information about the AG-UI protocol: https://docs.ag-ui.com/
```
