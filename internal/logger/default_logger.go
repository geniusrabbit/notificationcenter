package logger

import "log"

// DefaultLogger object for logging
var DefaultLogger defaultLogger

type defaultLogger struct{}

func (l defaultLogger) Error(params ...interface{}) {
	log.Println(append([]interface{}{`[error]`}, params...)...)
}

func (l defaultLogger) Debugf(msg string, params ...interface{}) {
	log.Printf(`[debug] `+msg+`\n`, params...)
}
