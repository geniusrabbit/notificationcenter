package natstream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	var m message
	assert.Equal(t, ``, m.ID())
	assert.Nil(t, m.Body())
	assert.Panics(t, func() { _ = m.Ack() })
}
