package models

import (
	"fmt"
	"time"
)

type RetirementRule struct {
	Id               int64
	DocType          string
	DateCreated      string
	LastUpdated      string
	Name             string
	SourceAsset      string
	ExchangeEndpoint string
	Factor           float32
	Fee              float32
}

func (aR *RetirementRule) Key() string {
	return fmt.Sprintf("RetirementRule%d", aR.Id)
}

func NewRetirementRule() RetirementRule {
	t := time.Now()

	retirementRule := RetirementRule{}
	retirementRule.DocType = "RETIREMENT_RULE"
	retirementRule.DateCreated = fmt.Sprintf("%s", t.Local())

	return retirementRule
}
