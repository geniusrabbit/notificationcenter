package redis

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

func TestRedisChecker(t *testing.T) {
	ctx := context.TODO()
	msg := struct{ s string }{s: `test`}
	mr, err := miniredis.Run()
	assert.NoError(t, err)
	defer mr.Close()

	checker, err := NewByURL("redis://"+mr.Addr(), time.Minute)
	assert.NoError(t, err)

	assert.False(t, checker.IsSkip(ctx, msg))
	assert.NoError(t, checker.MarkAsSent(ctx, msg))
	assert.True(t, checker.IsSkip(ctx, msg))
}
