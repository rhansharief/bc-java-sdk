package main

import (
	"encoding/json"
	"fmt"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type deleteAssetParams struct {
	Key string
}

func DeleteAsset(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := deleteAssetParams{}

	// parse params string to struct and throw an error if found
	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	// check if UsetAsset key doesn't exist
	assetAsBytes, _ := APIstub.GetState(params.Key)
	asset := models.Asset{}
	if assetAsBytes == nil {
		return shim.Error(fmt.Sprintf("Asset with key %s not found", params.Key))
	}
	if err := json.Unmarshal(assetAsBytes, &asset); err != nil {
		return shim.Error(err.Error())
	}

	if err := APIstub.DelState(params.Key); err != nil {
		return shim.Error(err.Error())
	}

	// maintain the index
	indexName := "Asset"
	assetIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{asset.Key(), fmt.Sprintf("%d", asset.Id), asset.Name})
	if err != nil {
		return shim.Error(err.Error())
	}

	//  Delete index entry to state.
	err = APIstub.DelState(assetIndexKey)
	if err != nil {
		return shim.Error("Failed to delete composite key:" + err.Error())
	}

	return shim.Success(nil)
}
