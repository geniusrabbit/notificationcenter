package logger

import "log"

// DefaultLogger object for logging
var DefaultLogger defaultLogger

type defaultLogger struct{}

func (l defaultLogger) Error(params ...any) {
	log.Println(append([]any{`[error]`}, params...)...)
}

func (l defaultLogger) Debugf(msg string, params ...any) {
	log.Printf(`[debug] `+msg+`\n`, params...)
}
