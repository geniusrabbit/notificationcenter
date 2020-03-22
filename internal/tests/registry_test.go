package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	nc "github.com/geniusrabbit/notificationcenter"
	"github.com/geniusrabbit/notificationcenter/dummy"
)

func TestRegistry(t *testing.T) {
	r := nc.NewRegistry()

	pub := dummy.Publisher{}
	sub := dummy.Subscriber{}

	err := r.Register(`test-pub`, pub, `test-sub`, sub)
	assert.NoError(t, err)
	assert.NotNil(t, r.Publisher(`test-pub`))
	assert.NotNil(t, r.Subscriber(`test-sub`))

	err = r.Subscribe(context.Background(), `test-sub`,
		nc.FuncReceiver(func(msg nc.Message) error { return nil }))
	assert.NoError(t, err, `subscribe`)

	err = r.Subscribe(context.Background(), `test-sub2`,
		nc.FuncReceiver(func(msg nc.Message) error { return nil }))
	assert.Error(t, err, `subscribe error`)

	err = r.Publish(context.Background(), `test-pub`, "test")
	assert.NoError(t, err, `publish`)

	err = r.Close()
	assert.NoError(t, err, `close`)
}
