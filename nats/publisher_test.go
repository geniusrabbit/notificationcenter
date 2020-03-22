package nats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPublisherPanic(t *testing.T) {
	defer func() {
		rec := recover()
		assert.NotNil(t, rec)
	}()
	pub := MustNewPublisher(nil, WithNatsURL(`nats://demo`), WithTokenHandler(nil))
	pub.Close()
}

func TestNewPublisher(t *testing.T) {
	pub, err := NewPublisher(nil,
		WithNatsURL(`nats://demo`),
		WithClientName(`test`),
		WithEncoder(nil),
		WithErrorHandler(nil),
		WithPanicHandler(nil),
		WithUserInfo(`test`, `test`),
		WithToken(`token`),
		WithUserCredentials(``),
	)
	if !assert.Error(t, err) {
		assert.NoError(t, pub.Close())
	}
}
