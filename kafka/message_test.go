package kafka

import (
	"context"
	"encoding/json"
	"math/rand"
	"sync"
	"testing"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"

	"github.com/geniusrabbit/notificationcenter/v2/encoder"
	"github.com/geniusrabbit/notificationcenter/v2/internal/bytebuffer"
)

type testGroupSession struct{}

func (*testGroupSession) Claims() map[string][]int32                                               { return nil }
func (*testGroupSession) MemberID() string                                                         { return "" }
func (*testGroupSession) GenerationID() int32                                                      { return 0 }
func (*testGroupSession) MarkOffset(topic string, partition int32, offset int64, metadata string)  {}
func (*testGroupSession) Commit()                                                                  {}
func (*testGroupSession) ResetOffset(topic string, partition int32, offset int64, metadata string) {}
func (*testGroupSession) MarkMessage(msg *sarama.ConsumerMessage, metadata string)                 {}
func (*testGroupSession) Context() context.Context                                                 { return nil }

func TestMessage(t *testing.T) {
	msg := &message{
		msg:     &sarama.ConsumerMessage{Value: []byte(`{"data": "test"}`)},
		session: &testGroupSession{},
	}
	assert.Equal(t, []byte(`{"data": "test"}`), msg.Body())
	assert.Equal(t, ``, msg.ID())
	assert.Nil(t, msg.Context())
	assert.Nil(t, msg.Ack())
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
