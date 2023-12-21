package serializer

import (
	"eventstore-intro/eventstore/aggregate"
	"eventstore-intro/eventstore/event"
)

type Serializer interface {
	SerializerEvent(aggregate aggregate.Aggregate, event event.EventModel) (event.EventModel, error)
	DeserializerEvent(event event.EventModel) (event.EventModel, error)
}
