package kafka

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
)

type testMessage struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

func Test_Sending(t *testing.T) {
	defaultCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	seedBroker := sarama.NewMockBroker(t, 1)
	leader := sarama.NewMockBroker(t, 2)
	defer leader.Close()
	defer seedBroker.Close()

	metadataResponse := new(sarama.MetadataResponse)
	metadataResponse.AddBroker(leader.Addr(), leader.BrokerID())
	metadataResponse.AddTopicPartition("my_topic", 0, leader.BrokerID(), nil, nil, nil, sarama.ErrNoError)
	seedBroker.Returns(metadataResponse)

	prodSuccess := new(sarama.ProduceResponse)
	prodSuccess.AddTopicPartition("my_topic", 0, sarama.ErrNoError)
	leader.Returns(prodSuccess)

	const messageCount = 10
	successCount := 0
	receiveChan := make(chan error, messageCount)

	// Connect the publisher
	publisher := MustNewPublisher(
		defaultCtx, []string{seedBroker.Addr()}, []string{"my_topic"},
		WithClientID("test"),
		WithPublisherSuccessHandler(func(*sarama.ProducerMessage) { receiveChan <- nil }),
		WithPublisherErrorHandler(func(err *sarama.ProducerError) { receiveChan <- err }),
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

func Test_NewPublisherError(t *testing.T) {
	_, err := NewPublisher(context.TODO(), nil, nil)
	assert.Error(t, err)
}

func Test_NewPublisherPanic(t *testing.T) {
	defer func() {
		assert.True(t, recover() != nil)
	}()
	MustNewPublisher(context.TODO(), nil, nil)
	time.Sleep(time.Millisecond * 100)
}

func strFromDict(strSize int) string {
	const dict = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz!@#$%^&*_+-=~"
	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dict[v%byte(len(dict))]
	}
	return string(bytes)
}
