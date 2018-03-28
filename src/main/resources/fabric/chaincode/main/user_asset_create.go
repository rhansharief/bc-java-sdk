package main

import (
	"encoding/json"
	"fmt"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// {"Args": ["createUserAsset", "{\"Key\":\"UserAsset1\", \"UserId\":\"User1\", \"AssetId\":\"Asset1\"}"]}
type createUserAssetParams struct {
	Key     string
	UserId  string
	AssetId string
}

// TODO check if working after create User and Asset is done
func CreateUserAsset(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := createUserAssetParams{}

	// parse params string to struct and throw an error if found
	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	// check if UserAsset key already exist
	userAssetAsBytes, _ := APIstub.GetState(params.Key)
	if userAssetAsBytes != nil {
		return shim.Error(fmt.Sprintf("UserAsset with key %s already exist", params.Key))
	}

	// get User
	userAsBytes, _ := APIstub.GetState(params.UserId)
	user := models.User{}

	if err := json.Unmarshal(userAsBytes, &user); err != nil {
		return shim.Error(err.Error())
	}

	if userAsBytes == nil {
		return shim.Error(fmt.Sprintf("User with key %s not found", params.UserId))
	}

	// get Asset
	assetAsBytes, _ := APIstub.GetState(params.AssetId)
	asset := models.Asset{}

	if err := json.Unmarshal(assetAsBytes, &asset); err != nil {
		return shim.Error(err.Error())
	}

	if assetAsBytes == nil {
		return shim.Error(fmt.Sprintf("Asset with key %s not found", params.AssetId))
	}

	userAsset := models.UserAsset{}
	userAsset.User = user
	userAsset.Asset = asset

	// Put new UserAsset to chain
	userAssetAsBytes, _ = json.Marshal(userAsset)
	if err := APIstub.PutState(userAsset.Key(), userAssetAsBytes); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
