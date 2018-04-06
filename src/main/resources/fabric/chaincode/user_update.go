package main

import (
	"encoding/json"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
	"fmt"
)

// {"Args": ["updateUser", "{\"Key\":\"User1\", \"FirstName\":\"Juan\", \"MiddleName\":\"dela\", \"LastName\":\"Cruz\", \"Mobile\":\"09171234567\", \"UserGroups\":[\"UserGroup2\"]}"]}
type updateUserParams struct {
	Key        string
	FirstName  string
	MiddleName string
	LastName   string
	Mobile     string
	UserGroups []string
}

func UpdateUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	// Parse updateUserParams
	params := updateUserParams{}
	userGroups := []string{}

	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	// get and check if user doesn't exist
	userAsBytes, _ := APIstub.GetState(params.Key)
	user := models.User{}
	if userAsBytes == nil {
		return shim.Error(fmt.Sprintf("User with key %s not found", params.Key))
	}
	if err := json.Unmarshal(userAsBytes, &user); err != nil {
		return shim.Error(err.Error())
	}

	// check if given list of UserGroup exist
	for _, element := range params.UserGroups {
		userGroupAsBytes, _ := APIstub.GetState(element)
		if userGroupAsBytes == nil {
			return shim.Error(fmt.Sprintf("UserGroup with key %s not found", element))
		}
		userGroups = append(userGroups, element)
	}

	user.FirstName = params.FirstName
	user.MiddleName = params.MiddleName
	user.LastName = params.LastName
	user.Mobile = params.Mobile
	user.UserGroups = userGroups

	// Put new user state to chain
	userAsBytes, _ = json.Marshal(user)
	if err := APIstub.PutState(user.Key(), userAsBytes); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
