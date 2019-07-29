//
// @project geniusrabbit::notificationcenter 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//
// DEPRECATED: Will be removed in the next iteration

package statsd

import (
	"fmt"
	"time"

	statsd "gopkg.in/alexcesaro/statsd.v2"

	"github.com/geniusrabbit/notificationcenter"
	"github.com/geniusrabbit/notificationcenter/metrics"
)

// StatsD client
type StatsD statsd.Client

// NewUDP makes statsd instance with specific params
func NewUDP(addr string, format metrics.Formater, opts ...statsd.Option) (notificationcenter.Streamer, error) {
	var (
		options = []statsd.Option{
			statsd.Address(addr),
			statsd.Network("udp"),
		}
	)

	if len(opts) > 0 {
		options = append(options, opts...)
	}

	c, err := statsd.New(options...)
	if err != nil {
		return nil, err
	}

	return metrics.NewWrapper((*StatsD)(c), format), err
}

// SendMetricas creates request with metrics
func (s *StatsD) SendMetricas(messages ...interface{}) (err error) {
	// Working without panics
	defer func() { recover() }()

	for _, m := range messages {
		switch msg := m.(type) {
		case string:
			s.client().Increment(msg)
		case []string:
			for _, m := range msg {
				s.client().Increment(m)
			}
		case map[string]int:
			for bucket, count := range msg {
				s.client().Count(bucket, count)
			}
		case metrics.TypeCount:
			for bucket, count := range msg {
				s.client().Count(bucket, count)
			}
		case map[string]func() error:
			for bucket, handler := range msg {
				t := s.client().NewTiming()
				if err = handler(); err != nil {
					break
				}
				t.Send(bucket)
			}
		case metrics.TypeTiming:
			for bucket, duration := range msg {
				s.client().Timing(bucket, duration)
			}
		case map[string]time.Duration:
			for bucket, duration := range msg {
				s.client().Timing(bucket, int(duration/time.Millisecond))
			}
		case metrics.TypeGauge:
			for bucket, value := range msg {
				s.client().Gauge(bucket, value)
			}
		case map[string]string:
			for bucket, value := range msg {
				s.client().Unique(bucket, value)
			}
		case metrics.TypeUnique:
			for bucket, value := range msg {
				s.client().Unique(bucket, value)
			}
		default:
			err = fmt.Errorf("Unexpected type of metric: %s", m)
		}
	}

	return
}

// Close closes client connection
func (s *StatsD) Close() error {
	s.client().Close()
	return nil
}

///////////////////////////////////////////////////////////////////////////////
/// Internal methods
///////////////////////////////////////////////////////////////////////////////

func (s *StatsD) client() *statsd.Client {
	return (*statsd.Client)(s)
}
