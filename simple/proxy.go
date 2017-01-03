//
// @project GeniusRabbit 2016
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016
//

package simple

import (
	"github.com/geniusrabbit/notificationcenter/subscriber"
)

// EventProxy struct
type EventProxy struct {
	subscriber.Base
}

// NewProxy object
func NewProxy() *EventProxy {
	return &EventProxy{}
}

// Send message
func (e *EventProxy) Send(msg ...interface{}) (err error) {
	for _, m := range msg {
		if err = e.Handle(m, true); nil != err {
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
