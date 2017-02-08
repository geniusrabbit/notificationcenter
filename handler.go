//
// @project geniusrabbit.com 2015 – 2016
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015 – 2016
//

package notificationcenter

import (
	"runtime"
)

// Handler interface
type Handler interface {
	Handle(item interface{}) error
}

// FuncHandler type
type FuncHandler func(item interface{}) error

// Handle this item
func (f FuncHandler) Handle(item interface{}) error {
	return f(item)
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
func (h multithread) Handle(item interface{}) error {
	return h.handler.Handle(item)
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
