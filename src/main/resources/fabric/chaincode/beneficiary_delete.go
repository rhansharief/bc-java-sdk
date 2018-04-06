package main

import (
	"encoding/json"
	"fmt"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// peer chaincode invoke -n mycc -c '{"Args": ["deleteBeneficiary", "{\"Key\":\"Beneficiary1\"}"]}' -C myc
type deleteBeneficiaryParams struct {
	Key string
}

func DeleteBeneficiary(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := deleteBeneficiaryParams{}

	// parse params string to struct and throw an error if found
	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	// check if Beneficiary key doesn't exist
	beneficiaryAsBytes, _ := APIstub.GetState(params.Key)
	if beneficiaryAsBytes == nil {
		return shim.Error(fmt.Sprintf("Beneficiary with key %s not found", params.Key))
	}

	beneficiary := models.Beneficiary{}

	if err := json.Unmarshal(beneficiaryAsBytes, &beneficiary); err != nil {
		return shim.Error(err.Error())
	}

	// couch query string - get all objects that has DocType="USER" AND Beneficiaries contains params.Key
	query := fmt.Sprintf(""+
		"{\"selector\": {"+
		"\"$and\": [{"+
		"\"DocType\": \"USER\","+
		"\"Beneficiaries\": {"+
		"\"$in\": [\"%s\"]"+
		"}"+
		"}]"+
		"}}", params.Key)

	resultsIterator, err := APIstub.GetQueryResult(query)

	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// loop through all users attach to a beneficiary
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		user := models.User{}
		if err := json.Unmarshal(queryResponse.Value, &user); err != nil {
			return shim.Error(err.Error())
		}

		// get index of the element we want to remove
		i := indexOf(params.Key, user.Beneficiaries)

		// remove element from array by index
		user.Beneficiaries = removeKey(user.Beneficiaries, i)

		// save updated user to state
		userAsBytes, _ := json.Marshal(user)
		if err := APIstub.PutState(user.Key(), userAsBytes); err != nil {
			return shim.Error(err.Error())
		}
	}

	if err := APIstub.DelState(params.Key); err != nil {
		return shim.Error(err.Error())
	}

	// maintain the index
	indexName := "Beneficiary"
	beneficiaryIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{beneficiary.Key(),
		fmt.Sprintf("%d", beneficiary.Id), beneficiary.User,
		beneficiary.FirstName, beneficiary.MiddleName, beneficiary.LastName, beneficiary.Birthdate,
		beneficiary.Mobile, beneficiary.Email, beneficiary.Address, beneficiary.City,
		beneficiary.ZipCode, beneficiary.Country})
	if err != nil {
		return shim.Error(err.Error())
	}

	//  Delete index entry to state.
	err = APIstub.DelState(beneficiaryIndexKey)
	if err != nil {
		return shim.Error("Failed to delete composite key:" + err.Error())
	}

	return shim.Success(nil)
}
