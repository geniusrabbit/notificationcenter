package notificationcenter

import "context"

// Message describes the access methods to the message original object
type Message interface {
	// Context of the message
	Context() context.Context

	// Unical message ID (depends on transport)
	ID() string

	// Body returns message data as bytes
	Body() []byte

	// Acknowledgment of the message processing
	Ack() error
}
