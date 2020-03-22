package dummy

import (
	"context"
	"testing"

	nc "github.com/geniusrabbit/notificationcenter"
	"github.com/stretchr/testify/assert"
)

func TestPublisherSubscriber(t *testing.T) {
	var sub Subscriber
	var pub Publisher
	var rec = nc.FuncReceiver(func(m nc.Message) error {
		return nil
	})
	assert.NoError(t, pub.Publish(context.Background(), "message"))
	assert.NoError(t, sub.Subscribe(context.Background(), rec))
	assert.NoError(t, sub.Listen(context.Background()))
	assert.NoError(t, sub.Close())
}
