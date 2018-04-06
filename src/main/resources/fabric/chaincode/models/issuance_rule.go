package models

import (
	"fmt"
	"time"
)

type IssuanceRule struct {
	Id               int64
	DocType          string
	DateCreated      string
	LastUpdated      string
	ExchangeEndpoint string
	Name             string
	DestinationAsset string
	Factor           float32
	Fee              float32
}

func (aR *IssuanceRule) Key() string {
	return fmt.Sprintf("IssuanceRule%d", aR.Id)
}

func NewIssuanceRule() IssuanceRule {
	t := time.Now()

	issuanceRule := IssuanceRule{}
	issuanceRule.DocType = "ISSUANCE_RULE"
	issuanceRule.DateCreated = fmt.Sprintf("%s", t.Local())

	return issuanceRule
}
