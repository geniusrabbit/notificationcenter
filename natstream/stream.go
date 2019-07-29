//
// @project GeniusRabbit 2018 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018 - 2019
//

package natstream

import (
	"bytes"
	"errors"

	"github.com/geniusrabbit/notificationcenter/encoder"
	"github.com/geniusrabbit/notificationcenter/subscriber"
	nstream "github.com/nats-io/stan.go"
)

// Error set
var (
	ErrConnectionIsClosed = errors.New("[notificationcenter:natstream] connection is closed")
)

// Stream for NATS queue
type Stream struct {
	subscriber.Base
	topic   string
	encoder encoder.Encoder
	conn    nstream.Conn
}

// NewStream of the NATS stream server
func NewStream(url, clusterID, clientID, topic string, options ...nstream.Option) (*Stream, error) {
	conn, err := nstream.Connect(clusterID, clientID, append(options, nstream.NatsURL(url))...)
	if err != nil || conn == nil {
		return nil, err
	}
	return &Stream{
		topic:   topic,
		conn:    conn,
		encoder: encoder.JSON,
	}, nil
}

// MustNewStream of the NATS stream server
func MustNewStream(url, clusterID, clientID, topic string, options ...nstream.Option) *Stream {
	stream, err := NewStream(url, clusterID, clientID, topic, options...)
	if err != nil || stream == nil {
		panic(err)
	}
	return stream
}

// Write binary message
func (s *Stream) Write(data []byte) (err error) {
	if s.conn == nil {
		return ErrConnectionIsClosed
	}
	return s.conn.Publish(s.topic, data)
}

// Send messages into the stream
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

// Close nats-stream client
func (s *Stream) Close() error {
	s.conn.Close()
	return nil
}
