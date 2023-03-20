package redis

import (
	"context"
	"time"

	"github.com/geniusrabbit/notificationcenter/v2/internal/objecthash"
	"github.com/redis/go-redis/v9"
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

// NewByURL returns checker by redis URL
// redis://[:password]@host:port/db?max_idle=3&max_active=5&idle_timeout=240s
func NewByURL(url string, lifetime time.Duration) (*Checker, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	return New(redis.NewClient(opt), lifetime), nil
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
