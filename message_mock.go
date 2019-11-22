package notificationcenter

// MessageMock simple object for testing reasons
type MessageMock struct {
	id   string
	body []byte
	err  error
}

// NewMessageMock object constructor
func NewMessageMock(id string, body []byte, ackErr error) *MessageMock {
	return &MessageMock{id: id, body: body, err: ackErr}
}

// ID Unical message ID (depends on transport)
func (m *MessageMock) ID() string {
	return m.id
}

// Body returns message data as bytes
func (m *MessageMock) Body() []byte {
	return m.body
}

// Ack - Acknowledgment of the message processing
func (m *MessageMock) Ack() error {
	return m.err
}
