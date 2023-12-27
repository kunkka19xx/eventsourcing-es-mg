package aggregate

import (
	"context"
	"eventstore-intro/pkg/eventstore/aggregate"
	"eventstore-intro/pkg/eventstore/event"
	"eventstore-intro/pkg/eventstore/z-external-app/account/events"
	"eventstore-intro/pkg/eventstore/z-external-app/account/models"
	"github.com/pkg/errors"
)

type AccountAggregate struct {
	*aggregate.AggregateBase
	Account *models.Account
}

func NewAccountAggregateWithID(id string) *AccountAggregate {
	accountAggregate := NewAccountAggregate()
	accountAggregate.SetID(id)
	accountAggregate.Account.CardNo = id
	return accountAggregate
}

func NewAccountAggregate() *AccountAggregate {
	accountAggregate := &AccountAggregate{Account: models.NewAccount()}
	aggregateBase := aggregate.NewAggregateBase(accountAggregate.When)
	aggregateBase.SetType("account")
	accountAggregate.AggregateBase = aggregateBase
	return accountAggregate
}

func (a *AccountAggregate) When(evt event.EventModel) error {
	switch evt.GetEventType() {
	case events.AccountActivated:
		return a.onAccountActivated(evt)

	case events.AccountUpdated:
		return a.onAccountUpdated(evt)

	case events.AccountDeactivated:
		return a.onAccountDeactivated(evt)
	default:
		return errors.New("Invalid event")
	}
}

func (a *AccountAggregate) onAccountActivated(evt event.EventModel) error {
	var eventData events.AccountActivatedEvent
	err := evt.GetJsonData(&eventData)
	if err != nil {
		return errors.New("Error while getting Event's Json Data")
	}
	a.Account.AccountNumber = eventData.AccountNumber
	a.Account.CardNo = eventData.CardNo
	a.Account.StatusControl = eventData.StatusControl
	a.Account.AvailableBalance = eventData.AvailableBalance
	return nil
}

func (a *AccountAggregate) onAccountUpdated(evt event.EventModel) error {
	return nil

}

func (a *AccountAggregate) onAccountDeactivated(evt event.EventModel) error {
	return nil

}

func (a *AccountAggregate) ActivateAccount(ctx context.Context, cardNo, accountNumber string) error {
	if cardNo == "" {
		return errors.New("Card No is required")
	}

	evt, err := events.NewAccountActivatedEvent(a, cardNo, accountNumber)
	if err != nil {
		return errors.New("Error while creating account activated event")
	}
	err = evt.SetMetadata(context.Background())
	if err != nil {
		return err
	}
	return a.Apply(evt)
}
