package notificationcenter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiPublisher(t *testing.T) {
	sum := 0
	pubs := MultiPublisher{
		FuncPublisher(func(ctx context.Context, messages ...any) error {
			for _, v := range messages {
				sum += v.(int)
			}
			return nil
		}),
		FuncPublisher(func(ctx context.Context, messages ...any) error {
			for _, v := range messages {
				sum += v.(int)
			}
			return nil
		}),
	}
	err := pubs.Publish(context.TODO(), int(1), int(2), int(3))
	assert.NoError(t, err, "publisher unexpected error")
	assert.Equal(t, int(1+2+3)*2, sum, "publisher unexpected behaviour")
}
