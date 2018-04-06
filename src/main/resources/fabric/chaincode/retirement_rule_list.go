package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// peer chaincode invoke -n mycc -c '{"Args": ["listRetirementRules", ""]}' -C myc
func ListRetirementRules(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	resultsIterator, err := APIstub.GetStateByPartialCompositeKey("RetirementRule", []string{})
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

		fmt.Printf("FOUND Retirement rule index:%s id:%s name:%s\n", objectType, compositeKeyParts[0], compositeKeyParts[1])

		// get user group from actual db
		retirementRuleAsBytes, _ := APIstub.GetState(compositeKeyParts[0])

		if retirementRuleAsBytes == nil {
			return shim.Error(fmt.Sprintf("Retirement Rule with key %s not found", compositeKeyParts[0]))
		}

		retirementRule := models.RetirementRule{}

		if err := json.Unmarshal(retirementRuleAsBytes, &retirementRule); err != nil {
			return shim.Error(err.Error())
		}

		// TODO check go for text template
		buffer.WriteString("\"Record\": {")
		buffer.WriteString(fmt.Sprintf("\"Key\": \"%s\", ", retirementRule.Key()))
		buffer.WriteString(fmt.Sprintf("\"Id\": \"%d\", ", retirementRule.Id))
		buffer.WriteString(fmt.Sprintf("\"Name\": \"%s\", ", retirementRule.Name))
		buffer.WriteString(fmt.Sprintf("\"Factor\": \"%f\", ", retirementRule.Factor))
		buffer.WriteString(fmt.Sprintf("\"Fee\": \"%f\", ", retirementRule.Fee))

		sourceAsBytes, _ := APIstub.GetState(retirementRule.SourceAsset)

		if sourceAsBytes == nil {
			return shim.Error(fmt.Sprintf("Source Asset with key %s not found", retirementRule.SourceAsset))
		}

		source := models.Asset{}

		if err := json.Unmarshal(sourceAsBytes, &source); err != nil {
			return shim.Error(err.Error())
		}

		buffer.WriteString("\"SourceAsset\": {")
		buffer.WriteString(fmt.Sprintf("\"Key\": \"%s\", ", source.Key()))
		buffer.WriteString(fmt.Sprintf("\"Id\": \"%d\", ", source.Id))
		buffer.WriteString(fmt.Sprintf("\"Name\": \"%s\"", source.Name))
		buffer.WriteString("}, ")

		endpointAsBytes, _ := APIstub.GetState(retirementRule.ExchangeEndpoint)

		if endpointAsBytes == nil {
			return shim.Error(fmt.Sprintf("ExchangeEndpoint with key %s not found", retirementRule.ExchangeEndpoint))
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
