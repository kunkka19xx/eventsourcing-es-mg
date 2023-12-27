package serializer

import (
	"example.com/main/cmd/eventstore/aggregate"
	"example.com/main/cmd/eventstore/event"
)

type Serializer interface {
	SerializerEvent(aggregate aggregate.Aggregate, event event.EventModel) (event.EventModel, error)
	DeserializerEvent(event event.EventModel) (event.EventModel, error)
}
