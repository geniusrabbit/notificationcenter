//
// @project geniusrabbit.com 2015, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015, 2019
//

package kafka

import (
	"github.com/Shopify/sarama"
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
func NewSubscriber(brokers []string, group string, topics []string, configArgs ...*cluster.Config) (*Subscriber, error) {
	var config *cluster.Config

	if len(configArgs) > 0 && configArgs[0] != nil {
		config = configArgs[0]
	} else {
		// Configure by default to receive the oldest messages
		config = cluster.NewConfig()
		config.Consumer.Return.Errors = true
		config.Group.Return.Notifications = true
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	consumer, err := cluster.NewConsumer(brokers, group, topics, config)

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
