package kafka

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"

	nc "github.com/geniusrabbit/notificationcenter"
	"github.com/geniusrabbit/notificationcenter/encoder"
	"github.com/geniusrabbit/notificationcenter/internal/logger"
)

// Options for publisher or subscriber
type Options struct {
	ClusterConfig cluster.Config

	// IsSynchronous type of producer
	// TODO: make it work for sync publisher
	IsSynchronous bool

	// Brokers contains list of broker hosts with port
	Brokers []string

	// Name of the subscription group
	GroupName string

	// Names of topics for subscribing or publishing
	Topics []string

	// ErrorHandler of message processing
	ErrorHandler nc.ErrorHandler

	// PanicHandler process panic
	PanicHandler nc.PanicHandler

	// PublisherErrorHandler provides handler of message send errors
	PublisherErrorHandler PublisherErrorHandler

	// PublisherSuccessHandler provides handler of message send success
	PublisherSuccessHandler PublisherSuccessHandler

	// SubscriberNotificationHandler provides handler of received messages
	SubscriberNotificationHandler SubscriberNotificationHandler

	// Message encoder interface
	Encoder encoder.Encoder

	// Logger of subscriber
	Logger loggerInterface
}

func (opt *Options) clusterConfig() *cluster.Config {
	return &opt.ClusterConfig
}

func (opt *Options) saramaConfig() *sarama.Config {
	return &opt.ClusterConfig.Config
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

// Option function type
type Option func(options *Options)

// WithSaramaConfig custom config
func WithSaramaConfig(streamConfig *sarama.Config) Option {
	return func(options *Options) {
		options.ClusterConfig.Config = *streamConfig
	}
}

// WithKafkaURL is an Option to set the URL the client should connect to.
// The url can contain username/password semantics. e.g. kafka://derek:pass@localhost:4222/{groupName}?topics=topic1,topic2
// Comma separated arrays are also supported, e.g. urlA, urlB.
func WithKafkaURL(urlString string) Option {
	return func(opt *Options) {
		u, err := url.Parse(urlString)
		if err != nil {
			panic(err)
		}
		if len(u.Path) > 1 {
			opt.GroupName = u.Path[1:]
		}
		if isSync := u.Query().Get(`sync`); isSync != `` {
			opt.IsSynchronous, _ = strconv.ParseBool(isSync)
		}
		topics := strings.Split(u.Query().Get(`topics`), `,`)
		if len(topics) == 1 && topics[0] == `` {
			topics = nil
		}
		if len(topics) > 0 {
			opt.Topics = topics
		}
		opt.Brokers = strings.Split(u.Host, ",")
		if u.User != nil && u.User.Username() != `` {
			opt.clusterConfig().Net.SASL.User = u.User.Username()
			opt.clusterConfig().Net.SASL.Password = u.User.Username()
		}
	}
}

// WithClientID value
func WithClientID(clientID string) Option {
	return func(options *Options) {
		options.ClusterConfig.Config.ClientID = clientID
	}
}

// WithBrokers overrides the list of brokers to connect
func WithBrokers(brokers ...string) Option {
	return func(options *Options) {
		options.Brokers = brokers
	}
}

// WithGroupName of the message consuming
func WithGroupName(name string) Option {
	return func(options *Options) {
		options.GroupName = name
	}
}

// WithTopics will set the list of topics for publishing or subscribing
func WithTopics(topics ...string) Option {
	return func(opt *Options) {
		opt.Topics = topics
	}
}

// WithKafkaVersion minimal version
func WithKafkaVersion(version sarama.KafkaVersion) Option {
	return func(options *Options) {
		options.ClusterConfig.Version = version
	}
}

// WithFlashFrequency of flushing
func WithFlashFrequency(frequency time.Duration) Option {
	return func(options *Options) {
		options.ClusterConfig.Producer.Flush.Frequency = frequency
	}
}

// WithFlashMessages minimal count
func WithFlashMessages(messageCount int) Option {
	return func(options *Options) {
		options.ClusterConfig.Producer.Flush.Messages = messageCount
	}
}

// WithCompression of the pipe
func WithCompression(codec sarama.CompressionCodec, level int) Option {
	return func(options *Options) {
		options.ClusterConfig.Producer.Compression = codec
		options.ClusterConfig.Producer.CompressionLevel = level
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

// WithPublisherErrorHandler set handler of the sarama errors
func WithPublisherErrorHandler(h PublisherErrorHandler) Option {
	return func(options *Options) {
		options.ClusterConfig.Producer.Return.Errors = h != nil
		options.PublisherErrorHandler = h
	}
}

// WithPublisherSuccessHandler set handler of the sarama success
func WithPublisherSuccessHandler(h PublisherSuccessHandler) Option {
	return func(options *Options) {
		options.ClusterConfig.Producer.Return.Successes = h != nil
		options.PublisherSuccessHandler = h
	}
}

// WithSubscriberNotificationHandler set handler of the cluster group notifications
func WithSubscriberNotificationHandler(h SubscriberNotificationHandler) Option {
	return func(options *Options) {
		options.ClusterConfig.Group.Return.Notifications = h != nil
		options.SubscriberNotificationHandler = h
	}
}

// WithClusterConfig custom config
func WithClusterConfig(clusterConfig *cluster.Config) Option {
	return func(options *Options) {
		options.ClusterConfig = *clusterConfig
	}
}
