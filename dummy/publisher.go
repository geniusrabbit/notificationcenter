package dummy

import (
	"context"
	"log"
)

// Publisher dummy implementation
type Publisher struct {
	Print bool
}

// Publish messages for dummy space
func (pub Publisher) Publish(ctx context.Context, messages ...any) error {
	if pub.Print {
		for _, msg := range messages {
			log.Println("dummy:publish", msg)
		}
	}
	return nil
}

// Close dummy publisher
func (pub Publisher) Close() error {
	return nil
}
