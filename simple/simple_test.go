package simple

import (
	"strings"
	"sync"
	"testing"

	"github.com/geniusrabbit/notificationcenter"
	"github.com/geniusrabbit/notificationcenter/encoder"
)

func Test_SimpleSubscriber(t *testing.T) {
	type m struct {
		Msg string `json:"msg"`
	}

	var (
		wg     sync.WaitGroup
		chanel = make(chan string, 1)
		tunel  = NewSimple(encoder.JSON)
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

	go func() {
		tunel.Listen()
	}()

	// Add message handler & test processing
	tunel.Subscribe(notificationcenter.FuncHandler(func(msg notificationcenter.Message) error {
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
		if err := tunel.Send(test.message); err != nil {
			t.Error(err)
		}
	}

	tunel.Close()
	wg.Wait()
}
