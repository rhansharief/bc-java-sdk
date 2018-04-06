package main

import (
	"bytes"
	"fmt"
	"encoding/json"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// {"Args": ["listAssetTransactionsByCompositeKey"]}
func ListAssetTransactionsByCompositeKey(APIstub shim.ChaincodeStubInterface) sc.Response {
	resultsIterator, err := APIstub.GetStateByPartialCompositeKey("AssetTransaction", []string{})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
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

		objectType, compositeKeyParts, err := APIstub.SplitCompositeKey(queryResponse.Key)
		if err != nil {
			return shim.Error(err.Error())
		}
		fmt.Printf("Found Asset Transaction index:%s id:%s \n", objectType, compositeKeyParts[0])
		fmt.Println(compositeKeyParts[0])
		fmt.Println("key")
		assetTransactionKey := compositeKeyParts[0]
		assetTransactionAsBytes, _ := APIstub.GetState(assetTransactionKey)
		assetTransaction := models.AssetTransaction{}
		if err := json.Unmarshal(assetTransactionAsBytes, &assetTransaction); err != nil {
			return shim.Error(err.Error())
		}

		sourceAssetAsBytes, _ := APIstub.GetState(assetTransaction.SourceAsset)
		sourceAsset := models.Asset{}
		if sourceAssetAsBytes == nil {
			return shim.Error(fmt.Sprintf("Source Asset with key %s not found", assetTransaction.SourceAsset))
		}
		if err := json.Unmarshal(sourceAssetAsBytes, &sourceAsset); err != nil {
			return shim.Error(err.Error())
		}

		destinationAssetAsBytes, _ := APIstub.GetState(assetTransaction.DestinationAsset)
		destinationAsset := models.Asset{}
		if destinationAssetAsBytes == nil {
			return shim.Error(fmt.Sprintf("Destination Asset with key %s not found", assetTransaction.DestinationAsset))
		}
		if err := json.Unmarshal(destinationAssetAsBytes, &destinationAsset); err != nil {
			return shim.Error(err.Error())
		}

		buffer.WriteString("{")

		buffer.WriteString(fmt.Sprintf("\"Key\": \"%s\", ", assetTransaction.Key()))
		buffer.WriteString(fmt.Sprintf("\"Id\": \"%d\", ", assetTransaction.Id))
		buffer.WriteString(fmt.Sprintf("\"Factor\": \"%f\", ", assetTransaction.Factor))
		buffer.WriteString(fmt.Sprintf("\"Fee\": \"%f\", ", assetTransaction.Fee))
		buffer.WriteString(fmt.Sprintf("\"Amount\": \"%f\", ", assetTransaction.Amount))

		buffer.WriteString("\"SourceAsset\": {")
		buffer.WriteString(fmt.Sprintf("\"Key\": \"%s\", ", sourceAsset.Key()))
		buffer.WriteString(fmt.Sprintf("\"Id\": \"%d\", ", sourceAsset.Id))
		buffer.WriteString(fmt.Sprintf("\"Name\": \"%s\"", sourceAsset.Name))
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
