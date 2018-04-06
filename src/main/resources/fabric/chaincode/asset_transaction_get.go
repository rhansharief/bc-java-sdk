package main

import (
	"encoding/json"
	"fmt"
	"bytes"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// {"Args": ["getAssetTransaction", "{\"Key\":\"AssetTransaction1\"}"]}
type getAssetTransactionParams struct {
	Key string
}

func GetAssetTransaction(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := getAssetTransactionParams{}

	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	assetTransactionAsBytes, _ := APIstub.GetState(params.Key)
	assetTransaction := models.AssetTransaction{}
	if assetTransactionAsBytes == nil {
		return shim.Error(fmt.Sprintf("Asset Transaction with key %s not found", params.Key))
	}
	if err := json.Unmarshal(assetTransactionAsBytes, &assetTransaction); err != nil {
		return shim.Error(err.Error())
	}

	sourceAssetAsBytes, _ := APIstub.GetState(assetTransaction.SourceAsset)
	sourceAsset := models.Asset{}
	if sourceAssetAsBytes == nil {
		return shim.Error(fmt.Sprintf("Asset with key %s not found", assetTransaction.SourceAsset))
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


	var buffer bytes.Buffer
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

	return shim.Success(buffer.Bytes())
}
