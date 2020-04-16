//
// @project geniusrabbit.com 2015, 2019 - 2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015, 2019 - 2020
//

package kafka

import (
	"context"
	"time"

	cluster "github.com/bsm/sarama-cluster"
	"github.com/geniusrabbit/notificationcenter"
)

type loggerInterface interface {
	Error(params ...interface{})
	Debugf(msg string, params ...interface{})
}

// SubscriberNotificationHandler callback function
type SubscriberNotificationHandler func(notification *cluster.Notification)

// Subscriber for kafka
type Subscriber struct {
	notificationcenter.ModelSubscriber

	// consumer object which receive the messages
	consumer *cluster.Consumer

	// notificationHandler callback
	notificationHandler SubscriberNotificationHandler

	// logger interface
	logger loggerInterface
}

// NewSubscriber connection to kafka "group" from list of topics
func NewSubscriber(options ...Option) (*Subscriber, error) {
	var opts Options
	opts.ClusterConfig = *cluster.NewConfig()
	opts.ClusterConfig.Consumer.Offsets.CommitInterval = time.Second
	for _, opt := range options {
		opt(&opts)
	}
	consumer, err := cluster.NewConsumer(opts.Brokers, opts.group(), opts.Topics, opts.clusterConfig())
	if err != nil {
		return nil, err
	}
	return &Subscriber{
		ModelSubscriber: notificationcenter.ModelSubscriber{
			ErrorHandler: opts.ErrorHandler,
			PanicHandler: opts.PanicHandler,
		},
		consumer: consumer,
		logger:   opts.logger(),
	}, nil
}

///////////////////////////////////////////////////////////////////////////////
/// Methods
///////////////////////////////////////////////////////////////////////////////

// Listen kafka consumer
func (s *Subscriber) Listen(ctx context.Context) (err error) {
loop:
	for {
		if s.consumer == nil {
			break
		}
		select {
		case msg, ok := <-s.consumer.Messages():
			if !ok {
				break loop
			}
			m := &message{msg: msg, consumer: s.consumer}
			if err := s.ProcessMessage(m); err != nil {
				s.logger.Error(err)
			}
		case err, ok := <-s.consumer.Errors():
			if !ok {
				break loop
			}
			if err != nil {
				s.processError(err)
			}
		case notification, ok := <-s.consumer.Notifications():
			if !ok {
				break loop
			}
			s.processNotification(notification)
		case <-ctx.Done():
			break loop
		}
	}
	return err
}

// Close kafka consumer
func (s *Subscriber) Close() error {
	if err := s.consumer.Close(); err != nil {
		_ = s.ModelSubscriber.Close()
		return err
	}
	return s.ModelSubscriber.Close()
}

func (s *Subscriber) processError(err error) {
	if s.ModelSubscriber.ErrorHandler != nil {
		s.ModelSubscriber.ErrorHandler(nil, err)
	} else {
		s.logger.Error(err)
	}
}

func (s *Subscriber) processNotification(notification *cluster.Notification) {
	if s.notificationHandler != nil {
		s.notificationHandler(notification)
	} else {
		s.logger.Debugf("consumer notification: %+v", notification)
	}
}
