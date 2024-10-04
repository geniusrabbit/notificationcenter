package redis

import (
	"crypto/tls"
	"net/url"
	"strings"
	"time"

	"github.com/demdxx/gocast/v2"
	"github.com/redis/go-redis/v9"

	nc "github.com/geniusrabbit/notificationcenter/v2"
	"github.com/geniusrabbit/notificationcenter/v2/encoder"
	"github.com/geniusrabbit/notificationcenter/v2/internal/logger"
)

type loggerInterface interface {
	Error(params ...any)
	Debugf(msg string, params ...any)
}

// Option of the Redis subscriber or publisher
type Option func(opt *Options)

// Options of the Redis wrapper
type Options struct {
	// Client of the redis server
	Client *redis.Client

	// Channels list
	Channels []string

	// RedisOptions to connect to the client
	RedisOptions *redis.Options

	// ErrorHandler of message processing
	ErrorHandler nc.ErrorHandler

	// PanicHandler process panic
	PanicHandler nc.PanicHandler

	// Message encoder interface
	Encoder encoder.Encoder

	// Logger of subscriber
	Logger loggerInterface
}

func (opt *Options) connect() *redis.Client {
	if opt.Client != nil {
		return opt.Client
	}
	opt.Client = redis.NewClient(opt.RedisOptions)
	return opt.Client
}

func (opt *Options) encoder() encoder.Encoder {
	if opt.Encoder == nil {
		return encoder.JSON
	}
	return opt.Encoder
}

func (opt *Options) channels() []string {
	if len(opt.Channels) == 0 {
		return []string{`default`}
	}
	return opt.Channels
}

func (opt *Options) logger() loggerInterface {
	if opt.Logger == nil {
		return logger.DefaultLogger
	}
	return opt.Logger
}

// WithRedisURL is an Option to set the URL the client should connect to.
// The url can contain username/password semantics. e.g. redis://derek:pass@localhost:4222/databaseNumber?topics=t1,t2
func WithRedisURL(urlString string) Option {
	return func(opt *Options) {
		u, err := url.Parse(urlString)
		if err != nil {
			panic(err)
		}
		dbNum := 0
		if len(u.Path) > 1 {
			dbNum = gocast.Number[int](u.Path[1:])
		}
		query := u.Query()
		channels := strings.Split(query.Get(`topics`), `,`)
		if len(channels) == 1 && channels[0] == `` {
			channels = nil
		}
		if len(channels) > 0 {
			opt.Channels = channels
		}
		password, _ := u.User.Password()
		network := "tcp"
		if u.Scheme == "unix" {
			network = "unix"
		}
		opt.RedisOptions = &redis.Options{
			Network:         network,
			Addr:            u.Host,
			Username:        u.User.Username(),
			Password:        password,
			DB:              dbNum,
			ClientName:      query.Get("client_name"),
			MaxRetries:      gocast.Number[int](query.Get("max_retries")),
			MinRetryBackoff: parseDuration(query.Get("min_retry_backoff")),
			MaxRetryBackoff: parseDuration(query.Get("max_retry_backoff")),
			DialTimeout:     parseDuration(query.Get("dial_timeout")),
			ReadTimeout:     parseDuration(query.Get("read_timeout")),
			WriteTimeout:    parseDuration(query.Get("write_timeout")),
			PoolFIFO:        gocast.Bool(query.Get("pool_fifo")),
			PoolSize:        gocast.Number[int](query.Get("pool_size")),
			PoolTimeout:     parseDuration(query.Get("pool_timeout")),
			MinIdleConns:    gocast.Number[int](query.Get("min_idle_conns")),
			MaxIdleConns:    gocast.Number[int](query.Get("max_idle_conns")),
		}

		if query.Has("conn_max_idle_time") {
			opt.RedisOptions.ConnMaxIdleTime = parseDuration(query.Get("conn_max_idle_time"))
		} else {
			opt.RedisOptions.ConnMaxIdleTime = parseDuration(query.Get("idle_timeout"))
		}
		if query.Has("conn_max_lifetime") {
			opt.RedisOptions.ConnMaxLifetime = parseDuration(query.Get("conn_max_lifetime"))
		} else {
			opt.RedisOptions.ConnMaxLifetime = parseDuration(query.Get("max_conn_age"))
		}

		if u.Scheme == "rediss" {
			opt.RedisOptions.TLSConfig = &tls.Config{
				ServerName: u.Hostname(),
				MinVersion: tls.VersionTLS12,
			}
		}
	}
}

// WithRedisOpts set redis options to connect
func WithRedisOpts(opts *redis.Options) Option {
	return func(opt *Options) {
		opt.RedisOptions = opts
	}
}

// WithRedisClient puts redis connect option
func WithRedisClient(cli *redis.Client) Option {
	return func(opt *Options) {
		opt.Client = cli
	}
}

// WithTopics of the pubsub stream
func WithTopics(topics ...string) Option {
	return func(opt *Options) {
		opt.Channels = topics
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

func parseDuration(s string) time.Duration {
	d, _ := time.ParseDuration(s)
	return d
}
