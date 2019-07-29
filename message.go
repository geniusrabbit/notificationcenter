package notificationcenter

// Message describes the access methods to the message original object
type Message interface {
	// Unical message ID (depends on transport)
	ID() string

	// Body returns message data as bytes
	Body() []byte

	// Acknowledgment of the message processing
	Ack() error
}
