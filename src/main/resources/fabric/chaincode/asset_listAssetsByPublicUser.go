package main

import (
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// {"Args": ["listAssetsByPublicUser", "{\"UserKey\":\"User1\"}"]}
type listAssetsByPublicUserParams struct {
	UserKey string
}

func ListAssetsByPublicUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := listAssetsByPublicUserParams{}

	// parse params string to struct and throw an error if found
	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	sourceUserAssets, errAsString := PrivateListAssetsByPublicUser(APIstub, params.UserKey)
	if errAsString != "" {
		return shim.Error(errAsString)
	}

	sourceUserAssetsAsBytes, err := json.Marshal(sourceUserAssets)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(sourceUserAssetsAsBytes)
}
