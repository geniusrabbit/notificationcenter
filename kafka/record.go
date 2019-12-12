//
// @project geniusrabbit.com 2015 – 2016, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015 – 2016, 2019
//

package kafka

import "sync"

const minimalMessageSize = 1000

var byteEncoderPool = sync.Pool{
	New: func() interface{} {
		return make(kafkaByteEncoder, 0, minimalMessageSize)
	},
}

type kafkaByteEncoder []byte

func (k kafkaByteEncoder) Encode() ([]byte, error) {
	return k, nil
}

func (k kafkaByteEncoder) Length() int {
	return len(k)
}

func (k kafkaByteEncoder) Release() {
	byteEncoderPool.Put(k[:0])
}

func byteEncoder(data []byte) kafkaByteEncoder {
	byteEncoder := byteEncoderPool.New().(kafkaByteEncoder)
	return append(byteEncoder, data...)
}
