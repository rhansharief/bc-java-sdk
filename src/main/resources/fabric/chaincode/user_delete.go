package main

import (
	"encoding/json"
	"fmt"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// peer chaincode invoke -n mycc -c '{"Args": ["deleteUser", "{\"Key\":\"User1\"}"]}' -C myc
type deleteUserParams struct {
	Key string
}

func DeleteUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := deleteUserParams{}

	// parse params string to struct and throw an error if found
	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	// check if User key doesn't exist
	userAsBytes, _ := APIstub.GetState(params.Key)
	if userAsBytes == nil {
		return shim.Error(fmt.Sprintf("User with key %s not found", params.Key))
	}

	user := models.User{}

	if err := json.Unmarshal(userAsBytes, &user); err != nil {
		return shim.Error(err.Error())
	}

	if err := APIstub.DelState(params.Key); err != nil {
		return shim.Error(err.Error())
	}

	// maintain the index
	indexName := "User"
	userIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{user.Key(), fmt.Sprintf("%d", user.Id), user.Email, user.Username, user.Address, user.FirstName, user.MiddleName, user.LastName, user.Mobile})
	if err != nil {
		return shim.Error(err.Error())
	}

	//  Delete index entry to state.
	err = APIstub.DelState(userIndexKey)
	if err != nil {
		return shim.Error("Failed to delete composite key:" + err.Error())
	}

	return shim.Success(nil)
}
