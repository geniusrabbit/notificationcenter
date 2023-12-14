package notificationcenter

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type msg struct{ b []byte }

func (m *msg) Context() context.Context { return context.Background() }
func (m *msg) ID() string               { return "" }
func (m *msg) Body() []byte             { return m.b }
func (m *msg) Ack() error               { return nil }

func TestReceiverFunc(t *testing.T) {
	t.Run("func", func(t *testing.T) {
		r := ReceiverFrom(func() error { return nil })
		if assert.NotNil(t, r, "receiver is nil") {
			assert.NoError(t, r.Receive(&msg{}), "receiver unexpected error")
		}
	})

	t.Run("func with context", func(t *testing.T) {
		r := ReceiverFrom(func(ctx context.Context) error { return nil })
		if assert.NotNil(t, r, "receiver is nil") {
			assert.NoError(t, r.Receive(&msg{}), "receiver unexpected error")
		}
	})

	t.Run("func with message", func(t *testing.T) {
		r := ReceiverFrom(func(m Message) error { return nil })
		if assert.NotNil(t, r, "receiver is nil") {
			assert.NoError(t, r.Receive(&msg{}), "receiver unexpected error")
		}
	})

	t.Run("func custom", func(t *testing.T) {
		type msgType struct {
			P1 int `json:"p1"`
			P2 int `json:"p2"`
		}
		r := ReceiverFrom(func(msg Message, m *msgType) error {
			if msg == nil {
				return fmt.Errorf("invalid message param")
			}
			if m.P1 >= m.P2 {
				return fmt.Errorf("invalid param value")
			}
			return nil
		})
		if assert.NotNil(t, r, "receiver is nil") {
			assert.NoError(t, r.Receive(&msg{b: []byte(`{"p1":1,"p2":2}`)}), "receiver unexpected error")
		}
	})
}
