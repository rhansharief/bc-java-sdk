package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// peer chaincode invoke -n mycc -c '{"Args": ["getUserGroup", "{\"Key\":\"UserGroup1\"}"]}' -C myc
type getUserGroupParams struct {
	Key string
}

func GetUserGroup(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := getUserGroupParams{}

	// parse params string to struct and throw an error if found
	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	groupAsBytes, _ := APIstub.GetState(params.Key)

	if groupAsBytes == nil {
		return shim.Error(fmt.Sprintf("UserGroup with key %s not found", params.Key))
	}

	group := models.UserGroup{}

	if err := json.Unmarshal(groupAsBytes, &group); err != nil {
		return shim.Error(err.Error())
	}

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

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

	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}
