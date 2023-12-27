package models

type AccountProjection struct {
	AccountNumber    string `json:"accountNumber" bson:"_id,omitempty"`
	AccountingTitle  string `json:"accountingTitle" bson:"accountingTitle,omitempty"`
	ClientID         string `json:"clientID" bson:"clientID,omitempty"`
	ClientName       string `json:"clientName" bson:"clientName,omitempty"`
	CardNo           string `json:"cardNo" bson:"cardNo,omitempty"`
	StatusControl    int    `json:"statusControl" bson:"statusControl,omitempty"`
	AvailableBalance int64  `json:"availableBalance" bson:"availableBalance,omitempty"`
}
