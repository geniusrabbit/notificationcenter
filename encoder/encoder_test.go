package encoder

import (
	"bytes"
	"encoding/xml"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testItem struct {
	XMLName xml.Name `json:"-" xml:"item"`
	Value   string   `json:"value" xml:"value"`
}

func TestEncoding(t *testing.T) {
	tests := []struct {
		obj    any
		target string
		enc    Encoder
	}{
		{
			obj:    testItem{Value: "target"},
			target: `{"value":"target"}`,
			enc:    JSON,
		},
		{
			obj:    testItem{Value: "target"},
			target: `<item><value>target</value></item>`,
			enc:    XML,
		},
	}

	var buff bytes.Buffer
	for _, test := range tests {
		buff.Reset()
		err := test.enc(test.obj, &buff)
		assert.NoError(t, err)
		assert.Equal(t, test.target, strings.TrimSpace(buff.String()))
	}
}
