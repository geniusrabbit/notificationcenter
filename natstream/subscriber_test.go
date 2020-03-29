package natstream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSubscriberPanic(t *testing.T) {
	defer func() {
		rec := recover()
		assert.NotNil(t, rec)
	}()
	sub := MustNewSubscriber(WithNatsURL(`nats://demo`))
	sub.Close()
}

func TestNewSubscriber(t *testing.T) {
	sub, err := NewSubscriber(
		WithNatsURL(`nats://demo`),
		WithGroupName(`test`),
		WithLogger(nil),
	)
	if !assert.Error(t, err) {
		assert.NoError(t, sub.Close())
	}
}
