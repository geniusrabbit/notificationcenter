//
// @project geniusrabbit.com 2015, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015, 2019
//

package kafka

import (
	cluster "github.com/bsm/sarama-cluster"
	"github.com/geniusrabbit/notificationcenter/subscriber"
)

type loggerInterface interface {
	Error(params ...interface{})
	Debugf(msg string, params ...interface{})
}

// Subscriber for kafka
type Subscriber struct {
	subscriber.Base

	// consumer object which receive the messages
	consumer *cluster.Consumer

	// logger interface
	logger loggerInterface
}

// NewSubscriber connection to kafka "group" from list of topics
func NewSubscriber(brokers []string, group string, topics []string, options ...OptionSubscriber) (*Subscriber, error) {
	conf := cluster.NewConfig()
	for _, opt := range options {
		opt(conf)
	}
	consumer, err := cluster.NewConsumer(brokers, group, topics, conf)
	if err != nil {
		return nil, err
	}
	return &Subscriber{consumer: consumer}, nil
}

// SetLogger interface
func (s *Subscriber) SetLogger(logger loggerInterface) {
	s.logger = logger
}

///////////////////////////////////////////////////////////////////////////////
/// Methods
///////////////////////////////////////////////////////////////////////////////

// Listen kafka consumer
func (s *Subscriber) Listen() (err error) {
loop:
	for {
		if s.consumer == nil {
			break
		}
		select {
		case msg, ok := <-s.consumer.Messages():
			if ok {
				m := &message{msg: msg, consumer: s.consumer}
				if err := s.Handle(m, false); err != nil {
					s.logError(err)
				}
			} else {
				break loop
			}
		case err, ok := <-s.consumer.Errors():
			if ok {
				s.logError(err)
			} else {
				break loop
			}
		case notification, ok := <-s.consumer.Notifications():
			if ok {
				s.logNotification(notification)
			} else {
				break loop
			}
		}
	}
	return err
}

// Close kafka consumer
func (s *Subscriber) Close() (err error) {
	if err = s.consumer.Close(); err != nil {
		s.CloseAll()
		return err
	}
	return s.CloseAll()
}

func (s *Subscriber) logError(err error) {
	if s.logger != nil && err != nil {
		s.logger.Error(err)
	}
}

func (s *Subscriber) logNotification(notification *cluster.Notification) {
	if s.logger != nil && notification != nil {
		s.logger.Debugf("consumer notification: %+v", notification)
	}
}
