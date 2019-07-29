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
	"bytes"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/geniusrabbit/notificationcenter/encoder"
)

// Stream implementation of Streamer interface for the Kafka driver
type Stream struct {
	producer sarama.AsyncProducer
	topics   []string
	encoder  encoder.Encoder
}

// NewStream to the kafka with some brokers and topics for sending
func NewStream(brokerList []string, topics []string, conf *sarama.Config) (*Stream, error) {
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

func MustNewStream(brokerList []string, topics []string, conf *sarama.Config) *Stream {
	stream, err := NewStream(brokerList, topics, conf)
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
	var buff bytes.Buffer
	for _, it := range messages {
		buff.Reset()
		if err = s.encoder(it, &buff); err == nil {
			err = s.sendByteMessage(buff.Bytes())
		}
		if err != nil {
			break
		}
	}
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

// Close kafka producer
func (s *Stream) Close() error {
	return s.producer.Close()
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

func newProducer(brokerList []string, config *sarama.Config) (sarama.AsyncProducer, error) {
	if config == nil {
		config = sarama.NewConfig()
	}

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		return nil, fmt.Errorf("[kafka] failed to start producer: %v", err)
	}

	// We will just log to STDOUT if we're not able to produce messages.
	// Note: messages will only be returned here after all retry attempts are exhausted.
	go func() {
		for err := range producer.Errors() {
			log.Println("Failed to write log entry:", err)
		}
	}()

	return producer, nil
}
