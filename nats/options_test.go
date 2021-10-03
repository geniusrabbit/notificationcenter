package nats

import (
	"context"
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

func TestOptionWithURL(t *testing.T) {
	var options Options

	WithTopics(`topic-test1`, `topic-test2`)(&options)
	assert.ElementsMatch(t, []string{`topic-test1`, `topic-test2`}, options.Topics)

	WithNatsURL(`nats://demo:4222/test?topics=topic1,topic2`)(&options)
	assert.ElementsMatch(t, []string{`topic1`, `topic2`}, options.Topics)
	assert.Equal(t, `test`, options.group())

	WithContext(context.Background())(&options)
	assert.NotNil(t, options.context(), `context`)
}
