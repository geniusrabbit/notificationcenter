package kafka

import (
	"testing"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
)

func Test_Message(t *testing.T) {
	msg := &message{
		msg:      &sarama.ConsumerMessage{Value: []byte(`{"data": "test"}`)},
		consumer: nil,
	}
	assert.Equal(t, []byte(`{"data": "test"}`), msg.Body())
	assert.Equal(t, ``, msg.ID())
	assert.Error(t, msg.Ack())
}
