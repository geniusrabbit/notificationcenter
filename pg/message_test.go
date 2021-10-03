package pg

import (
	"context"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func Test_Message(t *testing.T) {
	msg := fromPgNotify(context.TODO(),
		&pq.Notification{
			BePid:   1,
			Channel: "test",
			Extra:   `{"data": "test"}`,
		},
	)
	assert.NotNil(t, msg.Context())
	assert.NotNil(t, msg.Notification())
	assert.Equal(t, []byte(`{"data": "test"}`), msg.Body())
	assert.Equal(t, `test-1`, msg.ID())
	assert.Nil(t, msg.Ack())
}
