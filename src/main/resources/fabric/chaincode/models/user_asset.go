package models

import (
	"fmt"
	"strings"
	"strconv"
	"time"
)

type UserAsset struct {
	Id          int64
	DocType     string
	DateCreated string
	LastUpdated string
	User        string
	Asset       string
	Balance     float32
}

func (uA *UserAsset) Key() string {
	user, _:= strconv.Atoi(strings.Split(uA.User, "User")[1])
	userId := int64(user)

	asset, _:= strconv.Atoi(strings.Split(uA.Asset, "Asset")[1])
	assetId := int64(asset)

	return fmt.Sprintf("UserAsset%d%d", userId, assetId)
}

func NewUserAsset() UserAsset {
	t := time.Now()

	userAsset := UserAsset{}
	userAsset.DocType = "USER_ASSET"
	userAsset.DateCreated = fmt.Sprintf("%s", t.Local())

	return userAsset
}
