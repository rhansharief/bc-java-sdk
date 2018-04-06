package main

import (
	"encoding/json"
	"fmt"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// {"Args": ["createAssetRule", "{\"Id\":1, \"Name\":\"AssetRule\", \"SourceKey\":\"Asset1\", \"DestinationKey\":\"Asset2\", \"Factor\":2, \"Fee\":0.50}"]}
type createAssetRuleParams struct {
	Id    			    int64
	Name 				    string
	SourceKey 			string
	DestinationKey  string
	Factor 			    float32
	Fee 				    float32
}

func CreateAssetRule(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := createAssetRuleParams{}

	// parse params string to struct and throw an error if found
	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	// check Source Asset if it exist
	sourceAssetAsBytes, _ := APIstub.GetState(params.SourceKey)
	if sourceAssetAsBytes == nil {
		return shim.Error(fmt.Sprintf("Asset with key %s not found", params.SourceKey))
	}

	// check Destination Asset if it exist
	destinationAssetAsBytes, _ := APIstub.GetState(params.DestinationKey)
	if destinationAssetAsBytes == nil {
		return shim.Error(fmt.Sprintf("Asset with key %s not found", params.DestinationKey))
	}

	// check if AssetRule already exist based on Source and Destination Ids
	assetRuleAsBytes, err := PrivateGetAssetRuleByAssets(APIstub, params.SourceKey, params.DestinationKey)
	if err != "" {
		return shim.Error(err)
	}
	if assetRuleAsBytes != nil {
		return shim.Error(fmt.Sprintf("Asset Rule with source key=%s and destination key=%s already exist", params.SourceKey, params.DestinationKey))
	}

	assetRule := models.NewAssetRule()
	assetRule.Id = params.Id
	assetRule.Name = params.Name
	assetRule.Source = params.SourceKey
	assetRule.Destination = params.DestinationKey
	assetRule.Factor = params.Factor
	assetRule.Fee = params.Fee

	// check if AssetRule already exist based on Id
	if assetRuleAsBytes, _ := APIstub.GetState(assetRule.Key()); assetRuleAsBytes != nil {
		return shim.Error(fmt.Sprintf("AssetRule with key %d already exist", params.Id))
	}

	// put new AssetRule to chain
	assetRuleAsBytes, _ = json.Marshal(assetRule)
	if err := APIstub.PutState(assetRule.Key(), assetRuleAsBytes); err != nil {
		return shim.Error(err.Error())
	}

	// create composite key for AssetRule
	indexName := "AssetRule"
	assetRuleIndexKey, errCompositeKey := APIstub.CreateCompositeKey(indexName, []string{fmt.Sprintf("%d", assetRule.Id), assetRule.Name})
	if errCompositeKey != nil {
		return shim.Error(errCompositeKey.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	APIstub.PutState(assetRuleIndexKey, value)

	return shim.Success(nil)
}
