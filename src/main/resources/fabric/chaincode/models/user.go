package models

import (
	"fmt"
	"time"
)

type User struct {
	Id            int64
	UID           string
	DocType       string
	DateCreated   string
	LastUpdated   string
	Email         string
	Username      string
	Address       string
	FirstName     string
	MiddleName    string
	LastName      string
	Mobile        string
	Assets        []string
	UserGroups    []string
	Beneficiaries []string
}

func (u *User) Key() string {
	return fmt.Sprintf("User%d", u.Id)
}

func NewUser() User {
	t := time.Now()

	user := User{}
	user.DocType = "USER"
	user.DateCreated = fmt.Sprintf("%s", t.Local())

	return user
}
