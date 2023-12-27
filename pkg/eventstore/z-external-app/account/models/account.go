package models

type Account struct {
	AccountNumber    string
	AccountingTitle  string
	ClientID         string
	ClientName       string
	CardNo           string
	StatusControl    int
	AvailableBalance int64
}

func NewAccount() *Account {
	return &Account{
		AccountNumber:    "acc-123",
		AccountingTitle:  "acc-title",
		ClientID:         "id",
		ClientName:       "name",
		CardNo:           "",
		StatusControl:    1,
		AvailableBalance: 0,
	}
}
