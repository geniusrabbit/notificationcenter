package bigcache

import (
	"time"

	"github.com/allegro/bigcache"
	"github.com/geniusrabbit/notificationcenter/v2/internal/objecthash"
)

// Checker provides inmemory messages test
type Checker struct {
	cache *bigcache.BigCache
}

// New checker wraps bigcache
func New(cache *bigcache.BigCache) *Checker {
	return &Checker{cache: cache}
}

// NewByConfig checker
func NewByConfig(config bigcache.Config) (*Checker, error) {
	cache, err := bigcache.NewBigCache(config)
	if err != nil {
		return nil, err
	}
	return New(cache), nil
}

// NewDefault checker
func NewDefault(duration time.Duration) (*Checker, error) {
	config := bigcache.DefaultConfig(duration)
	return NewByConfig(config)
}

// IsSkip message if was sent
func (ch *Checker) IsSkip(msg any) bool {
	val, _ := ch.cache.Get(objecthash.Hash(msg))
	return len(val) == 1 && val[0] == 't'
}

// MarkAsSent message to the publisher
func (ch *Checker) MarkAsSent(msg any) error {
	return ch.cache.Set(objecthash.Hash(msg), []byte(`t`))
}
