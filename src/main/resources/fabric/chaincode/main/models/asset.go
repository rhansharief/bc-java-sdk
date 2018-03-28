package models

import "fmt"

type Asset struct {
	Id   int64
	Name string
}

func (u *Asset) Key() string {
	return fmt.Sprintf("Asset%d", u.Id)
}
