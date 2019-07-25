//
// @project GeniusRabbit 2016 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 - 2019
//

package nats

import (
	"bytes"

	"github.com/geniusrabbit/notificationcenter/encoder"
	"github.com/geniusrabbit/notificationcenter/subscriber"
	nats "github.com/nats-io/nats.go"
)

// Stream for NATS queue
type Stream struct {
	subscriber.Base
	topics  []string
	encoder encoder.Encoder
	conn    *nats.Conn
}

// NewStream object
func NewStream(topics []string, url string, options ...nats.Option) (*Stream, error) {
	var conn, err = nats.Connect(url, options...)
	if err != nil || conn == nil {
		return nil, err
	}
	return &Stream{
		topics:  topics,
		conn:    conn,
		encoder: encoder.JSON,
	}, nil
}

// MustNewStream object
func MustNewStream(topics []string, url string, options ...nats.Option) *Stream {
	stream, err := NewStream(topics, url, options...)
	if err != nil || stream == nil {
		panic(err)
	}
	return stream
}

// Write binary message
func (s *Stream) Write(data []byte) (err error) {
	for _, topic := range s.topics {
		if s.conn == nil {
			break
		}
		if err = s.conn.Publish(topic, data); err == nil {
			break
		}
	}
	return
}

// Send message
func (s *Stream) Send(messages ...interface{}) (err error) {
	var buff bytes.Buffer
	for _, it := range messages {
		buff.Reset()
		if err := s.encoder(it, &buff); err == nil {
			err = s.Write(buff.Bytes())
		}
		if err != nil {
			break
		}
	}
	return
}

// Close nats client
func (s *Stream) Close() error {
	s.conn.Close()
	return nil
}
