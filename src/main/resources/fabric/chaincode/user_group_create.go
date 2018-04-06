package main

import (
	"encoding/json"
	"fmt"

	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type createUserGroupParams struct {
	Name      string
	AssetKeys []string
	Id        int64
}

// peer chaincode invoke -n mycc -c '{"Args": ["createUserGroup", "{\"Id\":1, \"Name\":\"UserGroup 1\", \"AssetKeys\":[\"Asset1\"]}"]}' -C myc
func CreateUserGroup(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := createUserGroupParams{}

	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	group := models.NewUserGroup()

	for _, element := range params.AssetKeys {
		assetAsBytes, _ := APIstub.GetState(element)

		if assetAsBytes == nil {
			return shim.Error(fmt.Sprintf("Asset with key %s not found", element))
		}

		group.Assets = append(group.Assets, element)
	}

	group.Name = params.Name
	group.Id = params.Id

	// check if User group key already exist
	uGroupAsBytes, _ := APIstub.GetState(group.Key())
	if uGroupAsBytes != nil {
		return shim.Error(fmt.Sprintf("Group with key %s already exist", group.Key()))
	}

	// Put new endpoint state to chain
	groupAsBytes, _ := json.Marshal(group)
	if err := APIstub.PutState(group.Key(), groupAsBytes); err != nil {
		return shim.Error(err.Error())
	}

	// CREATE COMPOSITE KEY FOR USER GROUP
	indexName := "UserGroup"
	groupIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{group.Key(), fmt.Sprintf("%d", group.Id), group.Name})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	APIstub.PutState(groupIndexKey, value)

	return shim.Success(nil)
}
