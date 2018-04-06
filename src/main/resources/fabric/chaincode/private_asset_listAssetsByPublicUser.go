package main

import (
	"fmt"
	"encoding/json"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Get assets of a User
// Params:
//    userKey="User1"
func PrivateListAssetsByPublicUser(APIstub shim.ChaincodeStubInterface, userKey string) ([]models.Asset, string) {
	publicUserAssets := []models.Asset{}

	// check if User exists and get User model
	userAsBytes, _ := APIstub.GetState(userKey)
	user := models.User{}
	if userAsBytes == nil {
		return []models.Asset{}, fmt.Sprintf("User with key %s not found", userKey)
	}
	if err := json.Unmarshal(userAsBytes, &user); err != nil {
		return []models.Asset{}, err.Error()
	}

	// loop through all User.UserGroups
	for _, userGroupEl := range user.UserGroups {
		// check if UserGroup exist and get UserGroup model
		userGroupAsBytes, _ := APIstub.GetState(userGroupEl)
		userGroup := models.UserGroup{}
		if userGroupAsBytes == nil {
			return []models.Asset{}, fmt.Sprintf("UserGroup with key %s not found", userGroupEl)
		}
		if err := json.Unmarshal(userGroupAsBytes, &userGroup); err != nil {
			return []models.Asset{}, err.Error()
		}

		// loop through all UserGroup.Assets
		for _, assetEl := range userGroup.Assets {
			// check if Asset exist and get Asset model
			assetAsBytes, _ := APIstub.GetState(assetEl)
			asset := models.Asset{}
			if assetAsBytes == nil {
				return []models.Asset{}, fmt.Sprintf("Asser with key %s not found", assetEl)
			}
			if err := json.Unmarshal(assetAsBytes, &asset); err != nil {
				return []models.Asset{}, err.Error()
			}

			// push asset to Assets array
			publicUserAssets = append(publicUserAssets, asset)
		}
	}

	return publicUserAssets, ""
}
