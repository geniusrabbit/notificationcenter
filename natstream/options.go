package natstream

import (
	"context"
	"math/rand"
	"net/url"
	"strings"

	nats "github.com/nats-io/nats.go"
	nstream "github.com/nats-io/stan.go"

	nc "github.com/geniusrabbit/notificationcenter/v2"
	"github.com/geniusrabbit/notificationcenter/v2/encoder"
	"github.com/geniusrabbit/notificationcenter/v2/internal/logger"
)

type loggerInterface interface {
	Error(params ...any)
	Debugf(msg string, params ...any)
}

// Options of the NATS wrapper
type Options struct {
	Ctx context.Context

	// Raw options from the "github.com/nats-io/stan.go" module
	NatsOptions []nstream.Option

	// NatsSubscriptions suboptions of subscriptions
	NatsSubscriptions []nstream.Subscription

	// Name of the subscription group
	GroupName string

	// Names of topics for subscribing or publishing
	Topics []string

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

func (opt *Options) randomTopic() string {
	if len(opt.Topics) == 0 {
		return `default`
	}
	return opt.Topics[rand.Int()%len(opt.Topics)]
}

func (opt *Options) group() string {
	if opt.GroupName == `` {
		return `default`
	}
	return opt.GroupName
}

func (opt *Options) context() context.Context {
	if opt.Ctx == nil {
		return context.Background()
	}
	return opt.Ctx
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
func WithNatsURL(urlString string) Option {
	return func(opt *Options) {
		u, err := url.Parse(urlString)
		if err != nil {
			panic(err)
		}
		if len(u.Path) > 1 {
			opt.GroupName = u.Path[1:]
		}
		topics := strings.Split(u.Query().Get(`topics`), `,`)
		if len(topics) == 1 && topics[0] == `` {
			topics = nil
		}
		u.Path = ``
		u.RawQuery = ``
		if len(topics) > 0 {
			opt.Topics = topics
		}
		opt.NatsOptions = append(opt.NatsOptions, nstream.NatsURL(u.String()))
	}
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

// WithContext puts the client context value
func WithContext(ctx context.Context) Option {
	return func(opt *Options) {
		opt.Ctx = ctx
	}
}

// WithTopics will set the list of topics for publishing or subscribing
func WithTopics(topics ...string) Option {
	return func(opt *Options) {
		opt.Topics = topics
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
