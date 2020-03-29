//
// @project GeniusRabbit 2016 – 2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2020
//

package nats

import (
	"context"

	nats "github.com/nats-io/nats.go"

	"github.com/geniusrabbit/notificationcenter"
)

type loggerInterface interface {
	Error(params ...interface{})
	Debugf(msg string, params ...interface{})
}

// Subscriber for NATS queue
type Subscriber struct {
	notificationcenter.ModelSubscriber

	// Group name
	group string

	// List of subscibed topics
	topics []string

	// Subscriptions
	natsSubscriptions []*nats.Subscription

	// connection to the nats server
	conn *nats.Conn

	// Close connection event queue
	closeEvent chan bool

	// logger interface
	logger loggerInterface
}

// NewSubscriber creates new subscriber object
// @url nats://login:password@hostname:port/group?topics=topic1,topic2,topicN
func NewSubscriber(options ...Option) (*Subscriber, error) {
	var opts Options
	for _, opt := range options {
		opt(&opts)
	}
	conn, err := opts.clientConn()
	if err != nil || conn == nil {
		return nil, err
	}
	sub := &Subscriber{
		ModelSubscriber: notificationcenter.ModelSubscriber{
			ErrorHandler: opts.ErrorHandler,
			PanicHandler: opts.PanicHandler,
		},
		group:      opts.group(),
		topics:     opts.Topics,
		conn:       conn,
		closeEvent: make(chan bool, 1),
		logger:     opts.logger(),
	}
	return sub, sub.subscribe()
}

// MustNewSubscriber creates new subscriber object
func MustNewSubscriber(options ...Option) *Subscriber {
	sub, err := NewSubscriber(options...)
	if err != nil || sub == nil {
		panic(err)
	}
	return sub
}

// Listen kafka consumer
func (s *Subscriber) Listen(ctx context.Context) error {
	select {
	case <-ctx.Done():
	case <-s.closeEvent:
	}
	return nil
}

func (s *Subscriber) subscribe() (err error) {
	for _, topic := range s.topics {
		var sub *nats.Subscription
		if s.group != `` {
			sub, err = s.conn.QueueSubscribe(topic, s.group, s.message)
		} else {
			sub, err = s.conn.Subscribe(topic, s.message)
		}
		if err != nil {
			break
		}
		s.natsSubscriptions = append(s.natsSubscriptions, sub)
	}
	return err
}

// message execute
func (s *Subscriber) message(m *nats.Msg) {
	if err := s.ProcessMessage((*message)(m)); err != nil {
		s.logger.Error(err)
	}
}

// Close nats client
func (s *Subscriber) Close() error {
	s.conn.Close()
	s.closeEvent <- true
	return nil
}
