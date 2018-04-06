package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// peer chaincode invoke -n mycc -c '{"Args": ["getUser", "{\"Key\":\"User1\"}"]}' -C myc
type getUserParams struct {
	Key string
}

func GetUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := getUserParams{}

	// parse params string to struct and throw an error if found
	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	// get User
	userAsBytes, _ := APIstub.GetState(params.Key)

	if userAsBytes == nil {
		return shim.Error(fmt.Sprintf("User with key %s not found", params.Key))
	}

	user := models.User{}

	if err := json.Unmarshal(userAsBytes, &user); err != nil {
		return shim.Error(err.Error())
	}

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	buffer.WriteString("\"Record\": {")
	buffer.WriteString(fmt.Sprintf("\"Key\": \"%s\", ", user.Key()))
	buffer.WriteString(fmt.Sprintf("\"Id\": \"%d\", ", user.Id))
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

	buffer.WriteString("]")
	buffer.WriteString("}")

	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}
