package natstream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPublisherPanic(t *testing.T) {
	defer func() {
		rec := recover()
		assert.NotNil(t, rec)
	}()
	pub := MustNewPublisher(`test`, WithNatsURL(`nats://demo`))
	pub.Close()
}

func TestNewPublisher(t *testing.T) {
	pub, err := NewPublisher(`test`,
		WithNatsURL(`nats://demo`),
		WithClusterID(`test`),
		WithClientID(`test`),
		WithEncoder(nil),
		WithErrorHandler(nil),
		WithPanicHandler(nil),
	)
	if !assert.Error(t, err) {
		assert.NoError(t, pub.Close())
	}
}
