package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// peer chaincode invoke -n mycc -c '{"Args": ["getIssuanceRule", "{\"Key\":\"IssuanceRule1\"}"]}' -C myc
type getIssuanceRuleParams struct {
	Key string
}

func GetIssuanceRule(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := getIssuanceRuleParams{}

	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	issuanceAsBytes, _ := APIstub.GetState(params.Key)

	if issuanceAsBytes == nil {
		return shim.Error(fmt.Sprintf("Issuance Rule with key %s not found", params.Key))
	}

	return shim.Success(issuanceAsBytes)
}
