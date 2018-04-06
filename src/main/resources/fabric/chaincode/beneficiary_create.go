package main

import (
	"encoding/json"
	"fmt"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type createBeneficiaryParams struct {
	Id         int64
	FirstName  string
	MiddleName string
	LastName   string
	Birthdate  string
	Mobile     string
	Email      string
	Address    string
	City       string
	ZipCode    string
	Country    string
	UserKey    string
}

// peer chaincode invoke -n mycc -c '{"Args": ["createBeneficiary", "{\"Id\":1, \"FirstName\":\"Bene\",\"MiddleName\":\"San\",\"LastName\":\"Sebas\",\"Birthdate\":\"03/05/97\",\"Mobile\":\"09362331942\",\"Email\":\"mirvhin143@gmail.com\",\"Address\":\"Alley 21\",\"City\":\"Pateros\",\"ZipCode\":\"1016\",\"Country\":\"Philippines\",\"UserKey\":\"User1\"}"]}' -C myc
func CreateBeneficiary(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := createBeneficiaryParams{}

	// parse params string to struct and throw an error if found
	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	// get User
	userAsBytes, _ := APIstub.GetState(params.UserKey)
	user := models.User{}

	if userAsBytes == nil {
		return shim.Error(fmt.Sprintf("User with key %s not found", params.UserKey))
	}

	if err := json.Unmarshal(userAsBytes, &user); err != nil {
		return shim.Error(err.Error())
	}

	beneficiary := models.NewBeneficiary()
	beneficiary.Id = params.Id
	beneficiary.User = user.Key()
	beneficiary.FirstName = params.FirstName
	beneficiary.MiddleName = params.MiddleName
	beneficiary.LastName = params.LastName
	beneficiary.Birthdate = params.Birthdate
	beneficiary.Mobile = params.Mobile
	beneficiary.Email = params.Email
	beneficiary.Address = params.Email
	beneficiary.City = params.City
	beneficiary.ZipCode = params.ZipCode
	beneficiary.Country = params.Country

	// check if Beneficiary key already exist
	if beneficiaryAsBytes, _ := APIstub.GetState(beneficiary.Key()); beneficiaryAsBytes != nil {
		return shim.Error(fmt.Sprintf("Beneficiary with key %v already exist", params.Id))
	}

	// Put new Beneficiary to chain
	beneficiaryAsBytes, _ := json.Marshal(beneficiary)
	if err := APIstub.PutState(beneficiary.Key(), beneficiaryAsBytes); err != nil {
		return shim.Error(err.Error())
	}

	user.Beneficiaries = append(user.Beneficiaries, beneficiary.Key())
	userAsBytes, _ = json.Marshal(user)
	if err := APIstub.PutState(user.Key(), userAsBytes); err != nil {
		return shim.Error(err.Error())
	}

	// CREATE COMPOSITE KEY FOR USER
	indexName := "Beneficiary"
	beneficiaryIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{beneficiary.Key(),
		fmt.Sprintf("%d", beneficiary.Id), beneficiary.User,
		beneficiary.FirstName, beneficiary.MiddleName, beneficiary.LastName, beneficiary.Birthdate,
		beneficiary.Mobile, beneficiary.Email, beneficiary.Address, beneficiary.City,
		beneficiary.ZipCode, beneficiary.Country})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	APIstub.PutState(beneficiaryIndexKey, value)

	return shim.Success(nil)
}
