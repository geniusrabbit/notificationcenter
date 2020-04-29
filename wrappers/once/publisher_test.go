package once

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/geniusrabbit/notificationcenter/dummy"
	"github.com/geniusrabbit/notificationcenter/wrappers/once/bigcache"
)

func TestPublisher(t *testing.T) {
	ctx := context.TODO()
	cache, err := bigcache.NewDefault(time.Second)
	assert.NoError(t, err)
	assert.Panics(t, func() { MewPublisherWrapper(nil, nil) })
	assert.Panics(t, func() { MewPublisherWrapper(dummy.Publisher{}, nil) })
	pub := MewPublisherWrapper(dummy.Publisher{}, cache)
	assert.NoError(t, pub.Publish(ctx, `test`))
	assert.NoError(t, pub.Close())
}
