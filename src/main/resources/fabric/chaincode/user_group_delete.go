package main

import (
	"encoding/json"
	"fmt"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// {"Args": ["deleteUserGroup", "{\"Key\":\"UserGroup1\"}"]}
type deleteUserGroupParams struct {
	Key string
}

func DeleteUserGroup(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := deleteUserGroupParams{}

	// parse params string to struct and throw an error if found
	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	// check if UserGroup key doesn't exist
	userGroupAsBytes, _ := APIstub.GetState(params.Key)
	userGroup := models.UserGroup{}
	if userGroupAsBytes == nil {
		return shim.Error(fmt.Sprintf("UserGroup with key %s not found", params.Key))
	}
	if err := json.Unmarshal(userGroupAsBytes, &userGroup); err != nil {
		return shim.Error(err.Error())
	}

	// couch query string - get all objects that has DocType="USER" AND UserGroups contains params.Key
	query := fmt.Sprintf(""+
		"{\"selector\": {"+
		"\"$and\": [{"+
		"\"DocType\": \"USER\","+
		"\"UserGroups\": {"+
		"\"$in\": [\"%s\"]"+
		"}"+
		"}]"+
		"}}", params.Key)

	resultsIterator, err := APIstub.GetQueryResult(query)

	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// loop through all users attach to a group
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		user := models.User{}
		if err := json.Unmarshal(queryResponse.Value, &user); err != nil {
			return shim.Error(err.Error())
		}

		// get index of the element we want to remove
		i := indexOf(params.Key, user.UserGroups)

		// remove element from array by index
		user.UserGroups = removeKey(user.UserGroups, i)

		// save updated user to state
		userAsBytes, _ := json.Marshal(user)
		if err := APIstub.PutState(user.Key(), userAsBytes); err != nil {
			return shim.Error(err.Error())
		}
	}

	// delete UserGroup
	if err := APIstub.DelState(params.Key); err != nil {
		return shim.Error(err.Error())
	}

	// maintain the index
	indexName := "UserGroup"
	userGroupIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{userGroup.Key(), fmt.Sprintf("%d", userGroup.Id), userGroup.Name})
	if err != nil {
		return shim.Error(err.Error())
	}

	// Delete index entry to state.
	err = APIstub.DelState(userGroupIndexKey)
	if err != nil {
		return shim.Error("Failed to delete composite key:" + err.Error())
	}

	return shim.Success(nil)
}
