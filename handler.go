//
// @project geniusrabbit.com 2015 – 2016, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015 – 2016, 2019
//

package notificationcenter

import (
	"runtime"
)

// Handler interface
type Handler interface {
	Handle(msg Message) error
}

// FuncHandler type
type FuncHandler func(msg Message) error

// Handle this item
func (f FuncHandler) Handle(msg Message) error {
	return f(msg)
}

// MultithreadHandler it's ext interface
// which contains count of concurrent processes
type MultithreadHandler interface {
	Handler
	Concurrently() int
}

type multithread struct {
	count   int
	handler Handler
}

// Concurrently count
func (h multithread) Concurrently() int {
	return h.count
}

// Handle this item
func (h multithread) Handle(msg Message) error {
	return h.handler.Handle(msg)
}

// NewMultithreadHandler processor
func NewMultithreadHandler(count int, handler Handler) MultithreadHandler {
	if count < 0 {
		count = runtime.NumCPU() + count
	} else if count == 0 {
		count = runtime.NumCPU()
	}

	if count < 1 {
		count = 0
	}

	return multithread{
		count:   count,
		handler: handler,
	}
}
