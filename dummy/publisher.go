package dummy

import (
	"context"
	"log"
)

// Publisher dummy implementation
type Publisher struct{}

// Publish messages for dummy space
func (pub Publisher) Publish(ctx context.Context, messages ...any) error {
	for _, msg := range messages {
		log.Println("publish", msg)
	}
	return nil
}

// Close dummy publisher
func (pub Publisher) Close() error {
	return nil
}
