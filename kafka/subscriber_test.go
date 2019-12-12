package kafka

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewSubscriberError(t *testing.T) {
	_, err := NewSubscriber(nil, "", nil)
	assert.Error(t, err)
}
