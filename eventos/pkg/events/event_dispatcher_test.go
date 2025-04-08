package events

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestEvent struct {
	Name    string
	Payload interface{}
}

func (e *TestEvent) GetName() string {
	return e.Name
}

func (e *TestEvent) GetPayload() interface{} {
	return e.Payload
}

func (e *TestEvent) GetDateTime() time.Time {
	return time.Now()
}

type TestEventHandler struct {
	ID int
}

func (h *TestEventHandler) Handle(event EventInterface, wg *sync.WaitGroup) {

}

type EventDispatcherTestSuit struct {
	suite.Suite
	event           TestEvent
	event2          TestEvent
	handler1        TestEventHandler
	handler2        TestEventHandler
	handler3        TestEventHandler
	eventDispatcher *EventDispatcher
}

func (s *EventDispatcherTestSuit) SetupTest() {
	s.event = TestEvent{Name: "test", Payload: "test"}
	s.event2 = TestEvent{Name: "test2", Payload: "test2"}
	s.handler1 = TestEventHandler{ID: 1}
	s.handler2 = TestEventHandler{ID: 2}
	s.handler3 = TestEventHandler{ID: 3}
	s.eventDispatcher = NewEventDispatcher()
}

func (suit *EventDispatcherTestSuit) TestEventDispatcher_Register() {
	err := suit.eventDispatcher.Register(suit.event.GetName(), &suit.handler1)
	suit.Nil(err)
	suit.Equal(1, len(suit.eventDispatcher.handlers[suit.event.GetName()]))

	err = suit.eventDispatcher.Register(suit.event.GetName(), &suit.handler2)
	suit.Nil(err)
	suit.Equal(2, len(suit.eventDispatcher.handlers[suit.event.GetName()]))

	assert.Equal(suit.T(), &suit.handler1, suit.eventDispatcher.handlers[suit.event.GetName()][0])
	assert.Equal(suit.T(), &suit.handler2, suit.eventDispatcher.handlers[suit.event.GetName()][1])
}

func (suit *EventDispatcherTestSuit) TestEventDispatcher_Register_WithSameHandler() {
	err := suit.eventDispatcher.Register(suit.event.GetName(), &suit.handler1)
	suit.Nil(err)
	suit.Equal(1, len(suit.eventDispatcher.handlers[suit.event.GetName()]))

	err = suit.eventDispatcher.Register(suit.event.GetName(), &suit.handler1)
	suit.Equal(ErrHandlerAlreadyRegistered, err)
}

func (suit *EventDispatcherTestSuit) TestEventDispatcher_Clear() {
	//Event 1
	err := suit.eventDispatcher.Register(suit.event.GetName(), &suit.handler1)
	suit.Nil(err)
	err = suit.eventDispatcher.Register(suit.event.GetName(), &suit.handler2)
	suit.Nil(err)

	//Event 2
	err = suit.eventDispatcher.Register(suit.event2.GetName(), &suit.handler3)
	suit.Nil(err)

	err = suit.eventDispatcher.Clear()
	suit.Nil(err)

	suit.Equal(0, len(suit.eventDispatcher.handlers[suit.event.GetName()]))
	suit.Equal(0, len(suit.eventDispatcher.handlers[suit.event2.GetName()]))
}

func (suit *EventDispatcherTestSuit) TestEventDispatcher_Has() {
	err := suit.eventDispatcher.Register(suit.event.GetName(), &suit.handler1)
	suit.Nil(err)
	suit.Equal(1, len(suit.eventDispatcher.handlers[suit.event.GetName()]))

	err = suit.eventDispatcher.Register(suit.event.GetName(), &suit.handler2)
	suit.Nil(err)
	suit.Equal(2, len(suit.eventDispatcher.handlers[suit.event.GetName()]))

	assert.True(suit.T(), suit.eventDispatcher.Has(suit.event.GetName(), &suit.handler1))
	assert.True(suit.T(), suit.eventDispatcher.Has(suit.event.GetName(), &suit.handler2))

	assert.False(suit.T(), suit.eventDispatcher.Has(suit.event.GetName(), &suit.handler3))
}

type MockHandler struct {
	mock.Mock
}

// mock tem de implementar a interface
func (m *MockHandler) Handle(event EventInterface, wg *sync.WaitGroup) {
	m.Called(event)
	wg.Done()
}

func (suit *EventDispatcherTestSuit) TestEventDispatcher_Dispatch() {
	eh := &MockHandler{}
	eh.On("Handle", &suit.event)

	eh2 := &MockHandler{}
	eh2.On("Handle", &suit.event)

	suit.eventDispatcher.Register(suit.event.GetName(), eh)
	suit.eventDispatcher.Register(suit.event.GetName(), eh2)
	suit.eventDispatcher.Dispatch(&suit.event)
	eh.AssertExpectations(suit.T())
	eh2.AssertExpectations(suit.T())
	eh.AssertNumberOfCalls(suit.T(), "Handle", 1)
	eh2.AssertNumberOfCalls(suit.T(), "Handle", 1)
}

func (suit *EventDispatcherTestSuit) TestEventDispatcher_Remove() {
	err := suit.eventDispatcher.Register(suit.event.GetName(), &suit.handler1)
	suit.Nil(err)
	suit.Equal(1, len(suit.eventDispatcher.handlers[suit.event.GetName()]))

	err = suit.eventDispatcher.Register(suit.event.GetName(), &suit.handler2)
	suit.Nil(err)
	suit.Equal(2, len(suit.eventDispatcher.handlers[suit.event.GetName()]))

	suit.eventDispatcher.Remove(suit.event.GetName(), &suit.handler1)
	suit.Equal(1, len(suit.eventDispatcher.handlers[suit.event.GetName()]))
	assert.Equal(suit.T(), &suit.handler2, suit.eventDispatcher.handlers[suit.event.GetName()][0])

	suit.eventDispatcher.Remove(suit.event.GetName(), &suit.handler2)
	suit.Equal(0, len(suit.eventDispatcher.handlers[suit.event.GetName()]))

}

func TestSuite(t *testing.T) {
	suite.Run(t, new(EventDispatcherTestSuit))
}
