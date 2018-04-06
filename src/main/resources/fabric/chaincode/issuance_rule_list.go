package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// peer chaincode invoke -n mycc -c '{"Args": ["listIssuanceRules", ""]}' -C myc
func ListIssuanceRules(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	resultsIterator, err := APIstub.GetStateByPartialCompositeKey("IssuanceRule", []string{})
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

		fmt.Printf("FOUND Issuance rule index:%s id:%s name:%s\n", objectType, compositeKeyParts[0], compositeKeyParts[1])

		// get user group from actual db
		issuanceRuleAsBytes, _ := APIstub.GetState(compositeKeyParts[0])

		if issuanceRuleAsBytes == nil {
			return shim.Error(fmt.Sprintf("Issuance Rule with key %s not found", compositeKeyParts[0]))
		}

		issuanceRule := models.IssuanceRule{}

		if err := json.Unmarshal(issuanceRuleAsBytes, &issuanceRule); err != nil {
			return shim.Error(err.Error())
		}

		// TODO check go for text template
		buffer.WriteString("\"Record\": {")
		buffer.WriteString(fmt.Sprintf("\"Key\": \"%s\", ", issuanceRule.Key()))
		buffer.WriteString(fmt.Sprintf("\"Id\": \"%d\", ", issuanceRule.Id))
		buffer.WriteString(fmt.Sprintf("\"Name\": \"%s\", ", issuanceRule.Name))
		buffer.WriteString(fmt.Sprintf("\"Factor\": \"%f\", ", issuanceRule.Factor))
		buffer.WriteString(fmt.Sprintf("\"Fee\": \"%f\", ", issuanceRule.Fee))

		destinationAsBytes, _ := APIstub.GetState(issuanceRule.DestinationAsset)

		if destinationAsBytes == nil {
			return shim.Error(fmt.Sprintf("Destination Asset with key %s not found", issuanceRule.DestinationAsset))
		}

		destination := models.Asset{}

		if err := json.Unmarshal(destinationAsBytes, &destination); err != nil {
			return shim.Error(err.Error())
		}

		buffer.WriteString("\"DestinationAsset\": {")
		buffer.WriteString(fmt.Sprintf("\"Key\": \"%s\", ", destination.Key()))
		buffer.WriteString(fmt.Sprintf("\"Id\": \"%d\", ", destination.Id))
		buffer.WriteString(fmt.Sprintf("\"Name\": \"%s\", ", destination.Name))
		buffer.WriteString("}, ")

		endpointAsBytes, _ := APIstub.GetState(issuanceRule.ExchangeEndpoint)

		if endpointAsBytes == nil {
			return shim.Error(fmt.Sprintf("ExchangeEndpoint with key %s not found", issuanceRule.ExchangeEndpoint))
		}

		endpoint := models.ExchangeEndpoint{}

		if err := json.Unmarshal(endpointAsBytes, &endpoint); err != nil {
			return shim.Error(err.Error())
		}

		buffer.WriteString("\"ExchangeEndpoint\": {")
		buffer.WriteString(fmt.Sprintf("\"Key\": \"%s\", ", endpoint.Key()))
		buffer.WriteString(fmt.Sprintf("\"Id\": \"%d\", ", endpoint.Id))
		buffer.WriteString(fmt.Sprintf("\"Name\": \"%s\"", endpoint.Name))
		buffer.WriteString("}")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}
