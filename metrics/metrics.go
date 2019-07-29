//
// @project geniusrabbit::notificationcenter 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package metrics

import (
	"io"

	"github.com/geniusrabbit/notificationcenter"
)

// Metric types
type (
	TypeIncrement []string
	TypeGauge     map[string]interface{}
	TypeTiming    map[string]int
	TypeCount     map[string]int
	TypeUnique    map[string]string
)

type metricas interface {
	SendMetricas(metricas ...interface{}) error
}

type wrapper struct {
	formater Formater
	logger   metricas
}

// NewWrapper of metrica
func NewWrapper(mm metricas, format Formater) notificationcenter.Streamer {
	return wrapper{
		formater: format,
		logger:   mm,
	}
}

// Send data to statistic
func (w wrapper) Send(messages ...interface{}) (err error) {
	for _, msg := range messages {
		if w.formater == nil {
			err = w.logger.SendMetricas(msg)
		} else {
			err = w.logger.SendMetricas(w.formater.Format(msg))
		}
		if err != nil {
			break
		}
	}
	return
}

// Close metrica target
func (w wrapper) Close() error {
	if w.logger != nil {
		if cl, _ := w.logger.(io.Closer); cl != nil {
			return cl.Close()
		}
	}
	return nil
}
