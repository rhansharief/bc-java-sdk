package main

import (
	"bytes"
	"fmt"
	"encoding/json"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// {"Args": ["listIssuanceTransactionsByQuery"]}
func ListIssuanceTransactionsByQuery(APIstub shim.ChaincodeStubInterface) sc.Response {
	query := "{\"selector\":{\"DocType\":\"ISSUANCE_TRANSACTION\"}}"
	resultsIterator, err := APIstub.GetQueryResult(query)

	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(", ")
		}

		issuanceTransaction := models.IssuanceTransaction{}
		if err := json.Unmarshal(queryResponse.Value, &issuanceTransaction); err != nil {
			return shim.Error(err.Error())
		}

		brokerAsBytes, _ := APIstub.GetState(issuanceTransaction.Broker)
		broker := models.User{}
		if brokerAsBytes == nil {
			return shim.Error(fmt.Sprintf("User with key %s not found", issuanceTransaction.Broker))
		}
		if err := json.Unmarshal(brokerAsBytes, &broker); err != nil {
			return shim.Error(err.Error())
		}

		exchangeEndpointAsBytes, _ := APIstub.GetState(issuanceTransaction.ExchangeEndpoint)
		exchangeEndpoint := models.ExchangeEndpoint{}
		if exchangeEndpointAsBytes == nil {
			return shim.Error(fmt.Sprintf("Exchange Endpoint with key %s not found", issuanceTransaction.ExchangeEndpoint))
		}
		if err := json.Unmarshal(exchangeEndpointAsBytes, &exchangeEndpoint); err != nil {
			return shim.Error(err.Error())
		}

		destinationAssetAsBytes, _ := APIstub.GetState(issuanceTransaction.DestinationAsset)
		destinationAsset := models.Asset{}
		if destinationAssetAsBytes == nil {
			return shim.Error(fmt.Sprintf("Destination Asset with key %s not found", issuanceTransaction.DestinationAsset))
		}
		if err := json.Unmarshal(destinationAssetAsBytes, &destinationAsset); err != nil {
			return shim.Error(err.Error())
		}

		// TODO check go for text template
		buffer.WriteString("{")

		buffer.WriteString(fmt.Sprintf("\"Key\": \"%s\", ", issuanceTransaction.Key()))
		buffer.WriteString(fmt.Sprintf("\"Id\": \"%d\", ", issuanceTransaction.Id))
		buffer.WriteString(fmt.Sprintf("\"Factor\": \"%f\", ", issuanceTransaction.Factor))
		buffer.WriteString(fmt.Sprintf("\"Fee\": \"%f\", ", issuanceTransaction.Fee))
		buffer.WriteString(fmt.Sprintf("\"Amount\": \"%f\", ", issuanceTransaction.Amount))

		buffer.WriteString("\"Broker\": {")
		buffer.WriteString(fmt.Sprintf("\"Key\": \"%s\", ", broker.Key()))
		buffer.WriteString(fmt.Sprintf("\"Id\": \"%d\", ", broker.Id))
		buffer.WriteString(fmt.Sprintf("\"FirstName\": \"%s\", ", broker.FirstName))
		buffer.WriteString(fmt.Sprintf("\"MiddleName\": \"%s\", ", broker.MiddleName))
		buffer.WriteString(fmt.Sprintf("\"LastName\": \"%s\", ", broker.LastName))
		buffer.WriteString(fmt.Sprintf("\"Mobile\": \"%s\"", broker.Mobile))
		buffer.WriteString("}, ")

		buffer.WriteString("\"ExchangeEndpoint\": {")
		buffer.WriteString(fmt.Sprintf("\"Key\": \"%s\", ", exchangeEndpoint.Key()))
		buffer.WriteString(fmt.Sprintf("\"Id\": \"%d\", ", exchangeEndpoint.Id))
		buffer.WriteString(fmt.Sprintf("\"Name\": \"%s\"", exchangeEndpoint.Name))
		buffer.WriteString("}, ")

		buffer.WriteString("\"DestinationAsset\": {")
		buffer.WriteString(fmt.Sprintf("\"Key\": \"%s\", ", destinationAsset.Key()))
		buffer.WriteString(fmt.Sprintf("\"Id\": \"%d\", ", destinationAsset.Id))
		buffer.WriteString(fmt.Sprintf("\"Name\": \"%s\"", destinationAsset.Name))
		buffer.WriteString("}")

		buffer.WriteString("}")

		bArrayMemberAlreadyWritten = true
	}

	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}