package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// peer chaincode invoke -n mycc -c '{"Args": ["listUserGroups", ""]}' -C myc
func ListUserGroups(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	resultsIterator, err := APIstub.GetStateByPartialCompositeKey("UserGroup", []string{})
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

		fmt.Printf("FOUND User Group index:%s id:%s name:%s\n", objectType, compositeKeyParts[0], compositeKeyParts[1])

		// get user group from actual db
		groupAsBytes, _ := APIstub.GetState(compositeKeyParts[0])

		if groupAsBytes == nil {
			return shim.Error(fmt.Sprintf("UserGroup with key %s not found", compositeKeyParts[0]))
		}

		group := models.UserGroup{}

		if err := json.Unmarshal(groupAsBytes, &group); err != nil {
			return shim.Error(err.Error())
		}

		// TODO check go for text template
		buffer.WriteString("\"Record\": {")
		buffer.WriteString(fmt.Sprintf("\"Key\": \"%s\", ", group.Key()))
		buffer.WriteString(fmt.Sprintf("\"Id\": \"%d\", ", group.Id))
		buffer.WriteString(fmt.Sprintf("\"Name\": \"%s\", ", group.Name))
		buffer.WriteString("\"Assets:\" [")

		elementAlreadyWritten := false
		for _, element := range group.Assets {
			if elementAlreadyWritten == true {
				buffer.WriteString(",")
			}

			assetAsBytes, _ := APIstub.GetState(element)

			if assetAsBytes == nil {
				return shim.Error(fmt.Sprintf("Asset with key %s not found", element))
			}

			asset := models.Asset{}

			if err := json.Unmarshal(assetAsBytes, &asset); err != nil {
				return shim.Error(err.Error())
			}

			buffer.WriteString("{")
			buffer.WriteString(fmt.Sprintf("\"Id\": \"%d\", ", asset.Id))
			buffer.WriteString(fmt.Sprintf("\"Name\": \"%s\"", asset.Name))
			buffer.WriteString("}")

			elementAlreadyWritten = true
		}

		buffer.WriteString("]")
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}
