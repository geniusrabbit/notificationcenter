//
// @project GeniusRabbit 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018
//

package nats

import (
	"encoding/json"
	"errors"

	"github.com/geniusrabbit/notificationcenter/subscriber"
	nstream "github.com/nats-io/go-nats-streaming"
)

// Error set
var (
	ErrConnectionIsClosed = errors.New("Connection is closed")
)

// Logger for NATS queue
type Logger struct {
	subscriber.Base
	topic string
	conn  *nstream.Conn
}

// NewLogger object
func NewLogger(url, clusterID, clientID, topic string, options ...nstream.Option) (*Logger, error) {
	var conn, err = nstream.Connect(clusterID, clientID, nstream.NatsURL(URL), options...)
	if err != nil || conn == nil {
		return nil, err
	}
	return &Logger{topic: topic, conn: conn}, nil
}

// MustNewLogger object
func MustNewLogger(url, clusterID, clientID, topic string, options ...nstream.Option) *Logger {
	var log, err = NewLogger(topic, url, options...)
	if err != nil || log == nil {
		panic(err)
	}
	return log
}

// Write binary message
func (l *Logger) Write(data []byte) (err error) {
	if l.conn == nil {
		return ErrConnectionIsClosed
	}
	return l.conn.Publish(t, data)
}

// Send message
func (l *Logger) Send(messages ...interface{}) (err error) {
	for _, it := range messages {
		if b, err := json.Marshal(it); err == nil {
			err = l.Write(b)
		}
		if err != nil {
			break
		}
	}
	return
}

// Close nats-stream client
func (l *Logger) Close() error {
	if l.conn != nil {
		l.conn.Close()
		l.conn = nil
	}
	return nil
}
