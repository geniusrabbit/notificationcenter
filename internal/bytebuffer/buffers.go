package bytebuffer

import (
	"bytes"
	"sync"
)

var buffers = sync.Pool{New: func() any { return &bytes.Buffer{} }}

// AcquireBuffer from the pool
func AcquireBuffer() *bytes.Buffer {
	return buffers.Get().(*bytes.Buffer)
}

// ReleaseBuffer and return back to the pool
func ReleaseBuffer(buff *bytes.Buffer) {
	if buff != nil {
		buff.Reset()
		buffers.Put(buff)
	}
}
