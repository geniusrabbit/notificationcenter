//
// @project geniusrabbit.com 2015, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015, 2019
//

// Config example
// config := sarama.NewConfig()
// config.Producer.RequiredAcks = sarama.WaitForLocal   // Only wait for the leader to ack
// config.Producer.Compression = sarama.CompressionNone // Compress messages Snappy
// config.Producer.Flush.Frequency = flushFrequency     // Flush batches every ...

package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/geniusrabbit/notificationcenter/encoder"
)

// StreamErrorHandler callback function
type StreamErrorHandler func(*sarama.ProducerError)

// StreamSuccessHandler callback function
type StreamSuccessHandler func(*sarama.ProducerMessage)

// Stream implementation of Streamer interface for the Kafka driver
type Stream struct {
	producer sarama.AsyncProducer
	topics   []string
	encoder  encoder.Encoder
}

// NewStream to the kafka with some brokers and topics for sending
func NewStream(brokerList, topics []string, options ...OptionStream) (*Stream, error) {
	conf := sarama.NewConfig()
	for _, opt := range options {
		opt(conf)
	}
	producer, err := newProducer(brokerList, conf)
	if err != nil {
		return nil, err
	}
	return &Stream{
		producer: producer,
		topics:   topics,
		encoder:  encoder.JSON,
	}, nil
}

// MustNewStream connection or panic
func MustNewStream(brokerList, topics []string, options ...OptionStream) *Stream {
	stream, err := NewStream(brokerList, topics, options...)
	if err != nil {
		panic(err)
	}
	return stream
}

///////////////////////////////////////////////////////////////////////////////
/// Methods
///////////////////////////////////////////////////////////////////////////////

// Send messages to the KAFKA stream
func (s *Stream) Send(messages ...interface{}) (err error) {
	buff := acquireBuffer()
	for _, it := range messages {
		buff.Reset()
		if err = s.encoder(it, buff); err == nil {
			err = s.sendByteMessage(buff.Bytes())
		}
		if err != nil {
			break
		}
	}
	releaseBuffer(buff)
	return
}

func (s *Stream) sendByteMessage(data []byte) error {
	for _, topic := range s.topics {
		s.producer.Input() <- &sarama.ProducerMessage{
			Topic: topic,
			Value: byteEncoder(data),
		}
	}
	return nil
}

// Processor of events
func (s *Stream) Processor(ctx context.Context, errHandler StreamErrorHandler, successHandler StreamSuccessHandler) {
	if successHandler == nil {
		successHandler = func(msg *sarama.ProducerMessage) {
			log.Printf("%s [%s] success partition=%d offset=%d\n",
				msg.Timestamp, msg.Topic, msg.Partition, msg.Offset)
		}
	}

	if errHandler == nil {
		errHandler = func(err *sarama.ProducerError) {
			log.Printf("%s [%s] error %s", err.Msg.Timestamp, err.Msg.Topic, err.Error())
		}
	}

	// Note: messages will only be returned here after all retry attempts are exhausted.
loop:
	for {
		select {
		case errMsg, ok := <-s.producer.Errors():
			if ok {
				errHandler(errMsg)
			} else {
				break loop
			}
		case msg, ok := <-s.producer.Successes():
			if ok {
				successHandler(msg)
			} else {
				break loop
			}
		case <-ctx.Done():
			break loop
		}
	}
}

// Producer interface accessor
func (s *Stream) Producer() sarama.AsyncProducer {
	return s.producer
}

// Close kafka producer
func (s *Stream) Close() error {
	return s.producer.Close()
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

func newProducer(brokerList []string, config *sarama.Config) (sarama.AsyncProducer, error) {
	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		return nil, fmt.Errorf("[kafka] failed to start producer: %v", err)
	}
	return producer, nil
}
