package main

import (
	"bytes"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// {"Args": ["listAssetRuleByCompositeKey"]}
func ListAssetRulesByCompositeKey(APIstub shim.ChaincodeStubInterface) sc.Response {
	resultsIterator, err := APIstub.GetStateByPartialCompositeKey("AssetRule", []string{})
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
			buffer.WriteString(",")
		}

		objectType, compositeKeyParts, err := APIstub.SplitCompositeKey(queryResponse.Key)
		if err != nil {
			return shim.Error(err.Error())
		}

		fmt.Printf("Found Asset Rule index:%s id:%s name:%s\n", objectType, compositeKeyParts[0], compositeKeyParts[1])

		buffer.WriteString(fmt.Sprintf("{\"Key\":\"%s\",", compositeKeyParts[0]))
		buffer.WriteString("\"Record\":{")
		buffer.WriteString(fmt.Sprintf("\"Name\":\"%s\"", compositeKeyParts[1]))
		buffer.WriteString("}}")

		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}
