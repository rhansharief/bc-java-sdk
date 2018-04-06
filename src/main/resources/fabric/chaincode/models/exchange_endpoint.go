package models

import (
	"fmt"
	"time"
)

type ExchangeEndpoint struct {
	Id          int64
	DocType     string
	DateCreated string
	LastUpdated string
	Name        string
}

func (uA *ExchangeEndpoint) Key() string {
	return fmt.Sprintf("ExchangeEndpoint%d", uA.Id)
}

func NewExchangeEndpoint() ExchangeEndpoint {
	t := time.Now()

	ExchangeEndpoint := ExchangeEndpoint{}
	ExchangeEndpoint.DocType = "EXCHANGE_ENDPOINT"
	ExchangeEndpoint.DateCreated = fmt.Sprintf("%s", t.Local())

	return ExchangeEndpoint
}
