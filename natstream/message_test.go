package natstream

import (
	"testing"

	nstream "github.com/nats-io/stan.go"
	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	var m = message{msg: &nstream.Msg{}}
	assert.Equal(t, ``, m.ID())
	assert.Nil(t, m.Context())
	assert.Nil(t, m.Body())
	assert.Panics(t, func() { _ = m.Ack() })
}
