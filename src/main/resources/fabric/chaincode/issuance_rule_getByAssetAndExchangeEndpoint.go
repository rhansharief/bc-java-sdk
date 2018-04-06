package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// peer chaincode invoke -n mycc -c '{"Args": ["getIssuanceRuleByAssetAndExchangeEndpoint", "{\"ExchangeEndpointKey\": \"ExchangeEndpoint1\", \"DestinationAssetKey\":\"Asset1\"}"]}' -C myc
type getIssuanceRuleByAssetAndExchangeEndpointParams struct {
	ExchangeEndpointKey string
	DestinationAssetKey string
}

func GetIssuanceRuleByAssetAndExchangeEndpoint(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := getIssuanceRuleByAssetAndExchangeEndpointParams{}

	// parse params string to struct and throw an error if found
	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	issuanceRuleAsBytes, err := PrivateGetIssuanceRuleByAssetAndExchangeEndpoint(APIstub, params.ExchangeEndpointKey, params.DestinationAssetKey)
	if err != "" {
		return shim.Error(err)
	}
	if issuanceRuleAsBytes == nil {
		return shim.Error(fmt.Sprintf("Issuance Rule with endpoint %s and destination %s not found", params.ExchangeEndpointKey, params.DestinationAssetKey))
	}

	return shim.Success(issuanceRuleAsBytes)
}
