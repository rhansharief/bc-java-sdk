package main

import (
	"encoding/json"
	"fmt"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func PrivateCreateUserAsset(APIstub shim.ChaincodeStubInterface, params createUserAssetParams) (models.UserAsset, string) {
	// get User
	userAsBytes, _ := APIstub.GetState(params.UserId)
	user := models.User{}
	if userAsBytes == nil {
		return models.UserAsset{}, fmt.Sprintf("User with key %s not found", params.UserId)
	}
	if err := json.Unmarshal(userAsBytes, &user); err != nil {
		return models.UserAsset{}, err.Error()
	}

	// get Asset
	assetAsBytes, _ := APIstub.GetState(params.AssetId)
	asset := models.Asset{}
	if assetAsBytes == nil {
		return models.UserAsset{}, fmt.Sprintf("Asset with key %s not found", params.AssetId)
	}
	if err := json.Unmarshal(assetAsBytes, &asset); err != nil {
		return models.UserAsset{}, err.Error()
	}

	userAsset := models.NewUserAsset()
	userAsset.Id = params.Id
	userAsset.User = user.Key()
	userAsset.Asset = asset.Key()
	userAsset.Balance = params.Balance

	// check if UserAsset key already exist
	if userAssetAsBytes, _ := APIstub.GetState(userAsset.Key()); userAssetAsBytes != nil {
		return models.UserAsset{}, fmt.Sprintf("UserAsset with key %v already exist", params.Id)
	}

	// Put new UserAsset to chain
	userAssetAsBytes, _ := json.Marshal(userAsset)
	if err := APIstub.PutState(userAsset.Key(), userAssetAsBytes); err != nil {
		return models.UserAsset{}, err.Error()
	}

	// include newly created UserAsset to list of assets
	user.Assets = append(user.Assets, userAsset.Key())
	userAsBytes, _ = json.Marshal(user)
	if err := APIstub.PutState(user.Key(), userAsBytes); err != nil {
		return models.UserAsset{}, err.Error()
	}

	// CREATE COMPOSITE KEY FOR USER
	indexName := "UserAsset"
	userAssetIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{userAsset.Key(), fmt.Sprintf("%d", userAsset.Id), userAsset.User, userAsset.Asset, fmt.Sprintf("%d", userAsset.Balance)})
	if err != nil {
		return models.UserAsset{}, err.Error()
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	APIstub.PutState(userAssetIndexKey, value)

	return userAsset, ""
}
