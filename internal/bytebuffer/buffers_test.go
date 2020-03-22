package bytebuffer

import "testing"

func TestPool(t *testing.T) {
	for i := 0; i < 100; i++ {
		buff := AcquireBuffer()
		if buff == nil {
			t.Error(`invalid buffer allocation`)
		}
		ReleaseBuffer(buff)
		ReleaseBuffer(nil)
	}
}
