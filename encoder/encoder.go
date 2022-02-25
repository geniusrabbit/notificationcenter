package encoder

import (
	"encoding/json"
	"encoding/xml"
	"io"
)

// Encoder function type
type Encoder func(msg any, wr io.Writer) error

// JSON encoder implementation
func JSON(msg any, wr io.Writer) error {
	enc := json.NewEncoder(wr)
	enc.SetEscapeHTML(false)
	return enc.Encode(msg)
}

// XML encoder implementation
func XML(msg any, wr io.Writer) error {
	enc := xml.NewEncoder(wr)
	return enc.Encode(msg)
}
