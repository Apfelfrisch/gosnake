package client

type EventListener func(event Event)

type Event interface{}

type GameHasStarted struct{}
type GameHasEnded struct{}
type GameWas struct{}
type PlayerDashed struct{}
type PlayerCrashed struct{}
type PlayerHasEaten struct{}
type PlayerWalkedWall struct{}

type EventBus struct {
	lst map[Event][]EventListener
}

func NewEventBus() *EventBus {
	return &EventBus{make(map[Event][]EventListener)}
}

func (m *EventBus) Add(e Event, l EventListener) {
	m.lst[e] = append(m.lst[e], l)
}

func (m *EventBus) Dispatch(e Event) {
	for _, listener := range m.lst[e] {
		listener(e)
	}
}
