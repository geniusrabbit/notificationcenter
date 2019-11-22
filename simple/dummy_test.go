package simple

import (
	"testing"

	nc "github.com/geniusrabbit/notificationcenter"
	"github.com/stretchr/testify/assert"
)

func Test_Dummy(t *testing.T) {
	dummy := NewDummy()
	handler := nc.FuncHandler(func(m nc.Message) error {
		return nil
	})
	assert.NoError(t, dummy.Send("message"))
	assert.NoError(t, dummy.Subscribe(handler))
	assert.NoError(t, dummy.Unsubscribe(handler))
	assert.NoError(t, dummy.Listen())
	assert.NoError(t, dummy.Close())
}
