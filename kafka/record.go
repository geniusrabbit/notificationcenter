//
// @project geniusrabbit.com 2015 – 2016, 2019 – 2022, 2025
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015 – 2016, 2019 – 2022, 2025
//

package kafka

// Note: The encoder is intentionally not pooled to avoid data races with Sarama's
// asynchronous producer which may continue to read the message value on
// background goroutines even after a Success notification is delivered.
// The message value must remain immutable for the lifetime of the ProducerMessage.

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
	// no-op: keep data immutable and let GC reclaim memory when no longer referenced
}

func byteEncoder(data []byte) *kafkaByteEncoder {
	// Copy data to ensure immutability and avoid sharing underlying array
	cp := make([]byte, len(data))
	copy(cp, data)
	return &kafkaByteEncoder{data: cp}
}
