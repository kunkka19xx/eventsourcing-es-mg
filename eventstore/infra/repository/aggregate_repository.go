package infra

import (
	"context"
	"eventstore-intro/eventstore/aggregate"
	"eventstore-intro/eventstore/event"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"log"
	"math"
	"os"
)

type aggregateStore struct {
	db *esdb.Client
	//serializer serializer.Serializer
}

func NewAggregateStore(db *esdb.Client) *aggregateStore {
	return &aggregateStore{
		db: db,
	}
}

var tracer = otel.Tracer("event-store")

const count = math.MaxInt64

func (a *aggregateStore) Load(ctx context.Context, aggregate aggregate.Aggregate) error {
	ctx, span := tracer.Start(ctx, "aggregateStore.Load")
	defer span.End()

	kv := attribute.String("AggregateID", aggregate.GetID())
	span.SetAttributes(kv)

	stream, err := a.db.ReadStream(ctx, aggregate.GetID(), esdb.ReadStreamOptions{}, count)
	if err != nil {
		span.RecordError(err)
		return errors.Wrap(err, "db.ReadStream")
	}
	defer stream.Close()
	for {
		ev, err := stream.Recv()
		if errors.Is(err, esdb.ErrStreamNotFound) {
			span.RecordError(err)
			return errors.Wrap(err, "stream.Recv")
		}
		if errors.Is(err, os.ErrNotExist) {
			break
		}
		if err != nil {
			span.RecordError(err)
			return errors.Wrap(err, "stream.Recv")
		}

		esEvent := event.NewRecordedEvent(ev.Event)
		if err := aggregate.RaiseEvent(esEvent); err != nil {
			span.RecordError(err)
			return errors.Wrap(err, "RaiseEvent")
		}
		log.Printf("(Load) esEvent: {%s}", esEvent)
	}

	log.Printf("(Load) aggregate: {%s}", aggregate.String())
	return nil
}
