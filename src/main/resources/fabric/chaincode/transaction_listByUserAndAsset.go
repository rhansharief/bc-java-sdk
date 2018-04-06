package main

import (
	"fmt"
	"encoding/json"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// TODO replace AssetId with AssetKey
// {"Args": ["listTransactionByUser", "{\"UserKey\":\"User1\", \"AssetId\":1}"]}
type listTransactionByUserParams struct {
	UserKey string
	AssetId int64
}

func ListTransactionByUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	// TODO replace AssetTransaction with transaction
	transactions := []models.AssetTransaction{}
	params := listTransactionByUserParams{}

	// parse params string to struct and throw an error if found
	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	 // get User Assets
	userAssets, err := PrivateListAssetsByPublicUser(APIstub, params.UserKey)
	if err != "" {
		return shim.Error(err)
	}

	// loop through all user.Asset
	for _, assetEl := range userAssets {
		// filter assets
		if params.AssetId == assetEl.Id {
			assetKey := fmt.Sprintf("Asset%d", assetEl.Id)

			// search for AssetTransactions, IssuanceTransactions based on each assets using a query
			query := fmt.Sprintf("" +
				"{\"selector\":{" +
					"\"$or\": [" +
						"{ \"DocType\": \"ASSET_TRANSACTION\" }," +
						"{ \"DocType\": \"ISSUANCE_TRANSACTION\" }" +
					"], " +
					"\"$or\": [" +
						"{\"SourceAsset\":\"%[1]s\"}," +
						"{\"DestinationAssetId\": \"%[1]s\"}" +
					"]" +
				"}}", assetKey)

			resultsIterator, err := APIstub.GetQueryResult(query)

			if err != nil {
				return shim.Error(err.Error())
			}
			defer resultsIterator.Close()

			for resultsIterator.HasNext() {
				queryResponse, err := resultsIterator.Next()
				if err != nil {
					return shim.Error(err.Error())
				}

				assetTransaction := models.AssetTransaction{}
				if err := json.Unmarshal(queryResponse.Value, &assetTransaction); err != nil {
					return shim.Error(err.Error())
				}

				// push found assetTransaction to total assetTransaction list
				transactions = append(transactions, assetTransaction)
			}
		}
	}

	transactionsAsBytes, _ := json.Marshal(transactions)

	return shim.Success(transactionsAsBytes)
}
