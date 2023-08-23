package kafka

import (
	"context"

	"github.com/IBM/sarama"
)

type message struct {
	msg     *sarama.ConsumerMessage
	session sarama.ConsumerGroupSession
}

// Context of the message
func (m *message) Context() context.Context {
	return m.session.Context()
}

// ID returns unical message ID (depends on transport)
func (m *message) ID() string {
	return ""
}

// Body returns message data as bytes
func (m *message) Body() []byte {
	return m.msg.Value
}

// Acknowledgment of the message processing
func (m *message) Ack() error {
	m.session.MarkMessage(m.msg, "")
	return nil
}
