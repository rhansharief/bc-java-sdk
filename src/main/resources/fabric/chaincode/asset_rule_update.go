package main

import (
	"encoding/json"
	"fmt"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// {"Args": ["updateAssetRule", "{\"Key\":\"AssetRule1\", \"Factor\":50.50, \"Fee\":0.50}"]}
type updateAssetRuleParams struct {
	Key    	string
	Factor  float32
	Fee 		float32
}

func UpdateAssetRule(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := updateAssetRuleParams{}

	// parse params string to struct and throw an error if found
	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	// check if AssetRule exist
	assetRuleAsBytes, _ := APIstub.GetState(params.Key)
	assetRule := models.AssetRule{}
	if assetRuleAsBytes == nil {
		return shim.Error(fmt.Sprintf("AssetRule with key %s not found", params.Key))
	}
	if err := json.Unmarshal(assetRuleAsBytes, &assetRule); err != nil {
		return shim.Error(err.Error())
	}

	// update values
	assetRule.Factor = params.Factor
	assetRule.Fee = params.Fee

	// put new AssetRule to chain
	assetRuleAsBytes, _ = json.Marshal(assetRule)
	if err := APIstub.PutState(assetRule.Key(), assetRuleAsBytes); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
