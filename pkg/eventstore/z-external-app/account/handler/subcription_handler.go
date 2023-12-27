package handler

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"golang.org/x/sync/errgroup"
)

func (o *mongoProjection) Subscribe(ctx context.Context, prefixes []string, poolSize int, worker Worker) error {

	err := o.db.CreatePersistentSubscriptionAll(ctx, "mongo-prj", esdb.PersistentAllSubscriptionOptions{
		Filter: &esdb.SubscriptionFilter{Type: esdb.StreamFilterType, Prefixes: prefixes},
	})
	if err != nil {
		if subscriptionError, ok := err.(*esdb.PersistentSubscriptionError); !ok || ok && (subscriptionError.Code != 6) {
			return err

		}
	}

	stream, err := o.db.ConnectToPersistentSubscription(
		ctx,
		"$all",
		"mongo-prj",
		esdb.ConnectToPersistentSubscriptionOptions{},
	)
	if err != nil {
		return err
	}
	defer stream.Close()

	g, ctx := errgroup.WithContext(ctx)
	for i := 0; i < poolSize; i++ {
		g.Go(o.runWorker(ctx, worker, stream, i))
	}
	return g.Wait()
}
