package events

import (
	"eventstore-intro/pkg/eventstore/aggregate"
	"eventstore-intro/pkg/eventstore/event"
)

const (
	AccountActivated   = "ACCOUNT_ACTIVATED"
	AccountUpdated     = "ACCOUNT_UPDATED"
	AccountDeactivated = "ACCOUNT_DEACTIVATED"
)

type AccountActivatedEvent struct {
	AccountNumber    string
	CardNo           string
	StatusControl    int
	AvailableBalance int64
}

func NewAccountActivatedEvent(aggregate aggregate.Aggregate, cardNo, accountNumber string) (event.EventModel, error) {
	eventData := AccountActivatedEvent{
		CardNo:           cardNo,
		AccountNumber:    accountNumber,
		StatusControl:    1,
		AvailableBalance: 0,
	}
	evt := event.NewBaseEvent(aggregate, AccountActivated)
	if err := evt.SetJsonData(&eventData); err != nil {
		return event.EventModel{}, err
	}
	return evt, nil
}
