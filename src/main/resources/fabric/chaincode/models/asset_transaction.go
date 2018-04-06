package models

import (
	"fmt"
	"time"
)

type AssetTransaction struct {
	Id    	          int64
	DocType           string
	DateCreated       string
	LastUpdated       string
	SourceAsset       string
	DestinationAsset  string
	Factor            float32
	Fee               float32
	Amount 	          float32
}

func (t *AssetTransaction) Key() string {
	return fmt.Sprintf("Transaction%d", t.Id)
}

func NewAssetTransaction() AssetTransaction {
	t := time.Now()

	assetTransaction := AssetTransaction{}
	assetTransaction.DocType = "ASSET_TRANSACTION"
	assetTransaction.DateCreated = fmt.Sprintf("%s", t.Local())

	return assetTransaction
}
