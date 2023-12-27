package handler

import (
	"context"
	"eventstore-intro/pkg/eventstore/event"
	"eventstore-intro/pkg/eventstore/z-external-app/account/events"
	"eventstore-intro/pkg/eventstore/z-external-app/account/models"
	"github.com/pkg/errors"
)

func (o *mongoProjection) OnAccountActivate(ctx context.Context, evt event.EventModel) error {

	var eventData events.AccountActivatedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.New("evt.GetJsonData")
	}

	op := &models.AccountProjection{
		AccountNumber:    eventData.AccountNumber,
		AvailableBalance: 0,
		StatusControl:    1,
		CardNo:           eventData.CardNo,
	}

	_, err := o.accountRepo.ActivateAccount(ctx, op)
	if err != nil {
		return err
	}

	return nil
}
