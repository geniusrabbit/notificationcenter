//
// @project geniusrabbit.com 2015 – 2016
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015 – 2016
//

package kafka

type kafkaByteEncoder struct {
	data []byte
}

func (k *kafkaByteEncoder) Encode() ([]byte, error) {
	return k.data, nil
}

func (k *kafkaByteEncoder) Length() int {
	if nil == k || nil == k.data {
		return 0
	}
	return len(k.data)
}

func byteEncoder(data []byte) *kafkaByteEncoder {
	return &kafkaByteEncoder{data: data}
}
