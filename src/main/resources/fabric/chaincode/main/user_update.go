package main

import (
	"encoding/json"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// {"Args": ["updateUser", "{\"Key\":\"User1\"},{\"FirstName\":\"Juan\"},{\"MiddleName\":\"dela\"},{\"LastName\":\"Cruz\"},{\"Mobile\":\"09171234567\"}"]}
type updateUserParams struct {
	Key        string
	FirstName  string
	MiddleName string
	LastName   string
	Mobile     string
}

// TODO: Check if this is working when User create and get has been implemented
func UpdateUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	// Parse updateUserParams
	params := updateUserParams{}

	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	// Get user
	userAsBytes, _ := APIstub.GetState(params.Key)
	user := models.User{}

	if err := json.Unmarshal(userAsBytes, &user); err != nil {
		return shim.Error(err.Error())
	}

	user.FirstName = params.FirstName
	user.MiddleName = params.MiddleName
	user.LastName = params.LastName
	user.Mobile = params.Mobile

	// Put new user state to chain
	userAsBytes, _ = json.Marshal(user)
	if err := APIstub.PutState(user.Key(), userAsBytes); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
