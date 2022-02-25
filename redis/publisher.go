package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"go.uber.org/multierr"

	nc "github.com/geniusrabbit/notificationcenter"
	"github.com/geniusrabbit/notificationcenter/encoder"
	"github.com/geniusrabbit/notificationcenter/internal/bytebuffer"
)

// Publisher for NATS queue
type Publisher struct {
	// List of channels
	channels []string

	// Connection to the redis client
	cli *redis.Client

	// Message encoder
	encoder encoder.Encoder

	// errorHandler process errors
	errorHandler nc.ErrorHandler

	// panicHandler process panics
	panicHandler nc.PanicHandler
}

// NewPublisher of the NATS stream server
func NewPublisher(options ...Option) (*Publisher, error) {
	var opts Options
	for _, opt := range options {
		opt(&opts)
	}
	return &Publisher{
		channels:     opts.channels(),
		cli:          opts.connect(),
		encoder:      opts.encoder(),
		errorHandler: opts.ErrorHandler,
		panicHandler: opts.PanicHandler,
	}, nil
}

// MustNewPublisher of the NATS stream server
func MustNewPublisher(options ...Option) *Publisher {
	stream, err := NewPublisher(options...)
	if err != nil || stream == nil {
		panic(err)
	}
	return stream
}

// Publish one or more messages to the pub-service
func (s *Publisher) Publish(ctx context.Context, messages ...any) (err error) {
	buff := bytebuffer.AcquireBuffer()
	defer func() {
		bytebuffer.ReleaseBuffer(buff)
		if s.panicHandler == nil {
			return
		}
		if rec := recover(); rec != nil {
			s.panicHandler(nil, rec)
		}
	}()
	for _, msg := range messages {
		buff.Reset()
		if err = s.encoder(msg, buff); err == nil {
			for _, channel := range s.channels {
				if pubErr := s.cli.Publish(ctx, channel, buff.Bytes()).Err(); pubErr != nil {
					err = multierr.Append(err, pubErr)
				}
			}
		}
		if err != nil {
			if s.errorHandler == nil {
				break
			}
			// Massage wasn't encoded so we can process only error
			s.errorHandler(nil, err)
			err = nil
		}
	}
	return err
}

// Close nats-stream client
func (s *Publisher) Close() error {
	return s.cli.Close()
}
