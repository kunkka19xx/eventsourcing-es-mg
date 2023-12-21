package domain

import "time"

type BankAccount struct {
	Id            string
	AccountHolder string
	CreationDate  time.Time
	AccountType   AccountType
	Balance       float64
}
type AccountType int

const (
	UNDEFINED AccountType = iota + 1
	SAVINGS
	CURRENT
)

func (a AccountType) String() string {
	names := []string{"UNDEFINED", "SAVINGS", "CURRENT"}
	if int(a) < 1 || int(a) > len(names) {
		return "UNKNOWN"
	}
	return names[a-1]
}

func (a AccountType) EnumIndex() int {
	return int(a)
}
