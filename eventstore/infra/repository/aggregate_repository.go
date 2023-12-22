package infra

import (
	"context"
	"eventstore-intro/eventstore/aggregate"
	"eventstore-intro/eventstore/config"
	"eventstore-intro/eventstore/event"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"io"
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

func makeTracer(ctx context.Context, spanName, key, value string) trace.Span {
	ctx, span := tracer.Start(ctx, spanName)
	defer span.End()

	kv := attribute.String(key, value)
	span.SetAttributes(kv)
	return span
}

func (a *aggregateStore) LoadAggregate(ctx context.Context, aggregate aggregate.Aggregate) error {
	span := makeTracer(ctx, "aggregateStore.Load", "AggregateID", aggregate.GetID())
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
		esEvent := event.NewEventFromRecorded(ev.Event)
		if err := aggregate.RaiseEvent(esEvent); err != nil {
			span.RecordError(err)
			return errors.Wrap(err, "RaiseEvent")
		}
		log.Printf("(Load) esEvent: {%s}", esEvent)
	}

	log.Printf("(Load) aggregate: {%s}", aggregate.String())
	return nil
}

func (a *aggregateStore) IsAggregateExisting(ctx context.Context, streamID string) error {
	span := makeTracer(ctx, "aggregateStore.Exists", "AggregateID", streamID)

	readStreamOptions := esdb.ReadStreamOptions{Direction: esdb.Backwards, From: esdb.Revision(1)}

	stream, err := a.db.ReadStream(ctx, streamID, readStreamOptions, 1)
	if err != nil {
		return errors.Wrap(err, "db.ReadStream")
	}
	defer stream.Close()

	for {
		_, err := stream.Recv()
		if errors.Is(err, esdb.ErrStreamNotFound) {
			span.RecordError(err)
			return errors.Wrap(esdb.ErrStreamNotFound, "stream.Recv")
		}
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			span.RecordError(err)
			return errors.Wrap(err, "stream.Recv")
		}
	}

	return nil
}

func (a *aggregateStore) SaveAggregate(ctx context.Context, aggregate aggregate.Aggregate) error {
	span := makeTracer(ctx, "aggregateStore.Save", "AggregateID", aggregate.GetID())

	if len(aggregate.GetUncommittedEvents()) == 0 {
		log.Printf("(Save) [no uncommittedEvents] len: {%d}", len(aggregate.GetUncommittedEvents()))
		return nil
	}

	eventsData := make([]esdb.EventData, 0, len(aggregate.GetUncommittedEvents()))
	for _, ev := range aggregate.GetUncommittedEvents() {
		eventsData = append(eventsData, ev.NewEventStoreData())
	}

	var expectedRevision esdb.ExpectedRevision
	if aggregate.GetVersion() == config.AGGREGATE_START_VERSION {
		expectedRevision = esdb.NoStream{}
		log.Printf("(Save) expectedRevision: {%T}", expectedRevision)

		appendStream, err := a.db.AppendToStream(
			ctx,
			aggregate.GetID(),
			esdb.AppendToStreamOptions{ExpectedRevision: expectedRevision},
			eventsData...,
		)
		if err != nil {
			span.RecordError(err)

			return errors.Wrap(err, "db.AppendToStream")
		}

		log.Printf("(Save) stream: {%+v}", appendStream)
		return nil
	}

	readOps := esdb.ReadStreamOptions{Direction: esdb.Backwards, From: esdb.End{}}
	stream, err := a.db.ReadStream(context.Background(), aggregate.GetID(), readOps, 1)
	if err != nil {
		span.RecordError(err)

		return errors.Wrap(err, "db.ReadStream")
	}
	defer stream.Close()

	lastEvent, err := stream.Recv()
	if err != nil {
		span.RecordError(err)

		return errors.Wrap(err, "stream.Recv")
	}

	expectedRevision = esdb.Revision(lastEvent.OriginalEvent().EventNumber)
	log.Printf("(Save) expectedRevision: {%T}", expectedRevision)

	appendStream, err := a.db.AppendToStream(
		ctx,
		aggregate.GetID(),
		esdb.AppendToStreamOptions{ExpectedRevision: expectedRevision},
		eventsData...,
	)
	if err != nil {
		span.RecordError(err)

		return errors.Wrap(err, "db.AppendToStream")
	}

	log.Printf("(Save) stream: {%+v}", appendStream)
	aggregate.ClearChanges()
	return nil
}
