package main

import (
	"encoding/json"
	"fmt"

	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// peer chaincode invoke -n mycc -c '{"Args": ["createAsset", "{\"Id\":1, \"Name\":\"Asset1\"}"]}' -C myc
type createAssetParams struct {
	Id   int64
	Name string
}

func CreateAsset(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := createAssetParams{}

	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	asset := models.Asset{}

	asset.Name = params.Name
	asset.Id = params.Id

	getAssetAsBytes, _ := APIstub.GetState(asset.Key())
	if getAssetAsBytes != nil {
		return shim.Error(fmt.Sprintf("Asset with key %s already exist", asset.Key()))
	}

	assetAsBytes, _ := json.Marshal(asset)

	if err := APIstub.PutState(asset.Key(), assetAsBytes); err != nil {
		return shim.Error(err.Error())
	}

	// CREATE COMPOSITE KEY FOR USER GROUP
	indexName := "Asset"
	assetIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{asset.Key(), fmt.Sprintf("%d", asset.Id), asset.Name})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	APIstub.PutState(assetIndexKey, value)

	return shim.Success(nil)
}
