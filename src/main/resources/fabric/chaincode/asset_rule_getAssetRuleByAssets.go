package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// {"Args": ["getAssetRuleByAssets", "{\"Source\":\"Asset1\", \"Destination\":\"Asset2\"}"]}
type GetAssetRuleByAssetsParams struct {
	Source       string
	Destination  string
}

func GetAssetRuleByAssets(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := GetAssetRuleByAssetsParams{}

	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	assetRuleAsBytes, err := PrivateGetAssetRuleByAssets(APIstub, params.Source, params.Destination)
	if err != "" {
		return shim.Error(err)
	}
	if assetRuleAsBytes == nil {
		return shim.Error(fmt.Sprintf("Asset Rule with source key=%s and destination key=%s doesn't exist", params.Source, params.Destination))
	}

	return shim.Success(assetRuleAsBytes)
}






