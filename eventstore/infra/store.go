package infra

import (
	"context"
	"eventstore-intro/eventstore/aggregate"
	"eventstore-intro/eventstore/event"
)

type AggregateStore interface {
	//LoadAggregate gets latest version of an aggregate
	LoadAggregate(ctx context.Context, aggregate aggregate.Aggregate) error
	//SaveAggregate saves uncommitted events
	SaveAggregate(ctx context.Context, aggregate aggregate.Aggregate) error
	// IsAggregateExists check aggregate is existed or not by ID
	IsAggregateExisting(ctx context.Context, streamId string) error
}

type EventStore interface {
	//LoadEvents gets all events for the aggregate ID
	LoadEvents(ctx context.Context, streamId string) ([]event.EventModel, error)
	//SaveEvents saves all events to eventstoreDB
	SaveEvents(ctx context.Context, streamId string, events event.EventModel) error
}

type SnapshotStore interface {
	// SaveSnapshot save aggregate snapshot.
	SaveSnapshot(ctx context.Context, aggregate aggregate.Aggregate) error

	// GetSnapshot load aggregate snapshot.
	GetSnapshot(ctx context.Context, id string) (*aggregate.Snapshot, error)
}
