package models

import (
	"fmt"
	"time"
)

type Asset struct {
	Id          int64
	DocType     string
	DateCreated string
	LastUpdated string
	Name        string
}

func (u *Asset) Key() string {
	return fmt.Sprintf("Asset%d", u.Id)
}

func NewAsset() Asset {
	t := time.Now()

	asset := Asset{}
	asset.DocType = "ASSET"
	asset.DateCreated = fmt.Sprintf("%s", t.Local())

	return asset
}
