package main

import (
	"encoding/json"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type updateAssetParams struct {
	Key  string
	Name string
}

func UpdateAsset(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := updateAssetParams{}

	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	assetAsBytes, _ := APIstub.GetState(params.Key)
	asset := models.Asset{}

	if err := json.Unmarshal(assetAsBytes, &asset); err != nil {
		return shim.Error(err.Error())
	}

	asset.Name = params.Name

	assetAsBytes, _ = json.Marshal(asset)
	if err := APIstub.PutState(asset.Key(), assetAsBytes); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
