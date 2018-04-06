package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// peer chaincode invoke -n mycc -c '{"Args": ["listExchangeEndpoints", ""]}' -C myc
func ListExchangeEndpoints(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	resultsIterator, err := APIstub.GetStateByPartialCompositeKey("ExchangeEndpoint", []string{})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		objectType, compositeKeyParts, err := APIstub.SplitCompositeKey(queryResponse.Key)
		if err != nil {
			return shim.Error(err.Error())
		}

		fmt.Printf("Found Exchange Endpoint index:%s id:%s name:%s\n", objectType, compositeKeyParts[0], compositeKeyParts[1])

		// get user group from actual db
		endpointAsBytes, _ := APIstub.GetState(compositeKeyParts[0])

		if endpointAsBytes == nil {
			return shim.Error(fmt.Sprintf("ExchangeEndpoint with key %s not found", compositeKeyParts[0]))
		}

		endpoint := models.ExchangeEndpoint{}

		if err := json.Unmarshal(endpointAsBytes, &endpoint); err != nil {
			return shim.Error(err.Error())
		}

		// TODO check go for text template
		buffer.WriteString("\"Record\": {")
		buffer.WriteString(fmt.Sprintf("\"Key\": \"%s\", ", endpoint.Key()))
		buffer.WriteString(fmt.Sprintf("\"Id\": \"%d\", ", endpoint.Id))
		buffer.WriteString(fmt.Sprintf("\"Name\": \"%s\"", endpoint.Name))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}
