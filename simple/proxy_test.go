package simple

import (
	"strings"
	"sync"
	"testing"

	nc "github.com/geniusrabbit/notificationcenter"
	"github.com/geniusrabbit/notificationcenter/encoder"
	"github.com/stretchr/testify/assert"
)

func Test_ProxySubscriber(t *testing.T) {
	type m struct {
		Msg string `json:"msg"`
	}

	var (
		wg     sync.WaitGroup
		chanel = make(chan string, 1)
		proxy  = NewProxy(encoder.JSON)
		tests  = []struct {
			message m
			result  string
		}{
			{
				message: m{Msg: "Test"},
				result:  `{"msg":"Test"}`,
			},
		}
	)

	assert.NoError(t, proxy.Listen())

	// Add message handler & test processing
	proxy.Subscribe(nc.FuncHandler(func(msg nc.Message) error {
		result := <-chanel
		data := strings.TrimSpace(string(msg.Body()))
		if data != result {
			t.Errorf("invalid message encode: `%s` != `%s`", data, result)
		}
		wg.Done()
		return nil
	}))

	for _, test := range tests {
		chanel <- test.result
		wg.Add(1)
		if err := proxy.Send(test.message); err != nil {
			t.Error(err)
		}
	}

	proxy.Close()
	wg.Wait()
}
