package main

import (
	"encoding/json"
	"fmt"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// peer chaincode invoke -n mycc -c '{"Args": ["deleteExchangeEndpoint", "{\"Key\":\"ExchangeEndpoint1\"}"]}' -C myc
type deleteExchangeEndpointParams struct {
	Key string
}

func DeleteExchangeEndpoint(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := deleteExchangeEndpointParams{}

	// parse params string to struct and throw an error if found
	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	// check if UsetAsset key doesn't exist
	endpointAsBytes, _ := APIstub.GetState(params.Key)
	if endpointAsBytes == nil {
		return shim.Error(fmt.Sprintf("ExchangeEndpoint with key %s not found", params.Key))
	}

	endpoint := models.ExchangeEndpoint{}

	if err := json.Unmarshal(endpointAsBytes, &endpoint); err != nil {
		return shim.Error(err.Error())
	}

	if err := APIstub.DelState(params.Key); err != nil {
		return shim.Error(err.Error())
	}

	// maintain the index
	indexName := "ExchangeEndpoint"
	endpointIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{endpoint.Key(), fmt.Sprintf("%d", endpoint.Id), endpoint.Name})
	if err != nil {
		return shim.Error(err.Error())
	}

	//  Delete index entry to state.
	err = APIstub.DelState(endpointIndexKey)
	if err != nil {
		return shim.Error("Failed to delete composite key:" + err.Error())
	}

	return shim.Success(nil)
}
