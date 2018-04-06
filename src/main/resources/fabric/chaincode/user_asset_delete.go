package main

import (
	"encoding/json"
	"fmt"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// peer chaincode invoke -n mycc -c '{"Args": ["deleteUserAsset", "{\"Key\":\"UserAsset1\"}"]}' -C myc
type deleteUserAssetParams struct {
	Key string
}

func DeleteUserAsset(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := deleteUserAssetParams{}

	// parse params string to struct and throw an error if found
	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	// check if UsetAsset key doesn't exist
	userAssetAsBytes, _ := APIstub.GetState(params.Key)
	if userAssetAsBytes == nil {
		return shim.Error(fmt.Sprintf("UserAsset with key %s not found", params.Key))
	}

	userAsset := models.UserAsset{}

	if err := json.Unmarshal(userAssetAsBytes, &userAsset); err != nil {
		return shim.Error(err.Error())
	}

	if err := APIstub.DelState(params.Key); err != nil {
		return shim.Error(err.Error())
	}

	// maintain the index
	indexName := "UserAsset"
	userAssetIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{userAsset.Key(), fmt.Sprintf("%d", userAsset.Id), userAsset.User, userAsset.Asset, fmt.Sprintf("%d", userAsset.Balance)})
	if err != nil {
		return shim.Error(err.Error())
	}

	//  Delete index entry to state.
	err = APIstub.DelState(userAssetIndexKey)
	if err != nil {
		return shim.Error("Failed to delete composite key:" + err.Error())
	}

	return shim.Success(nil)
}
