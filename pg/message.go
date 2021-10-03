package pg

import (
	"context"
	"fmt"

	"github.com/lib/pq"
)

type message struct {
	ctx    context.Context
	notify *pq.Notification
}

func fromPgNotify(ctx context.Context, notify *pq.Notification) *message {
	return &message{ctx: ctx, notify: notify}
}

// Context of the message
func (m *message) Context() context.Context {
	return m.ctx
}

// Unical message ID (depends on transport)
func (m *message) ID() string {
	n := m.Notification()
	return fmt.Sprintf("%s-%d", n.Channel, n.BePid)
}

// Body returns message data as bytes
func (m *message) Body() []byte {
	return []byte(m.Notification().Extra)
}

// Acknowledgment of the message processing
func (m *message) Ack() error {
	return nil
}

// Notification type of message
func (m *message) Notification() *pq.Notification {
	return m.notify
}
