package projection

import (
	"context"
	"eventstore-intro/pkg/eventstore/event"
)

type Projection interface {
	When(ctx context.Context, event event.EventModel) error
}
