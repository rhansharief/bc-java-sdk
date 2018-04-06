package main

import (
	"encoding/json"
	"fmt"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// peer chaincode invoke -n mycc -c '{"Args": ["createRetirementRule", "{\"Id\":1, \"Name\":\"RetirementRule\", \"ExchangeEndpointKey\": \"ExchangeEndpoint1\", \"SourceAssetKey\":\"Asset1\", \"Factor\":50.50, \"Fee\":0.50}"]}' -C myc
type createRetirementRuleParams struct {
	Id                  int64
	Name                string
	ExchangeEndpointKey string
	SourceAssetKey      string
	Factor              float32
	Fee                 float32
}

func CreateRetirementRule(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := createRetirementRuleParams{}

	// parse params string to struct and throw an error if found
	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	// check Source Asset if it exist
	endpointAsBytes, _ := APIstub.GetState(params.ExchangeEndpointKey)
	endpoint := models.ExchangeEndpoint{}

	if endpointAsBytes == nil {
		return shim.Error(fmt.Sprintf("Exchange endpoint with key %s not found", params.ExchangeEndpointKey))
	}

	if err := json.Unmarshal(endpointAsBytes, &endpoint); err != nil {
		return shim.Error(err.Error())
	}

	// check Destination Asset if it exist
	sourceAssetAsBytes, _ := APIstub.GetState(params.SourceAssetKey)
	sourceAsset := models.Asset{}

	if sourceAssetAsBytes == nil {
		return shim.Error(fmt.Sprintf("Asset with key %s not found", params.SourceAssetKey))
	}

	if err := json.Unmarshal(sourceAssetAsBytes, &sourceAsset); err != nil {
		return shim.Error(err.Error())
	}

	retirementRule := models.NewRetirementRule()
	retirementRule.Id = params.Id
	retirementRule.Name = params.Name
	retirementRule.ExchangeEndpoint = endpoint.Key()
	retirementRule.SourceAsset = sourceAsset.Key()
	retirementRule.Factor = params.Factor
	retirementRule.Fee = params.Fee

	retirementRuleAsBytes, _ := APIstub.GetState(retirementRule.Key())

	// check if AssetRule already exist based on Id
	if retirementRuleAsBytes != nil {
		return shim.Error(fmt.Sprintf("RetirementRule with key %d already exist", params.Id))
	}

	// put new AssetRule to chain
	retirementRuleAsBytes, _ = json.Marshal(retirementRule)
	if err := APIstub.PutState(retirementRule.Key(), retirementRuleAsBytes); err != nil {
		return shim.Error(err.Error())
	}

	// create composite key for AssetRule
	indexName := "RetirementRule"
	retirementRuleIndexKey, errCompositeKey := APIstub.CreateCompositeKey(indexName, []string{retirementRule.Key(), fmt.Sprintf("%d", retirementRule.Id), retirementRule.Name})
	if errCompositeKey != nil {
		return shim.Error(errCompositeKey.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	APIstub.PutState(retirementRuleIndexKey, value)

	return shim.Success(nil)
}
