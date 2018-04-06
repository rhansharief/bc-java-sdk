package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// peer chaincode invoke -n mycc -c '{"Args": ["getUserByUid", "{\"UID\":\"176000120221\"}"]}' -C myc
type getUserByUIDParams struct {
	UID string
}

// GetUserByUID function
func GetUserByUID(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := getUserByUIDParams{}

	// parse params string to struct and throw an error if found
	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	// get User
	query := fmt.Sprintf("{\"selector\":{\"UID\":\"%s\"}}", params.UID)
	resultsIterator, err := APIstub.GetQueryResult(query)

	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON string
	var buffer bytes.Buffer

	if !resultsIterator.HasNext() {
		return shim.Error(fmt.Sprintf("User with UID %s", params.UID))
	}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		buffer.WriteString(string(queryResponse.Value))
	}

	return shim.Success(buffer.Bytes())
}
