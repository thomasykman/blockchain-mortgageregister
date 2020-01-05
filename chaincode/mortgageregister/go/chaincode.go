// Inspired by fabric-samples, in particular fabcar and marbles_chaincode_private.go
// Place this file in: fabric-sample/chaincode/marbles02_private/go

package main

/* Imports
 * 2 utility libraries for formatting, reading and writing JSON
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ===================================================================================
// Data structures (struct)
// ===================================================================================

// struct representing a mortgage register chaincode "SimpleChaincode"
type SimpleChaincode struct {
}

// struct representation of a mortgage loan
type loan struct {
	LoanUID   string `json:"loanUID"`   // unique ID of the loan, e.g. loan1, loan2, loan3
	Issuer    string `json:"issuer"`    // bank issuing the loan
	Buyer     string `json:"buyer"`     // hash of the national registry number of the person borrowing the loan (buyer of the house)
	Notary    string `json:"notary"`    // notary of the buyer
	Status    string `json:"status"`    // status of the loan
	StartDate string `json:"startDate"` // start date of the loan, e.g. 31/01/2020
	EndDate   string `json:"endDate"`   // end date of the loan, e.g. 31/01/2040
}

// struct representation of the private info of a loan, visible only by the bank and notary involved in the loan
type loanPrivateInfo struct {
	LoanUID      string  `json:"loanUID"`      // unique ID of the loan, e.g. loan1, loan2, loan3
	LoanValue    int     `json:"loanValue"`    // mortgage loan value (eg 300000)
	Currency     string  `json:"currency"`     // mortgage loan currency (eg EUR)
	InterestRate float64 `json:"interestRate"` // mortgage loan interest rate (eg 0.02)
}

// ===================================================================================
// Main - creates a mortgageRegister chaincode ("SimpleChaincode")
// ===================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting mortgageRegister chaincode: %s", err)
	}
}

// ===================================================================================
// Init - initializes the mortgageRegister chaincode
// ===================================================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// ===================================================================================
// Invoke - Entry point for invocations
// ===================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	// GetFunctionAndParameters returns the first argument as the function
	// name and the rest of the arguments as parameters in a string array.
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	switch function {
	case "issueLoan":
		return t.issueLoan(stub, args)
	case "readLoan":
		return t.readLoan(stub, args)
	case "changeLoanStatus":
		return t.changeLoanStatus(stub, args)
	default:
		//error
		fmt.Println("invoke did not find func: " + function)
		return shim.Error("Received unknown function invocation")
	}
}

// ===================================================================================
// getOrgNameFromMSPID - Extract the organisation name from the organisation MSP identity
// e.g. returns "Org1" from "Org1MSP"
// ===================================================================================
func getOrgNameFromMSPID(orgMSPID string) string {
	temp := strings.Split(orgMSPID, "MSP")
	orgName := temp[0]
	return orgName
}

// ===================================================================================
// issueLoan - Create a new mortgage loan asset, store it into chaincode state
// transient data: key = "loan"; value = LoanUID, Buyer, Notary, StartDate, EndDate, LoanValue, Currency, InterestRate
// ===================================================================================

func (t *SimpleChaincode) issueLoan(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("*** start issueLoan")

	// *** 0. declare variable "err", define new structure "loanTransientInput"
	var err error

	type loanTransientInput struct {
		LoanUID      string  `json:"loanUID"`
		Buyer        string  `json:"buyer"`
		Notary       string  `json:"notary"`
		StartDate    string  `json:"startDate"`
		EndDate      string  `json:"endDate"`
		LoanValue    int     `json:"loanValue"`
		Currency     string  `json:"currency"`
		InterestRate float64 `json:"interestRate"`
	}

	// *** 1. check : #arguments = 0 ?
	// All information to be stored on the private data collection is passed via transient data,
	// so there should be no argument.
	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Private loan data must be passed in transient map.")
	}

	// *** 2. get the name of the organisation of the transaction creator
	// (user submitting the current transaction "issueLoan")

	// GetMSPID returns the ID of the MSP associated with the identity that submitted the transaction
	creatorOrgMSPID, err := cid.GetMSPID(stub)
	if err != nil {
		return shim.Error("Impossible to get the transaction creator's organisation MSP identity")
	}
	creatorOrgName := getOrgNameFromMSPID(creatorOrgMSPID)

	// *** 3. transMap = get transient data
	// Get transient data passed when the chaincode function "issueLoan" is invoked by the CLI or client app
	transMap, err := stub.GetTransient()

	// Check if GetTransient doesn't give any error
	if err != nil {
		return shim.Error("Error getting transient data: " + err.Error())
	}

	// Check : key in transient map = "loan"
	if _, ok := transMap["loan"]; !ok {
		return shim.Error("loan must be a key in the transient map")
	}

	// Check: value for key "loan" in transient map is not empty
	if len(transMap["loan"]) == 0 {
		return shim.Error("loan value in the transient map must be a non-empty JSON string")
	}

	// *** 4. get loanInput from the transient map: loanInput ≃ transMap["loan"]
	var loanInput loanTransientInput

	// Parse transient data from transMap (JSON format) and store the result in loanInput
	err = json.Unmarshal(transMap["loan"], &loanInput)
	if err != nil {
		return shim.Error("Failed to decode JSON of: " + string(transMap["loan"]))
	}

	// Check: all fields have been filled correctly (e.g. not empty, loan value > 0, loan status right)
	if len(loanInput.LoanUID) == 0 {
		return shim.Error("The loan UID is a required field")
	}
	if len(loanInput.Buyer) == 0 {
		return shim.Error("The buyer name is a required field")
	}
	if len(loanInput.Notary) == 0 {
		return shim.Error("The notary name is a required field")
	}
	if len(loanInput.StartDate) == 0 {
		return shim.Error("The start date is a required field")
	}
	if len(loanInput.EndDate) == 0 {
		return shim.Error("The end date is a required field")
	}
	if loanInput.LoanValue == 0 {
		return shim.Error("The loan value is a required field")
	}
	if len(loanInput.Currency) == 0 {
		return shim.Error("The currency is a required field")
	}
	if loanInput.InterestRate == 0 {
		return shim.Error("The interest rate is a required field")
	}
	if loanInput.LoanValue <= 0 {
		return shim.Error("The loan value must be positive")
	}

	// *** 5. check : loan does not exist yet on the private data collection ?
	// Check in the private data collection "collectionLoans" if there is already a loan with the same loanUID as in loanInput
	existingLoanAsBytes, err := stub.GetPrivateData("collectionLoans", loanInput.LoanUID)
	if err != nil {
		return shim.Error("Failed to get loan: " + err.Error())
	} else if existingLoanAsBytes != nil {
		return shim.Error("This loan already exists: " + loanInput.LoanUID)
	}

	// *** 6. create new loan and save it on the private data collection "collectionLoans"
	loan := &loan{
		LoanUID:   loanInput.LoanUID,
		Issuer:    creatorOrgName,
		Buyer:     loanInput.Buyer,
		Notary:    loanInput.Notary,
		Status:    "issued",
		StartDate: loanInput.StartDate,
		EndDate:   loanInput.EndDate,
	}

	loanAsBytes, err := json.Marshal(loan)
	if err != nil {
		return shim.Error(err.Error())
	}

	// PutPrivateData(collection, key, value) puts the `key` and `value` into the private data collection
	err = stub.PutPrivateData("collectionLoans", loanInput.LoanUID, loanAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// *** 7. create new loan private data and store it on the private data collection "collectionLoanPrivateInfo"
	loanPrivateInfo := &loanPrivateInfo{
		LoanUID:      loanInput.LoanUID,
		LoanValue:    loanInput.LoanValue,
		Currency:     loanInput.Currency,
		InterestRate: loanInput.InterestRate,
	}

	loanPrivateInfoAsBytes, err := json.Marshal(loanPrivateInfo)
	if err != nil {
		return shim.Error(err.Error())
	}

	// PutPrivateData(collection, key, value) puts the `key` and `value` into the private data collection
	err = stub.PutPrivateData("collectionLoanPrivateInfo", loanInput.LoanUID, loanPrivateInfoAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// *** 8. return success
	fmt.Println("*** end issueLoan")
	return shim.Success(nil)
}

// ===================================================================================
// readLoan - read a loan from chaincode state for a specific private data collection
// arguments: loanUID, name of the private data collection
// ===================================================================================
func (t *SimpleChaincode) readLoan(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var loanUID, collection, jsonResp string
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting loan UID and private data collection")
	}

	// UID of the loan for which we look for information
	loanUID = args[0]
	// name of the private data collection
	collection = args[1]

	// get loan information from the private data collection "collection"
	valAsBytes, err := stub.GetPrivateData(collection, loanUID)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for loan " + loanUID + "\"}"
		return shim.Error(jsonResp)
	} else if valAsBytes == nil {
		jsonResp = "{\"Error\":\"Loan does not exist: " + loanUID + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsBytes)
}

// ===================================================================================
// changeLoanStatus - change the status of a loan to a new status
// transient data: key = "loan_status" ; value = loanUID, new status
// ===================================================================================
func (t *SimpleChaincode) changeLoanStatus(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("*** start changeLoanStatus")

	// *** 0. get the name of the organisation of the transaction creator
	// (user submitting the current transaction "changeLoanStatus")

	// GetMSPID returns the ID of the MSP associated with the identity that submitted the transaction
	creatorOrgMSPID, err := cid.GetMSPID(stub)
	if err != nil {
		return shim.Error("Impossible to get the transaction creator's organisation MSP identity")
	}
	creatorOrgName := getOrgNameFromMSPID(creatorOrgMSPID)

	// *** 1. define new structure + basic checks
	type loanStatusTransientInput struct {
		LoanUID string `json:"loanUID"`
		Status  string `json:"status"`
	}

	// check: no argument
	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Private loan data must be passed in transient map.")
	}

	// *** 2. transMap : get transient data
	transMap, err := stub.GetTransient()

	// check: no error
	if err != nil {
		return shim.Error("Error getting transient: " + err.Error())
	}

	// check: "loan_status" = key in transient map
	if _, ok := transMap["loan_status"]; !ok {
		return shim.Error("loan_status must be a key in the transient map")
	}

	// Check: value for key "loan" in transient map is not empty
	if len(transMap["loan_status"]) == 0 {
		return shim.Error("loan_status value in the transient map must be a non-empty JSON string")
	}

	// *** 3. loanStatusInput : unmarshal transient data : loanUID, new status
	var loanStatusInput loanStatusTransientInput
	err = json.Unmarshal(transMap["loan_status"], &loanStatusInput)

	// check: no error
	if err != nil {
		return shim.Error("Failed to decode JSON of: " + string(transMap["loan_status"]))
	}

	// check: loanUID is not empty
	if len(loanStatusInput.LoanUID) == 0 {
		return shim.Error("loan UID field must be a non-empty string")
	}

	// check: loan status is not empty
	if len(loanStatusInput.Status) == 0 {
		return shim.Error("loan status field must be a non-empty string")
	}

	// *** 4. loanAsBytes : get the loan "loanUID" from private data collection
	loanAsBytes, err := stub.GetPrivateData("collectionLoans", loanStatusInput.LoanUID)

	// check: no error, loan exists
	if err != nil {
		return shim.Error("Failed to get loan:" + err.Error())
	} else if loanAsBytes == nil {
		return shim.Error("Loan does not exist: " + loanStatusInput.LoanUID)
	}

	// *** 5. loanToChange : unmarshal loanAsBytes
	loanToChange := loan{}
	err = json.Unmarshal(loanAsBytes, &loanToChange)

	// check: no error
	if err != nil {
		return shim.Error(err.Error())
	}

	// checks:
	// 1. new status (loanStatusInput.Status) can only be "active", "inactive" or "cancelled"
	// 2. status can only change from "issued" to "active", "active" to "inactive", or "inactive" to "cancelled"
	// 3. only Notary/Bank can change status
	switch loanStatusInput.Status {
	case "active":
		if loanToChange.Status != "issued" {
			return shim.Error("Only 'issued' loan can be changed to 'active'. New status is:" + loanStatusInput.Status + ", but current status is: " + loanToChange.Status)
		} else if creatorOrgName != loanToChange.Notary {
			return shim.Error("Only notary can change the status to 'active'. Transaction creator is: " + creatorOrgName + ", but notary is: " + loanToChange.Notary)
		}
	case "inactive":
		if loanToChange.Status != "active" {
			return shim.Error("Only 'active' loan can be changed to 'inactive'. New status is:" + loanStatusInput.Status + ", but current status is: " + loanToChange.Status)
		} else if (creatorOrgName != loanToChange.Notary) && (creatorOrgName != loanToChange.Issuer) {
			return shim.Error("Only notary or bank can change the status to 'active'. Transaction creator is: " + creatorOrgName + ", but notary is: " + loanToChange.Notary + " and bank is: " + loanToChange.Issuer)
		}
	case "cancelled":
		if loanToChange.Status != "inactive" {
			return shim.Error("Only 'inactive' loan can be changed to 'cancelled'. New status is:" + loanStatusInput.Status + ", but current status is: " + loanToChange.Status)
		} else if creatorOrgName != loanToChange.Notary {
			return shim.Error("Only notary can change the status to 'active'. Transaction creator is: " + creatorOrgName + ", but notary is: " + loanToChange.Notary)
		}
	default:
		return shim.Error("New status must be 'active', 'inactive' or 'cancelled'. New status proposed is: " + loanStatusInput.Status)
	}

	// *** 6. update private data collection with the new status of the loan "loanUID"
	// change status of the loan
	loanToChange.Status = loanStatusInput.Status

	// encode loanToChange in JSON
	loanToChangeJSON, _ := json.Marshal(loanToChange)

	// update the loanToChange in the private data collection
	err = stub.PutPrivateData("collectionLoans", loanToChange.LoanUID, loanToChangeJSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	// *** 7. return success
	fmt.Println("*** end changeLoanStatus")
	return shim.Success(nil)
}
