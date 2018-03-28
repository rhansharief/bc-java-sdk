package main

import (
	"encoding/json"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type createAssetParams struct {
	Name string
	Id   int64
}

func CreateAsset(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := createAssetParams{}

	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	asset := models.Asset{}
	asset.Name = params.Name
	asset.Id = params.Id

	assetAsBytes, _ := json.Marshal(asset)

	if err := APIstub.PutState(asset.Key(), assetAsBytes); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
