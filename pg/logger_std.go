package pg

import (
	"log"
	"os"
)

type LoggerStd log.Logger

func nlog() *LoggerStd {
	return (*LoggerStd)(log.New(os.Stdout, "pgevent_", log.LstdFlags))
}

func (l *LoggerStd) Info(params ...interface{}) {
	l.lg().Println(params...)
}

func (l *LoggerStd) Error(params ...interface{}) {
	l.lg().Println(params...)
}

func (l *LoggerStd) Debugf(msg string, params ...interface{}) {
	l.lg().Printf(msg, params...)
}

func (l *LoggerStd) lg() *log.Logger {
	return (*log.Logger)(l)
}
