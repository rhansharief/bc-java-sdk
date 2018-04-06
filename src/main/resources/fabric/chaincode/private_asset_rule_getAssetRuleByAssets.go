package main

import (
	"bytes"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
)

// Get asset rule with property DocType=ASSET_RULE, Source=Asset1 and Destination=Asset2
// Params:
//    sourceAssetKey="Asset1"
//    destinationAssetKey="Asset2"
func PrivateGetAssetRuleByAssets(APIstub shim.ChaincodeStubInterface, sourceAssetKey string, destinationAssetKey string) ([]byte, string) {
	query := fmt.Sprintf("{\"selector\":{\"DocType\":\"ASSET_RULE\", \"Source\":\"%s\", \"Destination\":\"%s\"}}", sourceAssetKey, destinationAssetKey )
	resultsIterator, err := APIstub.GetQueryResult(query)

	if err != nil {
		return nil, err.Error()
	}
	defer resultsIterator.Close()

	// buffer is a JSON string
	var buffer bytes.Buffer

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err.Error()
		}

		buffer.WriteString(string(queryResponse.Value))
	}

	return buffer.Bytes(), ""
}