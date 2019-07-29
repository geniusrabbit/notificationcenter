package natstream

import (
	nstream "github.com/nats-io/stan.go"
)

type message nstream.Msg

func messageFromNats(msg *nstream.Msg) *message {
	return (*message)(msg)
}

// Unical message ID (depends on transport)
func (m *message) ID() string {
	return ""
}

// Body returns message data as bytes
func (m *message) Body() []byte {
	return (*nstream.Msg)(m).Data
}

// Acknowledgment of the message processing
func (m *message) Ack() error {
	return (*nstream.Msg)(m).Ack()
}
