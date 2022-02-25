//
// @project geniusrabbit.com 2015 – 2016, 2019 – 2022
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015 – 2016, 2019 – 2022
//

package kafka

import "sync"

const minimalMessageSize = 1000

var byteEncoderPool = sync.Pool{
	New: func() any {
		return &kafkaByteEncoder{data: make([]byte, 0, minimalMessageSize)}
	},
}

type kafkaByteEncoder struct {
	data []byte
}

func (k *kafkaByteEncoder) Encode() ([]byte, error) {
	return k.data, nil
}

func (k *kafkaByteEncoder) Length() int {
	return len(k.data)
}

func (k *kafkaByteEncoder) Release() {
	k.data = k.data[:0]
	byteEncoderPool.Put(k)
}

func byteEncoder(data []byte) *kafkaByteEncoder {
	byteEncoder := byteEncoderPool.New().(*kafkaByteEncoder)
	byteEncoder.data = append(byteEncoder.data, data...)
	return byteEncoder
}
