package redis

import (
	"testing"

	nc "github.com/geniusrabbit/notificationcenter"
	"github.com/geniusrabbit/notificationcenter/encoder"
	"github.com/geniusrabbit/notificationcenter/internal/logger"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestOption(t *testing.T) {
	var options Options
	assert.NotNil(t, options.encoder(), `encoder`)
	assert.ElementsMatch(t, []string{`default`}, options.channels())
	assert.NotNil(t, options.logger(), `logger`)
}

func TestOptionWithURL(t *testing.T) {
	var options Options

	WithTopics(`topic-test1`, `topic-test2`)(&options)
	assert.ElementsMatch(t, []string{`topic-test1`, `topic-test2`}, options.Channels)

	WithRedisOpts(&redis.Options{})(&options)
	assert.NotNil(t, options.RedisOptions, "invalid options")
	options.RedisOptions = nil

	WithRedisURL(`redis://demo:4222/1?topics=topic1,topic2`)(&options)
	assert.ElementsMatch(t, []string{`topic1`, `topic2`}, options.Channels)
	assert.ElementsMatch(t, []string{`topic1`, `topic2`}, options.channels())
	assert.NotNil(t, options.RedisOptions, "invalid options")

	assert.Panics(t, func() { WithRedisURL(`://`)(&options) })

	WithRedisClient(&redis.Client{})(&options)
	assert.NotNil(t, options.Client, "invalid client object")
	assert.NotNil(t, options.connect(), "invalid client object")

	WithEncoder(encoder.JSON)(&options)
	assert.NotNil(t, options.encoder(), "invalid encoder")

	WithErrorHandler(func(msg nc.Message, err error) {})(&options)
	assert.NotNil(t, options.ErrorHandler, "invalid error handler")

	WithPanicHandler(func(msg nc.Message, recoverData any) {})(&options)
	assert.NotNil(t, options.PanicHandler, "invalid panic handler")

	WithLogger(logger.DefaultLogger)(&options)
	assert.NotNil(t, options.logger(), "invalid logger object")
}
