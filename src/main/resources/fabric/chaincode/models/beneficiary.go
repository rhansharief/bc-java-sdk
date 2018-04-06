package models

import "fmt"

type Beneficiary struct {
	Id         int64
	DocType    string
	FirstName  string
	MiddleName string
	LastName   string
	Birthdate  string
	Mobile     string
	Email      string
	Address    string
	City       string
	ZipCode    string
	Country    string
	User       string
}

func (u *Beneficiary) Key() string {
	return fmt.Sprintf("Beneficiary%d", u.Id)
}

func NewBeneficiary() Beneficiary {
	Beneficiary := Beneficiary{}
	Beneficiary.DocType = "BENEFICIARY"

	return Beneficiary
}
