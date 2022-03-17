package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type targetInterface interface {
	Error(params ...any)
	Debugf(msg string, params ...any)
}

func TestDefaultLogger(t *testing.T) {
	var (
		def   any = DefaultLogger
		lg, _     = def.(targetInterface)
	)
	assert.NotNil(t, lg)
}
