package once

import (
	"context"
	"io"

	nc "github.com/geniusrabbit/notificationcenter"
)

// PublisherWrapper provides additional check before send message to the stream
type PublisherWrapper struct {
	checker Checker
	pub     nc.Publisher
}

// MewPublisherWrapper with checker
func MewPublisherWrapper(pub nc.Publisher, checker Checker) *PublisherWrapper {
	if pub == nil {
		panic("publisher is required")
	}
	if checker == nil {
		panic("check is required")
	}
	return &PublisherWrapper{
		pub:     pub,
		checker: checker,
	}
}

// Publish one or more messages to the pub-service if will pass conditions
func (wr *PublisherWrapper) Publish(ctx context.Context, messages ...interface{}) error {
	for _, msg := range messages {
		if wr.checker.IsSkip(msg) {
			continue
		}
		if err := wr.pub.Publish(ctx, msg); err != nil {
			return err
		}
		if err := wr.checker.MarkAsSent(msg); err != nil {
			return err
		}
	}
	return nil
}

// Close publisher if supported
func (wr *PublisherWrapper) Close() error {
	if closer, ok := wr.pub.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
