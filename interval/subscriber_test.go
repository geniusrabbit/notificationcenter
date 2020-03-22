package interval

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	nc "github.com/geniusrabbit/notificationcenter"
	"github.com/stretchr/testify/assert"
)

func TestInterval(t *testing.T) {
	var (
		counter     = 0
		ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond*100)
	)
	defer cancel()

	sub := NewSubscriber(time.Millisecond, WithHandler(func() interface{} { return "test" }))
	err := sub.Subscribe(ctx, nc.FuncReceiver(func(msg nc.Message) error {
		switch v := MessageValue(msg).(type) {
		case string:
			counter++
		default:
			return fmt.Errorf(`invalid message type %+v`, v)
		}
		return nil
	}))
	assert.NoError(t, err)

	err = sub.Listen(ctx)
	assert.NoError(t, err)

	err = sub.Close()
	assert.NoError(t, err)

	assert.True(t, counter > 10, `generated messages`)
}

func TestIntervalPanic(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	sub := NewSubscriber(time.Millisecond,
		WithHandler(func() interface{} { return "test" }),
		WithPanicHandler(func(msg nc.Message, recoverData interface{}) {
			assert.Equal(t, "test", recoverData)
		}),
	)
	err := sub.Subscribe(ctx, nc.FuncReceiver(func(msg nc.Message) error { panic("test") }))
	assert.NoError(t, err)

	err = sub.Listen(ctx)
	assert.NoError(t, err)

	err = sub.Close()
	assert.NoError(t, err)
}

func TestIntervalError(t *testing.T) {
	testError := errors.New(`test`)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	sub := NewSubscriber(time.Millisecond,
		WithHandler(func() interface{} { return "test" }),
		WithErrorHandler(func(msg nc.Message, err error) {
			assert.Equal(t, testError, err)
		}),
	)
	err := sub.Subscribe(ctx, nc.FuncReceiver(func(msg nc.Message) error {
		return testError
	}))
	assert.NoError(t, err)

	err = sub.Listen(ctx)
	assert.NoError(t, err)

	err = sub.Close()
	assert.NoError(t, err)
}
