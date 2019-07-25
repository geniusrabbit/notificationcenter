//
// @project GeniusRabbit 2016, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016, 2019
//

package simple

import (
	"bytes"

	"github.com/geniusrabbit/notificationcenter/encoder"
	"github.com/geniusrabbit/notificationcenter/subscriber"
)

// EventProxy implements basic event processor with message processing without chanel pool
type EventProxy struct {
	subscriber.Base
	encoder encoder.Encoder
}

// NewProxy object with encoding option
func NewProxy(enc ...encoder.Encoder) *EventProxy {
	proxy := &EventProxy{}
	if len(enc) > 0 && enc[0] != nil {
		proxy.encoder = enc[0]
	} else {
		proxy.encoder = encoder.JSON
	}
	return proxy
}

// Send message
func (e *EventProxy) Send(msgs ...interface{}) (err error) {
	for _, m := range msgs {
		var msg *message
		if msg, err = e.messageFrom(m); err != nil {
			break
		}
		if err = e.Handle(msg, true); err != nil {
			break
		}
	}
	return
}

// Listen process
func (e *EventProxy) Listen() error {
	return nil
}

// Close proxy handler
func (e *EventProxy) Close() error {
	return e.Base.CloseAll()
}

func (e *EventProxy) messageFrom(m interface{}) (*message, error) {
	var buff bytes.Buffer
	if err := e.encoder(m, &buff); err != nil {
		return nil, err
	}
	return &message{data: buff.Bytes()}, nil
}
