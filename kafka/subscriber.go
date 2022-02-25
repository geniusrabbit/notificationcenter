//
// @project geniusrabbit.com 2015, 2019 - 2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015, 2019 - 2020
//

package kafka

import (
	"context"
	"time"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/geniusrabbit/notificationcenter"
)

type loggerInterface interface {
	Error(params ...any)
	Debugf(msg string, params ...any)
}

// SubscriberNotificationHandler callback function
type SubscriberNotificationHandler func(notification *cluster.Notification)

// Subscriber for kafka
type Subscriber struct {
	notificationcenter.ModelSubscriber

	topics []string

	// consumer object which receive the messages
	// consumer      *cluster.Consumer
	consumerGroup sarama.ConsumerGroup

	// notificationHandler callback
	notificationHandler SubscriberNotificationHandler

	// logger interface
	logger loggerInterface
}

// NewSubscriber connection to kafka "group" from list of topics
func NewSubscriber(options ...Option) (*Subscriber, error) {
	var opts Options
	opts.ClusterConfig = *sarama.NewConfig()
	opts.ClusterConfig.Consumer.Offsets.CommitInterval = time.Second
	for _, opt := range options {
		opt(&opts)
	}
	consumerGroup, err := sarama.NewConsumerGroup(opts.Brokers, opts.group(), opts.clusterConfig())
	if err != nil {
		return nil, err
	}
	return &Subscriber{
		ModelSubscriber: notificationcenter.ModelSubscriber{
			ErrorHandler: opts.ErrorHandler,
			PanicHandler: opts.PanicHandler,
		},
		topics:        opts.Topics,
		consumerGroup: consumerGroup,
		logger:        opts.logger(),
	}, nil
}

///////////////////////////////////////////////////////////////////////////////
/// Methods
///////////////////////////////////////////////////////////////////////////////

// Listen kafka consumer
func (s *Subscriber) Listen(ctx context.Context) error {
	for {
		if err := s.consumerGroup.Consume(ctx, s.topics, s); err != nil {
			s.processError(err)
		}
		// check if context was cancelled, signaling that the consumer should stop
		if ctx.Err() != nil {
			return nil
		}
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim.
func (s *Subscriber) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
// but before the offsets are committed for the very last time.
func (s *Subscriber) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (s *Subscriber) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/main/consumer_group.go#L27-L29
	for msg := range claim.Messages() {
		m := &message{msg: msg, session: session}
		if err := s.ProcessMessage(m); err != nil {
			s.logger.Error(err)
		}
	}
	return nil
}

// Close kafka consumer
func (s *Subscriber) Close() error {
	if err := s.consumerGroup.Close(); err != nil {
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
