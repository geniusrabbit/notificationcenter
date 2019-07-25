//
// @project geniusrabbit.com 2015 – 2016, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015 – 2016, 2019
//

package kafka

type kafkaByteEncoder []byte

func (k kafkaByteEncoder) Encode() ([]byte, error) {
	return k, nil
}

func (k kafkaByteEncoder) Length() int {
	return len(k)
}

func byteEncoder(data []byte) kafkaByteEncoder {
	return kafkaByteEncoder(data)
}
