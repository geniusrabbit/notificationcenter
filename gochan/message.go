package gochan

import "context"

type message struct {
	ctx  context.Context
	data []byte
}

// Context of the message
func (m message) Context() context.Context {
	return m.ctx
}

// Unical message ID (depends on transport)
func (m message) ID() string {
	return ""
}

// Body returns message data as bytes
func (m message) Body() []byte {
	return m.data
}

// Acknowledgment of the message processing
func (m message) Ack() error {
	return nil
}
