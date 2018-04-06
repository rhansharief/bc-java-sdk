package main

import (
	"encoding/json"
	"fmt"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// {"Args": ["updateUserAsset", "{\"Key\":\"UserAsset1\", \"UserId\":\"User1\", \"AssetId\":\"Asset1\"}"]}
type updateUserAssetParams struct {
	Key     string
	UserId  string
	AssetId string
}

func UpdateUserAsset(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := updateUserAssetParams{}

	// parse params string to struct and throw an error if found
	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	// get UserAsset
	userAssetAsBytes, _ := APIstub.GetState(params.Key)
	userAsset := models.UserAsset{}

	if userAssetAsBytes == nil {
		return shim.Error(fmt.Sprintf("UserAsset with key %s not found", params.Key))
	}

	if err := json.Unmarshal(userAssetAsBytes, &userAsset); err != nil {
		return shim.Error(err.Error())
	}

	// get User
	userAsBytes, _ := APIstub.GetState(params.UserId)
	user := models.User{}

	if userAsBytes == nil {
		return shim.Error(fmt.Sprintf("User with key %s not found", params.UserId))
	}

	if err := json.Unmarshal(userAsBytes, &user); err != nil {
		return shim.Error(err.Error())
	}

	// get Asset
	assetAsBytes, _ := APIstub.GetState(params.AssetId)
	asset := models.Asset{}

	if assetAsBytes == nil {
		return shim.Error(fmt.Sprintf("Asset with key %s not found", params.AssetId))
	}

	if err := json.Unmarshal(assetAsBytes, &asset); err != nil {
		return shim.Error(err.Error())
	}

	userAsset.User = user.Key()
	userAsset.Asset = asset.Key()

	userAssetAsBytes, _ = json.Marshal(userAsset)
	if err := APIstub.PutState(userAsset.Key(), userAssetAsBytes); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
