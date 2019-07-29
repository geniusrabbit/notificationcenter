package kafka

import (
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
)

type message struct {
	msg      *sarama.ConsumerMessage
	consumer *cluster.Consumer
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
	m.consumer.MarkOffset(m.msg, "")
	return nil
}
