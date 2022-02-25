package concurrency

import (
	"github.com/demdxx/rpool/v2"
	nc "github.com/geniusrabbit/notificationcenter/v2"
)

type errHandlerFnk func(err error, msg nc.Message)

type receiver struct {
	receiver   nc.Receiver
	errHandler errHandlerFnk
	execPool   *rpool.PoolFunc[nc.Message]
}

// WithWorkers wraps receiver in concurrent pool
func WithWorkers(recr nc.Receiver, workers int, errHandler errHandlerFnk, options ...Option) nc.Receiver {
	if workers <= 1 && len(options) < 1 {
		return recr
	}
	poolOptions := []rpool.Option{rpool.WithWorkerPoolSize(workers)}
	for _, opt := range options {
		poolOptions = append(poolOptions, opt())
	}
	r := &receiver{receiver: recr, errHandler: errHandler}
	r.execPool = rpool.NewPoolFunc(r.executor, poolOptions...)
	return r
}

func (r *receiver) Receive(msg nc.Message) error {
	r.execPool.Call(msg)
	return nil
}

func (r *receiver) executor(msg nc.Message) {
	err := r.receiver.Receive(msg)
	if err != nil {
		if r.errHandler != nil {
			r.errHandler(err, msg)
		} else {
			panic(err)
		}
	}
}

func (r *receiver) Close() error {
	return r.execPool.Close()
}
