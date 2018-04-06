package models

import (
	"fmt"
	"time"
)

type AssetRule struct {
	Id					int64
	DocType     string
	DateCreated string
	LastUpdated string
	Name				string
	Source			string
	Destination	string
	Factor			float32
	Fee					float32
}

func (aR *AssetRule) Key() string {
	return fmt.Sprintf("AssetRule%d", aR.Id)
}

func NewAssetRule() AssetRule {
	t := time.Now()

	assetRule := AssetRule{}
	assetRule.DocType = "ASSET_RULE"
	assetRule.DateCreated = fmt.Sprintf("%s", t.Local())

	return assetRule
}
