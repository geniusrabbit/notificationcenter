package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type targetInterface interface {
	Error(params ...interface{})
	Debugf(msg string, params ...interface{})
}

func TestDefaultLogger(t *testing.T) {
	var (
		def   interface{} = DefaultLogger
		lg, _             = def.(targetInterface)
	)
	assert.NotNil(t, lg)
}
