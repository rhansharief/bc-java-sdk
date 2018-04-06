package main

import (
	"encoding/json"
	"fmt"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// peer chaincode invoke -n mycc -c '{"Args": ["createIssuanceRule", "{\"Id\":1, \"Name\":\"IssuanceRule\", \"ExchangeEndpointKey\": \"ExchangeEndpoint1\", \"DestinationAssetKey\":\"Asset1\", \"Factor\":2, \"Fee\":0.50}"]}' -C myc
type createIssuanceRuleParams struct {
	Id                  int64
	Name                string
	ExchangeEndpointKey string
	DestinationAssetKey string
	Factor              float32
	Fee                 float32
}

func CreateIssuanceRule(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := createIssuanceRuleParams{}

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
	destinationAssetAsBytes, _ := APIstub.GetState(params.DestinationAssetKey)
	destinationAsset := models.Asset{}

	if destinationAssetAsBytes == nil {
		return shim.Error(fmt.Sprintf("Asset with key %s not found", params.DestinationAssetKey))
	}

	if err := json.Unmarshal(destinationAssetAsBytes, &destinationAsset); err != nil {
		return shim.Error(err.Error())
	}

	issuanceRule := models.NewIssuanceRule()
	issuanceRule.Id = params.Id
	issuanceRule.Name = params.Name
	issuanceRule.ExchangeEndpoint = endpoint.Key()
	issuanceRule.DestinationAsset = destinationAsset.Key()
	issuanceRule.Factor = params.Factor
	issuanceRule.Fee = params.Fee

	issuanceRuleAsBytes, _ := APIstub.GetState(issuanceRule.Key())

	// check if AssetRule already exist based on Id
	if issuanceRuleAsBytes != nil {
		return shim.Error(fmt.Sprintf("IssuanceRule with key %d already exist", params.Id))
	}

	// put new AssetRule to chain
	issuanceRuleAsBytes, _ = json.Marshal(issuanceRule)
	if err := APIstub.PutState(issuanceRule.Key(), issuanceRuleAsBytes); err != nil {
		return shim.Error(err.Error())
	}

	// create composite key for AssetRule
	indexName := "IssuanceRule"
	issuanceRuleIndexKey, errCompositeKey := APIstub.CreateCompositeKey(indexName, []string{issuanceRule.Key(), fmt.Sprintf("%d", issuanceRule.Id), issuanceRule.Name})
	if errCompositeKey != nil {
		return shim.Error(errCompositeKey.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	APIstub.PutState(issuanceRuleIndexKey, value)

	return shim.Success(nil)
}
