package natstream

import (
	"context"

	nstream "github.com/nats-io/stan.go"
)

type message struct {
	ctx context.Context
	msg *nstream.Msg
}

func messageFromNats(ctx context.Context, msg *nstream.Msg) *message {
	return &message{ctx: ctx, msg: msg}
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
func (m *message) Body() []byte {
	return m.msg.Data
}

// Acknowledgment of the message processing
func (m *message) Ack() error {
	return m.msg.Ack()
}
