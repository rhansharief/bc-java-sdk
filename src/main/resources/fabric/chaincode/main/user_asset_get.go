package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// {"Args": ["getUserAsset", "{\"Key\":\"UserAsset1\"}"]}
type getUserAssetParams struct {
	Key string
}

// TODO check if working after create UserAsset is done
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

	return shim.Success(userAssetAsBytes)
}
