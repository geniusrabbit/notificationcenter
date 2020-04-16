package kafka

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOption(t *testing.T) {
	var options Options
	assert.Equal(t, `default`, options.group())
	assert.NotNil(t, options.logger(), `logger`)
}

func TestOptionWithURL(t *testing.T) {
	var options Options

	WithTopics(`topic-test1`, `topic-test2`)(&options)
	assert.ElementsMatch(t, []string{`topic-test1`, `topic-test2`}, options.Topics)

	WithKafkaURL(`nats://demo1:8000,demo2:8000/test?topics=topic1,topic2`)(&options)
	assert.ElementsMatch(t, []string{`demo1:8000`, `demo2:8000`}, options.Brokers)
	assert.ElementsMatch(t, []string{`topic1`, `topic2`}, options.Topics)
	assert.Equal(t, `test`, options.group())
}
