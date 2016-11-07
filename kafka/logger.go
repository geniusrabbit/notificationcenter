//
// @project geniusrabbit.com 2015
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015
//

// Config example
// config := sarama.NewConfig()
// config.Producer.RequiredAcks = sarama.WaitForLocal   // Only wait for the leader to ack
// config.Producer.Compression = sarama.CompressionNone // Compress messages Snappy
// config.Producer.Flush.Frequency = flushFrequency     // Flush batches every ...

package kafka

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

// Logger object
type Logger struct {
	sync.Mutex
	topics   []string
	producer sarama.AsyncProducer
}

// NewLogger for topics
func NewLogger(brokerList []string, topics []string, conf *sarama.Config, flushFrequency time.Duration) *Logger {
	return &Logger{
		topics:   topics,
		producer: newLogProducer(brokerList, conf, flushFrequency),
	}
}

///////////////////////////////////////////////////////////////////////////////
/// Methods
///////////////////////////////////////////////////////////////////////////////

// Close kafka producer
func (l *Logger) Close() (err error) {
	if nil != l {
		l.Lock()
		defer l.Unlock()
		if nil != l.producer {
			err = l.producer.Close()
		}
	}
	return
}

// Write binary message
func (l *Logger) Write(data []byte) error {
	if nil != data && len(data) > 0 {
		for _, topic := range l.topics {
			l.producer.Input() <- &sarama.ProducerMessage{
				Topic: topic,
				Value: byteEncoder(data),
			}
		}
	}
	return nil
}

// Send message
func (l *Logger) Send(messages ...interface{}) (err error) {
	if nil != messages {
		for _, it := range messages {
			if b, err := json.Marshal(it); nil == err {
				err = l.Write(b)
			}
			if nil != err {
				break
			}
		}
	}
	return
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

func newLogProducer(brokerList []string, config *sarama.Config, flushFrequency time.Duration) sarama.AsyncProducer {
	if nil == config {
		config = sarama.NewConfig()
	}

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start producer:", err)
	}

	// We will just log to STDOUT if we're not able to produce messages.
	// Note: messages will only be returned here after all retry attempts are exhausted.
	go func() {
		for err := range producer.Errors() {
			log.Println("Failed to write log entry:", err)
		}
	}()

	return producer
}
