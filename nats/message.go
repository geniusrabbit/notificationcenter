package nats

import (
	nats "github.com/nats-io/nats.go"
)

type message nats.Msg

// Unical message ID (depends on transport)
func (m *message) ID() string {
	return ""
}

// Body returns message data as bytes
func (m *message) Body() []byte {
	return (*nats.Msg)(m).Data
}

// Acknowledgment of the message processing
func (m *message) Ack() error {
	return nil
}
