package models

import (
	"fmt"
	"time"
)

type IssuanceTransaction struct {
	Id    	          int64
	DocType           string
	DateCreated       string
	LastUpdated       string
	Broker            string
	ExchangeEndpoint  string
	DestinationAsset  string
	Factor            float32
	Fee               float32
	Amount 	          float32
}

func (t *IssuanceTransaction) Key() string {
	return fmt.Sprintf("IssuanceTransaction%d", t.Id)
}

func NewIssuanceTransaction() IssuanceTransaction {
	t := time.Now()

	issuanceTransaction := IssuanceTransaction{}
	issuanceTransaction.DocType = "ISSUANCE_TRANSACTION"
	issuanceTransaction.DateCreated = fmt.Sprintf("%s", t.Local())

	return issuanceTransaction
}
