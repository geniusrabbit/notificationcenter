//
// @project geniusrabbit.com 2015
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015
//
//
// ./kafka-topics.sh --zookeeper localhost --create --partitions 1 --replication-factor 1 --topic KTEST --config cleanup.policy=compact
//
// ./kafka-topics.sh --zookeeper localhost --describe --topic KTEST
// Topic:KTEST PartitionCount:1  ReplicationFactor:1 Configs:cleanup.policy=compact
//   Topic: KTEST  Partition: 0  Leader: 0 Replicas: 0 Isr: 0
//
// ./kafka-topics.sh --zookeeper localhost --alter --partitions 3 --topic KTEST --config cleanup.policy=compact
//
//  ./kafka-topics.sh --zookeeper localhost --describe --topic KTEST
// Topic:KTEST PartitionCount:3  ReplicationFactor:1 Configs:
//   Topic: KTEST  Partition: 0  Leader: 0 Replicas: 0 Isr: 0
//   Topic: KTEST  Partition: 1  Leader: 0 Replicas: 0 Isr: 0
//   Topic: KTEST  Partition: 2  Leader: 0 Replicas: 0 Isr: 0

package kafka

import (
	"log"
	"time"

	"github.com/Shopify/sarama"
	"github.com/geniusrabbit/notificationcenter/subscriber"
	"github.com/wvanbergen/kafka/consumergroup"
)

// Subscriber for kafka
type Subscriber struct {
	subscriber.Base
	consumer *consumergroup.ConsumerGroup
}

// NewSubscriber for kafka
func NewSubscriber(addrs []string, group string, topics []string) (*Subscriber, error) {
	config := consumergroup.NewConfig()
	config.Offsets.Initial = sarama.OffsetNewest
	config.Offsets.ProcessingTimeout = 1 * time.Second
	config.Offsets.CommitInterval = 2 * time.Second
	return NewSubscriberByConfig(addrs, config, group, topics)
}

// NewSubscriberByConfig for kafka
func NewSubscriberByConfig(addrs []string, config *consumergroup.Config, group string, topics []string) (*Subscriber, error) {
	consumer, err := consumergroup.JoinConsumerGroup(group, topics, addrs, config)
	if err != nil {
		return nil, err
	}
	return &Subscriber{consumer: consumer}, nil
}

///////////////////////////////////////////////////////////////////////////////
/// Methods
///////////////////////////////////////////////////////////////////////////////

// Listen kafka consumer
func (s *Subscriber) Listen() (err error) {
	for {
		select {
		case err := <-s.consumer.Errors():
			log.Println(err)
		case m := <-s.consumer.Messages():
			if nil != s {
				if err = s.Handle(m, false); nil != err {
					break
				}
				s.consumer.CommitUpto(m)
			}
		}
	}
	return
}

// Close kafka consumer
func (s *Subscriber) Close() (err error) {
	if nil != s.consumer {
		err = s.consumer.Close()
		s.consumer = nil
		s.CloseAll()
	}
	return
}
