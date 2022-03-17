package gochan

import (
	"context"
	"encoding/json"
	"sync"
	"testing"

	nc "github.com/geniusrabbit/notificationcenter/v2"
	"github.com/stretchr/testify/assert"
)

type testMessage struct {
	Text string
}

func TestSubscriberPublishing(t *testing.T) {
	var (
		wg       sync.WaitGroup
		testText string
		proxy    = New(10)
		pub      = proxy.Publisher()
	)

	defer func() { _ = proxy.Close() }()
	err := proxy.Subscribe(context.Background(), nc.FuncReceiver(func(msg nc.Message) error {
		defer wg.Done()
		var m testMessage
		if err := json.Unmarshal(msg.Body(), &m); err != nil {
			return err
		}
		testText = m.Text
		return msg.Ack()
	}))

	assert.NoError(t, err, `Subscribe`)

	go func() {
		err = proxy.Listen(context.Background())
		assert.NoError(t, err, `Listen`)
	}()

	wg.Add(1)
	err = pub.Publish(context.Background(), testMessage{Text: "test-message"})
	assert.NoError(t, err, `Publish`)

	wg.Wait()
	assert.Equal(t, "test-message", testText)
}
