package notificationcenter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type testPublisher struct{}

func (p *testPublisher) Publish(ctx context.Context, messages ...any) error { return nil }

type testSubscriber struct{}

func (p *testSubscriber) Subscribe(ctx context.Context, receiver Receiver) error { return nil }
func (p *testSubscriber) Listen(ctx context.Context) error                       { return nil }
func (p *testSubscriber) Close() error                                           { return nil }

type registryTestSuite struct {
	suite.Suite
	// reg *Registry
}

func (suite *registryTestSuite) SetupSuite() {}

func (suite *registryTestSuite) TearDownSuite() {
	suite.NoError(Close())
}

func (suite *registryTestSuite) TestRegister() {
	err := Register(
		"test1", &testPublisher{},
		"test2", &testSubscriber{},
	)
	suite.NoError(err)
}

func (suite *registryTestSuite) TestNullPublisher() {
	suite.Nil(PublisherByName("undefined"))
}

func (suite *registryTestSuite) TestNullSubsciber() {
	suite.Nil(SubscriberByName("undefined"))
}

func (suite *registryTestSuite) TestPublish() {
	err := Register("test-pub-1", &testPublisher{})
	suite.NoError(err)
	err = Publish(context.TODO(), "test-pub-1", "hi")
	suite.NoError(err)
}

func (suite *registryTestSuite) TestSubscribe() {
	err := Register("test-sub-1", &testSubscriber{})
	suite.NoError(err)
	err = Subscribe(context.TODO(), "test-sub-1",
		FuncReceiver(func(msg Message) error { return nil }))
	suite.NoError(err)
}

func TestRegister(t *testing.T) {
	suite.Run(t, new(registryTestSuite))
}
