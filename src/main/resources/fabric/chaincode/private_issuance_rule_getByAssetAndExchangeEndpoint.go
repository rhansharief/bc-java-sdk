package main

import (
	"bytes"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func PrivateGetIssuanceRuleByAssetAndExchangeEndpoint(APIstub shim.ChaincodeStubInterface, exchangeEndpointKey string, destinationAssetKey string) ([]byte, string) {
	query := fmt.Sprintf("{\"selector\":{\"DocType\":\"ISSUANCE_RULE\", \"ExchangeEndpoint\":\"%s\", \"DestinationAsset\":\"%s\"}}", exchangeEndpointKey, destinationAssetKey)
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
