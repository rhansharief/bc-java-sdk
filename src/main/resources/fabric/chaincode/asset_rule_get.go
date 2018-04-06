package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// {"Args": ["getAssetRule", "{\"Key\":\"AssetRule1\"}"]}
type getAssetRuleParams struct {
	Key string
}

func GetAssetRule(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := getAssetRuleParams{}

	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	assetAsBytes, _ := APIstub.GetState(params.Key)

	if assetAsBytes == nil {
		return shim.Error(fmt.Sprintf("Asset Rule with key %s not found", params.Key))
	}

	return shim.Success(assetAsBytes)
}
