package infra

import (
	"context"
	"eventstore-intro/pkg/eventstore/aggregate"
	"eventstore-intro/pkg/eventstore/config"
	"eventstore-intro/pkg/eventstore/event"
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

func MakeTracer(ctx context.Context, spanName, key, value string) trace.Span {
	ctx, span := tracer.Start(ctx, spanName)
	defer span.End()

	kv := attribute.String(key, value)
	span.SetAttributes(kv)
	return span
}

func (a *aggregateStore) LoadAggregate(ctx context.Context, aggregate aggregate.Aggregate) error {
	span := MakeTracer(ctx, "aggregateStore.Load", "AggregateID", aggregate.GetID())
	defer span.End()

	stream, err := a.db.ReadStream(ctx, aggregate.GetID(), esdb.ReadStreamOptions{}, math.MaxInt64)
	if err != nil {
		return logAndWrapError(span, "db.ReadStream", err)
	}
	defer stream.Close()

	for {
		ev, err := stream.Recv()
		if err != nil {
			if errors.Is(err, esdb.ErrStreamNotFound) || errors.Is(err, os.ErrNotExist) {
				// Break the loop if the stream is not found or does not exist
				break
			}
			return logAndWrapError(span, "stream.Recv", err)
		}

		esEvent := event.NewEventFromRecorded(ev.Event)
		if err := aggregate.RaiseEvent(esEvent); err != nil {
			return logAndWrapError(span, "RaiseEvent", err)
		}
		log.Printf("(Load) esEvent: {%s}", esEvent)
	}

	return nil
}

func (a *aggregateStore) IsAggregateExisting(ctx context.Context, streamID string) (bool, error) {
	span := MakeTracer(ctx, "aggregateStore.Exists", "AggregateID", streamID)
	defer span.End()

	readStreamOptions := esdb.ReadStreamOptions{Direction: esdb.Backwards, From: esdb.Revision(1)}

	stream, err := a.db.ReadStream(ctx, streamID, readStreamOptions, 1)
	if err != nil {
		if errors.Is(err, esdb.ErrStreamNotFound) {
			span.RecordError(err)
			return false, nil // Aggregate does not exist
		}
		span.RecordError(err)
		return false, errors.Wrap(err, "db.ReadStream")
	}
	defer stream.Close()

	_, err = stream.Recv()
	if err != nil {
		if errors.Is(err, io.EOF) {
			// Reached the end of the stream without finding an event
			return false, nil
		}
		span.RecordError(err)
		return false, errors.Wrap(err, "stream.Recv")
	}

	// Aggregate exists
	return true, nil
}

func (a *aggregateStore) SaveAggregate(ctx context.Context, aggregate aggregate.Aggregate) error {
	span := MakeTracer(ctx, "aggregateStore.Save", "AggregateID", aggregate.GetID())
	defer span.End()

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
	} else {
		exists, err := a.IsAggregateExisting(ctx, aggregate.GetID())
		if err != nil {
			return logAndWrapError(span, "checking aggregate existence", err)
		}
		if !exists {
			expectedRevision = esdb.NoStream{}
		} else {
			// TODO: Handle the case where the stream already exists
		}
	}

	_, err := a.db.AppendToStream(
		ctx,
		aggregate.GetID(),
		esdb.AppendToStreamOptions{ExpectedRevision: expectedRevision},
		eventsData...,
	)
	if err != nil {
		return logAndWrapError(span, "db.AppendToStream", err)
	}

	aggregate.ClearUncommittedEvents()
	return nil
}

func logAndWrapError(span trace.Span, action string, err error) error {
	log.Printf("(%s) error: %v", action, err)
	span.RecordError(err)
	return errors.Wrap(err, action)
}
