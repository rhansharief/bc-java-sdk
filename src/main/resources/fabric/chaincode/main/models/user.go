package models

import (
	"fmt"
)

type User struct {
	Id         int64
	Email      string
	Username   string
	Address    string
	FirstName  string
	MiddleName string
	LastName   string
	Mobile     string
	Assets     []UserAsset
}

func (u *User) Key() string {
	return fmt.Sprintf("User%d", u.Id)
}
