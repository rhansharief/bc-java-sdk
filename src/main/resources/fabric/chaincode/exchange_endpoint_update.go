package main

import (
	"encoding/json"
	"fmt"

	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type updateExchangeEndpointParams struct {
	Key  string
	Name string
}

// Update ExchangeEndpoint
// peer chaincode invoke -n mycc -c '{"Args": ["updateExchangeEndpoint", "{\"Key\":\"ExchangeEndpoint1\", \"Name\":\"Endpoint 2\"}"]}' -C myc
func UpdateExchangeEndpoint(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	// Parse updateExchangeEndpointParams
	params := updateExchangeEndpointParams{}

	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	// Get endpoint
	endpointAsBytes, _ := APIstub.GetState(params.Key)
	if endpointAsBytes == nil {
		return shim.Error(fmt.Sprintf("ExchangeEndpoint with key %s not found", params.Key))
	}

	endpoint := models.ExchangeEndpoint{}

	if err := json.Unmarshal(endpointAsBytes, &endpoint); err != nil {
		return shim.Error(err.Error())
	}

	endpoint.Name = params.Name

	// Put new user state to chain
	endpointAsBytes, _ = json.Marshal(endpoint)
	if err := APIstub.PutState(endpoint.Key(), endpointAsBytes); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
