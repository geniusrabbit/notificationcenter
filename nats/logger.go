//
// @project GeniusRabbit 2016
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016
//

package nats

import (
	"encoding/json"

	"github.com/geniusrabbit/notificationcenter/subscriber"
	"github.com/nats-io/nats"
)

// Logger for NATS queue
type Logger struct {
	subscriber.Base
	topics []string
	conn   *nats.Conn
}

// NewLogger object
func NewLogger(topics []string, url string, options ...nats.Option) (*Logger, error) {
	var conn, err = nats.Connect(url, options...)
	if nil != err || nil == conn {
		return nil, err
	}

	return &Logger{topics: topics, conn: conn}, nil
}

// MustNewLogger object
func MustNewLogger(topics []string, url string, options ...nats.Option) *Logger {
	var log, err = NewLogger(topics, url, options...)
	if nil != err || nil == log {
		panic(err)
	}
	return log
}

// Write binary message
func (l *Logger) Write(data []byte) (err error) {
	for _, t := range l.topics {
		if nil == l.conn {
			break
		}
		if err = l.conn.Publish(t, data); nil == err {
			break
		}
	}
	return
}

// Send message
func (l *Logger) Send(messages ...interface{}) (err error) {
	if nil != messages {
		for _, it := range messages {
			if b, err := json.Marshal(it); nil == err {
				err = l.Write(b)
			}
			if nil != err {
				break
			}
		}
	}
	return
}

// Close nats client
func (l *Logger) Close() error {
	if nil != l.conn {
		l.conn.Close()
		l.conn = nil
	}
}
