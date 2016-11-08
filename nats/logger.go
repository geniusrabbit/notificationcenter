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

// Log for NATS queue
type Log struct {
	subscriber.Base
	topics []string
	conn   *nats.Conn
}

// NewLog object
func NewLog(topics []string, url string, options ...nats.Option) (*Log, error) {
	var conn, err = nats.Connect(url, options...)
	if nil != err || nil == conn {
		return nil, err
	}

	return &Log{topics: topics, conn: conn}, nil
}

// MustNewLog object
func MustNewLog(topics []string, url string, options ...nats.Option) *Log {
	var log, err = NewLog(topics, url, options...)
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
func (l *Log) Send(messages ...interface{}) (err error) {
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
func (l *Log) Close() error {
	if nil != l.conn {
		l.conn.Close()
		l.conn = nil
	}
}
