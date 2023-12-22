package charging

import (
	"eventstore-intro/pkg/query"
)

type OpenAccountCommand struct {
	AccountHolder  string
	AccountType    domain.AccountType
	OpeningBalance float64
}

type CloseAccountCommand struct {
}

type DepositFundsCommand struct {
	Amount float64
}

type WithdrawFundsCommand struct {
	Amount float64
}
