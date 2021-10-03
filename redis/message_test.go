package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	var m message
	assert.Equal(t, ``, m.ID())
	assert.Nil(t, (*message)(nil).Body())
	assert.Nil(t, m.Context())
	assert.Equal(t, ``, string(m.Body()))
	assert.NoError(t, m.Ack())
}
