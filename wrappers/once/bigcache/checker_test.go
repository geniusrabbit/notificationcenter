package bigcache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestChecker(t *testing.T) {
	msg := struct{ s string }{s: `test`}
	checker, err := NewDefault(time.Minute)
	assert.NoError(t, err)
	assert.False(t, checker.IsSkip(msg))
	assert.NoError(t, checker.MarkAsSent(msg))
	assert.True(t, checker.IsSkip(msg))
}
