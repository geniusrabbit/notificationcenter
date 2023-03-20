package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type message struct {
	ctx context.Context
	msg *redis.Message
}

func messageFromRedis(ctx context.Context, msg *redis.Message) *message {
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
	if m == nil || m.msg == nil {
		return nil
	}
	return []byte(m.msg.Payload)
}

// Acknowledgment of the message processing
func (m *message) Ack() error {
	return nil
}
