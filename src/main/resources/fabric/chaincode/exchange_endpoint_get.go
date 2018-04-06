package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// peer chaincode invoke -n mycc -c '{"Args": ["getExchangeEndpoint", "{\"Key\":\"ExchangeEndpoint1\"}"]}' -C myc
type getExchangeEndpointParams struct {
	Key string
}

// TODO check if working after create UserAsset is done
func GetExchangeEndpoint(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := getExchangeEndpointParams{}

	// parse params string to struct and throw an error if found
	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	endpointAsBytes, _ := APIstub.GetState(params.Key)

	if endpointAsBytes == nil {
		return shim.Error(fmt.Sprintf("ExchangeEndpoint with key %s not found", params.Key))
	}

	return shim.Success(endpointAsBytes)
}
