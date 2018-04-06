package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
)

// peer chaincode invoke -n mycc -c '{"Args": ["createUserAsset", "{\"Id\":1, \"UserId\":\"User1\", \"AssetId\":\"Asset1\", \"Balance\":100.00}"]}' -C myc
type createUserAssetParams struct {
	Id      int64
	UserId  string
	AssetId string
	Balance float32
}

func CreateUserAsset(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := createUserAssetParams{}

	// parse params string to struct and throw an error if found
	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	_, err := PrivateCreateUserAsset(APIstub, params)
	if err != "" {
		return shim.Error(err)
	}

	return shim.Success(nil)
}
