package notificationcenter

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Handler(t *testing.T) {
	var (
		resultID string
		data     []byte
		ackErr   error
		fh       = FuncHandler(func(msg Message) error {
			resultID = msg.ID()
			data = msg.Body()
			ackErr = msg.Ack()
			return nil
		})
	)
	fh.Handle(NewMessageMock("test", []byte("{data}"), fmt.Errorf("ack")))
	assert.Equal(t, "test", resultID, "execute FuncHandler")
	assert.Equal(t, []byte("{data}"), data, "execute MultithreadHandler")
	assert.Error(t, ackErr, "Acknowledgment error")
}

func Test_MultithreadHandler(t *testing.T) {
	var (
		resultID string
		data     []byte
		ackErr   error
		mh       = NewMultithreadHandler(10, FuncHandler(func(msg Message) error {
			resultID = msg.ID()
			data = msg.Body()
			ackErr = msg.Ack()
			return nil
		}))
	)
	mh.Handle(NewMessageMock("test_multithread", []byte("{data}"), fmt.Errorf("ack")))
	assert.Equal(t, "test_multithread", resultID, "execute MultithreadHandler")
	assert.Equal(t, []byte("{data}"), data, "execute MultithreadHandler")
	assert.Error(t, ackErr, "Acknowledgment error")
	assert.Equal(t, 10, mh.Concurrently(), "execute MultithreadHandler")
}
