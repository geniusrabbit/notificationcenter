package pg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Message(t *testing.T) {
	msg := &message{
		BePid:   1,
		Channel: "test",
		Extra:   `{"data": "test"}`,
	}
	assert.NotNil(t, msg.Notification())
	assert.Equal(t, []byte(`{"data": "test"}`), msg.Body())
	assert.Equal(t, `test_1`, msg.ID())
}
