package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// {"Args": ["deleteUserAsset", "{\"Key\":\"UserAsset1\"}"]}
type deleteUserAssetParams struct {
	Key string
}

// TODO check if working after create UserAsset is done
func DeleteUserAsset(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := deleteUserAssetParams{}

	// parse params string to struct and throw an error if found
	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	// check if UsetAsset key doesn't exist
	userAssetAsBytes, _ := APIstub.GetState(params.Key)
	if userAssetAsBytes == nil {
		return shim.Error(fmt.Sprintf("UserAsset with key %s not found", params.Key))
	}

	if err := APIstub.DelState(params.Key); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
