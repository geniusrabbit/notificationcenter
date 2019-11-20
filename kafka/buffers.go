package kafka

import (
	"bytes"
	"sync"
)

var buffers = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

func acquireBuffer() *bytes.Buffer {
	return buffers.Get().(*bytes.Buffer)
}

func releaseBuffer(buff *bytes.Buffer) {
	if buff != nil {
		buff.Reset()
		buffers.Put(buff)
	}
}
