package kafka

import (
	"encoding/json"
	"math/rand"
	"sync"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"

	"github.com/geniusrabbit/notificationcenter/encoder"
	"github.com/geniusrabbit/notificationcenter/internal/bytebuffer"
)

func TestMessage(t *testing.T) {
	msg := &message{
		msg:      &sarama.ConsumerMessage{Value: []byte(`{"data": "test"}`)},
		consumer: nil,
	}
	assert.Equal(t, []byte(`{"data": "test"}`), msg.Body())
	assert.Equal(t, ``, msg.ID())
	assert.Error(t, msg.Ack())
}

func TestAsyncEncode(t *testing.T) {
	var (
		wg     sync.WaitGroup
		msg    testMessage
		stream = make(chan *kafkaByteEncoder, 1000)
	)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < 1000; j++ {
				buff := bytebuffer.AcquireBuffer()
				message := &testMessage{ID: rand.Int(), Title: strFromDict(rand.Intn(100) + 1)}
				if err := encoder.JSON(message, buff); err == nil {
					stream <- byteEncoder(buff.Bytes())
				} else {
					t.Error(err)
				}
				bytebuffer.ReleaseBuffer(buff)
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(stream)
	}()

	for data := range stream {
		if err := json.Unmarshal(data.data, &msg); err != nil {
			t.Error(err, string(data.data))
		}
		data.Release()
	}
}
