package kafka

import (
	"encoding/json"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/geniusrabbit/notificationcenter/encoder"
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

func Test_asyncEncode(t *testing.T) {
	var (
		wg     sync.WaitGroup
		msg    testMessage
		stream = make(chan []byte, 1000)
	)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < 1000; j++ {
				buff := acquireBuffer()
				message := &testMessage{ID: rand.Int(), Title: strFromDict(rand.Intn(100) + 1)}
				if err := encoder.JSON(message, buff); err == nil {
					stream <- byteEncoder(buff.Bytes())
				} else {
					t.Error(err)
				}
				releaseBuffer(buff)
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(stream)
	}()

	for data := range stream {
		if err := json.Unmarshal(data, &msg); err != nil {
			t.Error(err, string(data))
		}
		(kafkaByteEncoder)(data).Release()
	}
}

func Test_NewStreamError(t *testing.T) {
	_, err := NewStream(nil, nil)
	assert.Error(t, err)
}

func Test_NewStreamPanic(t *testing.T) {
	defer func() {
		assert.True(t, recover() != nil)
	}()
	MustNewStream(nil, nil)
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
