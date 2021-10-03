package mocks

import "context"

// Message simple object for testing reasons
type Message struct {
	ctx  context.Context
	id   string
	body []byte
	err  error
}

// NewMessage object constructor
func NewMessage(ctx context.Context, id string, body []byte, ackErr error) *Message {
	return &Message{ctx: ctx, id: id, body: body, err: ackErr}
}

// Context of the message
func (m *Message) Context() context.Context {
	return m.ctx
}

// ID Unical message ID (depends on transport)
func (m *Message) ID() string {
	return m.id
}

// Body returns message data as bytes
func (m *Message) Body() []byte {
	return m.body
}

// Ack - Acknowledgment of the message processing
func (m *Message) Ack() error {
	return m.err
}
