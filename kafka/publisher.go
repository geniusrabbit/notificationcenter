//
// @project geniusrabbit.com 2015, 2019 - 2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015, 2019 - 2020
//

// Config example
// config := sarama.NewConfig()
// config.Producer.RequiredAcks = sarama.WaitForLocal   // Only wait for the leader to ack
// config.Producer.Compression = sarama.CompressionNone // Compress messages Snappy
// config.Producer.Flush.Frequency = flushFrequency     // Flush batches every ...

package kafka

import (
	"context"

	"github.com/Shopify/sarama"

	nc "github.com/geniusrabbit/notificationcenter"
	"github.com/geniusrabbit/notificationcenter/encoder"
	"github.com/geniusrabbit/notificationcenter/internal/bytebuffer"
)

// PublisherErrorHandler callback function
type PublisherErrorHandler func(*sarama.ProducerError)

// PublisherSuccessHandler callback function
type PublisherSuccessHandler func(*sarama.ProducerMessage)

// Publisher implementation of Publisher interface for the Kafka driver
type Publisher struct {
	producer sarama.AsyncProducer

	// List of the topic names
	topics []string

	// Message encoder
	encoder encoder.Encoder

	// errorHandler process errors
	errorHandler nc.ErrorHandler

	// panicHandler process panics
	panicHandler nc.PanicHandler

	// publisherErrorHandler provides handler of message send errors
	publisherErrorHandler PublisherErrorHandler

	// successHandler provides handler of message send success
	successHandler PublisherSuccessHandler
}

// NewPublisher to the kafka with some brokers and topics for sending
func NewPublisher(ctx context.Context, options ...Option) (*Publisher, error) {
	var opts Options
	opts.ClusterConfig = *sarama.NewConfig()
	for _, opt := range options {
		opt(&opts)
	}
	producer, err := sarama.NewAsyncProducer(opts.Brokers, opts.saramaConfig())
	if err != nil {
		return nil, err
	}
	pub := &Publisher{
		producer:              producer,
		topics:                opts.Topics,
		encoder:               opts.encoder(),
		errorHandler:          opts.ErrorHandler,
		panicHandler:          opts.PanicHandler,
		publisherErrorHandler: opts.PublisherErrorHandler,
		successHandler:        opts.PublisherSuccessHandler,
	}
	if ret := opts.clusterConfig().Producer.Return; ret.Errors || ret.Successes {
		go pub.listen(ctx)
	}
	return pub, nil
}

// MustNewPublisher connection or panic
func MustNewPublisher(ctx context.Context, options ...Option) *Publisher {
	publisher, err := NewPublisher(ctx, options...)
	if err != nil {
		panic(err)
	}
	return publisher
}

///////////////////////////////////////////////////////////////////////////////
/// Methods
///////////////////////////////////////////////////////////////////////////////

// Publish one or more messages to the pub-service
func (p *Publisher) Publish(ctx context.Context, messages ...any) (err error) {
	buff := bytebuffer.AcquireBuffer()
	defer func() {
		bytebuffer.ReleaseBuffer(buff)
		if p.panicHandler == nil {
			return
		}
		if rec := recover(); rec != nil {
			p.panicHandler(nil, rec)
		}
	}()
	for _, msg := range messages {
		buff.Reset()
		if err = p.encoder(msg, buff); err == nil {
			err = p.publishByteMessage(buff.Bytes())
		}
		if err != nil {
			if p.errorHandler == nil {
				break
			}
			// Massage wasn't encoded so we can process only error
			p.errorHandler(nil, err)
			err = nil
		}
	}
	return err
}

func (p *Publisher) publishByteMessage(data []byte) error {
	for _, topic := range p.topics {
		p.producer.Input() <- &sarama.ProducerMessage{
			Topic: topic,
			Value: byteEncoder(data),
		}
	}
	return nil
}

func (p *Publisher) listen(ctx context.Context) {
loop:
	for {
		select {
		case errMsg, ok := <-p.producer.Errors():
			if !ok {
				break loop
			}
			p.publisherErrorHandler(errMsg)
		case msg, ok := <-p.producer.Successes():
			if !ok {
				break loop
			}
			if p.successHandler != nil {
				p.successHandler(msg)
			}
			// Now we can to release the kafka buffer :)
			msg.Value.(*kafkaByteEncoder).Release()
		case <-ctx.Done():
			break loop
		}
	}
}

// Close kafka producer
func (p *Publisher) Close() error {
	return p.producer.Close()
}
