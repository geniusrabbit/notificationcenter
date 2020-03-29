//
// @project GeniusRabbit 2016 - 2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 - 2020
//

package nats

import (
	"context"

	nats "github.com/nats-io/nats.go"

	nc "github.com/geniusrabbit/notificationcenter"
	"github.com/geniusrabbit/notificationcenter/encoder"
	"github.com/geniusrabbit/notificationcenter/internal/bytebuffer"
)

// Publisher provides functionality to work with NATS queue
type Publisher struct {
	// List of topics
	topics []string

	// Connection to the nats server
	conn *nats.Conn

	// Message encoder
	encoder encoder.Encoder

	// errorHandler process errors
	errorHandler nc.ErrorHandler

	// panicHandler process panics
	panicHandler nc.PanicHandler
}

// NewPublisher object
func NewPublisher(options ...Option) (*Publisher, error) {
	var opts Options
	for _, opt := range options {
		opt(&opts)
	}
	conn, err := opts.clientConn()
	if err != nil || conn == nil {
		return nil, err
	}
	return &Publisher{
		topics:       opts.Topics,
		conn:         conn,
		errorHandler: opts.ErrorHandler,
		panicHandler: opts.PanicHandler,
		encoder:      opts.encoder(),
	}, nil
}

// MustNewPublisher object
func MustNewPublisher(options ...Option) *Publisher {
	Publisher, err := NewPublisher(options...)
	if err != nil || Publisher == nil {
		panic(err)
	}
	return Publisher
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
			err = s.publishByteMessage(buff.Bytes())
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

// Write binary message
func (s *Publisher) publishByteMessage(data []byte) (err error) {
	for _, topic := range s.topics {
		if err = s.conn.Publish(topic, data); err != nil {
			break
		}
	}
	return err
}

// Close nats client
func (s *Publisher) Close() error {
	s.conn.Close()
	return nil
}
