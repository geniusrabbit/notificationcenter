package kafka

import (
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

	// Name of the subscription group
	GroupName string

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

// WithClientID value
func WithClientID(clientID string) Option {
	return func(options *Options) {
		options.ClusterConfig.Config.ClientID = clientID
	}
}

// WithGroupName of the message consuming
func WithGroupName(name string) Option {
	return func(options *Options) {
		options.GroupName = name
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
