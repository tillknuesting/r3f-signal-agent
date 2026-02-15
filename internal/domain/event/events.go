package event

type Event interface {
	Type() string
	EventTimestamp() string
}

type EventDispatcher interface {
	Subscribe(eventType string, handler EventHandler)
	Dispatch(event Event)
}

type EventHandler func(event Event)

type TrendCollectedEvent struct {
	TrendID   string
	SourceID  string
	Timestamp string
}

func (e *TrendCollectedEvent) Type() string           { return "trend.collected" }
func (e *TrendCollectedEvent) EventTimestamp() string { return e.Timestamp }

type CollectionStartedEvent struct {
	JobID     string
	Profile   string
	Timestamp string
}

func (e *CollectionStartedEvent) Type() string           { return "collection.started" }
func (e *CollectionStartedEvent) EventTimestamp() string { return e.Timestamp }

type CollectionCompletedEvent struct {
	JobID      string
	ItemsCount int
	Timestamp  string
}

func (e *CollectionCompletedEvent) Type() string           { return "collection.completed" }
func (e *CollectionCompletedEvent) EventTimestamp() string { return e.Timestamp }
