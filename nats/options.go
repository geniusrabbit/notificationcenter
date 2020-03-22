package nats

import (
	nats "github.com/nats-io/nats.go"

	nc "github.com/geniusrabbit/notificationcenter"
	"github.com/geniusrabbit/notificationcenter/encoder"
	"github.com/geniusrabbit/notificationcenter/internal/logger"
)

// Options of the NATS wrapper
type Options struct {
	// Inited NATS connection
	Conn *nats.Conn

	// ConnURL of connection
	ConnURL string

	// Raw options from the "github.com/nats-io/nats.go" module
	NatsOptions []nats.Option

	// Name of the subscription group
	GroupName string

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
// The url can contain username/password semantics. e.g. nats://derek:pass@localhost:4222
// Comma separated arrays are also supported, e.g. urlA, urlB.
func WithNatsURL(url string) Option {
	return func(opt *Options) {
		opt.ConnURL = url
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
