package main

import (
	"encoding/json"
	"fmt"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// {"Args": ["createAssetTransaction", "{\"Id\":1, \"SourceUserKey\":\"User1\", \"SourceAssetKey\":\"Asset1\", \"DestinationUserKey\":\"User2\", \"DestinationAssetKey\":\"Asset2\", \"Amount\":50}"]}
type createAssetTransactionParams struct {
	Id    			        int64
	SourceUserKey 			string
	SourceAssetKey 			string
	DestinationUserKey  string
	DestinationAssetKey string
	Amount 			        float32
}

func CreateAssetTransaction(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := createAssetTransactionParams{}

	// parse params string to struct and throw an error if found
	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	sourceUser := models.User{}
	destinationUser := models.User{}

	// check if SourceUser exist
	sourceUserAsBytes, _ := APIstub.GetState(params.SourceUserKey)
	if sourceUserAsBytes == nil {
		return shim.Error(fmt.Sprintf("Source User with key %s not found", params.SourceUserKey))
	}
	if err := json.Unmarshal(sourceUserAsBytes, &sourceUser); err != nil {
		return shim.Error(err.Error())
	}

	// check if DestinationUser exist
	destinationUserAsBytes, _ := APIstub.GetState(params.DestinationUserKey)
	if destinationUserAsBytes == nil {
		return shim.Error(fmt.Sprintf("Destination User with key %s not found", params.DestinationUserKey))
	}
	if err := json.Unmarshal(destinationUserAsBytes, &destinationUser); err != nil {
		return shim.Error(err.Error())
	}

	//  search for asset Rule
	assetRuleAsBytes, err := PrivateGetAssetRuleByAssets(APIstub, params.SourceAssetKey, params.DestinationAssetKey)
	assetRule := models.AssetRule{}
	if err != "" {
		return shim.Error(err)
	}
	if assetRuleAsBytes == nil {
		return shim.Error(fmt.Sprintf("Asset Rule with source key=%s and destination key=%s doesn't exist", params.SourceAssetKey, params.DestinationAssetKey))
	}
	if err := json.Unmarshal(assetRuleAsBytes, &assetRule); err != nil {
		return shim.Error(err.Error())
	}

	// search for UserAsset of the sender
	sourceUserAsset, err := PrivateGetOrNewUserAsset(APIstub, params.SourceUserKey, params.SourceAssetKey)
	if err != "" {
		return shim.Error(err)
	}

	// search for the UserAsset of the recipient
	destinationUserAsset, err := PrivateGetOrNewUserAsset(APIstub, params.DestinationUserKey, params.DestinationAssetKey)
	if err != "" {
		return shim.Error(err)
	}

	// apply the add and subtract on each userAssets
	totalAmountDeduct := params.Amount + assetRule.Fee
	totalAmountIncrement := params.Amount * assetRule.Factor

	// deduct totalAmountDeduct to SourceUserAsset
	sourceUserAsset.Balance = sourceUserAsset.Balance - totalAmountDeduct
	sourceUserAssetAsBytes, _ := json.Marshal(sourceUserAsset)
	if err := APIstub.PutState(sourceUserAsset.Key(), sourceUserAssetAsBytes); err != nil {
		return shim.Error(err.Error())
	}

	// increment amount on DestinationUserAsset
	destinationUserAsset.Balance = destinationUserAsset.Balance + totalAmountIncrement
	destinationUserAssetAsBytes, _ := json.Marshal(destinationUserAsset)
	if err := APIstub.PutState(destinationUserAsset.Key(), destinationUserAssetAsBytes); err != nil {
		return shim.Error(err.Error())
	}

	// create the AssetTransaction object
	assetTransaction := models.NewAssetTransaction()
	assetTransaction.Id = params.Id
	assetTransaction.SourceAsset = params.SourceAssetKey
	assetTransaction.DestinationAsset = params.DestinationAssetKey
	assetTransaction.Factor = assetRule.Factor
	assetTransaction.Fee = assetRule.Fee
	assetTransaction.Amount = params.Amount

	// check if AssetTransaction key already exist
	if assetTransactionAsBytes, _ := APIstub.GetState(assetTransaction.Key()); assetTransactionAsBytes != nil {
		return shim.Error(fmt.Sprintf("AssetTransaction with key %v already exist", params.Id))
	}

	// put new AssetTransaction to ledger
	assetTransactionAsBytes, _ := json.Marshal(assetTransaction)
	if err := APIstub.PutState(assetTransaction.Key(), assetTransactionAsBytes); err != nil {
		return shim.Error(err.Error())
	}

	// create Asset Transaction composite key
	indexName := "AssetTransaction"
	assetTransactionIndexKey, errCompositeKey := APIstub.CreateCompositeKey(indexName, []string{assetTransaction.Key()})
	if errCompositeKey != nil {
		return shim.Error(errCompositeKey.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	APIstub.PutState(assetTransactionIndexKey, value)

	return shim.Success(nil)
}
