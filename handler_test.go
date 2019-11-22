package notificationcenter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Handler(t *testing.T) {
	var (
		resultID string
		fh       = FuncHandler(func(msg Message) error {
			resultID = msg.ID()
			return nil
		})
	)
	fh.Handle(NewMessageMock("test", nil, nil))
	assert.Equal(t, "test", resultID, "execute FuncHandler")
}

func Test_MultithreadHandler(t *testing.T) {
	var (
		resultID string
		mh       = NewMultithreadHandler(10, FuncHandler(func(msg Message) error {
			resultID = msg.ID()
			return nil
		}))
	)
	mh.Handle(NewMessageMock("test_multithread", nil, nil))
	assert.Equal(t, "test_multithread", resultID, "execute MultithreadHandler")
	assert.Equal(t, 10, mh.Concurrently(), "execute MultithreadHandler")
}
