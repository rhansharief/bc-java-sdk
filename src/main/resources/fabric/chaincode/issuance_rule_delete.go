package main

import (
	"encoding/json"
	"fmt"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// peer chaincode invoke -n mycc -c '{"Args": ["deleteIssuanceRule", "{\"Key\":\"IssuanceRule1\"}"]}' -C myc
type deleteIssuanceRuleParams struct {
	Key string
}

func DeleteIssuanceRule(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := deleteIssuanceRuleParams{}

	// parse params string to struct and throw an error if found
	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	// check if IssuanceRule key doesn't exist
	issuanceRuleAsBytes, _ := APIstub.GetState(params.Key)
	if issuanceRuleAsBytes == nil {
		return shim.Error(fmt.Sprintf("IssuanceRule with key %s not found", params.Key))
	}

	issuanceRule := models.IssuanceRule{}

	if err := json.Unmarshal(issuanceRuleAsBytes, &issuanceRule); err != nil {
		return shim.Error(err.Error())
	}

	if err := APIstub.DelState(params.Key); err != nil {
		return shim.Error(err.Error())
	}

	// maintain the index
	indexName := "IssuanceRule"
	issuanceRuleIndexKey, errCompositeKey := APIstub.CreateCompositeKey(indexName, []string{issuanceRule.Key(), fmt.Sprintf("%d", issuanceRule.Id), issuanceRule.Name})
	if errCompositeKey != nil {
		return shim.Error(errCompositeKey.Error())
	}

	//  Delete index entry to state.
	err := APIstub.DelState(issuanceRuleIndexKey)
	if err != nil {
		return shim.Error("Failed to delete composite key:" + err.Error())
	}

	return shim.Success(nil)
}
