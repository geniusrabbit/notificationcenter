package simple

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Message(t *testing.T) {
	msg := &message{
		data: []byte(`{"data": "test"}`),
	}
	assert.Equal(t, []byte(`{"data": "test"}`), msg.Body())
	assert.Equal(t, ``, msg.ID())
	assert.Nil(t, msg.Ack())
}
