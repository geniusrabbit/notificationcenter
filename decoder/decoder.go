package decoder

import (
	"encoding/json"
	"encoding/xml"
)

// Decoder function type
type Decoder func(data []byte, msg any) error

// JSON decoder implementation
func JSON(data []byte, msg any) error {
	return json.Unmarshal(data, msg)
}

// XML decoder implementation
func XML(data []byte, msg any) error {
	return xml.Unmarshal(data, msg)
}
