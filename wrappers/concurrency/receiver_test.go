package concurrency

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	nc "github.com/geniusrabbit/notificationcenter"
	"github.com/geniusrabbit/notificationcenter/mocks"
)

func TestConcurrency(t *testing.T) {
	var (
		wg       sync.WaitGroup
		lastErr  error
		receiver = nc.FuncReceiver(func(msg nc.Message) error {
			if bytes.Equal([]byte(`error`), msg.Body()) {
				return fmt.Errorf(`test error`)
			}
			wg.Done()
			return nil
		})
		errWrapper = func(err error, msg nc.Message) {
			defer wg.Done()
			lastErr = err
		}
		rc = WithWorkers(receiver, 10, errWrapper,
			WithWorkerPoolSize(5), WithRecoverHandler(nil))
	)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		body := []byte(`test`)
		if i == 73 {
			body = []byte(`error`)
		}
		err := rc.Receive(mocks.NewMessage(context.TODO(), `test`, body, nil))
		assert.NoError(t, err, `receive message`)
	}

	wg.Wait()
	assert.Error(t, lastErr, `test error`)
	assert.NoError(t, rc.(io.Closer).Close())
}
