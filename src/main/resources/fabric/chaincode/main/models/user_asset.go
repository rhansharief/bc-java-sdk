package models

import (
	"fmt"
)

type UserAsset struct {
	Id    int64
	User  User
	Asset Asset
}

func (uA *UserAsset) Key() string {
	return fmt.Sprintf("UserAsset%d", uA.Id)
}
