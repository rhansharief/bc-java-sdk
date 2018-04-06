package main

import (
	"encoding/json"
	"fmt"
	"bytes"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// {"Args": ["getIssuanceTransaction", "{\"Key\":\"IssuanceTransaction1\"}"]}
type getIssuanceTransactionParams struct {
	Key string
}

func GetIssuanceTransaction(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := getIssuanceTransactionParams{}

	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	issuanceTransactionAsBytes, _ := APIstub.GetState(params.Key)
	issuanceTransaction := models.IssuanceTransaction{}
	if issuanceTransactionAsBytes == nil {
		return shim.Error(fmt.Sprintf("Issuance Transaction with key %s not found", params.Key))
	}
	if err := json.Unmarshal(issuanceTransactionAsBytes, &issuanceTransaction); err != nil {
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


	var buffer bytes.Buffer
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

	return shim.Success(buffer.Bytes())
}
