package charging

import (
	domain2 "eventstore-intro/eventstore/command"
	"eventstore-intro/pkg/query"
)

type OpenAccountCommand struct {
	domain2.BaseCommand
	AccountHolder  string
	AccountType    domain.AccountType
	OpeningBalance float64
}

type CloseAccountCommand struct {
	domain2.BaseCommand
}

type DepositFundsCommand struct {
	domain2.BaseCommand
	Amount float64
}

type WithdrawFundsCommand struct {
	domain2.BaseCommand
	Amount float64
}
