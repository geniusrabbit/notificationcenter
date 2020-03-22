package mocks

// Message simple object for testing reasons
type Message struct {
	id   string
	body []byte
	err  error
}

// NewMessage object constructor
func NewMessage(id string, body []byte, ackErr error) *Message {
	return &Message{id: id, body: body, err: ackErr}
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
