package interval

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	msg1 := &message{v: "test"}
	msg2 := &message{v: []byte("test")}
	msg3 := &message{v: []int{1, 2, 3}}

	assert.Equal(t, "", msg1.ID(), "message 1")
	assert.Equal(t, "test", msg1.Value(), "message 1")
	assert.Equal(t, []byte("test"), msg1.Body(), "message 1")
	assert.Equal(t, "test", MessageValue(msg1), "message 1")
	assert.Nil(t, msg1.Ack(), "message 1")

	assert.Equal(t, "", msg2.ID(), "message 2")
	assert.Equal(t, []byte("test"), msg2.Value(), "message 2")
	assert.Equal(t, []byte("test"), msg2.Body(), "message 2")
	assert.Equal(t, []byte("test"), MessageValue(msg2), "message 2")
	assert.Nil(t, msg2.Ack(), "message 2")

	assert.Equal(t, "", msg3.ID(), "message 3")
	assert.Equal(t, []int{1, 2, 3}, msg3.Value(), "message 3")
	assert.Equal(t, []byte(`[1,2,3]`), msg3.Body(), "message 3")
	assert.Equal(t, []int{1, 2, 3}, MessageValue(msg3), "message 3")
	assert.Nil(t, msg3.Ack(), "message 3")

	assert.Panics(t, func() { MessageValue(`panic`) }, "message panic")
}
