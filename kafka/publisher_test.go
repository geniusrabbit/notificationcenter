package kafka

import (
	"context"
	crand "crypto/rand"
	"math/rand"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
)

type testMessage struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

func TestSending(t *testing.T) {
	defaultCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	seedBroker := sarama.NewMockBroker(t, 1)
	leader := sarama.NewMockBroker(t, 2)
	defer seedBroker.Close()
	defer leader.Close()

	// Use handler maps to always respond to requests without running out of expectations
	seedBroker.SetHandlerByMap(map[string]sarama.MockResponse{
		"ApiVersionsRequest": sarama.NewMockApiVersionsResponse(t),
		"MetadataRequest": sarama.NewMockMetadataResponse(t).
			SetBroker(leader.Addr(), leader.BrokerID()).
			SetLeader("my_topic", 0, leader.BrokerID()),
	})

	leader.SetHandlerByMap(map[string]sarama.MockResponse{
		"ApiVersionsRequest": sarama.NewMockApiVersionsResponse(t),
		"MetadataRequest": sarama.NewMockMetadataResponse(t).
			SetBroker(leader.Addr(), leader.BrokerID()).
			SetLeader("my_topic", 0, leader.BrokerID()),
		"ProduceRequest": sarama.NewMockProduceResponse(t).
			SetError("my_topic", 0, sarama.ErrNoError),
	})
	const messageCount = 10
	successCount := 0
	receiveChan := make(chan error, messageCount)

	config := sarama.NewConfig()
	config.Version = sarama.V2_6_0_0
	config.Producer.Flush.Messages = messageCount
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Retry.Backoff = 100 * time.Millisecond
	config.Net.DialTimeout = 3 * time.Second
	config.Net.ReadTimeout = 3 * time.Second
	config.Net.WriteTimeout = 3 * time.Second
	// Be conservative with metadata retries to reduce extra expectations
	config.Metadata.Retry.Max = 1

	// Connect the publisher
	publisher := MustNewPublisher(
		defaultCtx,
		WithBrokers(seedBroker.Addr()),
		WithTopics(`my_topic`),
		WithClientID("test"),
		WithSaramaConfig(config),
		WithPublisherSuccessHandler(func(*sarama.ProducerMessage) {
			select {
			case receiveChan <- nil:
			case <-defaultCtx.Done():
			}
		}),
		WithPublisherErrorHandler(func(err *sarama.ProducerError) {
			select {
			case receiveChan <- err:
			case <-defaultCtx.Done():
			}
		}),
		WithCompression(sarama.CompressionNone, 0),
		WithFlashMessages(messageCount),
		WithFlashFrequency(0),
	)

	for i := 0; i < messageCount; i++ {
		err := publisher.Publish(defaultCtx, &testMessage{ID: rand.Int(), Title: "test"})
		assert.NoError(t, err, "send message")
	}

loop:
	for i := 0; i < messageCount; i++ {
		select {
		case err := <-receiveChan:
			if err != nil {
				t.Errorf("kafka publisher test: %s", err)
			} else {
				successCount++
			}
		case <-time.After(time.Second):
			t.Errorf("Timeout waiting for msg #%d", i)
			break loop
		}
	}

	assert.NoError(t, publisher.Close(), "close connection")
	assert.Equal(t, messageCount, successCount, "not all messages are success")
}

func TestNewPublisherError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := NewPublisher(ctx)
	assert.Error(t, err)
}

func TestNewPublisherWithInvalidBroker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := NewPublisher(ctx, WithBrokers("invalid:9092"))
	assert.Error(t, err)
}

func TestNewPublisherPanic(t *testing.T) {
	defer func() {
		r := recover()
		assert.NotNil(t, r, "expected panic but none occurred")
	}()
	MustNewPublisher(context.TODO()) // No options provided, should panic
}

func strFromDict(strSize int) string {
	const dict = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz!@#$%^&*_+-=~"
	var bytes = make([]byte, strSize)
	_, _ = crand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dict[v%byte(len(dict))]
	}
	return string(bytes)
}
