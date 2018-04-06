package main

import (
	"encoding/json"
	"fmt"
	"main/models"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// {"Args": ["createIssuanceTransaction", "{\"Id\":1, \"BrokerKey\":\"User3\", \"ExchangeEndpointKey\":\"ExchangeEndpoint1\", \"DestinationUserKey\":\"User1\", \"DestinationAssetKey\":\"Asset1\", \"Amount\":50}"]}
type createIssuanceTransactionParams struct {
	Id    			        int64
	BrokerKey      			string
	ExchangeEndpointKey string
	DestinationUserKey  string
	DestinationAssetKey string
	Amount 			        float32
}

func CreateIssuanceTransaction(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	params := createIssuanceTransactionParams{}

	// parse params string to struct and throw an error if found
	if err := json.Unmarshal([]byte(args[0]), &params); err != nil {
		return shim.Error(err.Error())
	}

	broker := models.User{}
	destinationUser := models.User{}

	// check if Broker exist
	brokerAsBytes, _ := APIstub.GetState(params.BrokerKey)
	if brokerAsBytes == nil {
		return shim.Error(fmt.Sprintf("Broker with key %s not found", params.BrokerKey))
	}
	if err := json.Unmarshal(brokerAsBytes, &broker); err != nil {
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
	issuanceRuleAsBytes, err := PrivateGetIssuanceRuleByAssetAndExchangeEndpoint(APIstub, params.ExchangeEndpointKey, params.DestinationAssetKey)
	issuanceRule := models.IssuanceRule{}
	if err != "" {
		return shim.Error(err)
	}
	if issuanceRuleAsBytes == nil {
		return shim.Error(fmt.Sprintf("Issuance Rule with endpoint %s and destination %s not found", params.ExchangeEndpointKey, params.DestinationAssetKey))
	}
	if err := json.Unmarshal(issuanceRuleAsBytes, &issuanceRule); err != nil {
		return shim.Error(err.Error())
	}

	// search for the UserAsset of the recipient
	destinationUserAsset, err := PrivateGetOrNewUserAsset(APIstub, params.DestinationUserKey, params.DestinationAssetKey)
	if err != "" {
		return shim.Error(err)
	}

	// compute total amount
	totalAmountIncrement := params.Amount * issuanceRule.Factor

	// increment amount on DestinationUserAsset
	destinationUserAsset.Balance = destinationUserAsset.Balance + totalAmountIncrement
	destinationUserAssetAsBytes, _ := json.Marshal(destinationUserAsset)
	if err := APIstub.PutState(destinationUserAsset.Key(), destinationUserAssetAsBytes); err != nil {
		return shim.Error(err.Error())
	}

	// create the IssuanceTransaction object
	issuanceTransaction := models.NewIssuanceTransaction()
	issuanceTransaction.Id = params.Id
	issuanceTransaction.Broker = params.BrokerKey
	issuanceTransaction.ExchangeEndpoint = params.ExchangeEndpointKey
	issuanceTransaction.DestinationAsset = params.DestinationAssetKey
	issuanceTransaction.Factor = issuanceTransaction.Factor
	issuanceTransaction.Fee = issuanceTransaction.Fee
	issuanceTransaction.Amount = params.Amount

	// check if IssuanceTransaction key already exist
	if issuanceTransactionAsBytes, _ := APIstub.GetState(issuanceTransaction.Key()); issuanceTransactionAsBytes != nil {
		return shim.Error(fmt.Sprintf("IssuanceTransaction with key %v already exist", params.Id))
	}

	// put new IssuanceTransaction to ledger
	issuanceTransactionAsBytes, _ := json.Marshal(issuanceTransaction)
	if err := APIstub.PutState(issuanceTransaction.Key(), issuanceTransactionAsBytes); err != nil {
		return shim.Error(err.Error())
	}

	// create Issuance Transaction composite key
	indexName := "IssuanceTransaction"
	issuanceTransactionIndexKey, errCompositeKey := APIstub.CreateCompositeKey(indexName, []string{issuanceTransaction.Key()})
	if errCompositeKey != nil {
		return shim.Error(errCompositeKey.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	APIstub.PutState(issuanceTransactionIndexKey, value)

	return shim.Success(nil)
}
