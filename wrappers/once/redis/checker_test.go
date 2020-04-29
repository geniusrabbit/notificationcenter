package redis

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/elliotchance/redismock"
	"github.com/go-redis/redis"
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
	msg := struct{ s string }{s: `test`}
	mr, err := miniredis.Run()
	assert.NoError(t, err)
	defer mr.Close()
	checker := New(newTestRedis(mr), time.Minute)
	assert.False(t, checker.IsSkip(msg))
	assert.NoError(t, checker.MarkAsSent(msg))
	assert.True(t, checker.IsSkip(msg))
}
