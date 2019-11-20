package kafka

import (
	"time"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
)

// OptionStream function type
type OptionStream func(conf *sarama.Config)

// WithStreamSaramaConfig custom config
func WithStreamSaramaConfig(streamConfig *sarama.Config) OptionStream {
	return func(conf *sarama.Config) {
		*conf = *streamConfig
	}
}

// WithStreamClientID value
func WithStreamClientID(clientID string) OptionStream {
	return func(conf *sarama.Config) {
		conf.ClientID = clientID
	}
}

// WithStreamKafkaVersion minimal version
func WithStreamKafkaVersion(version sarama.KafkaVersion) OptionStream {
	return func(conf *sarama.Config) {
		conf.Version = version
	}
}

// WithStreamFlashFrequency of flushing
func WithStreamFlashFrequency(frequency time.Duration) OptionStream {
	return func(conf *sarama.Config) {
		conf.Producer.Flush.Frequency = frequency
	}
}

// WithStreamFlashMessages minimal count
func WithStreamFlashMessages(messageCount int) OptionStream {
	return func(conf *sarama.Config) {
		conf.Producer.Flush.Messages = messageCount
	}
}

// WithStreamCompression of the pipe
func WithStreamCompression(codec sarama.CompressionCodec, level int) OptionStream {
	return func(conf *sarama.Config) {
		conf.Producer.Compression = codec
		conf.Producer.CompressionLevel = level
	}
}

// OptionSubscriber function type
type OptionSubscriber func(conf *cluster.Config)

// WithSubscriberClusterConfig custom config
func WithSubscriberClusterConfig(clusterConfig *cluster.Config) OptionSubscriber {
	return func(conf *cluster.Config) {
		*conf = *clusterConfig
	}
}
