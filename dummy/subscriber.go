package dummy

import (
	"context"

	nc "github.com/geniusrabbit/notificationcenter/v2"
)

// Subscriber struct
type Subscriber struct{}

// Subscribe new handler
func (e Subscriber) Subscribe(ctx context.Context, r nc.Receiver) error {
	return nil
}

// Listen process
func (e Subscriber) Listen(ctx context.Context) error {
	return nil
}

// Close proxy handler
func (e Subscriber) Close() error {
	return nil
}

var _ nc.Subscriber = Subscriber{}
