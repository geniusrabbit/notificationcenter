package decoder

import (
	"encoding/xml"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testItem struct {
	XMLName xml.Name `json:"-" xml:"item"`
	Value   string   `json:"value" xml:"value"`
}

func TestDecoding(t *testing.T) {
	tests := []struct {
		target any
		data   string
		dec    Decoder
	}{
		{
			target: &testItem{Value: "target"},
			data:   `{"value":"target"}`,
			dec:    JSON,
		},
		{
			target: &testItem{XMLName: xml.Name{Local: "item"}, Value: "target"},
			data:   `<item><value>target</value></item>`,
			dec:    XML,
		},
	}

	for _, test := range tests {
		var it testItem
		err := test.dec([]byte(test.data), &it)
		assert.NoError(t, err)
		assert.True(t, reflect.DeepEqual(test.target, &it))
	}
}
