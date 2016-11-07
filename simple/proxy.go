//
// @project GeniusRabbit 2016
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016
//

package simple

// Handler interface
type Handler interface {
	Handle(item interface{}) error
}

// EventProxy struct
type EventProxy struct {
	handlers []Handler
}

// NewProxy object
func NewProxy() *EventProxy {
	return &EventProxy{}
}

// Subscribe handler
func (e *EventProxy) Subscribe(handler Handler) error {
	e.handlers = append(e.handlers, handler)
	return nil
}

// Send message
func (e *EventProxy) Send(msg ...interface{}) (err error) {
	if nil != e.handlers {
		for _, h := range e.handlers {
			for _, m := range msg {
				if err = h.Handle(m); nil != err {
					break
				}
			}
		} // end for
	}
	return
}

// Close proxy handler
func (e *EventProxy) Close() {
	e.Send(nil)
	e.handlers = nil
}
