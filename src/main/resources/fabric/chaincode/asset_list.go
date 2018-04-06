package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// peer chaincode invoke -n mycc -c '{"Args": ["listAssets", ""]}' -C myc
func ListAssets(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	resultsIterator, err := APIstub.GetStateByPartialCompositeKey("Asset", []string{})
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

		fmt.Printf("Found Asset index:%s id:%s name:%s\n", objectType, compositeKeyParts[0], compositeKeyParts[1])

		// get user group from actual db
		assetAsBytes, _ := APIstub.GetState(compositeKeyParts[0])

		if assetAsBytes == nil {
			return shim.Error(fmt.Sprintf("Asset with key %s not found", compositeKeyParts[0]))
		}

		asset := models.Asset{}

		if err := json.Unmarshal(assetAsBytes, &asset); err != nil {
			return shim.Error(err.Error())
		}

		// TODO check go for text template
		buffer.WriteString("\"Record\": {")
		buffer.WriteString(fmt.Sprintf("\"Key\": \"%s\", ", asset.Key()))
		buffer.WriteString(fmt.Sprintf("\"Id\": \"%d\", ", asset.Id))
		buffer.WriteString(fmt.Sprintf("\"Name\": \"%s\"", asset.Name))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}
