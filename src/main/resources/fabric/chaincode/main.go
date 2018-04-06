package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("chaincode")

// Define the Smart Contract structure
type SmartContract struct{}

/*
 * The Init method is called when the Smart Contract "transaction" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "transaction"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()

	logger.Info("########### Transaction Invoke ###########")
	logger.Info("Function Name", function)
	logger.Info("Arguments", args)
	logger.Info("##########################################")

	// Route to the appropriate handler function to interact with the ledger appropriately
	switch function {
	//Exchange endpoint
	case "createExchangeEndpoint":
		return CreateExchangeEndpoint(APIstub, args)

	case "getExchangeEndpoint":
		return GetExchangeEndpoint(APIstub, args)

	case "listExchangeEndpoints":
		return ListExchangeEndpoints(APIstub, args)

	case "listExchangeEndpointByRange":
		return ListExchangeEndpointByRange(APIstub, args)

	case "updateExchangeEndpoint":
		return UpdateExchangeEndpoint(APIstub, args)

	case "deleteExchangeEndpoint":
		return DeleteExchangeEndpoint(APIstub, args)

	//Asset
	case "createAsset":
		return CreateAsset(APIstub, args)
	case "getAsset":
		return GetAsset(APIstub, args)
	case "deleteAsset":
		return DeleteAsset(APIstub, args)
	case "updateAsset":
		return UpdateAsset(APIstub, args)

	case "listAssets":
		return ListAssets(APIstub, args)

	case "listAssetsByRange":
		return ListAssetsByRange(APIstub, args)
	case "listAssetsByPublicUser":
		return ListAssetsByPublicUser(APIstub, args)
	// Beneficiary
	case "createBeneficiary":
		return CreateBeneficiary(APIstub, args)
	case "deleteBeneficiary":
		return DeleteBeneficiary(APIstub, args)
	case "getBeneficiary":
		return GetBeneficiary(APIstub, args)
	case "listBeneficiaries":
		return ListBeneficiaries(APIstub, args)
	case "listBeneficiariesByRange":
		return ListBeneficiariesByRange(APIstub, args)
	case "listBeneficiariesByQuery":
		return ListBeneficiariesByQuery(APIstub, args)
	// User Asset
	case "createUserAsset":
		return CreateUserAsset(APIstub, args)
	case "updateUserAsset":
		return UpdateUserAsset(APIstub, args)
	case "deleteUserAsset":
		return DeleteUserAsset(APIstub, args)
	case "getUserAsset":
		return GetUserAsset(APIstub, args)
	case "listUserAssets":
		return ListUserAssets(APIstub, args)
	case "listUserAssetsByRange":
		return ListUserAssetsByRange(APIstub, args)
	case "listUserAssetsByQuery":
		return ListUserAssetsByQuery(APIstub, args)
	// User
	case "createUser":
		return CreateUser(APIstub, args)
	case "updateUser":
		return UpdateUser(APIstub, args)
	case "getUser":
		return GetUser(APIstub, args)
	case "getUserByUid":
		return GetUserByUID(APIstub, args)
	case "deleteUser":
		return DeleteUser(APIstub, args)
	case "listUsers":
		return ListUsers(APIstub, args)
	case "listUserByRange":
		return ListUserByRange(APIstub, args)
	case "listUserByQuery":
		return ListUserByQuery(APIstub, args)
	// USER GROUP
	case "createUserGroup":
		return CreateUserGroup(APIstub, args)
	case "updateUserGroup":
		return UpdateUserGroup(APIstub, args)

	case "listUserGroups":
		return ListUserGroups(APIstub, args)

	case "getUserGroup":
		return GetUserGroup(APIstub, args)
	case "deleteUserGroup":
		return DeleteUserGroup(APIstub, args)
	// ISSUANCE RULE
	case "createIssuanceRule":
		return CreateIssuanceRule(APIstub, args)

	case "listIssuanceRulesByQuery":
		return ListIssuanceRulesByQuery(APIstub)

	case "listIssuanceRules":
		return ListIssuanceRules(APIstub, args)

	case "deleteIssuanceRule":
		return DeleteIssuanceRule(APIstub, args)

	case "getIssuanceRule":
		return GetIssuanceRule(APIstub, args)

	case "updateIssuanceRule":
		return UpdateIssuanceRule(APIstub, args)

	case "getIssuanceRuleByAssetAndExchangeEndpoint":
		return GetIssuanceRuleByAssetAndExchangeEndpoint(APIstub, args)

	// Asset Rule
	case "createAssetRule":
		return CreateAssetRule(APIstub, args)
	case "updateAssetRule":
		return UpdateAssetRule(APIstub, args)
	case "getAssetRule":
		return GetAssetRule(APIstub, args)
	case "listAssetRulesByQuery":
		return ListAssetRulesByQuery(APIstub)
	case "listAssetRulesByCompositeKey":
		return ListAssetRulesByCompositeKey(APIstub)
	case "getAssetRuleByAssets":
		return GetAssetRuleByAssets(APIstub, args)
	// Asset Transaction
	case "createAssetTransaction":
		return CreateAssetTransaction(APIstub, args)
	case "getAssetTransaction":
		return GetAssetTransaction(APIstub, args)
	case "listAssetTransactionsByQuery":
		return ListAssetTransactionsByQuery(APIstub)
	case "listAssetTransactionsByCompositeKey":
		return ListAssetTransactionsByCompositeKey(APIstub)
	// Issuance Transaction
	case "createIssuanceTransaction":
		return CreateIssuanceTransaction(APIstub, args)
	case "getIssuanceTransaction":
		return GetIssuanceTransaction(APIstub, args)
	case "listIssuanceTransactionsByQuery":
		return ListIssuanceTransactionsByQuery(APIstub)
	case "listIssuanceTransactionsByCompositeKey":
		return ListIssuanceTransactionsByCompositeKey(APIstub)
	// Transactions
	case "listTransactionByUser":
		return ListTransactionByUser(APIstub, args)

	// Retirement RULE
	case "createRetirementRule":
		return CreateRetirementRule(APIstub, args)

	case "listRetirementRules":
		return ListRetirementRules(APIstub, args)

	case "updateRetirementRule":
		return UpdateRetirementRule(APIstub, args)

	default:
		return shim.Error("Invalid Smart Contract function name.")
	}
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {
	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
