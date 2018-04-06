package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// peer chaincode invoke -n mycc -c '{"Args": ["listBeneficiaries", ""]}' -C myc
func ListBeneficiaries(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	resultsIterator, err := APIstub.GetStateByPartialCompositeKey("Beneficiary", []string{})
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

		fmt.Printf("FOUND Beneficiary index:%s id:%s name:%s\n", objectType, compositeKeyParts[0], compositeKeyParts[1])

		// get user group from actual db
		beneficiaryAsBytes, _ := APIstub.GetState(compositeKeyParts[0])

		if beneficiaryAsBytes == nil {
			return shim.Error(fmt.Sprintf("Beneficiary with key %s not found", compositeKeyParts[0]))
		}

		beneficiary := models.Beneficiary{}

		if err := json.Unmarshal(beneficiaryAsBytes, &beneficiary); err != nil {
			return shim.Error(err.Error())
		}

		// TODO check go for text template
		buffer.WriteString("\"Record\": {")
		buffer.WriteString(fmt.Sprintf("\"Key\": \"%s\", ", beneficiary.Key()))
		buffer.WriteString(fmt.Sprintf("\"Id\": \"%d\", ", beneficiary.Id))
		buffer.WriteString(fmt.Sprintf("\"FirstName\": \"%s\", ", beneficiary.FirstName))
		buffer.WriteString(fmt.Sprintf("\"MiddleName\": \"%s\", ", beneficiary.MiddleName))
		buffer.WriteString(fmt.Sprintf("\"LastName\": \"%s\", ", beneficiary.LastName))
		buffer.WriteString(fmt.Sprintf("\"Birthdate\": \"%s\", ", beneficiary.Birthdate))
		buffer.WriteString(fmt.Sprintf("\"Mobile\": \"%s\", ", beneficiary.Mobile))
		buffer.WriteString(fmt.Sprintf("\"Email\": \"%s\", ", beneficiary.Email))
		buffer.WriteString(fmt.Sprintf("\"Address\": \"%s\", ", beneficiary.Address))
		buffer.WriteString(fmt.Sprintf("\"City\": \"%s\", ", beneficiary.City))
		buffer.WriteString(fmt.Sprintf("\"ZipCode\": \"%s\", ", beneficiary.ZipCode))
		buffer.WriteString(fmt.Sprintf("\"Country\": \"%s\", ", beneficiary.Country))

		userAsBytes, _ := APIstub.GetState(beneficiary.User)

		if userAsBytes == nil {
			return shim.Error(fmt.Sprintf("User with key %s not found", beneficiary.User))
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
