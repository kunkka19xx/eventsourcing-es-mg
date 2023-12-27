package aggregate

import (
	"context"
	"eventstore-intro/pkg/eventstore/infra"
)

func LoadAccountAggregate(ctx context.Context, eventStore infra.AggregateStore, aggregateID string) (*AccountAggregate, error) {

	return nil, nil
}
