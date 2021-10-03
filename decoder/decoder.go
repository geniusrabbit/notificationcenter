package decoder

import (
	"encoding/json"
	"encoding/xml"
)

// Decoder function type
type Decoder func(data []byte, msg interface{}) error

// JSON decoder implementation
func JSON(data []byte, msg interface{}) error {
	return json.Unmarshal(data, msg)
}

// XML decoder implementation
func XML(data []byte, msg interface{}) error {
	return xml.Unmarshal(data, msg)
}
