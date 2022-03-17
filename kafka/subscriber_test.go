package kafka

import (
	"strconv"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"

	nc "github.com/geniusrabbit/notificationcenter/v2"
)

const (
	testTopicName         = `my_topic`
	testConsumerGroupName = `default`
)

func TestSubscriberReceiveMessages(t *testing.T) {
	testMsg := sarama.StringEncoder("Foo")
	broker0 := sarama.NewMockBroker(t, 0)
	messageCount := 10

	mockFetchResponse := sarama.NewMockFetchResponse(t, 1)
	for i := 0; i < messageCount; i++ {
		mockFetchResponse.SetMessage(testTopicName, 0, int64(i+1234), testMsg)
	}

	broker0.SetHandlerByMap(map[string]sarama.MockResponse{
		"MetadataRequest": sarama.NewMockMetadataResponse(t).
			SetBroker(broker0.Addr(), broker0.BrokerID()).
			SetLeader(testTopicName, 0, broker0.BrokerID()),
		"FindCoordinatorRequest": sarama.NewMockFindCoordinatorResponse(t).
			SetCoordinator(sarama.CoordinatorGroup, testConsumerGroupName, broker0),
		"JoinGroupRequest": sarama.NewMockWrapper(&sarama.JoinGroupResponse{
			Version:       1,
			GroupProtocol: `consumer`,
			MemberId:      "OneProtocol",
			LeaderId:      strconv.Itoa(int(broker0.BrokerID())),
		}),
		"SyncGroupRequest":  sarama.NewMockWrapper(&sarama.SyncGroupResponse{}),
		"LeaveGroupRequest": sarama.NewMockWrapper(&sarama.LeaveGroupResponse{}),
		"HeartbeatRequest":  sarama.NewMockWrapper(&sarama.HeartbeatResponse{}),
		"OffsetRequest": sarama.NewMockOffsetResponse(t).
			SetOffset(testTopicName, 0, sarama.OffsetOldest, 1234).
			SetOffset(testTopicName, 0, sarama.OffsetNewest, 1234),
		"OffsetFetchRequest": sarama.NewMockOffsetFetchResponse(t).
			SetOffset(testConsumerGroupName, testTopicName, 0, 1234, ``, sarama.ErrNoError),
		"FetchRequest": mockFetchResponse,
	})

	defer broker0.Close()

	// Connect to the mock broker
	sub, err := NewSubscriber(
		WithBrokers(broker0.Addr()),
		WithTopics(testTopicName),
		WithGroupName(testConsumerGroupName),
	)
	assert.NoError(t, err, `new subscriber`)

	err = sub.Close()
	assert.NoError(t, err, `close connection`)
}

func TestNewSubscriberError(t *testing.T) {
	_, err := NewSubscriber(
		WithGroupName(``),
		WithKafkaVersion(sarama.V0_10_0_0),
		WithErrorHandler(func(msg nc.Message, err error) {}),
		WithPanicHandler(func(msg nc.Message, recoverData any) {}),
	)
	assert.Error(t, err)
}
