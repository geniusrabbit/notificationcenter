package nats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOption(t *testing.T) {
	var options Options
	conn, err := options.clientConn()
	assert.NotNil(t, options.encoder(), `encoder`)
	assert.Equal(t, `default`, options.group())
	assert.NotNil(t, options.logger(), `logger`)
	assert.Nil(t, conn)
	assert.Error(t, err)
}
