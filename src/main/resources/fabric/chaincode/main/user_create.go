package main

import (
	"main/models"

	"encoding/json"

	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type createUserParams struct {
	Id         int64
	Email      string
	Username   string
	Address    string
	FirstName  string
	MiddleName string
	LastName   string
	Mobile     string
}

// Create User
func CreateUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := createUserParams{}

	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	var user = models.User{}

	user.Id = params.Id
	user.Email = params.Email
	user.Username = params.Username
	user.Address = params.Address
	user.FirstName = params.FirstName
	user.MiddleName = params.MiddleName
	user.LastName = params.LastName
	user.Mobile = params.Mobile

	// check if User key already exist
	userAsBytes, _ := APIstub.GetState(user.Key())
	if userAsBytes != nil {
		return shim.Error(fmt.Sprintf("User with key %s already exist", user.Key()))
	}

	userAsBytes, _ = json.Marshal(user)

	if err := APIstub.PutState(user.Key(), userAsBytes); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
