package handler

import (
	"context"
	ev "eventstore-intro/pkg/eventstore/event"
	infra "eventstore-intro/pkg/eventstore/infra/repository"
	"eventstore-intro/pkg/eventstore/z-external-app/account/events"
	"eventstore-intro/pkg/eventstore/z-external-app/account/repository"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/pkg/errors"
	_ "golang.org/x/sync/errgroup"
	"log"
)

type mongoProjection struct {
	db          *esdb.Client
	accountRepo repository.AccountMongoRepository
	// add more repository if need
}

func (o *mongoProjection) runWorker(ctx context.Context, worker Worker, stream *esdb.PersistentSubscription, workerID int) func() error {
	return func() error {
		if err := worker(ctx, stream, workerID); err != nil {
			log.Printf("Worker %d encountered an error: %v", workerID, err)
			return err
		}
		return nil
	}
}

type Worker func(ctx context.Context, stream *esdb.PersistentSubscription, workerID int) error

func (o *mongoProjection) ProcessEvents(ctx context.Context, stream *esdb.PersistentSubscription, workerID int) error {

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			event := stream.Recv()

			if event.SubscriptionDropped != nil {
				return errors.Wrap(event.SubscriptionDropped.Error, "Subscription Dropped")
			}

			if event.EventAppeared != nil {

				err := o.When(ctx, ev.NewEventFromRecorded(event.EventAppeared.Event))
				if err != nil {

					if err := stream.Nack(err.Error(), esdb.Nack_Retry, event.EventAppeared); err != nil {
						return errors.Wrap(err, "stream.Nack")
					}
				}

				err = stream.Ack(event.EventAppeared)
				if err != nil {
					return errors.Wrap(err, "stream.Ack")
				}
			}
		}
	}
}

func (o *mongoProjection) When(ctx context.Context, event ev.EventModel) error {
	span := infra.MakeTracer(ctx, "mongoProjection.When", "EventID", event.ID)

	defer span.End()

	switch event.GetEventType() {
	case events.AccountActivated:
		return o.OnAccountActivate(ctx, event)
	default:
		return errors.New("Unknown Event Type")
	}
}
