package concurrency

import "github.com/demdxx/rpool"

// Option func type which adjust option values
type Option func() rpool.Option

// WithWorkerPoolSize setup maximal size of worker pool
func WithWorkerPoolSize(size int) Option {
	return func() rpool.Option {
		return rpool.WithWorkerPoolSize(size)
	}
}

// WithRecoverHandler defined error handler
func WithRecoverHandler(f func(interface{})) Option {
	return func() rpool.Option {
		return rpool.WithRecoverHandler(f)
	}
}
