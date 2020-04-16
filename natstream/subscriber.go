//
// @project GeniusRabbit 2018 - 2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018 - 2020
//

package natstream

import (
	"context"

	nstream "github.com/nats-io/stan.go"

	nc "github.com/geniusrabbit/notificationcenter"
)

type loggerInterface interface {
	Error(params ...interface{})
	Debugf(msg string, params ...interface{})
}

// Subscriber for NATS queue
type Subscriber struct {
	nc.ModelSubscriber

	// Additional options for subscribers
	natsSubscriptionOptions []nstream.SubscriptionOption

	// Group name
	group string

	// List of subscibed topics
	topics []string

	// Subscriptions
	natsSubscriptions []nstream.Subscription

	// connection to the NATS stream server
	conn nstream.Conn

	// Close connection event queue
	closeEvent chan bool

	// logger interface
	logger loggerInterface
}

// NewSubscriber creates new subscriber object
func NewSubscriber(options ...Option) (*Subscriber, error) {
	var opts Options
	for _, opt := range options {
		opt(&opts)
	}
	conn, err := nstream.Connect(opts.clusterID(), opts.clientID(), opts.NatsOptions...)
	if err != nil || conn == nil {
		return nil, err
	}
	sub := &Subscriber{
		ModelSubscriber: nc.ModelSubscriber{
			ErrorHandler: opts.ErrorHandler,
			PanicHandler: opts.PanicHandler,
		},
		conn:              conn,
		group:             opts.group(),
		topics:            opts.Topics,
		natsSubscriptions: opts.NatsSubscriptions,
		closeEvent:        make(chan bool, 1),
		logger:            opts.logger(),
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
		var sub nstream.Subscription
		if s.group != "" {
			sub, err = s.conn.QueueSubscribe(topic, s.group, s.message, s.natsSubscriptionOptions...)
		} else {
			sub, err = s.conn.Subscribe(topic, s.message, s.natsSubscriptionOptions...)
		}
		if err != nil {
			break
		}
		s.natsSubscriptions = append(s.natsSubscriptions, sub)
	}
	return err
}

// message execute
func (s *Subscriber) message(m *nstream.Msg) {
	if err := s.ProcessMessage(messageFromNats(m)); err != nil {
		s.logger.Error(err)
	}
}

// Close nstream client
func (s *Subscriber) Close() error {
	err := s.conn.Close()
	s.closeEvent <- true
	return err
}
