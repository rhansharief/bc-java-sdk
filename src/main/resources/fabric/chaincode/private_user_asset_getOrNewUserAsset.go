package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func PrivateGetOrNewUserAsset(APIstub shim.ChaincodeStubInterface, userKey string, assetKey string) (models.UserAsset, string) {
	userId := GetId(userKey, "User")
	assetId := GetId(assetKey, "Asset")

	// check if User has access to an Asset and return error if not
	userHasAsset := false
	userAssets, err := PrivateListAssetsByPublicUser(APIstub, userKey)
	if err != "" {
		return models.UserAsset{}, err
	}
	for _, assetEl := range userAssets {
		if assetId == assetEl.Id {
			userHasAsset = true
		}
	}
	if userHasAsset == false {
		return models.UserAsset{}, fmt.Sprintf("User with key %s doesn't have an asset with key %s", userKey, assetKey)
	}

	// get or create UserAsset
	userAssetKey := fmt.Sprintf("UserAsset%d%d", userId, assetId)
	userAssetAsBytes, _ := APIstub.GetState(userAssetKey)
	userAsset := models.UserAsset{}

	if userAssetAsBytes == nil {
		id, _ := strconv.Atoi(fmt.Sprintf("%d%d", userId, assetId))

		// create UserAsset
		params := createUserAssetParams{}
		params.Id = int64(id)
		params.UserId = userKey
		params.AssetId = assetKey
		params.Balance = 100

		userAsset, err = PrivateCreateUserAsset(APIstub, params)
		if err != "" {
			return models.UserAsset{}, err
		}
	} else {
		if err := json.Unmarshal(userAssetAsBytes, &userAsset); err != nil {
			return models.UserAsset{}, err.Error()
		}
	}

	return userAsset, ""
}