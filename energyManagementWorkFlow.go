package main

import (
	"bytes"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/rs/xid"
)

var logger = shim.NewLogger("energyManagement-chaincode-log")

//var logger = shim.NewLogger("dcot-chaincode")

// DcotWorkflowChaincode implementation
type EnergyManagementWorkFlow.go struct {
	testMode bool
}
func (t *EnergyManagementWorkFlow) Init(stub shim.ChaincodeStubInterface) pb.Response {

	logger.Info("Chaincode Interface - Init()\n")
	logger.SetLevel(shim.LogDebug)
	_, args := stub.GetFunctionAndParameters()
	//var err error

	// Upgrade Mode 1: leave ledger state as it was
	if len(args) == 0 {
		//logger.Info("Args correctly!!!")
		return shim.Success(nil)
	}

	return shim.Success(nil)
}

func (t *EnergyManagementWorkFlow) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	var creatorOrg, creatorCertIssuer string
	//var attrValue string
	var err error
	var isEnabled bool
	var callerRole string

	logger.Debug("Chaincode Interface - Invoke()\n")

	if !t.testMode {
		creatorOrg, creatorCertIssuer, err = getTxCreatorInfo(stub)
		if err != nil {
			logger.Error("Error extracting creator identity info: \n", err.Error())
			return shim.Error(err.Error())
		}
		logger.Info("EnergyManagementWorkFlow Invoke by '', ''\n", creatorOrg, creatorCertIssuer)
		callerRole, _, err = getTxCreatorInfo(stub)
		if err != nil {
			return shim.Error(err.Error())
		}

		isEnabled, _, err = isInvokerOperator(stub, callerRole)
		if err != nil {
			logger.Error("Error getting attribute info: \n", err.Error())
			return shim.Error(err.Error())
		} 
	}

	function, args := stub.GetFunctionAndParameters()

	if (!isEnabled){
		return shim.Error("Permission denied to access the chaincode!!!\n")
	}

	if function == "putData" {
		return t.putData(stub, args)
	} else if function == "getData" {
		return t.updateAnalyticsInstances(stub, args)
	}
	return shim.Error("Invalid invoke function name\n")
}


func (t *EnergyManagementWorkFlow) getData(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("EnergyManagementWorkFlow:      getData()\n")
	if len(args) != 1 {
		logger.Error("getData() ERROR: need exactly one arguments\n")
		return shim.Error("getData() ERROR: need exactly two arguments")
	}

	energyManagementKey, err := stub.CreateCompositeKey("POLIMI-Energy-Management", args[0])
	if err != nil {
		logger.Error("CreateCompositeKey() ERROR: " err.Error())
		return shim.Error("CreateCompositeKey() ERROR: " err.Error())
	}
	dataBytes, err1 := stub.GetState(energyManagementKey)
	if err1 != nil {
		logger.Error("GetState() ERROR: " err1.Error())
		return shim.Error("GetState() ERROR: " err1.Error())
	}
	dataString := string(dataBytes[:])
	logger.Info("getData(): " +dataString)
	err2 := stub.SetEvent("getData-Event", dataBytes)
	if err2 != nil {
		logger.Error("SetEvent() ERROR: " err2.Error())
	}
	return shim.Success(dataBytes)


}
func (t *EnergyManagementWorkFlow) putData(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	logger.Info("EnergyManagementWorkFlow:      putData()\n")
	if len(args) != 2 {
		logger.Error("putData() ERROR: need exactly two arguments \n")
		return shim.Error("putData() ERROR: need exactly two arguments!! [KEY + PAYLOAD]")
	}
	energyManagementKey, err := stub.CreateCompositeKey("POLIMI-Energy-Management", args[0])
	if err != nil {
		logger.Error("CreateCompositeKey() ERROR: " err.Error())
		return shim.Error("CreateCompositeKey() ERROR: " err.Error())
	}

	err1 := stub.PutState(energyManagementKey,  []byte(args[1]))
	if err1 != nil {
		logger.Error("PutState() ERROR: " err1.Error())
		return shim.Error("PutState() ERROR: " err1.Error())
	}
	err2 := stub.SetEvent("putData-Event", []byte(args[1]))
	if err2 != nil {
		logger.Error("SetEvent() ERROR: " err2.Error())
	}
	return shim.Success(nil)

}

func main() {
	twc := new(EnergyManagementWorkFlow)
	twc.testMode = true
	err := shim.Start(twc)
	if err != nil {
		logger.Error("Error starting Energy-Management-chaincode: ", err)
	}
}