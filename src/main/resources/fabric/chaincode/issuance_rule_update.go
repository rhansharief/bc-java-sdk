package main

import (
	"encoding/json"
	"fmt"

	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type updateIssuanceRuleParams struct {
	Key                 string
	Name                string
	ExchangeEndpointKey string
	DestinationAssetKey string
	Factor              float32
	Fee                 float32
}

// Update IssuanceRule
// peer chaincode invoke -n mycc -c '{"Args": ["updateIssuanceRule", "{\"Key\":\"IssuanceRule1\", \"Name\":\"IssuanceRule 111\", \"ExchangeEndpointKey\": \"ExchangeEndpoint1\", \"DestinationAssetKey\":\"Asset1\", \"Factor\":50.50, \"Fee\":0.50}"]}' -C myc
func UpdateIssuanceRule(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	// Parse updateExchangeEndpointParams
	params := updateIssuanceRuleParams{}

	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	// Get endpoint
	issuanceAsBytes, _ := APIstub.GetState(params.Key)
	if issuanceAsBytes == nil {
		return shim.Error(fmt.Sprintf("IssuanceRule with key %s not found", params.Key))
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

	issuance := models.IssuanceRule{}

	if err := json.Unmarshal(issuanceAsBytes, &issuance); err != nil {
		return shim.Error(err.Error())
	}

	issuance.Name = params.Name
	issuance.ExchangeEndpoint = endpoint.Key()
	issuance.DestinationAsset = destinationAsset.Key()
	issuance.Factor = params.Factor
	issuance.Fee = params.Fee

	// Put new user state to chain
	issuanceAsBytes, _ = json.Marshal(issuance)
	if err := APIstub.PutState(issuance.Key(), issuanceAsBytes); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
