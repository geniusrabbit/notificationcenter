package natstream

import (
	nats "github.com/nats-io/nats.go"
	nstream "github.com/nats-io/stan.go"

	nc "github.com/geniusrabbit/notificationcenter"
	"github.com/geniusrabbit/notificationcenter/encoder"
	"github.com/geniusrabbit/notificationcenter/internal/logger"
)

// Options of the NATS wrapper
type Options struct {
	// Raw options from the "github.com/nats-io/stan.go" module
	NatsOptions []nstream.Option

	// NatsSubscriptions suboptions of subscriptions
	NatsSubscriptions []nstream.Subscription

	// Name of the subscription group
	GroupName string

	// ClusterID common for the group of services
	ClusterID string

	// Client ID unical for the service
	ClientID string

	// ErrorHandler of message processing
	ErrorHandler nc.ErrorHandler

	// PanicHandler process panic
	PanicHandler nc.PanicHandler

	// Message encoder interface
	Encoder encoder.Encoder

	// Logger of subscriber
	Logger loggerInterface
}

func (opt *Options) encoder() encoder.Encoder {
	if opt.Encoder == nil {
		return encoder.JSON
	}
	return opt.Encoder
}

func (opt *Options) group() string {
	if opt.GroupName == `` {
		return `default`
	}
	return opt.GroupName
}

func (opt *Options) clusterID() string {
	if opt.ClusterID == `` {
		return `default`
	}
	return opt.ClusterID
}

func (opt *Options) clientID() string {
	if opt.ClientID == `` {
		return `default`
	}
	return opt.ClientID
}

func (opt *Options) logger() loggerInterface {
	if opt.Logger == nil {
		return logger.DefaultLogger
	}
	return opt.Logger
}

// Option of the NATS subscriber or publisher
type Option func(opt *Options)

// WithNatsURL is an Option to set the URL the client should connect to.
// The url can contain username/password semantics. e.g. nats://derek:pass@localhost:4222
// Comma separated arrays are also supported, e.g. urlA, urlB.
func WithNatsURL(url string) Option {
	return WithNatsOptions(nstream.NatsURL(url))
}

// WithNatsConn is an Option to set the underlying NATS connection to be used
// by a streaming connection object. When such option is set, closing the
// streaming connection does not close the provided NATS connection.
func WithNatsConn(nc *nats.Conn) Option {
	return WithNatsOptions(nstream.NatsConn(nc))
}

// WithClusterID puts the cluster ID value
func WithClusterID(id string) Option {
	return func(opt *Options) {
		opt.ClusterID = id
	}
}

// WithClientID puts the client ID value
func WithClientID(id string) Option {
	return func(opt *Options) {
		opt.ClientID = id
	}
}

// WithGroupName puts the name group of the subsciber or publisher
func WithGroupName(name string) Option {
	return func(opt *Options) {
		opt.GroupName = name
	}
}

// WithNatsOptions provides options of the NATS module
func WithNatsOptions(natsOpts ...nstream.Option) Option {
	return func(opt *Options) {
		opt.NatsOptions = append(opt.NatsOptions, natsOpts...)
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

// WithEncoder set the message encoder
func WithEncoder(encoder encoder.Encoder) Option {
	return func(opt *Options) {
		opt.Encoder = encoder
	}
}

// WithLogger provides logging interface
func WithLogger(logger loggerInterface) Option {
	return func(opt *Options) {
		opt.Logger = logger
	}
}
