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
	//Asset
	case "createAsset":
		return CreateAsset(APIstub, args)
	// User Asset
	case "createUserAsset":
		return CreateUserAsset(APIstub, args)
	case "updateUserAsset":
		return UpdateUserAsset(APIstub, args)
	case "deleteUserAsset":
		return DeleteUserAsset(APIstub, args)
	case "getUserAsset":
		return GetUserAsset(APIstub, args)
	// User
	case "createUser":
		return CreateUser(APIstub, args)
	case "updateUser":
		return UpdateUser(APIstub, args)
	case "listUserByRange":
		return ListUserByRange(APIstub, args)
	case "listUserByQuery":
		return ListUserByQuery(APIstub, args)
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
