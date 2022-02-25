package interval

import (
	nc "github.com/geniusrabbit/notificationcenter/v2"
)

type handler func() any

// Options time interval wrapper
type Options struct {
	// Handler generates message by time interval
	Handler handler

	// ErrorHandler of message processing
	ErrorHandler nc.ErrorHandler

	// PanicHandler process panic
	PanicHandler nc.PanicHandler
}

// Option func type
type Option func(opt *Options)

// WithHandler set custom message generation function
func WithHandler(h handler) Option {
	return func(options *Options) {
		options.Handler = h
	}
}

// WithErrorHandler set handler of error processing
func WithErrorHandler(h nc.ErrorHandler) Option {
	return func(options *Options) {
		options.ErrorHandler = h
	}
}

// WithPanicHandler set handler of panic processing
func WithPanicHandler(h nc.PanicHandler) Option {
	return func(options *Options) {
		options.PanicHandler = h
	}
}
