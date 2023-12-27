package command

import (
	"context"
	"eventstore-intro/pkg/eventstore/command"
	"eventstore-intro/pkg/eventstore/infra"
	"eventstore-intro/pkg/eventstore/z-external-app/account/aggregate"
	"github.com/pkg/errors"
)

type ActivateAccountCommand struct {
	command.BaseCommand
	CardNo        string
	AccountNumber string
}

func NewActivateAccountCommand(aggregateID string, ctx context.Context, cardNo, accountNumber string) *ActivateAccountCommand {
	return &ActivateAccountCommand{BaseCommand: command.NewBaseCommand(aggregateID), CardNo: cardNo, AccountNumber: accountNumber}
}

type activateAccountHandler struct {
	aggregateStore infra.AggregateStore
}

type ActivateAccountCommandHandler interface {
	Handle(ctx context.Context, command *ActivateAccountCommand) error
}

func (a *activateAccountHandler) Handle(ctx context.Context, command *ActivateAccountCommand) error {

	account := aggregate.NewAccountAggregateWithID(command.AggregateID)
	isAggregateExist, err := a.aggregateStore.IsAggregateExisting(ctx, account.ID)
	if err != nil {
		return errors.New("")
	}
	if isAggregateExist {
		return nil
	}
	err = account.ActivateAccount(ctx, command.CardNo, command.AccountNumber)

	if err != nil {
		return err
	}

	return a.aggregateStore.SaveAggregate(ctx, account)
}
