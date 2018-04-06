package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// peer chaincode invoke -n mycc -c '{"Args": ["getBeneficiary", "{\"Key\":\"Beneficiary1\"}"]}' -C myc
type getBeneficiaryParams struct {
	Key string
}

func GetBeneficiary(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := getBeneficiaryParams{}

	// parse params string to struct and throw an error if found
	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	// get Beneficiary
	beneficiaryAsBytes, _ := APIstub.GetState(params.Key)

	if beneficiaryAsBytes == nil {
		return shim.Error(fmt.Sprintf("Beneficiary with key %s not found", params.Key))
	}

	beneficiary := models.Beneficiary{}

	if err := json.Unmarshal(beneficiaryAsBytes, &beneficiary); err != nil {
		return shim.Error(err.Error())
	}

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")
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
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}
