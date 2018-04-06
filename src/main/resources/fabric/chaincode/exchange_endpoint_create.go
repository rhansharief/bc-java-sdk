package main

import (
	"encoding/json"
	"fmt"

	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type createExchangeEndpointParams struct {
	Name string
	Id   int64
}

// Create ExchangeEndpoint
// peer chaincode invoke -n mycc -c '{"Args": ["createExchangeEndpoint", "{\"Id\":1, \"Name\":\"Endpoint 1\"}"]}' -C myc
func CreateExchangeEndpoint(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	// Parse createExchangeEndpointParams
	params := createExchangeEndpointParams{}

	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	endpoint := models.NewExchangeEndpoint()

	endpoint.Name = params.Name
	endpoint.Id = params.Id

	// check if User key already exist
	exchangeAsBytes, _ := APIstub.GetState(endpoint.Key())
	if exchangeAsBytes != nil {
		return shim.Error(fmt.Sprintf("Exchange endpoint with key %s already exist", endpoint.Key()))
	}

	// Put new endpoint state to chain
	endpointAsBytes, _ := json.Marshal(endpoint)
	if err := APIstub.PutState(endpoint.Key(), endpointAsBytes); err != nil {
		return shim.Error(err.Error())
	}

	// CREATE COMPOSITE KEY FOR USER GROUP
	indexName := "ExchangeEndpoint"
	endpointIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{endpoint.Key(), fmt.Sprintf("%d", endpoint.Id), endpoint.Name})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	APIstub.PutState(endpointIndexKey, value)

	return shim.Success(nil)
}
