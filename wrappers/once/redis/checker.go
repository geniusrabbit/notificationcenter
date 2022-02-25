package redis

import (
	"context"
	"time"

	"github.com/geniusrabbit/notificationcenter/internal/objecthash"
	"github.com/go-redis/redis/v8"
)

// Checker provides inmemory messages test
type Checker struct {
	lifetime time.Duration
	client   redis.Cmdable
}

// New checker wraps bigcache
func New(client redis.Cmdable, lifetime time.Duration) *Checker {
	return &Checker{client: client, lifetime: lifetime}
}

// NewByHost checker
func NewByHost(host string, lifetime time.Duration) *Checker {
	return New(redis.NewClient(&redis.Options{Addr: host}), lifetime)
}

// IsSkip message if was sent
func (ch *Checker) IsSkip(ctx context.Context, msg any) bool {
	val, _ := ch.client.Get(ctx, objecthash.Hash(msg)).Result()
	return len(val) == 1 && val == `t`
}

// MarkAsSent message to the publisher
func (ch *Checker) MarkAsSent(ctx context.Context, msg any) error {
	return ch.client.SetNX(ctx, objecthash.Hash(msg), []byte(`t`), ch.lifetime).Err()
}
