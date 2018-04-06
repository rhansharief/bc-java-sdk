package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// peer chaincode invoke -n mycc -c '{"Args": ["listUsers", ""]}' -C myc
func ListUsers(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	resultsIterator, err := APIstub.GetStateByPartialCompositeKey("User", []string{})
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

		fmt.Printf("FOUND User index:%s id:%s name:%s\n", objectType, compositeKeyParts[0], compositeKeyParts[1])

		// get user group from actual db
		userAsBytes, _ := APIstub.GetState(compositeKeyParts[0])

		if userAsBytes == nil {
			return shim.Error(fmt.Sprintf("User with key %s not found", compositeKeyParts[0]))
		}

		user := models.User{}

		if err := json.Unmarshal(userAsBytes, &user); err != nil {
			return shim.Error(err.Error())
		}

		// TODO check go for text template
		buffer.WriteString("\"Record\": {")
		buffer.WriteString(fmt.Sprintf("\"Key\": \"%s\", ", user.Key()))
		buffer.WriteString(fmt.Sprintf("\"Id\": \"%d\", ", user.Id))
		buffer.WriteString(fmt.Sprintf("\"UID\": \"%s\", ", user.UID))
		buffer.WriteString(fmt.Sprintf("\"Email\": \"%s\", ", user.Email))
		buffer.WriteString(fmt.Sprintf("\"Username\": \"%s\", ", user.Username))
		buffer.WriteString(fmt.Sprintf("\"Address\": \"%s\", ", user.Address))
		buffer.WriteString(fmt.Sprintf("\"FirstName\": \"%s\", ", user.FirstName))
		buffer.WriteString(fmt.Sprintf("\"MiddleName\": \"%s\", ", user.MiddleName))
		buffer.WriteString(fmt.Sprintf("\"LastName\": \"%s\", ", user.LastName))
		buffer.WriteString(fmt.Sprintf("\"Mobile\": \"%s\", ", user.Mobile))
		buffer.WriteString("\"UserGroups:\" [")

		elementAlreadyWritten := false
		for _, element := range user.UserGroups {
			if elementAlreadyWritten == true {
				buffer.WriteString(",")
			}

			groupAsBytes, _ := APIstub.GetState(element)

			if groupAsBytes == nil {
				return shim.Error(fmt.Sprintf("Group with key %s not found", element))
			}

			group := models.UserGroup{}

			if err := json.Unmarshal(groupAsBytes, &group); err != nil {
				return shim.Error(err.Error())
			}

			buffer.WriteString("{")
			buffer.WriteString(fmt.Sprintf("\"Id\": \"%d\", ", group.Id))
			buffer.WriteString(fmt.Sprintf("\"Name\": \"%s\"", group.Name))
			buffer.WriteString("}")

			elementAlreadyWritten = true
		}

		buffer.WriteString("],")

		buffer.WriteString("\"Beneficiaries:\" [")

		beneficiaryAlreadyWritten := false
		for _, element := range user.Beneficiaries {
			if beneficiaryAlreadyWritten == true {
				buffer.WriteString(",")
			}

			beneficiaryAsBytes, _ := APIstub.GetState(element)

			if beneficiaryAsBytes == nil {
				return shim.Error(fmt.Sprintf("Beneficiary with key %s not found", element))
			}

			beneficiary := models.Beneficiary{}

			if err := json.Unmarshal(beneficiaryAsBytes, &beneficiary); err != nil {
				return shim.Error(err.Error())
			}

			buffer.WriteString("{")
			buffer.WriteString(fmt.Sprintf("\"Key\": \"%s\" ", beneficiary.Key()))
			buffer.WriteString("}")

			beneficiaryAlreadyWritten = true
		}

		buffer.WriteString("]")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}
