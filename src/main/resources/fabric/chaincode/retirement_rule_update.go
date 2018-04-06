package main

import (
	"encoding/json"
	"fmt"

	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type updateRetirementRuleParams struct {
	Key    string
	Factor float32
	Fee    float32
}

// peer chaincode invoke -n mycc -c '{"Args": ["updateRetirementRule", "{\"Key\":\"RetirementRule1\", \"Factor\":100.10, \"Fee\":10.10}"]}' -C myc
func UpdateRetirementRule(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	// Parse updateExchangeEndpointParams
	params := updateRetirementRuleParams{}

	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	// Get endpoint
	retirementAsBytes, _ := APIstub.GetState(params.Key)
	if retirementAsBytes == nil {
		return shim.Error(fmt.Sprintf("RetirementRule with key %s not found", params.Key))
	}

	retirement := models.RetirementRule{}

	if err := json.Unmarshal(retirementAsBytes, &retirement); err != nil {
		return shim.Error(err.Error())
	}

	retirement.Factor = params.Factor
	retirement.Fee = params.Fee

	// Put new user state to chain
	retirementAsBytes, _ = json.Marshal(retirement)
	if err := APIstub.PutState(retirement.Key(), retirementAsBytes); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
