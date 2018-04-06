package main

import (
	"encoding/json"
	"fmt"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type updateUserGroupParams struct {
	Key       string
	Name      string
	AssetKeys []string
}

// {"Args": ["updateUserGroup", "{\"Key\":\"UserGroup1\", \"Name\":\"Some other UserGroup 1\", \"AssetKeys\":[\"Asset2\"]}"]}
func UpdateUserGroup(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := updateUserGroupParams{}
	assets := []string{}

	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	// check if user group doesn't exist
	userGroupAsBytes, _ := APIstub.GetState(params.Key)
	userGroup := models.UserGroup{}
	if userGroupAsBytes == nil {
		return shim.Error(fmt.Sprintf("UserGroup with key %s not found", params.Key))
	}
	if err := json.Unmarshal(userGroupAsBytes, &userGroup); err != nil {
		return shim.Error(err.Error())
	}

	// check if given list of Asset exist
	for _, element := range params.AssetKeys {
		assetAsBytes, _ := APIstub.GetState(element)
		if assetAsBytes == nil {
			return shim.Error(fmt.Sprintf("Asset with key %s not found", element))
		}
		assets = append(assets, element)
	}

	userGroup.Name = params.Name
	userGroup.Assets = assets

	// put new endpoint state to chain
	groupAsBytes, _ := json.Marshal(userGroup)
	if err := APIstub.PutState(userGroup.Key(), groupAsBytes); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
