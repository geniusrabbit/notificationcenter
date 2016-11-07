//
// @project geniusrabbit.com 2016
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016
//

package statsd

import (
	"fmt"
	"time"

	statsd "gopkg.in/alexcesaro/statsd.v2"
)

// Metric types
type (
	TypeIncrement []string
	TypeGauge     map[string]int
	TypeTiming    map[string]int
	TypeCount     map[string]int
	TypeUnique    map[string]string
)

// StatsD client
type StatsD *statsd.Client

// NewUDP makes statsd instance with specific params
func NewUDP(addr string, opts ...statsd.Option) (*StatsD, error) {
	options := statsd.Option{
		statsd.Address(addr),
		statsd.Network("udp"),
	}

	if len(opts) > 0 {
		options = append(options, opts...)
	}

	c, err := statsd.New(options...)
	if err != nil {
		return nil, err
	}
	return (*StatsD)(c), nil
}

// Send creates request with metrics
func (s *StatsD) Send(messages ...interface{}) (err error) {
	// Working without panics
	defer func() { recover() }()

	for _, m := range messages {
		switch msg := m.(type) {
		case string:
			s.client().Increment(msg)
		case []string:
			for _, s := range msg {
				s.client().Increment(s)
			}
		case map[string]int, TypeCount:
			for bucket, count := range msg {
				s.client.Count(bucket, count)
			}
		case map[string]func() error:
			for bucket, handler := range msg {
				t := s.client.NewTiming()
				if err = handler(); err != nil {
					break
				}
				t.Send(bucket)
			}
		case TypeTiming:
			for bucket, duration := range msg {
				s.client.Timing(bucket, duration)
			}
		case map[string]time.Duration:
			for bucket, duration := range msg {
				s.client.Timing(bucket, int(duration/time.Millisecond))
			}
		case TypeGauge:
			for bucket, value := range msg {
				s.client.Gauge(bucket, value)
			}
		case map[string]string, TypeUnique:
			for bucket, value := range msg {
				s.client.Unique(bucket, value)
			}
		default:
			err = fmt.Errorf("Unexpected type of metric: %s", m)
		}
	}

	return
}

// Close closes client connection
func (s *StatsD) Close() {
	s.client().Close()
}

///////////////////////////////////////////////////////////////////////////////
/// Internal methods
///////////////////////////////////////////////////////////////////////////////

func (s *StatsD) client() *statsd.Client {
	return s.(*statsd.Client)
}
