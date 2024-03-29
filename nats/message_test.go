package nats

import (
	"testing"

	nats "github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	var m = message{msg: &nats.Msg{}}
	assert.Equal(t, ``, m.ID())
	assert.Nil(t, m.Body())
	assert.Nil(t, m.Context())
	assert.NoError(t, m.Ack())
}
