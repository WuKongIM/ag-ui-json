package agui

import (
	"encoding/json"
	"fmt"
	"io"
)

// Predefined encoding/decoding errors
var (
	ErrInvalidEventType   = fmt.Errorf("agui: invalid event type")
	ErrInvalidMessageType = fmt.Errorf("agui: invalid message type")
	ErrInvalidStructure   = fmt.Errorf("agui: invalid message structure")
	ErrUnmarshalFailed    = fmt.Errorf("agui: failed to unmarshal")
	ErrMarshalFailed      = fmt.Errorf("agui: failed to marshal")
	ErrValidationFailed   = fmt.Errorf("agui: validation failed")
)

// EventProbe is used to determine the type of an incoming event by examining the type field.
type EventProbe struct {
	Type      EventType       `json:"type"`
	RawData   json.RawMessage `json:"-"`
	Timestamp *int64          `json:"timestamp,omitempty"`
	RawEvent  interface{}     `json:"rawEvent,omitempty"`
}

// MessageProbe is used to determine the type of an incoming message by examining the role field.
type MessageProbe struct {
	Role    Role            `json:"role"`
	RawData json.RawMessage `json:"-"`
	ID      string          `json:"id"`
	Name    string          `json:"name,omitempty"`
}

// Encoder provides functionality to encode AG-UI protocol data structures to JSON.
type Encoder struct {
	writer io.Writer
}

// NewEncoder creates a new Encoder that writes to the provided io.Writer.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{writer: w}
}

// Encode marshals and writes an AG-UI data structure to the underlying writer.
func (e *Encoder) Encode(v interface{}) error {
	// Validate the data structure if it implements the Validate method
	if validator, ok := v.(interface{ Validate() error }); ok {
		if err := validator.Validate(); err != nil {
			return fmt.Errorf("%w: %v", ErrValidationFailed, err)
		}
	}

	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrMarshalFailed, err)
	}

	_, err = e.writer.Write(data)
	if err != nil {
		return fmt.Errorf("agui: failed to write encoded data: %w", err)
	}

	return nil
}

// EncodeEvent marshals an Event to JSON bytes.
func EncodeEvent(event Event) ([]byte, error) {
	if err := event.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrValidationFailed, err)
	}

	data, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrMarshalFailed, err)
	}

	return data, nil
}

// EncodeMessage marshals a Message to JSON bytes.
func EncodeMessage(message Message) ([]byte, error) {
	if err := message.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrValidationFailed, err)
	}

	data, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrMarshalFailed, err)
	}

	return data, nil
}

// Decoder provides functionality to decode AG-UI protocol data structures from JSON.
type Decoder struct {
	decoder *json.Decoder
}

// NewDecoder creates a new Decoder that reads from the provided io.Reader.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{decoder: json.NewDecoder(r)}
}

// DecodeEvent reads and decodes a single AG-UI event from the underlying reader.
func (d *Decoder) DecodeEvent() (Event, error) {
	var rawData json.RawMessage
	if err := d.decoder.Decode(&rawData); err != nil {
		if err == io.EOF {
			return nil, err
		}
		return nil, fmt.Errorf("%w: %v", ErrUnmarshalFailed, err)
	}

	var probe EventProbe
	if err := json.Unmarshal(rawData, &probe); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnmarshalFailed, err)
	}
	probe.RawData = rawData

	// Re-decode the raw data into the specific event type
	return decodeEventFromProbe(&probe)
}

// DecodeMessage reads and decodes a single AG-UI message from the underlying reader.
func (d *Decoder) DecodeMessage() (Message, error) {
	var rawData json.RawMessage
	if err := d.decoder.Decode(&rawData); err != nil {
		if err == io.EOF {
			return nil, err
		}
		return nil, fmt.Errorf("%w: %v", ErrUnmarshalFailed, err)
	}

	var probe MessageProbe
	if err := json.Unmarshal(rawData, &probe); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnmarshalFailed, err)
	}
	probe.RawData = rawData

	// Re-decode the raw data into the specific message type
	return decodeMessageFromProbe(&probe)
}

// DecodeEventFromBytes decodes an Event from JSON bytes.
func DecodeEventFromBytes(data []byte) (Event, error) {
	var probe EventProbe
	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnmarshalFailed, err)
	}

	probe.RawData = data
	return decodeEventFromProbe(&probe)
}

// DecodeMessageFromBytes decodes a Message from JSON bytes.
func DecodeMessageFromBytes(data []byte) (Message, error) {
	var probe MessageProbe
	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnmarshalFailed, err)
	}

	probe.RawData = data
	return decodeMessageFromProbe(&probe)
}

// decodeEventFromProbe decodes an event based on the probed type.
func decodeEventFromProbe(probe *EventProbe) (Event, error) {
	var data []byte
	if probe.RawData != nil {
		data = probe.RawData
	} else {
		// If no raw data, marshal the probe back to JSON
		var err error
		data, err = json.Marshal(probe)
		if err != nil {
			return nil, fmt.Errorf("%w: failed to marshal probe: %v", ErrUnmarshalFailed, err)
		}
	}

	switch probe.Type {
	case EventTypeRunStarted:
		var event RunStartedEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, fmt.Errorf("%w: RunStartedEvent: %v", ErrUnmarshalFailed, err)
		}
		return &event, event.Validate()

	case EventTypeRunFinished:
		var event RunFinishedEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, fmt.Errorf("%w: RunFinishedEvent: %v", ErrUnmarshalFailed, err)
		}
		return &event, event.Validate()

	case EventTypeRunError:
		var event RunErrorEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, fmt.Errorf("%w: RunErrorEvent: %v", ErrUnmarshalFailed, err)
		}
		return &event, event.Validate()

	case EventTypeStepStarted:
		var event StepStartedEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, fmt.Errorf("%w: StepStartedEvent: %v", ErrUnmarshalFailed, err)
		}
		return &event, event.Validate()

	case EventTypeStepFinished:
		var event StepFinishedEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, fmt.Errorf("%w: StepFinishedEvent: %v", ErrUnmarshalFailed, err)
		}
		return &event, event.Validate()

	case EventTypeTextMessageStart:
		var event TextMessageStartEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, fmt.Errorf("%w: TextMessageStartEvent: %v", ErrUnmarshalFailed, err)
		}
		return &event, event.Validate()

	case EventTypeTextMessageContent:
		var event TextMessageContentEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, fmt.Errorf("%w: TextMessageContentEvent: %v", ErrUnmarshalFailed, err)
		}
		return &event, event.Validate()

	case EventTypeTextMessageEnd:
		var event TextMessageEndEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, fmt.Errorf("%w: TextMessageEndEvent: %v", ErrUnmarshalFailed, err)
		}
		return &event, event.Validate()

	case EventTypeToolCallStart:
		var event ToolCallStartEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, fmt.Errorf("%w: ToolCallStartEvent: %v", ErrUnmarshalFailed, err)
		}
		return &event, event.Validate()

	case EventTypeToolCallArgs:
		var event ToolCallArgsEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, fmt.Errorf("%w: ToolCallArgsEvent: %v", ErrUnmarshalFailed, err)
		}
		return &event, event.Validate()

	case EventTypeToolCallEnd:
		var event ToolCallEndEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, fmt.Errorf("%w: ToolCallEndEvent: %v", ErrUnmarshalFailed, err)
		}
		return &event, event.Validate()

	case EventTypeToolCallResult:
		var event ToolCallResultEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, fmt.Errorf("%w: ToolCallResultEvent: %v", ErrUnmarshalFailed, err)
		}
		return &event, event.Validate()

	case EventTypeStateSnapshot:
		var event StateSnapshotEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, fmt.Errorf("%w: StateSnapshotEvent: %v", ErrUnmarshalFailed, err)
		}
		return &event, event.Validate()

	case EventTypeStateDelta:
		var event StateDeltaEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, fmt.Errorf("%w: StateDeltaEvent: %v", ErrUnmarshalFailed, err)
		}
		return &event, event.Validate()

	case EventTypeMessagesSnapshot:
		var event MessagesSnapshotEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, fmt.Errorf("%w: MessagesSnapshotEvent: %v", ErrUnmarshalFailed, err)
		}
		return &event, event.Validate()

	case EventTypeRaw:
		var event RawEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, fmt.Errorf("%w: RawEvent: %v", ErrUnmarshalFailed, err)
		}
		return &event, event.Validate()

	case EventTypeCustom:
		var event CustomEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, fmt.Errorf("%w: CustomEvent: %v", ErrUnmarshalFailed, err)
		}
		return &event, event.Validate()

	default:
		return nil, fmt.Errorf("%w: unknown event type: %s", ErrInvalidEventType, probe.Type)
	}
}

// decodeMessageFromProbe decodes a message based on the probed role.
func decodeMessageFromProbe(probe *MessageProbe) (Message, error) {
	var data []byte
	if probe.RawData != nil {
		data = probe.RawData
	} else {
		// If no raw data, marshal the probe back to JSON
		var err error
		data, err = json.Marshal(probe)
		if err != nil {
			return nil, fmt.Errorf("%w: failed to marshal probe: %v", ErrUnmarshalFailed, err)
		}
	}

	switch probe.Role {
	case RoleDeveloper:
		var message DeveloperMessage
		if err := json.Unmarshal(data, &message); err != nil {
			return nil, fmt.Errorf("%w: DeveloperMessage: %v", ErrUnmarshalFailed, err)
		}
		return &message, message.Validate()

	case RoleSystem:
		var message SystemMessage
		if err := json.Unmarshal(data, &message); err != nil {
			return nil, fmt.Errorf("%w: SystemMessage: %v", ErrUnmarshalFailed, err)
		}
		return &message, message.Validate()

	case RoleAssistant:
		var message AssistantMessage
		if err := json.Unmarshal(data, &message); err != nil {
			return nil, fmt.Errorf("%w: AssistantMessage: %v", ErrUnmarshalFailed, err)
		}
		return &message, message.Validate()

	case RoleUser:
		var message UserMessage
		if err := json.Unmarshal(data, &message); err != nil {
			return nil, fmt.Errorf("%w: UserMessage: %v", ErrUnmarshalFailed, err)
		}
		return &message, message.Validate()

	case RoleTool:
		var message ToolMessage
		if err := json.Unmarshal(data, &message); err != nil {
			return nil, fmt.Errorf("%w: ToolMessage: %v", ErrUnmarshalFailed, err)
		}
		return &message, message.Validate()

	default:
		return nil, fmt.Errorf("%w: unknown message role: %s", ErrInvalidMessageType, probe.Role)
	}
}

// StreamDecoder provides functionality for decoding streaming AG-UI events.
// This is particularly useful for the event-driven architecture of AG-UI.
type StreamDecoder struct {
	decoder *json.Decoder
}

// NewStreamDecoder creates a new StreamDecoder that reads from the provided io.Reader.
func NewStreamDecoder(r io.Reader) *StreamDecoder {
	return &StreamDecoder{decoder: json.NewDecoder(r)}
}

// DecodeEvents continuously decodes events from the stream until EOF or error.
// It returns a channel of events and a channel of errors.
func (s *StreamDecoder) DecodeEvents() (<-chan Event, <-chan error) {
	eventChan := make(chan Event, 10)
	errorChan := make(chan error, 1)

	go func() {
		defer close(eventChan)
		defer close(errorChan)

		for {
			var rawData json.RawMessage
			if err := s.decoder.Decode(&rawData); err != nil {
				if err == io.EOF {
					return // Normal end of stream
				}
				errorChan <- fmt.Errorf("%w: %v", ErrUnmarshalFailed, err)
				return
			}

			var probe EventProbe
			if err := json.Unmarshal(rawData, &probe); err != nil {
				errorChan <- fmt.Errorf("%w: %v", ErrUnmarshalFailed, err)
				return
			}
			probe.RawData = rawData

			event, err := decodeEventFromProbe(&probe)
			if err != nil {
				errorChan <- err
				return
			}

			eventChan <- event
		}
	}()

	return eventChan, errorChan
}

// DecodeMessages continuously decodes messages from the stream until EOF or error.
// It returns a channel of messages and a channel of errors.
func (s *StreamDecoder) DecodeMessages() (<-chan Message, <-chan error) {
	messageChan := make(chan Message, 10)
	errorChan := make(chan error, 1)

	go func() {
		defer close(messageChan)
		defer close(errorChan)

		for {
			var rawData json.RawMessage
			if err := s.decoder.Decode(&rawData); err != nil {
				if err == io.EOF {
					return // Normal end of stream
				}
				errorChan <- fmt.Errorf("%w: %v", ErrUnmarshalFailed, err)
				return
			}

			var probe MessageProbe
			if err := json.Unmarshal(rawData, &probe); err != nil {
				errorChan <- fmt.Errorf("%w: %v", ErrUnmarshalFailed, err)
				return
			}
			probe.RawData = rawData

			message, err := decodeMessageFromProbe(&probe)
			if err != nil {
				errorChan <- err
				return
			}

			messageChan <- message
		}
	}()

	return messageChan, errorChan
}
