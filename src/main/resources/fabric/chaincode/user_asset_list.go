package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// peer chaincode invoke -n mycc -c '{"Args": ["listUserAssets", ""]}' -C myc
func ListUserAssets(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	resultsIterator, err := APIstub.GetStateByPartialCompositeKey("UserAsset", []string{})
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

		fmt.Printf("FOUND UserAsset index:%s id:%s name:%s\n", objectType, compositeKeyParts[0], compositeKeyParts[1])

		// get user group from actual db
		userAssetAsBytes, _ := APIstub.GetState(compositeKeyParts[0])

		if userAssetAsBytes == nil {
			return shim.Error(fmt.Sprintf("User Asset with key %s not found", compositeKeyParts[0]))
		}

		userAsset := models.UserAsset{}

		if err := json.Unmarshal(userAssetAsBytes, &userAsset); err != nil {
			return shim.Error(err.Error())
		}

		// TODO check go for text template
		buffer.WriteString("\"Record\": {")
		buffer.WriteString(fmt.Sprintf("\"Key\": \"%s\", ", userAsset.Key()))
		buffer.WriteString(fmt.Sprintf("\"Id\": \"%d\", ", userAsset.Id))
		buffer.WriteString(fmt.Sprintf("\"Balance\": \"%f\", ", userAsset.Balance))

		assetAsBytes, _ := APIstub.GetState(userAsset.Asset)

		if assetAsBytes == nil {
			return shim.Error(fmt.Sprintf("Asset with key %s not found", userAsset.Asset))
		}

		asset := models.Asset{}

		if err := json.Unmarshal(assetAsBytes, &asset); err != nil {
			return shim.Error(err.Error())
		}

		buffer.WriteString("\"Asset\": {")
		buffer.WriteString(fmt.Sprintf("\"Key\": \"%s\", ", asset.Key()))
		buffer.WriteString(fmt.Sprintf("\"Id\": \"%d\", ", asset.Id))
		buffer.WriteString(fmt.Sprintf("\"Name\": \"%s\", ", asset.Name))
		buffer.WriteString("}, ")

		userAsBytes, _ := APIstub.GetState(userAsset.User)

		if userAsBytes == nil {
			return shim.Error(fmt.Sprintf("User with key %s not found", userAsset.User))
		}

		user := models.User{}

		if err := json.Unmarshal(userAsBytes, &user); err != nil {
			return shim.Error(err.Error())
		}

		buffer.WriteString("\"User\": {")
		buffer.WriteString(fmt.Sprintf("\"Key\": \"%s\", ", user.Key()))
		buffer.WriteString(fmt.Sprintf("\"Id\": \"%d\", ", user.Id))
		buffer.WriteString(fmt.Sprintf("\"FirstName\": \"%s\", ", user.FirstName))
		buffer.WriteString(fmt.Sprintf("\"MiddleName\": \"%s\", ", user.MiddleName))
		buffer.WriteString(fmt.Sprintf("\"LastName\": \"%s\", ", user.LastName))
		buffer.WriteString(fmt.Sprintf("\"Mobile\": \"%s\"", user.Mobile))
		buffer.WriteString("}")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}
