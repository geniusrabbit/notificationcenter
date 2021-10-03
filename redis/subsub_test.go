package redis

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"
	nc "github.com/geniusrabbit/notificationcenter"
	"github.com/stretchr/testify/assert"
)

func TestSubSub(t *testing.T) {
	type tmsg struct{ S string }
	var (
		rw      sync.RWMutex
		msg     tmsg
		mr, err = miniredis.Run()
	)
	assert.NoError(t, err)
	defer mr.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	sub := MustNewSubscriber(WithRedisURL("redis://" + mr.Addr()))
	pub := MustNewPublisher(WithRedisURL("redis://" + mr.Addr()))

	err = sub.Subscribe(ctx, nc.FuncReceiver(func(m nc.Message) error {
		rw.Lock()
		defer rw.Unlock()
		return json.Unmarshal(m.Body(), &msg)
	}))

	time.Sleep(time.Second)

	assert.NoError(t, err, "subscibe chanel")
	assert.NoError(t, pub.Publish(ctx, tmsg{S: "test"}), "publish message")

	<-ctx.Done()

	rw.RLock()
	defer rw.RUnlock()
	assert.Equal(t, "test", msg.S)
	assert.NoError(t, pub.Close(), "close pub redis connection")
	assert.NoError(t, sub.Close(), "close sub redis connection")
}
