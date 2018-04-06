package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// peer chaincode invoke -n mycc -c '{"Args": ["getUserAsset", "{\"Key\":\"UserAsset1\"}"]}' -C myc
type getUserAssetParams struct {
	Key string
}

func GetUserAsset(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := getUserAssetParams{}

	// parse params string to struct and throw an error if found
	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	// get UserAsset
	userAssetAsBytes, _ := APIstub.GetState(params.Key)

	if userAssetAsBytes == nil {
		return shim.Error(fmt.Sprintf("UserAsset with key %s not found", params.Key))
	}

	userAsset := models.UserAsset{}

	if err := json.Unmarshal(userAssetAsBytes, &userAsset); err != nil {
		return shim.Error(err.Error())
	}

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")
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
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}
