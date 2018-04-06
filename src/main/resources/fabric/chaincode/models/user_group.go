package models

import (
	"fmt"
	"time"
)

type UserGroup struct {
	Id          int64
	DocType     string
	DateCreated string
	LastUpdated string
	Name        string
	Assets      []string
}

func (uA *UserGroup) Key() string {
	return fmt.Sprintf("UserGroup%d", uA.Id)
}

func NewUserGroup() UserGroup {
	t := time.Now()

	userGroup := UserGroup{}
	userGroup.DocType = "USER_GROUP"
	userGroup.DateCreated = fmt.Sprintf("%s", t.Local())

	return userGroup
}