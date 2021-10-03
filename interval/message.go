package interval

import (
	"context"
	"encoding/json"
	"errors"
)

// ErrInvalidMessageType defines the message value extracting error
var ErrInvalidMessageType = errors.New(`invalid message type`)

type messageValue interface {
	Value() interface{}
}

type message struct {
	ctx context.Context
	v   interface{}
}

// Context of the message
func (m *message) Context() context.Context {
	return m.ctx
}

// Unical message ID (depends on transport)
func (m *message) ID() string {
	return ""
}

// Body returns message data as bytes
func (m *message) Body() (d []byte) {
	switch v := m.v.(type) {
	case []byte:
		d = v
	case string:
		d = []byte(v)
	default:
		d, _ = json.Marshal(m.v)
	}
	return d
}

// Acknowledgment of the message processing
func (m *message) Ack() error {
	return nil
}

// Value returns value of the message
func (m *message) Value() interface{} {
	return m.v
}

// MessageValue take the target value from the message
func MessageValue(m interface{}) interface{} {
	switch v := m.(type) {
	case messageValue:
		return v.Value()
	}
	panic(ErrInvalidMessageType)
}
