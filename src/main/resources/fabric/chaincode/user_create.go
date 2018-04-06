package main

import (
	"main/models"

	"encoding/json"

	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// peer chaincode invoke -n mycc -c '{"Args":["createUser", "{\"Id\":1,\"UID\":\"176000120221\",\"Email\":\"test@gmail.com\",\"Username\":\"test\",\"Address\":\"Ortigas\",\"FirstName\":\"Juan\",\"MiddleName\":\"Pedro\",\"LastName\":\"Corazon\",\"Mobile\":\"09992323232\", \"UserGroups\":[\"UserGroup1\"]}"]}' -C myc
type createUserParams struct {
	Id         int64
	UID        string
	Email      string
	Username   string
	Address    string
	FirstName  string
	MiddleName string
	LastName   string
	Mobile     string
	UserGroups []string
}

// Create User
func CreateUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := createUserParams{}

	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	var user = models.NewUser()

	user.Id = params.Id
	user.UID = params.UID
	user.Email = params.Email
	user.Username = params.Username
	user.Address = params.Address
	user.FirstName = params.FirstName
	user.MiddleName = params.MiddleName
	user.LastName = params.LastName
	user.Mobile = params.Mobile

	for _, element := range params.UserGroups {
		groupAsBytes, _ := APIstub.GetState(element)

		if groupAsBytes == nil {
			return shim.Error(fmt.Sprintf("Group with key %s not found", element))
		}

		user.UserGroups = append(user.UserGroups, element)
	}

	// check if User key already exist
	userAsBytes, _ := APIstub.GetState(user.Key())
	if userAsBytes != nil {
		return shim.Error(fmt.Sprintf("User with key %s already exist", user.Key()))
	}

	userAsBytes, _ = json.Marshal(user)

	if err := APIstub.PutState(user.Key(), userAsBytes); err != nil {
		return shim.Error(err.Error())
	}

	// CREATE COMPOSITE KEY FOR USER
	indexName := "User"
	userIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{user.Key(), fmt.Sprintf("%d", user.Id), user.UID, user.Email, user.Username, user.Address, user.FirstName, user.MiddleName, user.LastName, user.Mobile})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	APIstub.PutState(userIndexKey, value)

	return shim.Success(nil)
}
