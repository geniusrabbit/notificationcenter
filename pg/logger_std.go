package pg

import (
	"log"
	"os"
)

// LoggerStd implementation it's just dummy implementation without level checking
type LoggerStd log.Logger

func nlog() *LoggerStd {
	return (*LoggerStd)(log.New(os.Stdout, "pgevent_", log.LstdFlags))
}

// Info level printing
func (l *LoggerStd) Info(params ...any) {
	l.lg().Println(params...)
}

// Error level printing
func (l *LoggerStd) Error(params ...any) {
	l.lg().Println(params...)
}

// Debugf level printing
func (l *LoggerStd) Debugf(msg string, params ...any) {
	l.lg().Printf(msg, params...)
}

func (l *LoggerStd) lg() *log.Logger {
	return (*log.Logger)(l)
}
