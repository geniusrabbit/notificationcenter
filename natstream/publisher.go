//
// @project GeniusRabbit 2018 - 2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018 - 2020
//

package natstream

import (
	"context"

	nstream "github.com/nats-io/stan.go"

	nc "github.com/geniusrabbit/notificationcenter"
	"github.com/geniusrabbit/notificationcenter/encoder"
	"github.com/geniusrabbit/notificationcenter/internal/bytebuffer"
)

// Publisher for NATS queue
type Publisher struct {
	// List of topics
	topic string

	// Connection to the nats server
	conn nstream.Conn

	// Message encoder
	encoder encoder.Encoder

	// errorHandler process errors
	errorHandler nc.ErrorHandler

	// panicHandler process panics
	panicHandler nc.PanicHandler
}

// NewPublisher of the NATS stream server
func NewPublisher(topic string, options ...Option) (*Publisher, error) {
	var opts Options
	for _, opt := range options {
		opt(&opts)
	}
	conn, err := nstream.Connect(opts.clusterID(), opts.clientID(), opts.NatsOptions...)
	if err != nil || conn == nil {
		return nil, err
	}
	return &Publisher{
		topic:        topic,
		conn:         conn,
		encoder:      opts.encoder(),
		errorHandler: opts.ErrorHandler,
		panicHandler: opts.PanicHandler,
	}, nil
}

// MustNewPublisher of the NATS stream server
func MustNewPublisher(topic string, options ...Option) *Publisher {
	stream, err := NewPublisher(topic, options...)
	if err != nil || stream == nil {
		panic(err)
	}
	return stream
}

// Publish one or more messages to the pub-service
func (s *Publisher) Publish(ctx context.Context, messages ...interface{}) (err error) {
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
			err = s.conn.Publish(s.topic, buff.Bytes())
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
	s.conn.Close()
	return nil
}
