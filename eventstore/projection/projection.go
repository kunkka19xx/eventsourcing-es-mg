package projection

import (
	"context"
	"eventstore-intro/eventstore/event"
)

type Projection interface {
	When(ctx context.Context, event event.EventModel) error
}
