package redis

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/elliotchance/redismock/v8"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

// newTestRedis returns a redis.Cmdable.
func newTestRedis(mr *miniredis.Miniredis) *redismock.ClientMock {
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	return redismock.NewNiceMock(client)
}

func TestRedisChecker(t *testing.T) {
	ctx := context.TODO()
	msg := struct{ s string }{s: `test`}
	mr, err := miniredis.Run()
	assert.NoError(t, err)
	defer mr.Close()
	checker := New(newTestRedis(mr), time.Minute)
	assert.False(t, checker.IsSkip(ctx, msg))
	assert.NoError(t, checker.MarkAsSent(ctx, msg))
	assert.True(t, checker.IsSkip(ctx, msg))
}
