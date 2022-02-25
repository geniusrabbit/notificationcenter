package nats

import (
	"context"
	"net/url"
	"strings"

	nats "github.com/nats-io/nats.go"

	nc "github.com/geniusrabbit/notificationcenter/v2"
	"github.com/geniusrabbit/notificationcenter/v2/encoder"
	"github.com/geniusrabbit/notificationcenter/v2/internal/logger"
)

// Options of the NATS wrapper
type Options struct {
	Ctx context.Context

	// Inited NATS connection
	Conn *nats.Conn

	// ConnURL of connection
	ConnURL string

	// Raw options from the "github.com/nats-io/nats.go" module
	NatsOptions []nats.Option

	// Name of the subscription group
	GroupName string

	// Names of topics for subscribing or publishing
	Topics []string

	// ErrorHandler of message processing
	ErrorHandler nc.ErrorHandler

	// PanicHandler process panic
	PanicHandler nc.PanicHandler

	// Message encoder interface
	Encoder encoder.Encoder

	// Logger of subscriber
	Logger loggerInterface
}

func (opt *Options) clientConn() (*nats.Conn, error) {
	if opt.Conn != nil {
		return opt.Conn, nil
	}
	return nats.Connect(opt.ConnURL, opt.NatsOptions...)
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

func (opt *Options) context() context.Context {
	if opt.Ctx == nil {
		return context.Background()
	}
	return opt.Ctx
}

func (opt *Options) logger() loggerInterface {
	if opt.Logger == nil {
		return logger.DefaultLogger
	}
	return opt.Logger
}

// Option of the NATS subscriber or publisher
type Option func(opt *Options)

// WithNatsConn is an Option to inited client connection
func WithNatsConn(conn *nats.Conn) Option {
	return func(opt *Options) {
		opt.Conn = conn
	}
}

// WithNatsURL is an Option to set the URL the client should connect to.
// The url can contain username/password semantics. e.g. nats://derek:pass@localhost:4222/{groupName}?topics=topic1,topic2
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

		opt.ConnURL = u.String()
		if len(topics) > 0 {
			opt.Topics = topics
		}
	}
}

// WithClientName is an Option to set the client name.
func WithClientName(name string) Option {
	return WithNatsOptions(nats.Name(name))
}

// WithUserInfo is an Option to set the username and password to
// use when not included directly in the URLs.
func WithUserInfo(user, password string) Option {
	return WithNatsOptions(nats.UserInfo(user, password))
}

// WithToken is an Option to set the token to use
// when a token is not included directly in the URLs
// and when a token handler is not provided.
func WithToken(token string) Option {
	return WithNatsOptions(nats.Token(token))
}

// WithTokenHandler is an Option to set the token handler to use
// when a token is not included directly in the URLs
// and when a token is not set.
func WithTokenHandler(cb nats.AuthTokenHandler) Option {
	return WithNatsOptions(nats.TokenHandler(cb))
}

// WithUserCredentials is a convenience function that takes a filename
// for a user's JWT and a filename for the user's private Nkey seed.
func WithUserCredentials(userOrChainedFile string, seedFiles ...string) Option {
	return WithNatsOptions(nats.UserCredentials(userOrChainedFile, seedFiles...))
}

// WithUserJWT will set the callbacks to retrieve the user's JWT and
// the signature callback to sign the server nonce. This an the Nkey
// option are mutually exclusive.
func WithUserJWT(userCB nats.UserJWTHandler, sigCB nats.SignatureHandler) Option {
	return WithNatsOptions(nats.UserJWT(userCB, sigCB))
}

// WithNkey will set the public Nkey and the signature callback to
// sign the server nonce.
func WithNkey(pubKey string, sigCB nats.SignatureHandler) Option {
	return WithNatsOptions(nats.Nkey(pubKey, sigCB))
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
func WithNatsOptions(natsOpts ...nats.Option) Option {
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
