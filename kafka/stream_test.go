package kafka

import (
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
	seedBroker := sarama.NewMockBroker(t, 1)
	leader := sarama.NewMockBroker(t, 2)

	metadataResponse := new(sarama.MetadataResponse)
	metadataResponse.AddBroker(leader.Addr(), leader.BrokerID())
	metadataResponse.AddTopicPartition("my_topic", 0, leader.BrokerID(), nil, nil, nil, sarama.ErrNoError)
	seedBroker.Returns(metadataResponse)

	prodSuccess := new(sarama.ProduceResponse)
	prodSuccess.AddTopicPartition("my_topic", 0, sarama.ErrNoError)
	leader.Returns(prodSuccess)

	defer leader.Close()
	defer seedBroker.Close()

	// Connect the stream
	stream := MustNewStream([]string{seedBroker.Addr()}, []string{"my_topic"},
		func(conf *sarama.Config) {
			conf.Producer.Flush.Messages = 10
			conf.Producer.Return.Successes = true
		})

	const messageCount = 10
	successCount := 0

	for i := 0; i < messageCount; i++ {
		assert.NoError(t, stream.Send(&testMessage{ID: rand.Int(), Title: "test"}), "send message")
	}

loop:
	for i := 0; i < messageCount; i++ {
		select {
		case err := <-stream.Producer().Errors():
			t.Errorf("kafka2 stream test: %s", err)
		case _ = <-stream.Producer().Successes():
			successCount++
		case <-time.After(time.Second):
			t.Errorf("Timeout waiting for msg #%d", i)
			break loop
		}
	}

	assert.Equal(t, messageCount, successCount, "not all messages are success")
	assert.NoError(t, stream.Close(), "close connection")
}
